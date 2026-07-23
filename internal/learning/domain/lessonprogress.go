package domain

import "fmt"

// LessonProgress tracks one lesson's read/recall progress: which lesson (by
// LessonRef) and its current ChunkState. State only ever moves forward —
// locked -> in_progress -> passed — through the transition methods below,
// which return a new value and never mutate the receiver.
type LessonProgress struct {
	lessonRef LessonRef
	state     ChunkState
}

func NewLessonProgress(lessonRef LessonRef, state ChunkState) (LessonProgress, error) {
	if lessonRef.IsZero() {
		return LessonProgress{}, fmt.Errorf("learning: lesson progress: lesson ref: %w", ErrInvalidLessonRef)
	}
	if state.IsZero() {
		return LessonProgress{}, fmt.Errorf("learning: lesson progress %q: state: %w", lessonRef.String(), ErrInvalidChunkState)
	}
	return LessonProgress{lessonRef: lessonRef, state: state}, nil
}

func (p LessonProgress) LessonRef() LessonRef { return p.lessonRef }
func (p LessonProgress) State() ChunkState    { return p.state }

// IsZero reports whether p was never constructed via NewLessonProgress.
func (p LessonProgress) IsZero() bool { return p.lessonRef.IsZero() || p.state.IsZero() }

// Unlock transitions a locked lesson to in_progress, making it available to
// read. It is the only edge into in_progress: Gate decides WHEN a lesson
// becomes eligible, Unlock performs the transition once it is.
func (p LessonProgress) Unlock() (LessonProgress, error) {
	if !p.state.IsLocked() {
		return LessonProgress{}, fmt.Errorf("learning: lesson progress %q: unlock: %w", p.lessonRef.String(), ErrAlreadyUnlocked)
	}
	return LessonProgress{lessonRef: p.lessonRef, state: chunkStateInProgress}, nil
}

// MarkPassed transitions an in_progress lesson to passed once its recall
// checks are all graded correct. A locked lesson can never be marked
// passed — that invariant is the gating rule this bounded context exists
// to enforce.
func (p LessonProgress) MarkPassed() (LessonProgress, error) {
	switch {
	case p.state.IsLocked():
		return LessonProgress{}, fmt.Errorf("learning: lesson progress %q: mark passed: %w", p.lessonRef.String(), ErrLessonLocked)
	case p.state.IsPassed():
		return LessonProgress{}, fmt.Errorf("learning: lesson progress %q: mark passed: %w", p.lessonRef.String(), ErrAlreadyPassed)
	}
	return LessonProgress{lessonRef: p.lessonRef, state: chunkStatePassed}, nil
}
