package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewReviewCard(t *testing.T) {
	card := mustReviewCard(t, "What is a B-tree?")
	if card.IsZero() {
		t.Fatalf("NewReviewCard result IsZero() = true, want false")
	}
	if card.EaseFactor().Hundredths() != 250 {
		t.Errorf("EaseFactor().Hundredths() = %d, want 250 (SM-2's starting 2.5)", card.EaseFactor().Hundredths())
	}
	if card.IntervalDays() != 0 || card.Repetition() != 0 {
		t.Errorf("fresh card IntervalDays()/Repetition() = %d/%d, want 0/0", card.IntervalDays(), card.Repetition())
	}
}

func TestNewReviewCard_ZeroCheckKey(t *testing.T) {
	if _, err := NewReviewCard(CheckKey{}); !errors.Is(err, ErrInvalidCheckKey) {
		t.Errorf("NewReviewCard(zero CheckKey) error = %v, want wrapping ErrInvalidCheckKey", err)
	}
}

func TestReviewCardZeroValue(t *testing.T) {
	var c ReviewCard
	if !c.IsZero() {
		t.Errorf("zero-value ReviewCard.IsZero() = false, want true")
	}
}

func TestReviewCard_Advance_ZeroInputs(t *testing.T) {
	card := mustReviewCard(t, "What is a B-tree?")
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	if _, err := (ReviewCard{}).Advance(reviewQualityValue(t, 5), now); !errors.Is(err, ErrInvalidReviewCard) {
		t.Errorf("Advance on zero card: error = %v, want wrapping ErrInvalidReviewCard", err)
	}
	if _, err := card.Advance(ReviewQuality{}, now); !errors.Is(err, ErrInvalidReviewQuality) {
		t.Errorf("Advance with zero quality: error = %v, want wrapping ErrInvalidReviewQuality", err)
	}
}

// TestReviewCard_FirstThreeReviews walks a fresh card through three
// consecutive q=5 reviews and pins every resulting field against values
// hand-derived from the published algorithm (Wozniak 1990):
//
//	review 1: n=0 -> I(1)=1;  EF 2.50 -> 2.60 (delta +0.10 at q=5)
//	review 2: n=1 -> I(2)=6;  EF 2.60 -> 2.70
//	review 3: n=2 -> I(3)=round(I(2)*EF_before_review_3)=round(6*2.70)=16; EF 2.70 -> 2.80
//
// Each review happens well past the card's own previous due date, which
// also proves the new due date is computed from reviewedAt, not from the
// stale due date that review inherited.
func TestReviewCard_FirstThreeReviews(t *testing.T) {
	card := mustReviewCard(t, "What is a B-tree?")
	q5 := reviewQualityValue(t, 5)

	reviewedAt1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	after1, err := card.Advance(q5, reviewedAt1)
	if err != nil {
		t.Fatalf("review 1: unexpected error: %v", err)
	}
	assertCardState(t, "review 1", after1, 1, 1, 260, reviewedAt1.AddDate(0, 0, 1))

	reviewedAt2 := reviewedAt1.AddDate(0, 0, 40) // deliberately late vs. after1.DueAt()
	after2, err := after1.Advance(q5, reviewedAt2)
	if err != nil {
		t.Fatalf("review 2: unexpected error: %v", err)
	}
	assertCardState(t, "review 2", after2, 2, 6, 270, reviewedAt2.AddDate(0, 0, 6))

	reviewedAt3 := reviewedAt2.AddDate(0, 0, 90) // deliberately late vs. after2.DueAt()
	after3, err := after2.Advance(q5, reviewedAt3)
	if err != nil {
		t.Fatalf("review 3: unexpected error: %v", err)
	}
	assertCardState(t, "review 3", after3, 3, 16, 280, reviewedAt3.AddDate(0, 0, 16))
}

func assertCardState(t *testing.T, label string, card ReviewCard, wantRepetition, wantIntervalDays, wantEaseHundredths int, wantDueAt time.Time) {
	t.Helper()
	if card.Repetition() != wantRepetition {
		t.Errorf("%s: Repetition() = %d, want %d", label, card.Repetition(), wantRepetition)
	}
	if card.IntervalDays() != wantIntervalDays {
		t.Errorf("%s: IntervalDays() = %d, want %d", label, card.IntervalDays(), wantIntervalDays)
	}
	if got := card.EaseFactor().Hundredths(); got != wantEaseHundredths {
		t.Errorf("%s: EaseFactor().Hundredths() = %d, want %d", label, got, wantEaseHundredths)
	}
	if !card.DueAt().Equal(wantDueAt) {
		t.Errorf("%s: DueAt() = %v, want %v", label, card.DueAt(), wantDueAt)
	}
}

// TestReviewCard_Advance_DueDateBasedOnReviewedAt isolates the due-date
// base case: the new due date must come from THIS review's reviewedAt, not
// from the card's own previous due date, even when a review happens long
// after it was due.
func TestReviewCard_Advance_DueDateBasedOnReviewedAt(t *testing.T) {
	card := mustReviewCard(t, "What is a B-tree?")
	firstReviewedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	after1, err := card.Advance(reviewQualityValue(t, 5), firstReviewedAt)
	if err != nil {
		t.Fatalf("first Advance: unexpected error: %v", err)
	}
	wantFirstDue := firstReviewedAt.AddDate(0, 0, 1)
	if !after1.DueAt().Equal(wantFirstDue) {
		t.Fatalf("DueAt() = %v, want %v", after1.DueAt(), wantFirstDue)
	}

	secondReviewedAt := after1.DueAt().AddDate(0, 0, 40) // 40 days late
	after2, err := after1.Advance(reviewQualityValue(t, 5), secondReviewedAt)
	if err != nil {
		t.Fatalf("second Advance: unexpected error: %v", err)
	}
	wantSecondDue := secondReviewedAt.AddDate(0, 0, 6) // computed from reviewedAt, not after1.DueAt()
	if !after2.DueAt().Equal(wantSecondDue) {
		t.Fatalf("DueAt() = %v, want %v (based on reviewedAt, not the previous due date)", after2.DueAt(), wantSecondDue)
	}
}

// TestReviewCard_FailedReviewAfterLongStreak builds a 3-review correct
// streak, then fails the 4th review, and pins the full reset: repetition
// back to 0, interval back to the fixed 1 day, and the ease factor still
// adjusted by the failing quality's own penalty (not left untouched).
func TestReviewCard_FailedReviewAfterLongStreak(t *testing.T) {
	card := mustReviewCard(t, "What is a B-tree?")
	q5 := reviewQualityValue(t, 5)
	reviewedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	streak := card
	var err error
	for i := 0; i < 3; i++ {
		streak, err = streak.Advance(q5, reviewedAt)
		if err != nil {
			t.Fatalf("review %d: unexpected error: %v", i+1, err)
		}
		reviewedAt = reviewedAt.AddDate(0, 0, streak.IntervalDays())
	}
	if streak.Repetition() != 3 || streak.IntervalDays() != 16 || streak.EaseFactor().Hundredths() != 280 {
		t.Fatalf("after 3 correct reviews: repetition=%d interval=%d ease=%d, want 3/16/280",
			streak.Repetition(), streak.IntervalDays(), streak.EaseFactor().Hundredths())
	}

	failedAt := reviewedAt
	failed, err := streak.Advance(reviewQualityValue(t, 1), failedAt) // q=1: d=4, delta=-54
	if err != nil {
		t.Fatalf("failed review: unexpected error: %v", err)
	}
	if failed.Repetition() != 0 {
		t.Errorf("Repetition() = %d, want 0 (a failed review resets the streak)", failed.Repetition())
	}
	if failed.IntervalDays() != 1 {
		t.Errorf("IntervalDays() = %d, want 1 (SM-2 restarts at I(1) on failure)", failed.IntervalDays())
	}
	if wantEase := 280 - 54; failed.EaseFactor().Hundredths() != wantEase {
		t.Errorf("EaseFactor().Hundredths() = %d, want %d", failed.EaseFactor().Hundredths(), wantEase)
	}
	if wantDue := failedAt.AddDate(0, 0, 1); !failed.DueAt().Equal(wantDue) {
		t.Errorf("DueAt() = %v, want %v", failed.DueAt(), wantDue)
	}
}

// TestReviewCard_FirstReviewFails covers a fresh card's very first review
// failing outright: repetition must stay at 0 (not go negative or panic)
// and the interval must still be the fixed 1 day.
func TestReviewCard_FirstReviewFails(t *testing.T) {
	card := mustReviewCard(t, "What is a B-tree?")
	reviewedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	failed, err := card.Advance(reviewQualityValue(t, 0), reviewedAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if failed.Repetition() != 0 {
		t.Errorf("Repetition() = %d, want 0", failed.Repetition())
	}
	if failed.IntervalDays() != 1 {
		t.Errorf("IntervalDays() = %d, want 1", failed.IntervalDays())
	}
	if wantEase := 250 - 80; failed.EaseFactor().Hundredths() != wantEase {
		t.Errorf("EaseFactor().Hundredths() = %d, want %d", failed.EaseFactor().Hundredths(), wantEase)
	}
}
