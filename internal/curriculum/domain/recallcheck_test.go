package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewRecallCheck(t *testing.T) {
	tests := []struct {
		name            string
		position        Position
		kind            RecallKind
		question        string
		expectedAnswer  string
		options         []string
		explanation     string
		wantErr         error
		wantOptionsLen  int
		wantExplanation string
	}{
		{
			name: "valid short_answer", position: mustPosition(t, 1), kind: mustRecallKind(t, "short_answer"),
			question: "Aggregate root คืออะไร?", expectedAnswer: "entity ที่เป็นทางเข้าเดียวของ aggregate",
		},
		{
			name: "valid mcq", position: mustPosition(t, 1), kind: mustRecallKind(t, "mcq"),
			question: "ข้อใดถูกต้อง?", expectedAnswer: "b", options: []string{"a", "b", "c"}, wantOptionsLen: 3,
		},
		{
			name: "valid mcq with explanation", position: mustPosition(t, 1), kind: mustRecallKind(t, "mcq"),
			question: "ข้อใดถูกต้อง?", expectedAnswer: "b", options: []string{"a", "b", "c"}, wantOptionsLen: 3,
			explanation: "b ถูกเพราะตรง constraint ใน stem", wantExplanation: "b ถูกเพราะตรง constraint ใน stem",
		},
		{
			name: "explanation trimmed", position: mustPosition(t, 1), kind: mustRecallKind(t, "short_answer"),
			question: "q", expectedAnswer: "a", explanation: "  padded  ", wantExplanation: "padded",
		},
		{
			name: "whitespace-only explanation normalises to absent, not an error", position: mustPosition(t, 1), kind: mustRecallKind(t, "short_answer"),
			question: "q", expectedAnswer: "a", explanation: "   ", wantExplanation: "",
		},
		{
			name: "explanation exceeding max runes", position: mustPosition(t, 1), kind: mustRecallKind(t, "short_answer"),
			question: "q", expectedAnswer: "a", explanation: strings.Repeat("อ", maxRecallExplanationRunes+1),
			wantErr: ErrInvalidRecallExplanation,
		},
		{
			name: "explanation exactly at max runes is valid", position: mustPosition(t, 1), kind: mustRecallKind(t, "short_answer"),
			question: "q", expectedAnswer: "a", explanation: strings.Repeat("อ", maxRecallExplanationRunes),
			wantExplanation: strings.Repeat("อ", maxRecallExplanationRunes),
		},
		{
			name: "zero-value position", position: Position{}, kind: mustRecallKind(t, "short_answer"),
			question: "q", expectedAnswer: "a", wantErr: ErrInvalidPosition,
		},
		{
			name: "zero-value kind", position: mustPosition(t, 1), kind: RecallKind{},
			question: "q", expectedAnswer: "a", wantErr: ErrInvalidRecallKind,
		},
		{
			name: "empty question", position: mustPosition(t, 1), kind: mustRecallKind(t, "short_answer"),
			question: "", expectedAnswer: "a", wantErr: ErrInvalidRecallQuestion,
		},
		{
			name: "empty expected answer", position: mustPosition(t, 1), kind: mustRecallKind(t, "short_answer"),
			question: "q", expectedAnswer: "", wantErr: ErrInvalidRecallAnswer,
		},
		{
			name: "mcq with zero options", position: mustPosition(t, 1), kind: mustRecallKind(t, "mcq"),
			question: "q", expectedAnswer: "a", options: nil, wantErr: ErrInvalidRecallOptions,
		},
		{
			name: "mcq with one option", position: mustPosition(t, 1), kind: mustRecallKind(t, "mcq"),
			question: "q", expectedAnswer: "a", options: []string{"a"}, wantErr: ErrInvalidRecallOptions,
		},
		{
			name: "short_answer with options", position: mustPosition(t, 1), kind: mustRecallKind(t, "short_answer"),
			question: "q", expectedAnswer: "a", options: []string{"a", "b"}, wantErr: ErrInvalidRecallOptions,
		},
		{
			name: "mcq with an empty option string", position: mustPosition(t, 1), kind: mustRecallKind(t, "mcq"),
			question: "q", expectedAnswer: "a", options: []string{"a", ""}, wantErr: ErrInvalidRecallOptions,
		},
		{
			name: "mcq answer not among options", position: mustPosition(t, 1), kind: mustRecallKind(t, "mcq"),
			question: "q", expectedAnswer: "z", options: []string{"a", "b", "c"}, wantErr: ErrInvalidRecallAnswer,
		},
		{
			name: "mcq answer matches option only after trimming whitespace", position: mustPosition(t, 1), kind: mustRecallKind(t, "mcq"),
			question: "q", expectedAnswer: "  b  ", options: []string{"a", "b", "c"}, wantOptionsLen: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRecallCheck(tt.position, tt.kind, tt.question, tt.expectedAnswer, tt.options, tt.explanation)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewRecallCheck() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewRecallCheck() unexpected error: %v", err)
			}
			if got.Position() != tt.position {
				t.Errorf("Position() = %v, want %v", got.Position(), tt.position)
			}
			if got.Kind() != tt.kind {
				t.Errorf("Kind() = %v, want %v", got.Kind(), tt.kind)
			}
			if len(got.Options()) != tt.wantOptionsLen {
				t.Errorf("len(Options()) = %d, want %d", len(got.Options()), tt.wantOptionsLen)
			}
			if got.Explanation() != tt.wantExplanation {
				t.Errorf("Explanation() = %q, want %q", got.Explanation(), tt.wantExplanation)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for a validly constructed RecallCheck")
			}
		})
	}
}

func TestRecallCheckOptionsDefensiveCopy(t *testing.T) {
	options := []string{"a", "b"}
	rc := mustRecallCheck(t, 1, "mcq", "q", "a", options)

	options[0] = "mutated"
	got := rc.Options()
	if got[0] != "a" {
		t.Errorf("RecallCheck mutated after caller modified original slice: got %q, want %q", got[0], "a")
	}

	got[1] = "mutated"
	again := rc.Options()
	if again[1] != "b" {
		t.Errorf("Options() mutated externally: got %q, want %q", again[1], "b")
	}
}

func TestRecallCheckShortAnswerOptionsIsNil(t *testing.T) {
	rc := mustRecallCheck(t, 1, "short_answer", "q", "a", nil)
	if rc.Options() != nil {
		t.Errorf("Options() = %v, want nil for a short_answer check", rc.Options())
	}
}

func TestRecallCheckZeroValue(t *testing.T) {
	var rc RecallCheck
	if !rc.IsZero() {
		t.Errorf("zero-value RecallCheck.IsZero() = false, want true")
	}
}
