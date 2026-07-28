package domain

import (
	"strings"
	"testing"
)

func TestIsValidSlugShape(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want bool
	}{
		{name: "simple", raw: "aggregate", want: true},
		{name: "with hyphen", raw: "aggregate-root", want: true},
		{name: "with digits", raw: "cap-theorem-2", want: true},
		{name: "max length", raw: strings.Repeat("a", maxSlugShapeLen), want: true},
		{name: "empty", raw: "", want: false},
		{name: "uppercase", raw: "Aggregate", want: false},
		{name: "space", raw: "b trees", want: false},
		{name: "leading hyphen", raw: "-aggregate", want: false},
		{name: "trailing hyphen", raw: "aggregate-", want: false},
		{name: "consecutive hyphens", raw: "aggregate--root", want: false},
		{name: "too long", raw: strings.Repeat("a", maxSlugShapeLen+1), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSlugShape(tt.raw); got != tt.want {
				t.Errorf("IsValidSlugShape(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}
