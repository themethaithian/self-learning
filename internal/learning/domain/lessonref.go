package domain

import "fmt"

const maxLessonRefLen = 100

// LessonRef identifies which Lesson a LessonProgress tracks, by the
// underlying concept slug. It is a small VO local to this bounded context
// rather than curriculum.Slug imported directly: the two contexts talk only
// through shared ID values and events (see docs/design.md §1), never by
// importing each other's domain types, so learning stays free to evolve its
// own identity story (e.g. keying off a numeric lesson id later) without a
// compile-time dependency on curriculum's package.
type LessonRef struct {
	value string
}

func NewLessonRef(raw string) (LessonRef, error) {
	if !isValidLessonRef(raw) {
		return LessonRef{}, fmt.Errorf("learning: lesson ref %q: %w", raw, ErrInvalidLessonRef)
	}
	return LessonRef{value: raw}, nil
}

func (r LessonRef) String() string { return r.value }

// IsZero reports whether r was never constructed via NewLessonRef.
func (r LessonRef) IsZero() bool { return r.value == "" }

// isValidLessonRef mirrors curriculum.Slug's shape rules.
func isValidLessonRef(raw string) bool {
	n := len(raw)
	if n == 0 || n > maxLessonRefLen {
		return false
	}
	for i := 0; i < n; i++ {
		c := raw[i]
		switch {
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9':
			continue
		case c == '-':
			if i == 0 || i == n-1 || raw[i-1] == '-' {
				return false
			}
		default:
			return false
		}
	}
	return true
}
