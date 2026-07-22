package domain

import (
	"cmp"
	"fmt"
	"slices"
)

// Topic is the curriculum aggregate root: a top-level subject within a
// Track, made of an ordered set of Chapters, identified by its Slug.
type Topic struct {
	track    Track
	slug     Slug
	title    string
	position Position
	chapters []Chapter
}

// NewTopic constructs a validated Topic from track, slug, title, position,
// and chapters; it defensively copies chapters, so mutating the slice
// afterwards does not affect the result.
func NewTopic(track Track, slug Slug, title string, position Position, chapters []Chapter) (Topic, error) {
	if track.IsZero() {
		return Topic{}, fmt.Errorf("curriculum: topic: track: %w", ErrInvalidTrack)
	}
	if slug.IsZero() {
		return Topic{}, fmt.Errorf("curriculum: topic: slug: %w", ErrInvalidSlug)
	}
	trimmedTitle, err := validateTitle(title)
	if err != nil {
		return Topic{}, fmt.Errorf("curriculum: topic %q: title: %w", slug.String(), err)
	}
	if position.IsZero() {
		return Topic{}, fmt.Errorf("curriculum: topic %q: position: %w", slug.String(), ErrInvalidPosition)
	}
	if len(chapters) == 0 {
		return Topic{}, fmt.Errorf("curriculum: topic %q: %w", slug.String(), ErrNoChildren)
	}

	seenSlugs := make(map[string]struct{}, len(chapters))
	seenPositions := make(map[int]struct{}, len(chapters))
	// Concept slugs are unique per chapter (NewChapter's own invariant) but
	// must also be unique across the whole topic: lessons are stored at
	// content/lessons/<topic>/<concept>.json, a two-level path with no
	// chapter segment, so only the aggregate root can catch a collision.
	seenConceptSlugs := make(map[string]struct{})
	for _, ch := range chapters {
		if ch.IsZero() {
			return Topic{}, fmt.Errorf("curriculum: topic %q: %w", slug.String(), ErrZeroChild)
		}

		key := ch.slug.String()
		if _, dup := seenSlugs[key]; dup {
			return Topic{}, fmt.Errorf("curriculum: topic %q: chapter slug %q: %w", slug.String(), key, ErrDuplicateSlug)
		}
		seenSlugs[key] = struct{}{}

		pos := ch.position.Int()
		if _, dup := seenPositions[pos]; dup {
			return Topic{}, fmt.Errorf("curriculum: topic %q: chapter position %d: %w", slug.String(), pos, ErrDuplicatePosition)
		}
		seenPositions[pos] = struct{}{}

		for _, c := range ch.concepts {
			conceptKey := c.slug.String()
			if _, dup := seenConceptSlugs[conceptKey]; dup {
				return Topic{}, fmt.Errorf("curriculum: topic %q: chapter %q: concept slug %q: %w", slug.String(), key, conceptKey, ErrDuplicateSlug)
			}
			seenConceptSlugs[conceptKey] = struct{}{}
		}
	}

	copied := make([]Chapter, len(chapters))
	copy(copied, chapters)
	slices.SortFunc(copied, func(a, b Chapter) int {
		return cmp.Compare(a.position.Int(), b.position.Int())
	})
	return Topic{track: track, slug: slug, title: trimmedTitle, position: position, chapters: copied}, nil
}

func (t Topic) Track() Track       { return t.track }
func (t Topic) Slug() Slug         { return t.slug }
func (t Topic) Title() string      { return t.title }
func (t Topic) Position() Position { return t.position }

// Chapters returns a copy of the topic's chapters in ascending position
// order; mutating the result does not affect the Topic.
func (t Topic) Chapters() []Chapter {
	out := make([]Chapter, len(t.chapters))
	copy(out, t.chapters)
	return out
}

// IsZero reports whether t was never constructed via NewTopic.
func (t Topic) IsZero() bool {
	return t.track.IsZero() || t.slug.IsZero() || t.position.IsZero() || t.title == ""
}
