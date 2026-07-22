package domain

import "fmt"

const maxSlugLen = 100

// Slug is a URL-safe natural key for a Topic, Chapter, or Concept: 1-100
// characters, lowercase alphanumerics and hyphens, must start and end
// alphanumeric, no consecutive hyphens.
type Slug struct {
	value string
}

func NewSlug(raw string) (Slug, error) {
	if !isValidSlug(raw) {
		return Slug{}, fmt.Errorf("curriculum: slug %q: %w", raw, ErrInvalidSlug)
	}
	return Slug{value: raw}, nil
}

func (s Slug) String() string { return s.value }

// IsZero reports whether s was never constructed via NewSlug.
func (s Slug) IsZero() bool { return s.value == "" }

func isValidSlug(raw string) bool {
	n := len(raw)
	if n == 0 || n > maxSlugLen {
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
