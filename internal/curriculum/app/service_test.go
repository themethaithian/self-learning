package app

import (
	"context"
	"errors"
	"testing"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

type fakeRepository struct {
	topics []domain.Topic
	err    error
}

func (f fakeRepository) Topics(context.Context) ([]domain.Topic, error) {
	return f.topics, f.err
}

func mustTopic(t *testing.T, slug string) domain.Topic {
	t.Helper()
	track, err := domain.NewTrack("go")
	if err != nil {
		t.Fatalf("NewTrack: %v", err)
	}
	s, err := domain.NewSlug(slug)
	if err != nil {
		t.Fatalf("NewSlug: %v", err)
	}
	pos, err := domain.NewPosition(1)
	if err != nil {
		t.Fatalf("NewPosition: %v", err)
	}
	conceptSlug, err := domain.NewSlug("concept-a")
	if err != nil {
		t.Fatalf("NewSlug: %v", err)
	}
	concept, err := domain.NewConcept(conceptSlug, "Concept A", "outline", pos)
	if err != nil {
		t.Fatalf("NewConcept: %v", err)
	}
	chapterSlug, err := domain.NewSlug("chapter-a")
	if err != nil {
		t.Fatalf("NewSlug: %v", err)
	}
	chapter, err := domain.NewChapter(chapterSlug, "Chapter A", pos, []domain.Concept{concept})
	if err != nil {
		t.Fatalf("NewChapter: %v", err)
	}
	topic, err := domain.NewTopic(track, s, "Topic", pos, []domain.Chapter{chapter})
	if err != nil {
		t.Fatalf("NewTopic: %v", err)
	}
	return topic
}

func TestServiceTree_Success(t *testing.T) {
	want := []domain.Topic{mustTopic(t, "topic-a")}
	svc := NewService(fakeRepository{topics: want})

	got, err := svc.Tree(context.Background())
	if err != nil {
		t.Fatalf("Tree() unexpected error: %v", err)
	}
	if len(got) != len(want) || got[0].Slug().String() != want[0].Slug().String() {
		t.Fatalf("Tree() = %v, want %v", got, want)
	}
}

func TestServiceTree_RepositoryError(t *testing.T) {
	repoErr := errors.New("connection lost")
	svc := NewService(fakeRepository{err: repoErr})

	_, err := svc.Tree(context.Background())
	if err == nil {
		t.Fatal("Tree() expected error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("Tree() error = %v, want it to wrap %v", err, repoErr)
	}
}
