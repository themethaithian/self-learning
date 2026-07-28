package infra

import (
	"context"
	"errors"
	"strings"
	"testing"

	curriculumapp "github.com/themethaithian/self-learning/internal/curriculum/app"
)

func TestRepositoryTopics_Success(t *testing.T) {
	result := &stubResult{rows: []stubRow{
		stubFullRow("go", "go-basics", "Go Basics Title", 3, "syntax-chapter", "Syntax Chapter Title", 5, "variables-concept", "Variables Concept Title", "Variables Outline Text", 7),
	}}
	db := openStubDB(t, result)
	repo := NewRepository(db)

	got, _, err := repo.Topics(context.Background())
	if err != nil {
		t.Fatalf("Topics() unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("Topics() returned %d topics, want 1", len(got))
	}

	topic := got[0]
	if topic.Track().String() != "go" || topic.Slug().String() != "go-basics" ||
		topic.Title() != "Go Basics Title" || topic.Position().Int() != 3 {
		t.Fatalf("topic = {%q %q %q %d}, want {go go-basics \"Go Basics Title\" 3}",
			topic.Track().String(), topic.Slug().String(), topic.Title(), topic.Position().Int())
	}

	chapters := topic.Chapters()
	if len(chapters) != 1 {
		t.Fatalf("Chapters() = %d, want 1", len(chapters))
	}
	chapter := chapters[0]
	if chapter.Slug().String() != "syntax-chapter" || chapter.Title() != "Syntax Chapter Title" || chapter.Position().Int() != 5 {
		t.Fatalf("chapter = {%q %q %d}, want {syntax-chapter \"Syntax Chapter Title\" 5}",
			chapter.Slug().String(), chapter.Title(), chapter.Position().Int())
	}

	concepts := chapter.Concepts()
	if len(concepts) != 1 {
		t.Fatalf("Concepts() = %d, want 1", len(concepts))
	}
	concept := concepts[0]
	if concept.Slug().String() != "variables-concept" ||
		concept.Title() != "Variables Concept Title" ||
		concept.Outline() != "Variables Outline Text" ||
		concept.Position().Int() != 7 {
		t.Fatalf("concept = {%q %q %q %d}, want {variables-concept \"Variables Concept Title\" \"Variables Outline Text\" 7}",
			concept.Slug().String(), concept.Title(), concept.Outline(), concept.Position().Int())
	}

	if rows := result.driverRows(); rows == nil || !rows.closed {
		t.Fatal("Topics() did not close the driver rows")
	}
}

func TestRepositoryTopics_EmptyIsNonNilEmptySlice(t *testing.T) {
	db := openStubDB(t, &stubResult{})
	repo := NewRepository(db)

	got, _, err := repo.Topics(context.Background())
	if err != nil {
		t.Fatalf("Topics() unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("Topics() returned nil, want a non-nil empty slice")
	}
	if len(got) != 0 {
		t.Fatalf("Topics() = %d topics, want 0", len(got))
	}
}

// TestRepositoryTopics_ChildlessTopicFailsLoudly exercises the real
// database/sql NULL path (driver.Value nil -> sql.NullString{Valid: false}),
// which the assembleTree unit tests bypass by constructing NullString
// directly: a topic with no chapters must still reach assembleTree as a row
// and fail via ErrNoChildren, not vanish.
func TestRepositoryTopics_ChildlessTopicFailsLoudly(t *testing.T) {
	db := openStubDB(t, &stubResult{rows: []stubRow{
		stubTopicOnlyRow("go", "go-basics", "Go Basics", 1),
	}})
	repo := NewRepository(db)

	_, _, err := repo.Topics(context.Background())
	if err == nil {
		t.Fatal("Topics() expected an error for a childless topic, got nil")
	}
	if !strings.Contains(err.Error(), "go-basics") || !strings.Contains(err.Error(), "no children") {
		t.Fatalf("Topics() error = %q, want it to name the topic and say no children", err.Error())
	}
}

// TestRepositoryTopics_QueryUsesLeftJoins pins the one thing a stub that
// doesn't execute SQL can still check about join semantics: the query text
// says LEFT JOIN. A stub can't verify join *behavior* (it always returns
// whatever fixture rows it's given, regardless of query), so this is
// deliberately a text assertion, not a behavioral one. LEFT JOIN matters
// because a topic with no chapters (or a chapter with no concepts) must
// still produce a row that reaches assembleTree, so it fails loudly via
// ErrNoChildren instead of silently vanishing from an inner-joined result.
func TestRepositoryTopics_QueryUsesLeftJoins(t *testing.T) {
	result := &stubResult{rows: []stubRow{
		stubFullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1),
	}}
	db := openStubDB(t, result)
	repo := NewRepository(db)

	if _, _, err := repo.Topics(context.Background()); err != nil {
		t.Fatalf("Topics() unexpected error: %v", err)
	}

	query := result.query()
	if !strings.Contains(query, "LEFT JOIN chapters") {
		t.Errorf("query = %q, want it to LEFT JOIN chapters", query)
	}
	if !strings.Contains(query, "LEFT JOIN concepts") {
		t.Errorf("query = %q, want it to LEFT JOIN concepts", query)
	}
	if !strings.Contains(query, "LEFT JOIN lessons l ON l.concept_id = co.id") {
		t.Errorf("query = %q, want lessons joined on the UNIQUE concept_id (any other predicate risks fanning rows out)", query)
	}
}

// stubRowInvalidTopicPosition makes Scan fail on the *first* row, before the
// driver-level rows are ever exhausted to io.EOF — sql.Rows only auto-closes
// on a clean drain, so this is the shape that actually exercises a missing
// defer rows.Close().
func stubRowInvalidTopicPosition() stubRow {
	row := stubFullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1)
	row["t.position"] = "not-a-number"
	return row
}

func TestRepositoryTopics_ScanErrorClosesRows(t *testing.T) {
	result := &stubResult{rows: []stubRow{stubRowInvalidTopicPosition()}}
	db := openStubDB(t, result)
	repo := NewRepository(db)

	if _, _, err := repo.Topics(context.Background()); err == nil {
		t.Fatal("Topics() expected a scan error, got nil")
	}
	if rows := result.driverRows(); rows == nil || !rows.closed {
		t.Fatal("Topics() left the driver rows open after a mid-loop scan error (missing defer rows.Close())")
	}
}

func TestRepositoryTopics_BrokenRowStreamIsAnError(t *testing.T) {
	db := openStubDB(t, &stubResult{
		rows: []stubRow{
			stubFullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1),
		},
		nextErr: errors.New("stub: connection dropped mid-stream"),
	})
	repo := NewRepository(db)

	if _, _, err := repo.Topics(context.Background()); err == nil {
		t.Fatal("Topics() expected an error from a broken row stream (rows.Err()), got nil")
	}
}

func TestRepositoryTopics_QueryErrorWraps(t *testing.T) {
	db := openStubDB(t, &stubResult{queryErr: errors.New("stub: connection refused")})
	repo := NewRepository(db)

	if _, _, err := repo.Topics(context.Background()); err == nil {
		t.Fatal("Topics() expected an error, got nil")
	}
}

func TestRepositoryTopics_ContextPropagatesToDriver(t *testing.T) {
	db := openStubDB(t, &stubResult{rows: []stubRow{
		stubFullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1),
	}})
	repo := NewRepository(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, _, err := repo.Topics(ctx); err == nil {
		t.Fatal("Topics() expected an error for an already-canceled context, got nil (ctx is not reaching the driver)")
	}
}

// TestRepositoryTopics_AvailabilityMapMatchesTreeConcepts checks that
// Repository.Topics's two independent passes over the same flat rows —
// assembleTree and conceptAvailabilityFromRows — agree on which concepts
// exist and which of those have a lesson. This is a fixture-based stub
// test: it cannot prove a real LEFT JOIN in production MySQL is correct or
// that ORDER BY still holds after adding it (the stub returns whatever rows
// a test hands it, regardless of query text) — TestRepositoryTopics_
// QueryUsesLeftJoins pins the query text itself for that.
func TestRepositoryTopics_AvailabilityMapMatchesTreeConcepts(t *testing.T) {
	withLesson := stubFullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1)
	withLesson["l.est_minutes"] = 8
	withoutLesson := stubFullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "loops", "Loops", "outline", 2)

	db := openStubDB(t, &stubResult{rows: []stubRow{withLesson, withoutLesson}})
	repo := NewRepository(db)

	topics, availability, err := repo.Topics(context.Background())
	if err != nil {
		t.Fatalf("Topics() unexpected error: %v", err)
	}

	concepts := topics[0].Chapters()[0].Concepts()
	if len(concepts) != 2 || concepts[0].Slug().String() != "variables" || concepts[1].Slug().String() != "loops" {
		t.Fatalf("assembled concepts = %v, want [variables loops]", concepts)
	}

	variablesPath := curriculumapp.ConceptPath{TopicSlug: "go-basics", ConceptSlug: "variables"}
	loopsPath := curriculumapp.ConceptPath{TopicSlug: "go-basics", ConceptSlug: "loops"}
	if got, ok := availability[variablesPath]; !ok || got != 8 {
		t.Errorf("availability[%v] = (%d, %v), want (8, true)", variablesPath, got, ok)
	}
	if _, ok := availability[loopsPath]; ok {
		t.Errorf("availability contains %v, want absent (no lesson)", loopsPath)
	}
}
