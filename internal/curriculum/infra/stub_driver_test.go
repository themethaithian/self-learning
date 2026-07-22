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

// stubResult is one test's fixture, plus what the stub observed while
// serving it, so a test can assert on the query text and on cleanup calls a
// real driver would make. All access goes through stubMu — the same lock
// openStubDB uses for the driver registry — so this file has one locking
// rule, not two.
type stubResult struct {
	rows     []stubRow
	queryErr error
	nextErr  error

	lastQuery string
	lastRows  *stubRows
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
	return nil, errors.New("stub: transactions not supported")
}

func (c *stubConn) QueryContext(ctx context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	stubMu.Lock()
	c.result.lastQuery = query
	stubMu.Unlock()

	if c.result.queryErr != nil {
		return nil, c.result.queryErr
	}

	rows := &stubRows{
		columns: parseSelectColumns(query),
		rows:    c.result.rows,
		nextErr: c.result.nextErr,
	}
	stubMu.Lock()
	c.result.lastRows = rows
	stubMu.Unlock()
	return rows, nil
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
