package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/themethaithian/self-learning/migrations"
)

type migrationFile struct {
	version  int
	name     string
	filename string
	stmts    []string
}

var migrationFilenamePattern = regexp.MustCompile(`^(\d+)_([a-zA-Z0-9-]+)\.sql$`)

// Migrate applies every embedded migration not yet recorded in
// schema_migrations, each inside its own transaction, in ascending version
// order. Safe to call on every process start: already-applied versions are
// skipped, so re-running is a no-op.
func Migrate(ctx context.Context, db *sql.DB) error {
	return applyPending(ctx, db, migrations.FS)
}

func applyPending(ctx context.Context, db *sql.DB, fsys fs.FS) error {
	if err := ensureSchemaMigrationsTable(ctx, db); err != nil {
		return err
	}

	all, err := loadMigrations(fsys)
	if err != nil {
		return err
	}

	applied, err := appliedVersions(ctx, db)
	if err != nil {
		return err
	}

	for _, m := range pendingMigrations(all, applied) {
		if err := applyOne(ctx, db, m); err != nil {
			return fmt.Errorf("mysql: apply %s: %w", m.filename, err)
		}
	}
	return nil
}

// loadMigrations parses and orders every *.sql file in fsys. It touches no
// database, so the ordering/parsing logic is unit-testable without MySQL.
func loadMigrations(fsys fs.FS) ([]migrationFile, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("mysql: read migrations dir: %w", err)
	}

	seenVersions := make(map[int]string, len(entries))
	all := make([]migrationFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		m, err := parseMigrationFilename(entry.Name())
		if err != nil {
			return nil, err
		}
		if existing, ok := seenVersions[m.version]; ok {
			return nil, fmt.Errorf("mysql: duplicate migration version %d (%s and %s)", m.version, existing, m.filename)
		}
		seenVersions[m.version] = m.filename

		content, err := fs.ReadFile(fsys, m.filename)
		if err != nil {
			return nil, fmt.Errorf("mysql: read %s: %w", m.filename, err)
		}
		m.stmts = splitStatements(string(content))
		all = append(all, m)
	}

	sort.Slice(all, func(i, j int) bool { return all[i].version < all[j].version })
	return all, nil
}

func parseMigrationFilename(filename string) (migrationFile, error) {
	match := migrationFilenamePattern.FindStringSubmatch(filename)
	if match == nil {
		return migrationFile{}, fmt.Errorf("mysql: invalid migration filename %q, want NNN_name.sql", filename)
	}

	version, err := strconv.Atoi(match[1])
	if err != nil {
		return migrationFile{}, fmt.Errorf("mysql: invalid migration version in %q: %w", filename, err)
	}
	return migrationFile{version: version, name: match[2], filename: filename}, nil
}

// pendingMigrations returns the migrations in all whose version is absent
// from applied, preserving the ascending order loadMigrations produced.
func pendingMigrations(all []migrationFile, applied map[int]bool) []migrationFile {
	pending := make([]migrationFile, 0, len(all))
	for _, m := range all {
		if !applied[m.version] {
			pending = append(pending, m)
		}
	}
	return pending
}

// splitStatements splits a migration file into individual statements on ';'.
// Migration files are hand-written schema DDL, never string literals
// containing ';', so a plain split is safe here.
func splitStatements(sqlText string) []string {
	raw := strings.Split(sqlText, ";")
	stmts := make([]string, 0, len(raw))
	for _, s := range raw {
		s = strings.TrimSpace(s)
		if s != "" {
			stmts = append(stmts, s)
		}
	}
	return stmts
}

func ensureSchemaMigrationsTable(ctx context.Context, db *sql.DB) error {
	const stmt = `CREATE TABLE IF NOT EXISTS schema_migrations (
		version BIGINT PRIMARY KEY,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB`
	if _, err := db.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("mysql: create schema_migrations: %w", err)
	}
	return nil
}

func appliedVersions(ctx context.Context, db *sql.DB) (map[int]bool, error) {
	rows, err := db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("mysql: read schema_migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("mysql: scan schema_migrations: %w", err)
		}
		applied[version] = true
	}
	return applied, rows.Err()
}

func applyOne(ctx context.Context, db *sql.DB, m migrationFile) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	for _, stmt := range m.stmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("exec statement: %w", err)
		}
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES (?)", m.version); err != nil {
		return fmt.Errorf("record version: %w", err)
	}

	return tx.Commit()
}
