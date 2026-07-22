// Command import-curriculum loads every curriculum JSON file in a directory
// and upserts it into MySQL: each file is the source of truth for its topic,
// so any chapter or concept missing from the file is deleted inside the same
// transaction as the upsert, and the import fails if a lesson still
// references a concept being deleted.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	curriculumapp "github.com/themethaithian/self-learning/internal/curriculum/app"
	curriculuminfra "github.com/themethaithian/self-learning/internal/curriculum/infra"
	"github.com/themethaithian/self-learning/internal/platform/config"
	"github.com/themethaithian/self-learning/internal/platform/mysql"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dir := flag.String("dir", "content/curriculum", "directory of curriculum JSON files to import")
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

	files, err := curriculumFiles(dir)
	if err != nil {
		return fmt.Errorf("list curriculum files in %s: %w", dir, err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no *.json files found in %s (check -dir)", dir)
	}

	importSvc := curriculumapp.NewImportService(curriculuminfra.FileLoader{}, curriculuminfra.NewRepository(db))

	totalConcepts := 0
	for _, file := range files {
		topic, err := importSvc.ImportFile(ctx, file)
		if err != nil {
			return fmt.Errorf("import %s: %w", file, err)
		}

		concepts := 0
		for _, ch := range topic.Chapters() {
			concepts += len(ch.Concepts())
		}
		totalConcepts += concepts
		logger.Info("curriculum topic upserted",
			"file", file,
			"topic", topic.Slug().String(),
			"chapters", len(topic.Chapters()),
			"concepts", concepts,
		)
	}

	logger.Info("curriculum import complete", "files", len(files), "concepts", totalConcepts)
	return nil
}

// os.ReadDir already returns entries sorted by filename, so files come back
// in filename order with no extra sort needed.
func curriculumFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		files = append(files, filepath.Join(dir, e.Name()))
	}
	return files, nil
}
