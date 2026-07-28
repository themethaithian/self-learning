package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewLessonRef(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{name: "simple", raw: "aggregate"},
		{name: "with hyphen", raw: "aggregate-root"},
		{name: "with digits", raw: "cap-theorem-2"},
		{name: "empty", raw: "", wantErr: ErrInvalidLessonRef},
		{name: "uppercase", raw: "Aggregate", wantErr: ErrInvalidLessonRef},
		{name: "leading hyphen", raw: "-aggregate", wantErr: ErrInvalidLessonRef},
		{name: "trailing hyphen", raw: "aggregate-", wantErr: ErrInvalidLessonRef},
		{name: "consecutive hyphens", raw: "aggregate--root", wantErr: ErrInvalidLessonRef},
		{name: "too long", raw: strings.Repeat("a", maxSlugShapeLen+1), wantErr: ErrInvalidLessonRef},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLessonRef(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewLessonRef(%q) error = %v, want wrapping %v", tt.raw, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewLessonRef(%q) unexpected error: %v", tt.raw, err)
			}
			if got.String() != tt.raw {
				t.Errorf("String() = %q, want %q", got.String(), tt.raw)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for valid ref %q", tt.raw)
			}
		})
	}
}

func TestLessonRefZeroValue(t *testing.T) {
	var r LessonRef
	if !r.IsZero() {
		t.Errorf("zero-value LessonRef.IsZero() = false, want true")
	}
}
