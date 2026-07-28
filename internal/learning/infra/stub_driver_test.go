package infra

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"sort"
	"sync"
	"testing"
	"time"
)

// progressConceptKey mirrors selectConceptProgressSQL's and touchProgressSQL's
// own WHERE predicate — (topic slug, concept slug) — so a repository test
// only passes if the lookup actually discriminates by topic, the same way a
// real query missing "t.slug = ?" would fail it.
type progressConceptKey struct{ topicSlug, conceptSlug string }

type stubProgressRow struct {
	state         string
	firstPassedAt *time.Time
	lastReadAt    *time.Time
}

// stubData is a real in-memory table keyed exactly like the schema's actual
// keys (lesson existence by topic+concept slug, progress by lesson id), not
// a canned-rows fixture — so it can only answer ConceptProgress/AllProgress
// correctly if Upsert/Touch actually wrote through the same query text
// production uses.
type stubData struct {
	mu sync.Mutex

	lessons  map[progressConceptKey]int64
	progress map[int64]stubProgressRow

	queryErr error
	execErr  error
}

func newStubData() *stubData {
	return &stubData{lessons: map[progressConceptKey]int64{}, progress: map[int64]stubProgressRow{}}
}

func (d *stubData) seedLesson(topicSlug, conceptSlug string, lessonID int64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.lessons[progressConceptKey{topicSlug, conceptSlug}] = lessonID
}

func (d *stubData) seedProgress(topicSlug, conceptSlug string, lessonID int64, row stubProgressRow) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.lessons[progressConceptKey{topicSlug, conceptSlug}] = lessonID
	d.progress[lessonID] = row
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
	db.SetMaxOpenConns(1)
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

type stubConn struct {
	data *stubData
}

func (c *stubConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("learningstub: Prepare not supported, use QueryContext/ExecContext")
}

func (c *stubConn) Close() error { return nil }

func (c *stubConn) Begin() (driver.Tx, error) {
	return nil, errors.New("learningstub: transactions not supported")
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
// production sent, then applies real upsert/touch semantics against
// stubData — a mutation of upsertProgressSQL's or touchProgressSQL's own
// text (wrong table, dropped WHERE, dropped COALESCE) is invisible to
// dispatch itself; TestUpsertProgressSQLShape and
// TestTouchProgressSQLShape guard the literal text instead.
func (c *stubConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	c.data.mu.Lock()
	defer c.data.mu.Unlock()

	if c.data.execErr != nil {
		return nil, c.data.execErr
	}

	switch query {
	case upsertProgressSQL:
		lessonID := args[0].Value.(int64)
		state := argString(args[1])
		firstPassedAt := argOptionalTime(args[2])
		lastReadAt := args[3].Value.(time.Time)

		existing := c.data.progress[lessonID]
		if existing.firstPassedAt != nil {
			firstPassedAt = existing.firstPassedAt
		}
		c.data.progress[lessonID] = stubProgressRow{state: state, firstPassedAt: firstPassedAt, lastReadAt: &lastReadAt}
		return driver.RowsAffected(1), nil

	case touchProgressSQL:
		lastReadAt := args[0].Value.(time.Time)
		topicSlug, conceptSlug := argString(args[1]), argString(args[2])
		lessonID, ok := c.data.lessons[progressConceptKey{topicSlug, conceptSlug}]
		if !ok {
			return driver.RowsAffected(0), nil
		}
		row := c.data.progress[lessonID]
		row.lastReadAt = &lastReadAt
		c.data.progress[lessonID] = row
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
	defer c.data.mu.Unlock()

	if c.data.queryErr != nil {
		return nil, c.data.queryErr
	}

	switch query {
	case selectConceptProgressSQL:
		topicSlug, conceptSlug := argString(args[0]), argString(args[1])
		lessonID, ok := c.data.lessons[progressConceptKey{topicSlug, conceptSlug}]
		if !ok {
			return &conceptProgressRows{}, nil
		}
		row, hasProgress := c.data.progress[lessonID]
		return &conceptProgressRows{rows: []conceptProgressFixtureRow{{lessonID: lessonID, hasProgress: hasProgress, row: row}}}, nil

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
	var out []allProgressFixtureRow
	for key, lessonID := range c.data.lessons {
		row, ok := c.data.progress[lessonID]
		if !ok {
			continue
		}
		out = append(out, allProgressFixtureRow{topicSlug: key.topicSlug, conceptSlug: key.conceptSlug, row: row})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].topicSlug != out[j].topicSlug {
			return out[i].topicSlug < out[j].topicSlug
		}
		return out[i].conceptSlug < out[j].conceptSlug
	})
	return out
}

type conceptProgressFixtureRow struct {
	lessonID    int64
	hasProgress bool
	row         stubProgressRow
}

type conceptProgressRows struct {
	rows []conceptProgressFixtureRow
	pos  int
}

func (r *conceptProgressRows) Columns() []string {
	return []string{"l.id", "lp.state", "lp.first_passed_at", "lp.last_read_at"}
}
func (r *conceptProgressRows) Close() error { return nil }

func (r *conceptProgressRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.rows) {
		return io.EOF
	}
	fr := r.rows[r.pos]
	dest[0] = fr.lessonID
	if fr.hasProgress {
		dest[1] = []byte(fr.row.state)
		dest[2] = timeDriverValue(fr.row.firstPassedAt)
		dest[3] = timeDriverValue(fr.row.lastReadAt)
	} else {
		dest[1], dest[2], dest[3] = nil, nil, nil
	}
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
