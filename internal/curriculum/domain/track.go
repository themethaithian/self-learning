package domain

import (
	"fmt"
	"slices"
)

// Track is one of the curriculum's fixed learning tracks.
type Track struct {
	value string
}

// tracks is the single source of truth for which tracks exist and their
// canonical order; NewTrack and Tracks both derive from it so the two can
// never drift apart.
var tracks = [...]Track{{value: "ddd"}, {value: "distsys"}, {value: "aws"}, {value: "go"}, {value: "dsa"}, {value: "ddia"}, {value: "ai-systems"}}

func NewTrack(raw string) (Track, error) {
	for _, t := range tracks {
		if t.value == raw {
			return t, nil
		}
	}
	return Track{}, fmt.Errorf("curriculum: track %q: %w", raw, ErrInvalidTrack)
}

func (t Track) String() string { return t.value }

// IsZero reports whether t was never constructed via NewTrack.
func (t Track) IsZero() bool { return t.value == "" }

// Tracks returns the curriculum tracks in canonical display order — not an
// incidental order, callers (e.g. grouping an API response) depend on it.
func Tracks() []Track { return slices.Clone(tracks[:]) }
