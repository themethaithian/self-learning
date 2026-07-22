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
)
