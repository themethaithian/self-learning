package infra

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// stubRow maps a column expression exactly as it appears in the SELECT list
// (e.g. "co.title") to its value: a string or int for a real value, nil for
// SQL NULL. Deriving Columns() from the query text itself, rather than from
// a hardcoded list, keeps the stub from encoding the same column-order
// assumption production Scan already makes — the two must agree by
// evaluating the same source, not by construction.
type stubRow map[string]any

func stubFullRow(track, topicSlug, topicTitle string, topicPos int, chapterSlug, chapterTitle string, chapterPos int, conceptSlug, conceptTitle, conceptOutline string, conceptPos int) stubRow {
	return stubRow{
		"t.track": track, "t.slug": topicSlug, "t.title": topicTitle, "t.position": topicPos,
		"c.slug": chapterSlug, "c.title": chapterTitle, "c.position": chapterPos,
		"co.slug": conceptSlug, "co.title": conceptTitle, "co.outline": conceptOutline, "co.position": conceptPos,
		"l.est_minutes": nil,
	}
}

func stubTopicOnlyRow(track, slug, title string, pos int) stubRow {
	return stubRow{
		"t.track": track, "t.slug": slug, "t.title": title, "t.position": pos,
		"c.slug": nil, "c.title": nil, "c.position": nil,
		"co.slug": nil, "co.title": nil, "co.outline": nil, "co.position": nil,
		"l.est_minutes": nil,
	}
}

func stubChapterOnlyRow(track, topicSlug, topicTitle string, topicPos int, chapterSlug, chapterTitle string, chapterPos int) stubRow {
	return stubRow{
		"t.track": track, "t.slug": topicSlug, "t.title": topicTitle, "t.position": topicPos,
		"c.slug": chapterSlug, "c.title": chapterTitle, "c.position": chapterPos,
		"co.slug": nil, "co.title": nil, "co.outline": nil, "co.position": nil,
		"l.est_minutes": nil,
	}
}

// stubResult is one test's fixture, plus what the stub observed while
// serving it, so a test can assert on the query text and on cleanup calls a
// real driver would make. All access goes through stubMu — the same lock
// openStubDB uses for the driver registry — so this file has one locking
// rule, not two.
type stubResult struct {
	rows     []stubRow
	queryErr error
	nextErr  error

	beginErr error
	// execErrOnCall, if non-zero, is the 1-based ExecContext call number that
	// fails with execErr, simulating a write that fails partway through a
	// multi-statement import. 0 means every call succeeds.
	execErrOnCall int
	execErr       error

	lastQuery     string
	lastQueryArgs []driver.NamedValue
	lastRows      *stubRows

	execCalls    []stubExecCall
	txBegan      int
	txCommitted  int
	txRolledBack int

	// tables backs ExecContext with real upsert/delete semantics keyed by
	// the schema's actual unique constraints, so SaveTopic can be tested for
	// idempotency and orphan cleanup, not just "some write happened".
	// QueryContext reads through it too when a test leaves rows unset, so
	// Repository.Topics after Repository.SaveTopic sees what was actually
	// written.
	tables *stubTables
}

func (r *stubResult) query() string {
	stubMu.Lock()
	defer stubMu.Unlock()
	return r.lastQuery
}

func (r *stubResult) queryArgs() []driver.NamedValue {
	stubMu.Lock()
	defer stubMu.Unlock()
	return r.lastQueryArgs
}

// seedConcept registers a concept the selectConceptIDSQL lookup can
// resolve, keyed by exactly the columns its WHERE clause filters on. A
// SaveLesson test seeds this instead of a canned rows fixture so the
// lookup's topic-scoping is actually exercised: two concepts sharing a
// slug under different topics only work if resolveConceptID's args reach
// the right (topic_slug, concept_slug) pair.
func (r *stubResult) seedConcept(topicSlug, conceptSlug string, id int64) {
	stubMu.Lock()
	defer stubMu.Unlock()
	if r.tables == nil {
		r.tables = newStubTables()
	}
	r.tables.lessonConcepts[lessonConceptKey{topicSlug: topicSlug, conceptSlug: conceptSlug}] = id
}

func (r *stubResult) driverRows() *stubRows {
	stubMu.Lock()
	defer stubMu.Unlock()
	return r.lastRows
}

func (r *stubResult) execLog() []stubExecCall {
	stubMu.Lock()
	defer stubMu.Unlock()
	return append([]stubExecCall(nil), r.execCalls...)
}

func (r *stubResult) txCounts() (began, committed, rolledBack int) {
	stubMu.Lock()
	defer stubMu.Unlock()
	return r.txBegan, r.txCommitted, r.txRolledBack
}

var (
	stubMu       sync.Mutex
	stubData     = map[string]*stubResult{}
	registerStub sync.Once
)

// openStubDB gives Repository a *sql.DB backed entirely by an in-memory
// driver.Driver, so Topics can be exercised without Docker or a real MySQL
// connection.
func openStubDB(t *testing.T, result *stubResult) *sql.DB {
	t.Helper()
	registerStub.Do(func() {
		sql.Register("curriculumstub", stubDriver{})
	})

	name := t.Name()
	stubMu.Lock()
	stubData[name] = result
	stubMu.Unlock()
	t.Cleanup(func() {
		stubMu.Lock()
		delete(stubData, name)
		stubMu.Unlock()
	})

	db, err := sql.Open("curriculumstub", name)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	return db
}

type stubDriver struct{}

func (stubDriver) Open(name string) (driver.Conn, error) {
	stubMu.Lock()
	result, ok := stubData[name]
	stubMu.Unlock()
	if !ok {
		return nil, errors.New("stub: no fixture registered for " + name)
	}
	return &stubConn{result: result}, nil
}

type stubConn struct {
	result *stubResult
}

func (c *stubConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("stub: Prepare not supported, use QueryContext")
}

func (c *stubConn) Close() error { return nil }

func (c *stubConn) Begin() (driver.Tx, error) {
	stubMu.Lock()
	defer stubMu.Unlock()
	if c.result.beginErr != nil {
		return nil, c.result.beginErr
	}
	c.result.txBegan++
	return &stubTx{result: c.result}, nil
}

type stubTx struct {
	result *stubResult
}

func (tx *stubTx) Commit() error {
	stubMu.Lock()
	defer stubMu.Unlock()
	tx.result.txCommitted++
	return nil
}

func (tx *stubTx) Rollback() error {
	stubMu.Lock()
	defer stubMu.Unlock()
	tx.result.txRolledBack++
	return nil
}

// stubExecCall records one ExecContext invocation exactly as the driver saw
// it — query text, arguments through placeholders, and the id/rows-affected
// it answered with — so a test can assert statement order, argument values,
// and that a child's parent-id argument really is the id the parent's own
// upsert returned.
type stubExecCall struct {
	query        string
	args         []driver.NamedValue
	lastInsertID int64
	rowsAffected int64
}

// ExecContext backs both plain db.ExecContext and every statement run inside
// a transaction: with MaxOpenConns(1) the same *stubConn serves the
// connection Begin() was called on, so this single method is enough to
// observe a whole transaction's writes in order. Everything — the
// error-injection check, the table mutation, and the call log — happens
// under one lock acquisition, so a concurrent test can never observe a
// half-applied call.
func (c *stubConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	stubMu.Lock()
	defer stubMu.Unlock()

	callNum := len(c.result.execCalls) + 1
	if c.result.execErr != nil && (c.result.execErrOnCall == 0 || c.result.execErrOnCall == callNum) {
		c.result.execCalls = append(c.result.execCalls, stubExecCall{query: query, args: args})
		return nil, c.result.execErr
	}

	if c.result.tables == nil {
		c.result.tables = newStubTables()
	}
	res, err := c.result.tables.apply(query, args)
	if err != nil {
		c.result.execCalls = append(c.result.execCalls, stubExecCall{query: query, args: args})
		return nil, err
	}

	lastID, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()
	c.result.execCalls = append(c.result.execCalls, stubExecCall{
		query: query, args: args, lastInsertID: lastID, rowsAffected: rowsAffected,
	})
	return res, nil
}

type stubExecResult struct {
	lastInsertID int64
	rowsAffected int64
}

func (r stubExecResult) LastInsertId() (int64, error) { return r.lastInsertID, nil }
func (r stubExecResult) RowsAffected() (int64, error) { return r.rowsAffected, nil }

// chapterKey and conceptKey mirror the schema's real unique constraints —
// chapters on (topic_id, slug), concepts on (chapter_id, slug) — so the stub
// can tell "this is the same row, maybe changed" from "this is a new row"
// the same way MySQL's ON DUPLICATE KEY UPDATE does.
type chapterKey struct {
	topicID int64
	slug    string
}

type conceptKey struct {
	chapterID int64
	slug      string
}

type stubTopicRow struct {
	id       int64
	track    string
	title    string
	position int64
}

type stubChapterRow struct {
	id       int64
	title    string
	position int64
}

type stubConceptRow struct {
	id       int64
	title    string
	outline  string
	position int64
}

type stubLessonRow struct {
	id         int64
	version    int64
	titleEn    string
	estMinutes int64
	bodyMd     string
	refs       string
}

type stubRecallCheckRow struct {
	id             int64
	position       int64
	kind           string
	question       string
	expectedAnswer string
	options        *string
	explanation    *string
}

// lessonConceptKey mirrors resolveConceptID's own WHERE predicate — (topic
// slug, concept slug) — so the concept lookup a SaveLesson test exercises
// actually discriminates by topic, the way the real query's `t.slug = ?`
// must, rather than answering the same canned row for any topic.
type lessonConceptKey struct {
	topicSlug   string
	conceptSlug string
}

// stubTables is a tiny in-memory relational store standing in for MySQL:
// enough upsert and stale-row-delete semantics to exercise idempotency and
// orphan cleanup, without a real database. A single id counter spans all
// three tables (unlike real independent per-table AUTO_INCREMENT) so a
// mis-threaded parent id — e.g. a concept bound to its topic's id instead of
// its chapter's — always shows up as a wrong number rather than a
// coincidentally correct one.
type stubTables struct {
	nextID int64

	topics   map[string]stubTopicRow
	chapters map[chapterKey]stubChapterRow
	concepts map[conceptKey]stubConceptRow

	// lessons is keyed by concept_id, mirroring the schema's UNIQUE KEY
	// uniq_lessons_concept. recallChecks is keyed by lesson_id: a
	// delete-then-insert table, not an upsert, so it holds a slice per
	// lesson rather than a keyed row.
	lessons      map[int64]stubLessonRow
	recallChecks map[int64][]stubRecallCheckRow

	// lessonConcepts seeds the concept lookup SaveLesson issues before it
	// ever upserts a lesson. It is independent of the concepts map above
	// (keyed by chapter_id, populated only by SaveTopic writes) because a
	// SaveLesson test has no reason to build a whole topic/chapter/concept
	// tree just to seed one concept id.
	lessonConcepts map[lessonConceptKey]int64
}

func newStubTables() *stubTables {
	return &stubTables{
		topics:         make(map[string]stubTopicRow),
		chapters:       make(map[chapterKey]stubChapterRow),
		concepts:       make(map[conceptKey]stubConceptRow),
		lessons:        make(map[int64]stubLessonRow),
		recallChecks:   make(map[int64][]stubRecallCheckRow),
		lessonConcepts: make(map[lessonConceptKey]int64),
	}
}

// apply dispatches on exact query text for the three static upsert
// statements, and on prefix for the delete statements, whose placeholder
// count varies with the number of children being kept.
func (tb *stubTables) apply(query string, args []driver.NamedValue) (driver.Result, error) {
	switch {
	case query == upsertTopicSQL:
		return tb.upsertTopic(args)
	case query == upsertChapterSQL:
		return tb.upsertChapter(args)
	case query == upsertConceptSQL:
		return tb.upsertConcept(args)
	case query == upsertLessonSQL:
		return tb.upsertLesson(args)
	case query == deleteRecallChecksSQL:
		return tb.deleteRecallChecks(args)
	case query == insertRecallCheckSQL:
		return tb.insertRecallCheck(args)
	case strings.HasPrefix(query, "DELETE FROM concepts WHERE chapter_id = ?"):
		return tb.deleteStaleConcepts(args)
	case strings.HasPrefix(query, "DELETE FROM concepts WHERE chapter_id IN"):
		return tb.deleteConceptsOfStaleChapters(args)
	case strings.HasPrefix(query, "DELETE FROM chapters WHERE topic_id = ?"):
		return tb.deleteStaleChapters(args)
	default:
		return nil, fmt.Errorf("stub: unrecognised exec query: %s", query)
	}
}

func argString(a driver.NamedValue) string { return a.Value.(string) }
func argInt64(a driver.NamedValue) int64   { return a.Value.(int64) }

// argOptionalString distinguishes a genuine SQL NULL argument (the
// recall_checks.options column for a short_answer check) from a JSON
// string argument (an mcq's marshalled options) — the two must never be
// confused the way a plain argString would.
func argOptionalString(a driver.NamedValue) *string {
	if a.Value == nil {
		return nil
	}
	s := argString(a)
	return &s
}

func (tb *stubTables) upsertTopic(args []driver.NamedValue) (driver.Result, error) {
	track, slug, title, position := argString(args[0]), argString(args[1]), argString(args[2]), argInt64(args[3])
	row := stubTopicRow{track: track, title: title, position: position}

	if existing, ok := tb.topics[slug]; ok {
		row.id = existing.id
		tb.topics[slug] = row
		if existing.track == row.track && existing.title == row.title && existing.position == row.position {
			return stubExecResult{lastInsertID: row.id, rowsAffected: 0}, nil
		}
		return stubExecResult{lastInsertID: row.id, rowsAffected: 2}, nil
	}

	tb.nextID++
	row.id = tb.nextID
	tb.topics[slug] = row
	return stubExecResult{lastInsertID: row.id, rowsAffected: 1}, nil
}

func (tb *stubTables) upsertChapter(args []driver.NamedValue) (driver.Result, error) {
	topicID, slug, title, position := argInt64(args[0]), argString(args[1]), argString(args[2]), argInt64(args[3])
	key := chapterKey{topicID: topicID, slug: slug}
	row := stubChapterRow{title: title, position: position}

	if existing, ok := tb.chapters[key]; ok {
		row.id = existing.id
		tb.chapters[key] = row
		if existing.title == row.title && existing.position == row.position {
			return stubExecResult{lastInsertID: row.id, rowsAffected: 0}, nil
		}
		return stubExecResult{lastInsertID: row.id, rowsAffected: 2}, nil
	}

	tb.nextID++
	row.id = tb.nextID
	tb.chapters[key] = row
	return stubExecResult{lastInsertID: row.id, rowsAffected: 1}, nil
}

func (tb *stubTables) upsertConcept(args []driver.NamedValue) (driver.Result, error) {
	chapterID, slug := argInt64(args[0]), argString(args[1])
	title, outline, position := argString(args[2]), argString(args[3]), argInt64(args[4])
	key := conceptKey{chapterID: chapterID, slug: slug}
	row := stubConceptRow{title: title, outline: outline, position: position}

	if existing, ok := tb.concepts[key]; ok {
		row.id = existing.id
		tb.concepts[key] = row
		if existing.title == row.title && existing.outline == row.outline && existing.position == row.position {
			return stubExecResult{lastInsertID: row.id, rowsAffected: 0}, nil
		}
		return stubExecResult{lastInsertID: row.id, rowsAffected: 2}, nil
	}

	tb.nextID++
	row.id = tb.nextID
	tb.concepts[key] = row
	return stubExecResult{lastInsertID: row.id, rowsAffected: 1}, nil
}

func (tb *stubTables) upsertLesson(args []driver.NamedValue) (driver.Result, error) {
	conceptID := argInt64(args[0])
	row := stubLessonRow{
		version: argInt64(args[1]), titleEn: argString(args[2]),
		estMinutes: argInt64(args[3]), bodyMd: argString(args[4]), refs: argString(args[5]),
	}

	if existing, ok := tb.lessons[conceptID]; ok {
		row.id = existing.id
		tb.lessons[conceptID] = row
		if existing == row {
			return stubExecResult{lastInsertID: row.id, rowsAffected: 0}, nil
		}
		return stubExecResult{lastInsertID: row.id, rowsAffected: 2}, nil
	}

	tb.nextID++
	row.id = tb.nextID
	tb.lessons[conceptID] = row
	return stubExecResult{lastInsertID: row.id, rowsAffected: 1}, nil
}

func (tb *stubTables) deleteRecallChecks(args []driver.NamedValue) (driver.Result, error) {
	lessonID := argInt64(args[0])
	affected := int64(len(tb.recallChecks[lessonID]))
	delete(tb.recallChecks, lessonID)
	return stubExecResult{rowsAffected: affected}, nil
}

func (tb *stubTables) insertRecallCheck(args []driver.NamedValue) (driver.Result, error) {
	lessonID := argInt64(args[0])
	row := stubRecallCheckRow{
		position: argInt64(args[1]), kind: argString(args[2]),
		question: argString(args[3]), expectedAnswer: argString(args[4]), options: argOptionalString(args[5]),
		explanation: argOptionalString(args[6]),
	}

	tb.nextID++
	row.id = tb.nextID
	tb.recallChecks[lessonID] = append(tb.recallChecks[lessonID], row)
	return stubExecResult{lastInsertID: row.id, rowsAffected: 1}, nil
}

// lookupConceptID answers selectConceptIDSQL's "WHERE t.slug = ? AND
// co.slug = ?" by both args, not just the second — dropping the topic-slug
// predicate here would silently resolve to whichever topic seeded the
// slug first, exactly the bug a real query missing t.slug = ? would cause.
func (tb *stubTables) lookupConceptID(args []driver.NamedValue) []stubRow {
	if len(args) != 2 {
		return nil
	}
	key := lessonConceptKey{topicSlug: argString(args[0]), conceptSlug: argString(args[1])}
	id, ok := tb.lessonConcepts[key]
	if !ok {
		return nil
	}
	return []stubRow{{"co.id": int(id)}}
}

// lookupLessonByConcept answers selectLessonByConceptSQL by resolving
// (topic slug, concept slug) through the same lessonConcepts map
// resolveConceptID uses, then looking up that concept's lesson row — so a
// concept slug that only exists under a different topic (or a concept with
// no lesson saved) both correctly produce no row, the same way a query that
// dropped the t.slug predicate would NOT.
func (tb *stubTables) lookupLessonByConcept(args []driver.NamedValue) []stubRow {
	if len(args) != 2 {
		return nil
	}
	key := lessonConceptKey{topicSlug: argString(args[0]), conceptSlug: argString(args[1])}
	conceptID, ok := tb.lessonConcepts[key]
	if !ok {
		return nil
	}
	lesson, ok := tb.lessons[conceptID]
	if !ok {
		return nil
	}
	return []stubRow{{
		"l.id": int(lesson.id), "l.version": int(lesson.version), "l.title_en": lesson.titleEn,
		"l.est_minutes": int(lesson.estMinutes), "l.body_md": lesson.bodyMd, "l.refs": lesson.refs,
	}}
}

// lookupRecallChecksByLesson answers selectRecallChecksByLessonSQL, sorted
// by position — recallChecks is a delete-then-insert table (see
// stubTables' field comment), so insertion order alone cannot be trusted to
// reflect position order the way ORDER BY position does in production.
func (tb *stubTables) lookupRecallChecksByLesson(args []driver.NamedValue) []stubRow {
	if len(args) != 1 {
		return nil
	}
	lessonID := argInt64(args[0])
	checks := append([]stubRecallCheckRow(nil), tb.recallChecks[lessonID]...)
	sort.Slice(checks, func(i, j int) bool { return checks[i].position < checks[j].position })

	rows := make([]stubRow, len(checks))
	for i, c := range checks {
		var options any
		if c.options != nil {
			options = *c.options
		}
		var explanation any
		if c.explanation != nil {
			explanation = *c.explanation
		}
		rows[i] = stubRow{
			"position": int(c.position), "type": c.kind, "question": c.question,
			"expected_answer": c.expectedAnswer, "options": options, "explanation": explanation,
		}
	}
	return rows
}

func keepSet(args []driver.NamedValue) map[string]bool {
	keep := make(map[string]bool, len(args))
	for _, a := range args {
		keep[argString(a)] = true
	}
	return keep
}

func (tb *stubTables) deleteStaleConcepts(args []driver.NamedValue) (driver.Result, error) {
	chapterID := argInt64(args[0])
	keep := keepSet(args[1:])

	var affected int64
	for key := range tb.concepts {
		if key.chapterID == chapterID && !keep[key.slug] {
			delete(tb.concepts, key)
			affected++
		}
	}
	return stubExecResult{rowsAffected: affected}, nil
}

func (tb *stubTables) deleteConceptsOfStaleChapters(args []driver.NamedValue) (driver.Result, error) {
	topicID := argInt64(args[0])
	keep := keepSet(args[1:])

	staleChapterIDs := make(map[int64]bool)
	for key, row := range tb.chapters {
		if key.topicID == topicID && !keep[key.slug] {
			staleChapterIDs[row.id] = true
		}
	}

	var affected int64
	for key := range tb.concepts {
		if staleChapterIDs[key.chapterID] {
			delete(tb.concepts, key)
			affected++
		}
	}
	return stubExecResult{rowsAffected: affected}, nil
}

func (tb *stubTables) deleteStaleChapters(args []driver.NamedValue) (driver.Result, error) {
	topicID := argInt64(args[0])
	keep := keepSet(args[1:])

	var affected int64
	for key := range tb.chapters {
		if key.topicID == topicID && !keep[key.slug] {
			delete(tb.chapters, key)
			affected++
		}
	}
	return stubExecResult{rowsAffected: affected}, nil
}

// toRows replays the write-side table state as topicsQuery's LEFT JOIN would
// see it, so a test can call Repository.Topics right after Repository.
// SaveTopic and get an answer reflecting what SaveTopic actually persisted —
// including any reconciliation deletes — through the same assembleTree path
// production uses, instead of re-asserting on stub internals. A topic with
// no matching chapters, or a chapter with no matching concepts, must still
// produce a row (concept/chapter columns NULL) exactly like the real LEFT
// JOIN does — TestRepositoryTopics_ChildlessTopicFailsLoudly pins that
// contract, and a lingering childless row is exactly what a missing DELETE
// leaves behind.
func (tb *stubTables) toRows() []stubRow {
	var out []stubRow
	for topicSlug, tp := range tb.topics {
		chapterRows := 0
		for chKey, ch := range tb.chapters {
			if chKey.topicID != tp.id {
				continue
			}
			chapterRows++

			conceptRows := 0
			for coKey, co := range tb.concepts {
				if coKey.chapterID != ch.id {
					continue
				}
				conceptRows++
				out = append(out, stubFullRow(
					tp.track, topicSlug, tp.title, int(tp.position),
					chKey.slug, ch.title, int(ch.position),
					coKey.slug, co.title, co.outline, int(co.position),
				))
			}
			if conceptRows == 0 {
				out = append(out, stubChapterOnlyRow(tp.track, topicSlug, tp.title, int(tp.position), chKey.slug, ch.title, int(ch.position)))
			}
		}
		if chapterRows == 0 {
			out = append(out, stubTopicOnlyRow(tp.track, topicSlug, tp.title, int(tp.position)))
		}
	}
	return out
}

func (c *stubConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	stubMu.Lock()
	c.result.lastQuery = query
	c.result.lastQueryArgs = args
	rows := c.result.rows
	if rows == nil && c.result.tables != nil {
		switch query {
		case topicsQuery:
			rows = c.result.tables.toRows()
		case selectConceptIDSQL:
			rows = c.result.tables.lookupConceptID(args)
		case selectLessonByConceptSQL:
			rows = c.result.tables.lookupLessonByConcept(args)
		case selectRecallChecksByLessonSQL:
			rows = c.result.tables.lookupRecallChecksByLesson(args)
		}
	}
	queryErr := c.result.queryErr
	nextErr := c.result.nextErr
	stubMu.Unlock()

	if queryErr != nil {
		return nil, queryErr
	}

	driverRows := &stubRows{
		columns: parseSelectColumns(query),
		rows:    rows,
		nextErr: nextErr,
	}
	stubMu.Lock()
	c.result.lastRows = driverRows
	stubMu.Unlock()
	return driverRows, nil
}

// parseSelectColumns extracts the SELECT list of a "SELECT ... FROM ..."
// query as written, so the stub's column set can never drift from what the
// query actually asks for.
func parseSelectColumns(query string) []string {
	upper := strings.ToUpper(query)
	start := strings.Index(upper, "SELECT") + len("SELECT")
	end := strings.Index(upper, "FROM")

	parts := strings.Split(query[start:end], ",")
	cols := make([]string, len(parts))
	for i, p := range parts {
		cols[i] = strings.TrimSpace(p)
	}
	return cols
}

type stubRows struct {
	columns []string
	rows    []stubRow
	nextErr error
	pos     int
	closed  bool
}

func (r *stubRows) Columns() []string { return r.columns }

func (r *stubRows) Close() error {
	r.closed = true
	return nil
}

// Next returns nextErr, if set, once the fixture rows are exhausted — it
// simulates a connection that dies mid-stream, distinct from a clean io.EOF.
// A fixture row missing one of the query's columns is a fixture bug, not a
// NULL, so it fails with its own message rather than defaulting silently.
func (r *stubRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.rows) {
		if r.nextErr != nil {
			return r.nextErr
		}
		return io.EOF
	}

	row := r.rows[r.pos]
	for i, col := range r.columns {
		val, ok := row[col]
		if !ok {
			return fmt.Errorf("stub: fixture row missing column %q", col)
		}
		dest[i] = toDriverValue(val)
	}
	r.pos++
	return nil
}

func toDriverValue(v any) driver.Value {
	switch val := v.(type) {
	case nil:
		return nil
	case string:
		return []byte(val)
	case int:
		return []byte(strconv.Itoa(val))
	default:
		panic(fmt.Sprintf("stub: unsupported fixture value type %T", v))
	}
}
