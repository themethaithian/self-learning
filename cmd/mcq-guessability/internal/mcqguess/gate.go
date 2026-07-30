package mcqguess

import (
	"fmt"
	"sort"
)

// GateResult is the verdict of comparing a Report's heuristics against a
// maximum tolerated relative excess over baseline. It has three terminal
// states, not two: Requested is false when the gate was never asked to run
// (maxExcessRatio < 0 — report-only mode); when Requested is true, Judged
// counts how many heuristics actually had enough samples to compare against
// maxExcessRatio (see MinSampleSize) — Failed is only meaningful once
// Judged > 0. A Requested gate with Judged == 0 measured nothing at all
// (every heuristic's sample was too small), and ExitCode treats that the
// same as failure: a gate that could not measure anything must not report
// success, the same principle run() already applies when a directory has
// zero mcq questions in it. Insufficient lists the heuristics that were
// skipped rather than judged; those never contribute to Failed or Judged.
type GateResult struct {
	Requested    bool
	Failed       bool
	Judged       int
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
// but it would exclude tracks sized 100-549 from ever being gated at all —
// this repo's own ai-and-llm-systems (158 mcqs) and
// designing-data-intensive-applications (187 mcqs) tracks both clear 100
// today and would both lose gate coverage permanently at 550, and AWS-S2's
// 5th (multiple-response) option position will start small too. (A track
// as small as domain-driven-design, 11 mcqs today, is excluded by both 100
// and 550 — it is not the deciding example between them.) 100 is a
// deliberate middle ground between "too noisy to trust" and "too strict to
// ever fire on this repo's real tracks", not the statistically safest
// option.
const MinSampleSize = 100

// EvaluateGate checks report's longest-option, shortest-option, and
// middle-option heuristics, plus every position heuristic, each with
// Total >= MinSampleSize, against maxExcessRatio (e.g. 0.20 means "no
// heuristic may beat its own baseline by more than 20%, relatively"). A
// negative maxExcessRatio disables the gate — report-only mode, the
// default until a future ticket wires this into CI or a Makefile target.
//
// All three length heuristics and every position heuristic are checked,
// even though Q-1 (web/lib/shuffle.ts) shuffles mcq options at render time
// so a stored position bias never reaches the user: a position bias still
// signals something wrong in how distractors were authored, and shuffling
// only hides the symptom for this app's UI, not the underlying corpus
// defect. The longest/shortest/middle heuristics, unlike position, survive
// the shuffle and reach the user directly.
func EvaluateGate(report Report, maxExcessRatio float64) GateResult {
	if maxExcessRatio < 0 {
		return GateResult{}
	}

	g := GateResult{Requested: true}
	check := func(name string, h HeuristicResult) {
		if h.Total < MinSampleSize {
			g.Insufficient = append(g.Insufficient, fmt.Sprintf(
				"%s: n=%d < minimum %d, not gated", name, h.Total, MinSampleSize))
			return
		}
		g.Judged++
		excess := h.ExcessRatio()
		if excess > maxExcessRatio {
			g.Failed = true
			g.Violations = append(g.Violations, fmt.Sprintf(
				"%s: excess %.1f%% > max %.1f%% (hit rate %.1f%% vs baseline %.1f%%, n=%d)",
				name, excess*100, maxExcessRatio*100, h.HitRate()*100, h.AvgBaseline()*100, h.Total))
		}
	}

	check("longest-option heuristic", report.Longest)
	check("shortest-option heuristic", report.Shortest)
	check("middle-option heuristic", report.Middle)
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

	g := GateResult{Requested: true}
	for _, t := range tracks {
		tg := EvaluateGate(byTrack[t], maxExcessRatio)
		if tg.Failed {
			g.Failed = true
		}
		g.Judged += tg.Judged
		for _, v := range tg.Violations {
			g.Violations = append(g.Violations, fmt.Sprintf("[%s] %s", t, v))
		}
		for _, v := range tg.Insufficient {
			g.Insufficient = append(g.Insufficient, fmt.Sprintf("[%s] %s", t, v))
		}
	}
	return g
}

func ExitCode(g GateResult) int {
	if !g.Requested {
		return 0
	}
	if g.Failed || g.Judged == 0 {
		return 1
	}
	return 0
}
