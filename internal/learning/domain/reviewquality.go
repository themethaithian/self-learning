package domain

import "fmt"

// ReviewQuality is SM-2's 0-5 recall grade for one review; IsCorrect
// reports the algorithm's q >= 3 boundary. The published SM-2 algorithm
// (Wozniak 1990) only defines what to do once you have q — it says nothing
// about how to arrive at one from a real quiz answer. NewReviewQuality's
// mapping from (AttemptOutcome, Confidence) to q is THIS application's own
// design decision, not part of SM-2 itself.
type ReviewQuality struct {
	// grade holds q+1: q=0 is a legitimate SM-2 grade (confident-and-wrong,
	// see qualityMappings), so a raw 0-5 field could not double as "never
	// constructed" the way most VOs in this package use their zero value.
	grade int
}

type qualityMapping struct {
	outcome    AttemptOutcome
	confidence Confidence
	quality    int
}

// qualityMappings is the single source of truth for every (outcome,
// confidence) -> quality pairing; NewReviewQuality is the only constructor,
// so no other pairing can ever exist. Fields are the AttemptOutcome/
// Confidence VOs themselves (indexed from attemptOutcomes/confidences, the
// same arrays those types' own constructors validate against), not bare
// strings — a typo like "corect" here would fail to compile instead of
// silently falling into NewReviewQuality's "no mapping" branch.
//
// Confident-and-wrong maps to 0, BELOW guessed-and-wrong's 2: a confidently
// held wrong belief resists correction more than an acknowledged gap, so it
// should get SM-2's steepest ease-factor penalty (q=0 costs -0.80 vs q=2's
// -0.32 — see EaseFactor.Adjust), making it come back sooner as the
// schedule advances. This distinction only works because Q-1's Commit stage
// captures Confidence BEFORE Reveal — asked afterward, "how confident were
// you" is no longer a measurement of recall, it is a measurement of
// hindsight.
var qualityMappings = [...]qualityMapping{
	{outcome: attemptOutcomes[0], confidence: confidences[2], quality: 5}, // correct, confident
	{outcome: attemptOutcomes[0], confidence: confidences[1], quality: 4}, // correct, unsure
	{outcome: attemptOutcomes[0], confidence: confidences[0], quality: 3}, // correct, guessed
	{outcome: attemptOutcomes[1], confidence: confidences[0], quality: 2}, // incorrect, guessed
	{outcome: attemptOutcomes[1], confidence: confidences[1], quality: 1}, // incorrect, unsure
	{outcome: attemptOutcomes[1], confidence: confidences[2], quality: 0}, // incorrect, confident
}

// NewReviewQuality derives q from Q-1's two collected axes. Every
// (AttemptOutcome, Confidence) combination has an entry in qualityMappings
// (2 outcomes x 3 confidences = the 6 rows above), so the "no mapping"
// branch below can only be reached if a future member is added to one of
// those enums without a matching row here.
func NewReviewQuality(outcome AttemptOutcome, confidence Confidence) (ReviewQuality, error) {
	if outcome.IsZero() {
		return ReviewQuality{}, fmt.Errorf("learning: review quality: outcome is zero: %w", ErrInvalidAttemptOutcome)
	}
	if confidence.IsZero() {
		return ReviewQuality{}, fmt.Errorf("learning: review quality: confidence is zero: %w", ErrInvalidConfidence)
	}
	for _, m := range qualityMappings {
		if m.outcome == outcome && m.confidence == confidence {
			return ReviewQuality{grade: m.quality + 1}, nil
		}
	}
	return ReviewQuality{}, fmt.Errorf("learning: review quality: no mapping for outcome=%s confidence=%s: %w", outcome, confidence, ErrInvalidReviewQuality)
}

func (q ReviewQuality) Grade() int { return q.grade - 1 }

// IsCorrect reports SM-2's q >= 3 boundary: the point at which a review
// advances the repetition count instead of resetting it.
func (q ReviewQuality) IsCorrect() bool { return q.Grade() >= 3 }

// IsZero reports whether q was never constructed via NewReviewQuality.
func (q ReviewQuality) IsZero() bool { return q.grade == 0 }
