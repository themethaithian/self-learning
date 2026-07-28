package infra

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/themethaithian/self-learning/migrations"
)

func TestRepositoryGet_Unset(t *testing.T) {
	db := openStubPrefsDB(t, newStubPrefsData())
	repo := NewRepository(db)

	value, ok, err := repo.Get(context.Background(), "focus_track")
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("Get() ok = true, want false; value=%q", value)
	}
}

func TestRepositorySetThenGet(t *testing.T) {
	db := openStubPrefsDB(t, newStubPrefsData())
	repo := NewRepository(db)

	if err := repo.Set(context.Background(), "focus_track", "ai-systems"); err != nil {
		t.Fatalf("Set() unexpected error: %v", err)
	}
	value, ok, err := repo.Get(context.Background(), "focus_track")
	if err != nil || !ok || value != "ai-systems" {
		t.Fatalf("Get() = (%q, %v, %v), want (ai-systems, true, nil)", value, ok, err)
	}
}

// TestRepositorySet_IdempotentUpsert proves repeated Set calls don't error
// and leave the right single value — it does not prove ON DUPLICATE KEY
// UPDATE is present in the SQL (the stub's map would behave identically for
// a plain INSERT); TestUpsertPrefSQLShape checks that instead.
func TestRepositorySet_IdempotentUpsert(t *testing.T) {
	db := openStubPrefsDB(t, newStubPrefsData())
	repo := NewRepository(db)

	for i := 0; i < 2; i++ {
		if err := repo.Set(context.Background(), "focus_track", "ai-systems"); err != nil {
			t.Fatalf("Set() call %d unexpected error: %v", i, err)
		}
	}
	value, ok, err := repo.Get(context.Background(), "focus_track")
	if err != nil || !ok || value != "ai-systems" {
		t.Fatalf("Get() = (%q, %v, %v), want (ai-systems, true, nil)", value, ok, err)
	}
}

func TestRepositorySet_UpdatesExistingValue(t *testing.T) {
	db := openStubPrefsDB(t, newStubPrefsData())
	repo := NewRepository(db)

	if err := repo.Set(context.Background(), "focus_track", "ai-systems"); err != nil {
		t.Fatalf("Set() unexpected error: %v", err)
	}
	if err := repo.Set(context.Background(), "focus_track", "go"); err != nil {
		t.Fatalf("Set() unexpected error: %v", err)
	}

	value, ok, err := repo.Get(context.Background(), "focus_track")
	if err != nil || !ok || value != "go" {
		t.Fatalf("Get() = (%q, %v, %v), want (go, true, nil)", value, ok, err)
	}
}

func TestRepositoryDelete_ClearsRow(t *testing.T) {
	db := openStubPrefsDB(t, newStubPrefsData())
	repo := NewRepository(db)
	if err := repo.Set(context.Background(), "focus_track", "ai-systems"); err != nil {
		t.Fatalf("Set() unexpected error: %v", err)
	}

	if err := repo.Delete(context.Background(), "focus_track"); err != nil {
		t.Fatalf("Delete() unexpected error: %v", err)
	}

	_, ok, err := repo.Get(context.Background(), "focus_track")
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if ok {
		t.Fatal("Get() ok = true after Delete, want false")
	}
}

func TestRepositoryDelete_AbsentRowIsNoOp(t *testing.T) {
	db := openStubPrefsDB(t, newStubPrefsData())
	repo := NewRepository(db)

	if err := repo.Delete(context.Background(), "focus_track"); err != nil {
		t.Fatalf("Delete() on an absent row unexpected error: %v", err)
	}
}

func TestRepositoryGet_PropagatesError(t *testing.T) {
	data := newStubPrefsData()
	data.queryErr = errors.New("connection refused")
	db := openStubPrefsDB(t, data)
	repo := NewRepository(db)

	_, _, err := repo.Get(context.Background(), "focus_track")
	if err == nil {
		t.Fatal("Get() expected error, got nil")
	}
}

func TestRepositorySet_PropagatesError(t *testing.T) {
	data := newStubPrefsData()
	data.execErr = errors.New("connection refused")
	db := openStubPrefsDB(t, data)
	repo := NewRepository(db)

	if err := repo.Set(context.Background(), "focus_track", "ai-systems"); err == nil {
		t.Fatal("Set() expected error, got nil")
	}
}

func TestRepositoryDelete_PropagatesError(t *testing.T) {
	data := newStubPrefsData()
	data.execErr = errors.New("connection refused")
	db := openStubPrefsDB(t, data)
	repo := NewRepository(db)

	if err := repo.Delete(context.Background(), "focus_track"); err == nil {
		t.Fatal("Delete() expected error, got nil")
	}
}

// TestPrefsTableNameConsistency guards against the migration and the Go-side
// SQL drifting onto different table names — nothing else would notice, since
// no test here touches a real MySQL schema.
func TestPrefsTableNameConsistency(t *testing.T) {
	migrationSQL, err := migrations.FS.ReadFile("005_prefs.sql")
	if err != nil {
		t.Fatalf("read 005_prefs.sql: %v", err)
	}
	if !strings.Contains(string(migrationSQL), "CREATE TABLE IF NOT EXISTS prefs") {
		t.Fatalf("005_prefs.sql does not contain %q", "CREATE TABLE IF NOT EXISTS prefs")
	}

	stmts := map[string]string{
		"selectPrefValueSQL": selectPrefValueSQL,
		"upsertPrefSQL":      upsertPrefSQL,
		"deletePrefSQL":      deletePrefSQL,
	}
	for name, stmt := range stmts {
		if !strings.Contains(stmt, "prefs") {
			t.Errorf("%s = %q, want it to reference the prefs table", name, stmt)
		}
	}
}

// TestUpsertPrefSQLShape asserts the statement text itself, independent of
// the stub above: the stub interprets args[0]/args[1] positionally as
// (name, value) regardless of what the SQL's own column list says, so a
// mutation that swaps the column list order (INSERT INTO prefs (value,
// name)) would pass every stub-backed test here while silently writing rows
// backwards against a real MySQL server. Only asserting on the literal SQL
// text catches that.
func TestUpsertPrefSQLShape(t *testing.T) {
	mustContain := []string{
		"INSERT INTO prefs (name, value)",
		"VALUES (?, ?) AS new",
		"ON DUPLICATE KEY UPDATE",
		"value = new.value",
	}
	for _, want := range mustContain {
		if !strings.Contains(upsertPrefSQL, want) {
			t.Errorf("upsertPrefSQL = %q, want it to contain %q", upsertPrefSQL, want)
		}
	}
}
