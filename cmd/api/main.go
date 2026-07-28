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

	curriculumapp "github.com/themethaithian/self-learning/internal/curriculum/app"
	curriculuminfra "github.com/themethaithian/self-learning/internal/curriculum/infra"
	learningapp "github.com/themethaithian/self-learning/internal/learning/app"
	learninginfra "github.com/themethaithian/self-learning/internal/learning/infra"
	"github.com/themethaithian/self-learning/internal/platform/config"
	"github.com/themethaithian/self-learning/internal/platform/httpserver"
	"github.com/themethaithian/self-learning/internal/platform/middleware"
	"github.com/themethaithian/self-learning/internal/platform/mysql"
	prefsapp "github.com/themethaithian/self-learning/internal/prefs/app"
	prefsinfra "github.com/themethaithian/self-learning/internal/prefs/infra"
)

// Shutdown returns as soon as in-flight requests finish; this deadline is
// only a cap, and must exceed the server's 60s WriteTimeout so a redeploy
// never kills a long-running LLM-grading request mid-flight.
const shutdownTimeout = 65 * time.Second

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

	curriculumRepo := curriculuminfra.NewRepository(db)
	curriculumService := curriculumapp.NewService(curriculumRepo)
	curriculumHandler := curriculuminfra.NewHandler(curriculumService, logger)

	prefsRepo := prefsinfra.NewRepository(db)
	prefsService := prefsapp.NewService(prefsRepo)
	prefsHandler := prefsinfra.NewHandler(prefsService, logger)

	learningRepo := learninginfra.NewRepository(db)
	learningService := learningapp.NewService(learningRepo)
	learningHandler := learninginfra.NewHandler(learningService, logger)

	mux := httpserver.NewMux(db, cfg.APIBearerToken, curriculumHandler, prefsHandler, learningHandler)
	handler := middleware.Logging(logger)(middleware.Recovery(logger)(middleware.CORS(cfg.CORSAllowedOrigin)(mux)))
	srv := httpserver.New(fmt.Sprintf(":%d", cfg.APIPort), handler)

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
