package domain

import (
	"errors"
	"testing"
)

func TestNewChunkState(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{name: "locked", raw: "locked"},
		{name: "in_progress", raw: "in_progress"},
		{name: "passed", raw: "passed"},
		{name: "empty", raw: "", wantErr: ErrInvalidChunkState},
		{name: "wrong case", raw: "Locked", wantErr: ErrInvalidChunkState},
		{name: "unknown state", raw: "skipped", wantErr: ErrInvalidChunkState},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChunkState(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewChunkState(%q) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewChunkState(%q) unexpected error: %v", tt.raw, err)
			}
			if got.String() != tt.raw {
				t.Errorf("String() = %q, want %q", got.String(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid state %q", tt.raw)
			}

			wantLocked := tt.raw == "locked"
			wantInProgress := tt.raw == "in_progress"
			wantPassed := tt.raw == "passed"
			if got.IsLocked() != wantLocked {
				t.Errorf("IsLocked() = %v, want %v", got.IsLocked(), wantLocked)
			}
			if got.IsInProgress() != wantInProgress {
				t.Errorf("IsInProgress() = %v, want %v", got.IsInProgress(), wantInProgress)
			}
			if got.IsPassed() != wantPassed {
				t.Errorf("IsPassed() = %v, want %v", got.IsPassed(), wantPassed)
			}
		})
	}
}

func TestChunkStateZeroValue(t *testing.T) {
	var s ChunkState
	if !s.IsZero() {
		t.Errorf("zero-value ChunkState.IsZero() = false, want true")
	}
}
