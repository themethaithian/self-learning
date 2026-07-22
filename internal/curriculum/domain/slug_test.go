package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewSlug(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{name: "valid single char", raw: "a"},
		{name: "valid with hyphen and digits", raw: "layered-architecture-101"},
		{name: "valid exactly 100 chars", raw: strings.Repeat("a", 100)},
		{name: "empty", raw: "", wantErr: ErrInvalidSlug},
		{name: "uppercase", raw: "DDD", wantErr: ErrInvalidSlug},
		{name: "leading hyphen", raw: "-ddd", wantErr: ErrInvalidSlug},
		{name: "trailing hyphen", raw: "ddd-", wantErr: ErrInvalidSlug},
		{name: "double hyphen", raw: "dd--d", wantErr: ErrInvalidSlug},
		{name: "101 chars", raw: strings.Repeat("a", 101), wantErr: ErrInvalidSlug},
		{name: "underscore not allowed", raw: "d_d", wantErr: ErrInvalidSlug},
		{name: "space not allowed", raw: "d d", wantErr: ErrInvalidSlug},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSlug(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewSlug(%q) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewSlug(%q) unexpected error: %v", tt.raw, err)
			}
			if got.String() != tt.raw {
				t.Errorf("String() = %q, want %q", got.String(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid slug %q", tt.raw)
			}
		})
	}
}

func TestSlugZeroValue(t *testing.T) {
	var s Slug
	if !s.IsZero() {
		t.Errorf("zero-value Slug.IsZero() = false, want true")
	}
}
