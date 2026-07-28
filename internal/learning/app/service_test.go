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
// canned-response fixture, so a bug in Service's decision closure (e.g.
// upserting when it should have only touched last_read_at) shows up as a
// wrong stored value, the same way a real MySQL row would.
type fakeRepository struct {
	lessons map[conceptKey]bool
	rows    map[conceptKey]fakeRow
	all     []ProgressEntry

	transitionErr error
	allErr        error

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

// Transition mirrors infra's own contract: a missing lesson is answered
// directly, never routed through decide, the same way Repository.Transition
// never calls decide for a lesson it never locked.
func (f *fakeRepository) Transition(_ context.Context, topic, concept string, decide func(ProgressEntry, bool) (ProgressDecision, error)) (ProgressEntry, error) {
	if f.transitionErr != nil {
		return ProgressEntry{}, f.transitionErr
	}

	k := conceptKey{topic, concept}
	if !f.lessons[k] {
		return ProgressEntry{}, ErrLessonNotFound
	}
	row, hasProgress := f.rows[k]
	current := ProgressEntry{Topic: topic, Concept: concept}
	if hasProgress {
		current.State, current.FirstPassedAt, current.LastReadAt = row.state, row.firstPassedAt, row.lastReadAt
	}

	decision, err := decide(current, hasProgress)
	if err != nil {
		return ProgressEntry{}, err
	}

	now := time.Now()
	if decision.TouchOnly {
		f.touchCalls++
		row.lastReadAt = &now
	} else {
		f.upsertCalls++
		if decision.State.IsPassed() && row.firstPassedAt == nil {
			row.firstPassedAt = &now
		}
		row.state = decision.State
		row.lastReadAt = &now
	}
	f.rows[k] = row

	return ProgressEntry{
		Topic: topic, Concept: concept,
		State: row.state, FirstPassedAt: row.firstPassedAt, LastReadAt: row.lastReadAt,
	}, nil
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

func chunkStatePtr(t *testing.T, raw string) *domain.ChunkState {
	t.Helper()
	s := mustChunkState(t, raw)
	return &s
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

// TestServiceSetProgress_InvalidSlug pins R1(a): a malformed topic or
// concept must 400 before any repository round trip — never reach a DB
// query that MySQL's ai_ci collation could match against the wrong casing.
func TestServiceSetProgress_InvalidSlug(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		concept string
	}{
		{name: "uppercase concept", topic: "ddia", concept: "B-Trees"},
		{name: "uppercase topic", topic: "DDIA", concept: "b-trees"},
		{name: "space in concept", topic: "ddia", concept: "b trees"},
		{name: "empty topic", topic: "", concept: "b-trees"},
		{name: "empty concept", topic: "ddia", concept: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepository().withLesson(tt.topic, tt.concept)
			svc := NewService(repo)

			_, err := svc.SetProgress(context.Background(), tt.topic, tt.concept, "in_progress")
			if !errors.Is(err, ErrInvalidSlug) {
				t.Fatalf("SetProgress(%q, %q) error = %v, want it to wrap ErrInvalidSlug", tt.topic, tt.concept, err)
			}
			if repo.upsertCalls != 0 || repo.touchCalls != 0 {
				t.Fatalf("SetProgress(%q, %q) reached the repository despite a malformed slug", tt.topic, tt.concept)
			}
		})
	}
}

// TestServiceSetProgress_StateMatrix collapses what used to be five
// near-identical single-scenario tests into one table, so a state this API
// can actually see (locked included — a fresh concept never seeds it, but a
// future gating ticket's rows will) can't go untested by simply not having a
// function for it.
func TestServiceSetProgress_StateMatrix(t *testing.T) {
	firstPassed := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name              string
		seed              *domain.ChunkState
		seedFirstPassedAt *time.Time
		requested         string
		wantState         string
		wantUpserts       int
		wantTouches       int
	}{
		{name: "none -> in_progress creates fresh", seed: nil, requested: "in_progress", wantState: "in_progress", wantUpserts: 1},
		{name: "none -> passed creates fresh", seed: nil, requested: "passed", wantState: "passed", wantUpserts: 1},
		{name: "in_progress -> in_progress is idempotent", seed: chunkStatePtr(t, "in_progress"), requested: "in_progress", wantState: "in_progress", wantTouches: 1},
		{name: "in_progress -> passed is a real transition", seed: chunkStatePtr(t, "in_progress"), requested: "passed", wantState: "passed", wantUpserts: 1},
		{name: "passed -> in_progress never downgrades", seed: chunkStatePtr(t, "passed"), seedFirstPassedAt: &firstPassed, requested: "in_progress", wantState: "passed", wantTouches: 1},
		{name: "passed -> passed is idempotent", seed: chunkStatePtr(t, "passed"), seedFirstPassedAt: &firstPassed, requested: "passed", wantState: "passed", wantTouches: 1},
		{name: "locked -> in_progress unlocks for real", seed: chunkStatePtr(t, "locked"), requested: "in_progress", wantState: "in_progress", wantUpserts: 1},
		{name: "locked -> passed unlocks then passes", seed: chunkStatePtr(t, "locked"), requested: "passed", wantState: "passed", wantUpserts: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepository().withLesson("ddia", "b-trees")
			if tt.seed != nil {
				repo.withProgress("ddia", "b-trees", *tt.seed, tt.seedFirstPassedAt, nil)
			}
			svc := NewService(repo)

			entry, err := svc.SetProgress(context.Background(), "ddia", "b-trees", tt.requested)
			if err != nil {
				t.Fatalf("SetProgress() unexpected error: %v", err)
			}
			if entry.State.String() != tt.wantState {
				t.Errorf("State = %v, want %v", entry.State, tt.wantState)
			}
			if repo.upsertCalls != tt.wantUpserts || repo.touchCalls != tt.wantTouches {
				t.Errorf("upsertCalls=%d touchCalls=%d, want upsertCalls=%d touchCalls=%d", repo.upsertCalls, repo.touchCalls, tt.wantUpserts, tt.wantTouches)
			}
			if entry.LastReadAt == nil {
				t.Error("LastReadAt = nil, want it refreshed on every call")
			}
			if tt.seedFirstPassedAt != nil && (entry.FirstPassedAt == nil || !entry.FirstPassedAt.Equal(*tt.seedFirstPassedAt)) {
				t.Errorf("FirstPassedAt = %v, want unchanged %v", entry.FirstPassedAt, *tt.seedFirstPassedAt)
			}
			if tt.seedFirstPassedAt == nil && tt.wantState == "passed" && entry.FirstPassedAt == nil {
				t.Error("FirstPassedAt = nil, want it set on the first-ever pass")
			}
		})
	}
}

func TestServiceSetProgress_TransitionErrorPropagates(t *testing.T) {
	repoErr := errors.New("connection lost")
	repo := newFakeRepository().withLesson("ddia", "b-trees")
	repo.transitionErr = repoErr
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
