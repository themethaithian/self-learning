package domain

import (
	"fmt"
	"slices"
)

const minMCQOptions = 2

// RecallCheck is one recall-check question within a Lesson: a short_answer
// prompt graded against expectedAnswer, or an mcq prompt with its options.
// explanation is optional (see NewRecallCheck) and, when present, must never
// reach the client before the answer itself is revealed.
type RecallCheck struct {
	position       Position
	kind           RecallKind
	question       string
	expectedAnswer string
	options        []string
	explanation    string
}

// NewRecallCheck constructs a validated RecallCheck. An mcq kind requires at
// least 2 options; a short_answer kind must have none — the two are mutually
// exclusive by construction, never mixed at read time. explanation is
// optional: an empty or whitespace-only string means "no explanation", not
// an error.
func NewRecallCheck(position Position, kind RecallKind, question, expectedAnswer string, options []string, explanation string) (RecallCheck, error) {
	if position.IsZero() {
		return RecallCheck{}, fmt.Errorf("curriculum: recall check: position: %w", ErrInvalidPosition)
	}
	if kind.IsZero() {
		return RecallCheck{}, fmt.Errorf("curriculum: recall check: kind: %w", ErrInvalidRecallKind)
	}
	trimmedQuestion, err := validateRecallQuestion(question)
	if err != nil {
		return RecallCheck{}, fmt.Errorf("curriculum: recall check %d: question: %w", position.Int(), err)
	}
	trimmedAnswer, err := validateRecallAnswer(expectedAnswer)
	if err != nil {
		return RecallCheck{}, fmt.Errorf("curriculum: recall check %d: expected answer: %w", position.Int(), err)
	}
	trimmedExplanation, err := validateRecallExplanation(explanation)
	if err != nil {
		return RecallCheck{}, fmt.Errorf("curriculum: recall check %d: explanation: %w", position.Int(), err)
	}

	if kind.IsMCQ() {
		if len(options) < minMCQOptions {
			return RecallCheck{}, fmt.Errorf("curriculum: recall check %d: mcq with %d options: %w", position.Int(), len(options), ErrInvalidRecallOptions)
		}
	} else if len(options) != 0 {
		return RecallCheck{}, fmt.Errorf("curriculum: recall check %d: short_answer with %d options: %w", position.Int(), len(options), ErrInvalidRecallOptions)
	}

	trimmedOptions := make([]string, len(options))
	for i, o := range options {
		trimmed, err := validateRecallOption(o)
		if err != nil {
			return RecallCheck{}, fmt.Errorf("curriculum: recall check %d: option %d: %w", position.Int(), i, err)
		}
		trimmedOptions[i] = trimmed
	}

	// An mcq's expectedAnswer must be one of its options: the client marks
	// the correct option by exact string equality (self-grade reveal UI),
	// so a mismatch here would silently render with no correct marker.
	if kind.IsMCQ() && !slices.Contains(trimmedOptions, trimmedAnswer) {
		return RecallCheck{}, fmt.Errorf("curriculum: recall check %d: expected answer %q not among mcq options: %w", position.Int(), trimmedAnswer, ErrInvalidRecallAnswer)
	}

	return RecallCheck{
		position: position, kind: kind, question: trimmedQuestion,
		expectedAnswer: trimmedAnswer, options: trimmedOptions, explanation: trimmedExplanation,
	}, nil
}

func (rc RecallCheck) Position() Position     { return rc.position }
func (rc RecallCheck) Kind() RecallKind       { return rc.kind }
func (rc RecallCheck) Question() string       { return rc.question }
func (rc RecallCheck) ExpectedAnswer() string { return rc.expectedAnswer }
func (rc RecallCheck) Explanation() string    { return rc.explanation }

// Options returns a copy of the mcq options, or nil for a short_answer
// check; mutating the result does not affect the RecallCheck.
func (rc RecallCheck) Options() []string {
	if len(rc.options) == 0 {
		return nil
	}
	out := make([]string, len(rc.options))
	copy(out, rc.options)
	return out
}

// IsZero reports whether rc was never constructed via NewRecallCheck.
func (rc RecallCheck) IsZero() bool {
	return rc.position.IsZero() || rc.kind.IsZero() || rc.question == ""
}
