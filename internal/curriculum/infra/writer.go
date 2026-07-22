package infra

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

// id = LAST_INSERT_ID(id) makes MySQL report the existing row's id on the
// update path too — plain ON DUPLICATE KEY UPDATE would otherwise leave
// LAST_INSERT_ID() unset. new.col reads the row alias; VALUES() is deprecated.
const (
	upsertTopicSQL = `
INSERT INTO topics (track, slug, title, position)
VALUES (?, ?, ?, ?) AS new
ON DUPLICATE KEY UPDATE
    id = LAST_INSERT_ID(id),
    track = new.track,
    title = new.title,
    position = new.position`

	upsertChapterSQL = `
INSERT INTO chapters (topic_id, slug, title, position)
VALUES (?, ?, ?, ?) AS new
ON DUPLICATE KEY UPDATE
    id = LAST_INSERT_ID(id),
    title = new.title,
    position = new.position`

	upsertConceptSQL = `
INSERT INTO concepts (chapter_id, slug, title, outline, position)
VALUES (?, ?, ?, ?, ?) AS new
ON DUPLICATE KEY UPDATE
    id = LAST_INSERT_ID(id),
    title = new.title,
    outline = new.outline,
    position = new.position`

	deleteStaleConceptsSQL = `DELETE FROM concepts WHERE chapter_id = ? AND slug NOT IN (%s)`

	// Concepts under a chapter that is itself about to be deleted (a renamed
	// or removed chapter) are never visited by deleteStaleConceptsSQL, since
	// that only runs for chapters still present in the topic. This clears
	// them first so deleteStaleChaptersSQL doesn't hit fk_concepts_chapter.
	deleteStaleChapterConceptsSQL = `DELETE FROM concepts WHERE chapter_id IN (SELECT id FROM chapters WHERE topic_id = ? AND slug NOT IN (%s))`

	deleteStaleChaptersSQL = `DELETE FROM chapters WHERE topic_id = ? AND slug NOT IN (%s)`
)

// SaveTopic upserts t and its whole chapter/concept tree in one transaction,
// then deletes any chapter or concept row under t that the current t no
// longer lists — a moved, renamed, or deleted node in the source file must
// not leave an orphan row behind, since Topics rehydrates through NewTopic
// and rejects a topic with a duplicate or gapped slug/position. The
// aggregate is the transaction boundary, so a topic is never visible with
// only some of its children applied.
func (r *Repository) SaveTopic(ctx context.Context, t domain.Topic) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("infra: save topic %q: begin tx: %w", t.Slug().String(), err)
	}
	defer tx.Rollback()

	if err := saveTopicTx(ctx, tx, t); err != nil {
		return fmt.Errorf("infra: save topic %q: %w", t.Slug().String(), err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("infra: save topic %q: commit: %w", t.Slug().String(), err)
	}
	return nil
}

func saveTopicTx(ctx context.Context, tx *sql.Tx, t domain.Topic) error {
	topicID, err := upsertTopic(ctx, tx, t)
	if err != nil {
		return fmt.Errorf("topic: %w", err)
	}

	chapters := t.Chapters()
	chapterSlugs := make([]string, 0, len(chapters))
	for _, ch := range chapters {
		chapterID, err := upsertChapter(ctx, tx, topicID, ch)
		if err != nil {
			return fmt.Errorf("chapter %q: %w", ch.Slug().String(), err)
		}
		chapterSlugs = append(chapterSlugs, ch.Slug().String())

		concepts := ch.Concepts()
		conceptSlugs := make([]string, 0, len(concepts))
		for _, c := range concepts {
			if err := upsertConcept(ctx, tx, chapterID, c); err != nil {
				return fmt.Errorf("chapter %q: concept %q: %w", ch.Slug().String(), c.Slug().String(), err)
			}
			conceptSlugs = append(conceptSlugs, c.Slug().String())
		}

		if err := deleteNotIn(ctx, tx, deleteStaleConceptsSQL, chapterID, conceptSlugs); err != nil {
			return fmt.Errorf("chapter %q: delete stale concepts: %w", ch.Slug().String(), err)
		}
	}

	if err := deleteNotIn(ctx, tx, deleteStaleChapterConceptsSQL, topicID, chapterSlugs); err != nil {
		return fmt.Errorf("delete concepts of stale chapters: %w", err)
	}
	if err := deleteNotIn(ctx, tx, deleteStaleChaptersSQL, topicID, chapterSlugs); err != nil {
		return fmt.Errorf("delete stale chapters: %w", err)
	}
	return nil
}

// deleteNotIn runs one of the delete*SQL templates against parentID, keeping
// only rows whose slug is in keepSlugs. keepSlugs is never empty: NewChapter
// and NewTopic both reject an empty child list, so every parent SaveTopic
// writes has at least one surviving child by the time this runs.
func deleteNotIn(ctx context.Context, tx *sql.Tx, template string, parentID int64, keepSlugs []string) error {
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(keepSlugs)), ",")
	query := fmt.Sprintf(template, placeholders)

	args := make([]any, 0, len(keepSlugs)+1)
	args = append(args, parentID)
	for _, slug := range keepSlugs {
		args = append(args, slug)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func upsertTopic(ctx context.Context, tx *sql.Tx, t domain.Topic) (int64, error) {
	res, err := tx.ExecContext(ctx, upsertTopicSQL, t.Track().String(), t.Slug().String(), t.Title(), t.Position().Int())
	if err != nil {
		return 0, fmt.Errorf("upsert: %w", err)
	}
	return lastInsertID(res, "topic")
}

func upsertChapter(ctx context.Context, tx *sql.Tx, topicID int64, ch domain.Chapter) (int64, error) {
	res, err := tx.ExecContext(ctx, upsertChapterSQL, topicID, ch.Slug().String(), ch.Title(), ch.Position().Int())
	if err != nil {
		return 0, fmt.Errorf("upsert: %w", err)
	}
	return lastInsertID(res, "chapter")
}

func upsertConcept(ctx context.Context, tx *sql.Tx, chapterID int64, c domain.Concept) error {
	res, err := tx.ExecContext(ctx, upsertConceptSQL, chapterID, c.Slug().String(), c.Title(), c.Outline(), c.Position().Int())
	if err != nil {
		return fmt.Errorf("upsert: %w", err)
	}
	_, err = lastInsertID(res, "concept")
	return err
}

// lastInsertID rejects 0 explicitly: a genuine row id is never 0, so a 0
// here means the ON DUPLICATE KEY UPDATE clause failed to set
// id = LAST_INSERT_ID(id) — a regression that would otherwise only surface
// two layers away, as a foreign key violation on the next insert.
func lastInsertID(res sql.Result, what string) (int64, error) {
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: last insert id: %w", what, err)
	}
	if id == 0 {
		return 0, fmt.Errorf("%s: last insert id was 0, want id = LAST_INSERT_ID(id) to have set it", what)
	}
	return id, nil
}
