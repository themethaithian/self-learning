package mcqguess

import "testing"

func TestQuestionBaseline(t *testing.T) {
	tests := []struct {
		name    string
		options []string
		want    float64
	}{
		{name: "three options", options: []string{"a", "b", "c"}, want: 1.0 / 3.0},
		{name: "four options", options: []string{"a", "b", "c", "d"}, want: 0.25},
		{name: "two options", options: []string{"a", "b"}, want: 0.5},
		{name: "zero options guards against divide-by-zero", options: nil, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Question{Options: tt.options}
			if got := q.Baseline(); got != tt.want {
				t.Errorf("Baseline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionExpectedIndex(t *testing.T) {
	tests := []struct {
		name           string
		options        []string
		expectedAnswer string
		wantIndex      int
		wantOK         bool
	}{
		{name: "found at start", options: []string{"a", "b", "c"}, expectedAnswer: "a", wantIndex: 0, wantOK: true},
		{name: "found in middle", options: []string{"a", "b", "c"}, expectedAnswer: "b", wantIndex: 1, wantOK: true},
		{name: "not among options", options: []string{"a", "b", "c"}, expectedAnswer: "z", wantIndex: -1, wantOK: false},
		{name: "no options at all", options: nil, expectedAnswer: "a", wantIndex: -1, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Question{Options: tt.options, ExpectedAnswer: tt.expectedAnswer}
			gotIndex, gotOK := q.ExpectedIndex()
			if gotIndex != tt.wantIndex || gotOK != tt.wantOK {
				t.Errorf("ExpectedIndex() = (%d, %v), want (%d, %v)", gotIndex, gotOK, tt.wantIndex, tt.wantOK)
			}
		})
	}
}
