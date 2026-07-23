// Command import-lessons loads every lesson JSON file under a directory
// (content/lessons/<topic>/<concept>.json) and upserts it into MySQL. It
// must run after cmd/import-curriculum: a lesson resolves its concept by
// (topic slug, concept slug) and fails loudly if that concept does not
// exist yet.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	curriculuminfra "github.com/themethaithian/self-learning/internal/curriculum/infra"
	"github.com/themethaithian/self-learning/internal/platform/config"
	"github.com/themethaithian/self-learning/internal/platform/mysql"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dir := flag.String("dir", "content/lessons", "directory of lesson JSON files to import")
	flag.Parse()

	if err := run(context.Background(), logger, *dir); err != nil {
		logger.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger, dir string) error {
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

	files, err := lessonFiles(dir)
	if err != nil {
		return fmt.Errorf("list lesson files in %s: %w", dir, err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no *.json files found under %s (check -dir)", dir)
	}

	repo := curriculuminfra.NewRepository(db)

	for _, file := range files {
		topicSlug, lesson, err := curriculuminfra.LoadLesson(file)
		if err != nil {
			return fmt.Errorf("import %s: %w", file, err)
		}
		inserted, err := repo.SaveLesson(ctx, topicSlug, lesson)
		if err != nil {
			return fmt.Errorf("import %s: %w", file, err)
		}

		status := "updated"
		if inserted {
			status = "inserted"
		}
		logger.Info("lesson upserted",
			"file", file,
			"topic", topicSlug,
			"concept", lesson.Slug().String(),
			"recall_checks", len(lesson.RecallChecks()),
			"status", status,
		)
	}

	logger.Info("lesson import complete", "files", len(files))
	return nil
}

// os.ReadDir returns entries sorted by name at both levels, so lessonFiles
// returns paths in a stable order with no extra sort needed.
func lessonFiles(dir string) ([]string, error) {
	topicDirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, td := range topicDirs {
		if !td.IsDir() {
			continue
		}
		topicPath := filepath.Join(dir, td.Name())
		entries, err := os.ReadDir(topicPath)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
				continue
			}
			files = append(files, filepath.Join(topicPath, e.Name()))
		}
	}
	return files, nil
}
