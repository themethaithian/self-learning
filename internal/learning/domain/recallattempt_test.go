package domain

import (
	"errors"
	"testing"
)

func TestNewCheckKind(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{name: "short_answer", raw: "short_answer"},
		{name: "mcq", raw: "mcq"},
		{name: "empty", raw: "", wantErr: ErrInvalidCheckKind},
		{name: "wrong case", raw: "MCQ", wantErr: ErrInvalidCheckKind},
		{name: "unknown value", raw: "fill_in_blank", wantErr: ErrInvalidCheckKind},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCheckKind(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewCheckKind(%q) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewCheckKind(%q) unexpected error: %v", tt.raw, err)
			}
			if got.String() != tt.raw {
				t.Errorf("String() = %q, want %q", got.String(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid kind %q", tt.raw)
			}
			if got.IsMCQ() != (tt.raw == "mcq") {
				t.Errorf("IsMCQ() = %v, want %v", got.IsMCQ(), tt.raw == "mcq")
			}
		})
	}
}

func TestCheckKindZeroValue(t *testing.T) {
	var k CheckKind
	if !k.IsZero() {
		t.Errorf("zero-value CheckKind.IsZero() = false, want true")
	}
}

func mustCheckKey(t *testing.T, topic, conceptSlug, question string) CheckKey {
	t.Helper()
	k, err := NewCheckKey(topic, mustLessonRef(t, conceptSlug), question)
	if err != nil {
		t.Fatalf("NewCheckKey(%q, %q, %q) failed: %v", topic, conceptSlug, question, err)
	}
	return k
}

func mustCheckKind(t *testing.T, raw string) CheckKind {
	t.Helper()
	k, err := NewCheckKind(raw)
	if err != nil {
		t.Fatalf("NewCheckKind(%q) failed: %v", raw, err)
	}
	return k
}

func mustConfidence(t *testing.T, raw string) Confidence {
	t.Helper()
	c, err := NewConfidence(raw)
	if err != nil {
		t.Fatalf("NewConfidence(%q) failed: %v", raw, err)
	}
	return c
}

func mustAttemptOutcome(t *testing.T, raw string) AttemptOutcome {
	t.Helper()
	o, err := NewAttemptOutcome(raw)
	if err != nil {
		t.Fatalf("NewAttemptOutcome(%q) failed: %v", raw, err)
	}
	return o
}

func strPtr(s string) *string { return &s }

func TestNewRecallAttempt(t *testing.T) {
	validKey := mustCheckKey(t, "ddia", "b-trees", "What is a B-tree?")

	tests := []struct {
		name           string
		checkKey       CheckKey
		kind           CheckKind
		confidence     Confidence
		outcome        AttemptOutcome
		selectedOption *string
		gradedBy       GradedBy
		wantErr        error
	}{
		{
			name:     "valid short_answer, no selected option",
			checkKey: validKey, kind: mustCheckKind(t, "short_answer"),
			confidence: mustConfidence(t, "unsure"), outcome: mustAttemptOutcome(t, "correct"),
			gradedBy: GradedBySelf,
		},
		{
			name:     "valid mcq with selected option",
			checkKey: validKey, kind: mustCheckKind(t, "mcq"),
			confidence: mustConfidence(t, "confident"), outcome: mustAttemptOutcome(t, "incorrect"),
			selectedOption: strPtr("Option B"), gradedBy: GradedBySelf,
		},
		{
			name:     "valid mcq with no selected option",
			checkKey: validKey, kind: mustCheckKind(t, "mcq"),
			confidence: mustConfidence(t, "guessed"), outcome: mustAttemptOutcome(t, "correct"),
			gradedBy: GradedBySelf,
		},
		{
			name:     "selected option on short_answer is rejected",
			checkKey: validKey, kind: mustCheckKind(t, "short_answer"),
			confidence: mustConfidence(t, "unsure"), outcome: mustAttemptOutcome(t, "correct"),
			selectedOption: strPtr("Option B"), gradedBy: GradedBySelf,
			wantErr: ErrSelectedOptionNotMCQ,
		},
		{
			name: "zero check key", checkKey: CheckKey{}, kind: mustCheckKind(t, "mcq"),
			confidence: mustConfidence(t, "unsure"), outcome: mustAttemptOutcome(t, "correct"),
			gradedBy: GradedBySelf, wantErr: ErrInvalidCheckKey,
		},
		{
			name: "zero kind", checkKey: validKey, kind: CheckKind{},
			confidence: mustConfidence(t, "unsure"), outcome: mustAttemptOutcome(t, "correct"),
			gradedBy: GradedBySelf, wantErr: ErrInvalidCheckKind,
		},
		{
			name: "zero confidence", checkKey: validKey, kind: mustCheckKind(t, "mcq"),
			confidence: Confidence{}, outcome: mustAttemptOutcome(t, "correct"),
			gradedBy: GradedBySelf, wantErr: ErrInvalidConfidence,
		},
		{
			name: "zero outcome", checkKey: validKey, kind: mustCheckKind(t, "mcq"),
			confidence: mustConfidence(t, "unsure"), outcome: AttemptOutcome{},
			gradedBy: GradedBySelf, wantErr: ErrInvalidAttemptOutcome,
		},
		{
			name: "zero graded by", checkKey: validKey, kind: mustCheckKind(t, "mcq"),
			confidence: mustConfidence(t, "unsure"), outcome: mustAttemptOutcome(t, "correct"),
			gradedBy: GradedBy{}, wantErr: ErrInvalidGradedBy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRecallAttempt(tt.checkKey, tt.kind, tt.confidence, tt.outcome, tt.selectedOption, tt.gradedBy)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewRecallAttempt() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewRecallAttempt() unexpected error: %v", err)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for a valid attempt")
			}
			if got.CheckKey() != tt.checkKey {
				t.Errorf("CheckKey() = %v, want %v", got.CheckKey(), tt.checkKey)
			}
			if got.Kind() != tt.kind {
				t.Errorf("Kind() = %v, want %v", got.Kind(), tt.kind)
			}
			if got.Confidence() != tt.confidence {
				t.Errorf("Confidence() = %v, want %v", got.Confidence(), tt.confidence)
			}
			if got.Outcome() != tt.outcome {
				t.Errorf("Outcome() = %v, want %v", got.Outcome(), tt.outcome)
			}
			if got.GradedBy() != tt.gradedBy {
				t.Errorf("GradedBy() = %v, want %v", got.GradedBy(), tt.gradedBy)
			}
			gotOpt := got.SelectedOption()
			switch {
			case tt.selectedOption == nil && gotOpt != nil:
				t.Errorf("SelectedOption() = %v, want nil", *gotOpt)
			case tt.selectedOption != nil && (gotOpt == nil || *gotOpt != *tt.selectedOption):
				t.Errorf("SelectedOption() = %v, want %v", gotOpt, *tt.selectedOption)
			}
		})
	}
}

func TestRecallAttemptZeroValue(t *testing.T) {
	var a RecallAttempt
	if !a.IsZero() {
		t.Errorf("zero-value RecallAttempt.IsZero() = false, want true")
	}
}

// TestRecallAttempt_SelectedOptionIsImmutable proves neither the caller's
// input pointer nor the getter's output pointer can be used to mutate a
// RecallAttempt's private state, matching every other VO in this package.
func TestRecallAttempt_SelectedOptionIsImmutable(t *testing.T) {
	validKey := mustCheckKey(t, "ddia", "b-trees", "What is a B-tree?")
	original := "Option B"
	input := original

	a, err := NewRecallAttempt(validKey, mustCheckKind(t, "mcq"), mustConfidence(t, "unsure"), mustAttemptOutcome(t, "correct"), &input, GradedBySelf)
	if err != nil {
		t.Fatalf("NewRecallAttempt() unexpected error: %v", err)
	}

	input = "mutated after construction"
	if got := a.SelectedOption(); got == nil || *got != original {
		t.Fatalf("SelectedOption() = %v after mutating the caller's input pointer, want unaffected %q", got, original)
	}

	got := a.SelectedOption()
	*got = "mutated via the getter's returned pointer"
	if again := a.SelectedOption(); again == nil || *again != original {
		t.Fatalf("SelectedOption() = %v after mutating a previously-returned pointer, want unaffected %q", again, original)
	}
}
