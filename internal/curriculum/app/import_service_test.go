package app

import (
	"context"
	"errors"
	"testing"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

type fakeLoader struct {
	topics map[string]domain.Topic
	err    error
}

func (f *fakeLoader) LoadTopic(path string) (domain.Topic, error) {
	if f.err != nil {
		return domain.Topic{}, f.err
	}
	return f.topics[path], nil
}

type fakeWriter struct {
	saved []domain.Topic
	err   error
}

func (f *fakeWriter) SaveTopic(_ context.Context, t domain.Topic) error {
	f.saved = append(f.saved, t)
	return f.err
}

func TestImportServiceImportFile_Success(t *testing.T) {
	topic := mustTopic(t, "topic-a")
	loader := &fakeLoader{topics: map[string]domain.Topic{"ddd.json": topic}}
	writer := &fakeWriter{}
	svc := NewImportService(loader, writer)

	got, err := svc.ImportFile(context.Background(), "ddd.json")
	if err != nil {
		t.Fatalf("ImportFile() unexpected error: %v", err)
	}
	if got.Slug().String() != "topic-a" {
		t.Fatalf("ImportFile() returned topic %q, want topic-a", got.Slug().String())
	}
	if len(writer.saved) != 1 || writer.saved[0].Slug().String() != "topic-a" {
		t.Fatalf("ImportFile() did not pass the loaded topic to Writer.SaveTopic: saved = %v", writer.saved)
	}
}

func TestImportServiceImportFile_LoaderError(t *testing.T) {
	loaderErr := errors.New("malformed file")
	loader := &fakeLoader{err: loaderErr}
	writer := &fakeWriter{}
	svc := NewImportService(loader, writer)

	_, err := svc.ImportFile(context.Background(), "broken.json")
	if err == nil {
		t.Fatal("ImportFile() expected error, got nil")
	}
	if !errors.Is(err, loaderErr) {
		t.Fatalf("ImportFile() error = %v, want it to wrap %v", err, loaderErr)
	}
	if len(writer.saved) != 0 {
		t.Fatalf("ImportFile() called Writer.SaveTopic after a loader error: saved = %v", writer.saved)
	}
}

func TestImportServiceImportFile_WriterError(t *testing.T) {
	topic := mustTopic(t, "topic-a")
	loader := &fakeLoader{topics: map[string]domain.Topic{"ddd.json": topic}}
	writerErr := errors.New("connection lost")
	writer := &fakeWriter{err: writerErr}
	svc := NewImportService(loader, writer)

	_, err := svc.ImportFile(context.Background(), "ddd.json")
	if err == nil {
		t.Fatal("ImportFile() expected error, got nil")
	}
	if !errors.Is(err, writerErr) {
		t.Fatalf("ImportFile() error = %v, want it to wrap %v", err, writerErr)
	}
}
