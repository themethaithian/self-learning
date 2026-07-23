package domain

import (
	"errors"
	"testing"
)

func TestNewEstMinutes(t *testing.T) {
	tests := []struct {
		name    string
		raw     int
		wantErr error
	}{
		{name: "min boundary", raw: 5},
		{name: "mid", raw: 7},
		{name: "max boundary", raw: 10},
		{name: "below min", raw: 4, wantErr: ErrInvalidEstMinutes},
		{name: "above max", raw: 11, wantErr: ErrInvalidEstMinutes},
		{name: "zero", raw: 0, wantErr: ErrInvalidEstMinutes},
		{name: "negative", raw: -1, wantErr: ErrInvalidEstMinutes},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEstMinutes(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewEstMinutes(%d) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewEstMinutes(%d) unexpected error: %v", tt.raw, err)
			}
			if got.Int() != tt.raw {
				t.Errorf("Int() = %d, want %d", got.Int(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid est minutes %d", tt.raw)
			}
		})
	}
}

func TestEstMinutesZeroValue(t *testing.T) {
	var e EstMinutes
	if !e.IsZero() {
		t.Errorf("zero-value EstMinutes.IsZero() = false, want true")
	}
}
