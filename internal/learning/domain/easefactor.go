package domain

import "fmt"

const (
	// easeFactorFloorHundredths is SM-2's published hard floor, EF=1.3.
	easeFactorFloorHundredths = 130
	// easeFactorStartHundredths is SM-2's published starting value, EF=2.5.
	easeFactorStartHundredths = 250
)

// EaseFactor is SM-2's per-card difficulty multiplier, represented as an
// integer count of hundredths (250 = 2.50) rather than a float. Every
// ReviewQuality is an integer 0-5, so the published recurrence's
// adjustment (see Adjust) always lands on an exact multiple of 0.01 —
// hundredths-as-int keeps every update exact integer arithmetic, so the
// binary-fraction rounding a float64 read-modify-write would accumulate
// across hundreds of reviews cannot happen here.
//
// Deliberately no float64 accessor: a fixed-point type that hands out a
// lossy escape hatch invites exactly the bug that escape hatch causes. Add
// one only when something outside this package genuinely needs a decimal
// string, and name it for what it produces (e.g. DecimalString), not as a
// numeric type that reads as safe to compute with.
type EaseFactor struct {
	hundredths int
}

// StartingEaseFactor is a fresh card's ease factor before its first review.
var StartingEaseFactor = EaseFactor{hundredths: easeFactorStartHundredths}

// NewEaseFactor reconstructs an already-valid, previously computed
// EaseFactor from its stored hundredths value. It exists for loading
// persisted state; Adjust is the only way to produce a new value mid
// algorithm.
func NewEaseFactor(hundredths int) (EaseFactor, error) {
	if hundredths < easeFactorFloorHundredths {
		return EaseFactor{}, fmt.Errorf("learning: ease factor %d: below floor %d: %w", hundredths, easeFactorFloorHundredths, ErrInvalidEaseFactor)
	}
	return EaseFactor{hundredths: hundredths}, nil
}

func (e EaseFactor) Hundredths() int { return e.hundredths }

// Adjust applies SM-2's published recurrence for one review's quality:
// EF' = EF + (0.1 - (5-q)*(0.08+(5-q)*0.02)), floored at 1.3. d is 5-q,
// ranging over 0-5; the expression below is that same formula scaled by
// 100 and reordered into pure integer arithmetic (10 = 0.1*100, 8 =
// 0.08*100, 2 = 0.02*100) — no float anywhere in the computation.
func (e EaseFactor) Adjust(quality ReviewQuality) EaseFactor {
	d := 5 - quality.Grade()
	delta := 10 - d*(8+2*d)
	next := e.hundredths + delta
	if next < easeFactorFloorHundredths {
		next = easeFactorFloorHundredths
	}
	return EaseFactor{hundredths: next}
}

// IsZero reports whether e was never constructed.
func (e EaseFactor) IsZero() bool { return e.hundredths == 0 }
