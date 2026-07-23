package infra

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

func TestLoadLesson_Valid(t *testing.T) {
	topicSlug, lesson, err := LoadLesson(filepath.Join("testdata", "domain-driven-design", "ubiquitous-language.json"))
	if err != nil {
		t.Fatalf("LoadLesson() unexpected error: %v", err)
	}

	if topicSlug != "domain-driven-design" {
		t.Errorf("topicSlug = %q, want domain-driven-design", topicSlug)
	}
	if lesson.Slug().String() != "ubiquitous-language" {
		t.Errorf("Slug() = %q, want ubiquitous-language", lesson.Slug().String())
	}
	if lesson.Version() != 1 {
		t.Errorf("Version() = %d, want 1", lesson.Version())
	}
	if lesson.TitleEn() != "Ubiquitous Language" {
		t.Errorf("TitleEn() = %q, want Ubiquitous Language", lesson.TitleEn())
	}
	if lesson.EstMinutes().Int() != 7 {
		t.Errorf("EstMinutes() = %d, want 7", lesson.EstMinutes().Int())
	}
	if len(lesson.References()) != 3 {
		t.Errorf("len(References()) = %d, want 3", len(lesson.References()))
	}

	checks := lesson.RecallChecks()
	if len(checks) != 4 {
		t.Fatalf("len(RecallChecks()) = %d, want 4", len(checks))
	}
	for i, c := range checks {
		if c.Position().Int() != i+1 {
			t.Errorf("RecallChecks()[%d].Position() = %d, want %d", i, c.Position().Int(), i+1)
		}
	}
	if checks[1].Kind().String() != "mcq" || len(checks[1].Options()) != 3 {
		t.Errorf("RecallChecks()[1] = kind %q options %v, want mcq with 3 options", checks[1].Kind().String(), checks[1].Options())
	}
	if checks[0].Kind().String() != "short_answer" || len(checks[0].Options()) != 0 {
		t.Errorf("RecallChecks()[0] = kind %q options %v, want short_answer with no options", checks[0].Kind().String(), checks[0].Options())
	}
}

func TestLoadLesson_DefaultsVersionWhenOmitted(t *testing.T) {
	_, lesson, err := LoadLesson(filepath.Join("testdata", "domain-driven-design", "version-omitted.json"))
	if err != nil {
		t.Fatalf("LoadLesson() unexpected error: %v", err)
	}
	if lesson.Version() != 1 {
		t.Errorf("Version() = %d, want default 1", lesson.Version())
	}
}

// Duplicate recall-check positions are exercised at the domain layer only
// (domain/lesson_test.go's "duplicate recall check positions" case): a
// loader-level fixture can't reproduce it because position is the array
// index, not a file field, so two entries in the same recall_checks array
// can never collide.
func TestLoadLesson_Errors(t *testing.T) {
	tests := []struct {
		name        string
		file        string
		wantErr     error
		wantContain []string
	}{
		{name: "malformed JSON", file: "malformed.json", wantContain: []string{"malformed.json"}},
		{name: "unknown top-level field", file: "unknown-field.json", wantContain: []string{"unknown-field.json"}},
		{name: "unknown field nested in a recall_check", file: "unknown-field-nested.json", wantContain: []string{"unknown-field-nested.json"}},
		{name: "trailing content after JSON value", file: "trailing-content.json", wantContain: []string{"trailing-content.json"}},
		{name: "mcq missing options", file: "mcq-missing-options.json", wantErr: domain.ErrInvalidRecallOptions, wantContain: []string{"mcq-missing-options.json"}},
		{name: "short_answer with options", file: "short-answer-with-options.json", wantErr: domain.ErrInvalidRecallOptions, wantContain: []string{"short-answer-with-options.json"}},
		{name: "est_minutes too low", file: "est-minutes-too-low.json", wantErr: domain.ErrInvalidEstMinutes, wantContain: []string{"est-minutes-too-low.json"}},
		{name: "est_minutes too high", file: "est-minutes-too-high.json", wantErr: domain.ErrInvalidEstMinutes, wantContain: []string{"est-minutes-too-high.json"}},
		{name: "one reference", file: "one-reference.json", wantErr: domain.ErrInvalidReferenceCount, wantContain: []string{"one-reference.json"}},
		{name: "five references", file: "five-references.json", wantErr: domain.ErrInvalidReferenceCount, wantContain: []string{"five-references.json"}},
		{name: "two recall checks", file: "two-recall-checks.json", wantErr: domain.ErrInvalidRecallCheckCount, wantContain: []string{"two-recall-checks.json"}},
		{name: "six recall checks", file: "six-recall-checks.json", wantErr: domain.ErrInvalidRecallCheckCount, wantContain: []string{"six-recall-checks.json"}},
		{name: "concept_id does not match filename", file: "mismatched-concept.json", wantContain: []string{"mismatched-concept.json", "concept_id", "something-else"}},
		{name: "topic does not match folder", file: "mismatched-topic.json", wantContain: []string{"mismatched-topic.json", "wrong-topic"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := LoadLesson(filepath.Join("testdata", "domain-driven-design", tt.file))
			if err == nil {
				t.Fatal("LoadLesson() expected error, got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("LoadLesson() error = %v, want it to wrap %v", err, tt.wantErr)
			}
			for _, want := range tt.wantContain {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("LoadLesson() error = %q, want it to contain %q", err.Error(), want)
				}
			}
		})
	}
}

func TestLoadLesson_MissingFile(t *testing.T) {
	_, _, err := LoadLesson(filepath.Join("testdata", "domain-driven-design", "does-not-exist.json"))
	if err == nil {
		t.Fatal("LoadLesson() expected an error for a missing file, got nil")
	}
}
