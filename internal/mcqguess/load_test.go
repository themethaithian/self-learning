package mcqguess

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func writeFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

func TestLoadQuestions_FiltersShortAnswer(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, filepath.Join(dir, "topic-a", "concept.json"), `{
		"topic": "topic-a",
		"recall_checks": [
			{"q": "sa", "expected_answer": "x", "type": "short_answer"},
			{"q": "mcq1", "expected_answer": "b", "type": "mcq", "options": ["a", "b", "c"]},
			{"q": "sa2", "expected_answer": "y", "type": "short_answer"}
		]
	}`)

	questions, err := LoadQuestions(dir)
	if err != nil {
		t.Fatalf("LoadQuestions() error = %v", err)
	}
	if len(questions) != 1 {
		t.Fatalf("LoadQuestions() returned %d questions, want 1 (short_answer must be excluded)", len(questions))
	}
	if questions[0].ExpectedAnswer != "b" {
		t.Errorf("questions[0].ExpectedAnswer = %q, want %q", questions[0].ExpectedAnswer, "b")
	}
	if questions[0].Track != "topic-a" {
		t.Errorf("questions[0].Track = %q, want %q", questions[0].Track, "topic-a")
	}
}

func TestLoadQuestions_WalksNestedTopicDirectories(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, filepath.Join(dir, "topic-a", "c1.json"), `{
		"topic": "topic-a",
		"recall_checks": [{"q": "q1", "expected_answer": "a", "type": "mcq", "options": ["a", "b"]}]
	}`)
	writeFixture(t, filepath.Join(dir, "topic-b", "c2.json"), `{
		"topic": "topic-b",
		"recall_checks": [{"q": "q2", "expected_answer": "c", "type": "mcq", "options": ["c", "d"]}]
	}`)

	questions, err := LoadQuestions(dir)
	if err != nil {
		t.Fatalf("LoadQuestions() error = %v", err)
	}
	if len(questions) != 2 {
		t.Fatalf("LoadQuestions() returned %d questions, want 2", len(questions))
	}

	tracks := []string{questions[0].Track, questions[1].Track}
	sort.Strings(tracks)
	if tracks[0] != "topic-a" || tracks[1] != "topic-b" {
		t.Errorf("tracks = %v, want [topic-a topic-b]", tracks)
	}
}

func TestLoadQuestions_IgnoresNonJSONFiles(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, filepath.Join(dir, "topic-a", "concept.json"), `{
		"topic": "topic-a",
		"recall_checks": [{"q": "q1", "expected_answer": "a", "type": "mcq", "options": ["a", "b"]}]
	}`)
	writeFixture(t, filepath.Join(dir, "topic-a", "README.md"), "not json")

	questions, err := LoadQuestions(dir)
	if err != nil {
		t.Fatalf("LoadQuestions() error = %v", err)
	}
	if len(questions) != 1 {
		t.Errorf("LoadQuestions() returned %d questions, want 1 (non-.json file must be ignored)", len(questions))
	}
}

func TestLoadQuestions_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	questions, err := LoadQuestions(dir)
	if err != nil {
		t.Fatalf("LoadQuestions() error = %v", err)
	}
	if len(questions) != 0 {
		t.Errorf("LoadQuestions() returned %d questions, want 0", len(questions))
	}
}
