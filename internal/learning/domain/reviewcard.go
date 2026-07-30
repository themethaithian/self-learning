package domain

import (
	"fmt"
	"math"
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

// Advance is SM-2's per-review transition (Wozniak 1990): a pure function
// of the card's current state, this review's quality, and when it
// happened — no clock, no I/O, so it is testable without mocking time.
// reviewedAt is the base the new due date is computed from — never the
// previous due date, so a card reviewed late does not compound that
// lateness into its next interval.
//
// The ease factor is adjusted on every review, including a failed one:
// that is the only place this app's confident-vs-guessed-wrong distinction
// (see ReviewQuality) has any effect, since every failure resets the
// interval to the same fixed 1 day regardless of how wrong it was — a
// bigger ease-factor penalty is what makes a confidently-wrong card come
// back sooner once it is passed again and the schedule resumes growing.
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
func nextInterval(previousDays int, ease EaseFactor) int {
	return int(math.Round(float64(previousDays) * ease.Float64()))
}
