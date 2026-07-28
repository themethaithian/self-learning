package domain

import (
	"errors"
	"strings"
	"testing"
)

// TestMaxFocusTrackLenFitsValueColumn guards against raising the constant
// past what prefs.value (VARCHAR(255)) can hold under MySQL strict mode.
func TestMaxFocusTrackLenFitsValueColumn(t *testing.T) {
	const prefsValueColumnWidth = 255
	if maxFocusTrackLen > prefsValueColumnWidth {
		t.Fatalf("maxFocusTrackLen = %d, must not exceed the prefs.value column width %d", maxFocusTrackLen, prefsValueColumnWidth)
	}
}

func TestNewFocusTrack(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{name: "valid lowercase with hyphen", raw: "ai-systems", wantErr: false},
		{name: "valid single char", raw: "a", wantErr: false},
		{name: "valid digits and hyphens", raw: "go-101", wantErr: false},
		{name: "valid at max length", raw: strings.Repeat("a", maxFocusTrackLen), wantErr: false},
		{name: "empty", raw: "", wantErr: true},
		{name: "over max length", raw: strings.Repeat("a", maxFocusTrackLen+1), wantErr: true},
		{name: "uppercase", raw: "AI-Systems", wantErr: true},
		{name: "space", raw: "ai systems", wantErr: true},
		{name: "punctuation", raw: "ai-systems!", wantErr: true},
		{name: "leading hyphen", raw: "-ai-systems", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track, err := NewFocusTrack(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("NewFocusTrack(%q) = nil error, want error", tt.raw)
				}
				if !errors.Is(err, ErrInvalidFocusTrack) {
					t.Fatalf("NewFocusTrack(%q) error = %v, want it to wrap ErrInvalidFocusTrack", tt.raw, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewFocusTrack(%q) unexpected error: %v", tt.raw, err)
			}
			if track.String() != tt.raw {
				t.Fatalf("String() = %q, want %q", track.String(), tt.raw)
			}
		})
	}
}
