package infra

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	learningapp "github.com/themethaithian/self-learning/internal/learning/app"
	"github.com/themethaithian/self-learning/internal/learning/domain"
	"github.com/themethaithian/self-learning/migrations"
)

func TestRepositoryConceptProgress_LessonNotFound(t *testing.T) {
	db := openStubDB(t, newStubData())
	repo := NewRepository(db)

	_, lessonExists, hasProgress, err := repo.ConceptProgress(context.Background(), "ddia", "b-trees")
	if err != nil {
		t.Fatalf("ConceptProgress() unexpected error: %v", err)
	}
	if lessonExists {
		t.Fatal("lessonExists = true, want false for an unseeded concept")
	}
	if hasProgress {
		t.Fatal("hasProgress = true, want false for an unseeded concept")
	}
}

func TestRepositoryConceptProgress_LessonExistsNoProgress(t *testing.T) {
	data := newStubData()
	data.seedLesson("ddia", "b-trees", 1)
	db := openStubDB(t, data)
	repo := NewRepository(db)

	_, lessonExists, hasProgress, err := repo.ConceptProgress(context.Background(), "ddia", "b-trees")
	if err != nil {
		t.Fatalf("ConceptProgress() unexpected error: %v", err)
	}
	if !lessonExists {
		t.Fatal("lessonExists = false, want true")
	}
	if hasProgress {
		t.Fatal("hasProgress = true, want false: no lesson_progress row written yet")
	}
}

func TestRepositoryConceptProgress_WithProgress(t *testing.T) {
	firstPassed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	lastRead := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	data := newStubData()
	data.seedProgress("ddia", "b-trees", 1, stubProgressRow{state: "passed", firstPassedAt: &firstPassed, lastReadAt: &lastRead})
	db := openStubDB(t, data)
	repo := NewRepository(db)

	entry, lessonExists, hasProgress, err := repo.ConceptProgress(context.Background(), "ddia", "b-trees")
	if err != nil {
		t.Fatalf("ConceptProgress() unexpected error: %v", err)
	}
	if !lessonExists || !hasProgress {
		t.Fatalf("lessonExists=%v hasProgress=%v, want both true", lessonExists, hasProgress)
	}
	if !entry.State.IsPassed() {
		t.Errorf("State = %v, want passed", entry.State)
	}
	if entry.FirstPassedAt == nil || !entry.FirstPassedAt.Equal(firstPassed) {
		t.Errorf("FirstPassedAt = %v, want %v", entry.FirstPassedAt, firstPassed)
	}
	if entry.LastReadAt == nil || !entry.LastReadAt.Equal(lastRead) {
		t.Errorf("LastReadAt = %v, want %v", entry.LastReadAt, lastRead)
	}
}

// TestRepositoryConceptProgress_ScopedByBothSlugs is the concept-slug
// analogue of curriculum's ScopedByBothSlugs test: a concept slug shared by
// two topics must not leak the other topic's progress. Note this does NOT
// catch a dropped "t.slug = ?" predicate in selectConceptProgressSQL
// itself — the stub dispatches on query-string identity, so an edit to the
// SQL text is invisible to this behavioral test either way.
// TestSelectConceptProgressSQLShape is the actual guard for the predicate.
func TestRepositoryConceptProgress_ScopedByBothSlugs(t *testing.T) {
	data := newStubData()
	data.seedProgress("topic-a", "shared-slug", 101, stubProgressRow{state: "passed"})
	data.seedLesson("topic-b", "shared-slug", 202)
	db := openStubDB(t, data)
	repo := NewRepository(db)

	entryB, existsB, hasProgressB, err := repo.ConceptProgress(context.Background(), "topic-b", "shared-slug")
	if err != nil {
		t.Fatalf("ConceptProgress(topic-b) unexpected error: %v", err)
	}
	if !existsB || hasProgressB {
		t.Fatalf("ConceptProgress(topic-b) = exists=%v hasProgress=%v, want exists=true hasProgress=false (never leak topic-a's passed state)", existsB, hasProgressB)
	}
	if entryB.State.IsPassed() {
		t.Error("ConceptProgress(topic-b) leaked topic-a's passed state")
	}

	entryA, existsA, hasProgressA, err := repo.ConceptProgress(context.Background(), "topic-a", "shared-slug")
	if err != nil {
		t.Fatalf("ConceptProgress(topic-a) unexpected error: %v", err)
	}
	if !existsA || !hasProgressA || !entryA.State.IsPassed() {
		t.Fatalf("ConceptProgress(topic-a) = exists=%v hasProgress=%v state=%v, want true/true/passed", existsA, hasProgressA, entryA.State)
	}
}

func TestRepositoryUpsertProgress_CreatesFreshInProgress(t *testing.T) {
	data := newStubData()
	data.seedLesson("ddia", "b-trees", 1)
	db := openStubDB(t, data)
	repo := NewRepository(db)
	ctx := context.Background()

	if err := repo.UpsertProgress(ctx, "ddia", "b-trees", mustChunkState(t, "in_progress")); err != nil {
		t.Fatalf("UpsertProgress() unexpected error: %v", err)
	}

	entry, _, hasProgress, err := repo.ConceptProgress(ctx, "ddia", "b-trees")
	if err != nil {
		t.Fatalf("ConceptProgress() unexpected error: %v", err)
	}
	if !hasProgress || !entry.State.IsInProgress() {
		t.Fatalf("entry = %+v, want hasProgress=true state=in_progress", entry)
	}
	if entry.FirstPassedAt != nil {
		t.Errorf("FirstPassedAt = %v, want nil for a fresh in_progress row", entry.FirstPassedAt)
	}
	if entry.LastReadAt == nil {
		t.Error("LastReadAt = nil, want it set")
	}
}

func TestRepositoryUpsertProgress_CreatesFreshPassed(t *testing.T) {
	data := newStubData()
	data.seedLesson("ddia", "b-trees", 1)
	db := openStubDB(t, data)
	repo := NewRepository(db)
	ctx := context.Background()

	if err := repo.UpsertProgress(ctx, "ddia", "b-trees", mustChunkState(t, "passed")); err != nil {
		t.Fatalf("UpsertProgress() unexpected error: %v", err)
	}

	entry, _, hasProgress, err := repo.ConceptProgress(ctx, "ddia", "b-trees")
	if err != nil {
		t.Fatalf("ConceptProgress() unexpected error: %v", err)
	}
	if !hasProgress || !entry.State.IsPassed() {
		t.Fatalf("entry = %+v, want hasProgress=true state=passed", entry)
	}
	if entry.FirstPassedAt == nil {
		t.Error("FirstPassedAt = nil, want it set on first pass")
	}
}

// TestRepositoryUpsertProgress_PreservesFirstPassedAt pins the COALESCE
// clause in upsertProgressSQL: a re-finish must never move first_passed_at.
func TestRepositoryUpsertProgress_PreservesFirstPassedAt(t *testing.T) {
	firstPassed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	data := newStubData()
	data.seedProgress("ddia", "b-trees", 1, stubProgressRow{state: "passed", firstPassedAt: &firstPassed})
	db := openStubDB(t, data)
	repo := NewRepository(db)
	ctx := context.Background()

	if err := repo.UpsertProgress(ctx, "ddia", "b-trees", mustChunkState(t, "passed")); err != nil {
		t.Fatalf("UpsertProgress() unexpected error: %v", err)
	}

	entry, _, _, err := repo.ConceptProgress(ctx, "ddia", "b-trees")
	if err != nil {
		t.Fatalf("ConceptProgress() unexpected error: %v", err)
	}
	if entry.FirstPassedAt == nil || !entry.FirstPassedAt.Equal(firstPassed) {
		t.Errorf("FirstPassedAt = %v, want unchanged %v", entry.FirstPassedAt, firstPassed)
	}
}

func TestRepositoryUpsertProgress_LessonNotFound(t *testing.T) {
	db := openStubDB(t, newStubData())
	repo := NewRepository(db)

	err := repo.UpsertProgress(context.Background(), "ddia", "b-trees", mustChunkState(t, "in_progress"))
	if !errors.Is(err, learningapp.ErrLessonNotFound) {
		t.Fatalf("UpsertProgress() error = %v, want it to wrap ErrLessonNotFound", err)
	}
}

// TestRepositoryTouchProgress_RefreshesLastReadOnly pins touchProgressSQL's
// entire reason to exist: it must move last_read_at and nothing else.
func TestRepositoryTouchProgress_RefreshesLastReadOnly(t *testing.T) {
	firstPassed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	lastRead := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	data := newStubData()
	data.seedProgress("ddia", "b-trees", 1, stubProgressRow{state: "passed", firstPassedAt: &firstPassed, lastReadAt: &lastRead})
	db := openStubDB(t, data)
	repo := NewRepository(db)
	ctx := context.Background()

	if err := repo.TouchProgress(ctx, "ddia", "b-trees"); err != nil {
		t.Fatalf("TouchProgress() unexpected error: %v", err)
	}

	entry, _, _, err := repo.ConceptProgress(ctx, "ddia", "b-trees")
	if err != nil {
		t.Fatalf("ConceptProgress() unexpected error: %v", err)
	}
	if !entry.State.IsPassed() {
		t.Errorf("State = %v, want unchanged passed", entry.State)
	}
	if entry.FirstPassedAt == nil || !entry.FirstPassedAt.Equal(firstPassed) {
		t.Errorf("FirstPassedAt = %v, want unchanged %v", entry.FirstPassedAt, firstPassed)
	}
	if entry.LastReadAt == nil || entry.LastReadAt.Equal(lastRead) {
		t.Errorf("LastReadAt = %v, want it refreshed away from %v", entry.LastReadAt, lastRead)
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

func TestRepositoryConceptProgress_QueryErrorWraps(t *testing.T) {
	data := newStubData()
	data.queryErr = errors.New("stub: connection refused")
	db := openStubDB(t, data)
	repo := NewRepository(db)

	_, _, _, err := repo.ConceptProgress(context.Background(), "ddia", "b-trees")
	if err == nil {
		t.Fatal("ConceptProgress() expected an error, got nil")
	}
}

func TestRepositoryUpsertProgress_ExecErrorWraps(t *testing.T) {
	data := newStubData()
	data.seedLesson("ddia", "b-trees", 1)
	data.execErr = errors.New("stub: connection refused")
	db := openStubDB(t, data)
	repo := NewRepository(db)

	err := repo.UpsertProgress(context.Background(), "ddia", "b-trees", mustChunkState(t, "in_progress"))
	if err == nil {
		t.Fatal("UpsertProgress() expected an error, got nil")
	}
}

func mustChunkState(t *testing.T, raw string) domain.ChunkState {
	t.Helper()
	s, err := domain.NewChunkState(raw)
	if err != nil {
		t.Fatalf("NewChunkState(%q) failed: %v", raw, err)
	}
	return s
}

func TestSelectConceptProgressSQLShape(t *testing.T) {
	for _, want := range []string{"t.slug = ?", "co.slug = ?", "LEFT JOIN lesson_progress lp"} {
		if !strings.Contains(selectConceptProgressSQL, want) {
			t.Errorf("selectConceptProgressSQL missing %q:\n%s", want, selectConceptProgressSQL)
		}
	}
}

func TestUpsertProgressSQLShape(t *testing.T) {
	for _, want := range []string{
		"INSERT INTO lesson_progress",
		"ON DUPLICATE KEY UPDATE",
		"state = new.state",
		"first_passed_at = COALESCE(lesson_progress.first_passed_at, new.first_passed_at)",
	} {
		if !strings.Contains(upsertProgressSQL, want) {
			t.Errorf("upsertProgressSQL missing %q:\n%s", want, upsertProgressSQL)
		}
	}
}

// TestTouchProgressSQLShape guards touchProgressSQL against ever touching
// state or first_passed_at — a mutation that widened its SET clause would
// silently break the forward-only rule this ticket exists to enforce.
func TestTouchProgressSQLShape(t *testing.T) {
	for _, want := range []string{"t.slug = ?", "co.slug = ?", "SET lp.last_read_at = ?"} {
		if !strings.Contains(touchProgressSQL, want) {
			t.Errorf("touchProgressSQL missing %q:\n%s", want, touchProgressSQL)
		}
	}
	for _, mustNotContain := range []string{"lp.state", "lp.first_passed_at"} {
		if strings.Contains(touchProgressSQL, mustNotContain) {
			t.Errorf("touchProgressSQL contains %q, want it to touch last_read_at only:\n%s", mustNotContain, touchProgressSQL)
		}
	}
}

func TestSelectAllProgressSQLShape(t *testing.T) {
	for _, want := range []string{"FROM lesson_progress lp", "ORDER BY t.slug, co.slug"} {
		if !strings.Contains(selectAllProgressSQL, want) {
			t.Errorf("selectAllProgressSQL missing %q:\n%s", want, selectAllProgressSQL)
		}
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

	stmts := map[string]string{
		"selectConceptProgressSQL": selectConceptProgressSQL,
		"upsertProgressSQL":        upsertProgressSQL,
		"touchProgressSQL":         touchProgressSQL,
		"selectAllProgressSQL":     selectAllProgressSQL,
	}
	for name, stmt := range stmts {
		if !strings.Contains(stmt, "lesson_progress") {
			t.Errorf("%s = %q, want it to reference the lesson_progress table", name, stmt)
		}
	}
}
