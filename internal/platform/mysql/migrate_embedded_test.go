package mysql

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/themethaithian/self-learning/migrations"
)

func TestLoadMigrationsFromEmbeddedFS(t *testing.T) {
	all, err := loadMigrations(migrations.FS)
	if err != nil {
		t.Fatalf("loadMigrations(migrations.FS) unexpected error: %v", err)
	}
	if len(all) == 0 {
		t.Fatal("loadMigrations(migrations.FS) returned no migrations, want at least 001_curriculum.sql")
	}
	if all[0].version != 1 || all[0].filename != "001_curriculum.sql" {
		t.Fatalf("loadMigrations(migrations.FS)[0] = %+v, want version 1 / 001_curriculum.sql", all[0])
	}
	if len(all[0].stmts) != 5 {
		t.Fatalf("001_curriculum.sql parsed into %d statements, want 5 (one CREATE TABLE per curriculum table)", len(all[0].stmts))
	}
}

// TestAllMigrationsCreateTableIsIdempotent pins Migrate's own doc comment:
// MySQL DDL is non-transactional, so a migration that fails partway through
// (or crashes before its version row is recorded) leaves some of its DDL
// already applied — every CREATE TABLE must converge on re-run. This covers
// every embedded migration file, not just one.
func TestAllMigrationsCreateTableIsIdempotent(t *testing.T) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		t.Fatalf("read migrations dir: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		content, err := fs.ReadFile(migrations.FS, entry.Name())
		if err != nil {
			t.Fatalf("read %s: %v", entry.Name(), err)
		}
		assertEveryCreateTableIsIdempotent(t, entry.Name(), string(content))
	}
}

func assertEveryCreateTableIsIdempotent(t *testing.T, filename, content string) {
	t.Helper()
	const marker = "CREATE TABLE"

	for searchFrom := 0; ; {
		idx := strings.Index(content[searchFrom:], marker)
		if idx == -1 {
			return
		}
		idx += searchFrom

		rest := strings.TrimLeft(content[idx+len(marker):], " \t\n")
		if !strings.HasPrefix(rest, "IF NOT EXISTS") {
			t.Errorf("%s: %q at byte offset %d is missing IF NOT EXISTS", filename, marker, idx)
		}
		searchFrom = idx + len(marker)
	}
}

// TestAllMigrationsAlterTableAddColumnIsGuarded is ALTER's counterpart to
// TestAllMigrationsCreateTableIsIdempotent above, not a literal extension of
// it: CREATE TABLE's idempotency proof is a single following token (IF NOT
// EXISTS), but MySQL 8 has no ADD COLUMN IF NOT EXISTS at all, so there is no
// equivalent token to check for. The only safe pattern (see
// migrations/008_recall-check-explanation.sql) is a PREPARE/EXECUTE built
// from an information_schema existence check, so this test checks for that
// pattern's markers instead of a suffix keyword. It proves the guard's
// *syntax* is present, the same static-text limit the CREATE TABLE check
// already has — the guard's actual MySQL *behavior* (a restart after a
// crash mid-migration converges) is proven live, not by this test; see
// docs/tickets/aws-cert.md's AWS-S1 section for that evidence.
func TestAllMigrationsAlterTableAddColumnIsGuarded(t *testing.T) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		t.Fatalf("read migrations dir: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		content, err := fs.ReadFile(migrations.FS, entry.Name())
		if err != nil {
			t.Fatalf("read %s: %v", entry.Name(), err)
		}
		assertAddColumnIsGuarded(t, entry.Name(), string(content))
	}
}

func assertAddColumnIsGuarded(t *testing.T, filename, content string) {
	t.Helper()
	if !strings.Contains(content, "ADD COLUMN") {
		return
	}
	for _, want := range []string{"information_schema.COLUMNS", "PREPARE ", "EXECUTE ", "DEALLOCATE PREPARE"} {
		if !strings.Contains(content, want) {
			t.Errorf("%s: contains ADD COLUMN but missing guard marker %q — a bare ALTER TABLE ADD COLUMN "+
				"cannot survive a crash between the DDL and the schema_migrations INSERT (see migrate.go's doc comment)",
				filename, want)
		}
	}
}
