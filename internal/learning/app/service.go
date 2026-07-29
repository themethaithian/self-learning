// Package app orchestrates the learning bounded context. It imports only
// domain and stdlib; infra depends on it, never the other way around.
package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

// ErrCheckNotInLesson is returned by RecordAttempt when topicSlug/conceptSlug
// resolve to a real lesson but question does not match any of that lesson's
// recall_checks rows. A client-side typo or a stale cached question must
// 400 here rather than silently minting a check_key no future SRS card will
// ever match (see domain.CheckKey's doc comment).
var ErrCheckNotInLesson = errors.New("learning: recall check not found in this lesson")

// ErrInvalidAttempt is returned when confidence, outcome, question
// non-emptiness, or the selected-option/kind pairing fails domain
// validation. confidence/outcome/question are checked before any repository
// round trip, the same way ErrInvalidSlug is; the selected-option/kind
// pairing can only be checked once kind is resolved from recall_checks.type,
// inside RecordAttempt's build closure.
var ErrInvalidAttempt = errors.New("learning: invalid recall attempt")

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

// AttemptRecord is one persisted recall_attempts row, echoed back to the
// caller exactly as stored. Question is the CANONICAL text CheckKey was
// hashed from, not necessarily what the caller submitted — echoing it lets a
// client that sent a case- or accent-variant question notice it was
// normalised, instead of quietly assuming CheckKey matches its own bytes.
type AttemptRecord struct {
	CheckKey       string
	Question       string
	Kind           domain.CheckKind
	Confidence     domain.Confidence
	Outcome        domain.AttemptOutcome
	SelectedOption *string
	GradedBy       domain.GradedBy
	CreatedAt      time.Time
}

// Repository is the learning read/write port. infra provides the MySQL
// adapter; tests provide a fake.
type Repository interface {
	// Transition calls decide with topicSlug/conceptSlug's current progress.
	// decide only ever runs once the lesson is confirmed to exist — a
	// missing lesson is Transition's own ErrLessonNotFound, never routed
	// through decide. hasProgress is false when the lesson exists but no
	// lesson_progress row has been written yet.
	Transition(ctx context.Context, topicSlug, conceptSlug string, decide func(current ProgressEntry, hasProgress bool) (ProgressDecision, error)) (ProgressEntry, error)

	AllProgress(ctx context.Context) ([]ProgressEntry, error)

	// RecordAttempt resolves the recall check identified by (topicSlug,
	// conceptSlug, question), then calls build with that check's CANONICAL
	// question and its real kind — never question itself, and never a
	// client-supplied kind. Canonical question matters because
	// recall_checks.question is matched under MySQL's case/accent-
	// insensitive collation, but domain.CheckKey hashes byte-exact: hashing
	// the caller's raw question instead of the canonical one would let a
	// case-variant phrasing pass the lesson-membership check yet mint a
	// DIFFERENT check_key than the canonical text would, silently forking
	// one question's attempt history (see domain.CanonicalQuestion's doc
	// comment). kind is server-resolved for the same reason CheckKey is:
	// it is static curriculum content the server already owns, not a
	// client judgement like outcome is — a trusted client-supplied kind
	// could store an mcq's selected_option against what the database says
	// is a short_answer check. It rejects a question that does not belong
	// to that lesson (ErrCheckNotInLesson) or a lesson that does not exist
	// (ErrLessonNotFound) without ever calling build. If build itself
	// errors (e.g. domain.NewRecallAttempt's mcq-only invariant), that
	// error is returned with nothing inserted — mirroring how Transition
	// never writes when decide errors. Attempts are append-only, so this is
	// a resolve-then-insert, never a Transition-style locked
	// read-modify-write.
	RecordAttempt(ctx context.Context, topicSlug, conceptSlug, question string, build func(canonicalQuestion domain.CanonicalQuestion, kind domain.CheckKind) (domain.RecallAttempt, error)) (AttemptRecord, error)
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

	entry, err := s.repo.Transition(ctx, topicSlug, conceptSlug, func(current ProgressEntry, hasProgress bool) (ProgressDecision, error) {
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

func validateTopicSlugShape(topicSlug string) error {
	if !domain.IsValidSlugShape(topicSlug) {
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

// RecordAttempt persists one self-graded recall-check attempt. confidence,
// outcome, and question-non-emptiness are all validated here, before any
// repository round trip. kind and the CheckKey cannot be validated here: kind
// is resolved server-side from recall_checks.type and CheckKey must hash the
// DB's own canonical question text, not rawQuestion — a case-variant
// question can pass "belongs to this lesson" under MySQL's collation yet
// hash differently from the canonical text (see domain.CanonicalQuestion's
// doc comment) — so build defers both until the repository hands back what
// it actually matched. trimmedQuestion (not rawQuestion) is what gets sent
// to the repository, so the "belongs to this lesson" lookup compares the
// same text CheckKey will hash, not a version that still carries whitespace
// domain.NewCanonicalQuestion would have trimmed away.
func (s Service) RecordAttempt(ctx context.Context, topicSlug, conceptSlug, rawQuestion, rawConfidence, rawOutcome string, rawSelectedOption *string) (AttemptRecord, error) {
	if err := validateTopicSlugShape(topicSlug); err != nil {
		return AttemptRecord{}, err
	}
	conceptRef, err := newConceptRef(conceptSlug)
	if err != nil {
		return AttemptRecord{}, err
	}

	confidence, err := domain.NewConfidence(rawConfidence)
	if err != nil {
		return AttemptRecord{}, fmt.Errorf("%w: %w", ErrInvalidAttempt, err)
	}
	outcome, err := domain.NewAttemptOutcome(rawOutcome)
	if err != nil {
		return AttemptRecord{}, fmt.Errorf("%w: %w", ErrInvalidAttempt, err)
	}
	trimmedQuestion := strings.TrimSpace(rawQuestion)
	if trimmedQuestion == "" {
		return AttemptRecord{}, fmt.Errorf("%w: %w", ErrInvalidAttempt, domain.ErrInvalidCheckKey)
	}
	selectedOption := normalizeSelectedOption(rawSelectedOption)

	build := func(canonicalQuestion domain.CanonicalQuestion, kind domain.CheckKind) (domain.RecallAttempt, error) {
		checkKey, err := domain.NewCheckKey(topicSlug, conceptRef, canonicalQuestion)
		if err != nil {
			return domain.RecallAttempt{}, fmt.Errorf("%w: %w", ErrInvalidAttempt, err)
		}
		attempt, err := domain.NewRecallAttempt(checkKey, kind, confidence, outcome, selectedOption, domain.GradedBySelf)
		if err != nil {
			return domain.RecallAttempt{}, fmt.Errorf("%w: %w", ErrInvalidAttempt, err)
		}
		return attempt, nil
	}

	entry, err := s.repo.RecordAttempt(ctx, topicSlug, conceptSlug, trimmedQuestion, build)
	if err != nil {
		return AttemptRecord{}, fmt.Errorf("learning: record attempt %s/%s: %w", topicSlug, conceptSlug, err)
	}
	return entry, nil
}

// normalizeSelectedOption treats an empty (post-trim) selected option the
// same as an absent one, so a client sending "" behaves identically to
// omitting the field rather than tripping the mcq-only invariant on
// meaningless input.
func normalizeSelectedOption(raw *string) *string {
	if raw == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
