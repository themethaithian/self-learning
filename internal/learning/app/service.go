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

// ErrInvalidSlug is returned when topicSlug or conceptSlug is not
// slug-shaped. This must be checked before any DB round trip: MySQL's
// ai_ci collation matches e.g. "B-Trees" to a stored "b-trees" row, so a
// malformed-but-DB-matching slug would otherwise reach domain.NewLessonRef
// downstream and fail there instead — a 500, not a 400.
var ErrInvalidSlug = errors.New("learning: invalid topic or concept slug")

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

// ProgressDecision is what a Transition callback decided to persist:
// TouchOnly refreshes last_read_at and leaves State/FirstPassedAt exactly as
// they are; otherwise State is written (and infra stamps first_passed_at the
// first time State is passed).
type ProgressDecision struct {
	State     domain.ChunkState
	TouchOnly bool
}

// Repository is the learning read/write port. infra provides the MySQL
// adapter; tests provide a fake.
type Repository interface {
	// Transition locks topicSlug/conceptSlug's row for one transaction and
	// calls decide with its current progress, so a decision can never be
	// made from a stale read — e.g. two racing PUTs, one from opening a
	// lesson and one from Finish, always see whatever the other one already
	// committed instead of interleaving into a state/first_passed_at
	// contradiction. lessonExists is false when there is no lesson row;
	// decide must then return an error, since there is nothing to write.
	// hasProgress is false when the lesson exists but no lesson_progress row
	// has been written yet.
	Transition(ctx context.Context, topicSlug, conceptSlug string, decide func(current ProgressEntry, lessonExists, hasProgress bool) (ProgressDecision, error)) (ProgressEntry, error)

	AllProgress(ctx context.Context) ([]ProgressEntry, error)
}

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
// domain.LessonProgress's: decide below calls Unlock/MarkPassed for real and
// relies on alreadyAtOrPastRequestedState to turn their errors into a
// touch-only no-op, rather than weakening the domain invariants to make
// HTTP convenient.
func (s Service) SetProgress(ctx context.Context, topicSlug, conceptSlug, rawState string) (ProgressEntry, error) {
	requested, err := parseRequestedState(rawState)
	if err != nil {
		return ProgressEntry{}, err
	}
	if err := validateTopicSlugShape(topicSlug); err != nil {
		return ProgressEntry{}, err
	}
	ref, err := newConceptRef(conceptSlug)
	if err != nil {
		return ProgressEntry{}, err
	}

	entry, err := s.repo.Transition(ctx, topicSlug, conceptSlug, func(current ProgressEntry, lessonExists, hasProgress bool) (ProgressDecision, error) {
		if !lessonExists {
			return ProgressDecision{}, ErrLessonNotFound
		}
		if !hasProgress {
			return ProgressDecision{State: requested}, nil
		}
		return decideTransition(ref, current.State, requested)
	})
	if err != nil {
		return ProgressEntry{}, fmt.Errorf("learning: set progress %s/%s: %w", topicSlug, conceptSlug, err)
	}
	return entry, nil
}

// decideTransition is the only place domain.LessonProgress's transition
// methods are called for an existing row. A requested passed unlocks first
// (ignoring ErrAlreadyUnlocked) because MarkPassed alone rejects a locked
// lesson outright — this is the sole way a locked row can ever reach passed.
func decideTransition(ref domain.LessonRef, currentState, requested domain.ChunkState) (ProgressDecision, error) {
	progress, err := domain.NewLessonProgress(ref, currentState)
	if err != nil {
		return ProgressDecision{}, err
	}

	if !requested.IsPassed() {
		next, err := progress.Unlock()
		if err != nil {
			if alreadyAtOrPastRequestedState(err) {
				return ProgressDecision{TouchOnly: true}, nil
			}
			return ProgressDecision{}, err
		}
		return ProgressDecision{State: next.State()}, nil
	}

	if unlocked, err := progress.Unlock(); err == nil {
		progress = unlocked
	} else if !errors.Is(err, domain.ErrAlreadyUnlocked) {
		return ProgressDecision{}, err
	}

	next, err := progress.MarkPassed()
	if err != nil {
		if alreadyAtOrPastRequestedState(err) {
			return ProgressDecision{TouchOnly: true}, nil
		}
		return ProgressDecision{}, err
	}
	return ProgressDecision{State: next.State()}, nil
}

func alreadyAtOrPastRequestedState(err error) bool {
	return errors.Is(err, domain.ErrAlreadyUnlocked) || errors.Is(err, domain.ErrAlreadyPassed)
}

func parseRequestedState(raw string) (domain.ChunkState, error) {
	if raw != "in_progress" && raw != "passed" {
		return domain.ChunkState{}, fmt.Errorf("learning: state %q: %w", raw, ErrInvalidState)
	}
	return domain.NewChunkState(raw)
}

// validateTopicSlugShape borrows domain.LessonRef's shape check (mirrors
// curriculum.Slug's rules) since this bounded context has no separate topic
// slug type — a topic slug and a concept slug are the same shape.
func validateTopicSlugShape(topicSlug string) error {
	if _, err := domain.NewLessonRef(topicSlug); err != nil {
		return fmt.Errorf("%w: %q", ErrInvalidSlug, topicSlug)
	}
	return nil
}

func newConceptRef(conceptSlug string) (domain.LessonRef, error) {
	ref, err := domain.NewLessonRef(conceptSlug)
	if err != nil {
		return domain.LessonRef{}, fmt.Errorf("%w: %q", ErrInvalidSlug, conceptSlug)
	}
	return ref, nil
}
