package mcqguess

// Question is one mcq recall_check read directly from a lesson content
// file. It intentionally bypasses the curriculum domain: RecallCheck's own
// constructor already rejects an expected answer that is not among its
// options, but this package must be able to report that condition itself
// (see Measure) rather than have the domain silently reject the file first.
type Question struct {
	Track          string
	Source         string
	Options        []string
	ExpectedAnswer string
}

// Baseline is the random-guess hit rate for this question alone. It is
// always 1/len(Options), never a fixed 33% or 25%: a corpus that mixes
// 3-option and 4-option questions has a different correct baseline per
// question, and collapsing that into one shared number is the exact
// measurement bug this package exists to avoid.
func (q Question) Baseline() float64 {
	return 1 / float64(len(q.Options))
}

// ExpectedIndex returns the 0-based index of ExpectedAnswer within Options.
// ok is false when the expected answer is not among the options at all —
// callers must treat that as a corpus error, not a silently-skipped row,
// since a missing index makes every heuristic's hit/miss undefined for that
// question.
func (q Question) ExpectedIndex() (index int, ok bool) {
	for i, opt := range q.Options {
		if opt == q.ExpectedAnswer {
			return i, true
		}
	}
	return -1, false
}
