package mysql

import (
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
