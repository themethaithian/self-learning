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

func TestMeasure_CorrectAnswerAlwaysLongest_LengthHeuristicIsHundredPercent(t *testing.T) {
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
	if got := overall.Length.HitRate(); got != 1.0 {
		t.Errorf("Length.HitRate() = %v, want 1.0 (100%%)", got)
	}
}

func TestMeasure_AllOptionsEqualLength_LengthHeuristicNearBaseline(t *testing.T) {
	// All three options tie for longest, so LongestOptionIndex always picks
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
	if got, want := overall.Length.HitRate(), overall.Length.AvgBaseline(); !almostEqual(got, want, floatEps) {
		t.Errorf("Length.HitRate() = %v, want == AvgBaseline() %v (within %v)", got, want, floatEps)
	}
	if got, want := overall.Length.ExcessRatio(), 0.0; !almostEqual(got, want, floatEps) {
		t.Errorf("Length.ExcessRatio() = %v, want %v", got, want)
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

	if got, want := byTrack["three-opt"].Length.AvgBaseline(), 1.0/3.0; !almostEqual(got, want, floatEps) {
		t.Errorf("three-opt track baseline = %v, want %v", got, want)
	}
	if got, want := byTrack["four-opt"].Length.AvgBaseline(), 0.25; !almostEqual(got, want, floatEps) {
		t.Errorf("four-opt track baseline = %v, want %v", got, want)
	}

	wantOverall := (3.0*(1.0/3.0) + 3.0*0.25) / 6.0 // 0.291666...
	if got := overall.Length.AvgBaseline(); !almostEqual(got, wantOverall, floatEps) {
		t.Errorf("overall baseline = %v, want %v (weighted mean, not either track's own baseline)", got, wantOverall)
	}
	if almostEqual(overall.Length.AvgBaseline(), 1.0/3.0, floatEps) {
		t.Errorf("overall baseline equals the three-opt track's baseline verbatim — per-track aggregation looks collapsed into one shared value")
	}
	if almostEqual(overall.Length.AvgBaseline(), 0.25, floatEps) {
		t.Errorf("overall baseline equals the four-opt track's baseline verbatim — per-track aggregation looks collapsed into one shared value")
	}
}

func TestMeasure_BaselineForThreeOptionCorpusIsOneThird(t *testing.T) {
	questions := questionsWithExpectedIndices("t", []string{"a", "b", "c"}, []int{0, 1, 2, 0, 1, 2})

	overall, _, err := Measure(questions)
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got, want := overall.Length.AvgBaseline(), 1.0/3.0; !almostEqual(got, want, floatEps) {
		t.Errorf("baseline = %v, want %v (1/len(options), not a hardcoded 25%%)", got, want)
	}
}

func TestMeasure_BaselineForFourOptionCorpusIsOneQuarter(t *testing.T) {
	questions := questionsWithExpectedIndices("t", []string{"a", "b", "c", "d"}, []int{0, 1, 2, 3, 0, 1, 2, 3})

	overall, _, err := Measure(questions)
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got, want := overall.Length.AvgBaseline(), 0.25; !almostEqual(got, want, floatEps) {
		t.Errorf("baseline = %v, want %v (1/len(options), not a hardcoded 33%%)", got, want)
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
