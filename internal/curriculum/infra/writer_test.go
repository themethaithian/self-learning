package infra

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	curriculumapp "github.com/themethaithian/self-learning/internal/curriculum/app"
	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

var _ curriculumapp.Writer = (*Repository)(nil)

func newTestConcept(t *testing.T, slug string, pos int) domain.Concept {
	t.Helper()
	s, err := domain.NewSlug(slug)
	if err != nil {
		t.Fatalf("NewSlug(%q): %v", slug, err)
	}
	p, err := domain.NewPosition(pos)
	if err != nil {
		t.Fatalf("NewPosition(%d): %v", pos, err)
	}
	c, err := domain.NewConcept(s, "Title "+slug, "Outline "+slug, p)
	if err != nil {
		t.Fatalf("NewConcept(%q): %v", slug, err)
	}
	return c
}

func newTestChapter(t *testing.T, slug string, pos int, concepts ...domain.Concept) domain.Chapter {
	t.Helper()
	s, err := domain.NewSlug(slug)
	if err != nil {
		t.Fatalf("NewSlug(%q): %v", slug, err)
	}
	p, err := domain.NewPosition(pos)
	if err != nil {
		t.Fatalf("NewPosition(%d): %v", pos, err)
	}
	ch, err := domain.NewChapter(s, "Title "+slug, p, concepts)
	if err != nil {
		t.Fatalf("NewChapter(%q): %v", slug, err)
	}
	return ch
}

func newTestTopic(t *testing.T, slug, title string, chapters ...domain.Chapter) domain.Topic {
	t.Helper()
	track, err := domain.NewTrack("go")
	if err != nil {
		t.Fatalf("NewTrack: %v", err)
	}
	s, err := domain.NewSlug(slug)
	if err != nil {
		t.Fatalf("NewSlug(%q): %v", slug, err)
	}
	pos, err := domain.NewPosition(1)
	if err != nil {
		t.Fatalf("NewPosition: %v", err)
	}
	topic, err := domain.NewTopic(track, s, title, pos, chapters)
	if err != nil {
		t.Fatalf("NewTopic: %v", err)
	}
	return topic
}

// twoChapterTopic: chapter "syntax" (2 concepts) then chapter "concurrency"
// (1 concept) — enough shape to prove SaveTopic threads parent ids
// correctly across more than one chapter.
func twoChapterTopic(t *testing.T) domain.Topic {
	t.Helper()
	return newTestTopic(t, "go-basics", "Go Basics",
		newTestChapter(t, "syntax", 1, newTestConcept(t, "variables", 1), newTestConcept(t, "loops", 2)),
		newTestChapter(t, "concurrency", 2, newTestConcept(t, "goroutines", 1)),
	)
}

func argValues(args []driver.NamedValue) []any {
	out := make([]any, len(args))
	for i, a := range args {
		out[i] = a.Value
	}
	return out
}

func TestRepositorySaveTopic_HappyPathCommitsOnce(t *testing.T) {
	result := &stubResult{}
	db := openStubDB(t, result)
	repo := NewRepository(db)

	if err := repo.SaveTopic(context.Background(), twoChapterTopic(t)); err != nil {
		t.Fatalf("SaveTopic() unexpected error: %v", err)
	}

	began, committed, rolledBack := result.txCounts()
	if began != 1 {
		t.Errorf("tx began = %d, want 1", began)
	}
	if committed != 1 {
		t.Errorf("tx committed = %d, want 1", committed)
	}
	if rolledBack != 0 {
		t.Errorf("tx rolled back = %d, want 0", rolledBack)
	}
}

// TestRepositorySaveTopic_StatementSequenceAndArgs pins the whole write:
// statement identity, argument values in order, and — critically — that
// each child's parent-id argument equals the LastInsertId the stub returned
// for the row that must be its parent. chapter_id is the only thing that
// decides which chapter a concept belongs to, so a swapped or hardcoded
// parent id here is a "valid but wrong" curriculum, not a crash.
func TestRepositorySaveTopic_StatementSequenceAndArgs(t *testing.T) {
	result := &stubResult{}
	db := openStubDB(t, result)
	repo := NewRepository(db)
	ctx := context.Background()

	// Warm up the id counter with an unrelated topic first, so the topic
	// under test does NOT land on id=1 — otherwise a mutation that
	// hardcodes "topic_id = 1" instead of threading the real id would
	// coincidentally match and this test would not catch it.
	warmup := newTestTopic(t, "warmup-topic", "Warmup", newTestChapter(t, "warmup-chapter", 1, newTestConcept(t, "warmup-concept", 1)))
	if err := repo.SaveTopic(ctx, warmup); err != nil {
		t.Fatalf("warm-up SaveTopic() unexpected error: %v", err)
	}
	callsBeforeTest := len(result.execLog())

	if err := repo.SaveTopic(ctx, twoChapterTopic(t)); err != nil {
		t.Fatalf("SaveTopic() unexpected error: %v", err)
	}
	calls := result.execLog()[callsBeforeTest:]

	want := []struct {
		name  string
		query string
		args  []any
	}{
		{"upsert topic", upsertTopicSQL, []any{"go", "go-basics", "Go Basics", int64(1)}},
		{"upsert chapter syntax", upsertChapterSQL, nil},    // topic_id filled in below
		{"upsert concept variables", upsertConceptSQL, nil}, // chapter_id filled in below
		{"upsert concept loops", upsertConceptSQL, nil},
		{"delete stale concepts of syntax", fmt.Sprintf(deleteStaleConceptsSQL, "?,?"), nil},
		{"upsert chapter concurrency", upsertChapterSQL, nil},
		{"upsert concept goroutines", upsertConceptSQL, nil},
		{"delete stale concepts of concurrency", fmt.Sprintf(deleteStaleConceptsSQL, "?"), nil},
		{"delete concepts of stale chapters", fmt.Sprintf(deleteStaleChapterConceptsSQL, "?,?"), nil}, // topic_id filled in below
		{"delete stale chapters", fmt.Sprintf(deleteStaleChaptersSQL, "?,?"), nil},
	}
	if len(calls) != len(want) {
		t.Fatalf("exec call count = %d, want %d: %+v", len(calls), len(want), calls)
	}

	topicID := calls[0].lastInsertID
	if calls[0].query != want[0].query || !reflect.DeepEqual(argValues(calls[0].args), want[0].args) {
		t.Fatalf("call 0 = %q %v, want %q %v", calls[0].query, argValues(calls[0].args), want[0].query, want[0].args)
	}

	want[1].args = []any{topicID, "syntax", "Title syntax", int64(1)}
	want[5].args = []any{topicID, "concurrency", "Title concurrency", int64(2)}
	want[8].args = []any{topicID, "syntax", "concurrency"}
	want[9].args = []any{topicID, "syntax", "concurrency"}

	chapterSyntaxID := calls[1].lastInsertID
	want[2].args = []any{chapterSyntaxID, "variables", "Title variables", "Outline variables", int64(1)}
	want[3].args = []any{chapterSyntaxID, "loops", "Title loops", "Outline loops", int64(2)}
	want[4].args = []any{chapterSyntaxID, "variables", "loops"}

	chapterConcurrencyID := calls[5].lastInsertID
	want[6].args = []any{chapterConcurrencyID, "goroutines", "Title goroutines", "Outline goroutines", int64(1)}
	want[7].args = []any{chapterConcurrencyID, "goroutines"}

	if topicID == 1 {
		t.Fatal("topicID = 1; the warm-up write must run first so a hardcoded topic_id = 1 mutation is distinguishable from the real id")
	}
	if chapterSyntaxID == topicID || chapterConcurrencyID == topicID || chapterSyntaxID == chapterConcurrencyID {
		t.Fatalf("expected distinct ids, got topic=%d syntax=%d concurrency=%d", topicID, chapterSyntaxID, chapterConcurrencyID)
	}

	for i, w := range want {
		t.Run(w.name, func(t *testing.T) {
			if calls[i].query != w.query {
				t.Errorf("query = %q, want %q", calls[i].query, w.query)
			}
			got := argValues(calls[i].args)
			if !reflect.DeepEqual(got, w.args) {
				t.Errorf("args = %#v, want %#v", got, w.args)
			}
		})
	}
}

// TestUpsertSQLShape asserts the statement text itself, independent of any
// stub: deleting ON DUPLICATE KEY UPDATE, dropping id = LAST_INSERT_ID(id),
// or dropping an update assignment for a mutable column must all fail here,
// even though the happy-path stub tests above would keep passing (the stub
// only sees whatever text production sends; it can't tell "no ODKU clause"
// from "no update needed").
func TestUpsertSQLShape(t *testing.T) {
	tests := []struct {
		name        string
		stmt        string
		mustContain []string
	}{
		{
			name: "topic upsert",
			stmt: upsertTopicSQL,
			mustContain: []string{
				") AS new",
				"ON DUPLICATE KEY UPDATE",
				"id = LAST_INSERT_ID(id)",
				"track = new.track",
				"title = new.title",
				"position = new.position",
			},
		},
		{
			name: "chapter upsert",
			stmt: upsertChapterSQL,
			mustContain: []string{
				") AS new",
				"ON DUPLICATE KEY UPDATE",
				"id = LAST_INSERT_ID(id)",
				"title = new.title",
				"position = new.position",
			},
		},
		{
			name: "concept upsert",
			stmt: upsertConceptSQL,
			mustContain: []string{
				") AS new",
				"ON DUPLICATE KEY UPDATE",
				"id = LAST_INSERT_ID(id)",
				"title = new.title",
				"outline = new.outline",
				"position = new.position",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, want := range tt.mustContain {
				if !strings.Contains(tt.stmt, want) {
					t.Errorf("statement missing %q:\n%s", want, tt.stmt)
				}
			}
		})
	}
}

// TestDeleteSQLShape mirrors TestUpsertSQLShape for the three stale-row
// DELETE templates: it asserts against the raw %s-templated constants, not
// a Sprintf'd copy of them, so flipping NOT IN to IN (or dropping a
// predicate) fails here even though the stub — which reimplements delete
// semantics in Go rather than parsing SQL — could not otherwise tell the
// difference, and even though a query-text assertion built from the same
// constant would just move together with the bug.
func TestDeleteSQLShape(t *testing.T) {
	tests := []struct {
		name        string
		stmt        string
		mustContain []string
	}{
		{
			name: "delete stale concepts",
			stmt: deleteStaleConceptsSQL,
			mustContain: []string{
				"DELETE FROM concepts",
				"chapter_id = ?",
				"slug NOT IN (%s)",
			},
		},
		{
			name: "delete concepts of stale chapters",
			stmt: deleteStaleChapterConceptsSQL,
			mustContain: []string{
				"DELETE FROM concepts",
				"chapter_id IN (SELECT id FROM chapters WHERE topic_id = ?",
				"slug NOT IN (%s)",
			},
		},
		{
			name: "delete stale chapters",
			stmt: deleteStaleChaptersSQL,
			mustContain: []string{
				"DELETE FROM chapters",
				"topic_id = ?",
				"slug NOT IN (%s)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, want := range tt.mustContain {
				if !strings.Contains(tt.stmt, want) {
					t.Errorf("statement missing %q:\n%s", want, tt.stmt)
				}
			}
		})
	}
}

func TestRepositorySaveTopic_SecondRunIsIdempotent(t *testing.T) {
	result := &stubResult{}
	db := openStubDB(t, result)
	repo := NewRepository(db)

	topic := twoChapterTopic(t)
	if err := repo.SaveTopic(context.Background(), topic); err != nil {
		t.Fatalf("first SaveTopic() unexpected error: %v", err)
	}
	firstCalls := result.execLog()

	if err := repo.SaveTopic(context.Background(), topic); err != nil {
		t.Fatalf("second SaveTopic() unexpected error: %v", err)
	}
	allCalls := result.execLog()
	secondCalls := allCalls[len(firstCalls):]

	if len(secondCalls) != len(firstCalls) {
		t.Fatalf("second run made %d exec calls, want %d (same shape as the first run)", len(secondCalls), len(firstCalls))
	}
	for i := range firstCalls {
		if firstCalls[i].query != secondCalls[i].query {
			t.Errorf("call %d: query changed between runs: %q vs %q", i, firstCalls[i].query, secondCalls[i].query)
		}
		if !reflect.DeepEqual(argValues(firstCalls[i].args), argValues(secondCalls[i].args)) {
			t.Errorf("call %d: args changed between runs: %v vs %v", i, argValues(firstCalls[i].args), argValues(secondCalls[i].args))
		}
		if firstCalls[i].lastInsertID != secondCalls[i].lastInsertID {
			t.Errorf("call %d: id changed between runs (%d vs %d) — re-import must converge on the same rows", i, firstCalls[i].lastInsertID, secondCalls[i].lastInsertID)
		}
	}
}

func TestRepositorySaveTopic_SecondRunAppliesChangedTitle(t *testing.T) {
	result := &stubResult{}
	db := openStubDB(t, result)
	repo := NewRepository(db)

	if err := repo.SaveTopic(context.Background(), twoChapterTopic(t)); err != nil {
		t.Fatalf("first SaveTopic() unexpected error: %v", err)
	}
	callsAfterFirst := len(result.execLog())

	updated := newTestTopic(t, "go-basics", "Go Basics Updated",
		newTestChapter(t, "syntax", 1, newTestConcept(t, "variables", 1), newTestConcept(t, "loops", 2)),
		newTestChapter(t, "concurrency", 2, newTestConcept(t, "goroutines", 1)),
	)
	if err := repo.SaveTopic(context.Background(), updated); err != nil {
		t.Fatalf("second SaveTopic() unexpected error: %v", err)
	}

	secondCalls := result.execLog()[callsAfterFirst:]
	topicCall := secondCalls[0]
	if topicCall.query != upsertTopicSQL {
		t.Fatalf("first call of second run = %q, want the topic upsert", topicCall.query)
	}

	wantArgs := []any{"go", "go-basics", "Go Basics Updated", int64(1)}
	if got := argValues(topicCall.args); !reflect.DeepEqual(got, wantArgs) {
		t.Errorf("second-run topic upsert args = %v, want %v", got, wantArgs)
	}
	if topicCall.rowsAffected != 2 {
		t.Errorf("second-run topic upsert rowsAffected = %d, want 2 (the changed title must actually reach the DB)", topicCall.rowsAffected)
	}
	if topicCall.lastInsertID != 1 {
		t.Errorf("second-run topic upsert id = %d, want 1 (same row as the first run)", topicCall.lastInsertID)
	}
}

func TestRepositorySaveTopic_MidwayErrorRollsBackNeverCommits(t *testing.T) {
	execErr := errors.New("stub: connection reset")
	result := &stubResult{execErrOnCall: 3, execErr: execErr} // 3rd call: topic, chapter, then the first concept
	db := openStubDB(t, result)
	repo := NewRepository(db)

	err := repo.SaveTopic(context.Background(), twoChapterTopic(t))
	if err == nil {
		t.Fatal("SaveTopic() expected an error, got nil")
	}
	if !errors.Is(err, execErr) {
		t.Fatalf("SaveTopic() error = %v, want it to wrap %v", err, execErr)
	}

	began, committed, rolledBack := result.txCounts()
	if began != 1 {
		t.Errorf("tx began = %d, want 1", began)
	}
	if committed != 0 {
		t.Errorf("tx committed = %d, want 0 (must never commit a partial write)", committed)
	}
	if rolledBack != 1 {
		t.Errorf("tx rolled back = %d, want 1", rolledBack)
	}
}

// TestRepositorySaveTopic_DeleteFailureRollsBackNeverCommits fails on the
// first stale-concept DELETE (call 5) rather than an upsert: this is the
// path that becomes real once lessons exist and fk_lessons_concept rejects
// deleting a concept a lesson still references — it must roll back the
// whole topic, not leave the upserts from earlier in the same import committed.
func TestRepositorySaveTopic_DeleteFailureRollsBackNeverCommits(t *testing.T) {
	execErr := errors.New("stub: fk_lessons_concept violation")
	result := &stubResult{execErrOnCall: 5, execErr: execErr} // 5th call: first stale-concept DELETE
	db := openStubDB(t, result)
	repo := NewRepository(db)

	err := repo.SaveTopic(context.Background(), twoChapterTopic(t))
	if err == nil {
		t.Fatal("SaveTopic() expected an error, got nil")
	}
	if !errors.Is(err, execErr) {
		t.Fatalf("SaveTopic() error = %v, want it to wrap %v", err, execErr)
	}

	began, committed, rolledBack := result.txCounts()
	if began != 1 {
		t.Errorf("tx began = %d, want 1", began)
	}
	if committed != 0 {
		t.Errorf("tx committed = %d, want 0 (must never commit a partial write)", committed)
	}
	if rolledBack != 1 {
		t.Errorf("tx rolled back = %d, want 1", rolledBack)
	}
}

func TestRepositorySaveTopic_BeginError(t *testing.T) {
	beginErr := errors.New("stub: pool exhausted")
	result := &stubResult{beginErr: beginErr}
	db := openStubDB(t, result)
	repo := NewRepository(db)

	err := repo.SaveTopic(context.Background(), twoChapterTopic(t))
	if err == nil {
		t.Fatal("SaveTopic() expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "go-basics") {
		t.Errorf("SaveTopic() error = %q, want it to name the topic", err.Error())
	}
}

// The four tests below are R4's headline correctness requirement: the
// import is additive-only for topics across files, but within one topic a
// moved, renamed, or deleted node must not leave an orphan row that
// permanently breaks Repository.Topics (and so GET /api/v1/curriculum) via
// ErrDuplicateSlug or ErrDuplicatePosition. Each test writes an initial
// shape, writes a modified shape, then proves the read path — the same
// Topics()/assembleTree code production uses — still succeeds and reflects
// the new shape.

func TestRepositorySaveTopic_ReconcilesMovedConcept(t *testing.T) {
	result := &stubResult{}
	db := openStubDB(t, result)
	repo := NewRepository(db)
	ctx := context.Background()

	initial := newTestTopic(t, "go-basics", "Go Basics",
		newTestChapter(t, "syntax", 1, newTestConcept(t, "variables", 1), newTestConcept(t, "loops", 2)),
		newTestChapter(t, "concurrency", 2, newTestConcept(t, "goroutines", 1)),
	)
	if err := repo.SaveTopic(ctx, initial); err != nil {
		t.Fatalf("initial SaveTopic() unexpected error: %v", err)
	}

	moved := newTestTopic(t, "go-basics", "Go Basics",
		newTestChapter(t, "syntax", 1, newTestConcept(t, "variables", 1)),
		newTestChapter(t, "concurrency", 2, newTestConcept(t, "goroutines", 1), newTestConcept(t, "loops", 2)),
	)
	if err := repo.SaveTopic(ctx, moved); err != nil {
		t.Fatalf("SaveTopic() with moved concept unexpected error: %v", err)
	}

	got, err := repo.Topics(ctx)
	if err != nil {
		t.Fatalf("Topics() after reconciliation unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("Topics() returned %d topics, want 1", len(got))
	}

	chapters := got[0].Chapters()
	syntaxConcepts := chapterConceptSlugs(chapters, "syntax")
	concurrencyConcepts := chapterConceptSlugs(chapters, "concurrency")
	if !reflect.DeepEqual(syntaxConcepts, []string{"variables"}) {
		t.Errorf("syntax concepts = %v, want [variables] (loops must have moved out)", syntaxConcepts)
	}
	if !reflect.DeepEqual(concurrencyConcepts, []string{"goroutines", "loops"}) {
		t.Errorf("concurrency concepts = %v, want [goroutines loops]", concurrencyConcepts)
	}
}

func TestRepositorySaveTopic_ReconcilesRenamedChapter(t *testing.T) {
	result := &stubResult{}
	db := openStubDB(t, result)
	repo := NewRepository(db)
	ctx := context.Background()

	initial := newTestTopic(t, "go-basics", "Go Basics", newTestChapter(t, "modeling", 1, newTestConcept(t, "aggregate", 1)))
	if err := repo.SaveTopic(ctx, initial); err != nil {
		t.Fatalf("initial SaveTopic() unexpected error: %v", err)
	}

	renamed := newTestTopic(t, "go-basics", "Go Basics", newTestChapter(t, "domain-modeling", 1, newTestConcept(t, "aggregate", 1)))
	if err := repo.SaveTopic(ctx, renamed); err != nil {
		t.Fatalf("SaveTopic() with renamed chapter unexpected error: %v", err)
	}

	got, err := repo.Topics(ctx)
	if err != nil {
		t.Fatalf("Topics() after reconciliation unexpected error: %v", err)
	}
	chapters := got[0].Chapters()
	if len(chapters) != 1 {
		t.Fatalf("len(Chapters()) = %d, want 1 (old \"modeling\" chapter must be gone)", len(chapters))
	}
	if chapters[0].Slug().String() != "domain-modeling" {
		t.Errorf("Chapters()[0].Slug() = %q, want domain-modeling", chapters[0].Slug().String())
	}
	if len(chapters[0].Concepts()) != 1 {
		t.Errorf("len(Concepts()) = %d, want 1 (no duplicate \"aggregate\" left behind under the old chapter)", len(chapters[0].Concepts()))
	}
}

func TestRepositorySaveTopic_ReconcilesRenamedConcept(t *testing.T) {
	result := &stubResult{}
	db := openStubDB(t, result)
	repo := NewRepository(db)
	ctx := context.Background()

	initial := newTestTopic(t, "go-basics", "Go Basics", newTestChapter(t, "modeling", 1, newTestConcept(t, "aggregate", 1)))
	if err := repo.SaveTopic(ctx, initial); err != nil {
		t.Fatalf("initial SaveTopic() unexpected error: %v", err)
	}

	renamed := newTestTopic(t, "go-basics", "Go Basics", newTestChapter(t, "modeling", 1, newTestConcept(t, "aggregates", 1)))
	if err := repo.SaveTopic(ctx, renamed); err != nil {
		t.Fatalf("SaveTopic() with renamed concept unexpected error: %v", err)
	}

	got, err := repo.Topics(ctx)
	if err != nil {
		t.Fatalf("Topics() after reconciliation unexpected error: %v", err)
	}
	concepts := got[0].Chapters()[0].Concepts()
	if len(concepts) != 1 {
		t.Fatalf("len(Concepts()) = %d, want 1 (old \"aggregate\" slug must be gone)", len(concepts))
	}
	if concepts[0].Slug().String() != "aggregates" {
		t.Errorf("Concepts()[0].Slug() = %q, want aggregates", concepts[0].Slug().String())
	}
}

func TestRepositorySaveTopic_ReconcilesDeletedConcept(t *testing.T) {
	result := &stubResult{}
	db := openStubDB(t, result)
	repo := NewRepository(db)
	ctx := context.Background()

	initial := newTestTopic(t, "go-basics", "Go Basics",
		newTestChapter(t, "syntax", 1, newTestConcept(t, "variables", 1), newTestConcept(t, "loops", 2)),
	)
	if err := repo.SaveTopic(ctx, initial); err != nil {
		t.Fatalf("initial SaveTopic() unexpected error: %v", err)
	}

	trimmed := newTestTopic(t, "go-basics", "Go Basics", newTestChapter(t, "syntax", 1, newTestConcept(t, "variables", 1)))
	if err := repo.SaveTopic(ctx, trimmed); err != nil {
		t.Fatalf("SaveTopic() with deleted concept unexpected error: %v", err)
	}

	got, err := repo.Topics(ctx)
	if err != nil {
		t.Fatalf("Topics() after reconciliation unexpected error: %v", err)
	}
	concepts := got[0].Chapters()[0].Concepts()
	if len(concepts) != 1 {
		t.Fatalf("len(Concepts()) = %d, want 1 (deleted \"loops\" must be gone)", len(concepts))
	}
	if concepts[0].Slug().String() != "variables" {
		t.Errorf("Concepts()[0].Slug() = %q, want variables", concepts[0].Slug().String())
	}
}

func chapterConceptSlugs(chapters []domain.Chapter, chapterSlug string) []string {
	for _, ch := range chapters {
		if ch.Slug().String() != chapterSlug {
			continue
		}
		out := make([]string, 0, len(ch.Concepts()))
		for _, c := range ch.Concepts() {
			out = append(out, c.Slug().String())
		}
		return out
	}
	return nil
}
