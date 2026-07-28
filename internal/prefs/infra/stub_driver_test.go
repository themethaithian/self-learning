package infra

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"sync"
	"testing"
)

// stubPrefsData is one test's in-memory prefs table, keyed exactly like the
// real schema's PRIMARY KEY (name) — this is a real key-value store, not a
// canned-rows fixture, so it can only answer Get correctly if Set/Delete
// actually wrote through the same query text production uses.
type stubPrefsData struct {
	mu    sync.Mutex
	table map[string]string

	queryErr error
	execErr  error
}

func newStubPrefsData() *stubPrefsData {
	return &stubPrefsData{table: make(map[string]string)}
}

var (
	stubPrefsMu       sync.Mutex
	stubPrefsRegistry = map[string]*stubPrefsData{}
	registerPrefsStub sync.Once
)

func openStubPrefsDB(t *testing.T, data *stubPrefsData) *sql.DB {
	t.Helper()
	registerPrefsStub.Do(func() {
		sql.Register("prefsstub", prefsStubDriver{})
	})

	name := t.Name()
	stubPrefsMu.Lock()
	stubPrefsRegistry[name] = data
	stubPrefsMu.Unlock()
	t.Cleanup(func() {
		stubPrefsMu.Lock()
		delete(stubPrefsRegistry, name)
		stubPrefsMu.Unlock()
	})

	db, err := sql.Open("prefsstub", name)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	return db
}

type prefsStubDriver struct{}

func (prefsStubDriver) Open(name string) (driver.Conn, error) {
	stubPrefsMu.Lock()
	data, ok := stubPrefsRegistry[name]
	stubPrefsMu.Unlock()
	if !ok {
		return nil, errors.New("prefsstub: no fixture registered for " + name)
	}
	return &prefsStubConn{data: data}, nil
}

type prefsStubConn struct {
	data *stubPrefsData
}

func (c *prefsStubConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prefsstub: Prepare not supported, use QueryContext/ExecContext")
}

func (c *prefsStubConn) Close() error { return nil }

func (c *prefsStubConn) Begin() (driver.Tx, error) {
	return nil, errors.New("prefsstub: transactions not supported")
}

// ExecContext dispatches by matching query against the same Go constant
// reference production sent, then maps args[0]/args[1] positionally as
// (name, value) regardless of what the SQL text's own column list says — a
// column-order mutation in upsertPrefSQL is invisible here. See
// TestUpsertPrefSQLShape in repository_test.go, which checks the literal
// SQL text for exactly that.
func (c *prefsStubConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	c.data.mu.Lock()
	defer c.data.mu.Unlock()

	if c.data.execErr != nil {
		return nil, c.data.execErr
	}

	switch query {
	case upsertPrefSQL:
		name, value := args[0].Value.(string), args[1].Value.(string)
		c.data.table[name] = value
		return driver.RowsAffected(1), nil
	case deletePrefSQL:
		name := args[0].Value.(string)
		delete(c.data.table, name)
		return driver.RowsAffected(1), nil
	default:
		return nil, fmt.Errorf("prefsstub: unrecognised exec query: %s", query)
	}
}

func (c *prefsStubConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	c.data.mu.Lock()
	defer c.data.mu.Unlock()

	if c.data.queryErr != nil {
		return nil, c.data.queryErr
	}

	switch query {
	case selectPrefValueSQL:
		name := args[0].Value.(string)
		value, ok := c.data.table[name]
		if !ok {
			return &prefsStubRows{}, nil
		}
		return &prefsStubRows{values: []string{value}}, nil
	default:
		return nil, fmt.Errorf("prefsstub: unrecognised query: %s", query)
	}
}

type prefsStubRows struct {
	values []string
	pos    int
}

func (r *prefsStubRows) Columns() []string { return []string{"value"} }
func (r *prefsStubRows) Close() error      { return nil }

func (r *prefsStubRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.values) {
		return io.EOF
	}
	dest[0] = []byte(r.values[r.pos])
	r.pos++
	return nil
}
