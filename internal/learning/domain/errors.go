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

	ErrInvalidCheckKey       = errors.New("invalid check key")
	ErrInvalidConfidence     = errors.New("invalid confidence")
	ErrInvalidAttemptOutcome = errors.New("invalid attempt outcome")
	ErrInvalidCheckKind      = errors.New("invalid check kind")
	ErrInvalidGradedBy       = errors.New("invalid graded by")
	ErrSelectedOptionNotMCQ  = errors.New("selected option only valid for mcq checks")

	ErrInvalidReviewQuality = errors.New("invalid review quality")
	ErrInvalidEaseFactor    = errors.New("invalid ease factor")
	ErrInvalidReviewCard    = errors.New("invalid review card")
)
