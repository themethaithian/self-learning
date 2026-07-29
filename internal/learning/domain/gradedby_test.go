package domain

import (
	"errors"
	"testing"
)

func TestNewGradedBy(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{name: "self", raw: "self"},
		{name: "llm", raw: "llm"},
		{name: "empty", raw: "", wantErr: ErrInvalidGradedBy},
		{name: "wrong case", raw: "Self", wantErr: ErrInvalidGradedBy},
		{name: "unknown value", raw: "human", wantErr: ErrInvalidGradedBy},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGradedBy(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewGradedBy(%q) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewGradedBy(%q) unexpected error: %v", tt.raw, err)
			}
			if got.String() != tt.raw {
				t.Errorf("String() = %q, want %q", got.String(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid graded by %q", tt.raw)
			}
		})
	}
}

func TestGradedBySelfSingleton(t *testing.T) {
	if GradedBySelf.String() != "self" {
		t.Errorf("GradedBySelf.String() = %q, want %q", GradedBySelf.String(), "self")
	}
}

func TestGradedByZeroValue(t *testing.T) {
	var g GradedBy
	if !g.IsZero() {
		t.Errorf("zero-value GradedBy.IsZero() = false, want true")
	}
}
