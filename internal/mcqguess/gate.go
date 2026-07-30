package mcqguess

import (
	"fmt"
	"sort"
)

// GateResult is the pass/fail verdict of comparing a Report's heuristics
// against a maximum tolerated relative excess over baseline. Insufficient
// lists heuristics that were skipped rather than judged because their
// sample was below MinSampleSize — those never contribute to Failed.
type GateResult struct {
	Failed       bool
	Violations   []string
	Insufficient []string
}

// MinSampleSize is the smallest heuristic sample EvaluateGate will judge.
// Below it, a heuristic is reported as insufficient-n instead of pass or
// fail, because its excess ratio is mostly sampling noise, not signal.
// Re-derived directly against this package's own gate semantics (exact
// one-sided binomial tail, p = baseline = 0.25, threshold = baseline *
// (1 + max-excess), at max-excess = 0.25): a heuristic that is truly fair
// (hit rate == baseline) still false-FAILs about 19.7% of the time at
// n=30, 6.9% at n=100, and 0.05% at n=550. n=550 would be the safe choice,
// but it would permanently exclude the domain-driven-design track (11
// mcqs today) — and any future small track, plus AWS-S2's 5th
// (multiple-response) option position, which will start small too — from
// ever being gated at all. 100 is a deliberate middle ground between "too
// noisy to trust" and "too strict to ever fire on this repo's real
// tracks", not the statistically safest option.
const MinSampleSize = 100

// EvaluateGate checks report's length heuristic and every position
// heuristic with Total >= MinSampleSize against maxExcessRatio (e.g. 0.20
// means "no heuristic may beat its own baseline by more than 20%,
// relatively"). A negative maxExcessRatio disables the gate — report-only
// mode, the default until a future ticket wires this into CI or a
// Makefile target.
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
		if h.Total < MinSampleSize {
			g.Insufficient = append(g.Insufficient, fmt.Sprintf(
				"%s: n=%d < minimum %d, not gated", name, h.Total, MinSampleSize))
			return
		}
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

// EvaluateTrackGates evaluates every track's Report independently against
// maxExcessRatio and fails if any single one of them does. Gating on a
// pooled "overall" Report instead is deliberately not offered: pooling a
// large, well-behaved track with a small bad one dilutes the bad track's
// excess below threshold — a 550-question track at +60% excess pooled
// with a fair, larger legacy corpus reads only about +32% pooled (worked
// example in docs/tickets/aws-cert.md's AWS-S3 section) — which would let
// exactly the batch this gate exists to catch through silently.
func EvaluateTrackGates(byTrack map[string]Report, maxExcessRatio float64) GateResult {
	if maxExcessRatio < 0 {
		return GateResult{}
	}

	tracks := make([]string, 0, len(byTrack))
	for t := range byTrack {
		tracks = append(tracks, t)
	}
	sort.Strings(tracks)

	var g GateResult
	for _, t := range tracks {
		tg := EvaluateGate(byTrack[t], maxExcessRatio)
		if tg.Failed {
			g.Failed = true
		}
		for _, v := range tg.Violations {
			g.Violations = append(g.Violations, fmt.Sprintf("[%s] %s", t, v))
		}
		for _, v := range tg.Insufficient {
			g.Insufficient = append(g.Insufficient, fmt.Sprintf("[%s] %s", t, v))
		}
	}
	return g
}

// ExitCode maps a gate verdict to a process exit code, 1 on failure and 0
// otherwise.
func ExitCode(g GateResult) int {
	if g.Failed {
		return 1
	}
	return 0
}
