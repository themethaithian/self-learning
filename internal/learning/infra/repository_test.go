package infra

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	learningapp "github.com/themethaithian/self-learning/internal/learning/app"
	"github.com/themethaithian/self-learning/internal/learning/domain"
	"github.com/themethaithian/self-learning/migrations"
)

func mustChunkState(t *testing.T, raw string) domain.ChunkState {
	t.Helper()
	s, err := domain.NewChunkState(raw)
	if err != nil {
		t.Fatalf("NewChunkState(%q) failed: %v", raw, err)
	}
	return s
}

func TestRepositoryTransition_LessonNotFound(t *testing.T) {
	db := openStubDB(t, newStubData())
	svc := learningapp.NewService(NewRepository(db))

	_, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if !errors.Is(err, learningapp.ErrLessonNotFound) {
		t.Fatalf("SetProgress() error = %v, want it to wrap ErrLessonNotFound", err)
	}
}

// TestRepositoryTransition_ScopedByBothSlugs is the concept-slug analogue of
// curriculum's ScopedByBothSlugs test: a concept slug shared by two topics
// must not leak the other topic's progress. Note this does NOT catch a
// dropped "AND co.slug = ?" predicate in selectLessonForUpdateSQL itself —
// the stub dispatches on query-string identity, so an edit to the SQL text
// is invisible to this behavioral test either way.
// TestSelectLessonForUpdateSQLShape is the actual guard for the predicate.
func TestRepositoryTransition_ScopedByBothSlugs(t *testing.T) {
	data := newStubData()
	data.seedProgress("topic-a", "shared-slug", 101, stubProgressRow{state: "passed"})
	data.seedLesson("topic-b", "shared-slug", 202)
	db := openStubDB(t, data)
	svc := learningapp.NewService(NewRepository(db))
	ctx := context.Background()

	entryB, err := svc.SetProgress(ctx, "topic-b", "shared-slug", "in_progress")
	if err != nil {
		t.Fatalf("SetProgress(topic-b) unexpected error: %v", err)
	}
	if !entryB.State.IsInProgress() {
		t.Fatalf("SetProgress(topic-b) State = %v, want in_progress (must never leak topic-a's passed state)", entryB.State)
	}

	row, ok := data.progressOf(101)
	if !ok || row.state != "passed" {
		t.Fatalf("topic-a's row = %+v (ok=%v), want unchanged passed", row, ok)
	}
}

// TestRepositoryTransition_DuplicateConceptSlugFailsLoudly pins S7:
// concepts is unique per (chapter_id, slug), not per topic, so two chapters
// of the same topic can legally share a concept slug in the real schema.
// lockLesson must fail loudly (an error, never a write) rather than
// silently picking one of the matching lessons the way QueryRow would.
func TestRepositoryTransition_DuplicateConceptSlugFailsLoudly(t *testing.T) {
	data := newStubData()
	data.seedLesson("ddia", "duplicated-slug", 1)
	data.seedLesson("ddia", "duplicated-slug", 2)
	db := openStubDB(t, data)
	svc := learningapp.NewService(NewRepository(db))

	_, err := svc.SetProgress(context.Background(), "ddia", "duplicated-slug", "in_progress")
	if err == nil {
		t.Fatal("SetProgress() expected an error for a concept slug matching two lessons, got nil")
	}

	if _, ok := data.progressOf(1); ok {
		t.Error("lesson 1 got a progress row written despite the ambiguous match")
	}
	if _, ok := data.progressOf(2); ok {
		t.Error("lesson 2 got a progress row written despite the ambiguous match")
	}
}

func TestRepositoryTransition_FreshInProgress(t *testing.T) {
	data := newStubData()
	data.seedLesson("ddia", "b-trees", 1)
	db := openStubDB(t, data)
	svc := learningapp.NewService(NewRepository(db))

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsInProgress() {
		t.Fatalf("State = %v, want in_progress", entry.State)
	}
	if entry.FirstPassedAt != nil {
		t.Errorf("FirstPassedAt = %v, want nil for a fresh in_progress row", entry.FirstPassedAt)
	}
	if entry.LastReadAt == nil {
		t.Error("LastReadAt = nil, want it set")
	}
}

func TestRepositoryTransition_FreshPassed(t *testing.T) {
	data := newStubData()
	data.seedLesson("ddia", "b-trees", 1)
	db := openStubDB(t, data)
	svc := learningapp.NewService(NewRepository(db))

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "passed")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsPassed() {
		t.Fatalf("State = %v, want passed", entry.State)
	}
	if entry.FirstPassedAt == nil {
		t.Error("FirstPassedAt = nil, want it set on first pass")
	}
}

// TestRepositoryTransition_PreservesFirstPassedAt pins the COALESCE clause
// in upsertProgressByLessonIDSQL: a re-finish must never move
// first_passed_at.
func TestRepositoryTransition_PreservesFirstPassedAt(t *testing.T) {
	firstPassed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	data := newStubData()
	data.seedProgress("ddia", "b-trees", 1, stubProgressRow{state: "passed", firstPassedAt: &firstPassed})
	db := openStubDB(t, data)
	svc := learningapp.NewService(NewRepository(db))

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "passed")
	if err != nil {
		t.Fatalf("SetProgress() re-finish unexpected error: %v", err)
	}
	if entry.FirstPassedAt == nil || !entry.FirstPassedAt.Equal(firstPassed) {
		t.Errorf("FirstPassedAt = %v, want unchanged %v", entry.FirstPassedAt, firstPassed)
	}
}

// TestRepositoryTransition_ForwardOnlyClamp pins the whole reason
// refreshLastReadAtSQL exists: it must move last_read_at and nothing else.
func TestRepositoryTransition_ForwardOnlyClamp(t *testing.T) {
	firstPassed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	lastRead := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	data := newStubData()
	data.seedProgress("ddia", "b-trees", 1, stubProgressRow{state: "passed", firstPassedAt: &firstPassed, lastReadAt: &lastRead})
	db := openStubDB(t, data)
	svc := learningapp.NewService(NewRepository(db))

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsPassed() {
		t.Errorf("State = %v, want unchanged passed (forward-only)", entry.State)
	}
	if entry.FirstPassedAt == nil || !entry.FirstPassedAt.Equal(firstPassed) {
		t.Errorf("FirstPassedAt = %v, want unchanged %v", entry.FirstPassedAt, firstPassed)
	}
	if entry.LastReadAt == nil || entry.LastReadAt.Equal(lastRead) {
		t.Errorf("LastReadAt = %v, want it refreshed away from %v", entry.LastReadAt, lastRead)
	}
}

// TestRepositoryTransition_LockedThenPassed pins R1(b): a state='locked' row
// (a future gating ticket's edge; this ticket never writes it) receiving
// {"state":"passed"} must succeed via unlock-then-pass, not 500.
func TestRepositoryTransition_LockedThenPassed(t *testing.T) {
	data := newStubData()
	data.seedProgress("ddia", "b-trees", 1, stubProgressRow{state: "locked"})
	db := openStubDB(t, data)
	svc := learningapp.NewService(NewRepository(db))

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "passed")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsPassed() {
		t.Fatalf("State = %v, want passed", entry.State)
	}
	if entry.FirstPassedAt == nil {
		t.Error("FirstPassedAt = nil, want it set")
	}
}

func TestRepositoryTransition_LockedThenInProgress(t *testing.T) {
	data := newStubData()
	data.seedProgress("ddia", "b-trees", 1, stubProgressRow{state: "locked"})
	db := openStubDB(t, data)
	svc := learningapp.NewService(NewRepository(db))

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsInProgress() {
		t.Fatalf("State = %v, want in_progress", entry.State)
	}
}

// TestRepositoryTransition_CanonicalSlugsReturned pins S1 directly against
// Repository, bypassing Service's stricter shape validation (which now
// rejects mixed-case input before any DB round trip): even called with the
// wrong case, the DB's own slugs must come back, never the caller's.
func TestRepositoryTransition_CanonicalSlugsReturned(t *testing.T) {
	data := newStubData()
	data.seedLesson("ddia", "b-trees", 1)
	db := openStubDB(t, data)
	repo := NewRepository(db)

	entry, err := repo.Transition(context.Background(), "DDIA", "B-Trees", func(_ learningapp.ProgressEntry, _ bool) (learningapp.ProgressDecision, error) {
		return learningapp.ProgressDecision{State: mustChunkState(t, "in_progress")}, nil
	})
	if err != nil {
		t.Fatalf("Transition() unexpected error: %v", err)
	}
	if entry.Topic != "ddia" || entry.Concept != "b-trees" {
		t.Errorf("Topic/Concept = %q/%q, want the DB's canonical ddia/b-trees, not the caller's DDIA/B-Trees", entry.Topic, entry.Concept)
	}
}

func TestRepositoryTransition_QueryErrorWraps(t *testing.T) {
	data := newStubData()
	data.queryErr = errors.New("stub: connection refused")
	db := openStubDB(t, data)
	svc := learningapp.NewService(NewRepository(db))

	_, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if err == nil {
		t.Fatal("SetProgress() expected an error, got nil")
	}
}

func TestRepositoryAllProgress_OnlyReturnsProgressRows(t *testing.T) {
	data := newStubData()
	data.seedLesson("ddia", "never-started", 1)
	data.seedProgress("ddia", "b-trees", 2, stubProgressRow{state: "in_progress"})
	data.seedProgress("ai-systems", "prompting", 3, stubProgressRow{state: "passed"})
	db := openStubDB(t, data)
	repo := NewRepository(db)

	entries, err := repo.AllProgress(context.Background())
	if err != nil {
		t.Fatalf("AllProgress() unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("AllProgress() returned %d entries, want 2 (never-started must be absent)", len(entries))
	}
	if entries[0].Topic != "ai-systems" || entries[0].Concept != "prompting" {
		t.Errorf("entries[0] = %+v, want ai-systems/prompting first (ORDER BY t.slug, co.slug)", entries[0])
	}
	if entries[1].Topic != "ddia" || entries[1].Concept != "b-trees" {
		t.Errorf("entries[1] = %+v, want ddia/b-trees", entries[1])
	}
}

func TestRepositoryAllProgress_Empty(t *testing.T) {
	db := openStubDB(t, newStubData())
	repo := NewRepository(db)

	entries, err := repo.AllProgress(context.Background())
	if err != nil {
		t.Fatalf("AllProgress() unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("AllProgress() = %v, want empty", entries)
	}
}

// TestRepositoryTransition_ConcurrentRaceNeverContradicts checks a
// conditional, not an unconditional guarantee: IF InnoDB really serializes
// two transactions at the SELECT ... FOR UPDATE lock on the same row — a
// documented locking-read guarantee this test cannot itself exercise, since
// the stub's per-lesson sync.Mutex only stands in for it — THEN
// decideTransition's logic converges on state=passed under either ordering
// for two racing writers on a fresh concept: PUT in_progress (as fired when
// a lesson opens) and PUT passed (Finish). Before the fix, the read and the
// write were two separate, unsynchronized statements, so whichever wrote
// last won unconditionally (state = new.state) regardless of what the other
// had just committed — an in_progress write landing after a passed write
// left state=in_progress with first_passed_at already non-nil. Neither this
// test nor this ticket's e2e curl runs establish the InnoDB-side premise:
// the critical section is a few milliseconds, so shell-launched concurrent
// requests essentially never actually interleave: that premise rests on
// InnoDB's documented FOR UPDATE semantics, not on any test here.
func TestRepositoryTransition_ConcurrentRaceNeverContradicts(t *testing.T) {
	const iterations = 50
	for i := 0; i < iterations; i++ {
		data := newStubData()
		data.seedLesson("ddia", "b-trees", 1)
		db := openStubDB(t, data)
		svc := learningapp.NewService(NewRepository(db))
		ctx := context.Background()

		var wg sync.WaitGroup
		var inProgressErr, passedErr error
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, inProgressErr = svc.SetProgress(ctx, "ddia", "b-trees", "in_progress")
		}()
		go func() {
			defer wg.Done()
			_, passedErr = svc.SetProgress(ctx, "ddia", "b-trees", "passed")
		}()
		wg.Wait()
		db.Close()

		if inProgressErr != nil {
			t.Fatalf("iteration %d: in_progress goroutine error: %v", i, inProgressErr)
		}
		if passedErr != nil {
			t.Fatalf("iteration %d: passed goroutine error: %v", i, passedErr)
		}

		row, ok := data.progressOf(1)
		if !ok {
			t.Fatalf("iteration %d: no progress row after concurrent transitions", i)
		}
		if row.state != "passed" {
			t.Fatalf("iteration %d: state = %q, want passed — both orderings must converge here (an in_progress goroutine erroring silently would otherwise leave state=in_progress and pass this test)", i, row.state)
		}
	}
}

func TestSelectLessonForUpdateSQLShape(t *testing.T) {
	want := `
SELECT l.id, t.slug, co.slug
FROM lessons l
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
WHERE t.slug = ? AND co.slug = ?
FOR UPDATE OF l`
	if selectLessonForUpdateSQL != want {
		t.Errorf("selectLessonForUpdateSQL =\n%q\nwant\n%q", selectLessonForUpdateSQL, want)
	}
}

func TestSelectProgressByLessonIDSQLShape(t *testing.T) {
	want := `
SELECT state, first_passed_at, last_read_at
FROM lesson_progress
WHERE lesson_id = ?`
	if selectProgressByLessonIDSQL != want {
		t.Errorf("selectProgressByLessonIDSQL =\n%q\nwant\n%q", selectProgressByLessonIDSQL, want)
	}
}

func TestUpsertProgressByLessonIDSQLShape(t *testing.T) {
	want := `
INSERT INTO lesson_progress (lesson_id, state, first_passed_at, last_read_at)
VALUES (?, ?, ?, ?)
AS new
ON DUPLICATE KEY UPDATE
    state = new.state,
    first_passed_at = COALESCE(lesson_progress.first_passed_at, new.first_passed_at),
    last_read_at = new.last_read_at`
	if upsertProgressByLessonIDSQL != want {
		t.Errorf("upsertProgressByLessonIDSQL =\n%q\nwant\n%q", upsertProgressByLessonIDSQL, want)
	}
}

func TestRefreshLastReadAtSQLShape(t *testing.T) {
	want := `
UPDATE lesson_progress
SET last_read_at = ?
WHERE lesson_id = ?`
	if refreshLastReadAtSQL != want {
		t.Errorf("refreshLastReadAtSQL =\n%q\nwant\n%q", refreshLastReadAtSQL, want)
	}
}

func TestSelectAllProgressSQLShape(t *testing.T) {
	want := `
SELECT t.slug, co.slug, lp.state, lp.first_passed_at, lp.last_read_at
FROM lesson_progress lp
JOIN lessons l ON l.id = lp.lesson_id
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
ORDER BY t.slug, co.slug`
	if selectAllProgressSQL != want {
		t.Errorf("selectAllProgressSQL =\n%q\nwant\n%q", selectAllProgressSQL, want)
	}
}

// TestLessonProgressTableNameConsistency guards against the migration and
// the Go-side SQL drifting onto different table names — nothing else would
// notice, since no test here touches a real MySQL schema.
func TestLessonProgressTableNameConsistency(t *testing.T) {
	migrationSQL, err := migrations.FS.ReadFile("002_learning.sql")
	if err != nil {
		t.Fatalf("read 002_learning.sql: %v", err)
	}
	if !strings.Contains(string(migrationSQL), "CREATE TABLE IF NOT EXISTS lesson_progress") {
		t.Fatalf("002_learning.sql does not contain %q", "CREATE TABLE IF NOT EXISTS lesson_progress")
	}

	// selectLessonForUpdateSQL is deliberately excluded: it locks the
	// lessons row and never names lesson_progress at all (a fresh concept
	// has no lesson_progress row yet to lock) — TestSelectLessonForUpdateSQLShape
	// pins its exact text instead. Every statement below does reference
	// lesson_progress and has its own *SQLShape test besides.
	stmts := map[string]string{
		"selectProgressByLessonIDSQL": selectProgressByLessonIDSQL,
		"upsertProgressByLessonIDSQL": upsertProgressByLessonIDSQL,
		"refreshLastReadAtSQL":        refreshLastReadAtSQL,
		"selectAllProgressSQL":        selectAllProgressSQL,
	}
	for name, stmt := range stmts {
		if !strings.Contains(stmt, "lesson_progress") {
			t.Errorf("%s = %q, want it to reference the lesson_progress table", name, stmt)
		}
	}
}
