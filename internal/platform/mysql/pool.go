// Package mysql wires the MySQL connection pool and runs schema migrations.
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/themethaithian/self-learning/internal/platform/config"
)

// Sized for a single-user app behind a $6/mo 1GB VPS, not a high-traffic
// service: few concurrent requests, but connections must be recycled well
// under MySQL's default wait_timeout so they don't go stale.
const (
	maxOpenConns    = 10
	maxIdleConns    = 5
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = 2 * time.Minute
	pingTimeout     = 5 * time.Second
)

// Open creates the pooled *sql.DB and verifies connectivity with a ping.
func Open(ctx context.Context, cfg config.DB) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("mysql: open: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("mysql: ping: %w", err)
	}
	return db, nil
}

func dsn(cfg config.DB) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}
