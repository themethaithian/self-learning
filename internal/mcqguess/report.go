package mcqguess

import "fmt"

// Report aggregates heuristic results over a group of mcq questions: either
// one track (Track != "") or the whole corpus (Track == ""). Position is
// keyed by 0-based guessed index; a HeuristicResult for index p is only
// built from questions that actually have more than p options, so a corpus
// mixing option counts never inflates or deflates a position's own
// baseline with questions the guess "index p" could not even apply to.
type Report struct {
	Track    string
	NumMCQs  int
	Length   HeuristicResult
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

	longest := LongestOptionIndex(q.Options)
	r.Length.add(longest == expectedIdx, baseline)

	if last := len(q.Options) - 1; last > r.MaxIndex {
		r.MaxIndex = last
	}
	for pos := range q.Options {
		hr := r.Position[pos]
		hr.add(pos == expectedIdx, baseline)
		r.Position[pos] = hr
	}
}
