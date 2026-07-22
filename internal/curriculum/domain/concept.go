package domain

import "fmt"

// Concept is the leaf of the curriculum tree: one lesson topic within a
// Chapter, identified by its Slug.
type Concept struct {
	slug     Slug
	title    string
	outline  string
	position Position
}

func NewConcept(slug Slug, title string, outline string, position Position) (Concept, error) {
	if slug.IsZero() {
		return Concept{}, fmt.Errorf("curriculum: concept: slug: %w", ErrInvalidSlug)
	}
	trimmedTitle, err := validateTitle(title)
	if err != nil {
		return Concept{}, fmt.Errorf("curriculum: concept %q: title: %w", slug.String(), err)
	}
	trimmedOutline, err := validateOutline(outline)
	if err != nil {
		return Concept{}, fmt.Errorf("curriculum: concept %q: outline: %w", slug.String(), err)
	}
	if position.IsZero() {
		return Concept{}, fmt.Errorf("curriculum: concept %q: position: %w", slug.String(), ErrInvalidPosition)
	}
	return Concept{slug: slug, title: trimmedTitle, outline: trimmedOutline, position: position}, nil
}

func (c Concept) Slug() Slug         { return c.slug }
func (c Concept) Title() string      { return c.title }
func (c Concept) Outline() string    { return c.outline }
func (c Concept) Position() Position { return c.position }

// IsZero reports whether c was never constructed via NewConcept.
func (c Concept) IsZero() bool {
	return c.slug.IsZero() || c.position.IsZero() || c.title == ""
}
