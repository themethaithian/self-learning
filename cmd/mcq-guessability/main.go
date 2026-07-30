// Command mcq-guessability measures how easily mcq recall_checks under a
// content tree can be guessed without reading the lesson, using simple
// test-taking heuristics (longest option, fixed position) compared against
// each question's own 1/len(options) baseline — never one hardcoded
// percentage, since content/lessons mixes a 3-option track today and will
// mix a 4-option AWS track soon. See docs/tickets/aws-cert.md's "AWS-S3"
// section for why this exists and docs/tickets/mcq-quality.md for the rules
// it measures.
//
// Usage:
//
//	go run ./cmd/mcq-guessability -dir content/lessons
//	go run ./cmd/mcq-guessability -dir content/lessons -max-excess 0.20
//
// A negative -max-excess (the default) reports only, exit code 0. A
// non-negative -max-excess turns this into a gate, evaluated per track (not
// on the pooled "Overall" figure, which pooling can dilute — see
// mcqguess.EvaluateTrackGates): exit code 1 if any track's length heuristic
// or any position heuristic beats its own baseline by more than that
// relative amount, skipping (not judging) any heuristic below
// mcqguess.MinSampleSize samples.
//
// Scope: Question.ExpectedAnswer is a single string, matched against one
// correct option — this measures single-answer mcq only. AWS-S2's
// multiple-response items ("Select TWO") are not yet representable in the
// content schema this package reads, so they cannot silently slip through
// unmeasured; once AWS-S2 lands, this package's Question/ExpectedIndex will
// need to grow multiple-answer support before this tool's gate can cover
// AWS-C1..C4 in full.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/themethaithian/self-learning/internal/mcqguess"
)

func main() {
	dir := flag.String("dir", "content/lessons", "root directory to scan for lesson JSON files (recurses)")
	maxExcess := flag.Float64("max-excess", -1, "max tolerated relative excess over baseline (e.g. 0.2 = 20%); negative disables the gate; the process exit code is 1 on failure so a CI job or Makefile target can key off it")
	flag.Parse()

	os.Exit(run(*dir, *maxExcess, os.Stdout))
}

func run(dir string, maxExcess float64, out io.Writer) int {
	questions, err := mcqguess.LoadQuestions(dir)
	if err != nil {
		fmt.Fprintln(out, "error:", err)
		return 1
	}
	if len(questions) == 0 {
		fmt.Fprintf(out, "no mcq recall_checks found under %s\n", dir)
		return 1
	}

	overall, byTrack, err := mcqguess.Measure(questions)
	if err != nil {
		fmt.Fprintln(out, "error:", err)
		return 1
	}

	fmt.Fprintln(out, "=== MCQ Guessability Report ===")
	fmt.Fprintf(out, "Root: %s\n\n", dir)
	printByTrack(out, byTrack)
	fmt.Fprintln(out, "--- Overall (all tracks pooled, context only — see Gate note below) ---")
	printReport(out, overall)

	gate := mcqguess.EvaluateTrackGates(byTrack, maxExcess)
	printGate(out, gate, maxExcess)
	return mcqguess.ExitCode(gate)
}

func printByTrack(out io.Writer, byTrack map[string]mcqguess.Report) {
	tracks := make([]string, 0, len(byTrack))
	for t := range byTrack {
		tracks = append(tracks, t)
	}
	sort.Strings(tracks)

	for _, t := range tracks {
		fmt.Fprintf(out, "--- Track: %s ---\n", t)
		printReport(out, byTrack[t])
		fmt.Fprintln(out)
	}
}

func printReport(out io.Writer, r mcqguess.Report) {
	fmt.Fprintf(out, "MCQs: %d\n", r.NumMCQs)
	fmt.Fprintf(out, "Length heuristic   (survives shuffle — reaches the user): %.1f%% actual vs %.1f%% baseline avg (excess %+.1f%% relative)\n",
		r.Length.HitRate()*100, r.Length.AvgBaseline()*100, r.Length.ExcessRatio()*100)
	fmt.Fprintln(out, "Position heuristic (shuffled at render, web/lib/shuffle.ts — content-quality signal only, not user-facing):")
	for pos := 0; pos <= r.MaxIndex; pos++ {
		hr, ok := r.Position[pos]
		if !ok {
			continue
		}
		fmt.Fprintf(out, "  index %d: %.1f%% actual vs %.1f%% baseline avg (excess %+.1f%% relative, n=%d)\n",
			pos, hr.HitRate()*100, hr.AvgBaseline()*100, hr.ExcessRatio()*100, hr.Total)
	}
}

func printGate(out io.Writer, gate mcqguess.GateResult, maxExcess float64) {
	if maxExcess < 0 {
		fmt.Fprintln(out, "\nGate: disabled (pass -max-excess to enable)")
		return
	}
	fmt.Fprintf(out, "\nGate: max relative excess over baseline = %.1f%%\n", maxExcess*100)
	fmt.Fprintf(out, "Gate: evaluated per track (fails if ANY track fails) — the pooled \"Overall\"\n")
	fmt.Fprintf(out, "      figure above is not gated: pooling a large fair track with a small bad\n")
	fmt.Fprintf(out, "      one can dilute the bad track's excess under threshold (see AWS-S3 in\n")
	fmt.Fprintf(out, "      docs/tickets/aws-cert.md).\n")
	if !gate.Failed {
		fmt.Fprintln(out, "Gate: PASS")
	} else {
		fmt.Fprintln(out, "Gate: FAIL")
		for _, v := range gate.Violations {
			fmt.Fprintln(out, " -", v)
		}
	}
	for _, v := range gate.Insufficient {
		fmt.Fprintln(out, " ~ insufficient n, not judged:", v)
	}
}
