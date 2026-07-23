package domain

import (
	"errors"
	"testing"
)

func TestNewLessonProgress(t *testing.T) {
	tests := []struct {
		name      string
		lessonRef LessonRef
		state     ChunkState
		wantErr   error
	}{
		{name: "valid locked", lessonRef: mustLessonRef(t, "aggregate"), state: mustChunkState(t, "locked")},
		{name: "valid in_progress", lessonRef: mustLessonRef(t, "aggregate"), state: mustChunkState(t, "in_progress")},
		{name: "valid passed", lessonRef: mustLessonRef(t, "aggregate"), state: mustChunkState(t, "passed")},
		{name: "zero-value lesson ref", lessonRef: LessonRef{}, state: mustChunkState(t, "locked"), wantErr: ErrInvalidLessonRef},
		{name: "zero-value state", lessonRef: mustLessonRef(t, "aggregate"), state: ChunkState{}, wantErr: ErrInvalidChunkState},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLessonProgress(tt.lessonRef, tt.state)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewLessonProgress() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewLessonProgress() unexpected error: %v", err)
			}
			if got.LessonRef() != tt.lessonRef {
				t.Errorf("LessonRef() = %v, want %v", got.LessonRef(), tt.lessonRef)
			}
			if got.State() != tt.state {
				t.Errorf("State() = %v, want %v", got.State(), tt.state)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for a validly constructed LessonProgress")
			}
		})
	}
}

func TestLessonProgressZeroValue(t *testing.T) {
	var p LessonProgress
	if !p.IsZero() {
		t.Errorf("zero-value LessonProgress.IsZero() = false, want true")
	}
}

func TestLessonProgressUnlock(t *testing.T) {
	tests := []struct {
		name    string
		state   string
		wantErr error
	}{
		{name: "from locked succeeds", state: "locked"},
		{name: "from in_progress rejected", state: "in_progress", wantErr: ErrAlreadyUnlocked},
		{name: "from passed rejected", state: "passed", wantErr: ErrAlreadyUnlocked},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := mustLessonProgress(t, "aggregate", tt.state)
			got, err := p.Unlock()
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Unlock() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unlock() unexpected error: %v", err)
			}
			if !got.State().IsInProgress() {
				t.Errorf("State() = %v, want in_progress", got.State())
			}
			if !p.State().IsLocked() {
				t.Errorf("receiver mutated: State() = %v, want still locked", p.State())
			}
		})
	}
}

func TestLessonProgressMarkPassed(t *testing.T) {
	tests := []struct {
		name    string
		state   string
		wantErr error
	}{
		{name: "from locked rejected", state: "locked", wantErr: ErrLessonLocked},
		{name: "from in_progress succeeds", state: "in_progress"},
		{name: "from passed rejected", state: "passed", wantErr: ErrAlreadyPassed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := mustLessonProgress(t, "aggregate", tt.state)
			got, err := p.MarkPassed()
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("MarkPassed() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("MarkPassed() unexpected error: %v", err)
			}
			if !got.State().IsPassed() {
				t.Errorf("State() = %v, want passed", got.State())
			}
			if !p.State().IsInProgress() {
				t.Errorf("receiver mutated: State() = %v, want still in_progress", p.State())
			}
		})
	}
}
