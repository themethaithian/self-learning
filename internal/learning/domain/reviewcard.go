package domain

import (
	"fmt"
	"time"
)

// ReviewCard is one check_key's SM-2 scheduling state. Its identity is the
// CheckKey alone (see checkkey.go) — content-addressed and deliberately
// carrying no link to any recall_checks row, so a card's schedule survives
// a content re-import that deletes and reinserts recall_checks with new
// ids.
type ReviewCard struct {
	checkKey       CheckKey
	easeFactor     EaseFactor
	intervalDays   int
	repetition     int
	dueAt          time.Time
	lastReviewedAt time.Time
}

// NewReviewCard is a fresh, never-reviewed card: SM-2's starting ease
// factor, zero repetitions and interval, and due at the zero time.Time —
// guaranteed to be <= any real "now", so a brand-new card needs no special
// case in a "due now" query.
func NewReviewCard(checkKey CheckKey) (ReviewCard, error) {
	if checkKey.IsZero() {
		return ReviewCard{}, fmt.Errorf("learning: review card: check key is zero: %w", ErrInvalidCheckKey)
	}
	return ReviewCard{checkKey: checkKey, easeFactor: StartingEaseFactor}, nil
}

func (c ReviewCard) CheckKey() CheckKey        { return c.checkKey }
func (c ReviewCard) EaseFactor() EaseFactor    { return c.easeFactor }
func (c ReviewCard) IntervalDays() int         { return c.intervalDays }
func (c ReviewCard) Repetition() int           { return c.repetition }
func (c ReviewCard) DueAt() time.Time          { return c.dueAt }
func (c ReviewCard) LastReviewedAt() time.Time { return c.lastReviewedAt }

// IsZero reports whether c was never constructed via NewReviewCard.
func (c ReviewCard) IsZero() bool { return c.checkKey.IsZero() }

// Advance is SM-2's per-review transition, a pure function of the card's
// current state, this review's quality, and when it happened — no clock,
// no I/O, so it is testable without mocking time. reviewedAt is the base
// the new due date is computed from — never the previous due date, so a
// card reviewed late does not compound that lateness into its next
// interval.
//
// Deliberate deviation from Wozniak 1990's step 6: the published algorithm
// says a failed review (q<3) restarts the interval schedule "without
// changing the E-Factor" — a plain instruction, not an ambiguity. This
// implementation adjusts the ease factor on every review instead, failed
// ones included. The reason: under the spec's own rule, q=0/1/2 would
// produce byte-identical next state (repetition and interval already reset
// to the same values regardless of which failing grade was given), which
// makes the Confidence axis a UI-only distinction with zero effect on
// scheduling. Adjusting EF on failure too is what lets a confidently-wrong
// answer (q=0, EF penalty -0.80) earn a harder ease-factor hit than a
// guessed-wrong one (q=2, -0.32, see EaseFactor.Adjust) — the schedule
// grows back more slowly for the card that needs closer attention, once it
// is passed again. This is the only channel that distinction has any
// effect through, since both failures reset today's interval to the same
// fixed 1 day either way.
func (c ReviewCard) Advance(quality ReviewQuality, reviewedAt time.Time) (ReviewCard, error) {
	if c.IsZero() {
		return ReviewCard{}, fmt.Errorf("learning: review card: advance: card is zero: %w", ErrInvalidReviewCard)
	}
	if quality.IsZero() {
		return ReviewCard{}, fmt.Errorf("learning: review card: advance: quality is zero: %w", ErrInvalidReviewQuality)
	}

	next := ReviewCard{
		checkKey:       c.checkKey,
		easeFactor:     c.easeFactor.Adjust(quality),
		lastReviewedAt: reviewedAt,
	}

	if !quality.IsCorrect() {
		next.repetition = 0
		next.intervalDays = 1
		next.dueAt = reviewedAt.AddDate(0, 0, next.intervalDays)
		return next, nil
	}

	switch c.repetition {
	case 0:
		next.intervalDays = 1
	case 1:
		next.intervalDays = 6
	default:
		next.intervalDays = nextInterval(c.intervalDays, c.easeFactor)
	}
	next.repetition = c.repetition + 1
	next.dueAt = reviewedAt.AddDate(0, 0, next.intervalDays)
	return next, nil
}

// nextInterval computes I(n) := I(n-1) * EF for a card's third and later
// consecutive correct reviews, using the ease factor as it stood BEFORE
// this review's adjustment — the standard SM-2 ordering: this review's
// grade changes the ease factor used for the review AFTER next, not this
// one's own interval.
//
// Rounds to the nearest day, not up: Wozniak's step 3 specifies ceiling
// ("if interval is a fraction, round it up"), a deliberate deviation here
// because ceiling is a one-directional bias that lengthens every
// multi-review interval — the wrong direction for exam-prep drilling,
// where a card surfacing a day early costs nothing and one surfacing a day
// late risks the answer being gone.
//
// previousDays*ease.Hundredths() is exact integer arithmetic (see
// EaseFactor's type doc); adding half of 100 before truncating with /100 is
// round-half-up without ever converting to float64. A float64 multiply
// (previousDays * ease as a decimal) can round the wrong way on an exact
// .5 tie: 95*2.30 is exactly 218.5, but float64 cannot represent 2.30
// exactly, so that path silently produces 218 instead of 219 (see
// TestNextInterval's exact-tie case).
func nextInterval(previousDays int, ease EaseFactor) int {
	return (previousDays*ease.Hundredths() + 50) / 100
}
