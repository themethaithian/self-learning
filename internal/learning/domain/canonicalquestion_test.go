package domain

import (
	"errors"
	"testing"
)

func TestNewCanonicalQuestion(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr error
	}{
		{name: "valid", raw: "What is a B-tree?", want: "What is a B-tree?"},
		{name: "trims leading/trailing whitespace", raw: "  What is a B-tree?  ", want: "What is a B-tree?"},
		{name: "empty", raw: "", wantErr: ErrInvalidCheckKey},
		{name: "whitespace-only", raw: "   ", wantErr: ErrInvalidCheckKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCanonicalQuestion(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewCanonicalQuestion(%q) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewCanonicalQuestion(%q) unexpected error: %v", tt.raw, err)
			}
			if got.String() != tt.want {
				t.Errorf("String() = %q, want %q", got.String(), tt.want)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for a valid question")
			}
		})
	}
}

func TestCanonicalQuestionZeroValue(t *testing.T) {
	var q CanonicalQuestion
	if !q.IsZero() {
		t.Errorf("zero-value CanonicalQuestion.IsZero() = false, want true")
	}
}
