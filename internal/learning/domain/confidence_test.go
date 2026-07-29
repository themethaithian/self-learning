package domain

import (
	"errors"
	"testing"
)

func TestNewConfidence(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{name: "guessed", raw: "guessed"},
		{name: "unsure", raw: "unsure"},
		{name: "confident", raw: "confident"},
		{name: "empty", raw: "", wantErr: ErrInvalidConfidence},
		{name: "wrong case", raw: "Guessed", wantErr: ErrInvalidConfidence},
		{name: "unknown value", raw: "certain", wantErr: ErrInvalidConfidence},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfidence(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewConfidence(%q) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewConfidence(%q) unexpected error: %v", tt.raw, err)
			}
			if got.String() != tt.raw {
				t.Errorf("String() = %q, want %q", got.String(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid confidence %q", tt.raw)
			}
		})
	}
}

func TestConfidenceZeroValue(t *testing.T) {
	var c Confidence
	if !c.IsZero() {
		t.Errorf("zero-value Confidence.IsZero() = false, want true")
	}
}
