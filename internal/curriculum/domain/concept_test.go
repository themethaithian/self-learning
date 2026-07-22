package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewConcept(t *testing.T) {
	const thaiChar = "ก"

	tests := []struct {
		name     string
		slug     Slug
		title    string
		outline  string
		position Position
		wantErr  error
	}{
		{
			name: "valid", slug: mustSlug(t, "aggregate"), title: "Aggregate",
			outline: "Intro to aggregates", position: mustPosition(t, 1),
		},
		{
			name: "title with surrounding whitespace trimmed", slug: mustSlug(t, "aggregate"),
			title: "  Aggregate  ", outline: "Outline text", position: mustPosition(t, 1),
		},
		{
			name: "title exactly 255 runes", slug: mustSlug(t, "aggregate"),
			title: strings.Repeat("a", 255), outline: "Outline text", position: mustPosition(t, 1),
		},
		{
			name: "title exactly 255 Thai runes", slug: mustSlug(t, "aggregate"),
			title: strings.Repeat(thaiChar, 255), outline: "Outline text", position: mustPosition(t, 1),
		},
		{
			name: "outline exactly 4000 runes", slug: mustSlug(t, "aggregate"),
			title: "Aggregate", outline: strings.Repeat("a", 4000), position: mustPosition(t, 1),
		},
		{
			name: "outline exactly 4000 Thai runes", slug: mustSlug(t, "aggregate"),
			title: "Aggregate", outline: strings.Repeat(thaiChar, 4000), position: mustPosition(t, 1),
		},
		{
			name: "zero-value slug", slug: Slug{}, title: "Aggregate",
			outline: "Outline text", position: mustPosition(t, 1), wantErr: ErrInvalidSlug,
		},
		{
			name: "whitespace-only title", slug: mustSlug(t, "aggregate"), title: "   ",
			outline: "Outline text", position: mustPosition(t, 1), wantErr: ErrInvalidTitle,
		},
		{
			name: "256-rune title", slug: mustSlug(t, "aggregate"), title: strings.Repeat("a", 256),
			outline: "Outline text", position: mustPosition(t, 1), wantErr: ErrInvalidTitle,
		},
		{
			name: "256 Thai-rune title", slug: mustSlug(t, "aggregate"), title: strings.Repeat(thaiChar, 256),
			outline: "Outline text", position: mustPosition(t, 1), wantErr: ErrInvalidTitle,
		},
		{
			name: "whitespace-only outline", slug: mustSlug(t, "aggregate"), title: "Aggregate",
			outline: "   ", position: mustPosition(t, 1), wantErr: ErrInvalidOutline,
		},
		{
			name: "4001-rune outline", slug: mustSlug(t, "aggregate"), title: "Aggregate",
			outline: strings.Repeat("a", 4001), position: mustPosition(t, 1), wantErr: ErrInvalidOutline,
		},
		{
			name: "4001 Thai-rune outline", slug: mustSlug(t, "aggregate"), title: "Aggregate",
			outline: strings.Repeat(thaiChar, 4001), position: mustPosition(t, 1), wantErr: ErrInvalidOutline,
		},
		{
			name: "zero-value position", slug: mustSlug(t, "aggregate"), title: "Aggregate",
			outline: "Outline text", position: Position{}, wantErr: ErrInvalidPosition,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConcept(tt.slug, tt.title, tt.outline, tt.position)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewConcept() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewConcept() unexpected error: %v", err)
			}
			if got.Slug() != tt.slug {
				t.Errorf("Slug() = %v, want %v", got.Slug(), tt.slug)
			}
			if want := strings.TrimSpace(tt.title); got.Title() != want {
				t.Errorf("Title() = %q, want %q", got.Title(), want)
			}
			if want := strings.TrimSpace(tt.outline); got.Outline() != want {
				t.Errorf("Outline() = %q, want %q", got.Outline(), want)
			}
			if got.Position() != tt.position {
				t.Errorf("Position() = %v, want %v", got.Position(), tt.position)
			}
		})
	}
}

func TestConceptZeroValue(t *testing.T) {
	var c Concept
	if !c.IsZero() {
		t.Errorf("zero-value Concept.IsZero() = false, want true")
	}
}
