// Command api runs the self-learning HTTP API: migrate the DB, serve
// requests, shut down gracefully on SIGINT/SIGTERM.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/themethaithian/self-learning/internal/platform/config"
	"github.com/themethaithian/self-learning/internal/platform/httpserver"
	"github.com/themethaithian/self-learning/internal/platform/middleware"
	"github.com/themethaithian/self-learning/internal/platform/mysql"
)

const shutdownTimeout = 10 * time.Second

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := mysql.Open(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if err := mysql.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	mux := httpserver.NewMux(db)
	handler := middleware.Logging(logger)(middleware.Recovery(logger)(mux))
	srv := httpserver.New(":"+apiPort(), handler)

	return serve(ctx, logger, srv)
}

func serve(ctx context.Context, logger *slog.Logger, srv *http.Server) error {
	errCh := make(chan error, 1)
	go func() {
		logger.Info("http server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	logger.Info("server stopped")
	return nil
}

func apiPort() string {
	if p := os.Getenv("API_PORT"); p != "" {
		return p
	}
	return "8080"
}
