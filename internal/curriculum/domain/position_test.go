package domain

import (
	"errors"
	"math"
	"testing"
)

func TestNewPosition(t *testing.T) {
	// Computed through a variable, not a constant expression: math.MaxInt32+1
	// must truncate at runtime on 32-bit `int` rather than fail `go vet`
	// with "constant overflows int" at compile time.
	var maxInt32AsInt64 int64 = math.MaxInt32
	aboveMaxInt32 := int(maxInt32AsInt64 + 1)

	tests := []struct {
		name    string
		raw     int
		wantErr error
	}{
		{name: "one", raw: 1},
		{name: "large", raw: 9999},
		{name: "max int32 boundary", raw: math.MaxInt32},
		{name: "zero", raw: 0, wantErr: ErrInvalidPosition},
		{name: "negative", raw: -1, wantErr: ErrInvalidPosition},
		{name: "above max int32", raw: aboveMaxInt32, wantErr: ErrInvalidPosition},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPosition(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewPosition(%d) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewPosition(%d) unexpected error: %v", tt.raw, err)
			}
			if got.Int() != tt.raw {
				t.Errorf("Int() = %d, want %d", got.Int(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid position %d", tt.raw)
			}
		})
	}
}

func TestPositionZeroValue(t *testing.T) {
	var p Position
	if !p.IsZero() {
		t.Errorf("zero-value Position.IsZero() = false, want true")
	}
}
