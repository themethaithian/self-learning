package mysql

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestParseMigrationFilename(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		wantVersion int
		wantName    string
		wantErr     bool
	}{
		{name: "valid three-digit version", filename: "001_curriculum.sql", wantVersion: 1, wantName: "curriculum"},
		{name: "valid multi-digit version", filename: "042_lessons-learning.sql", wantVersion: 42, wantName: "lessons-learning"},
		{name: "missing version prefix", filename: "curriculum.sql", wantErr: true},
		{name: "missing name suffix", filename: "001.sql", wantErr: true},
		{name: "wrong extension", filename: "001_curriculum.txt", wantErr: true},
		{name: "no underscore separator", filename: "001curriculum.sql", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := parseMigrationFilename(tt.filename)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseMigrationFilename(%q) expected error, got nil", tt.filename)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseMigrationFilename(%q) unexpected error: %v", tt.filename, err)
			}
			if m.version != tt.wantVersion || m.name != tt.wantName {
				t.Fatalf("parseMigrationFilename(%q) = {version: %d, name: %q}, want {version: %d, name: %q}",
					tt.filename, m.version, m.name, tt.wantVersion, tt.wantName)
			}
		})
	}
}

func TestLoadMigrationsOrdering(t *testing.T) {
	fsys := fstest.MapFS{
		"010_ten.sql": {Data: []byte("CREATE TABLE ten (id INT);")},
		"002_two.sql": {Data: []byte("CREATE TABLE two (id INT);")},
		"001_one.sql": {Data: []byte("CREATE TABLE one (id INT);")},
	}

	got, err := loadMigrations(fsys)
	if err != nil {
		t.Fatalf("loadMigrations() unexpected error: %v", err)
	}

	wantVersions := []int{1, 2, 10}
	if len(got) != len(wantVersions) {
		t.Fatalf("loadMigrations() returned %d migrations, want %d", len(got), len(wantVersions))
	}
	for i, want := range wantVersions {
		if got[i].version != want {
			t.Fatalf("loadMigrations()[%d].version = %d, want %d (numeric order, not lexical)", i, got[i].version, want)
		}
	}
}

func TestLoadMigrationsDuplicateVersion(t *testing.T) {
	fsys := fstest.MapFS{
		"001_first.sql":  {Data: []byte("CREATE TABLE first (id INT);")},
		"001_second.sql": {Data: []byte("CREATE TABLE second (id INT);")},
	}

	_, err := loadMigrations(fsys)
	if err == nil {
		t.Fatal("loadMigrations() expected error for duplicate version, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate migration version 1") {
		t.Fatalf("loadMigrations() error = %q, want it to mention duplicate version 1", err.Error())
	}
}

func TestLoadMigrationsInvalidFilename(t *testing.T) {
	fsys := fstest.MapFS{
		"not-a-migration.sql": {Data: []byte("CREATE TABLE x (id INT);")},
	}

	_, err := loadMigrations(fsys)
	if err == nil {
		t.Fatal("loadMigrations() expected error for invalid filename, got nil")
	}
}

func TestLoadMigrationsIgnoresNonSQLFiles(t *testing.T) {
	fsys := fstest.MapFS{
		"001_one.sql": {Data: []byte("CREATE TABLE one (id INT);")},
		"README.md":   {Data: []byte("not a migration")},
	}

	got, err := loadMigrations(fsys)
	if err != nil {
		t.Fatalf("loadMigrations() unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("loadMigrations() returned %d migrations, want 1 (README.md should be ignored)", len(got))
	}
}

func TestPendingMigrations(t *testing.T) {
	all := []migrationFile{
		{version: 1, filename: "001_one.sql"},
		{version: 2, filename: "002_two.sql"},
		{version: 3, filename: "003_three.sql"},
	}

	tests := []struct {
		name    string
		applied map[int]bool
		want    []int
	}{
		{name: "none applied", applied: map[int]bool{}, want: []int{1, 2, 3}},
		{name: "first applied", applied: map[int]bool{1: true}, want: []int{2, 3}},
		{name: "all applied", applied: map[int]bool{1: true, 2: true, 3: true}, want: []int{}},
		{name: "middle applied", applied: map[int]bool{2: true}, want: []int{1, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pendingMigrations(all, tt.applied)
			if len(got) != len(tt.want) {
				t.Fatalf("pendingMigrations() = %d migrations, want %d", len(got), len(tt.want))
			}
			for i, wantVersion := range tt.want {
				if got[i].version != wantVersion {
					t.Fatalf("pendingMigrations()[%d].version = %d, want %d", i, got[i].version, wantVersion)
				}
			}
		})
	}
}

func TestSplitStatements(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "single statement",
			sql:  "CREATE TABLE foo (id INT)",
			want: []string{"CREATE TABLE foo (id INT)"},
		},
		{
			name: "multiple statements",
			sql:  "CREATE TABLE foo (id INT);\nCREATE TABLE bar (id INT);",
			want: []string{"CREATE TABLE foo (id INT)", "CREATE TABLE bar (id INT)"},
		},
		{
			name: "trailing semicolon and blank lines ignored",
			sql:  "CREATE TABLE foo (id INT);\n\n\n",
			want: []string{"CREATE TABLE foo (id INT)"},
		},
		{
			name: "empty input",
			sql:  "",
			want: []string{},
		},
		{
			// Known limitation, checked in on purpose: a plain split cannot
			// tell a statement-ending ';' apart from one inside a string
			// literal or comment. Migration files must avoid both.
			name: "known limitation: semicolon inside a string literal splits mid-statement",
			sql:  "INSERT INTO foo (name) VALUES ('a;b');",
			want: []string{"INSERT INTO foo (name) VALUES ('a", "b')"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitStatements(tt.sql)
			if len(got) != len(tt.want) {
				t.Fatalf("splitStatements(%q) = %v, want %v", tt.sql, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("splitStatements(%q)[%d] = %q, want %q", tt.sql, i, got[i], tt.want[i])
				}
			}
		})
	}
}
