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
