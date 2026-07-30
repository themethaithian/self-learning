package app

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

type fakeRepository struct {
	topics       []domain.Topic
	availability map[ConceptPath]int
	err          error

	lesson    domain.Lesson
	lessonErr error
}

func (f fakeRepository) Topics(context.Context) ([]domain.Topic, map[ConceptPath]int, error) {
	return f.topics, f.availability, f.err
}

func (f fakeRepository) LessonByConcept(context.Context, string, string) (domain.Lesson, error) {
	return f.lesson, f.lessonErr
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

func mustLesson(t *testing.T, slug string) domain.Lesson {
	t.Helper()
	s, err := domain.NewSlug(slug)
	if err != nil {
		t.Fatalf("NewSlug: %v", err)
	}
	est, err := domain.NewEstMinutes(7)
	if err != nil {
		t.Fatalf("NewEstMinutes: %v", err)
	}
	refA, err := domain.NewReference("Source A", "https://example.com/a", "why a")
	if err != nil {
		t.Fatalf("NewReference: %v", err)
	}
	refB, err := domain.NewReference("Source B", "https://example.com/b", "why b")
	if err != nil {
		t.Fatalf("NewReference: %v", err)
	}
	kind, err := domain.NewRecallKind("short_answer")
	if err != nil {
		t.Fatalf("NewRecallKind: %v", err)
	}
	checks := make([]domain.RecallCheck, 3)
	for i := range checks {
		pos, err := domain.NewPosition(i + 1)
		if err != nil {
			t.Fatalf("NewPosition: %v", err)
		}
		checks[i], err = domain.NewRecallCheck(pos, kind, "q", "a", nil, "")
		if err != nil {
			t.Fatalf("NewRecallCheck: %v", err)
		}
	}

	l, err := domain.NewLesson(s, 1, "Title", est, "body", []domain.Reference{refA, refB}, checks)
	if err != nil {
		t.Fatalf("NewLesson: %v", err)
	}
	return l
}

func TestServiceLesson_Success(t *testing.T) {
	want := mustLesson(t, "concept-a")
	svc := NewService(fakeRepository{lesson: want})

	got, err := svc.Lesson(context.Background(), "topic-a", "concept-a")
	if err != nil {
		t.Fatalf("Lesson() unexpected error: %v", err)
	}
	if got.Slug() != want.Slug() || got.TitleEn() != want.TitleEn() {
		t.Fatalf("Lesson() = %v, want %v", got, want)
	}
}

func TestServiceLesson_NotFound(t *testing.T) {
	svc := NewService(fakeRepository{lessonErr: fmt.Errorf("infra: lesson not found: %w", ErrLessonNotFound)})

	_, err := svc.Lesson(context.Background(), "topic-a", "concept-a")
	if !errors.Is(err, ErrLessonNotFound) {
		t.Fatalf("Lesson() error = %v, want it to wrap ErrLessonNotFound", err)
	}
}

func TestServiceLesson_RepositoryError(t *testing.T) {
	repoErr := errors.New("connection lost")
	svc := NewService(fakeRepository{lessonErr: repoErr})

	_, err := svc.Lesson(context.Background(), "topic-a", "concept-a")
	if err == nil {
		t.Fatal("Lesson() expected error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("Lesson() error = %v, want it to wrap %v", err, repoErr)
	}
}

func TestServiceTree_Success(t *testing.T) {
	want := []domain.Topic{mustTopic(t, "topic-a")}
	path := ConceptPath{TopicSlug: "topic-a", ConceptSlug: "concept-a"}
	wantAvailability := map[ConceptPath]int{path: 7}
	svc := NewService(fakeRepository{topics: want, availability: wantAvailability})

	got, availability, err := svc.Tree(context.Background())
	if err != nil {
		t.Fatalf("Tree() unexpected error: %v", err)
	}
	if len(got) != len(want) || got[0].Slug().String() != want[0].Slug().String() {
		t.Fatalf("Tree() = %v, want %v", got, want)
	}
	if availability[path] != 7 {
		t.Fatalf("Tree() availability = %v, want %v", availability, wantAvailability)
	}
}

func TestServiceTree_RepositoryError(t *testing.T) {
	repoErr := errors.New("connection lost")
	svc := NewService(fakeRepository{err: repoErr})

	_, _, err := svc.Tree(context.Background())
	if err == nil {
		t.Fatal("Tree() expected error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("Tree() error = %v, want it to wrap %v", err, repoErr)
	}
}
