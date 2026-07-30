package mcqguess

import "testing"

func TestLongestOptionIndex(t *testing.T) {
	tests := []struct {
		name    string
		options []string
		want    int
	}{
		{
			name:    "unambiguous longest",
			options: []string{"short", "a bit longer", "the very longest option here"},
			want:    2,
		},
		{
			name:    "tie resolves to first (lowest) index",
			options: []string{"aaa", "bbb", "ccc"},
			want:    0,
		},
		{
			name: "tie between two, not all three, resolves to first tied index",
			// index 0 and 2 tie at 5 runes; index 1 is shorter.
			options: []string{"aaaaa", "bb", "ccccc"},
			want:    0,
		},
		{
			name: "rune count, not byte count: Thai text is ~3 bytes/char",
			// "สั้น" is 4 Thai runes / 12 bytes. "abcdefgh" is 8 runes / 8
			// bytes. By rune count "abcdefgh" (index 1) is longest; a
			// byte-based measurement would wrongly pick index 0 (12 > 8).
			options: []string{"สั้น", "abcdefgh"},
			want:    1,
		},
		{
			name:    "single option",
			options: []string{"only"},
			want:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LongestOptionIndex(tt.options); got != tt.want {
				t.Errorf("LongestOptionIndex(%v) = %d, want %d", tt.options, got, tt.want)
			}
		})
	}
}

func TestHeuristicResultRates(t *testing.T) {
	var h HeuristicResult
	h.add(true, 1.0/3.0)
	h.add(false, 1.0/3.0)
	h.add(false, 1.0/3.0)

	if got, want := h.HitRate(), 1.0/3.0; !almostEqual(got, want, 1e-9) {
		t.Errorf("HitRate() = %v, want %v", got, want)
	}
	if got, want := h.AvgBaseline(), 1.0/3.0; !almostEqual(got, want, 1e-9) {
		t.Errorf("AvgBaseline() = %v, want %v", got, want)
	}
	if got, want := h.ExcessRatio(), 0.0; !almostEqual(got, want, 1e-9) {
		t.Errorf("ExcessRatio() = %v, want %v", got, want)
	}
}

func TestHeuristicResultExcessRatioIsRelativeNotAbsolute(t *testing.T) {
	// 30% hit rate against a 25% baseline: (0.30-0.25)/0.25 = 0.20 (20%
	// relative), not 0.05 (5 percentage points) — the whole point of using
	// a ratio is that it stays comparable across corpora with different
	// baselines.
	var h HeuristicResult
	for i := 0; i < 30; i++ {
		h.add(true, 0.25)
	}
	for i := 0; i < 70; i++ {
		h.add(false, 0.25)
	}

	if got, want := h.ExcessRatio(), 0.20; !almostEqual(got, want, 1e-9) {
		t.Errorf("ExcessRatio() = %v, want %v", got, want)
	}
}

func TestHeuristicResultZeroTotal(t *testing.T) {
	var h HeuristicResult
	if got := h.HitRate(); got != 0 {
		t.Errorf("HitRate() on empty result = %v, want 0", got)
	}
	if got := h.AvgBaseline(); got != 0 {
		t.Errorf("AvgBaseline() on empty result = %v, want 0", got)
	}
	if got := h.ExcessRatio(); got != 0 {
		t.Errorf("ExcessRatio() on empty result = %v, want 0", got)
	}
}

// almostEqual is used only for assertions built from deterministic,
// exact-fraction fixtures (no sampling anywhere in this package's tests),
// so eps only needs to absorb float64 rounding, never statistical
// variance — see report_test.go's header comment for the full reasoning.
func almostEqual(a, b, eps float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= eps
}
