// Package app orchestrates the curriculum bounded context. It imports only
// domain and stdlib; infra depends on it, never the other way around.
package app

import (
	"context"
	"fmt"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

// Repository is the curriculum read port. infra provides the MySQL adapter;
// tests provide a fake.
type Repository interface {
	Topics(ctx context.Context) ([]domain.Topic, error)
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
