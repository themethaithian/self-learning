package mcqguess

import (
	"fmt"
	"strings"
	"testing"
)

// All fixtures in this file are built by hand, not sampled at random: every
// expected-answer index is chosen deterministically so a "within tolerance"
// assertion below is only absorbing float64 rounding, never statistical
// variance. That is why floatEps (1e-9) can stay tight without any risk of
// flaking — there is no randomness anywhere in these corpora to flake. A
// real bug this package needs to catch (e.g. an off-by-one dropping or
// gaining a handful of hits out of a few hundred) moves these numbers by at
// least ~0.3 percentage points, six orders of magnitude above floatEps.
const floatEps = 1e-9

func questionsWithExpectedIndices(track string, options []string, expectedIdxs []int) []Question {
	qs := make([]Question, len(expectedIdxs))
	for i, idx := range expectedIdxs {
		qs[i] = Question{
			Track:          track,
			Source:         fmt.Sprintf("synthetic/%s/%d.json", track, i),
			Options:        options,
			ExpectedAnswer: options[idx],
		}
	}
	return qs
}

func TestMeasure_CorrectAnswerAlwaysLongest_LongestHeuristicIsHundredPercent(t *testing.T) {
	options := []string{"short", "a bit longer", "the very longest option here"}
	idxs := make([]int, 10)
	for i := range idxs {
		idxs[i] = 2 // the longest option, unambiguously
	}
	questions := questionsWithExpectedIndices("t", options, idxs)

	overall, _, err := Measure(questions)
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got := overall.Longest.HitRate(); got != 1.0 {
		t.Errorf("Longest.HitRate() = %v, want 1.0 (100%%)", got)
	}
}

func TestMeasure_AllOptionsEqualLength_LongestHeuristicNearBaseline(t *testing.T) {
	// All three options tie for longest, so longestOptionIndex always picks
	// index 0 (documented tie-break). The correct answer cycles evenly
	// through indices 0/1/2 (100 of each across 300 questions), so a
	// heuristic that always guesses index 0 hits exactly 100/300 = 1/3 of
	// the time — precisely the baseline, not by chance but by construction.
	options := []string{"opt-A", "opt-B", "opt-C"} // 5 runes each
	idxs := make([]int, 300)
	for i := range idxs {
		idxs[i] = i % 3
	}
	questions := questionsWithExpectedIndices("t", options, idxs)

	overall, _, err := Measure(questions)
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got, want := overall.Longest.HitRate(), overall.Longest.AvgBaseline(); !almostEqual(got, want, floatEps) {
		t.Errorf("Longest.HitRate() = %v, want == AvgBaseline() %v (within %v)", got, want, floatEps)
	}
	if got, want := overall.Longest.ExcessRatio(), 0.0; !almostEqual(got, want, floatEps) {
		t.Errorf("Longest.ExcessRatio() = %v, want %v", got, want)
	}
}

func TestMeasure_CorrectAnswerAlwaysShortest_ShortestHeuristicIsHundredPercent(t *testing.T) {
	options := []string{"the very shortest one", "a bit longer than that", "long"}
	idxs := make([]int, 10)
	for i := range idxs {
		idxs[i] = 2 // "long" — the shortest option, unambiguously
	}
	questions := questionsWithExpectedIndices("t", options, idxs)

	overall, _, err := Measure(questions)
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got := overall.Shortest.HitRate(); got != 1.0 {
		t.Errorf("Shortest.HitRate() = %v, want 1.0 (100%%)", got)
	}
}

func TestMeasure_CorrectAnswerAlwaysMiddle_MiddleHeuristicIsHundredPercent(t *testing.T) {
	// "สั้น" (4 runes) < "กลางกลาง" (8 runes) < "ยาวมากมากมากมากมาก" (19 runes)
	// — three distinct lengths, so index 1 is the sole middle option: not
	// the longest, not the shortest. This is the exact shape of rule
	// mcq-quality.md's "บทเรียนที่ 1" names as the one that made an earlier
	// corpus 89.6% guessable (see docs/tickets/aws-cert.md's AWS-S3 "N3"
	// finding): forcing the answer to be non-extreme in a 3-option question
	// pins it to a single option.
	options := []string{"สั้น", "กลางกลาง", "ยาวมากมากมากมากมาก"}
	idxs := make([]int, 10)
	for i := range idxs {
		idxs[i] = 1 // the middle-length option, unambiguously
	}
	questions := questionsWithExpectedIndices("t", options, idxs)

	overall, _, err := Measure(questions)
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got := overall.Middle.HitRate(); got != 1.0 {
		t.Errorf("Middle.HitRate() = %v, want 1.0 (100%%)", got)
	}
	if got := overall.Longest.HitRate(); got != 0.0 {
		t.Errorf("Longest.HitRate() = %v, want 0.0: the always-correct answer here is never the longest option", got)
	}
	if got := overall.Shortest.HitRate(); got != 0.0 {
		t.Errorf("Shortest.HitRate() = %v, want 0.0: the always-correct answer here is never the shortest option", got)
	}
}

func TestMeasure_MiddleOptionIndex_BoundaryCases(t *testing.T) {
	// Mirrors middleOptionIndex's own unit tests in heuristics_test.go, but
	// through the full Measure pipeline: a question with no middle option
	// still counts toward Middle's sample (Total == NumMCQs, same as
	// Longest/Shortest) — it is an automatic miss, not an exclusion. See
	// Report's doc comment for why this differs from Position.
	twoOption := questionsWithExpectedIndices("two-opt", []string{"a", "bb"}, []int{0, 1, 0, 1})
	threeOptionTwoLengths := questionsWithExpectedIndices("tied-extremes", []string{"aa", "bb", "cccc"}, []int{0, 1, 2})

	overall, _, err := Measure(append(twoOption, threeOptionTwoLengths...))
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got := overall.Middle.Total; got != 7 {
		t.Errorf("Middle.Total = %d, want 7 (== NumMCQs): a question with no valid middle is a miss, not an exclusion", got)
	}
	if got := overall.Middle.Hits; got != 0 {
		t.Errorf("Middle.Hits = %d, want 0: none of these fixtures has a question with 3 distinct lengths, so none can hit", got)
	}
	if got := overall.Longest.Total; got != 7 {
		t.Errorf("Longest.Total = %d, want 7", got)
	}
}

func TestMeasure_CorrectAnswerAlwaysIndex1_PositionHeuristic(t *testing.T) {
	options := []string{"a", "b", "c", "d"}
	idxs := make([]int, 40)
	for i := range idxs {
		idxs[i] = 1
	}
	questions := questionsWithExpectedIndices("t", options, idxs)

	overall, _, err := Measure(questions)
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got := overall.Position[1].HitRate(); got != 1.0 {
		t.Errorf("Position[1].HitRate() = %v, want 1.0 (100%%)", got)
	}
	for _, idx := range []int{0, 2, 3} {
		if got := overall.Position[idx].HitRate(); got != 0.0 {
			t.Errorf("Position[%d].HitRate() = %v, want 0.0", idx, got)
		}
	}
}

func TestMeasure_MixedOptionCounts_PerTrackBaselinesDifferAndOverallIsNotEither(t *testing.T) {
	threeOpt := questionsWithExpectedIndices("three-opt", []string{"a", "b", "c"}, []int{0, 1, 2})
	fourOpt := questionsWithExpectedIndices("four-opt", []string{"a", "b", "c", "d"}, []int{0, 1, 2})

	overall, byTrack, err := Measure(append(threeOpt, fourOpt...))
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}

	if got, want := byTrack["three-opt"].Longest.AvgBaseline(), 1.0/3.0; !almostEqual(got, want, floatEps) {
		t.Errorf("three-opt track baseline = %v, want %v", got, want)
	}
	if got, want := byTrack["four-opt"].Longest.AvgBaseline(), 0.25; !almostEqual(got, want, floatEps) {
		t.Errorf("four-opt track baseline = %v, want %v", got, want)
	}

	wantOverall := (3.0*(1.0/3.0) + 3.0*0.25) / 6.0 // 0.291666...
	if got := overall.Longest.AvgBaseline(); !almostEqual(got, wantOverall, floatEps) {
		t.Errorf("overall baseline = %v, want %v (weighted mean, not either track's own baseline)", got, wantOverall)
	}
	if almostEqual(overall.Longest.AvgBaseline(), 1.0/3.0, floatEps) {
		t.Errorf("overall baseline equals the three-opt track's baseline verbatim — per-track aggregation looks collapsed into one shared value")
	}
	if almostEqual(overall.Longest.AvgBaseline(), 0.25, floatEps) {
		t.Errorf("overall baseline equals the four-opt track's baseline verbatim — per-track aggregation looks collapsed into one shared value")
	}
}

func TestMeasure_BaselineForThreeOptionCorpusIsOneThird(t *testing.T) {
	questions := questionsWithExpectedIndices("t", []string{"a", "b", "c"}, []int{0, 1, 2, 0, 1, 2})

	overall, _, err := Measure(questions)
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got, want := overall.Longest.AvgBaseline(), 1.0/3.0; !almostEqual(got, want, floatEps) {
		t.Errorf("baseline = %v, want %v (1/len(options), not a hardcoded 25%%)", got, want)
	}
}

func TestMeasure_BaselineForFourOptionCorpusIsOneQuarter(t *testing.T) {
	questions := questionsWithExpectedIndices("t", []string{"a", "b", "c", "d"}, []int{0, 1, 2, 3, 0, 1, 2, 3})

	overall, _, err := Measure(questions)
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got, want := overall.Longest.AvgBaseline(), 0.25; !almostEqual(got, want, floatEps) {
		t.Errorf("baseline = %v, want %v (1/len(options), not a hardcoded 33%%)", got, want)
	}
}

func TestMeasure_TooFewOptionsAborts(t *testing.T) {
	// A single-option question would otherwise pass ExpectedIndex trivially
	// (there is only one option to match) and then silently contribute a
	// Baseline() of 1.0 to AvgBaseline, understating the whole corpus's
	// computed excess without anyone noticing — the same "quietly wrong"
	// failure mode Measure already refuses for a missing expected answer.
	questions := []Question{
		{Track: "t", Source: "one-option.json", Options: []string{"only-choice"}, ExpectedAnswer: "only-choice"},
	}

	_, _, err := Measure(questions)
	if err == nil {
		t.Fatal("Measure() error = nil, want an error naming the too-few-options question")
	}
	if !strings.Contains(err.Error(), "one-option.json") {
		t.Errorf("Measure() error = %q, want it to name the offending source file", err.Error())
	}
	if !strings.Contains(err.Error(), "1 option") {
		t.Errorf("Measure() error = %q, want it to name the actual option count", err.Error())
	}
}

func TestMeasure_ExpectedAnswerNotInOptionsAborts(t *testing.T) {
	questions := []Question{
		{Track: "t", Source: "broken.json", Options: []string{"a", "b", "c"}, ExpectedAnswer: "z"},
	}

	_, _, err := Measure(questions)
	if err == nil {
		t.Fatal("Measure() error = nil, want an error naming the broken question")
	}
	if !strings.Contains(err.Error(), "broken.json") {
		t.Errorf("Measure() error = %q, want it to name the offending source file", err.Error())
	}
	if !strings.Contains(err.Error(), `"z"`) {
		t.Errorf("Measure() error = %q, want it to name the offending expected answer", err.Error())
	}
}
