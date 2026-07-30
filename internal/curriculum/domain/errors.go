package domain

import "errors"

// Sentinel errors returned by this package's constructors. Callers match on
// these with errors.Is; the wrapping fmt.Errorf calls add positional context.
var (
	ErrInvalidSlug       = errors.New("invalid slug")
	ErrInvalidTrack      = errors.New("invalid track")
	ErrInvalidPosition   = errors.New("invalid position")
	ErrInvalidTitle      = errors.New("invalid title")
	ErrInvalidOutline    = errors.New("invalid outline")
	ErrNoChildren        = errors.New("no children")
	ErrZeroChild         = errors.New("zero-value child")
	ErrDuplicateSlug     = errors.New("duplicate child slug")
	ErrDuplicatePosition = errors.New("duplicate child position")

	ErrInvalidLessonVersion     = errors.New("invalid lesson version")
	ErrInvalidEstMinutes        = errors.New("invalid est minutes")
	ErrInvalidBodyMd            = errors.New("invalid body md")
	ErrInvalidReference         = errors.New("invalid reference")
	ErrInvalidReferenceCount    = errors.New("invalid reference count")
	ErrInvalidRecallKind        = errors.New("invalid recall check kind")
	ErrInvalidRecallQuestion    = errors.New("invalid recall check question")
	ErrInvalidRecallAnswer      = errors.New("invalid recall check expected answer")
	ErrInvalidRecallOptions     = errors.New("invalid recall check options")
	ErrInvalidRecallCheckCount  = errors.New("invalid recall check count")
	ErrInvalidRecallExplanation = errors.New("invalid recall check explanation")
)
