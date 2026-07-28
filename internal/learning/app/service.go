// Package app orchestrates the learning bounded context. It imports only
// domain and stdlib; infra depends on it, never the other way around.
package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/themethaithian/self-learning/internal/learning/domain"
)

// ErrLessonNotFound is returned by Service.SetProgress when topicSlug and
// conceptSlug have no lesson yet — the caller must answer 404, never write a
// phantom lesson_progress row.
var ErrLessonNotFound = errors.New("learning: lesson not found")

// ErrInvalidState is returned when the requested state is not one this API
// accepts. "locked" is a valid domain.ChunkState but is never settable
// through this API: the product decision is soft-guide, not gating, so
// domain.Gate stays unused here and no caller can ever write "locked".
var ErrInvalidState = errors.New("learning: invalid requested state")

// ProgressEntry is one concept's persisted progress, addressed by topic and
// concept slug — the same identity GET /api/v1/lessons/{topic}/{concept}
// already uses, never a lesson_id the curriculum API never exposes.
type ProgressEntry struct {
	Topic         string
	Concept       string
	State         domain.ChunkState
	FirstPassedAt *time.Time
	LastReadAt    *time.Time
}

// Repository is the learning read/write port. infra provides the MySQL
// adapter; tests provide a fake.
type Repository interface {
	// ConceptProgress resolves one concept's progress. lessonExists is false
	// when topicSlug/conceptSlug has no lesson row. hasProgress is false when
	// the lesson exists but no lesson_progress row has been written yet.
	ConceptProgress(ctx context.Context, topicSlug, conceptSlug string) (entry ProgressEntry, lessonExists, hasProgress bool, err error)

	// UpsertProgress writes state for topicSlug/conceptSlug, creating the row
	// on first write. Infra sets first_passed_at the first time state is
	// passed and never overwrites it afterward.
	UpsertProgress(ctx context.Context, topicSlug, conceptSlug string, state domain.ChunkState) error

	// TouchProgress refreshes last_read_at only, leaving state and
	// first_passed_at untouched — the idempotent / forward-only no-op path.
	TouchProgress(ctx context.Context, topicSlug, conceptSlug string) error

	// AllProgress returns every concept that has a lesson_progress row.
	AllProgress(ctx context.Context) ([]ProgressEntry, error)
}

// Service is the learning use-case layer.
type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) ListProgress(ctx context.Context) ([]ProgressEntry, error) {
	entries, err := s.repo.AllProgress(ctx)
	if err != nil {
		return nil, fmt.Errorf("learning: list progress: %w", err)
	}
	return entries, nil
}

// SetProgress moves one concept's progress toward rawState ("in_progress" or
// "passed"). Idempotency and the forward-only rule are this layer's job, not
// domain.LessonProgress's: Unlock and MarkPassed deliberately error when the
// lesson is already past the state they transition into, and that is
// exactly "already at (or past) the requested state" — this method catches
// ErrAlreadyUnlocked/ErrAlreadyPassed and turns them into a successful
// no-op that still refreshes last_read_at, rather than weakening the domain
// invariants to make HTTP convenient. This is also why a passed lesson never
// downgrades on a later "in_progress" PUT: Unlock only succeeds from locked,
// so it errors identically whether the current state is in_progress or
// passed, and the no-op path never changes state either way.
func (s Service) SetProgress(ctx context.Context, topicSlug, conceptSlug, rawState string) (ProgressEntry, error) {
	requested, err := parseRequestedState(rawState)
	if err != nil {
		return ProgressEntry{}, err
	}

	current, lessonExists, hasProgress, err := s.repo.ConceptProgress(ctx, topicSlug, conceptSlug)
	if err != nil {
		return ProgressEntry{}, fmt.Errorf("learning: set progress %s/%s: %w", topicSlug, conceptSlug, err)
	}
	if !lessonExists {
		return ProgressEntry{}, fmt.Errorf("learning: set progress %s/%s: %w", topicSlug, conceptSlug, ErrLessonNotFound)
	}

	ref, err := domain.NewLessonRef(conceptSlug)
	if err != nil {
		return ProgressEntry{}, fmt.Errorf("learning: set progress %s/%s: %w", topicSlug, conceptSlug, err)
	}

	if !hasProgress {
		if _, err := domain.NewLessonProgress(ref, requested); err != nil {
			return ProgressEntry{}, fmt.Errorf("learning: set progress %s/%s: %w", topicSlug, conceptSlug, err)
		}
		if err := s.repo.UpsertProgress(ctx, topicSlug, conceptSlug, requested); err != nil {
			return ProgressEntry{}, fmt.Errorf("learning: set progress %s/%s: %w", topicSlug, conceptSlug, err)
		}
		return s.reload(ctx, topicSlug, conceptSlug)
	}

	progress, err := domain.NewLessonProgress(ref, current.State)
	if err != nil {
		return ProgressEntry{}, fmt.Errorf("learning: set progress %s/%s: %w", topicSlug, conceptSlug, err)
	}

	var next domain.LessonProgress
	if requested.IsPassed() {
		next, err = progress.MarkPassed()
	} else {
		next, err = progress.Unlock()
	}
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyPassed) || errors.Is(err, domain.ErrAlreadyUnlocked) {
			if err := s.repo.TouchProgress(ctx, topicSlug, conceptSlug); err != nil {
				return ProgressEntry{}, fmt.Errorf("learning: set progress %s/%s: %w", topicSlug, conceptSlug, err)
			}
			return s.reload(ctx, topicSlug, conceptSlug)
		}
		return ProgressEntry{}, fmt.Errorf("learning: set progress %s/%s: %w", topicSlug, conceptSlug, err)
	}

	if err := s.repo.UpsertProgress(ctx, topicSlug, conceptSlug, next.State()); err != nil {
		return ProgressEntry{}, fmt.Errorf("learning: set progress %s/%s: %w", topicSlug, conceptSlug, err)
	}
	return s.reload(ctx, topicSlug, conceptSlug)
}

func (s Service) reload(ctx context.Context, topicSlug, conceptSlug string) (ProgressEntry, error) {
	entry, _, _, err := s.repo.ConceptProgress(ctx, topicSlug, conceptSlug)
	if err != nil {
		return ProgressEntry{}, fmt.Errorf("learning: reload progress %s/%s: %w", topicSlug, conceptSlug, err)
	}
	return entry, nil
}

// parseRequestedState accepts only the two states this API ever writes.
func parseRequestedState(raw string) (domain.ChunkState, error) {
	if raw != "in_progress" && raw != "passed" {
		return domain.ChunkState{}, fmt.Errorf("learning: state %q: %w", raw, ErrInvalidState)
	}
	return domain.NewChunkState(raw)
}
