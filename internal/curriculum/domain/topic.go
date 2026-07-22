package domain

import "fmt"

// Topic is the curriculum aggregate root: a top-level subject within a
// Track, made of an ordered set of Chapters, identified by its Slug.
type Topic struct {
	track    Track
	slug     Slug
	title    string
	position Position
	chapters []Chapter
}

// NewTopic defensively copies chapters; the returned Topic is unaffected
// by later mutations to the slice.
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
	}

	copied := make([]Chapter, len(chapters))
	copy(copied, chapters)
	return Topic{track: track, slug: slug, title: trimmedTitle, position: position, chapters: copied}, nil
}

func (t Topic) Track() Track       { return t.track }
func (t Topic) Slug() Slug         { return t.slug }
func (t Topic) Title() string      { return t.title }
func (t Topic) Position() Position { return t.position }

// Chapters returns a copy of the topic's chapters; mutating the result
// does not affect the Topic.
func (t Topic) Chapters() []Chapter {
	out := make([]Chapter, len(t.chapters))
	copy(out, t.chapters)
	return out
}

// IsZero reports whether t was never constructed via NewTopic.
func (t Topic) IsZero() bool {
	return t.track.IsZero() || t.slug.IsZero() || t.position.IsZero() || t.title == ""
}
