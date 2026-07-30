package mcqguess

import "testing"

func reportWithLengthResult(hits, total int, baselinePerQuestion float64) Report {
	r := newReport("")
	for i := 0; i < total; i++ {
		r.Length.add(i < hits, baselinePerQuestion)
	}
	return r
}

func TestEvaluateGate_NegativeMaxExcessDisablesGate(t *testing.T) {
	r := reportWithLengthResult(100, 100, 0.25) // maximally over baseline
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
	r := reportWithLengthResult(30, 100, 0.25)
	got := EvaluateGate(r, 0.21)
	if got.Failed {
		t.Errorf("EvaluateGate() Failed = true, want false: excess 0.20 <= max 0.21")
	}
}

func TestEvaluateGate_FailsAboveThreshold(t *testing.T) {
	r := reportWithLengthResult(30, 100, 0.25) // excess ratio 0.20 exactly
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
	r := reportWithLengthResult(30, 100, 0.25) // excess ratio 0.20 exactly
	got := EvaluateGate(r, 0.20)
	if got.Failed {
		t.Errorf("EvaluateGate() Failed = true, want false: excess 0.20 == max 0.20 should pass")
	}
}

func TestMeasure_CorpusStraddlingGate_ExitCodeFlips(t *testing.T) {
	// 100 four-option questions (baseline 0.25 each). Options are built so
	// index 3 is unambiguously the longest for every question. 30 questions
	// have their correct answer at index 3 (the length heuristic "hits"),
	// 70 do not (24 at index 0, 23 at index 1, 23 at index 2) — so the
	// length heuristic's hit rate is exactly 30/100 = 0.30 against a 0.25
	// baseline: excess ratio (0.30-0.25)/0.25 = 0.20 exactly. Position
	// heuristics stay comfortably under any positive threshold (index 3
	// mirrors the length heuristic at +0.20 by construction; indices 0-2
	// sit below baseline), so only this one number decides the flip.
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
	if got, want := overall.Length.ExcessRatio(), 0.20; !almostEqual(got, want, floatEps) {
		t.Fatalf("fixture sanity check failed: Length.ExcessRatio() = %v, want %v", got, want)
	}

	justBelow := EvaluateGate(overall, 0.21)
	if ExitCode(justBelow) != 0 {
		t.Errorf("ExitCode() at max=0.21 (just above the 0.20 excess) = %d, want 0", ExitCode(justBelow))
	}

	justAbove := EvaluateGate(overall, 0.19)
	if ExitCode(justAbove) != 1 {
		t.Errorf("ExitCode() at max=0.19 (just below the 0.20 excess) = %d, want 1", ExitCode(justAbove))
	}
}

func TestExitCode(t *testing.T) {
	tests := []struct {
		name string
		g    GateResult
		want int
	}{
		{name: "passed", g: GateResult{Failed: false}, want: 0},
		{name: "failed", g: GateResult{Failed: true}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitCode(tt.g); got != tt.want {
				t.Errorf("ExitCode(%+v) = %d, want %d", tt.g, got, tt.want)
			}
		})
	}
}
