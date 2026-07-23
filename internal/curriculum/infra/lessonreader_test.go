package infra

import (
	"context"
	"errors"
	"strings"
	"testing"

	curriculumapp "github.com/themethaithian/self-learning/internal/curriculum/app"
)

// TestRepositoryLessonByConcept_HappyPath reads back exactly what
// TestRepositorySaveLesson_HappyPathCommitsOnce wrote, through the real
// relational stub rather than a canned row — proving the join, not just the
// Scan. It also pins that expected_answer and mcq options reach the domain
// value: the design decision this ticket implements (self-grade needs them).
func TestRepositoryLessonByConcept_HappyPath(t *testing.T) {
	result := seededConceptResult("domain-driven-design", "aggregate", 42)
	db := openStubDB(t, result)
	repo := NewRepository(db)
	ctx := context.Background()

	want := newTestLesson(t, "aggregate", 1)
	if _, err := repo.SaveLesson(ctx, "domain-driven-design", want); err != nil {
		t.Fatalf("SaveLesson() unexpected error: %v", err)
	}

	got, err := repo.LessonByConcept(ctx, "domain-driven-design", "aggregate")
	if err != nil {
		t.Fatalf("LessonByConcept() unexpected error: %v", err)
	}

	if got.Slug() != want.Slug() || got.Version() != want.Version() || got.TitleEn() != want.TitleEn() ||
		got.EstMinutes() != want.EstMinutes() || got.BodyMd() != want.BodyMd() {
		t.Fatalf("LessonByConcept() = %+v, want fields matching %+v", got, want)
	}

	gotRefs, wantRefs := got.References(), want.References()
	if len(gotRefs) != len(wantRefs) {
		t.Fatalf("References() = %d, want %d", len(gotRefs), len(wantRefs))
	}
	for i := range wantRefs {
		if gotRefs[i] != wantRefs[i] {
			t.Errorf("References()[%d] = %+v, want %+v", i, gotRefs[i], wantRefs[i])
		}
	}

	gotChecks, wantChecks := got.RecallChecks(), want.RecallChecks()
	if len(gotChecks) != len(wantChecks) {
		t.Fatalf("RecallChecks() = %d, want %d", len(gotChecks), len(wantChecks))
	}
	for i := range wantChecks {
		g, w := gotChecks[i], wantChecks[i]
		if g.Position() != w.Position() || g.Kind() != w.Kind() || g.Question() != w.Question() || g.ExpectedAnswer() != w.ExpectedAnswer() {
			t.Errorf("RecallChecks()[%d] = %+v, want %+v", i, g, w)
		}
		if gotOpts, wantOpts := g.Options(), w.Options(); len(gotOpts) != len(wantOpts) {
			t.Errorf("RecallChecks()[%d].Options() = %v, want %v", i, gotOpts, wantOpts)
		} else {
			for j := range wantOpts {
				if gotOpts[j] != wantOpts[j] {
					t.Errorf("RecallChecks()[%d].Options()[%d] = %q, want %q", i, j, gotOpts[j], wantOpts[j])
				}
			}
		}
	}
}

func TestRepositoryLessonByConcept_NoLessonForConcept(t *testing.T) {
	result := seededConceptResult("domain-driven-design", "aggregate", 42) // concept exists, no lesson saved
	db := openStubDB(t, result)
	repo := NewRepository(db)

	_, err := repo.LessonByConcept(context.Background(), "domain-driven-design", "aggregate")
	if !errors.Is(err, curriculumapp.ErrLessonNotFound) {
		t.Fatalf("LessonByConcept() error = %v, want it to wrap ErrLessonNotFound", err)
	}
}

// TestRepositoryLessonByConcept_ScopedByBothSlugs is the concept-slug
// analogue of TestRepositorySaveLesson_ResolvesConceptScopedToTopic: a
// concept slug shared by two topics must not leak the wrong topic's lesson
// when queried under the other topic. Note this does NOT catch a dropped
// "t.slug = ?" predicate in selectLessonByConceptSQL itself — the stub
// dispatches on query-string identity, so an edit to the SQL text is
// invisible to this behavioral test either way. TestSelectLessonByConceptSQLShape
// is the actual guard for the predicate; this test only guards the stub's
// (and by extension resolveConceptID's own) lookup-key semantics.
func TestRepositoryLessonByConcept_ScopedByBothSlugs(t *testing.T) {
	result := &stubResult{}
	result.seedConcept("topic-a", "shared-slug", 101)
	result.seedConcept("topic-b", "shared-slug", 202)
	db := openStubDB(t, result)
	repo := NewRepository(db)
	ctx := context.Background()

	if _, err := repo.SaveLesson(ctx, "topic-a", newTestLesson(t, "shared-slug", 1)); err != nil {
		t.Fatalf("SaveLesson(topic-a) unexpected error: %v", err)
	}

	if _, err := repo.LessonByConcept(ctx, "topic-b", "shared-slug"); !errors.Is(err, curriculumapp.ErrLessonNotFound) {
		t.Fatalf("LessonByConcept(topic-b) error = %v, want ErrLessonNotFound (topic-b never got a lesson)", err)
	}

	got, err := repo.LessonByConcept(ctx, "topic-a", "shared-slug")
	if err != nil {
		t.Fatalf("LessonByConcept(topic-a) unexpected error: %v", err)
	}
	if got.TitleEn() != "Title shared-slug" {
		t.Errorf("LessonByConcept(topic-a) TitleEn() = %q, want %q", got.TitleEn(), "Title shared-slug")
	}
}

// TestRepositoryLessonByConcept_MalformedRecallCheckIsNotConfusedWithNotFound
// seeds a lesson row directly (bypassing SaveLesson, which would itself
// reject this data through domain.NewLesson) to simulate a row already in
// the database that violates a recall-check invariant — an mcq with NULL
// options. Rehydration must fail with a non-ErrLessonNotFound error, so the
// handler answers 500 (logged, no leak) instead of a misleading 404.
func TestRepositoryLessonByConcept_MalformedRecallCheckIsNotConfusedWithNotFound(t *testing.T) {
	result := &stubResult{}
	result.seedConcept("domain-driven-design", "bad-lesson", 7)
	result.tables.lessons[7] = stubLessonRow{
		id: 99, version: 1, titleEn: "Bad Lesson", estMinutes: 7, bodyMd: "body",
		refs: `[{"title":"t1","source":"s1","why":"w1"},{"title":"t2","source":"s2","why":"w2"}]`,
	}
	result.tables.recallChecks[99] = []stubRecallCheckRow{
		{position: 1, kind: "mcq", question: "q", expectedAnswer: "a", options: nil},
	}
	db := openStubDB(t, result)
	repo := NewRepository(db)

	_, err := repo.LessonByConcept(context.Background(), "domain-driven-design", "bad-lesson")
	if err == nil {
		t.Fatal("LessonByConcept() expected an error for a malformed mcq recall check, got nil")
	}
	if errors.Is(err, curriculumapp.ErrLessonNotFound) {
		t.Fatal("LessonByConcept() wrapped ErrLessonNotFound for a malformed row — the handler would wrongly 404 instead of 500")
	}
}

func TestRepositoryLessonByConcept_QueryErrorWraps(t *testing.T) {
	result := seededConceptResult("domain-driven-design", "aggregate", 42)
	db := openStubDB(t, result)
	repo := NewRepository(db)
	ctx := context.Background()
	if _, err := repo.SaveLesson(ctx, "domain-driven-design", newTestLesson(t, "aggregate", 1)); err != nil {
		t.Fatalf("SaveLesson() unexpected error: %v", err)
	}

	result.queryErr = errors.New("stub: connection refused")

	_, err := repo.LessonByConcept(ctx, "domain-driven-design", "aggregate")
	if err == nil {
		t.Fatal("LessonByConcept() expected an error, got nil")
	}
	if errors.Is(err, curriculumapp.ErrLessonNotFound) {
		t.Fatal("LessonByConcept() wrapped ErrLessonNotFound for a real query error, want a different error")
	}
}

func TestSelectLessonByConceptSQLShape(t *testing.T) {
	for _, want := range []string{"t.slug = ?", "co.slug = ?", "JOIN concepts", "JOIN chapters", "JOIN topics"} {
		if !strings.Contains(selectLessonByConceptSQL, want) {
			t.Errorf("selectLessonByConceptSQL missing %q:\n%s", want, selectLessonByConceptSQL)
		}
	}
}
