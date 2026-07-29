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

func mustCheckKind(t *testing.T, raw string) domain.CheckKind {
	t.Helper()
	k, err := domain.NewCheckKind(raw)
	if err != nil {
		t.Fatalf("NewCheckKind(%q) failed: %v", raw, err)
	}
	return k
}

func mustConfidence(t *testing.T, raw string) domain.Confidence {
	t.Helper()
	c, err := domain.NewConfidence(raw)
	if err != nil {
		t.Fatalf("NewConfidence(%q) failed: %v", raw, err)
	}
	return c
}

func mustAttemptOutcome(t *testing.T, raw string) domain.AttemptOutcome {
	t.Helper()
	o, err := domain.NewAttemptOutcome(raw)
	if err != nil {
		t.Fatalf("NewAttemptOutcome(%q) failed: %v", raw, err)
	}
	return o
}

func mustCanonicalQuestion(t *testing.T, raw string) domain.CanonicalQuestion {
	t.Helper()
	q, err := domain.NewCanonicalQuestion(raw)
	if err != nil {
		t.Fatalf("NewCanonicalQuestion(%q) failed: %v", raw, err)
	}
	return q
}

func mustCheckKey(t *testing.T, topic, conceptSlug, question string) domain.CheckKey {
	t.Helper()
	concept, err := domain.NewLessonRef(conceptSlug)
	if err != nil {
		t.Fatalf("NewLessonRef(%q) failed: %v", conceptSlug, err)
	}
	k, err := domain.NewCheckKey(topic, concept, mustCanonicalQuestion(t, question))
	if err != nil {
		t.Fatalf("NewCheckKey(%q, %q, %q) failed: %v", topic, conceptSlug, question, err)
	}
	return k
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

func TestSelectRecallCheckExistsSQLShape(t *testing.T) {
	want := `
SELECT l.id, rc.question, rc.type
FROM recall_checks rc
JOIN lessons l ON l.id = rc.lesson_id
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
WHERE t.slug = ? AND co.slug = ? AND rc.question = ?`
	if selectRecallCheckExistsSQL != want {
		t.Errorf("selectRecallCheckExistsSQL =\n%q\nwant\n%q", selectRecallCheckExistsSQL, want)
	}
}

func TestSelectLessonExistsSQLShape(t *testing.T) {
	want := `
SELECT l.id
FROM lessons l
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
WHERE t.slug = ? AND co.slug = ?`
	if selectLessonExistsSQL != want {
		t.Errorf("selectLessonExistsSQL =\n%q\nwant\n%q", selectLessonExistsSQL, want)
	}
}

func TestInsertRecallAttemptSQLShape(t *testing.T) {
	want := `
INSERT INTO recall_attempts (lesson_id, check_key, type, confidence, outcome, selected_option, graded_by, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if insertRecallAttemptSQL != want {
		t.Errorf("insertRecallAttemptSQL =\n%q\nwant\n%q", insertRecallAttemptSQL, want)
	}
}

// buildAttemptFunc returns a build closure matching Repository.RecordAttempt's
// contract: it derives the CheckKey from whatever canonicalQuestion it is
// CALLED with, never from a question captured at closure-creation time, and
// it uses whatever kind it is CALLED with, never a kind supplied by the
// test — exactly what app.Service.RecordAttempt's own build closure does
// (R1, both for the question and for kind).
func buildAttemptFunc(t *testing.T, topic, conceptSlug, confidence, outcome string, selectedOption *string, gradedBy domain.GradedBy) func(domain.CanonicalQuestion, domain.CheckKind) (domain.RecallAttempt, error) {
	t.Helper()
	concept, err := domain.NewLessonRef(conceptSlug)
	if err != nil {
		t.Fatalf("NewLessonRef(%q) failed: %v", conceptSlug, err)
	}
	return func(canonicalQuestion domain.CanonicalQuestion, kind domain.CheckKind) (domain.RecallAttempt, error) {
		checkKey, err := domain.NewCheckKey(topic, concept, canonicalQuestion)
		if err != nil {
			return domain.RecallAttempt{}, err
		}
		return domain.NewRecallAttempt(checkKey, kind, mustConfidence(t, confidence), mustAttemptOutcome(t, outcome), selectedOption, gradedBy)
	}
}

// neverCalledBuild fails the test immediately if RecordAttempt ever invokes
// build — used by tests expecting resolution to fail before build runs.
func neverCalledBuild(t *testing.T) func(domain.CanonicalQuestion, domain.CheckKind) (domain.RecallAttempt, error) {
	return func(canonicalQuestion domain.CanonicalQuestion, kind domain.CheckKind) (domain.RecallAttempt, error) {
		t.Fatalf("build must not be called when the lesson/question does not resolve (got canonicalQuestion=%q, kind=%q)", canonicalQuestion.String(), kind.String())
		return domain.RecallAttempt{}, nil
	}
}

func TestRepositoryRecordAttempt_LessonNotFound(t *testing.T) {
	db := openStubDB(t, newStubData())
	repo := NewRepository(db)

	_, err := repo.RecordAttempt(context.Background(), "ddia", "b-trees", "What is a B-tree?", neverCalledBuild(t))
	if !errors.Is(err, learningapp.ErrLessonNotFound) {
		t.Fatalf("RecordAttempt() error = %v, want it to wrap ErrLessonNotFound", err)
	}
}

// TestRepositoryRecordAttempt_CheckNotInLesson pins the ticket's core
// requirement: a lesson that exists but a question that does not match any
// of its recall_checks must reject, not silently mint an orphan check_key.
func TestRepositoryRecordAttempt_CheckNotInLesson(t *testing.T) {
	data := newStubData()
	data.seedRecallCheck("ddia", "b-trees", 1, "What is a B-tree?", "short_answer")
	db := openStubDB(t, data)
	repo := NewRepository(db)

	_, err := repo.RecordAttempt(context.Background(), "ddia", "b-trees", "A different question entirely", neverCalledBuild(t))
	if !errors.Is(err, learningapp.ErrCheckNotInLesson) {
		t.Fatalf("RecordAttempt() error = %v, want it to wrap ErrCheckNotInLesson", err)
	}
	if len(data.attemptsFor(1)) != 0 {
		t.Fatalf("RecordAttempt() inserted a row despite the question not belonging to the lesson")
	}
}

// TestRepositoryRecordAttempt_ScopedByBothSlugs is the RecordAttempt
// analogue of TestRepositoryTransition_ScopedByBothSlugs: a question
// belonging to one topic's lesson must not validate against a
// same-concept-slug lesson under a different topic.
func TestRepositoryRecordAttempt_ScopedByBothSlugs(t *testing.T) {
	data := newStubData()
	data.seedRecallCheck("topic-a", "shared-slug", 101, "Only in topic-a", "short_answer")
	data.seedLesson("topic-b", "shared-slug", 202)
	db := openStubDB(t, data)
	repo := NewRepository(db)

	_, err := repo.RecordAttempt(context.Background(), "topic-b", "shared-slug", "Only in topic-a", neverCalledBuild(t))
	if !errors.Is(err, learningapp.ErrCheckNotInLesson) {
		t.Fatalf("RecordAttempt() error = %v, want it to wrap ErrCheckNotInLesson (must not leak topic-a's question)", err)
	}
}

func TestRepositoryRecordAttempt_Success(t *testing.T) {
	data := newStubData()
	data.seedRecallCheck("ddia", "b-trees", 1, "What is a B-tree?", "mcq")
	db := openStubDB(t, data)
	repo := NewRepository(db)
	selected := "Option B"
	wantKey := mustCheckKey(t, "ddia", "b-trees", "What is a B-tree?")

	entry, err := repo.RecordAttempt(context.Background(), "ddia", "b-trees", "What is a B-tree?",
		buildAttemptFunc(t, "ddia", "b-trees", "confident", "incorrect", &selected, domain.GradedBySelf))
	if err != nil {
		t.Fatalf("RecordAttempt() unexpected error: %v", err)
	}
	if entry.CheckKey != wantKey.String() {
		t.Errorf("CheckKey = %q, want %q", entry.CheckKey, wantKey.String())
	}
	if entry.Kind.String() != "mcq" || entry.Confidence.String() != "confident" || entry.Outcome.String() != "incorrect" {
		t.Errorf("entry = %+v, want kind=mcq confidence=confident outcome=incorrect", entry)
	}
	if entry.SelectedOption == nil || *entry.SelectedOption != "Option B" {
		t.Errorf("SelectedOption = %v, want %q", entry.SelectedOption, "Option B")
	}
	if entry.GradedBy.String() != "self" {
		t.Errorf("GradedBy = %q, want %q", entry.GradedBy.String(), "self")
	}

	stored := data.attemptsFor(1)
	if len(stored) != 1 {
		t.Fatalf("stored attempts = %d, want 1", len(stored))
	}
	if stored[0].checkKey != wantKey.String() || stored[0].kind != "mcq" || stored[0].confidence != "confident" || stored[0].outcome != "incorrect" {
		t.Errorf("stored row = %+v, want it to match the submitted attempt", stored[0])
	}
	if stored[0].selectedOption == nil || *stored[0].selectedOption != "Option B" {
		t.Errorf("stored selected_option = %v, want %q", stored[0].selectedOption, "Option B")
	}
	if stored[0].gradedBy != "self" {
		t.Errorf("stored graded_by = %q, want %q", stored[0].gradedBy, "self")
	}
}

// TestRepositoryRecordAttempt_CaseVariantResolvesToCanonicalCheckKey pins
// R1: recall_checks.question matches under MySQL's case-insensitive
// collation, but CheckKey hashing is byte-exact. A case-variant submission
// that still resolves to the same recall_checks row must hash to the SAME
// check_key as the canonical (as-stored) casing would — never a different
// one derived from the caller's own bytes, which would silently fork that
// question's attempt history.
func TestRepositoryRecordAttempt_CaseVariantResolvesToCanonicalCheckKey(t *testing.T) {
	data := newStubData()
	data.seedRecallCheck("ddia", "b-trees", 1, "What is a B-tree?", "short_answer")
	db := openStubDB(t, data)
	repo := NewRepository(db)
	wantKey := mustCheckKey(t, "ddia", "b-trees", "What is a B-tree?")

	entry, err := repo.RecordAttempt(context.Background(), "ddia", "b-trees", "what is a b-tree?",
		buildAttemptFunc(t, "ddia", "b-trees", "unsure", "correct", nil, domain.GradedBySelf))
	if err != nil {
		t.Fatalf("RecordAttempt() unexpected error: %v", err)
	}
	if entry.CheckKey != wantKey.String() {
		t.Errorf("CheckKey = %q, want %q (the canonical key, not one derived from the case-variant submission)", entry.CheckKey, wantKey.String())
	}

	rawKey := mustCheckKey(t, "ddia", "b-trees", "what is a b-tree?")
	if entry.CheckKey == rawKey.String() {
		t.Errorf("CheckKey matched the raw-submitted-casing hash — must be derived from the canonical stored question instead")
	}
}

// TestRepositoryRecordAttempt_AmbiguousCaseVariantRejected pins the other
// half of R1: MySQL's case/accent-insensitive collation could in principle
// match TWO recall_checks rows in the same lesson that differ only by case —
// content that should never be authored, but must not be silently resolved
// by picking whichever row comes back first. resolveRecallCheck must error
// and build must never run, mirroring lockLesson's own duplicate-match guard.
func TestRepositoryRecordAttempt_AmbiguousCaseVariantRejected(t *testing.T) {
	data := newStubData()
	data.seedRecallCheck("ddia", "b-trees", 1, "What is a B-tree?", "short_answer")
	data.seedRecallCheck("ddia", "b-trees", 1, "what is a b-tree?", "short_answer")
	db := openStubDB(t, data)
	repo := NewRepository(db)

	_, err := repo.RecordAttempt(context.Background(), "ddia", "b-trees", "WHAT IS A B-TREE?", neverCalledBuild(t))
	if err == nil {
		t.Fatal("RecordAttempt() error = nil, want an error when two recall_checks rows collide under case-insensitive collation")
	}
	if errors.Is(err, learningapp.ErrLessonNotFound) || errors.Is(err, learningapp.ErrCheckNotInLesson) {
		t.Fatalf("RecordAttempt() error = %v, want a distinct ambiguous-match error, not the not-found/not-in-lesson ones", err)
	}

	if stored := data.attemptsFor(1); len(stored) != 0 {
		t.Fatalf("RecordAttempt() inserted %d rows despite the ambiguous match, want 0", len(stored))
	}
}

// TestRepositoryRecordAttempt_AppendOnly proves (not just asserts) that the
// same question submitted twice produces two rows, not an upsert collapsing
// to one.
func TestRepositoryRecordAttempt_AppendOnly(t *testing.T) {
	data := newStubData()
	data.seedRecallCheck("ddia", "b-trees", 1, "What is a B-tree?", "short_answer")
	db := openStubDB(t, data)
	repo := NewRepository(db)
	build := buildAttemptFunc(t, "ddia", "b-trees", "unsure", "correct", nil, domain.GradedBySelf)

	if _, err := repo.RecordAttempt(context.Background(), "ddia", "b-trees", "What is a B-tree?", build); err != nil {
		t.Fatalf("first RecordAttempt() unexpected error: %v", err)
	}
	if _, err := repo.RecordAttempt(context.Background(), "ddia", "b-trees", "What is a B-tree?", build); err != nil {
		t.Fatalf("second RecordAttempt() unexpected error: %v", err)
	}

	stored := data.attemptsFor(1)
	if len(stored) != 2 {
		t.Fatalf("stored attempts = %d, want 2 (append-only, not upserted)", len(stored))
	}
}

// TestRepositoryRecordAttempt_GradedByThreadsThrough constructs a
// domain.RecallAttempt with GradedBy=llm directly (bypassing app.Service,
// which today only ever passes GradedBySelf) to prove the repository itself
// stores whatever GradedBy value the attempt carries, rather than a
// hardcoded "self" — see quiz.md's mutation table, mutation #6.
func TestRepositoryRecordAttempt_GradedByThreadsThrough(t *testing.T) {
	data := newStubData()
	data.seedRecallCheck("ddia", "b-trees", 1, "What is a B-tree?", "short_answer")
	db := openStubDB(t, data)
	repo := NewRepository(db)
	llmGraded, err := domain.NewGradedBy("llm")
	if err != nil {
		t.Fatalf("NewGradedBy(llm) failed: %v", err)
	}

	entry, err := repo.RecordAttempt(context.Background(), "ddia", "b-trees", "What is a B-tree?",
		buildAttemptFunc(t, "ddia", "b-trees", "unsure", "correct", nil, llmGraded))
	if err != nil {
		t.Fatalf("RecordAttempt() unexpected error: %v", err)
	}
	if entry.GradedBy.String() != "llm" {
		t.Errorf("GradedBy = %q, want %q", entry.GradedBy.String(), "llm")
	}
	stored := data.attemptsFor(1)
	if len(stored) != 1 || stored[0].gradedBy != "llm" {
		t.Fatalf("stored graded_by = %+v, want \"llm\"", stored)
	}
}

// TestRepositoryRecordAttempt_CreatedAtTruncatedToSeconds pins R2: the
// echoed created_at must be truncated to whole seconds so it always agrees
// with the TIMESTAMP(0) column recall_attempts.created_at actually stores —
// otherwise a sub-second Go timestamp could round differently than the
// truncated value RFC3339 formatting reports, off by up to a second.
func TestRepositoryRecordAttempt_CreatedAtTruncatedToSeconds(t *testing.T) {
	data := newStubData()
	data.seedRecallCheck("ddia", "b-trees", 1, "What is a B-tree?", "short_answer")
	db := openStubDB(t, data)
	repo := NewRepository(db)

	entry, err := repo.RecordAttempt(context.Background(), "ddia", "b-trees", "What is a B-tree?",
		buildAttemptFunc(t, "ddia", "b-trees", "unsure", "correct", nil, domain.GradedBySelf))
	if err != nil {
		t.Fatalf("RecordAttempt() unexpected error: %v", err)
	}
	if entry.CreatedAt.Nanosecond() != 0 {
		t.Errorf("CreatedAt = %v, want truncated to whole seconds (nanosecond component 0)", entry.CreatedAt)
	}

	stored := data.attemptsFor(1)
	if len(stored) != 1 || !stored[0].createdAt.Equal(entry.CreatedAt) {
		t.Fatalf("stored created_at = %v, want it to equal the echoed %v exactly", stored, entry.CreatedAt)
	}
}

// TestRepositoryRecordAttempt_KindResolvedFromRecallChecksType pins R1's
// kind fix: kind is no longer part of Service/Repository.RecordAttempt's
// signature at all, so a client cannot submit one that disagrees with
// recall_checks.type — this is now structurally impossible, not merely
// rejected. This test proves the positive side: whatever kind is stored
// always matches the seeded recall_checks row, regardless of two different
// lessons seeding two different kinds side by side.
func TestRepositoryRecordAttempt_KindResolvedFromRecallChecksType(t *testing.T) {
	data := newStubData()
	data.seedRecallCheck("ddia", "b-trees", 1, "What is a B-tree?", "mcq")
	data.seedRecallCheck("ddia", "another-concept", 2, "A short-answer question", "short_answer")
	db := openStubDB(t, data)
	repo := NewRepository(db)

	mcqEntry, err := repo.RecordAttempt(context.Background(), "ddia", "b-trees", "What is a B-tree?",
		buildAttemptFunc(t, "ddia", "b-trees", "confident", "correct", nil, domain.GradedBySelf))
	if err != nil {
		t.Fatalf("RecordAttempt() unexpected error: %v", err)
	}
	if mcqEntry.Kind.String() != "mcq" {
		t.Errorf("Kind = %q, want %q (resolved from recall_checks.type)", mcqEntry.Kind.String(), "mcq")
	}

	shortAnswerEntry, err := repo.RecordAttempt(context.Background(), "ddia", "another-concept", "A short-answer question",
		buildAttemptFunc(t, "ddia", "another-concept", "confident", "correct", nil, domain.GradedBySelf))
	if err != nil {
		t.Fatalf("RecordAttempt() unexpected error: %v", err)
	}
	if shortAnswerEntry.Kind.String() != "short_answer" {
		t.Errorf("Kind = %q, want %q (resolved from recall_checks.type)", shortAnswerEntry.Kind.String(), "short_answer")
	}

	stored1 := data.attemptsFor(1)
	if len(stored1) != 1 || stored1[0].kind != "mcq" {
		t.Fatalf("stored kind for lesson 1 = %+v, want %q", stored1, "mcq")
	}
	stored2 := data.attemptsFor(2)
	if len(stored2) != 1 || stored2[0].kind != "short_answer" {
		t.Fatalf("stored kind for lesson 2 = %+v, want %q", stored2, "short_answer")
	}
}
