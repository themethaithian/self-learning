package domain

import "fmt"

// AttemptOutcome is the one unified correct/incorrect shape a RecallAttempt
// reduces to regardless of check kind: an mcq's outcome comes from the
// client matching the selected option against expected_answer, a
// short_answer's from the user's self-rated Pass/Not-yet — both arrive here
// as the same (Confidence, AttemptOutcome) pair, which is what Q-2c's SM-2
// quality score will be derived from.
type AttemptOutcome struct {
	value string
}

// attemptOutcomes is the single source of truth for which values exist;
// NewAttemptOutcome is the only constructor, so no other value can ever
// exist.
var attemptOutcomes = [...]AttemptOutcome{{value: "correct"}, {value: "incorrect"}}

func NewAttemptOutcome(raw string) (AttemptOutcome, error) {
	for _, o := range attemptOutcomes {
		if o.value == raw {
			return o, nil
		}
	}
	return AttemptOutcome{}, fmt.Errorf("learning: attempt outcome %q: %w", raw, ErrInvalidAttemptOutcome)
}

func (o AttemptOutcome) String() string { return o.value }

// IsZero reports whether o was never constructed via NewAttemptOutcome.
func (o AttemptOutcome) IsZero() bool { return o.value == "" }
