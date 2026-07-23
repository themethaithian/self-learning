package infra

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

func newTestReference(t *testing.T, title string) domain.Reference {
	t.Helper()
	r, err := domain.NewReference(title, "source "+title, "why "+title)
	if err != nil {
		t.Fatalf("NewReference(%q): %v", title, err)
	}
	return r
}

func newTestRecallCheck(t *testing.T, pos int, kind string, options []string) domain.RecallCheck {
	t.Helper()
	p, err := domain.NewPosition(pos)
	if err != nil {
		t.Fatalf("NewPosition(%d): %v", pos, err)
	}
	k, err := domain.NewRecallKind(kind)
	if err != nil {
		t.Fatalf("NewRecallKind(%q): %v", kind, err)
	}
	rc, err := domain.NewRecallCheck(p, k, fmt.Sprintf("question %d", pos), fmt.Sprintf("answer %d", pos), options)
	if err != nil {
		t.Fatalf("NewRecallCheck(%d): %v", pos, err)
	}
	return rc
}

// newTestLesson builds a lesson with 3 recall checks — short_answer, mcq,
// short_answer, in that order — so tests can assert the options column is
// NULL for the short_answer rows and a JSON array for the mcq row.
func newTestLesson(t *testing.T, slug string, version int) domain.Lesson {
	t.Helper()
	s, err := domain.NewSlug(slug)
	if err != nil {
		t.Fatalf("NewSlug(%q): %v", slug, err)
	}
	est, err := domain.NewEstMinutes(7)
	if err != nil {
		t.Fatalf("NewEstMinutes: %v", err)
	}
	l, err := domain.NewLesson(s, version, "Title "+slug, est, "body "+slug,
		[]domain.Reference{newTestReference(t, "ref-a"), newTestReference(t, "ref-b")},
		[]domain.RecallCheck{
			newTestRecallCheck(t, 1, "short_answer", nil),
			newTestRecallCheck(t, 2, "mcq", []string{"opt-a", "opt-b"}),
			newTestRecallCheck(t, 3, "short_answer", nil),
		},
	)
	if err != nil {
		t.Fatalf("NewLesson(%q): %v", slug, err)
	}
	return l
}

// seededConceptResult returns a stubResult whose concept lookup resolves
// (topicSlug, conceptSlug) to id through the relational lessonConcepts
// table, exercising SaveLesson's real WHERE predicate rather than a canned
// row that would answer the same id regardless of topic.
func seededConceptResult(topicSlug, conceptSlug string, id int64) *stubResult {
	result := &stubResult{}
	result.seedConcept(topicSlug, conceptSlug, id)
	return result
}

func TestRepositorySaveLesson_HappyPathCommitsOnce(t *testing.T) {
	result := seededConceptResult("domain-driven-design", "aggregate", 42)
	db := openStubDB(t, result)
	repo := NewRepository(db)

	inserted, err := repo.SaveLesson(context.Background(), "domain-driven-design", newTestLesson(t, "aggregate", 1))
	if err != nil {
		t.Fatalf("SaveLesson() unexpected error: %v", err)
	}
	if !inserted {
		t.Error("inserted = false, want true for a brand-new lesson")
	}

	began, committed, rolledBack := result.txCounts()
	if began != 1 || committed != 1 || rolledBack != 0 {
		t.Errorf("tx (began, committed, rolledBack) = (%d, %d, %d), want (1, 1, 0)", began, committed, rolledBack)
	}
}

func TestRepositorySaveLesson_ConceptLookupUsesTopicAndConceptSlug(t *testing.T) {
	result := seededConceptResult("domain-driven-design", "aggregate", 42)
	db := openStubDB(t, result)
	repo := NewRepository(db)

	if _, err := repo.SaveLesson(context.Background(), "domain-driven-design", newTestLesson(t, "aggregate", 1)); err != nil {
		t.Fatalf("SaveLesson() unexpected error: %v", err)
	}

	if got := result.query(); got != selectConceptIDSQL {
		t.Fatalf("query = %q, want the concept lookup", got)
	}
	args := result.queryArgs()
	if len(args) != 2 {
		t.Fatalf("query args = %v, want 2 (topic slug, concept slug)", args)
	}
	if args[0].Value != "domain-driven-design" || args[1].Value != "aggregate" {
		t.Errorf("query args = %v, want [domain-driven-design aggregate]", []any{args[0].Value, args[1].Value})
	}
}

func TestRepositorySaveLesson_ConceptNotFound(t *testing.T) {
	result := &stubResult{} // no rows fixture: the lookup finds nothing
	db := openStubDB(t, result)
	repo := NewRepository(db)

	_, err := repo.SaveLesson(context.Background(), "domain-driven-design", newTestLesson(t, "aggregate", 1))
	if err == nil {
		t.Fatal("SaveLesson() expected error, got nil")
	}
	for _, want := range []string{"domain-driven-design", "aggregate", "not found", "import curriculum first"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("SaveLesson() error = %q, want it to contain %q", err.Error(), want)
		}
	}

	began, committed, rolledBack := result.txCounts()
	if began != 1 || committed != 0 || rolledBack != 1 {
		t.Errorf("tx (began, committed, rolledBack) = (%d, %d, %d), want (1, 0, 1)", began, committed, rolledBack)
	}
}

// TestRepositorySaveLesson_ResolvesConceptScopedToTopic pins the reason the
// lookup filters on t.slug at all: concept slugs are unique only within a
// topic (NewTopic's own invariant, enforced at curriculum-import time), not
// across topics. Two tracks can each have their own "aggregate" concept, so
// a lookup that dropped the topic predicate would silently resolve to
// whichever one got seeded — importing a lesson under the wrong topic's
// concept with no error.
func TestRepositorySaveLesson_ResolvesConceptScopedToTopic(t *testing.T) {
	result := &stubResult{}
	result.seedConcept("topic-a", "shared-slug", 101)
	result.seedConcept("topic-b", "shared-slug", 202)
	db := openStubDB(t, result)
	repo := NewRepository(db)

	if _, err := repo.SaveLesson(context.Background(), "topic-a", newTestLesson(t, "shared-slug", 1)); err != nil {
		t.Fatalf("SaveLesson(topic-a) unexpected error: %v", err)
	}

	calls := result.execLog()
	if len(calls) == 0 || calls[0].query != upsertLessonSQL {
		t.Fatalf("call 0 = %+v, want the lesson upsert", calls)
	}
	gotConceptID := argValues(calls[0].args)[0]
	if gotConceptID != int64(101) {
		t.Errorf("SaveLesson(topic-a) upserted lesson under concept_id = %v, want 101 (topic-a's concept, not topic-b's 202)", gotConceptID)
	}
}

func TestRepositorySaveLesson_ConceptOnlyUnderAnotherTopicIsNotFound(t *testing.T) {
	result := &stubResult{}
	result.seedConcept("topic-b", "shared-slug", 202) // never seeded under topic-a
	db := openStubDB(t, result)
	repo := NewRepository(db)

	_, err := repo.SaveLesson(context.Background(), "topic-a", newTestLesson(t, "shared-slug", 1))
	if err == nil {
		t.Fatal("SaveLesson(topic-a) expected error, got nil (the concept only exists under topic-b)")
	}
	for _, want := range []string{"topic-a", "shared-slug", "not found"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("SaveLesson() error = %q, want it to contain %q", err.Error(), want)
		}
	}
}

// TestSelectConceptIDSQLShape mirrors TestUpsertSQLShape/TestDeleteSQLShape:
// it pins the raw statement text itself, independent of the stub. The
// relational tests above prove resolveConceptID threads both args
// correctly; this proves the WHERE clause still names both columns, since
// a stub that dispatches on query-text identity can't detect the text of
// that same constant being edited out from under it.
func TestSelectConceptIDSQLShape(t *testing.T) {
	for _, want := range []string{"t.slug = ?", "co.slug = ?"} {
		if !strings.Contains(selectConceptIDSQL, want) {
			t.Errorf("selectConceptIDSQL missing %q:\n%s", want, selectConceptIDSQL)
		}
	}
}

func TestUpsertLessonSQLShape(t *testing.T) {
	mustContain := []string{
		") AS new",
		"ON DUPLICATE KEY UPDATE",
		"id = LAST_INSERT_ID(id)",
		"version = new.version",
		"title_en = new.title_en",
		"est_minutes = new.est_minutes",
		"body_md = new.body_md",
		"refs = new.refs",
	}
	for _, want := range mustContain {
		if !strings.Contains(upsertLessonSQL, want) {
			t.Errorf("upsertLessonSQL missing %q:\n%s", want, upsertLessonSQL)
		}
	}
}

// TestRepositorySaveLesson_RecallChecksReplacedInOrder pins the write after
// the lesson upsert: a DELETE of every existing recall_checks row for the
// lesson, then one INSERT per check in position order, with the options
// column NULL for short_answer and a JSON array string for mcq.
func TestRepositorySaveLesson_RecallChecksReplacedInOrder(t *testing.T) {
	result := seededConceptResult("domain-driven-design", "aggregate", 42)
	db := openStubDB(t, result)
	repo := NewRepository(db)

	if _, err := repo.SaveLesson(context.Background(), "domain-driven-design", newTestLesson(t, "aggregate", 1)); err != nil {
		t.Fatalf("SaveLesson() unexpected error: %v", err)
	}
	calls := result.execLog()

	if len(calls) != 5 {
		t.Fatalf("exec call count = %d, want 5 (upsert lesson, delete recall_checks, insert x3)", len(calls))
	}
	if calls[0].query != upsertLessonSQL {
		t.Fatalf("call 0 query = %q, want the lesson upsert", calls[0].query)
	}
	lessonID := calls[0].lastInsertID

	if calls[1].query != deleteRecallChecksSQL {
		t.Fatalf("call 1 query = %q, want the recall_checks delete", calls[1].query)
	}
	if got := argValues(calls[1].args); !reflect.DeepEqual(got, []any{lessonID}) {
		t.Errorf("delete args = %v, want [%v]", got, lessonID)
	}

	wantInserts := []struct {
		position int64
		kind     string
		options  any
	}{
		{1, "short_answer", nil},
		{2, "mcq", `["opt-a","opt-b"]`},
		{3, "short_answer", nil},
	}
	for i, want := range wantInserts {
		call := calls[2+i]
		if call.query != insertRecallCheckSQL {
			t.Fatalf("call %d query = %q, want the recall_check insert", 2+i, call.query)
		}
		got := argValues(call.args)
		wantArgs := []any{lessonID, want.position, want.kind, fmt.Sprintf("question %d", want.position), fmt.Sprintf("answer %d", want.position), want.options}
		if !reflect.DeepEqual(got, wantArgs) {
			t.Errorf("insert recall_check %d args = %v, want %v", want.position, got, wantArgs)
		}
	}
}

func TestRepositorySaveLesson_SecondRunIsIdempotent(t *testing.T) {
	result := seededConceptResult("domain-driven-design", "aggregate", 42)
	db := openStubDB(t, result)
	repo := NewRepository(db)
	ctx := context.Background()

	lesson := newTestLesson(t, "aggregate", 1)
	firstInserted, err := repo.SaveLesson(ctx, "domain-driven-design", lesson)
	if err != nil {
		t.Fatalf("first SaveLesson() unexpected error: %v", err)
	}
	if !firstInserted {
		t.Fatal("first SaveLesson() inserted = false, want true")
	}
	firstLessonID := result.execLog()[0].lastInsertID

	secondInserted, err := repo.SaveLesson(ctx, "domain-driven-design", lesson)
	if err != nil {
		t.Fatalf("second SaveLesson() unexpected error: %v", err)
	}
	if secondInserted {
		t.Error("second SaveLesson() inserted = true, want false (row already existed)")
	}

	secondCalls := result.execLog()[5:]
	if secondCalls[0].lastInsertID != firstLessonID {
		t.Errorf("second run lesson id = %d, want %d (same row as the first run)", secondCalls[0].lastInsertID, firstLessonID)
	}
	if len(secondCalls) != 5 {
		t.Fatalf("second run made %d exec calls, want 5", len(secondCalls))
	}
}

func TestRepositorySaveLesson_MidwayErrorRollsBackNeverCommits(t *testing.T) {
	execErr := errors.New("stub: connection reset")
	result := &stubResult{execErrOnCall: 3, execErr: execErr} // 3rd call: the first recall_check insert
	result.seedConcept("domain-driven-design", "aggregate", 42)
	db := openStubDB(t, result)
	repo := NewRepository(db)

	_, err := repo.SaveLesson(context.Background(), "domain-driven-design", newTestLesson(t, "aggregate", 1))
	if err == nil {
		t.Fatal("SaveLesson() expected an error, got nil")
	}
	if !errors.Is(err, execErr) {
		t.Fatalf("SaveLesson() error = %v, want it to wrap %v", err, execErr)
	}

	began, committed, rolledBack := result.txCounts()
	if began != 1 || committed != 0 || rolledBack != 1 {
		t.Errorf("tx (began, committed, rolledBack) = (%d, %d, %d), want (1, 0, 1)", began, committed, rolledBack)
	}
}
