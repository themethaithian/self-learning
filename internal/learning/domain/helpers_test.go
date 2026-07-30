package domain

import "testing"

func mustLessonRef(t *testing.T, raw string) LessonRef {
	t.Helper()
	r, err := NewLessonRef(raw)
	if err != nil {
		t.Fatalf("NewLessonRef(%q) failed: %v", raw, err)
	}
	return r
}

func mustCanonicalQuestion(t *testing.T, raw string) CanonicalQuestion {
	t.Helper()
	q, err := NewCanonicalQuestion(raw)
	if err != nil {
		t.Fatalf("NewCanonicalQuestion(%q) failed: %v", raw, err)
	}
	return q
}

func mustChunkState(t *testing.T, raw string) ChunkState {
	t.Helper()
	s, err := NewChunkState(raw)
	if err != nil {
		t.Fatalf("NewChunkState(%q) failed: %v", raw, err)
	}
	return s
}

func mustLessonProgress(t *testing.T, ref string, state string) LessonProgress {
	t.Helper()
	p, err := NewLessonProgress(mustLessonRef(t, ref), mustChunkState(t, state))
	if err != nil {
		t.Fatalf("NewLessonProgress(%q, %q) failed: %v", ref, state, err)
	}
	return p
}

// reviewQualityValue builds a ReviewQuality directly from a raw 0-5 grade,
// bypassing NewReviewQuality's (outcome, confidence) mapping so
// EaseFactor/ReviewCard tests exercise the SM-2 algorithm independently of
// that mapping's own correctness.
func reviewQualityValue(t *testing.T, v int) ReviewQuality {
	t.Helper()
	if v < 0 || v > 5 {
		t.Fatalf("reviewQualityValue(%d): out of range 0-5", v)
	}
	return ReviewQuality{grade: v + 1}
}

func mustReviewCard(t *testing.T, question string) ReviewCard {
	t.Helper()
	c, err := NewReviewCard(mustCheckKey(t, "ddia", "b-trees", question))
	if err != nil {
		t.Fatalf("NewReviewCard failed: %v", err)
	}
	return c
}
