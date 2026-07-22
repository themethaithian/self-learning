package domain

import (
	"cmp"
	"fmt"
	"slices"
)

// Chapter groups an ordered set of Concepts within a Topic, identified by
// its Slug.
type Chapter struct {
	slug     Slug
	title    string
	position Position
	concepts []Concept
}

// NewChapter constructs a validated Chapter from slug, title, position, and
// concepts; it defensively copies concepts, so mutating the slice afterwards
// does not affect the result.
func NewChapter(slug Slug, title string, position Position, concepts []Concept) (Chapter, error) {
	if slug.IsZero() {
		return Chapter{}, fmt.Errorf("curriculum: chapter: slug: %w", ErrInvalidSlug)
	}
	trimmedTitle, err := validateTitle(title)
	if err != nil {
		return Chapter{}, fmt.Errorf("curriculum: chapter %q: title: %w", slug.String(), err)
	}
	if position.IsZero() {
		return Chapter{}, fmt.Errorf("curriculum: chapter %q: position: %w", slug.String(), ErrInvalidPosition)
	}
	if len(concepts) == 0 {
		return Chapter{}, fmt.Errorf("curriculum: chapter %q: %w", slug.String(), ErrNoChildren)
	}

	seenSlugs := make(map[string]struct{}, len(concepts))
	seenPositions := make(map[int]struct{}, len(concepts))
	for _, c := range concepts {
		if c.IsZero() {
			return Chapter{}, fmt.Errorf("curriculum: chapter %q: %w", slug.String(), ErrZeroChild)
		}

		key := c.slug.String()
		if _, dup := seenSlugs[key]; dup {
			return Chapter{}, fmt.Errorf("curriculum: chapter %q: concept slug %q: %w", slug.String(), key, ErrDuplicateSlug)
		}
		seenSlugs[key] = struct{}{}

		pos := c.position.Int()
		if _, dup := seenPositions[pos]; dup {
			return Chapter{}, fmt.Errorf("curriculum: chapter %q: concept position %d: %w", slug.String(), pos, ErrDuplicatePosition)
		}
		seenPositions[pos] = struct{}{}
	}

	copied := make([]Concept, len(concepts))
	copy(copied, concepts)
	slices.SortFunc(copied, func(a, b Concept) int {
		return cmp.Compare(a.position.Int(), b.position.Int())
	})
	return Chapter{slug: slug, title: trimmedTitle, position: position, concepts: copied}, nil
}

func (ch Chapter) Slug() Slug         { return ch.slug }
func (ch Chapter) Title() string      { return ch.title }
func (ch Chapter) Position() Position { return ch.position }

// Concepts returns a copy of the chapter's concepts in ascending position
// order; mutating the result does not affect the Chapter.
func (ch Chapter) Concepts() []Concept {
	out := make([]Concept, len(ch.concepts))
	copy(out, ch.concepts)
	return out
}

// IsZero reports whether ch was never constructed via NewChapter.
func (ch Chapter) IsZero() bool {
	return ch.slug.IsZero() || ch.position.IsZero() || ch.title == ""
}
