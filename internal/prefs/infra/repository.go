// Package infra adapts the prefs domain and app layers to MySQL and HTTP —
// the vendor-facing details the domain and app layers must never see.
package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	selectPrefValueSQL = `SELECT value FROM prefs WHERE name = ?`

	// new.value reads the row alias; VALUES() is deprecated.
	upsertPrefSQL = `
INSERT INTO prefs (name, value)
VALUES (?, ?) AS new
ON DUPLICATE KEY UPDATE
    value = new.value`

	deletePrefSQL = `DELETE FROM prefs WHERE name = ?`
)

// Repository is the MySQL adapter for app.Repository.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Get returns the value stored under name. ok is false when no row exists —
// distinct from a stored empty string (the domain layer rejects an empty
// slug before Set is ever called).
func (r *Repository) Get(ctx context.Context, name string) (string, bool, error) {
	var value string
	err := r.db.QueryRowContext(ctx, selectPrefValueSQL, name).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("infra: get pref %q: %w", name, err)
	}
	return value, true, nil
}

func (r *Repository) Set(ctx context.Context, name, value string) error {
	if _, err := r.db.ExecContext(ctx, upsertPrefSQL, name, value); err != nil {
		return fmt.Errorf("infra: set pref %q: %w", name, err)
	}
	return nil
}

// Delete removes the row for name. Deleting an already-absent name affects
// zero rows and is not an error, which is what makes clearing a pref twice
// idempotent.
func (r *Repository) Delete(ctx context.Context, name string) error {
	if _, err := r.db.ExecContext(ctx, deletePrefSQL, name); err != nil {
		return fmt.Errorf("infra: delete pref %q: %w", name, err)
	}
	return nil
}
