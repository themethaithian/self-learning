package domain

import "errors"

// Sentinel errors returned by this package's constructors and transition
// methods. Callers match on these with errors.Is; the wrapping fmt.Errorf
// calls add identity context.
var (
	ErrInvalidLessonRef  = errors.New("invalid lesson ref")
	ErrInvalidChunkState = errors.New("invalid chunk state")
	ErrAlreadyUnlocked   = errors.New("lesson already unlocked")
	ErrLessonLocked      = errors.New("lesson is locked")
	ErrAlreadyPassed     = errors.New("lesson already passed")
)
