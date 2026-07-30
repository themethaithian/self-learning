package mcqguess

import "unicode/utf8"

// LongestOptionIndex returns the index of the option with the most Unicode
// code points (runes), not bytes: Thai text runs about 3 bytes per
// character, so a byte-based length would rank a short Thai option as
// "longer" than a longer English one, which is a different and wrong
// metric. Ties resolve to the first (lowest-index) option: a real
// "always pick the longest-looking one" test-taker scans top to bottom and
// commits to the first option that looks longest, and a deterministic rule
// keeps this tool's output reproducible run to run — a random tie-break
// would make the gate flaky for no reason.
func LongestOptionIndex(options []string) int {
	best := 0
	bestLen := utf8.RuneCountInString(options[0])
	for i := 1; i < len(options); i++ {
		if l := utf8.RuneCountInString(options[i]); l > bestLen {
			best, bestLen = i, l
		}
	}
	return best
}

// HeuristicResult is one heuristic's hit count against a set of questions,
// plus the sum of those same questions' own baselines. Summing per-question
// baselines (rather than storing one shared value) is what lets a mixed
// corpus report a correct average baseline for the exact subset of
// questions a heuristic was actually evaluated against.
type HeuristicResult struct {
	Hits        int
	Total       int
	baselineSum float64
}

func (h *HeuristicResult) add(hit bool, baseline float64) {
	h.Total++
	h.baselineSum += baseline
	if hit {
		h.Hits++
	}
}

// HitRate is the fraction of considered questions the heuristic guessed
// correctly.
func (h HeuristicResult) HitRate() float64 {
	if h.Total == 0 {
		return 0
	}
	return float64(h.Hits) / float64(h.Total)
}

// AvgBaseline is the mean random-guess baseline across the same questions
// HitRate was computed over.
func (h HeuristicResult) AvgBaseline() float64 {
	if h.Total == 0 {
		return 0
	}
	return h.baselineSum / float64(h.Total)
}

// ExcessRatio is (HitRate-AvgBaseline)/AvgBaseline: a relative measure, not
// percentage points. Percentage points are not comparable across corpora
// with different baselines — 5 points over a 25% baseline is a 20%
// relative jump, but the same 5 points over a 33.3% baseline is only a 15%
// jump — and reintroducing that at the excess-metric level would repeat,
// one layer up, the exact bug (a single shared baseline standing in for
// corpora that do not share one) this package exists to fix.
func (h HeuristicResult) ExcessRatio() float64 {
	b := h.AvgBaseline()
	if b == 0 {
		return 0
	}
	return (h.HitRate() - b) / b
}
