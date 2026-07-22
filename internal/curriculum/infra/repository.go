// Package infra adapts the curriculum domain and app layers to MySQL and
// HTTP — the vendor-facing details the domain and app layers must never see.
package infra

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

// topicsQuery is one statement across topics, chapters, and concepts to
// avoid N+1 round-trips. LEFT JOINs so a topic with no chapters (or a
// chapter with no concepts) still produces a row instead of silently
// vanishing from the result — assembleTree turns that into a named error.
// ORDER BY makes the result deterministic to read while debugging;
// assembleTree groups rows by slug and does not depend on this order.
const topicsQuery = `
SELECT
	t.track,
	t.slug,
	t.title,
	t.position,
	c.slug,
	c.title,
	c.position,
	co.slug,
	co.title,
	co.outline,
	co.position
FROM topics t
LEFT JOIN chapters c ON c.topic_id = t.id
LEFT JOIN concepts co ON co.chapter_id = c.id
ORDER BY t.position, t.id, c.position, c.id, co.position, co.id`

// Repository is the MySQL adapter for app.Repository.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Topics keeps the fold logic in assembleTree, which is testable without a
// database.
func (r *Repository) Topics(ctx context.Context) ([]domain.Topic, error) {
	rows, err := r.db.QueryContext(ctx, topicsQuery)
	if err != nil {
		return nil, fmt.Errorf("infra: query topics: %w", err)
	}
	defer rows.Close()

	var flat []conceptRow
	for rows.Next() {
		var row conceptRow
		if err := rows.Scan(
			&row.topicTrack, &row.topicSlug, &row.topicTitle, &row.topicPosition,
			&row.chapterSlug, &row.chapterTitle, &row.chapterPosition,
			&row.conceptSlug, &row.conceptTitle, &row.conceptOutline, &row.conceptPosition,
		); err != nil {
			return nil, fmt.Errorf("infra: scan topics row: %w", err)
		}
		flat = append(flat, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("infra: read topics rows: %w", err)
	}

	topics, err := assembleTree(flat)
	if err != nil {
		return nil, fmt.Errorf("infra: assemble topics: %w", err)
	}
	return topics, nil
}

// conceptRow is one row of topicsQuery. Chapter and concept columns are
// nullable: the LEFT JOIN leaves them NULL when a topic has no chapters, or
// a chapter has no concepts.
type conceptRow struct {
	topicTrack      string
	topicSlug       string
	topicTitle      string
	topicPosition   int
	chapterSlug     sql.NullString
	chapterTitle    sql.NullString
	chapterPosition sql.NullInt64
	conceptSlug     sql.NullString
	conceptTitle    sql.NullString
	conceptOutline  sql.NullString
	conceptPosition sql.NullInt64
}

type chapterAcc struct {
	slug     string
	title    string
	position int
	concepts []domain.Concept
}

type topicAcc struct {
	track        string
	slug         string
	title        string
	position     int
	chapterOrder []string
	chapters     map[string]*chapterAcc
}

// assembleTree folds flat, possibly scrambled-order rows into the
// curriculum aggregate by calling the domain constructors bottom-up, so a
// row that violates a domain invariant fails loudly instead of building an
// invalid tree. An empty input is not an error: it returns an empty slice.
func assembleTree(rows []conceptRow) ([]domain.Topic, error) {
	order := make([]string, 0)
	bySlug := make(map[string]*topicAcc)

	for _, row := range rows {
		ta, ok := bySlug[row.topicSlug]
		if !ok {
			ta = &topicAcc{
				track:    row.topicTrack,
				slug:     row.topicSlug,
				title:    row.topicTitle,
				position: row.topicPosition,
				chapters: make(map[string]*chapterAcc),
			}
			bySlug[row.topicSlug] = ta
			order = append(order, row.topicSlug)
		}
		if err := ta.addRow(row); err != nil {
			return nil, fmt.Errorf("topic %q: %w", ta.slug, err)
		}
	}

	topics := make([]domain.Topic, 0, len(order))
	for _, slug := range order {
		topic, err := bySlug[slug].build()
		if err != nil {
			return nil, fmt.Errorf("topic %q: %w", slug, err)
		}
		topics = append(topics, topic)
	}

	// topics.position is not unique across topics (each track's curriculum
	// JSON starts numbering at 1), so slug breaks ties for a total order.
	slices.SortStableFunc(topics, func(a, b domain.Topic) int {
		return cmp.Or(
			cmp.Compare(a.Position().Int(), b.Position().Int()),
			cmp.Compare(a.Slug().String(), b.Slug().String()),
		)
	})
	return topics, nil
}

func (ta *topicAcc) addRow(row conceptRow) error {
	if !row.chapterSlug.Valid {
		return nil
	}

	ca, ok := ta.chapters[row.chapterSlug.String]
	if !ok {
		ca = &chapterAcc{
			slug:     row.chapterSlug.String,
			title:    row.chapterTitle.String,
			position: int(row.chapterPosition.Int64),
		}
		ta.chapters[row.chapterSlug.String] = ca
		ta.chapterOrder = append(ta.chapterOrder, row.chapterSlug.String)
	}

	if !row.conceptSlug.Valid {
		return nil
	}

	concept, err := newConcept(row)
	if err != nil {
		return fmt.Errorf("chapter %q: %w", ca.slug, err)
	}
	ca.concepts = append(ca.concepts, concept)
	return nil
}

func (ta *topicAcc) build() (domain.Topic, error) {
	chapters := make([]domain.Chapter, 0, len(ta.chapterOrder))
	for _, slug := range ta.chapterOrder {
		chapter, err := ta.chapters[slug].build()
		if err != nil {
			return domain.Topic{}, err
		}
		chapters = append(chapters, chapter)
	}

	track, err := domain.NewTrack(ta.track)
	if err != nil {
		return domain.Topic{}, err
	}
	slug, err := domain.NewSlug(ta.slug)
	if err != nil {
		return domain.Topic{}, err
	}
	position, err := domain.NewPosition(ta.position)
	if err != nil {
		return domain.Topic{}, err
	}
	return domain.NewTopic(track, slug, ta.title, position, chapters)
}

func (ca *chapterAcc) build() (domain.Chapter, error) {
	slug, err := domain.NewSlug(ca.slug)
	if err != nil {
		return domain.Chapter{}, fmt.Errorf("chapter: %w", err)
	}
	position, err := domain.NewPosition(ca.position)
	if err != nil {
		return domain.Chapter{}, fmt.Errorf("chapter %q: %w", ca.slug, err)
	}
	return domain.NewChapter(slug, ca.title, position, ca.concepts)
}

func newConcept(row conceptRow) (domain.Concept, error) {
	slug, err := domain.NewSlug(row.conceptSlug.String)
	if err != nil {
		return domain.Concept{}, fmt.Errorf("concept: %w", err)
	}
	position, err := domain.NewPosition(int(row.conceptPosition.Int64))
	if err != nil {
		return domain.Concept{}, fmt.Errorf("concept %q: %w", row.conceptSlug.String, err)
	}
	return domain.NewConcept(slug, row.conceptTitle.String, row.conceptOutline.String, position)
}
