package infra

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

// progressConceptKey is the case-folded lookup key for stubData.lessons —
// selectLessonForUpdateSQL relies on MySQL's ai_ci collation matching
// "B-Trees" to a stored "b-trees" row, so the stub folds case the same way
// and separately carries the canonical (as-stored) slugs to prove the
// repository echoes those back, not the caller's casing.
type progressConceptKey struct{ topicSlug, conceptSlug string }

func foldKey(topicSlug, conceptSlug string) progressConceptKey {
	return progressConceptKey{strings.ToLower(topicSlug), strings.ToLower(conceptSlug)}
}

type stubLessonRef struct {
	id      int64
	topic   string
	concept string
}

type stubProgressRow struct {
	state         string
	firstPassedAt *time.Time
	lastReadAt    *time.Time
}

// stubData is a real in-memory table keyed exactly like the schema's actual
// keys (lesson existence by topic+concept slug, progress by lesson id), not
// a canned-rows fixture — so it can only answer Transition/AllProgress
// correctly if the write path actually went through the same query text
// production uses. locks models MySQL's SELECT ... FOR UPDATE row lock: one
// real sync.Mutex per lesson id, held for the lifetime of one transaction.
type stubData struct {
	mu sync.Mutex

	lessons  map[progressConceptKey]stubLessonRef
	progress map[int64]stubProgressRow

	locksMu sync.Mutex
	locks   map[int64]*sync.Mutex

	queryErr error
	execErr  error
}

func newStubData() *stubData {
	return &stubData{
		lessons:  map[progressConceptKey]stubLessonRef{},
		progress: map[int64]stubProgressRow{},
	}
}

func (d *stubData) seedLesson(topicSlug, conceptSlug string, lessonID int64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.lessons[foldKey(topicSlug, conceptSlug)] = stubLessonRef{id: lessonID, topic: topicSlug, concept: conceptSlug}
}

func (d *stubData) seedProgress(topicSlug, conceptSlug string, lessonID int64, row stubProgressRow) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.lessons[foldKey(topicSlug, conceptSlug)] = stubLessonRef{id: lessonID, topic: topicSlug, concept: conceptSlug}
	d.progress[lessonID] = row
}

func (d *stubData) progressOf(lessonID int64) (stubProgressRow, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	row, ok := d.progress[lessonID]
	return row, ok
}

func (d *stubData) lockFor(lessonID int64) *sync.Mutex {
	d.locksMu.Lock()
	defer d.locksMu.Unlock()
	if d.locks == nil {
		d.locks = map[int64]*sync.Mutex{}
	}
	l, ok := d.locks[lessonID]
	if !ok {
		l = &sync.Mutex{}
		d.locks[lessonID] = l
	}
	return l
}

var (
	stubRegistryMu sync.Mutex
	stubRegistry   = map[string]*stubData{}
	registerStub   sync.Once
)

func openStubDB(t *testing.T, data *stubData) *sql.DB {
	t.Helper()
	registerStub.Do(func() {
		sql.Register("learningstub", stubDriver{})
	})

	name := t.Name()
	stubRegistryMu.Lock()
	stubRegistry[name] = data
	stubRegistryMu.Unlock()
	t.Cleanup(func() {
		stubRegistryMu.Lock()
		delete(stubRegistry, name)
		stubRegistryMu.Unlock()
	})

	db, err := sql.Open("learningstub", name)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

type stubDriver struct{}

func (stubDriver) Open(name string) (driver.Conn, error) {
	stubRegistryMu.Lock()
	data, ok := stubRegistry[name]
	stubRegistryMu.Unlock()
	if !ok {
		return nil, errors.New("learningstub: no fixture registered for " + name)
	}
	return &stubConn{data: data}, nil
}

// stubConn's tx field is non-nil exactly while a transaction opened on this
// connection is in flight, so ExecContext/QueryContext can record which
// per-lesson locks to release on Commit/Rollback.
type stubConn struct {
	data *stubData
	tx   *stubTx
}

func (c *stubConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("learningstub: Prepare not supported, use QueryContext/ExecContext")
}

func (c *stubConn) Close() error { return nil }

func (c *stubConn) Begin() (driver.Tx, error) {
	return nil, errors.New("learningstub: use BeginTx")
}

func (c *stubConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	tx := &stubTx{conn: c}
	c.tx = tx
	return tx, nil
}

type stubTx struct {
	conn      *stubConn
	lockedIDs []int64
}

func (tx *stubTx) Commit() error   { tx.release(); return nil }
func (tx *stubTx) Rollback() error { tx.release(); return nil }

func (tx *stubTx) release() {
	for _, id := range tx.lockedIDs {
		tx.conn.data.lockFor(id).Unlock()
	}
	tx.lockedIDs = nil
	tx.conn.tx = nil
}

func argString(a driver.NamedValue) string { return a.Value.(string) }

func argOptionalTime(a driver.NamedValue) *time.Time {
	if a.Value == nil {
		return nil
	}
	t := a.Value.(time.Time)
	return &t
}

// ExecContext dispatches by matching query against the same Go constant
// production sent, then applies real upsert/refresh semantics against
// stubData — a mutation of upsertProgressByLessonIDSQL's or
// refreshLastReadAtSQL's own text (wrong table, dropped WHERE, dropped
// COALESCE) is invisible to dispatch itself; the SQL-shape tests in
// repository_test.go guard the literal text instead.
func (c *stubConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	c.data.mu.Lock()
	execErr := c.data.execErr
	c.data.mu.Unlock()
	if execErr != nil {
		return nil, execErr
	}

	switch query {
	case upsertProgressByLessonIDSQL:
		lessonID := args[0].Value.(int64)
		state := argString(args[1])
		firstPassedAt := argOptionalTime(args[2])
		lastReadAt := args[3].Value.(time.Time)

		c.data.mu.Lock()
		existing := c.data.progress[lessonID]
		if existing.firstPassedAt != nil {
			firstPassedAt = existing.firstPassedAt
		}
		c.data.progress[lessonID] = stubProgressRow{state: state, firstPassedAt: firstPassedAt, lastReadAt: &lastReadAt}
		c.data.mu.Unlock()
		return driver.RowsAffected(1), nil

	case refreshLastReadAtSQL:
		lastReadAt := args[0].Value.(time.Time)
		lessonID := args[1].Value.(int64)

		c.data.mu.Lock()
		row, ok := c.data.progress[lessonID]
		if ok {
			row.lastReadAt = &lastReadAt
			c.data.progress[lessonID] = row
		}
		c.data.mu.Unlock()
		if !ok {
			return driver.RowsAffected(0), nil
		}
		return driver.RowsAffected(1), nil

	default:
		return nil, fmt.Errorf("learningstub: unrecognised exec query: %s", query)
	}
}

func (c *stubConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	c.data.mu.Lock()
	queryErr := c.data.queryErr
	c.data.mu.Unlock()
	if queryErr != nil {
		return nil, queryErr
	}

	switch query {
	case selectLessonForUpdateSQL:
		topicSlug, conceptSlug := argString(args[0]), argString(args[1])
		c.data.mu.Lock()
		ref, ok := c.data.lessons[foldKey(topicSlug, conceptSlug)]
		c.data.mu.Unlock()
		if !ok {
			return &lessonForUpdateRows{}, nil
		}

		// The lock acquired here (never inside c.data.mu) models MySQL's row
		// lock: it can block this goroutine until a concurrent transaction
		// on the same lesson commits or rolls back.
		c.data.lockFor(ref.id).Lock()
		if c.tx != nil {
			c.tx.lockedIDs = append(c.tx.lockedIDs, ref.id)
		}
		return &lessonForUpdateRows{rows: []stubLessonRef{ref}}, nil

	case selectProgressByLessonIDSQL:
		lessonID := args[0].Value.(int64)
		row, ok := c.data.progressOf(lessonID)
		if !ok {
			return &progressByIDRows{}, nil
		}
		return &progressByIDRows{rows: []stubProgressRow{row}}, nil

	case selectAllProgressSQL:
		return &allProgressRows{rows: c.allProgressRows()}, nil

	default:
		return nil, fmt.Errorf("learningstub: unrecognised query: %s", query)
	}
}

// allProgressRows joins the write-side stubData back the way
// selectAllProgressSQL's real INNER JOINs would, sorted by (topic, concept)
// to match its ORDER BY.
func (c *stubConn) allProgressRows() []allProgressFixtureRow {
	c.data.mu.Lock()
	defer c.data.mu.Unlock()

	var out []allProgressFixtureRow
	for _, ref := range c.data.lessons {
		row, ok := c.data.progress[ref.id]
		if !ok {
			continue
		}
		out = append(out, allProgressFixtureRow{topicSlug: ref.topic, conceptSlug: ref.concept, row: row})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].topicSlug != out[j].topicSlug {
			return out[i].topicSlug < out[j].topicSlug
		}
		return out[i].conceptSlug < out[j].conceptSlug
	})
	return out
}

type lessonForUpdateRows struct {
	rows []stubLessonRef
	pos  int
}

func (r *lessonForUpdateRows) Columns() []string { return []string{"l.id", "t.slug", "co.slug"} }
func (r *lessonForUpdateRows) Close() error      { return nil }

func (r *lessonForUpdateRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.rows) {
		return io.EOF
	}
	ref := r.rows[r.pos]
	dest[0] = ref.id
	dest[1] = []byte(ref.topic)
	dest[2] = []byte(ref.concept)
	r.pos++
	return nil
}

type progressByIDRows struct {
	rows []stubProgressRow
	pos  int
}

func (r *progressByIDRows) Columns() []string {
	return []string{"state", "first_passed_at", "last_read_at"}
}
func (r *progressByIDRows) Close() error { return nil }

func (r *progressByIDRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.rows) {
		return io.EOF
	}
	row := r.rows[r.pos]
	dest[0] = []byte(row.state)
	dest[1] = timeDriverValue(row.firstPassedAt)
	dest[2] = timeDriverValue(row.lastReadAt)
	r.pos++
	return nil
}

type allProgressFixtureRow struct {
	topicSlug, conceptSlug string
	row                    stubProgressRow
}

type allProgressRows struct {
	rows []allProgressFixtureRow
	pos  int
}

func (r *allProgressRows) Columns() []string {
	return []string{"t.slug", "co.slug", "lp.state", "lp.first_passed_at", "lp.last_read_at"}
}
func (r *allProgressRows) Close() error { return nil }

func (r *allProgressRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.rows) {
		return io.EOF
	}
	fr := r.rows[r.pos]
	dest[0] = []byte(fr.topicSlug)
	dest[1] = []byte(fr.conceptSlug)
	dest[2] = []byte(fr.row.state)
	dest[3] = timeDriverValue(fr.row.firstPassedAt)
	dest[4] = timeDriverValue(fr.row.lastReadAt)
	r.pos++
	return nil
}

func timeDriverValue(t *time.Time) driver.Value {
	if t == nil {
		return nil
	}
	return *t
}
