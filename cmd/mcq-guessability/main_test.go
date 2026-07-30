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
// construction and the arithmetic behind it. Running it through the full
// CLI pipeline (not just EvaluateGate directly) proves main.go's own exit
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
	if !strings.Contains(failBuf.String(), "Gate: FAIL") {
		t.Errorf("output missing \"Gate: FAIL\":\n%s", failBuf.String())
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
