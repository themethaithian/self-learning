package domain

import (
	"errors"
	"testing"
)

func TestNewTrack(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{name: "ddd", raw: "ddd"},
		{name: "distsys", raw: "distsys"},
		{name: "aws", raw: "aws"},
		{name: "go", raw: "go"},
		{name: "dsa", raw: "dsa"},
		{name: "empty", raw: "", wantErr: ErrInvalidTrack},
		{name: "wrong case", raw: "DDD", wantErr: ErrInvalidTrack},
		{name: "unknown track", raw: "kubernetes", wantErr: ErrInvalidTrack},
		{name: "not trimmed", raw: " ddd", wantErr: ErrInvalidTrack},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTrack(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewTrack(%q) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewTrack(%q) unexpected error: %v", tt.raw, err)
			}
			if got.String() != tt.raw {
				t.Errorf("String() = %q, want %q", got.String(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid track %q", tt.raw)
			}
		})
	}
}

func TestTrackZeroValue(t *testing.T) {
	var tr Track
	if !tr.IsZero() {
		t.Errorf("zero-value Track.IsZero() = false, want true")
	}
}

func TestTracks(t *testing.T) {
	got := Tracks()
	want := []string{"ddd", "distsys", "aws", "go", "dsa"}
	if len(got) != len(want) {
		t.Fatalf("Tracks() returned %d tracks, want %d", len(got), len(want))
	}

	seen := make(map[string]bool, len(got))
	for i, tr := range got {
		if tr.IsZero() {
			t.Fatalf("Tracks()[%d] is zero-value", i)
		}
		if tr.String() != want[i] {
			t.Fatalf("Tracks()[%d] = %q, want %q", i, tr.String(), want[i])
		}
		seen[tr.String()] = true
	}
	if len(seen) != 5 {
		t.Fatalf("Tracks() contains %d distinct tracks, want exactly 5", len(seen))
	}
}
