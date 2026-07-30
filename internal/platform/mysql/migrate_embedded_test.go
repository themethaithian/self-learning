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

// TestAllMigrationsNonIdempotentAlterIsGuarded replaces an earlier version
// of this test (TestAllMigrationsAlterTableAddColumnIsGuarded) that did four
// whole-file strings.Contains checks and claimed parity with
// TestAllMigrationsCreateTableIsIdempotent above — that claim was false, and
// the check was bypassable at least five ways: a bare ALTER appended after a
// real guard elsewhere in the file (whole-file Contains sees the guard
// markers and never notices the extra unguarded statement); lowercase SQL
// (MySQL accepts it, Contains("ADD COLUMN") does not); a `--` comment merely
// quoting the guard's marker words next to a bare, unguarded ALTER; a bare
// CREATE INDEX; a bare ALTER TABLE ... DROP COLUMN — none of ADD COLUMN,
// DROP COLUMN, or CREATE INDEX has an IF NOT EXISTS/IF EXISTS form in MySQL
// 8, so all three need the same guard style as migrations/008's ADD COLUMN.
//
// This test is a genuinely different instrument from the CREATE TABLE one,
// not the same kind of check applied to ALTER: it parses each file into the
// exact statements applyOne actually executes (splitStatements, after
// stripping line comments) and rejects any TOP-LEVEL statement — not any
// substring anywhere in the file — that is a non-idempotent ALTER/CREATE
// INDEX. A correctly guarded ADD COLUMN lives inside the string literal of a
// SET ... = IF(...) statement, which starts with "SET", not "ALTER TABLE",
// so it never matches. The one blind spot left is the same one
// splitStatements already documents: a string literal (e.g. an
// explanation's text) containing ';' would split mid-statement — out of
// scope for this test the same way it is out of scope for splitStatements
// itself.
func TestAllMigrationsNonIdempotentAlterIsGuarded(t *testing.T) {
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
		assertNoUnguardedNonIdempotentDDL(t, entry.Name(), string(content))
	}
}

// stripLineComments removes "-- ..." line comments before splitting into
// statements, so a comment that merely quotes a guard's marker words next
// to a real, unguarded statement cannot pass this test — the exact "C"
// bypass a reviewer's probe file demonstrated against the whole-file
// Contains version this test replaces.
func stripLineComments(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if idx := strings.Index(line, "--"); idx != -1 {
			lines[i] = line[:idx]
		}
	}
	return strings.Join(lines, "\n")
}

func assertNoUnguardedNonIdempotentDDL(t *testing.T, filename, content string) {
	t.Helper()
	for _, stmt := range splitStatements(stripLineComments(content)) {
		up := strings.ToUpper(strings.Join(strings.Fields(stmt), " "))
		switch {
		case strings.HasPrefix(up, "ALTER TABLE") && strings.Contains(up, "ADD COLUMN"):
			t.Errorf("%s: bare top-level %q — MySQL 8 has no ADD COLUMN IF NOT EXISTS, so this "+
				"cannot survive a crash between the DDL and the schema_migrations INSERT (see "+
				"migrate.go's doc comment); guard it via information_schema.COLUMNS + "+
				"PREPARE/EXECUTE, the way migrations/008_recall-check-explanation.sql does",
				filename, stmt)
		case strings.HasPrefix(up, "ALTER TABLE") && strings.Contains(up, "DROP COLUMN"):
			t.Errorf("%s: bare top-level %q — MySQL 8 has no DROP COLUMN IF EXISTS, same "+
				"non-idempotency risk as ADD COLUMN", filename, stmt)
		case strings.HasPrefix(up, "CREATE INDEX"):
			t.Errorf("%s: bare top-level %q — MySQL has no CREATE INDEX IF NOT EXISTS", filename, stmt)
		}
	}
}
