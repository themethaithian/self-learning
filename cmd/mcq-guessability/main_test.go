package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fixtureCheck struct {
	Q              string   `json:"q"`
	ExpectedAnswer string   `json:"expected_answer"`
	Type           string   `json:"type"`
	Options        []string `json:"options,omitempty"`
}

type fixtureLesson struct {
	Topic        string         `json:"topic"`
	RecallChecks []fixtureCheck `json:"recall_checks"`
}

func writeGateStraddlingFixture(t *testing.T, dir string) {
	t.Helper()

	options := []string{"short-a", "short-b", "short-c", "the much longer fourth option here"}
	pick := func(idx int) fixtureCheck {
		return fixtureCheck{Q: "q", ExpectedAnswer: options[idx], Type: "mcq", Options: options}
	}

	var checks []fixtureCheck
	for i := 0; i < 30; i++ {
		checks = append(checks, pick(3))
	}
	for i := 0; i < 24; i++ {
		checks = append(checks, pick(0))
	}
	for i := 0; i < 23; i++ {
		checks = append(checks, pick(1))
	}
	for i := 0; i < 23; i++ {
		checks = append(checks, pick(2))
	}

	lesson := fixtureLesson{Topic: "gate-test", RecallChecks: checks}
	raw, err := json.Marshal(lesson)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	topicDir := filepath.Join(dir, "gate-test")
	if err := os.MkdirAll(topicDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(topicDir, "concept.json"), raw, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

// This fixture's length-heuristic excess is exactly 0.20 (30% hit rate vs
// 25% baseline) — see internal/mcqguess's gate_test.go for the identical
// construction and the exact arithmetic (including why the position-3
// heuristic moves together with the length heuristic here, not
// independently — they are the same guess on this fixture by construction,
// so both cross the threshold at once). Running it through the full CLI
// pipeline (not just EvaluateTrackGates directly) proves main.go's own exit
// code wiring, not just the internal package's gate logic in isolation.
func TestRun_GateStraddlingCorpus_ExitCodeFlips(t *testing.T) {
	dir := t.TempDir()
	writeGateStraddlingFixture(t, dir)

	var passBuf bytes.Buffer
	if code := run(dir, 0.21, &passBuf); code != 0 {
		t.Errorf("run() at max-excess=0.21 = %d, want 0", code)
	}
	if !strings.Contains(passBuf.String(), "Gate: PASS") {
		t.Errorf("output missing \"Gate: PASS\":\n%s", passBuf.String())
	}

	var failBuf bytes.Buffer
	if code := run(dir, 0.19, &failBuf); code != 1 {
		t.Errorf("run() at max-excess=0.19 = %d, want 1", code)
	}
	out := failBuf.String()
	if !strings.Contains(out, "Gate: FAIL") {
		t.Errorf("output missing \"Gate: FAIL\":\n%s", out)
	}
	if !strings.Contains(out, "length heuristic") || !strings.Contains(out, "index 3") {
		t.Errorf("output missing both expected violations (length heuristic AND position index 3):\n%s", out)
	}
}

func TestRun_PerTrackGateCatchesWhatPooledWouldMiss(t *testing.T) {
	// End-to-end version of the internal package's dilution test: a large
	// fair track (0% excess) plus a small bad track (60% excess) pools to
	// 15% overall — comfortably under a 20% threshold the bad track alone
	// fails badly. If run() gated on the pooled figure instead of per
	// track, this would exit 0.
	dir := t.TempDir()

	fairOptions := []string{"opt-A", "opt-B", "opt-C", "opt-D"}
	var fairChecks []fixtureCheck
	for i := 0; i < 300; i++ {
		fairChecks = append(fairChecks, fixtureCheck{Q: "q", ExpectedAnswer: fairOptions[i%4], Type: "mcq", Options: fairOptions})
	}
	writeLessonFixture(t, dir, "fair-track", fairChecks)

	badOptions := []string{"short-a", "short-b", "short-c", "the much longer fourth option here"}
	var badChecks []fixtureCheck
	pick := func(idx int) fixtureCheck {
		return fixtureCheck{Q: "q", ExpectedAnswer: badOptions[idx], Type: "mcq", Options: badOptions}
	}
	for i := 0; i < 40; i++ {
		badChecks = append(badChecks, pick(3))
	}
	for i := 0; i < 20; i++ {
		badChecks = append(badChecks, pick(0))
	}
	for i := 0; i < 20; i++ {
		badChecks = append(badChecks, pick(1))
	}
	for i := 0; i < 20; i++ {
		badChecks = append(badChecks, pick(2))
	}
	writeLessonFixture(t, dir, "bad-track", badChecks)

	var buf bytes.Buffer
	code := run(dir, 0.20, &buf)
	if code != 1 {
		t.Errorf("run() at max-excess=0.20 = %d, want 1 (bad-track's 60%% excess must fail even though pooling would hide it)\noutput:\n%s", code, buf.String())
	}
	if !strings.Contains(buf.String(), "[bad-track]") {
		t.Errorf("output missing a violation naming bad-track:\n%s", buf.String())
	}
}

func writeLessonFixture(t *testing.T, dir, topic string, checks []fixtureCheck) {
	t.Helper()
	lesson := fixtureLesson{Topic: topic, RecallChecks: checks}
	raw, err := json.Marshal(lesson)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	topicDir := filepath.Join(dir, topic)
	if err := os.MkdirAll(topicDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(topicDir, "concept.json"), raw, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func TestRun_GateDisabledByDefault(t *testing.T) {
	dir := t.TempDir()
	writeGateStraddlingFixture(t, dir)

	var buf bytes.Buffer
	if code := run(dir, -1, &buf); code != 0 {
		t.Errorf("run() at max-excess=-1 = %d, want 0 (gate disabled)", code)
	}
	if !strings.Contains(buf.String(), "Gate: disabled") {
		t.Errorf("output missing \"Gate: disabled\":\n%s", buf.String())
	}
}

func TestRun_NoMCQsFound(t *testing.T) {
	dir := t.TempDir()

	var buf bytes.Buffer
	if code := run(dir, -1, &buf); code != 1 {
		t.Errorf("run() on empty dir = %d, want 1", code)
	}
	if !strings.Contains(buf.String(), "no mcq recall_checks found") {
		t.Errorf("output missing the no-mcq message:\n%s", buf.String())
	}
}

func TestRun_MalformedCorpusAborts(t *testing.T) {
	dir := t.TempDir()
	topicDir := filepath.Join(dir, "broken-topic")
	if err := os.MkdirAll(topicDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	lesson := fixtureLesson{
		Topic: "broken-topic",
		RecallChecks: []fixtureCheck{
			{Q: "q", ExpectedAnswer: "not-an-option", Type: "mcq", Options: []string{"a", "b", "c"}},
		},
	}
	raw, err := json.Marshal(lesson)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if err := os.WriteFile(filepath.Join(topicDir, "concept.json"), raw, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var buf bytes.Buffer
	if code := run(dir, -1, &buf); code != 1 {
		t.Errorf("run() on malformed corpus = %d, want 1", code)
	}
	if !strings.Contains(buf.String(), "not among its") {
		t.Errorf("output missing the malformed-question error:\n%s", buf.String())
	}
}
