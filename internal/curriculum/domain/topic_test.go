package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewTopic(t *testing.T) {
	validChapters := func(t *testing.T) []Chapter {
		return []Chapter{
			mustChapter(t, "modeling", "Modeling", 1, []Concept{
				mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1),
			}),
			mustChapter(t, "strategic-design", "Strategic Design", 2, []Concept{
				mustConcept(t, "bounded-context", "Bounded Context Concept Title", "Bounded Context Concept Outline Body", 1),
			}),
		}
	}

	tests := []struct {
		name     string
		track    Track
		slug     Slug
		title    string
		position Position
		chapters func(t *testing.T) []Chapter
		wantErr  error
	}{
		{
			name: "valid", track: mustTrack(t, "ddd"), slug: mustSlug(t, "tactical-patterns"),
			title: "Tactical Patterns", position: mustPosition(t, 1), chapters: validChapters,
		},
		{
			name: "title with surrounding whitespace trimmed", track: mustTrack(t, "ddd"), slug: mustSlug(t, "tactical-patterns"),
			title: "  Tactical Patterns  ", position: mustPosition(t, 1), chapters: validChapters,
		},
		{
			name: "zero-value track", track: Track{}, slug: mustSlug(t, "tactical-patterns"),
			title: "Tactical Patterns", position: mustPosition(t, 1), chapters: validChapters,
			wantErr: ErrInvalidTrack,
		},
		{
			name: "zero-value slug", track: mustTrack(t, "ddd"), slug: Slug{},
			title: "Tactical Patterns", position: mustPosition(t, 1), chapters: validChapters,
			wantErr: ErrInvalidSlug,
		},
		{
			name: "whitespace-only title", track: mustTrack(t, "ddd"), slug: mustSlug(t, "tactical-patterns"),
			title: "   ", position: mustPosition(t, 1), chapters: validChapters,
			wantErr: ErrInvalidTitle,
		},
		{
			name: "zero-value position", track: mustTrack(t, "ddd"), slug: mustSlug(t, "tactical-patterns"),
			title: "Tactical Patterns", position: Position{}, chapters: validChapters,
			wantErr: ErrInvalidPosition,
		},
		{
			name: "empty children", track: mustTrack(t, "ddd"), slug: mustSlug(t, "tactical-patterns"),
			title: "Tactical Patterns", position: mustPosition(t, 1),
			chapters: func(t *testing.T) []Chapter { return nil },
			wantErr:  ErrNoChildren,
		},
		{
			name: "single zero-value child", track: mustTrack(t, "ddd"), slug: mustSlug(t, "tactical-patterns"),
			title: "Tactical Patterns", position: mustPosition(t, 1),
			chapters: func(t *testing.T) []Chapter {
				return []Chapter{
					mustChapter(t, "modeling", "Modeling", 1, []Concept{mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1)}),
					{},
				}
			},
			wantErr: ErrZeroChild,
		},
		{
			name: "two zero-value children report zero-child not duplicate", track: mustTrack(t, "ddd"), slug: mustSlug(t, "tactical-patterns"),
			title: "Tactical Patterns", position: mustPosition(t, 1),
			chapters: func(t *testing.T) []Chapter { return []Chapter{{}, {}} },
			wantErr:  ErrZeroChild,
		},
		{
			name: "duplicate chapter slug", track: mustTrack(t, "ddd"), slug: mustSlug(t, "tactical-patterns"),
			title: "Tactical Patterns", position: mustPosition(t, 1),
			chapters: func(t *testing.T) []Chapter {
				return []Chapter{
					mustChapter(t, "modeling", "Modeling", 1, []Concept{mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1)}),
					mustChapter(t, "modeling", "Modeling Again", 2, []Concept{mustConcept(t, "entity", "Entity Concept Title", "Entity Concept Outline Body", 1)}),
				}
			},
			wantErr: ErrDuplicateSlug,
		},
		{
			name: "duplicate chapter position", track: mustTrack(t, "ddd"), slug: mustSlug(t, "tactical-patterns"),
			title: "Tactical Patterns", position: mustPosition(t, 1),
			chapters: func(t *testing.T) []Chapter {
				return []Chapter{
					mustChapter(t, "modeling", "Modeling", 1, []Concept{mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1)}),
					mustChapter(t, "strategic-design", "Strategic Design", 1, []Concept{mustConcept(t, "bounded-context", "Bounded Context Concept Title", "Bounded Context Concept Outline Body", 1)}),
				}
			},
			wantErr: ErrDuplicatePosition,
		},
		{
			name: "duplicate concept slug across chapters", track: mustTrack(t, "ddd"), slug: mustSlug(t, "tactical-patterns"),
			title: "Tactical Patterns", position: mustPosition(t, 1),
			chapters: func(t *testing.T) []Chapter {
				return []Chapter{
					mustChapter(t, "modeling", "Modeling", 1, []Concept{mustConcept(t, "idempotency", "Idempotency Concept Title", "Idempotency Concept Outline Body", 1)}),
					mustChapter(t, "strategic-design", "Strategic Design", 2, []Concept{mustConcept(t, "idempotency", "Idempotency Again Concept Title", "Idempotency Again Concept Outline Body", 1)}),
				}
			},
			wantErr: ErrDuplicateSlug,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chapters := tt.chapters(t)
			got, err := NewTopic(tt.track, tt.slug, tt.title, tt.position, chapters)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewTopic() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewTopic() unexpected error: %v", err)
			}

			if got.Track() != tt.track {
				t.Errorf("Track() = %v, want %v", got.Track(), tt.track)
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

			gotChapters := got.Chapters()
			if len(gotChapters) != len(chapters) {
				t.Fatalf("len(Chapters()) = %d, want %d", len(gotChapters), len(chapters))
			}
			for i, ch := range chapters {
				if gotChapters[i].Slug() != ch.Slug() {
					t.Errorf("Chapters()[%d].Slug() = %v, want %v (not in ascending position order)", i, gotChapters[i].Slug(), ch.Slug())
				}
			}
		})
	}
}

func TestNewTopicCopiesInputSlice(t *testing.T) {
	chapters := []Chapter{
		mustChapter(t, "modeling", "Modeling", 1, []Concept{mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1)}),
		mustChapter(t, "strategic-design", "Strategic Design", 2, []Concept{mustConcept(t, "bounded-context", "Bounded Context Concept Title", "Bounded Context Concept Outline Body", 1)}),
	}
	topic, err := NewTopic(mustTrack(t, "ddd"), mustSlug(t, "tactical-patterns"), "Tactical Patterns", mustPosition(t, 1), chapters)
	if err != nil {
		t.Fatalf("NewTopic() failed: %v", err)
	}

	chapters[0] = mustChapter(t, "mutated", "Mutated", 99, []Concept{mustConcept(t, "mutated-concept", "Mutated Concept Title", "Mutated Concept Outline Body", 1)})

	got := topic.Chapters()
	if got[0].Slug().String() != "modeling" {
		t.Errorf("Topic mutated after caller modified original slice: got slug %q, want %q", got[0].Slug().String(), "modeling")
	}
}

func TestTopicChaptersDefensiveCopyAtDepth(t *testing.T) {
	topic, err := NewTopic(mustTrack(t, "ddd"), mustSlug(t, "tactical-patterns"), "Tactical Patterns", mustPosition(t, 1), []Chapter{
		mustChapter(t, "modeling", "Modeling", 1, []Concept{mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1)}),
		mustChapter(t, "strategic-design", "Strategic Design", 2, []Concept{mustConcept(t, "bounded-context", "Bounded Context Concept Title", "Bounded Context Concept Outline Body", 1)}),
	})
	if err != nil {
		t.Fatalf("NewTopic() failed: %v", err)
	}

	got := topic.Chapters()
	got[0] = mustChapter(t, "mutated", "Mutated", 99, []Concept{mustConcept(t, "mutated-concept", "Mutated Concept Title", "Mutated Concept Outline Body", 1)})
	got = append(got, mustChapter(t, "appended", "Appended", 3, []Concept{mustConcept(t, "appended-concept", "Appended Concept Title", "Appended Concept Outline Body", 1)}))

	again := topic.Chapters()
	if len(again) != 2 {
		t.Fatalf("len(Chapters()) after external mutation = %d, want 2", len(again))
	}
	if again[0].Slug().String() != "modeling" {
		t.Errorf("Chapters()[0] mutated externally: got slug %q, want %q", again[0].Slug().String(), "modeling")
	}
	if again[1].Slug().String() != "strategic-design" {
		t.Errorf("Chapters()[1] mutated externally: got slug %q, want %q", again[1].Slug().String(), "strategic-design")
	}
}

func TestNewTopicSortsChaptersByPosition(t *testing.T) {
	topic, err := NewTopic(mustTrack(t, "ddd"), mustSlug(t, "tactical-patterns"), "Tactical Patterns", mustPosition(t, 1), []Chapter{
		mustChapter(t, "strategic-design", "Strategic Design", 2, []Concept{
			mustConcept(t, "bounded-context", "Bounded Context Concept Title", "Bounded Context Concept Outline Body", 1),
		}),
		mustChapter(t, "modeling", "Modeling", 1, []Concept{
			mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1),
		}),
	})
	if err != nil {
		t.Fatalf("NewTopic() failed: %v", err)
	}

	chapters := topic.Chapters()
	if len(chapters) != 2 {
		t.Fatalf("len(Chapters()) = %d, want 2", len(chapters))
	}
	if chapters[0].Slug().String() != "modeling" {
		t.Errorf("Chapters()[0].Slug() = %q, want %q (not sorted by position)", chapters[0].Slug().String(), "modeling")
	}
	if chapters[1].Slug().String() != "strategic-design" {
		t.Errorf("Chapters()[1].Slug() = %q, want %q (not sorted by position)", chapters[1].Slug().String(), "strategic-design")
	}

	concepts := chapters[0].Concepts()
	if len(concepts) != 1 {
		t.Fatalf("len(Chapters()[0].Concepts()) = %d, want 1", len(concepts))
	}
	if got, want := concepts[0].Slug().String(), "aggregate"; got != want {
		t.Errorf("Chapters()[0].Concepts()[0].Slug() = %q, want %q", got, want)
	}
	if got, want := concepts[0].Title(), "Aggregate Concept Title"; got != want {
		t.Errorf("Chapters()[0].Concepts()[0].Title() = %q, want %q", got, want)
	}
	if got, want := concepts[0].Outline(), "Aggregate Concept Outline Body"; got != want {
		t.Errorf("Chapters()[0].Concepts()[0].Outline() = %q, want %q", got, want)
	}
}

func TestTopicZeroValue(t *testing.T) {
	var top Topic
	if !top.IsZero() {
		t.Errorf("zero-value Topic.IsZero() = false, want true")
	}
}
