package infra

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"slices"
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
// lessons maps to a slice, not a single ref: concepts are unique per
// (chapter_id, slug), not per topic, so a test can seed two lessons under
// the same (topicSlug, conceptSlug) to exercise the duplicate-match guard
// lockLesson enforces — a real fixture MySQL's schema itself would allow.
type stubData struct {
	mu sync.Mutex

	lessons  map[progressConceptKey][]stubLessonRef
	progress map[int64]stubProgressRow

	// recallChecks models recall_checks: lessonID -> (case-folded question
	// -> every canonical row folding to it, each carrying its own kind).
	// Folded because real MySQL matches recall_checks.question under
	// utf8mb4_0900_ai_ci (case/accent insensitive) — the stub must fold the
	// same way, or a regression that hashes the caller's raw casing instead
	// of the resolved canonical text (R1) would go undetected by every test
	// here. A slice, not a single value, because two recall_checks rows in
	// the same lesson can fold to the same key (differing only by
	// case/accent) — real MySQL would return both, and resolveRecallCheck
	// must refuse to pick one.
	recallChecks map[int64]map[string][]seededRecallCheck
	attempts     []stubRecallAttempt

	locksMu sync.Mutex
	locks   map[int64]*sync.Mutex

	queryErr error
	execErr  error
}

// seededRecallCheck is one recall_checks row as the stub models it: the
// canonical (as-seeded) question text and its kind — both are what
// selectRecallCheckExistsSQL's real JOIN would return, never anything a
// caller of RecordAttempt supplies (R1's question fix and R1's kind fix
// are the same shape of bug, fixed the same way).
type seededRecallCheck struct {
	question string
	kind     string
}

type stubRecallAttempt struct {
	lessonID       int64
	checkKey       string
	kind           string
	confidence     string
	outcome        string
	selectedOption *string
	gradedBy       string
	createdAt      time.Time
}

func newStubData() *stubData {
	return &stubData{
		lessons:  map[progressConceptKey][]stubLessonRef{},
		progress: map[int64]stubProgressRow{},
	}
}

func (d *stubData) seedLesson(topicSlug, conceptSlug string, lessonID int64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.addLessonLocked(topicSlug, conceptSlug, lessonID)
}

func (d *stubData) seedProgress(topicSlug, conceptSlug string, lessonID int64, row stubProgressRow) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.addLessonLocked(topicSlug, conceptSlug, lessonID)
	d.progress[lessonID] = row
}

func (d *stubData) addLessonLocked(topicSlug, conceptSlug string, lessonID int64) {
	key := foldKey(topicSlug, conceptSlug)
	d.lessons[key] = append(d.lessons[key], stubLessonRef{id: lessonID, topic: topicSlug, concept: conceptSlug})
}

// foldQuestion approximates MySQL's utf8mb4_0900_ai_ci matching for
// recall_checks.question with a plain case fold — good enough to exercise
// the same case-insensitivity real MySQL exhibits (see
// TestRepositoryRecordAttempt_CaseVariantResolvesToCanonicalCheckKey),
// without reimplementing full accent-insensitive Unicode collation here.
func foldQuestion(question string) string { return strings.ToLower(question) }

// seedRecallCheck registers lessonID (creating its lesson row if needed) as
// owning question with the given kind, the same shape
// selectRecallCheckExistsSQL's JOIN verifies against a real recall_checks
// table. Looked up by foldQuestion, same as the real column's
// case-insensitive collation, but the CANONICAL (as-seeded) question and
// kind are what get returned to a matching query — never the caller's
// casing, and never a client-supplied kind (R1). Calling this twice for the
// same lesson with questions that fold to the same key (e.g. differing only
// by case) seeds BOTH — modelling two distinct recall_checks rows that
// happen to collide under MySQL's collation, for the ambiguous-match test.
func (d *stubData) seedRecallCheck(topicSlug, conceptSlug string, lessonID int64, question, kind string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.addLessonLocked(topicSlug, conceptSlug, lessonID)
	if d.recallChecks == nil {
		d.recallChecks = map[int64]map[string][]seededRecallCheck{}
	}
	if d.recallChecks[lessonID] == nil {
		d.recallChecks[lessonID] = map[string][]seededRecallCheck{}
	}
	key := foldQuestion(question)
	d.recallChecks[lessonID][key] = append(d.recallChecks[lessonID][key], seededRecallCheck{question: question, kind: kind})
}

func (d *stubData) attemptsFor(lessonID int64) []stubRecallAttempt {
	d.mu.Lock()
	defer d.mu.Unlock()
	var out []stubRecallAttempt
	for _, a := range d.attempts {
		if a.lessonID == lessonID {
			out = append(out, a)
		}
	}
	return out
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

func argOptionalString(a driver.NamedValue) *string {
	if a.Value == nil {
		return nil
	}
	s := a.Value.(string)
	return &s
}

// recallAttemptEnumColumns mirrors migrations/006_recall.sql's ENUM column
// definitions exactly. Without this, the stub would accept any string for
// type/confidence/outcome/graded_by — real MySQL under STRICT_TRANS_TABLES
// rejects an out-of-set ENUM value with error 1265 ("Data truncated"), so a
// Go-side array widened without a matching migration change would insert
// successfully here and only fail against the real database, exactly the
// blind spot that let R1's kind bug and a bad confidence/outcome/graded_by
// value alike go undetected by every test that only exercises this stub.
var recallAttemptEnumColumns = map[string][]string{
	"type":       {"short_answer", "mcq"},
	"confidence": {"guessed", "unsure", "confident"},
	"outcome":    {"correct", "incorrect"},
	"graded_by":  {"self", "llm"},
}

func validateRecallAttemptEnums(kind, confidence, outcome, gradedBy string) error {
	values := map[string]string{"type": kind, "confidence": confidence, "outcome": outcome, "graded_by": gradedBy}
	for column, value := range values {
		if !slices.Contains(recallAttemptEnumColumns[column], value) {
			return fmt.Errorf("learningstub: Data truncated for column '%s' at row 1 (simulated MySQL error 1265): %q is not in ENUM%v", column, value, recallAttemptEnumColumns[column])
		}
	}
	return nil
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

	case insertRecallAttemptSQL:
		lessonID := args[0].Value.(int64)
		checkKey := argString(args[1])
		kind := argString(args[2])
		confidence := argString(args[3])
		outcome := argString(args[4])
		selectedOption := argOptionalString(args[5])
		gradedBy := argString(args[6])
		createdAt := args[7].Value.(time.Time)

		if err := validateRecallAttemptEnums(kind, confidence, outcome, gradedBy); err != nil {
			return nil, err
		}

		c.data.mu.Lock()
		c.data.attempts = append(c.data.attempts, stubRecallAttempt{
			lessonID: lessonID, checkKey: checkKey, kind: kind, confidence: confidence,
			outcome: outcome, selectedOption: selectedOption, gradedBy: gradedBy, createdAt: createdAt,
		})
		c.data.mu.Unlock()
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
		// A regression that drops the surrounding transaction must fail
		// immediately, not acquire a lock nothing will ever release —
		// checked before locking anything, so this never hangs the suite.
		if c.tx == nil {
			return nil, errors.New("learningstub: FOR UPDATE requires an active transaction")
		}

		topicSlug, conceptSlug := argString(args[0]), argString(args[1])
		c.data.mu.Lock()
		refs := append([]stubLessonRef(nil), c.data.lessons[foldKey(topicSlug, conceptSlug)]...)
		c.data.mu.Unlock()
		if len(refs) == 0 {
			return &lessonForUpdateRows{}, nil
		}

		// The lock acquired here (never inside c.data.mu) models MySQL's row
		// lock: it can block this goroutine until a concurrent transaction
		// on the same lesson commits or rolls back.
		for _, ref := range refs {
			c.data.lockFor(ref.id).Lock()
			c.tx.lockedIDs = append(c.tx.lockedIDs, ref.id)
		}
		return &lessonForUpdateRows{rows: refs}, nil

	case selectProgressByLessonIDSQL:
		lessonID := args[0].Value.(int64)
		row, ok := c.data.progressOf(lessonID)
		if !ok {
			return &progressByIDRows{}, nil
		}
		return &progressByIDRows{rows: []stubProgressRow{row}}, nil

	case selectAllProgressSQL:
		return &allProgressRows{rows: c.allProgressRows()}, nil

	case selectRecallCheckExistsSQL:
		topicSlug, conceptSlug, question := argString(args[0]), argString(args[1]), argString(args[2])
		c.data.mu.Lock()
		refs := append([]stubLessonRef(nil), c.data.lessons[foldKey(topicSlug, conceptSlug)]...)
		c.data.mu.Unlock()

		var matches []recallCheckMatch
		for _, ref := range refs {
			c.data.mu.Lock()
			seeds := append([]seededRecallCheck(nil), c.data.recallChecks[ref.id][foldQuestion(question)]...)
			c.data.mu.Unlock()
			for _, seed := range seeds {
				matches = append(matches, recallCheckMatch{lessonID: ref.id, question: seed.question, kind: seed.kind})
			}
		}
		return &recallCheckExistsRows{rows: matches}, nil

	case selectLessonExistsSQL:
		topicSlug, conceptSlug := argString(args[0]), argString(args[1])
		c.data.mu.Lock()
		refs := c.data.lessons[foldKey(topicSlug, conceptSlug)]
		c.data.mu.Unlock()
		if len(refs) == 0 {
			return &singleIDRows{}, nil
		}
		return &singleIDRows{id: refs[0].id, has: true}, nil

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
	for _, refs := range c.data.lessons {
		for _, ref := range refs {
			row, ok := c.data.progress[ref.id]
			if !ok {
				continue
			}
			out = append(out, allProgressFixtureRow{topicSlug: ref.topic, conceptSlug: ref.concept, row: row})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].topicSlug != out[j].topicSlug {
			return out[i].topicSlug < out[j].topicSlug
		}
		return out[i].conceptSlug < out[j].conceptSlug
	})
	return out
}

// singleIDRows answers a single-column, at-most-one-row query
// (selectLessonExistsSQL) — has false means zero rows, mirroring
// sql.ErrNoRows on the real QueryRowContext.Scan path.
type singleIDRows struct {
	id   int64
	has  bool
	done bool
}

func (r *singleIDRows) Columns() []string { return []string{"l.id"} }
func (r *singleIDRows) Close() error      { return nil }

func (r *singleIDRows) Next(dest []driver.Value) error {
	if !r.has || r.done {
		return io.EOF
	}
	dest[0] = r.id
	r.done = true
	return nil
}

// recallCheckMatch is one row selectRecallCheckExistsSQL's real JOIN would
// return: l.id plus the CANONICAL (as-seeded) rc.question and rc.type —
// never the caller's casing, and never a client-supplied kind.
type recallCheckMatch struct {
	lessonID int64
	question string
	kind     string
}

// recallCheckExistsRows answers selectRecallCheckExistsSQL. It can hold more
// than one row: two recall_checks folding to the same key under
// case/accent-insensitive collation both come back, exactly like real MySQL
// would, so resolveRecallCheck's ambiguous-match guard has something real to
// detect rather than only ever seeing at most one candidate by construction.
type recallCheckExistsRows struct {
	rows []recallCheckMatch
	pos  int
}

func (r *recallCheckExistsRows) Columns() []string {
	return []string{"l.id", "rc.question", "rc.type"}
}
func (r *recallCheckExistsRows) Close() error { return nil }

func (r *recallCheckExistsRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.rows) {
		return io.EOF
	}
	m := r.rows[r.pos]
	dest[0] = m.lessonID
	dest[1] = []byte(m.question)
	dest[2] = []byte(m.kind)
	r.pos++
	return nil
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

// TestStubInsertRecallAttemptRejectsInvalidEnumValue proves the fidelity fix
// itself: without it, the stub accepted any string for
// type/confidence/outcome/graded_by, so a Go-side enum array widened without
// a matching migrations/006_recall.sql change would insert here without
// error and only fail against real MySQL under STRICT_TRANS_TABLES (error
// 1265). Exercised directly against db.ExecContext, bypassing the domain
// layer's own validation entirely — this is testing the STUB's fidelity to
// the real schema, not RecordAttempt's request-handling path.
func TestStubInsertRecallAttemptRejectsInvalidEnumValue(t *testing.T) {
	tests := []struct {
		name                                string
		kind, confidence, outcome, gradedBy string
	}{
		{name: "invalid type", kind: "essay", confidence: "guessed", outcome: "correct", gradedBy: "self"},
		{name: "invalid confidence", kind: "mcq", confidence: "maybe", outcome: "correct", gradedBy: "self"},
		{name: "invalid outcome", kind: "mcq", confidence: "guessed", outcome: "skipped", gradedBy: "self"},
		{name: "invalid graded_by", kind: "mcq", confidence: "guessed", outcome: "correct", gradedBy: "robot"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := openStubDB(t, newStubData())
			_, err := db.ExecContext(context.Background(), insertRecallAttemptSQL,
				int64(1), "checkkey", tt.kind, tt.confidence, tt.outcome, nil, tt.gradedBy, time.Now())
			if err == nil {
				t.Fatalf("ExecContext() error = nil, want a simulated MySQL 1265 error for an out-of-ENUM value")
			}
		})
	}
}
