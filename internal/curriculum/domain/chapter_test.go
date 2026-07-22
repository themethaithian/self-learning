package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewChapter(t *testing.T) {
	tests := []struct {
		name     string
		slug     Slug
		title    string
		position Position
		concepts []Concept
		wantErr  error
	}{
		{
			name: "valid", slug: mustSlug(t, "modeling"), title: "Modeling", position: mustPosition(t, 1),
			concepts: []Concept{
				mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1),
				mustConcept(t, "entity", "Entity Concept Title", "Entity Concept Outline Body", 2),
			},
		},
		{
			name: "title with surrounding whitespace trimmed", slug: mustSlug(t, "modeling"),
			title: "  Modeling  ", position: mustPosition(t, 1),
			concepts: []Concept{
				mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1),
			},
		},
		{
			name: "zero-value slug", slug: Slug{}, title: "Modeling", position: mustPosition(t, 1),
			concepts: []Concept{mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1)},
			wantErr:  ErrInvalidSlug,
		},
		{
			name: "whitespace-only title", slug: mustSlug(t, "modeling"), title: "   ", position: mustPosition(t, 1),
			concepts: []Concept{mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1)},
			wantErr:  ErrInvalidTitle,
		},
		{
			name: "zero-value position", slug: mustSlug(t, "modeling"), title: "Modeling", position: Position{},
			concepts: []Concept{mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1)},
			wantErr:  ErrInvalidPosition,
		},
		{
			name: "empty children", slug: mustSlug(t, "modeling"), title: "Modeling", position: mustPosition(t, 1),
			concepts: nil, wantErr: ErrNoChildren,
		},
		{
			name: "single zero-value child", slug: mustSlug(t, "modeling"), title: "Modeling", position: mustPosition(t, 1),
			concepts: []Concept{
				mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1),
				{},
			},
			wantErr: ErrZeroChild,
		},
		{
			name: "two zero-value children report zero-child not duplicate", slug: mustSlug(t, "modeling"),
			title: "Modeling", position: mustPosition(t, 1),
			concepts: []Concept{{}, {}},
			wantErr:  ErrZeroChild,
		},
		{
			name: "duplicate concept slug", slug: mustSlug(t, "modeling"), title: "Modeling", position: mustPosition(t, 1),
			concepts: []Concept{
				mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1),
				mustConcept(t, "aggregate", "Aggregate Again Title", "Aggregate Again Outline Body", 2),
			},
			wantErr: ErrDuplicateSlug,
		},
		{
			name: "duplicate concept position", slug: mustSlug(t, "modeling"), title: "Modeling", position: mustPosition(t, 1),
			concepts: []Concept{
				mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1),
				mustConcept(t, "entity", "Entity Concept Title", "Entity Concept Outline Body", 1),
			},
			wantErr: ErrDuplicatePosition,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChapter(tt.slug, tt.title, tt.position, tt.concepts)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewChapter() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewChapter() unexpected error: %v", err)
			}

			if got.Slug() != tt.slug {
				t.Errorf("Slug() = %v, want %v", got.Slug(), tt.slug)
			}
			if want := strings.TrimSpace(tt.title); got.Title() != want {
				t.Errorf("Title() = %q, want %q", got.Title(), want)
			}
			if got.Position() != tt.position {
				t.Errorf("Position() = %v, want %v", got.Position(), tt.position)
			}

			concepts := got.Concepts()
			if len(concepts) != len(tt.concepts) {
				t.Fatalf("len(Concepts()) = %d, want %d", len(concepts), len(tt.concepts))
			}
			for i, c := range tt.concepts {
				if concepts[i].Slug() != c.Slug() {
					t.Errorf("Concepts()[%d].Slug() = %v, want %v (not in ascending position order)", i, concepts[i].Slug(), c.Slug())
				}
			}
		})
	}
}

func TestNewChapterCopiesInputSlice(t *testing.T) {
	concepts := []Concept{
		mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1),
		mustConcept(t, "entity", "Entity Concept Title", "Entity Concept Outline Body", 2),
	}
	ch := mustChapter(t, "modeling", "Modeling", 1, concepts)

	concepts[0] = mustConcept(t, "mutated", "Mutated Concept Title", "Mutated Concept Outline Body", 99)

	got := ch.Concepts()
	if got[0].Slug().String() != "aggregate" {
		t.Errorf("Chapter mutated after caller modified original slice: got slug %q, want %q", got[0].Slug().String(), "aggregate")
	}
}

func TestChapterConceptsDefensiveCopyAtDepth(t *testing.T) {
	ch := mustChapter(t, "modeling", "Modeling", 1, []Concept{
		mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1),
		mustConcept(t, "entity", "Entity Concept Title", "Entity Concept Outline Body", 2),
	})

	got := ch.Concepts()
	got[0] = mustConcept(t, "mutated", "Mutated Concept Title", "Mutated Concept Outline Body", 99)
	got = append(got, mustConcept(t, "appended", "Appended Concept Title", "Appended Concept Outline Body", 3))

	again := ch.Concepts()
	if len(again) != 2 {
		t.Fatalf("len(Concepts()) after external mutation = %d, want 2", len(again))
	}
	if again[0].Slug().String() != "aggregate" {
		t.Errorf("Concepts()[0] mutated externally: got slug %q, want %q", again[0].Slug().String(), "aggregate")
	}
	if again[1].Slug().String() != "entity" {
		t.Errorf("Concepts()[1] mutated externally: got slug %q, want %q", again[1].Slug().String(), "entity")
	}
}

func TestNewChapterSortsConceptsByPosition(t *testing.T) {
	ch := mustChapter(t, "modeling", "Modeling", 1, []Concept{
		mustConcept(t, "entity", "Entity Concept Title", "Entity Concept Outline Body", 2),
		mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1),
	})

	concepts := ch.Concepts()
	if len(concepts) != 2 {
		t.Fatalf("len(Concepts()) = %d, want 2", len(concepts))
	}
	if concepts[0].Slug().String() != "aggregate" {
		t.Errorf("Concepts()[0].Slug() = %q, want %q (not sorted by position)", concepts[0].Slug().String(), "aggregate")
	}
	if concepts[1].Slug().String() != "entity" {
		t.Errorf("Concepts()[1].Slug() = %q, want %q (not sorted by position)", concepts[1].Slug().String(), "entity")
	}
}

func TestChapterZeroValue(t *testing.T) {
	var ch Chapter
	if !ch.IsZero() {
		t.Errorf("zero-value Chapter.IsZero() = false, want true")
	}
}
