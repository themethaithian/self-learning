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
