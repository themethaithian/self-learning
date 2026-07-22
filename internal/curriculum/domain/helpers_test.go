package domain

import "testing"

func mustSlug(t *testing.T, raw string) Slug {
	t.Helper()
	s, err := NewSlug(raw)
	if err != nil {
		t.Fatalf("NewSlug(%q) failed: %v", raw, err)
	}
	return s
}

func mustTrack(t *testing.T, raw string) Track {
	t.Helper()
	tr, err := NewTrack(raw)
	if err != nil {
		t.Fatalf("NewTrack(%q) failed: %v", raw, err)
	}
	return tr
}

func mustPosition(t *testing.T, raw int) Position {
	t.Helper()
	p, err := NewPosition(raw)
	if err != nil {
		t.Fatalf("NewPosition(%d) failed: %v", raw, err)
	}
	return p
}

func mustConcept(t *testing.T, slug, title, outline string, position int) Concept {
	t.Helper()
	c, err := NewConcept(mustSlug(t, slug), title, outline, mustPosition(t, position))
	if err != nil {
		t.Fatalf("NewConcept(%q) failed: %v", slug, err)
	}
	return c
}

func mustChapter(t *testing.T, slug, title string, position int, concepts []Concept) Chapter {
	t.Helper()
	ch, err := NewChapter(mustSlug(t, slug), title, mustPosition(t, position), concepts)
	if err != nil {
		t.Fatalf("NewChapter(%q) failed: %v", slug, err)
	}
	return ch
}

// TestMustConceptRoles guards mustConcept's own argument order against literal
// expected strings — not against another mustConcept call — so a title/outline
// swap inside mustConcept fails here even though every other test in this
// package builds its expectations by calling mustConcept the same (buggy) way.
func TestMustConceptRoles(t *testing.T) {
	c := mustConcept(t, "aggregate", "Aggregate Concept Title", "Aggregate Concept Outline Body", 1)
	if c.Title() != "Aggregate Concept Title" {
		t.Errorf("mustConcept Title() = %q, want %q", c.Title(), "Aggregate Concept Title")
	}
	if c.Outline() != "Aggregate Concept Outline Body" {
		t.Errorf("mustConcept Outline() = %q, want %q", c.Outline(), "Aggregate Concept Outline Body")
	}
}
