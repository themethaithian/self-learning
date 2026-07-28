// Package app orchestrates the prefs bounded context. It imports only domain
// and stdlib; infra depends on it, never the other way around.
package app

import (
	"context"
	"fmt"

	"github.com/themethaithian/self-learning/internal/prefs/domain"
)

const focusTrackPrefName = "focus_track"

// Repository is the prefs read/write port over the generic name/value
// table. infra provides the MySQL adapter; tests provide a fake.
type Repository interface {
	Get(ctx context.Context, name string) (value string, ok bool, err error)
	Set(ctx context.Context, name, value string) error
	Delete(ctx context.Context, name string) error
}

// Service is the prefs use-case layer.
type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

// FocusTrack returns the currently set focus track slug. ok is false when no
// preference has been set (or it was cleared).
func (s Service) FocusTrack(ctx context.Context) (value string, ok bool, err error) {
	value, ok, err = s.repo.Get(ctx, focusTrackPrefName)
	if err != nil {
		return "", false, fmt.Errorf("prefs: get focus track: %w", err)
	}
	return value, ok, nil
}

// SetFocusTrack sets the focus track preference to raw, or clears it when
// raw is nil.
func (s Service) SetFocusTrack(ctx context.Context, raw *string) error {
	if raw == nil {
		if err := s.repo.Delete(ctx, focusTrackPrefName); err != nil {
			return fmt.Errorf("prefs: clear focus track: %w", err)
		}
		return nil
	}

	// Format only: whether this track exists in the curriculum is that
	// bounded context's concern — prefs must not import internal/curriculum.
	track, err := domain.NewFocusTrack(*raw)
	if err != nil {
		return fmt.Errorf("prefs: set focus track: %w", err)
	}
	if err := s.repo.Set(ctx, focusTrackPrefName, track.String()); err != nil {
		return fmt.Errorf("prefs: set focus track: %w", err)
	}
	return nil
}
