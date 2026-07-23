package app

import (
	"context"
	"fmt"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

// Loader is the curriculum content-file read port. infra provides the file
// adapter; tests provide a fake.
type Loader interface {
	LoadTopic(path string) (domain.Topic, error)
}

// Writer is the curriculum write port. infra provides the MySQL adapter;
// tests provide a fake.
type Writer interface {
	SaveTopic(ctx context.Context, t domain.Topic) error
}

// ImportService owns the import use case: load one content file, then write
// the resulting topic. cmd/import-curriculum wires adapters and iterates
// files; it does not otherwise touch domain or infra.
type ImportService struct {
	loader Loader
	writer Writer
}

func NewImportService(loader Loader, writer Writer) ImportService {
	return ImportService{loader: loader, writer: writer}
}

// ImportFile loads the topic at path and saves it, returning the topic so
// the caller can report what was imported (slug, chapter/concept counts).
func (s ImportService) ImportFile(ctx context.Context, path string) (domain.Topic, error) {
	topic, err := s.loader.LoadTopic(path)
	if err != nil {
		return domain.Topic{}, fmt.Errorf("curriculum: import %s: %w", path, err)
	}
	if err := s.writer.SaveTopic(ctx, topic); err != nil {
		return domain.Topic{}, fmt.Errorf("curriculum: import %s: %w", path, err)
	}
	return topic, nil
}
