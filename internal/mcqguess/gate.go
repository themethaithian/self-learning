package mcqguess

import "fmt"

// GateResult is the pass/fail verdict of comparing a Report's heuristics
// against a maximum tolerated relative excess over baseline.
type GateResult struct {
	Failed     bool
	Violations []string
}

// EvaluateGate checks report's length heuristic and every position
// heuristic against maxExcessRatio (e.g. 0.20 means "no heuristic may beat
// its own baseline by more than 20%, relatively"). A negative
// maxExcessRatio disables the gate — report-only mode, the default until a
// future ticket wires this into CI or a Makefile target.
//
// Both the length and every position heuristic are checked, even though
// Q-1 (web/lib/shuffle.ts) shuffles mcq options at render time so a stored
// position bias never reaches the user: a position bias still signals
// something wrong in how distractors were authored, and shuffling only
// hides the symptom for this app's UI, not the underlying corpus defect.
func EvaluateGate(report Report, maxExcessRatio float64) GateResult {
	if maxExcessRatio < 0 {
		return GateResult{}
	}

	var g GateResult
	check := func(name string, h HeuristicResult) {
		excess := h.ExcessRatio()
		if excess > maxExcessRatio {
			g.Failed = true
			g.Violations = append(g.Violations, fmt.Sprintf(
				"%s: excess %.1f%% > max %.1f%% (hit rate %.1f%% vs baseline %.1f%%, n=%d)",
				name, excess*100, maxExcessRatio*100, h.HitRate()*100, h.AvgBaseline()*100, h.Total))
		}
	}

	check("length heuristic", report.Length)
	for pos := 0; pos <= report.MaxIndex; pos++ {
		if hr, ok := report.Position[pos]; ok {
			check(fmt.Sprintf("position heuristic (index %d)", pos), hr)
		}
	}
	return g
}

// ExitCode maps a gate verdict to a process exit code: 1 on failure so a CI
// job or Makefile target can key off it, 0 otherwise.
func ExitCode(g GateResult) int {
	if g.Failed {
		return 1
	}
	return 0
}
