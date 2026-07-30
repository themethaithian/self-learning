package mcqguess

import "fmt"

// minMeasurableOptions mirrors the curriculum domain's minMCQOptions
// (internal/curriculum/domain/recallcheck.go): an mcq with fewer than 2
// options is not a guessing problem at all. This package bypasses that
// domain check on purpose (see Measure's doc comment), so it must enforce
// the same floor itself — a 1-option question would otherwise pass the
// "expected answer is among the options" check trivially (there is only one
// option to match) and then silently drag AvgBaseline up (Baseline() = 1.0
// for a single option), making the whole corpus's computed excess look
// smaller than it really is.
const minMeasurableOptions = 2

// Report aggregates heuristic results over a group of mcq questions: either
// one track (Track != "") or the whole corpus (Track == ""). Longest,
// Shortest, and Middle are the three length-based heuristics
// docs/tickets/mcq-quality.md's ตัวชี้วัด requires (ยาว/สั้น/กลาง) — all
// three survive Q-1's render-time option shuffle and reach the user, unlike
// Position, and all three have Total == NumMCQs always, including Middle: a
// question whose options collapse to only two distinct lengths has no
// middleOptionIndex answer (see its doc comment), and that counts as a miss
// for Middle rather than an exclusion. This is deliberately unlike Position,
// where a question with too few options is excluded rather than counted as
// a miss — "guess index 3" is undefined for a 2-option question (there is no
// index 3 to be right or wrong about), but "guess the non-extreme option"
// is well-defined and simply unwinnable when every option ties into an
// extreme; excluding it instead of counting the miss would let a corpus
// hide a real middle-length tell behind a shrinking denominator, which is
// exactly how the docs/tickets/mcq-quality.md incident this heuristic
// exists to catch was missed the first time (see docs/tickets/aws-cert.md's
// AWS-S3 "N3" finding: verified against real history, this choice is what
// makes 120/356 and 319/356 reproduce exactly). Position is keyed by
// 0-based guessed index; a HeuristicResult for index p is only built from
// questions that actually have more than p options, so a corpus mixing
// option counts never inflates or deflates a position's own baseline with
// questions the guess "index p" could not even apply to.
type Report struct {
	Track    string
	NumMCQs  int
	Longest  HeuristicResult
	Shortest HeuristicResult
	Middle   HeuristicResult
	Position map[int]HeuristicResult
	MaxIndex int
}

func newReport(track string) Report {
	return Report{Track: track, Position: map[int]HeuristicResult{}}
}

// Measure builds one overall Report plus one Report per track from
// questions. It returns an error the moment a question's expected answer is
// not found among its own options, naming the question's source and
// position in the input: the curriculum domain rejects this at import time
// (RecallCheck requires the answer to be one of the options), but this
// package reads content files directly rather than through the domain, so
// nothing upstream of Measure catches a malformed file first. Silently
// skipping such a question would shrink the sample size without saying so —
// the same kind of quiet wrongness this package exists to catch, not
// commit.
func Measure(questions []Question) (overall Report, byTrack map[string]Report, err error) {
	overall = newReport("")
	byTrack = map[string]Report{}

	for i, q := range questions {
		if len(q.Options) < minMeasurableOptions {
			return Report{}, nil, fmt.Errorf(
				"mcqguess: question %d (source %q): mcq has %d option(s), need at least %d to measure guessability",
				i, q.Source, len(q.Options), minMeasurableOptions)
		}
		expectedIdx, ok := q.ExpectedIndex()
		if !ok {
			return Report{}, nil, fmt.Errorf(
				"mcqguess: question %d (source %q): expected answer %q not among its %d options",
				i, q.Source, q.ExpectedAnswer, len(q.Options))
		}

		track := byTrack[q.Track]
		if track.Position == nil {
			track = newReport(q.Track)
		}

		accumulate(&overall, q, expectedIdx)
		accumulate(&track, q, expectedIdx)
		byTrack[q.Track] = track
	}
	return overall, byTrack, nil
}

func accumulate(r *Report, q Question, expectedIdx int) {
	r.NumMCQs++
	baseline := q.Baseline()

	longest := longestOptionIndex(q.Options)
	r.Longest.add(longest == expectedIdx, baseline)

	shortest := shortestOptionIndex(q.Options)
	r.Shortest.add(shortest == expectedIdx, baseline)

	middle, middleOK := middleOptionIndex(q.Options)
	r.Middle.add(middleOK && middle == expectedIdx, baseline)

	if last := len(q.Options) - 1; last > r.MaxIndex {
		r.MaxIndex = last
	}
	for pos := range q.Options {
		hr := r.Position[pos]
		hr.add(pos == expectedIdx, baseline)
		r.Position[pos] = hr
	}
}
