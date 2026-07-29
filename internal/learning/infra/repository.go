// Package infra adapts the learning domain and app layers to MySQL and
// HTTP — the vendor-facing details the domain and app layers must never see.
package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	learningapp "github.com/themethaithian/self-learning/internal/learning/app"
	"github.com/themethaithian/self-learning/internal/learning/domain"
)

const (
	// selectLessonForUpdateSQL locks the lessons row for the rest of the
	// transaction, serializing every concurrent Transition for this concept
	// through one writer at a time — two racing PUTs (e.g. one from opening
	// a lesson, one from Finish) can never both decide from the same stale
	// read. "FOR UPDATE OF l" (not bare FOR UPDATE, which would lock every
	// joined table — every chapter row of the topic, via concepts/chapters)
	// scopes the lock to lessons alone, so two PUTs on different concepts in
	// the same topic never serialize against each other. It also returns
	// the DB's own slugs, not the caller's, so case-insensitive collation
	// can never surface the caller's casing back out of a response.
	selectLessonForUpdateSQL = `
SELECT l.id, t.slug, co.slug
FROM lessons l
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
WHERE t.slug = ? AND co.slug = ?
FOR UPDATE OF l`

	selectProgressByLessonIDSQL = `
SELECT state, first_passed_at, last_read_at
FROM lesson_progress
WHERE lesson_id = ?`

	// upsertProgressByLessonIDSQL writes state on first write or on a real
	// transition. first_passed_at is preserved via COALESCE against the
	// table's own (pre-update) value, never the incoming one, so a caller
	// can pass "now" every time state is passed and still never move an
	// already-set first_passed_at.
	upsertProgressByLessonIDSQL = `
INSERT INTO lesson_progress (lesson_id, state, first_passed_at, last_read_at)
VALUES (?, ?, ?, ?)
AS new
ON DUPLICATE KEY UPDATE
    state = new.state,
    first_passed_at = COALESCE(lesson_progress.first_passed_at, new.first_passed_at),
    last_read_at = new.last_read_at`

	refreshLastReadAtSQL = `
UPDATE lesson_progress
SET last_read_at = ?
WHERE lesson_id = ?`

	// selectAllProgressSQL: absence means "not started", per the ticket's
	// response contract.
	selectAllProgressSQL = `
SELECT t.slug, co.slug, lp.state, lp.first_passed_at, lp.last_read_at
FROM lesson_progress lp
JOIN lessons l ON l.id = lp.lesson_id
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
ORDER BY t.slug, co.slug`

	// selectRecallCheckExistsSQL is how RecordAttempt enforces "the question
	// belongs to this lesson": the same lessons/concepts/chapters/topics
	// join selectLessonForUpdateSQL uses, extended with recall_checks so a
	// client-supplied question that does not match any recall_checks row
	// for this exact (topic, concept) never inserts a recall_attempts row
	// (see domain.CheckKey's doc comment on why an orphan check_key is
	// unrecoverable). It also returns rc.question itself: recall_checks.question
	// matches under MySQL's case/accent-insensitive collation, so the row
	// this WHERE clause matched can have different bytes than what the
	// caller sent — the caller must hash THIS returned text, never its own
	// input, or a case-variant question would mint a different check_key
	// than the canonical one (see domain.CheckKey's doc comment).
	selectRecallCheckExistsSQL = `
SELECT l.id, rc.question
FROM recall_checks rc
JOIN lessons l ON l.id = rc.lesson_id
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
WHERE t.slug = ? AND co.slug = ? AND rc.question = ?`

	// selectLessonExistsSQL disambiguates selectRecallCheckExistsSQL's zero
	// rows: this lesson not existing at all (ErrLessonNotFound, 404) versus
	// the lesson existing but the question not belonging to it
	// (ErrCheckNotInLesson, 400).
	selectLessonExistsSQL = `
SELECT l.id
FROM lessons l
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
WHERE t.slug = ? AND co.slug = ?`

	// insertRecallAttemptSQL is append-only by design: no ON DUPLICATE KEY
	// UPDATE and no unique constraint on check_key, so many rows can
	// accumulate per check over time (Q-2c reads them back newest-first).
	insertRecallAttemptSQL = `
INSERT INTO recall_attempts (lesson_id, check_key, type, confidence, outcome, selected_option, graded_by, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
)

// Repository is the MySQL adapter for app.Repository.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Transition(
	ctx context.Context,
	topicSlug, conceptSlug string,
	decide func(current learningapp.ProgressEntry, hasProgress bool) (learningapp.ProgressDecision, error),
) (learningapp.ProgressEntry, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return learningapp.ProgressEntry{}, fmt.Errorf("infra: transition %s/%s: begin tx: %w", topicSlug, conceptSlug, err)
	}
	defer tx.Rollback()

	entry, err := transitionTx(ctx, tx, topicSlug, conceptSlug, decide)
	if err != nil {
		return learningapp.ProgressEntry{}, err
	}
	if err := tx.Commit(); err != nil {
		return learningapp.ProgressEntry{}, fmt.Errorf("infra: transition %s/%s: commit: %w", topicSlug, conceptSlug, err)
	}
	return entry, nil
}

func transitionTx(
	ctx context.Context,
	tx *sql.Tx,
	topicSlug, conceptSlug string,
	decide func(current learningapp.ProgressEntry, hasProgress bool) (learningapp.ProgressDecision, error),
) (learningapp.ProgressEntry, error) {
	lessonID, canonicalTopic, canonicalConcept, err := lockLesson(ctx, tx, topicSlug, conceptSlug)
	if errors.Is(err, sql.ErrNoRows) {
		return learningapp.ProgressEntry{}, fmt.Errorf("infra: transition %s/%s: %w", topicSlug, conceptSlug, learningapp.ErrLessonNotFound)
	}
	if err != nil {
		return learningapp.ProgressEntry{}, fmt.Errorf("infra: transition %s/%s: lock lesson: %w", topicSlug, conceptSlug, err)
	}

	current, hasProgress, err := queryProgressByLessonID(ctx, tx, lessonID, canonicalTopic, canonicalConcept)
	if err != nil {
		return learningapp.ProgressEntry{}, fmt.Errorf("infra: transition %s/%s: %w", topicSlug, conceptSlug, err)
	}

	decision, err := decide(current, hasProgress)
	if err != nil {
		return learningapp.ProgressEntry{}, err
	}

	now := time.Now().UTC()
	if decision.TouchOnly {
		if _, err := tx.ExecContext(ctx, refreshLastReadAtSQL, now, lessonID); err != nil {
			return learningapp.ProgressEntry{}, fmt.Errorf("infra: transition %s/%s: refresh last_read_at: %w", topicSlug, conceptSlug, err)
		}
	} else {
		var firstPassedAt any
		if decision.State.IsPassed() {
			firstPassedAt = now
		}
		if _, err := tx.ExecContext(ctx, upsertProgressByLessonIDSQL, lessonID, decision.State.String(), firstPassedAt, now); err != nil {
			return learningapp.ProgressEntry{}, fmt.Errorf("infra: transition %s/%s: upsert: %w", topicSlug, conceptSlug, err)
		}
	}

	final, hasFinal, err := queryProgressByLessonID(ctx, tx, lessonID, canonicalTopic, canonicalConcept)
	if err != nil {
		return learningapp.ProgressEntry{}, fmt.Errorf("infra: transition %s/%s: reload: %w", topicSlug, conceptSlug, err)
	}
	if !hasFinal {
		return learningapp.ProgressEntry{}, fmt.Errorf("infra: transition %s/%s: progress row missing immediately after write", topicSlug, conceptSlug)
	}
	return final, nil
}

// lockLesson resolves and locks the lesson for (topicSlug, conceptSlug),
// returning sql.ErrNoRows when none matches. concepts is unique per
// (chapter_id, slug), not per topic, so two chapters of the same topic
// sharing a concept slug would make this query match more than one lesson —
// Query (not QueryRow, which would silently pick one) lets that be detected
// and fail loudly instead of ever locking and writing the wrong lesson.
func lockLesson(ctx context.Context, tx *sql.Tx, topicSlug, conceptSlug string) (id int64, canonicalTopic, canonicalConcept string, err error) {
	rows, err := tx.QueryContext(ctx, selectLessonForUpdateSQL, topicSlug, conceptSlug)
	if err != nil {
		return 0, "", "", err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		if found {
			return 0, "", "", fmt.Errorf("concept %s/%s matches more than one lesson — concept slugs are unique per chapter, not per topic", topicSlug, conceptSlug)
		}
		if err := rows.Scan(&id, &canonicalTopic, &canonicalConcept); err != nil {
			return 0, "", "", err
		}
		found = true
	}
	if err := rows.Err(); err != nil {
		return 0, "", "", err
	}
	if !found {
		return 0, "", "", sql.ErrNoRows
	}
	return id, canonicalTopic, canonicalConcept, nil
}

func queryProgressByLessonID(ctx context.Context, tx *sql.Tx, lessonID int64, topicSlug, conceptSlug string) (learningapp.ProgressEntry, bool, error) {
	var (
		stateRaw                  string
		firstPassedAt, lastReadAt sql.NullTime
	)
	err := tx.QueryRowContext(ctx, selectProgressByLessonIDSQL, lessonID).Scan(&stateRaw, &firstPassedAt, &lastReadAt)
	if errors.Is(err, sql.ErrNoRows) {
		return learningapp.ProgressEntry{Topic: topicSlug, Concept: conceptSlug}, false, nil
	}
	if err != nil {
		return learningapp.ProgressEntry{}, false, err
	}

	state, err := domain.NewChunkState(stateRaw)
	if err != nil {
		return learningapp.ProgressEntry{}, false, err
	}
	return learningapp.ProgressEntry{
		Topic: topicSlug, Concept: conceptSlug,
		State: state, FirstPassedAt: nullTimePtr(firstPassedAt), LastReadAt: nullTimePtr(lastReadAt),
	}, true, nil
}

func (r *Repository) AllProgress(ctx context.Context) ([]learningapp.ProgressEntry, error) {
	rows, err := r.db.QueryContext(ctx, selectAllProgressSQL)
	if err != nil {
		return nil, fmt.Errorf("infra: all progress: %w", err)
	}
	defer rows.Close()

	var entries []learningapp.ProgressEntry
	for rows.Next() {
		var (
			topicSlug, conceptSlug, stateRaw string
			firstPassedAt, lastReadAt        sql.NullTime
		)
		if err := rows.Scan(&topicSlug, &conceptSlug, &stateRaw, &firstPassedAt, &lastReadAt); err != nil {
			return nil, fmt.Errorf("infra: scan progress row: %w", err)
		}
		state, err := domain.NewChunkState(stateRaw)
		if err != nil {
			return nil, fmt.Errorf("infra: progress row %s/%s: %w", topicSlug, conceptSlug, err)
		}
		entries = append(entries, learningapp.ProgressEntry{
			Topic:         topicSlug,
			Concept:       conceptSlug,
			State:         state,
			FirstPassedAt: nullTimePtr(firstPassedAt),
			LastReadAt:    nullTimePtr(lastReadAt),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("infra: read progress rows: %w", err)
	}
	return entries, nil
}

func nullTimePtr(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	v := t.Time
	return &v
}

// resolveRecallCheck returns sql.ErrNoRows when no recall_checks row matches,
// mirroring lockLesson. It also refuses to pick a winner when the
// case/accent-insensitive collation matches MORE than one row — two
// recall_checks in the same lesson whose text differs only by case or accent
// would otherwise silently collapse onto whichever row MySQL happened to
// return first, forking or merging attempt history in a way no caller asked
// for. That should never happen given how lessons are authored, but Query
// (not QueryRow) lets it fail loudly instead of guessing.
func resolveRecallCheck(ctx context.Context, db *sql.DB, topicSlug, conceptSlug, question string) (lessonID int64, canonicalQuestion string, err error) {
	rows, err := db.QueryContext(ctx, selectRecallCheckExistsSQL, topicSlug, conceptSlug, question)
	if err != nil {
		return 0, "", err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		if found {
			return 0, "", fmt.Errorf("infra: question %q matches more than one recall_checks row in %s/%s under case/accent-insensitive collation — cannot resolve a single canonical check_key", question, topicSlug, conceptSlug)
		}
		if err := rows.Scan(&lessonID, &canonicalQuestion); err != nil {
			return 0, "", err
		}
		found = true
	}
	if err := rows.Err(); err != nil {
		return 0, "", err
	}
	if !found {
		return 0, "", sql.ErrNoRows
	}
	return lessonID, canonicalQuestion, nil
}

// RecordAttempt is a resolve-then-insert, not Transition's locked
// read-modify-write: recall_attempts is append-only, so there is no stored
// state a concurrent writer could race against — only the "question
// belongs to this lesson" resolution needs the database at all. build is
// called with the CANONICAL question selectRecallCheckExistsSQL matched,
// never question itself — see that constant's own doc comment for why.
func (r *Repository) RecordAttempt(
	ctx context.Context,
	topicSlug, conceptSlug, question string,
	build func(canonicalQuestion string) (domain.RecallAttempt, error),
) (learningapp.AttemptRecord, error) {
	lessonID, canonicalQuestion, err := resolveRecallCheck(ctx, r.db, topicSlug, conceptSlug, question)
	if errors.Is(err, sql.ErrNoRows) {
		exists, existsErr := r.lessonExists(ctx, topicSlug, conceptSlug)
		if existsErr != nil {
			return learningapp.AttemptRecord{}, fmt.Errorf("infra: record attempt %s/%s: check lesson exists: %w", topicSlug, conceptSlug, existsErr)
		}
		if !exists {
			return learningapp.AttemptRecord{}, fmt.Errorf("infra: record attempt %s/%s: %w", topicSlug, conceptSlug, learningapp.ErrLessonNotFound)
		}
		return learningapp.AttemptRecord{}, fmt.Errorf("infra: record attempt %s/%s: %w", topicSlug, conceptSlug, learningapp.ErrCheckNotInLesson)
	}
	if err != nil {
		return learningapp.AttemptRecord{}, fmt.Errorf("infra: record attempt %s/%s: resolve check: %w", topicSlug, conceptSlug, err)
	}

	attempt, err := build(canonicalQuestion)
	if err != nil {
		return learningapp.AttemptRecord{}, err
	}

	var selectedOption any
	if opt := attempt.SelectedOption(); opt != nil {
		selectedOption = *opt
	}

	// Truncated to whole seconds before storing (not just before echoing
	// back) so the Go value and the row created_at's TIMESTAMP(0) column
	// actually holds always agree — MySQL rounds half-up on store, but
	// handler.go's RFC3339 formatting truncates, so an un-truncated now
	// could report a created_at up to a second earlier than the real row.
	now := time.Now().UTC().Truncate(time.Second)
	_, err = r.db.ExecContext(ctx, insertRecallAttemptSQL,
		lessonID, attempt.CheckKey().String(), attempt.Kind().String(), attempt.Confidence().String(),
		attempt.Outcome().String(), selectedOption, attempt.GradedBy().String(), now,
	)
	if err != nil {
		return learningapp.AttemptRecord{}, fmt.Errorf("infra: record attempt %s/%s: insert: %w", topicSlug, conceptSlug, err)
	}

	return learningapp.AttemptRecord{
		CheckKey:       attempt.CheckKey().String(),
		Question:       canonicalQuestion,
		Kind:           attempt.Kind(),
		Confidence:     attempt.Confidence(),
		Outcome:        attempt.Outcome(),
		SelectedOption: attempt.SelectedOption(),
		GradedBy:       attempt.GradedBy(),
		CreatedAt:      now,
	}, nil
}

func (r *Repository) lessonExists(ctx context.Context, topicSlug, conceptSlug string) (bool, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, selectLessonExistsSQL, topicSlug, conceptSlug).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
