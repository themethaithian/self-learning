package domain

import "fmt"

const (
	minEstMinutes = 5
	maxEstMinutes = 10
)

// EstMinutes is a Lesson's estimated reading time in minutes, constrained to
// the content pipeline's 5-10 minute chunk size (CLAUDE.md content model).
type EstMinutes struct {
	value int
}

func NewEstMinutes(raw int) (EstMinutes, error) {
	if raw < minEstMinutes || raw > maxEstMinutes {
		return EstMinutes{}, fmt.Errorf("curriculum: est minutes %d: %w", raw, ErrInvalidEstMinutes)
	}
	return EstMinutes{value: raw}, nil
}

func (e EstMinutes) Int() int { return e.value }

// IsZero reports whether e was never constructed via NewEstMinutes.
func (e EstMinutes) IsZero() bool { return e.value == 0 }
