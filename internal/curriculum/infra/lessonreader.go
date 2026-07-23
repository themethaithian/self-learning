package infra

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	curriculumapp "github.com/themethaithian/self-learning/internal/curriculum/app"
	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

const (
	selectLessonByConceptSQL = `
SELECT l.id, l.version, l.title_en, l.est_minutes, l.body_md, l.refs
FROM lessons l
JOIN concepts co ON co.id = l.concept_id
JOIN chapters ch ON ch.id = co.chapter_id
JOIN topics t ON t.id = ch.topic_id
WHERE t.slug = ? AND co.slug = ?`

	selectRecallChecksByLessonSQL = `
SELECT position, type, question, expected_answer, options
FROM recall_checks
WHERE lesson_id = ?
ORDER BY position`
)

type lessonRow struct {
	id         int64
	version    int
	titleEn    string
	estMinutes int
	bodyMd     string
	refs       string
}

// LessonByConcept resolves the lesson by joining lessons -> concepts ->
// chapters -> topics on (t.slug, co.slug) — both slugs are required because
// concept slugs are only topic-wide unique, the same reason SaveLesson's
// concept lookup filters on both. Returns curriculumapp.ErrLessonNotFound
// when the concept has no lesson yet, so the handler can 404 vs 500.
func (r *Repository) LessonByConcept(ctx context.Context, topicSlug, conceptSlug string) (domain.Lesson, error) {
	var row lessonRow
	err := r.db.QueryRowContext(ctx, selectLessonByConceptSQL, topicSlug, conceptSlug).Scan(
		&row.id, &row.version, &row.titleEn, &row.estMinutes, &row.bodyMd, &row.refs,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Lesson{}, fmt.Errorf("infra: lesson %s/%s: %w", topicSlug, conceptSlug, curriculumapp.ErrLessonNotFound)
	}
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("infra: query lesson %s/%s: %w", topicSlug, conceptSlug, err)
	}

	references, err := unmarshalReferences(row.refs)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("infra: lesson %s/%s: %w", topicSlug, conceptSlug, err)
	}
	recallChecks, err := r.recallChecksByLesson(ctx, row.id)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("infra: lesson %s/%s: %w", topicSlug, conceptSlug, err)
	}

	slug, err := domain.NewSlug(conceptSlug)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("infra: lesson %s/%s: %w", topicSlug, conceptSlug, err)
	}
	estMinutes, err := domain.NewEstMinutes(row.estMinutes)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("infra: lesson %s/%s: %w", topicSlug, conceptSlug, err)
	}

	lesson, err := domain.NewLesson(slug, row.version, row.titleEn, estMinutes, row.bodyMd, references, recallChecks)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("infra: lesson %s/%s: %w", topicSlug, conceptSlug, err)
	}
	return lesson, nil
}

func unmarshalReferences(refsJSON string) ([]domain.Reference, error) {
	var rows []referenceRowDTO
	if err := json.Unmarshal([]byte(refsJSON), &rows); err != nil {
		return nil, fmt.Errorf("unmarshal references: %w", err)
	}
	references := make([]domain.Reference, 0, len(rows))
	for i, row := range rows {
		ref, err := domain.NewReference(row.Title, row.Source, row.Why)
		if err != nil {
			return nil, fmt.Errorf("reference %d: %w", i, err)
		}
		references = append(references, ref)
	}
	return references, nil
}

func (r *Repository) recallChecksByLesson(ctx context.Context, lessonID int64) ([]domain.RecallCheck, error) {
	rows, err := r.db.QueryContext(ctx, selectRecallChecksByLessonSQL, lessonID)
	if err != nil {
		return nil, fmt.Errorf("query recall checks: %w", err)
	}
	defer rows.Close()

	var checks []domain.RecallCheck
	for rows.Next() {
		var (
			position       int
			kind           string
			question       string
			expectedAnswer string
			options        sql.NullString
		)
		if err := rows.Scan(&position, &kind, &question, &expectedAnswer, &options); err != nil {
			return nil, fmt.Errorf("scan recall check: %w", err)
		}

		check, err := toRecallCheck(position, kind, question, expectedAnswer, options)
		if err != nil {
			return nil, fmt.Errorf("recall check %d: %w", position, err)
		}
		checks = append(checks, check)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read recall checks: %w", err)
	}
	return checks, nil
}

func toRecallCheck(position int, kind, question, expectedAnswer string, options sql.NullString) (domain.RecallCheck, error) {
	var optionValues []string
	if options.Valid {
		if err := json.Unmarshal([]byte(options.String), &optionValues); err != nil {
			return domain.RecallCheck{}, fmt.Errorf("unmarshal options: %w", err)
		}
	}

	pos, err := domain.NewPosition(position)
	if err != nil {
		return domain.RecallCheck{}, err
	}
	recallKind, err := domain.NewRecallKind(kind)
	if err != nil {
		return domain.RecallCheck{}, err
	}
	return domain.NewRecallCheck(pos, recallKind, question, expectedAnswer, optionValues)
}
