package domain

import (
	"errors"
	"slices"
	"testing"
)

// TestNewReviewQuality pins all 6 (outcome, confidence) -> quality pairings
// against the table from CLAUDE.md / the ticket, hand-written here rather
// than derived from qualityMappings itself.
func TestNewReviewQuality(t *testing.T) {
	tests := []struct {
		name       string
		outcome    string
		confidence string
		want       int
	}{
		{name: "correct confident", outcome: "correct", confidence: "confident", want: 5},
		{name: "correct unsure", outcome: "correct", confidence: "unsure", want: 4},
		{name: "correct guessed", outcome: "correct", confidence: "guessed", want: 3},
		{name: "incorrect guessed", outcome: "incorrect", confidence: "guessed", want: 2},
		{name: "incorrect unsure", outcome: "incorrect", confidence: "unsure", want: 1},
		{name: "incorrect confident", outcome: "incorrect", confidence: "confident", want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewReviewQuality(mustAttemptOutcome(t, tt.outcome), mustConfidence(t, tt.confidence))
			if err != nil {
				t.Fatalf("NewReviewQuality(%s, %s) unexpected error: %v", tt.outcome, tt.confidence, err)
			}
			if got := q.Value(); got != tt.want {
				t.Errorf("NewReviewQuality(%s, %s).Value() = %d, want %d", tt.outcome, tt.confidence, got, tt.want)
			}
		})
	}
}

// TestNewReviewQuality_ConfidentIncorrectIsLowestGrade pins this mapping's
// entire reason for existing: confident-and-wrong must grade strictly BELOW
// guessed-and-wrong, not merely differently and not above it.
func TestNewReviewQuality_ConfidentIncorrectIsLowestGrade(t *testing.T) {
	confidentWrong, err := NewReviewQuality(mustAttemptOutcome(t, "incorrect"), mustConfidence(t, "confident"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if confidentWrong.Value() != 0 {
		t.Fatalf("confident-incorrect quality = %d, want 0", confidentWrong.Value())
	}

	guessedWrong, err := NewReviewQuality(mustAttemptOutcome(t, "incorrect"), mustConfidence(t, "guessed"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if confidentWrong.Value() >= guessedWrong.Value() {
		t.Fatalf("confident-incorrect quality (%d) must be strictly below guessed-incorrect quality (%d)",
			confidentWrong.Value(), guessedWrong.Value())
	}
}

func TestNewReviewQuality_ZeroInputs(t *testing.T) {
	confident := mustConfidence(t, "confident")
	correct := mustAttemptOutcome(t, "correct")

	if _, err := NewReviewQuality(AttemptOutcome{}, confident); !errors.Is(err, ErrInvalidAttemptOutcome) {
		t.Errorf("zero outcome: error = %v, want wrapping ErrInvalidAttemptOutcome", err)
	}
	if _, err := NewReviewQuality(correct, Confidence{}); !errors.Is(err, ErrInvalidConfidence) {
		t.Errorf("zero confidence: error = %v, want wrapping ErrInvalidConfidence", err)
	}
}

func TestReviewQualityZeroValue(t *testing.T) {
	var q ReviewQuality
	if !q.IsZero() {
		t.Errorf("zero-value ReviewQuality.IsZero() = false, want true")
	}
}

// TestReviewQualityIsCorrect_Boundary pins SM-2's q >= 3 boundary exactly
// at 3: a mutant changing it to > 3 or >= 2 must fail one of these two
// cases.
func TestReviewQualityIsCorrect_Boundary(t *testing.T) {
	tests := []struct {
		outcome, confidence string
		wantValue           int
		wantCorrect         bool
	}{
		{outcome: "correct", confidence: "guessed", wantValue: 3, wantCorrect: true},
		{outcome: "incorrect", confidence: "guessed", wantValue: 2, wantCorrect: false},
	}
	for _, tt := range tests {
		q, err := NewReviewQuality(mustAttemptOutcome(t, tt.outcome), mustConfidence(t, tt.confidence))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if q.Value() != tt.wantValue {
			t.Fatalf("Value() = %d, want %d", q.Value(), tt.wantValue)
		}
		if got := q.IsCorrect(); got != tt.wantCorrect {
			t.Errorf("quality %d IsCorrect() = %v, want %v", tt.wantValue, got, tt.wantCorrect)
		}
	}
}

// TestReviewQualityIsCorrect_AllSixValues extends the boundary check across
// every grade, not only the two nearest the cutoff.
func TestReviewQualityIsCorrect_AllSixValues(t *testing.T) {
	want := map[int]bool{5: true, 4: true, 3: true, 2: false, 1: false, 0: false}
	for value, wantCorrect := range want {
		q := reviewQualityValue(t, value)
		if got := q.IsCorrect(); got != wantCorrect {
			t.Errorf("quality %d IsCorrect() = %v, want %v", value, got, wantCorrect)
		}
	}
}

// TestQualityMappingsIsExactly pins the exhaustive set of 6 (outcome,
// confidence) -> quality rows against a hand-written literal — catching a
// mapping ADDED, REMOVED, or altered by one grade, not just the specific
// combinations TestNewReviewQuality happens to enumerate.
func TestQualityMappingsIsExactly(t *testing.T) {
	want := []qualityMapping{
		{outcome: "correct", confidence: "confident", quality: 5},
		{outcome: "correct", confidence: "unsure", quality: 4},
		{outcome: "correct", confidence: "guessed", quality: 3},
		{outcome: "incorrect", confidence: "guessed", quality: 2},
		{outcome: "incorrect", confidence: "unsure", quality: 1},
		{outcome: "incorrect", confidence: "confident", quality: 0},
	}
	if got := qualityMappings[:]; !slices.Equal(got, want) {
		t.Errorf("qualityMappings = %+v, want exactly %+v", got, want)
	}
}
