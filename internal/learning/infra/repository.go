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
	// selectConceptProgressSQL resolves lesson_progress by joining lessons ->
	// concepts -> chapters -> topics on (t.slug, co.slug) — both slugs are
	// required because concept slugs are only topic-wide unique. The LEFT
	// JOIN means a lesson with no progress row yet still produces one row
	// (lp.* NULL) rather than vanishing, so a caller can tell "no lesson"
	// (zero rows) from "lesson exists, never started" (one row, NULL state).
	selectConceptProgressSQL = `
SELECT l.id, lp.state, lp.first_passed_at, lp.last_read_at
FROM lessons l
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
LEFT JOIN lesson_progress lp ON lp.lesson_id = l.id
WHERE t.slug = ? AND co.slug = ?`

	// upsertProgressSQL writes state on first write or on a real transition.
	// first_passed_at is preserved via COALESCE against the table's own
	// (pre-update) value — not overwritten by the incoming value — so a
	// caller can safely pass "now" every time state is passed and still
	// never move an already-set first_passed_at.
	upsertProgressSQL = `
INSERT INTO lesson_progress (lesson_id, state, first_passed_at, last_read_at)
VALUES (?, ?, ?, ?)
AS new
ON DUPLICATE KEY UPDATE
    state = new.state,
    first_passed_at = COALESCE(lesson_progress.first_passed_at, new.first_passed_at),
    last_read_at = new.last_read_at`

	// touchProgressSQL is the idempotent / forward-only no-op path: it
	// refreshes last_read_at only, never state or first_passed_at.
	touchProgressSQL = `
UPDATE lesson_progress lp
JOIN lessons l ON l.id = lp.lesson_id
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
SET lp.last_read_at = ?
WHERE t.slug = ? AND co.slug = ?`

	// selectAllProgressSQL starts from lesson_progress itself (inner joins
	// outward), so only concepts with a progress row are ever returned —
	// absence means "not started", per the ticket's response contract.
	selectAllProgressSQL = `
SELECT t.slug, co.slug, lp.state, lp.first_passed_at, lp.last_read_at
FROM lesson_progress lp
JOIN lessons l ON l.id = lp.lesson_id
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
ORDER BY t.slug, co.slug`
)

// Repository is the MySQL adapter for app.Repository.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type progressRow struct {
	lessonID      int64
	state         sql.NullString
	firstPassedAt sql.NullTime
	lastReadAt    sql.NullTime
}

func (r *Repository) queryConceptRow(ctx context.Context, topicSlug, conceptSlug string) (progressRow, bool, error) {
	var row progressRow
	err := r.db.QueryRowContext(ctx, selectConceptProgressSQL, topicSlug, conceptSlug).
		Scan(&row.lessonID, &row.state, &row.firstPassedAt, &row.lastReadAt)
	if errors.Is(err, sql.ErrNoRows) {
		return progressRow{}, false, nil
	}
	if err != nil {
		return progressRow{}, false, err
	}
	return row, true, nil
}

func (r *Repository) ConceptProgress(ctx context.Context, topicSlug, conceptSlug string) (learningapp.ProgressEntry, bool, bool, error) {
	row, found, err := r.queryConceptRow(ctx, topicSlug, conceptSlug)
	if err != nil {
		return learningapp.ProgressEntry{}, false, false, fmt.Errorf("infra: concept progress %s/%s: %w", topicSlug, conceptSlug, err)
	}
	if !found {
		return learningapp.ProgressEntry{}, false, false, nil
	}
	if !row.state.Valid {
		return learningapp.ProgressEntry{Topic: topicSlug, Concept: conceptSlug}, true, false, nil
	}

	state, err := domain.NewChunkState(row.state.String)
	if err != nil {
		return learningapp.ProgressEntry{}, false, false, fmt.Errorf("infra: concept progress %s/%s: %w", topicSlug, conceptSlug, err)
	}
	return learningapp.ProgressEntry{
		Topic:         topicSlug,
		Concept:       conceptSlug,
		State:         state,
		FirstPassedAt: nullTimePtr(row.firstPassedAt),
		LastReadAt:    nullTimePtr(row.lastReadAt),
	}, true, true, nil
}

func (r *Repository) UpsertProgress(ctx context.Context, topicSlug, conceptSlug string, state domain.ChunkState) error {
	row, found, err := r.queryConceptRow(ctx, topicSlug, conceptSlug)
	if err != nil {
		return fmt.Errorf("infra: upsert progress %s/%s: %w", topicSlug, conceptSlug, err)
	}
	if !found {
		return fmt.Errorf("infra: upsert progress %s/%s: %w", topicSlug, conceptSlug, learningapp.ErrLessonNotFound)
	}

	now := time.Now().UTC()
	var firstPassedAt any
	if state.IsPassed() {
		firstPassedAt = now
	}
	if _, err := r.db.ExecContext(ctx, upsertProgressSQL, row.lessonID, state.String(), firstPassedAt, now); err != nil {
		return fmt.Errorf("infra: upsert progress %s/%s: %w", topicSlug, conceptSlug, err)
	}
	return nil
}

func (r *Repository) TouchProgress(ctx context.Context, topicSlug, conceptSlug string) error {
	now := time.Now().UTC()
	if _, err := r.db.ExecContext(ctx, touchProgressSQL, now, topicSlug, conceptSlug); err != nil {
		return fmt.Errorf("infra: touch progress %s/%s: %w", topicSlug, conceptSlug, err)
	}
	return nil
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
