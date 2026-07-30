package domain

import (
	"errors"
	"testing"
)

func TestNewEaseFactor(t *testing.T) {
	tests := []struct {
		name       string
		hundredths int
		wantErr    bool
	}{
		{name: "floor exactly", hundredths: 130},
		{name: "starting value", hundredths: 250},
		{name: "above starting", hundredths: 400},
		{name: "just below floor", hundredths: 129, wantErr: true},
		{name: "zero", hundredths: 0, wantErr: true},
		{name: "negative", hundredths: -10, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEaseFactor(tt.hundredths)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidEaseFactor) {
					t.Fatalf("NewEaseFactor(%d) error = %v, want wrapping ErrInvalidEaseFactor", tt.hundredths, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewEaseFactor(%d) unexpected error: %v", tt.hundredths, err)
			}
			if got.Hundredths() != tt.hundredths {
				t.Errorf("Hundredths() = %d, want %d", got.Hundredths(), tt.hundredths)
			}
		})
	}
}

func TestEaseFactorZeroValue(t *testing.T) {
	var e EaseFactor
	if !e.IsZero() {
		t.Errorf("zero-value EaseFactor.IsZero() = false, want true")
	}
}

func TestStartingEaseFactor(t *testing.T) {
	if got := StartingEaseFactor.Hundredths(); got != 250 {
		t.Errorf("StartingEaseFactor.Hundredths() = %d, want 250 (SM-2's published starting EF of 2.5)", got)
	}
}

// TestEaseFactorAdjust_AllSixQualities pins the recurrence's output for
// every ReviewQuality against values hand-derived from the published
// formula EF' = EF + (0.1 - (5-q)*(0.08+(5-q)*0.02)) starting at EF=2.50 —
// not from running this package's own code:
//
//	q=5 (d=0): delta = 0.10                    -> 2.60
//	q=4 (d=1): delta = 0.10 - 0.08 - 0.02       -> 2.50
//	q=3 (d=2): delta = 0.10 - 0.16 - 0.08 = -0.14 -> 2.36
//	q=2 (d=3): delta = 0.10 - 0.24 - 0.18 = -0.32 -> 2.18
//	q=1 (d=4): delta = 0.10 - 0.32 - 0.32 = -0.54 -> 1.96
//	q=0 (d=5): delta = 0.10 - 0.40 - 0.50 = -0.80 -> 1.70
func TestEaseFactorAdjust_AllSixQualities(t *testing.T) {
	tests := []struct {
		quality        int
		wantHundredths int
	}{
		{quality: 5, wantHundredths: 260},
		{quality: 4, wantHundredths: 250},
		{quality: 3, wantHundredths: 236},
		{quality: 2, wantHundredths: 218},
		{quality: 1, wantHundredths: 196},
		{quality: 0, wantHundredths: 170},
	}
	for _, tt := range tests {
		q := reviewQualityValue(t, tt.quality)
		got := StartingEaseFactor.Adjust(q)
		if got.Hundredths() != tt.wantHundredths {
			t.Errorf("StartingEaseFactor.Adjust(q=%d) = %d, want %d", tt.quality, got.Hundredths(), tt.wantHundredths)
		}
	}
}

// TestEaseFactorAdjust_FloorReachedAndHeld pins the hard floor: repeated
// q=0 reviews (each -0.80) push a card below 1.30, and once there, another
// failing review must hold at exactly 1.30, never lower.
func TestEaseFactorAdjust_FloorReachedAndHeld(t *testing.T) {
	q0 := reviewQualityValue(t, 0)

	first := StartingEaseFactor.Adjust(q0)
	if first.Hundredths() != 170 {
		t.Fatalf("after 1st failing review: Hundredths() = %d, want 170 (2.50 - 0.80)", first.Hundredths())
	}
	second := first.Adjust(q0)
	if second.Hundredths() != 130 {
		t.Fatalf("after 2nd failing review: Hundredths() = %d, want 130 (1.70 - 0.80 = 0.90, floored)", second.Hundredths())
	}
	third := second.Adjust(q0)
	if third.Hundredths() != 130 {
		t.Fatalf("after 3rd failing review at the floor: Hundredths() = %d, want 130 (held, not below)", third.Hundredths())
	}
}

// TestNextInterval pins I(n) := I(n-1) * EF's rounding rule, including an
// exact .5 tie, independent of any ReviewCard plumbing.
func TestNextInterval(t *testing.T) {
	tests := []struct {
		name         string
		previousDays int
		hundredths   int
		want         int
	}{
		{name: "rounds down", previousDays: 6, hundredths: 270, want: 16},    // 16.2
		{name: "half rounds up", previousDays: 5, hundredths: 250, want: 13}, // 12.5
		{name: "exact", previousDays: 10, hundredths: 200, want: 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ease, err := NewEaseFactor(tt.hundredths)
			if err != nil {
				t.Fatalf("NewEaseFactor(%d) unexpected error: %v", tt.hundredths, err)
			}
			if got := nextInterval(tt.previousDays, ease); got != tt.want {
				t.Errorf("nextInterval(%d, %.2f) = %d, want %d", tt.previousDays, ease.Float64(), got, tt.want)
			}
		})
	}
}
