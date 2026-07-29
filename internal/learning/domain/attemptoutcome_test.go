package domain

import (
	"errors"
	"testing"
)

func TestNewAttemptOutcome(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{name: "correct", raw: "correct"},
		{name: "incorrect", raw: "incorrect"},
		{name: "empty", raw: "", wantErr: ErrInvalidAttemptOutcome},
		{name: "wrong case", raw: "Correct", wantErr: ErrInvalidAttemptOutcome},
		{name: "unknown value", raw: "partial", wantErr: ErrInvalidAttemptOutcome},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAttemptOutcome(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewAttemptOutcome(%q) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewAttemptOutcome(%q) unexpected error: %v", tt.raw, err)
			}
			if got.String() != tt.raw {
				t.Errorf("String() = %q, want %q", got.String(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid outcome %q", tt.raw)
			}
		})
	}
}

func TestAttemptOutcomeZeroValue(t *testing.T) {
	var o AttemptOutcome
	if !o.IsZero() {
		t.Errorf("zero-value AttemptOutcome.IsZero() = false, want true")
	}
}
