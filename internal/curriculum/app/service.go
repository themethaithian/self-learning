// Package app orchestrates the curriculum bounded context. It imports only
// domain and stdlib; infra depends on it, never the other way around.
package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

// ErrLessonNotFound is returned by Service.Lesson when a concept exists but
// has no lesson content yet — distinct from any other repository failure so
// the HTTP handler can answer 404 instead of 500.
var ErrLessonNotFound = errors.New("curriculum: lesson not found")

// Repository is the curriculum read port. infra provides the MySQL adapter;
// tests provide a fake.
type Repository interface {
	Topics(ctx context.Context) ([]domain.Topic, error)
	LessonByConcept(ctx context.Context, topicSlug, conceptSlug string) (domain.Lesson, error)
}

// Service is the curriculum use-case layer.
type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Tree(ctx context.Context) ([]domain.Topic, error) {
	topics, err := s.repo.Topics(ctx)
	if err != nil {
		return nil, fmt.Errorf("curriculum: tree: %w", err)
	}
	return topics, nil
}

// Lesson returns the full lesson for a concept, identified by its topic and
// concept slugs. Returns ErrLessonNotFound if the concept has no lesson yet.
func (s Service) Lesson(ctx context.Context, topicSlug, conceptSlug string) (domain.Lesson, error) {
	lesson, err := s.repo.LessonByConcept(ctx, topicSlug, conceptSlug)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("curriculum: lesson %s/%s: %w", topicSlug, conceptSlug, err)
	}
	return lesson, nil
}
