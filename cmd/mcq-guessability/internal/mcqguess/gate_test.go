package mcqguess

import (
	"strings"
	"testing"
)

func reportWithLongestResult(hits, total int, baselinePerQuestion float64) Report {
	r := newReport("")
	for i := 0; i < total; i++ {
		r.Longest.add(i < hits, baselinePerQuestion)
	}
	return r
}

func reportWithLongestAndPosition(longestHits, longestTotal, positionIndex, positionHits, positionTotal int, baselinePerQuestion float64) Report {
	r := reportWithLongestResult(longestHits, longestTotal, baselinePerQuestion)
	hr := r.Position[positionIndex]
	for i := 0; i < positionTotal; i++ {
		hr.add(i < positionHits, baselinePerQuestion)
	}
	r.Position[positionIndex] = hr
	if positionIndex > r.MaxIndex {
		r.MaxIndex = positionIndex
	}
	return r
}

func TestEvaluateGate_NegativeMaxExcessDisablesGate(t *testing.T) {
	r := reportWithLongestResult(100, 100, 0.25) // maximally over baseline
	got := EvaluateGate(r, -1)
	if got.Failed {
		t.Errorf("EvaluateGate() with negative maxExcessRatio Failed = true, want false (gate disabled)")
	}
	if len(got.Violations) != 0 {
		t.Errorf("EvaluateGate() with negative maxExcessRatio Violations = %v, want none", got.Violations)
	}
}

func TestEvaluateGate_PassesBelowThreshold(t *testing.T) {
	// hit rate 30% vs baseline 25% -> excess ratio 0.20 exactly.
	r := reportWithLongestResult(30, 100, 0.25)
	got := EvaluateGate(r, 0.21)
	if got.Failed {
		t.Errorf("EvaluateGate() Failed = true, want false: excess 0.20 <= max 0.21")
	}
	if !got.Requested {
		t.Error("Requested = false, want true: maxExcessRatio >= 0 means the gate was asked to run")
	}
	if got.Judged != 1 {
		t.Errorf("Judged = %d, want 1: Longest had enough samples and must have been compared, not just passed by default", got.Judged)
	}
}

func TestEvaluateGate_FailsAboveThreshold(t *testing.T) {
	r := reportWithLongestResult(30, 100, 0.25) // excess ratio 0.20 exactly
	got := EvaluateGate(r, 0.19)
	if !got.Failed {
		t.Errorf("EvaluateGate() Failed = false, want true: excess 0.20 > max 0.19")
	}
	if len(got.Violations) == 0 {
		t.Errorf("EvaluateGate() Violations empty, want at least one explaining the failure")
	}
}

func TestEvaluateGate_ExactlyAtThresholdPasses(t *testing.T) {
	// excess strictly greater than max fails; excess == max must pass, so
	// the boundary itself is not accidentally inclusive on the wrong side.
	r := reportWithLongestResult(30, 100, 0.25) // excess ratio 0.20 exactly
	got := EvaluateGate(r, 0.20)
	if got.Failed {
		t.Errorf("EvaluateGate() Failed = true, want false: excess 0.20 == max 0.20 should pass")
	}
}

func TestEvaluateGate_PositionAloneCanFailWithLongestClean(t *testing.T) {
	// Regression for a mutation that deletes EvaluateGate's entire position
	// loop: every pass/fail test above builds its Report via
	// reportWithLongestResult, which leaves Position empty, so deleting the
	// position-checking loop still leaves every one of them green. This
	// Report has a perfectly clean Longest (hit rate == baseline, excess 0)
	// and a dirty Position[2] (excess 0.28 against a 0.20 threshold), so it
	// can only fail if the position loop actually runs.
	r := reportWithLongestAndPosition(25, 100, 2, 32, 100, 0.25)
	if got, want := r.Longest.ExcessRatio(), 0.0; !almostEqual(got, want, floatEps) {
		t.Fatalf("fixture sanity check: Longest.ExcessRatio() = %v, want %v (must be clean)", got, want)
	}
	if got, want := r.Position[2].ExcessRatio(), 0.28; !almostEqual(got, want, floatEps) {
		t.Fatalf("fixture sanity check: Position[2].ExcessRatio() = %v, want %v", got, want)
	}

	g := EvaluateGate(r, 0.20)
	if !g.Failed {
		t.Fatal("EvaluateGate() Failed = false, want true: position index 2 alone should fail the gate even though length is clean")
	}
	found := false
	for _, v := range g.Violations {
		if strings.Contains(v, "index 2") {
			found = true
		}
	}
	if !found {
		t.Errorf("Violations = %v, want one naming position index 2", g.Violations)
	}
}

func TestEvaluateGate_BelowMinSampleSizeIsInsufficientNotFailedNorJudgedPass(t *testing.T) {
	// A heuristic with a huge excess but too few samples must not fail the
	// gate — it must be reported as insufficient-n instead, so a track as
	// small as this repo's domain-driven-design (11 mcqs) never gets a
	// false FAIL purely from sampling noise. But it also must not report a
	// clean PASS: nothing here was actually judged (Longest is this
	// Report's only heuristic, and it was skipped), so ExitCode must treat
	// this the same as a gate that measured nothing — exit 1, not 0 (N1).
	r := reportWithLongestResult(MinSampleSize-1, MinSampleSize-1, 0.25) // 100% hit rate, n just below minimum
	g := EvaluateGate(r, 0.20)
	if g.Failed {
		t.Errorf("EvaluateGate() Failed = true, want false: n=%d is below MinSampleSize=%d", MinSampleSize-1, MinSampleSize)
	}
	if len(g.Insufficient) == 0 {
		t.Error("Insufficient is empty, want an entry explaining the sample was too small to judge")
	}
	if g.Judged != 0 {
		t.Errorf("Judged = %d, want 0: nothing met MinSampleSize", g.Judged)
	}
	if got := ExitCode(g); got != 1 {
		t.Errorf("ExitCode() = %d, want 1: a gate that judged nothing must not report success", got)
	}
}

func TestEvaluateTrackGates_NegativeMaxExcessDisables(t *testing.T) {
	byTrack := map[string]Report{"t": reportWithLongestResult(100, 100, 0.25)}
	g := EvaluateTrackGates(byTrack, -1)
	if g.Failed {
		t.Error("EvaluateTrackGates() with negative maxExcessRatio Failed = true, want false")
	}
}

func TestEvaluateTrackGates_PassesWhenAllTracksPass(t *testing.T) {
	byTrack := map[string]Report{
		"a": reportWithLongestResult(25, 100, 0.25),
		"b": reportWithLongestResult(26, 100, 0.25),
	}
	g := EvaluateTrackGates(byTrack, 0.20)
	if g.Failed {
		t.Errorf("EvaluateTrackGates() Failed = true, want false: both tracks are near baseline. Violations=%v", g.Violations)
	}
}

func TestEvaluateTrackGates_FailsWhenAnyTrackFails(t *testing.T) {
	byTrack := map[string]Report{
		"clean": reportWithLongestResult(25, 100, 0.25),
		"bad":   reportWithLongestResult(60, 100, 0.25), // excess 140%
	}
	g := EvaluateTrackGates(byTrack, 0.20)
	if !g.Failed {
		t.Fatal(`EvaluateTrackGates() Failed = false, want true: track "bad" is far over threshold`)
	}
	found := false
	for _, v := range g.Violations {
		if strings.Contains(v, "[bad]") {
			found = true
		}
	}
	if !found {
		t.Errorf("Violations = %v, want one prefixed with the failing track's name", g.Violations)
	}
}

func TestEvaluateTrackGates_DilutionExample_PerTrackCatchesWhatPooledWouldMiss(t *testing.T) {
	// Reproduces the code-review's dilution scenario: a large fair track
	// (0% excess) pooled with a smaller bad track (60% excess) reads only
	// 15% pooled — comfortably under a 20% threshold the bad track alone
	// fails badly. Gating per track catches it; gating on the pooled
	// Report alone would not (asserted below as the negative case).
	fairOptions := []string{"opt-A", "opt-B", "opt-C", "opt-D"} // equal length -> tie, tie-break picks index 0
	fairIdxs := make([]int, 300)
	for i := range fairIdxs {
		fairIdxs[i] = i % 4 // 75 each of 0,1,2,3 -> hit rate exactly baseline
	}
	fair := questionsWithExpectedIndices("fair-track", fairOptions, fairIdxs)

	badOptions := []string{"short-a", "short-b", "short-c", "the much longer fourth option here"}
	var badIdxs []int
	for i := 0; i < 40; i++ {
		badIdxs = append(badIdxs, 3) // hits: correct answer is the unique longest option
	}
	for i := 0; i < 20; i++ {
		badIdxs = append(badIdxs, 0)
	}
	for i := 0; i < 20; i++ {
		badIdxs = append(badIdxs, 1)
	}
	for i := 0; i < 20; i++ {
		badIdxs = append(badIdxs, 2)
	}
	bad := questionsWithExpectedIndices("bad-track", badOptions, badIdxs)

	overall, byTrack, err := Measure(append(fair, bad...))
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}

	if got, want := byTrack["fair-track"].Longest.ExcessRatio(), 0.0; !almostEqual(got, want, floatEps) {
		t.Fatalf("fixture sanity check: fair-track excess = %v, want %v", got, want)
	}
	if got, want := byTrack["bad-track"].Longest.ExcessRatio(), 0.60; !almostEqual(got, want, floatEps) {
		t.Fatalf("fixture sanity check: bad-track excess = %v, want %v", got, want)
	}
	if got, want := overall.Longest.ExcessRatio(), 0.15; !almostEqual(got, want, floatEps) {
		t.Fatalf("fixture sanity check: pooled excess = %v, want %v", got, want)
	}

	const maxExcess = 0.20
	pooledGate := EvaluateGate(overall, maxExcess)
	if pooledGate.Failed {
		t.Fatal("fixture sanity check: pooled gate Failed = true, want false — that dilution is exactly what this test demonstrates")
	}

	perTrackGate := EvaluateTrackGates(byTrack, maxExcess)
	if !perTrackGate.Failed {
		t.Fatal("EvaluateTrackGates() Failed = false, want true: bad-track's 60% excess must fail on its own even though pooling hides it")
	}
}

func TestMeasure_CorpusStraddlingGate_ExitCodeFlips(t *testing.T) {
	// 100 four-option questions (baseline 0.25 each). Options are built so
	// index 3 is unambiguously the longest for every question — which means
	// "guess longest" and "guess index 3" are the same strategy on this
	// fixture, so the longest-option heuristic and the position-3 heuristic
	// move together by construction, not by coincidence. 30 questions have
	// their correct answer at index 3 (both heuristics "hit" there), 70 do
	// not (24 at index 0, 23 at index 1, 23 at index 2) — so both the
	// longest-option heuristic and the position-3 heuristic have hit rate
	// exactly 30/100 = 0.30 against a 0.25 baseline: excess ratio 0.20
	// exactly for BOTH. Indices 0-2 sit below baseline (no risk of crossing
	// a positive threshold), but at max=0.19 EvaluateGate reports two
	// violations (longest-option AND position index 3), not one — asserted
	// below. See TestEvaluateGate_PositionAloneCanFailWithLongestClean for a
	// fixture that isolates the position check from the longest-option
	// heuristic instead.
	options := []string{"short-a", "short-b", "short-c", "the much longer fourth option here"}
	var questions []Question
	addQuestion := func(idx int) {
		questions = append(questions, Question{
			Track:          "gate-test",
			Options:        options,
			ExpectedAnswer: options[idx],
		})
	}
	for i := 0; i < 30; i++ {
		addQuestion(3)
	}
	for i := 0; i < 24; i++ {
		addQuestion(0)
	}
	for i := 0; i < 23; i++ {
		addQuestion(1)
	}
	for i := 0; i < 23; i++ {
		addQuestion(2)
	}

	overall, _, err := Measure(questions)
	if err != nil {
		t.Fatalf("Measure() error = %v", err)
	}
	if got, want := overall.Longest.ExcessRatio(), 0.20; !almostEqual(got, want, floatEps) {
		t.Fatalf("fixture sanity check failed: Longest.ExcessRatio() = %v, want %v", got, want)
	}

	justBelow := EvaluateGate(overall, 0.21)
	if ExitCode(justBelow) != 0 {
		t.Errorf("ExitCode() at max=0.21 (just above the 0.20 excess) = %d, want 0", ExitCode(justBelow))
	}

	justAbove := EvaluateGate(overall, 0.19)
	if ExitCode(justAbove) != 1 {
		t.Errorf("ExitCode() at max=0.19 (just below the 0.20 excess) = %d, want 1", ExitCode(justAbove))
	}
	if len(justAbove.Violations) != 2 {
		t.Errorf("EvaluateGate() at max=0.19 Violations = %v, want exactly 2 (longest-option AND position index 3 — they move together by construction on this fixture)", justAbove.Violations)
	}
}

func TestExitCode(t *testing.T) {
	tests := []struct {
		name string
		g    GateResult
		want int
	}{
		{name: "gate not requested (disabled)", g: GateResult{Requested: false, Failed: true}, want: 0},
		{name: "requested, judged, passed", g: GateResult{Requested: true, Judged: 3, Failed: false}, want: 0},
		{name: "requested, judged, failed", g: GateResult{Requested: true, Judged: 3, Failed: true}, want: 1},
		{name: "requested, nothing judged (N1)", g: GateResult{Requested: true, Judged: 0, Failed: false}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitCode(tt.g); got != tt.want {
				t.Errorf("ExitCode(%+v) = %d, want %d", tt.g, got, tt.want)
			}
		})
	}
}
