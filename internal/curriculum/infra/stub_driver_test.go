package infra

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
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
	}
}

func stubTopicOnlyRow(track, slug, title string, pos int) stubRow {
	return stubRow{
		"t.track": track, "t.slug": slug, "t.title": title, "t.position": pos,
		"c.slug": nil, "c.title": nil, "c.position": nil,
		"co.slug": nil, "co.title": nil, "co.outline": nil, "co.position": nil,
	}
}

func stubChapterOnlyRow(track, topicSlug, topicTitle string, topicPos int, chapterSlug, chapterTitle string, chapterPos int) stubRow {
	return stubRow{
		"t.track": track, "t.slug": topicSlug, "t.title": topicTitle, "t.position": topicPos,
		"c.slug": chapterSlug, "c.title": chapterTitle, "c.position": chapterPos,
		"co.slug": nil, "co.title": nil, "co.outline": nil, "co.position": nil,
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

	lastQuery string
	lastRows  *stubRows

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
}

func newStubTables() *stubTables {
	return &stubTables{
		topics:   make(map[string]stubTopicRow),
		chapters: make(map[chapterKey]stubChapterRow),
		concepts: make(map[conceptKey]stubConceptRow),
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

func (c *stubConn) QueryContext(ctx context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	stubMu.Lock()
	c.result.lastQuery = query
	rows := c.result.rows
	if rows == nil && c.result.tables != nil && query == topicsQuery {
		rows = c.result.tables.toRows()
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
