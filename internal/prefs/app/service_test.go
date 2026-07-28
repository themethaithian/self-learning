package app

import (
	"context"
	"errors"
	"testing"

	"github.com/themethaithian/self-learning/internal/prefs/domain"
)

// TestFocusTrackPrefNameLiteral pins the literal MySQL storage key: Get and
// Set both derive from the same constant, so every round-trip test above
// stays green even if the constant is renamed — but any row already written
// under the old literal becomes permanently unreachable in production, with
// no error and no test failure. This is the only thing that guards the
// literal itself, the same way TestUpsertPrefSQLShape guards the upsert SQL.
func TestFocusTrackPrefNameLiteral(t *testing.T) {
	if focusTrackPrefName != "focus_track" {
		t.Fatalf("focusTrackPrefName = %q, want %q", focusTrackPrefName, "focus_track")
	}
}

type fakeRepository struct {
	values  map[string]string
	getErr  error
	setErr  error
	delErr  error
	setCall int
	delCall int
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{values: make(map[string]string)}
}

func (f *fakeRepository) Get(_ context.Context, name string) (string, bool, error) {
	if f.getErr != nil {
		return "", false, f.getErr
	}
	v, ok := f.values[name]
	return v, ok, nil
}

func (f *fakeRepository) Set(_ context.Context, name, value string) error {
	f.setCall++
	if f.setErr != nil {
		return f.setErr
	}
	f.values[name] = value
	return nil
}

func (f *fakeRepository) Delete(_ context.Context, name string) error {
	f.delCall++
	if f.delErr != nil {
		return f.delErr
	}
	delete(f.values, name)
	return nil
}

func strPtr(s string) *string { return &s }

func TestServiceFocusTrack_Unset(t *testing.T) {
	svc := NewService(newFakeRepository())

	value, ok, err := svc.FocusTrack(context.Background())
	if err != nil {
		t.Fatalf("FocusTrack() unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("FocusTrack() ok = true, want false for an unset pref; value=%q", value)
	}
}

func TestServiceSetFocusTrack_ThenGet(t *testing.T) {
	svc := NewService(newFakeRepository())

	if err := svc.SetFocusTrack(context.Background(), strPtr("ai-systems")); err != nil {
		t.Fatalf("SetFocusTrack() unexpected error: %v", err)
	}
	value, ok, err := svc.FocusTrack(context.Background())
	if err != nil || !ok || value != "ai-systems" {
		t.Fatalf("FocusTrack() = (%q, %v, %v), want (ai-systems, true, nil)", value, ok, err)
	}
}

// TestServiceSetFocusTrack_Idempotent pins the ticket's core requirement:
// PUTting the same value twice must succeed both times with no error.
func TestServiceSetFocusTrack_Idempotent(t *testing.T) {
	repo := newFakeRepository()
	svc := NewService(repo)

	for i := 0; i < 2; i++ {
		if err := svc.SetFocusTrack(context.Background(), strPtr("ai-systems")); err != nil {
			t.Fatalf("SetFocusTrack() call %d unexpected error: %v", i, err)
		}
	}
	if repo.setCall != 2 {
		t.Fatalf("repo.Set called %d times, want 2", repo.setCall)
	}
	value, ok, err := svc.FocusTrack(context.Background())
	if err != nil || !ok || value != "ai-systems" {
		t.Fatalf("FocusTrack() = (%q, %v, %v), want (ai-systems, true, nil)", value, ok, err)
	}
}

func TestServiceSetFocusTrack_ClearDeletesRow(t *testing.T) {
	svc := NewService(newFakeRepository())
	if err := svc.SetFocusTrack(context.Background(), strPtr("ai-systems")); err != nil {
		t.Fatalf("SetFocusTrack() unexpected error: %v", err)
	}

	if err := svc.SetFocusTrack(context.Background(), nil); err != nil {
		t.Fatalf("SetFocusTrack(nil) unexpected error: %v", err)
	}

	_, ok, err := svc.FocusTrack(context.Background())
	if err != nil {
		t.Fatalf("FocusTrack() unexpected error: %v", err)
	}
	if ok {
		t.Fatal("FocusTrack() ok = true after clear, want false")
	}
}

// TestServiceSetFocusTrack_ClearIsIdempotent: clearing an already-unset pref
// (delete of an absent row) must not error either.
func TestServiceSetFocusTrack_ClearIsIdempotent(t *testing.T) {
	svc := NewService(newFakeRepository())

	for i := 0; i < 2; i++ {
		if err := svc.SetFocusTrack(context.Background(), nil); err != nil {
			t.Fatalf("SetFocusTrack(nil) call %d unexpected error: %v", i, err)
		}
	}
}

func TestServiceSetFocusTrack_InvalidFormat(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{name: "empty", raw: ""},
		{name: "uppercase", raw: "AI-Systems"},
		{name: "space", raw: "ai systems"},
		{name: "exclamation", raw: "ai-systems!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepository()
			svc := NewService(repo)

			err := svc.SetFocusTrack(context.Background(), strPtr(tt.raw))
			if !errors.Is(err, domain.ErrInvalidFocusTrack) {
				t.Fatalf("SetFocusTrack(%q) error = %v, want it to wrap ErrInvalidFocusTrack", tt.raw, err)
			}
			if repo.setCall != 0 {
				t.Fatalf("SetFocusTrack(%q) called repo.Set despite a validation failure", tt.raw)
			}
		})
	}
}

func TestServiceFocusTrack_RepositoryError(t *testing.T) {
	repoErr := errors.New("connection lost")
	repo := newFakeRepository()
	repo.getErr = repoErr
	svc := NewService(repo)

	_, _, err := svc.FocusTrack(context.Background())
	if !errors.Is(err, repoErr) {
		t.Fatalf("FocusTrack() error = %v, want it to wrap %v", err, repoErr)
	}
}

func TestServiceSetFocusTrack_RepositoryError(t *testing.T) {
	repoErr := errors.New("connection lost")
	repo := newFakeRepository()
	repo.setErr = repoErr
	svc := NewService(repo)

	err := svc.SetFocusTrack(context.Background(), strPtr("ai-systems"))
	if !errors.Is(err, repoErr) {
		t.Fatalf("SetFocusTrack() error = %v, want it to wrap %v", err, repoErr)
	}
}
