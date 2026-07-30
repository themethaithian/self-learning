package mcqguess

import "unicode/utf8"

// longestOptionIndex returns the index of the option with the most Unicode
// code points (runes), not bytes: Thai text runs about 3 bytes per
// character, so a byte-based length would rank a short Thai option as
// "longer" than a longer English one, which is a different and wrong
// metric. Ties resolve to the first (lowest-index) option: a real
// "always pick the longest-looking one" test-taker scans top to bottom and
// commits to the first option that looks longest, and a deterministic rule
// keeps this tool's output reproducible run to run — a random tie-break
// would make the gate flaky for no reason. Unexported (rather than guarded
// against an empty slice) because Measure is its only caller and Measure
// already rejects fewer than 2 options before this ever runs — there is no
// legitimate external caller for this to protect against.
func longestOptionIndex(options []string) int {
	best := 0
	bestLen := utf8.RuneCountInString(options[0])
	for i := 1; i < len(options); i++ {
		if l := utf8.RuneCountInString(options[i]); l > bestLen {
			best, bestLen = i, l
		}
	}
	return best
}

// shortestOptionIndex is longestOptionIndex's mirror: same rune-based
// length, same lowest-index tie-break, same reasons. "The correct answer
// is always the shortest option" is exactly as guessable a defect as
// "always the longest" — nothing about the length-tell is specific to one
// direction.
func shortestOptionIndex(options []string) int {
	best := 0
	bestLen := utf8.RuneCountInString(options[0])
	for i := 1; i < len(options); i++ {
		if l := utf8.RuneCountInString(options[i]); l < bestLen {
			best, bestLen = i, l
		}
	}
	return best
}

// middleOptionIndex returns the index of an option whose rune length is
// strictly between the shortest and longest option's length in the same
// question — neither extreme. This is the guess a test-taker makes under
// the rule "the correct answer is never the longest or shortest option",
// which docs/tickets/mcq-quality.md's "บทเรียนที่ 1" names as the exact
// rule that made an earlier version of this repo's own MCQ corpus 89.6%
// guessable: forcing the answer to be non-extreme in a 3-option question
// leaves exactly one option it can be.
//
// ok is false when no option qualifies: with exactly 2 options, every
// option is simultaneously the longest and the shortest, so there is no
// middle to guess; with 3+ options whose lengths collapse to only two
// distinct values (e.g. two short options and one long, or the reverse),
// every option is again an extreme and none is strictly between. Callers
// must exclude such questions from this heuristic's sample rather than
// counting them as a guaranteed miss — the guess is undefined, not wrong.
//
// With 4+ options there can be more than one qualifying "middle" option
// (e.g. lengths 3, 6, 7, 9 has two: 6 and 7). Ties there resolve to the
// lowest index, the same deterministic rule as longestOptionIndex and
// shortestOptionIndex, for the same reason: no RNG, reproducible output.
func middleOptionIndex(options []string) (index int, ok bool) {
	lens := make([]int, len(options))
	minLen, maxLen := 0, 0
	for i, o := range options {
		lens[i] = utf8.RuneCountInString(o)
		if i == 0 || lens[i] < minLen {
			minLen = lens[i]
		}
		if i == 0 || lens[i] > maxLen {
			maxLen = lens[i]
		}
	}
	for i, l := range lens {
		if l > minLen && l < maxLen {
			return i, true
		}
	}
	return -1, false
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
