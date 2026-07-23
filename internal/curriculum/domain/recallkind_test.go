package domain

import (
	"errors"
	"testing"
)

func TestNewRecallKind(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{name: "short_answer", raw: "short_answer"},
		{name: "mcq", raw: "mcq"},
		{name: "empty", raw: "", wantErr: ErrInvalidRecallKind},
		{name: "wrong case", raw: "MCQ", wantErr: ErrInvalidRecallKind},
		{name: "unknown kind", raw: "essay", wantErr: ErrInvalidRecallKind},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRecallKind(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewRecallKind(%q) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewRecallKind(%q) unexpected error: %v", tt.raw, err)
			}
			if got.String() != tt.raw {
				t.Errorf("String() = %q, want %q", got.String(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid kind %q", tt.raw)
			}
			if want := tt.raw == "mcq"; got.IsMCQ() != want {
				t.Errorf("IsMCQ() = %v, want %v", got.IsMCQ(), want)
			}
		})
	}
}

func TestRecallKindZeroValue(t *testing.T) {
	var k RecallKind
	if !k.IsZero() {
		t.Errorf("zero-value RecallKind.IsZero() = false, want true")
	}
}
