// Package domain holds the prefs bounded context's business rules: pure Go,
// no HTTP or database types.
package domain

import (
	"errors"
	"fmt"
)

var ErrInvalidFocusTrack = errors.New("invalid focus track")

// maxFocusTrackLen is a generous, arbitrary bound — not derived from the
// prefs.value column width (VARCHAR(255)) — so it's not "more correct" to
// raise it to 255 later.
const maxFocusTrackLen = 64

// FocusTrack is the slug of the curriculum track /today currently focuses
// on. It only enforces slug shape, not that the track exists — see
// app.Service.SetFocusTrack for why.
type FocusTrack struct {
	value string
}

func NewFocusTrack(raw string) (FocusTrack, error) {
	if !isValidFocusTrack(raw) {
		return FocusTrack{}, fmt.Errorf("prefs: focus track %q: %w", raw, ErrInvalidFocusTrack)
	}
	return FocusTrack{value: raw}, nil
}

func (f FocusTrack) String() string { return f.value }

func isValidFocusTrack(raw string) bool {
	n := len(raw)
	if n == 0 || n > maxFocusTrackLen {
		return false
	}
	for i := 0; i < n; i++ {
		c := raw[i]
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '-' {
			return false
		}
	}
	return true
}
