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

func mustEstMinutes(t *testing.T, raw int) EstMinutes {
	t.Helper()
	e, err := NewEstMinutes(raw)
	if err != nil {
		t.Fatalf("NewEstMinutes(%d) failed: %v", raw, err)
	}
	return e
}

func mustReference(t *testing.T, title, source, why string) Reference {
	t.Helper()
	r, err := NewReference(title, source, why)
	if err != nil {
		t.Fatalf("NewReference(%q) failed: %v", title, err)
	}
	return r
}

func mustRecallKind(t *testing.T, raw string) RecallKind {
	t.Helper()
	k, err := NewRecallKind(raw)
	if err != nil {
		t.Fatalf("NewRecallKind(%q) failed: %v", raw, err)
	}
	return k
}

func mustRecallCheck(t *testing.T, position int, kind string, question, expectedAnswer string, options []string) RecallCheck {
	t.Helper()
	return mustRecallCheckExplained(t, position, kind, question, expectedAnswer, options, "")
}

func mustRecallCheckExplained(t *testing.T, position int, kind string, question, expectedAnswer string, options []string, explanation string) RecallCheck {
	t.Helper()
	rc, err := NewRecallCheck(mustPosition(t, position), mustRecallKind(t, kind), question, expectedAnswer, options, explanation)
	if err != nil {
		t.Fatalf("NewRecallCheck(%d) failed: %v", position, err)
	}
	return rc
}

func validReferences(t *testing.T) []Reference {
	t.Helper()
	return []Reference{
		mustReference(t, "Evans ch. 4", "Domain-Driven Design, chapter 4", "core building blocks ของ aggregate"),
		mustReference(t, "Go blog: errors", "https://go.dev/blog/error-handling-and-go", "แนวทาง error wrapping มาตรฐาน"),
	}
}

func validRecallChecks(t *testing.T) []RecallCheck {
	t.Helper()
	return []RecallCheck{
		mustRecallCheck(t, 1, "short_answer", "Aggregate root คืออะไร?", "entity ที่เป็นทางเข้าเดียวของ aggregate", nil),
		mustRecallCheck(t, 2, "mcq", "ข้อใดคือ invariant ที่ถูกต้อง?", "b", []string{"a", "b", "c"}),
		mustRecallCheck(t, 3, "short_answer", "ทำไมต้อง encapsulate invariant?", "เพื่อป้องกัน state ที่ผิดกฎ", nil),
	}
}

func mustLesson(
	t *testing.T, slug string, version int, titleEn string, estMinutes int, bodyMd string,
	references []Reference, recallChecks []RecallCheck,
) Lesson {
	t.Helper()
	l, err := NewLesson(mustSlug(t, slug), version, titleEn, mustEstMinutes(t, estMinutes), bodyMd, references, recallChecks)
	if err != nil {
		t.Fatalf("NewLesson(%q) failed: %v", slug, err)
	}
	return l
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

// TestMustReferenceRoles guards mustReference's argument order the same way
// TestMustConceptRoles guards mustConcept's: three interchangeable strings
// (title/source/why) are easy to silently swap inside the helper.
func TestMustReferenceRoles(t *testing.T) {
	r := mustReference(t, "Reference Title", "Reference Source", "Reference Why")
	if r.Title() != "Reference Title" {
		t.Errorf("mustReference Title() = %q, want %q", r.Title(), "Reference Title")
	}
	if r.Source() != "Reference Source" {
		t.Errorf("mustReference Source() = %q, want %q", r.Source(), "Reference Source")
	}
	if r.Why() != "Reference Why" {
		t.Errorf("mustReference Why() = %q, want %q", r.Why(), "Reference Why")
	}
}

// TestMustRecallCheckRoles guards mustRecallCheck's argument order the same
// way TestMustConceptRoles guards mustConcept's: question and expectedAnswer
// are interchangeable strings easy to silently swap inside the helper.
func TestMustRecallCheckRoles(t *testing.T) {
	rc := mustRecallCheck(t, 1, "short_answer", "Recall Check Question", "Recall Check Expected Answer", nil)
	if rc.Question() != "Recall Check Question" {
		t.Errorf("mustRecallCheck Question() = %q, want %q", rc.Question(), "Recall Check Question")
	}
	if rc.ExpectedAnswer() != "Recall Check Expected Answer" {
		t.Errorf("mustRecallCheck ExpectedAnswer() = %q, want %q", rc.ExpectedAnswer(), "Recall Check Expected Answer")
	}
}
