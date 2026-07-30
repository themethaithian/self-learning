package infra

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

const (
	selectConceptIDSQL = `
SELECT co.id
FROM concepts co
JOIN chapters ch ON co.chapter_id = ch.id
JOIN topics t ON ch.topic_id = t.id
WHERE t.slug = ? AND co.slug = ?`

	upsertLessonSQL = `
INSERT INTO lessons (concept_id, version, title_en, est_minutes, body_md, refs)
VALUES (?, ?, ?, ?, ?, ?) AS new
ON DUPLICATE KEY UPDATE
    id = LAST_INSERT_ID(id),
    version = new.version,
    title_en = new.title_en,
    est_minutes = new.est_minutes,
    body_md = new.body_md,
    refs = new.refs`

	deleteRecallChecksSQL = `DELETE FROM recall_checks WHERE lesson_id = ?`

	insertRecallCheckSQL = `
INSERT INTO recall_checks (lesson_id, position, type, question, expected_answer, options, explanation)
VALUES (?, ?, ?, ?, ?, ?, ?)`
)

type referenceRowDTO struct {
	Title  string `json:"title"`
	Source string `json:"source"`
	Why    string `json:"why"`
}

// SaveLesson resolves l's concept by (topicSlug, l.Slug()), upserts the
// lesson row keyed by that concept, and fully replaces its recall checks —
// all in one transaction. recall_checks are owned by the Lesson aggregate,
// so delete-then-insert is correct: the rows carry no identity a caller
// depends on across imports. inserted reports whether this created a new
// lesson row (true) or updated an existing one (false), for the importer's
// per-file summary.
func (r *Repository) SaveLesson(ctx context.Context, topicSlug string, l domain.Lesson) (inserted bool, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("infra: save lesson %q: begin tx: %w", l.Slug().String(), err)
	}
	defer tx.Rollback()

	inserted, err = saveLessonTx(ctx, tx, topicSlug, l)
	if err != nil {
		return false, fmt.Errorf("infra: save lesson %q: %w", l.Slug().String(), err)
	}
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("infra: save lesson %q: commit: %w", l.Slug().String(), err)
	}
	return inserted, nil
}

func saveLessonTx(ctx context.Context, tx *sql.Tx, topicSlug string, l domain.Lesson) (bool, error) {
	conceptID, err := resolveConceptID(ctx, tx, topicSlug, l.Slug().String())
	if err != nil {
		return false, err
	}

	refsJSON, err := marshalReferences(l.References())
	if err != nil {
		return false, fmt.Errorf("marshal references: %w", err)
	}

	res, err := tx.ExecContext(ctx, upsertLessonSQL,
		conceptID, l.Version(), l.TitleEn(), l.EstMinutes().Int(), l.BodyMd(), refsJSON)
	if err != nil {
		return false, fmt.Errorf("upsert lesson: %w", err)
	}
	lessonID, err := lastInsertID(res, "lesson")
	if err != nil {
		return false, err
	}
	rowsAffected, _ := res.RowsAffected()

	if _, err := tx.ExecContext(ctx, deleteRecallChecksSQL, lessonID); err != nil {
		return false, fmt.Errorf("delete recall checks: %w", err)
	}
	for _, rc := range l.RecallChecks() {
		if err := insertRecallCheck(ctx, tx, lessonID, rc); err != nil {
			return false, fmt.Errorf("recall check %d: %w", rc.Position().Int(), err)
		}
	}

	return rowsAffected == 1, nil
}

// resolveConceptID is why lessons import after curriculum: a lesson file
// names its concept by (topic slug, concept slug), and that pair must
// already exist as a concepts row.
func resolveConceptID(ctx context.Context, tx *sql.Tx, topicSlug, conceptSlug string) (int64, error) {
	var id int64
	err := tx.QueryRowContext(ctx, selectConceptIDSQL, topicSlug, conceptSlug).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("concept %s/%s not found — import curriculum first", topicSlug, conceptSlug)
	}
	if err != nil {
		return 0, fmt.Errorf("resolve concept id: %w", err)
	}
	return id, nil
}

func marshalReferences(refs []domain.Reference) (string, error) {
	rows := make([]referenceRowDTO, len(refs))
	for i, r := range refs {
		rows[i] = referenceRowDTO{Title: r.Title(), Source: r.Source(), Why: r.Why()}
	}
	b, err := json.Marshal(rows)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func insertRecallCheck(ctx context.Context, tx *sql.Tx, lessonID int64, rc domain.RecallCheck) error {
	var options any
	if rc.Kind().IsMCQ() {
		b, err := json.Marshal(rc.Options())
		if err != nil {
			return fmt.Errorf("marshal options: %w", err)
		}
		options = string(b)
	}

	// NULL, not "", for an absent explanation — the same "optional value ->
	// SQL NULL" convention the options column above already uses.
	var explanation any
	if rc.Explanation() != "" {
		explanation = rc.Explanation()
	}

	_, err := tx.ExecContext(ctx, insertRecallCheckSQL,
		lessonID, rc.Position().Int(), rc.Kind().String(), rc.Question(), rc.ExpectedAnswer(), options, explanation)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	return nil
}
