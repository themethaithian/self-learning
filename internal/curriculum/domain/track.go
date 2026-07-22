package domain

import "fmt"

// Track is one of the five fixed learning tracks in the curriculum.
type Track struct {
	value string
}

func NewTrack(raw string) (Track, error) {
	switch raw {
	case "ddd", "distsys", "aws", "go", "dsa":
		return Track{value: raw}, nil
	default:
		return Track{}, fmt.Errorf("curriculum: track %q: %w", raw, ErrInvalidTrack)
	}
}

func (t Track) String() string { return t.value }

// IsZero reports whether t was never constructed via NewTrack.
func (t Track) IsZero() bool { return t.value == "" }
