package domain

import "fmt"

// GradedBy records who graded a RecallAttempt's outcome. Both enum members
// exist so the type is ready for the future LLM-grading adapter (the
// Grader port described in CLAUDE.md), but 'llm' is unreachable through
// this ticket's write path on purpose: Service.RecordAttempt always passes
// GradedBySelf and never parses graded_by from client input at all.
type GradedBy struct {
	value string
}

// gradedBys is the single source of truth for which values exist; both
// NewGradedBy and GradedBySelf below derive from it so the two can never
// drift apart.
var gradedBys = [...]GradedBy{{value: "self"}, {value: "llm"}}

// GradedBySelf is the only GradedBy this ticket's endpoint can ever produce.
var GradedBySelf = gradedBys[0]

func NewGradedBy(raw string) (GradedBy, error) {
	for _, g := range gradedBys {
		if g.value == raw {
			return g, nil
		}
	}
	return GradedBy{}, fmt.Errorf("learning: graded by %q: %w", raw, ErrInvalidGradedBy)
}

func (g GradedBy) String() string { return g.value }

// IsZero reports whether g was never constructed via NewGradedBy.
func (g GradedBy) IsZero() bool { return g.value == "" }
