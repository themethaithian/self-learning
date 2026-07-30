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

// Baseline is 1/len(Options), never a fixed 33% or 25%: a corpus that mixes
// 3-option and 4-option questions has a different correct baseline per
// question, and collapsing that into one shared number is the exact
// measurement bug this package exists to avoid. Callers must not invoke
// this on a Question with fewer than 2 options (Measure enforces that
// floor before calling it) — 0 options divides by zero and 1 option
// silently reports a baseline of 1.0, both meaningless as a guessing
// problem, so this returns 0 rather than +Inf or NaN for that unreachable
// case.
func (q Question) Baseline() float64 {
	if len(q.Options) == 0 {
		return 0
	}
	return 1 / float64(len(q.Options))
}

// ExpectedIndex reports whether ExpectedAnswer is among Options and, if so,
// at what index. Callers must treat ok == false as a corpus error, not a
// silently-skipped row, since a missing index makes every heuristic's
// hit/miss undefined for that question.
func (q Question) ExpectedIndex() (index int, ok bool) {
	for i, opt := range q.Options {
		if opt == q.ExpectedAnswer {
			return i, true
		}
	}
	return -1, false
}
