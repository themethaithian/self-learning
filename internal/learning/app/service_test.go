package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/themethaithian/self-learning/internal/learning/domain"
)

type conceptKey struct{ topic, concept string }

type fakeRow struct {
	state         domain.ChunkState
	firstPassedAt *time.Time
	lastReadAt    *time.Time
}

// fakeRepository is a real in-memory table keyed by (topic, concept), not a
// canned-response fixture, so a bug in Service's branching (e.g. calling
// UpsertProgress when it should have called TouchProgress) shows up as a
// wrong stored value, the same way a real MySQL row would.
type fakeRepository struct {
	lessons map[conceptKey]bool
	rows    map[conceptKey]fakeRow
	all     []ProgressEntry

	conceptProgressErr error
	upsertErr          error
	touchErr           error
	allErr             error

	upsertCalls int
	touchCalls  int
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{lessons: map[conceptKey]bool{}, rows: map[conceptKey]fakeRow{}}
}

func (f *fakeRepository) withLesson(topic, concept string) *fakeRepository {
	f.lessons[conceptKey{topic, concept}] = true
	return f
}

func (f *fakeRepository) withProgress(topic, concept string, state domain.ChunkState, firstPassedAt, lastReadAt *time.Time) *fakeRepository {
	f.lessons[conceptKey{topic, concept}] = true
	f.rows[conceptKey{topic, concept}] = fakeRow{state: state, firstPassedAt: firstPassedAt, lastReadAt: lastReadAt}
	return f
}

func (f *fakeRepository) ConceptProgress(_ context.Context, topic, concept string) (ProgressEntry, bool, bool, error) {
	if f.conceptProgressErr != nil {
		return ProgressEntry{}, false, false, f.conceptProgressErr
	}
	k := conceptKey{topic, concept}
	if !f.lessons[k] {
		return ProgressEntry{}, false, false, nil
	}
	row, ok := f.rows[k]
	if !ok {
		return ProgressEntry{Topic: topic, Concept: concept}, true, false, nil
	}
	return ProgressEntry{
		Topic: topic, Concept: concept,
		State: row.state, FirstPassedAt: row.firstPassedAt, LastReadAt: row.lastReadAt,
	}, true, true, nil
}

func (f *fakeRepository) UpsertProgress(_ context.Context, topic, concept string, state domain.ChunkState) error {
	f.upsertCalls++
	if f.upsertErr != nil {
		return f.upsertErr
	}
	k := conceptKey{topic, concept}
	now := time.Now()
	row := f.rows[k]
	row.state = state
	if state.IsPassed() && row.firstPassedAt == nil {
		row.firstPassedAt = &now
	}
	row.lastReadAt = &now
	f.rows[k] = row
	return nil
}

func (f *fakeRepository) TouchProgress(_ context.Context, topic, concept string) error {
	f.touchCalls++
	if f.touchErr != nil {
		return f.touchErr
	}
	k := conceptKey{topic, concept}
	row := f.rows[k]
	now := time.Now()
	row.lastReadAt = &now
	f.rows[k] = row
	return nil
}

func (f *fakeRepository) AllProgress(context.Context) ([]ProgressEntry, error) {
	if f.allErr != nil {
		return nil, f.allErr
	}
	return f.all, nil
}

func mustChunkState(t *testing.T, raw string) domain.ChunkState {
	t.Helper()
	s, err := domain.NewChunkState(raw)
	if err != nil {
		t.Fatalf("NewChunkState(%q) failed: %v", raw, err)
	}
	return s
}

func TestServiceSetProgress_LessonNotFound(t *testing.T) {
	svc := NewService(newFakeRepository())

	_, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if !errors.Is(err, ErrLessonNotFound) {
		t.Fatalf("SetProgress() error = %v, want it to wrap ErrLessonNotFound", err)
	}
}

func TestServiceSetProgress_InvalidState(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{name: "locked is a domain state but never settable", raw: "locked"},
		{name: "empty", raw: ""},
		{name: "unknown", raw: "skipped"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepository().withLesson("ddia", "b-trees")
			svc := NewService(repo)

			_, err := svc.SetProgress(context.Background(), "ddia", "b-trees", tt.raw)
			if !errors.Is(err, ErrInvalidState) {
				t.Fatalf("SetProgress(%q) error = %v, want it to wrap ErrInvalidState", tt.raw, err)
			}
			if repo.upsertCalls != 0 || repo.touchCalls != 0 {
				t.Fatalf("SetProgress(%q) touched the repository despite invalid state", tt.raw)
			}
		})
	}
}

func TestServiceSetProgress_FreshConceptInProgress(t *testing.T) {
	repo := newFakeRepository().withLesson("ddia", "b-trees")
	svc := NewService(repo)

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsInProgress() {
		t.Errorf("State = %v, want in_progress", entry.State)
	}
	if entry.FirstPassedAt != nil {
		t.Errorf("FirstPassedAt = %v, want nil for a fresh in_progress concept", entry.FirstPassedAt)
	}
	if entry.LastReadAt == nil {
		t.Error("LastReadAt = nil, want it set")
	}
}

func TestServiceSetProgress_FreshConceptPassed(t *testing.T) {
	repo := newFakeRepository().withLesson("ddia", "b-trees")
	svc := NewService(repo)

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "passed")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsPassed() {
		t.Errorf("State = %v, want passed", entry.State)
	}
	if entry.FirstPassedAt == nil {
		t.Error("FirstPassedAt = nil, want it set for a lesson passed on first write")
	}
}

func TestServiceSetProgress_TransitionInProgressToPassed(t *testing.T) {
	repo := newFakeRepository().withProgress("ddia", "b-trees", mustChunkState(t, "in_progress"), nil, nil)
	svc := NewService(repo)

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "passed")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsPassed() {
		t.Errorf("State = %v, want passed", entry.State)
	}
	if entry.FirstPassedAt == nil {
		t.Error("FirstPassedAt = nil, want it set on the real in_progress -> passed transition")
	}
	if repo.upsertCalls != 1 || repo.touchCalls != 0 {
		t.Errorf("upsertCalls=%d touchCalls=%d, want a real Upsert not a Touch", repo.upsertCalls, repo.touchCalls)
	}
}

func TestServiceSetProgress_IdempotentInProgress(t *testing.T) {
	repo := newFakeRepository().withProgress("ddia", "b-trees", mustChunkState(t, "in_progress"), nil, nil)
	svc := NewService(repo)

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsInProgress() {
		t.Errorf("State = %v, want it to stay in_progress", entry.State)
	}
	if repo.upsertCalls != 0 || repo.touchCalls != 1 {
		t.Errorf("upsertCalls=%d touchCalls=%d, want a no-op Touch not an Upsert", repo.upsertCalls, repo.touchCalls)
	}
}

// TestServiceSetProgress_ForwardOnlyNeverDowngrades is the ticket's core
// invariant: PUTting in_progress on an already-passed lesson must not move
// it backward, and must not disturb first_passed_at.
func TestServiceSetProgress_ForwardOnlyNeverDowngrades(t *testing.T) {
	firstPassed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	repo := newFakeRepository().withProgress("ddia", "b-trees", mustChunkState(t, "passed"), &firstPassed, nil)
	svc := NewService(repo)

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsPassed() {
		t.Errorf("State = %v, want it to stay passed (forward-only)", entry.State)
	}
	if entry.FirstPassedAt == nil || !entry.FirstPassedAt.Equal(firstPassed) {
		t.Errorf("FirstPassedAt = %v, want unchanged %v", entry.FirstPassedAt, firstPassed)
	}
	if entry.LastReadAt == nil {
		t.Error("LastReadAt = nil, want it refreshed even on the no-op path")
	}
	if repo.upsertCalls != 0 || repo.touchCalls != 1 {
		t.Errorf("upsertCalls=%d touchCalls=%d, want a no-op Touch not an Upsert", repo.upsertCalls, repo.touchCalls)
	}
}

func TestServiceSetProgress_IdempotentPassed(t *testing.T) {
	firstPassed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	repo := newFakeRepository().withProgress("ddia", "b-trees", mustChunkState(t, "passed"), &firstPassed, nil)
	svc := NewService(repo)

	entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "passed")
	if err != nil {
		t.Fatalf("SetProgress() unexpected error: %v", err)
	}
	if !entry.State.IsPassed() {
		t.Errorf("State = %v, want passed", entry.State)
	}
	if entry.FirstPassedAt == nil || !entry.FirstPassedAt.Equal(firstPassed) {
		t.Errorf("FirstPassedAt = %v, want unchanged %v (never overwritten on re-finish)", entry.FirstPassedAt, firstPassed)
	}
	if repo.upsertCalls != 0 || repo.touchCalls != 1 {
		t.Errorf("upsertCalls=%d touchCalls=%d, want a no-op Touch not an Upsert", repo.upsertCalls, repo.touchCalls)
	}
}

func TestServiceSetProgress_ConceptProgressErrorPropagates(t *testing.T) {
	repoErr := errors.New("connection lost")
	repo := newFakeRepository()
	repo.conceptProgressErr = repoErr
	svc := NewService(repo)

	_, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if !errors.Is(err, repoErr) {
		t.Fatalf("SetProgress() error = %v, want it to wrap %v", err, repoErr)
	}
}

func TestServiceSetProgress_UpsertErrorPropagates(t *testing.T) {
	repoErr := errors.New("connection lost")
	repo := newFakeRepository().withLesson("ddia", "b-trees")
	repo.upsertErr = repoErr
	svc := NewService(repo)

	_, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if !errors.Is(err, repoErr) {
		t.Fatalf("SetProgress() error = %v, want it to wrap %v", err, repoErr)
	}
}

func TestServiceSetProgress_TouchErrorPropagates(t *testing.T) {
	repoErr := errors.New("connection lost")
	repo := newFakeRepository().withProgress("ddia", "b-trees", mustChunkState(t, "in_progress"), nil, nil)
	repo.touchErr = repoErr
	svc := NewService(repo)

	_, err := svc.SetProgress(context.Background(), "ddia", "b-trees", "in_progress")
	if !errors.Is(err, repoErr) {
		t.Fatalf("SetProgress() error = %v, want it to wrap %v", err, repoErr)
	}
}

func TestServiceListProgress(t *testing.T) {
	want := []ProgressEntry{{Topic: "ddia", Concept: "b-trees", State: mustChunkState(t, "in_progress")}}
	repo := newFakeRepository()
	repo.all = want
	svc := NewService(repo)

	got, err := svc.ListProgress(context.Background())
	if err != nil {
		t.Fatalf("ListProgress() unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Topic != "ddia" || got[0].Concept != "b-trees" {
		t.Fatalf("ListProgress() = %+v, want %+v", got, want)
	}
}

func TestServiceListProgress_ErrorPropagates(t *testing.T) {
	repoErr := errors.New("connection lost")
	repo := newFakeRepository()
	repo.allErr = repoErr
	svc := NewService(repo)

	_, err := svc.ListProgress(context.Background())
	if !errors.Is(err, repoErr) {
		t.Fatalf("ListProgress() error = %v, want it to wrap %v", err, repoErr)
	}
}
