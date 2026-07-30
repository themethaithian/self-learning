package domain

import (
	"cmp"
	"fmt"
	"slices"
)

const (
	minReferences = 2
	maxReferences = 4

	// minRecallChecks stays at 3 — AWS-S1 raised only the ceiling (the
	// blocker was 5 x 50 concepts = 250 against a 550-question target). The
	// floor has never actually been binding: counted across all 118 existing
	// lessons, the lowest is 4 recall checks, none has 3, so this floor is
	// not close to constraining any real content today, and nothing in this
	// ticket's scope calls for changing that contract regardless.
	minRecallChecks = 3
	// maxRecallChecks 15 gives headroom above the AWS plan's ~11 checks per
	// concept (docs/tickets/aws-cert.md) without opening the door to an
	// unbounded per-lesson quiz.
	maxRecallChecks = 15
)

// Lesson is the curriculum's teaching content for one Concept: a self-study
// chunk of Markdown body, estimated reading time, recall checks, and
// references. It is identified by the concept's Slug plus a revision
// Version, so re-generating a concept's content produces a new Lesson
// rather than mutating the one already imported.
type Lesson struct {
	slug         Slug
	version      int
	titleEn      string
	estMinutes   EstMinutes
	bodyMd       string
	references   []Reference
	recallChecks []RecallCheck
}

// NewLesson constructs a validated Lesson from a concept slug, a revision
// version, and its content. It defensively copies references and
// recallChecks and sorts recallChecks by Position, so mutating the input
// slices afterwards does not affect the result.
func NewLesson(
	slug Slug,
	version int,
	titleEn string,
	estMinutes EstMinutes,
	bodyMd string,
	references []Reference,
	recallChecks []RecallCheck,
) (Lesson, error) {
	if slug.IsZero() {
		return Lesson{}, fmt.Errorf("curriculum: lesson: slug: %w", ErrInvalidSlug)
	}
	if version < 1 {
		return Lesson{}, fmt.Errorf("curriculum: lesson %q: version %d: %w", slug.String(), version, ErrInvalidLessonVersion)
	}
	trimmedTitle, err := validateTitle(titleEn)
	if err != nil {
		return Lesson{}, fmt.Errorf("curriculum: lesson %q: title: %w", slug.String(), err)
	}
	if estMinutes.IsZero() {
		return Lesson{}, fmt.Errorf("curriculum: lesson %q: est minutes: %w", slug.String(), ErrInvalidEstMinutes)
	}
	trimmedBody, err := validateBodyMd(bodyMd)
	if err != nil {
		return Lesson{}, fmt.Errorf("curriculum: lesson %q: body: %w", slug.String(), err)
	}

	if len(references) < minReferences || len(references) > maxReferences {
		return Lesson{}, fmt.Errorf("curriculum: lesson %q: %d references: %w", slug.String(), len(references), ErrInvalidReferenceCount)
	}
	for _, r := range references {
		if r.IsZero() {
			return Lesson{}, fmt.Errorf("curriculum: lesson %q: reference: %w", slug.String(), ErrZeroChild)
		}
	}

	if len(recallChecks) < minRecallChecks || len(recallChecks) > maxRecallChecks {
		return Lesson{}, fmt.Errorf("curriculum: lesson %q: %d recall checks: %w", slug.String(), len(recallChecks), ErrInvalidRecallCheckCount)
	}
	seenPositions := make(map[int]struct{}, len(recallChecks))
	for _, rc := range recallChecks {
		if rc.IsZero() {
			return Lesson{}, fmt.Errorf("curriculum: lesson %q: recall check: %w", slug.String(), ErrZeroChild)
		}
		pos := rc.position.Int()
		if _, dup := seenPositions[pos]; dup {
			return Lesson{}, fmt.Errorf("curriculum: lesson %q: recall check position %d: %w", slug.String(), pos, ErrDuplicatePosition)
		}
		seenPositions[pos] = struct{}{}
	}

	copiedRefs := make([]Reference, len(references))
	copy(copiedRefs, references)

	copiedChecks := make([]RecallCheck, len(recallChecks))
	copy(copiedChecks, recallChecks)
	slices.SortFunc(copiedChecks, func(a, b RecallCheck) int {
		return cmp.Compare(a.position.Int(), b.position.Int())
	})

	return Lesson{
		slug: slug, version: version, titleEn: trimmedTitle, estMinutes: estMinutes,
		bodyMd: trimmedBody, references: copiedRefs, recallChecks: copiedChecks,
	}, nil
}

func (l Lesson) Slug() Slug             { return l.slug }
func (l Lesson) Version() int           { return l.version }
func (l Lesson) TitleEn() string        { return l.titleEn }
func (l Lesson) EstMinutes() EstMinutes { return l.estMinutes }
func (l Lesson) BodyMd() string         { return l.bodyMd }

// References returns a copy of the lesson's references; mutating the result
// does not affect the Lesson.
func (l Lesson) References() []Reference {
	out := make([]Reference, len(l.references))
	copy(out, l.references)
	return out
}

// RecallChecks returns a copy of the lesson's recall checks in ascending
// position order; mutating the result does not affect the Lesson.
func (l Lesson) RecallChecks() []RecallCheck {
	out := make([]RecallCheck, len(l.recallChecks))
	copy(out, l.recallChecks)
	return out
}

// IsZero reports whether l was never constructed via NewLesson.
func (l Lesson) IsZero() bool {
	return l.slug.IsZero() || l.version == 0 || l.titleEn == ""
}
