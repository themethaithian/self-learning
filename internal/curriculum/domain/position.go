package domain

import (
	"fmt"
	"math"
)

// Position is a 1-based ordering index among sibling entities.
type Position struct {
	value int
}

func NewPosition(raw int) (Position, error) {
	// Upper bound matches the INT column position is stored in.
	if raw < 1 || raw > math.MaxInt32 {
		return Position{}, fmt.Errorf("curriculum: position %d: %w", raw, ErrInvalidPosition)
	}
	return Position{value: raw}, nil
}

func (p Position) Int() int { return p.value }

// IsZero reports whether p was never constructed via NewPosition.
func (p Position) IsZero() bool { return p.value == 0 }
