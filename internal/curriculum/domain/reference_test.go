package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewReference(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		source  string
		why     string
		wantErr error
	}{
		{name: "valid", title: "Evans ch. 4", source: "Domain-Driven Design, chapter 4", why: "core building blocks"},
		{
			name: "fields with surrounding whitespace trimmed", title: "  Evans ch. 4  ",
			source: "  Domain-Driven Design, chapter 4  ", why: "  core building blocks  ",
		},
		{name: "empty title", title: "", source: "Domain-Driven Design, chapter 4", why: "core building blocks", wantErr: ErrInvalidReference},
		{name: "whitespace-only title", title: "   ", source: "Domain-Driven Design, chapter 4", why: "core building blocks", wantErr: ErrInvalidReference},
		{name: "empty source", title: "Evans ch. 4", source: "", why: "core building blocks", wantErr: ErrInvalidReference},
		{name: "empty why", title: "Evans ch. 4", source: "Domain-Driven Design, chapter 4", why: "", wantErr: ErrInvalidReference},
		{
			name: "title too long", title: strings.Repeat("a", maxReferenceTitleRunes+1),
			source: "Domain-Driven Design, chapter 4", why: "core building blocks", wantErr: ErrInvalidReference,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReference(tt.title, tt.source, tt.why)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewReference() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewReference() unexpected error: %v", err)
			}
			if want := strings.TrimSpace(tt.title); got.Title() != want {
				t.Errorf("Title() = %q, want %q", got.Title(), want)
			}
			if want := strings.TrimSpace(tt.source); got.Source() != want {
				t.Errorf("Source() = %q, want %q", got.Source(), want)
			}
			if want := strings.TrimSpace(tt.why); got.Why() != want {
				t.Errorf("Why() = %q, want %q", got.Why(), want)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for a validly constructed Reference")
			}
		})
	}
}

func TestReferenceZeroValue(t *testing.T) {
	var r Reference
	if !r.IsZero() {
		t.Errorf("zero-value Reference.IsZero() = false, want true")
	}
}
