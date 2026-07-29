package domain

import "fmt"

// CheckKind is the answer format a RecallAttempt was taken against. It
// mirrors curriculum.RecallKind's two-value shape locally rather than
// importing it: this bounded context never imports curriculum's domain
// types directly (see LessonRef's doc comment), so it carries its own copy
// of the same closed enum instead of a compile-time dependency.
type CheckKind struct {
	value string
}

// checkKinds is the single source of truth for which values exist;
// NewCheckKind is the only constructor, so no other value can ever exist.
var checkKinds = [...]CheckKind{{value: "short_answer"}, {value: "mcq"}}

func NewCheckKind(raw string) (CheckKind, error) {
	for _, k := range checkKinds {
		if k.value == raw {
			return k, nil
		}
	}
	return CheckKind{}, fmt.Errorf("learning: check kind %q: %w", raw, ErrInvalidCheckKind)
}

func (k CheckKind) String() string { return k.value }

// IsZero reports whether k was never constructed via NewCheckKind.
func (k CheckKind) IsZero() bool { return k.value == "" }

func (k CheckKind) IsMCQ() bool { return k.value == "mcq" }

// RecallAttempt is one self-graded (today) attempt at a recall check: the
// user's Confidence and AttemptOutcome for one CheckKey, optionally naming
// which mcq option they picked. mcq's outcome comes from the client
// matching the selected option to expected_answer, short_answer's from a
// self-rated Pass/Not-yet — both reduce to this same (Confidence,
// AttemptOutcome) pair before reaching this aggregate.
type RecallAttempt struct {
	checkKey       CheckKey
	kind           CheckKind
	confidence     Confidence
	outcome        AttemptOutcome
	selectedOption *string
	gradedBy       GradedBy
}

// NewRecallAttempt validates a RecallAttempt. selectedOption may only be
// non-nil when kind is mcq: a selected option on a short_answer check (a
// free-text answer, never one of a fixed option list) can never mean
// anything, so it is rejected rather than silently stored.
func NewRecallAttempt(checkKey CheckKey, kind CheckKind, confidence Confidence, outcome AttemptOutcome, selectedOption *string, gradedBy GradedBy) (RecallAttempt, error) {
	if checkKey.IsZero() {
		return RecallAttempt{}, fmt.Errorf("learning: recall attempt: check key is zero: %w", ErrInvalidCheckKey)
	}
	if kind.IsZero() {
		return RecallAttempt{}, fmt.Errorf("learning: recall attempt: kind is zero: %w", ErrInvalidCheckKind)
	}
	if confidence.IsZero() {
		return RecallAttempt{}, fmt.Errorf("learning: recall attempt: confidence is zero: %w", ErrInvalidConfidence)
	}
	if outcome.IsZero() {
		return RecallAttempt{}, fmt.Errorf("learning: recall attempt: outcome is zero: %w", ErrInvalidAttemptOutcome)
	}
	if gradedBy.IsZero() {
		return RecallAttempt{}, fmt.Errorf("learning: recall attempt: graded by is zero: %w", ErrInvalidGradedBy)
	}
	if selectedOption != nil && !kind.IsMCQ() {
		return RecallAttempt{}, fmt.Errorf("learning: recall attempt: selected option on a %s check: %w", kind.String(), ErrSelectedOptionNotMCQ)
	}

	return RecallAttempt{
		checkKey: checkKey, kind: kind, confidence: confidence, outcome: outcome,
		selectedOption: copyStringPtr(selectedOption), gradedBy: gradedBy,
	}, nil
}

func (a RecallAttempt) CheckKey() CheckKey      { return a.checkKey }
func (a RecallAttempt) Kind() CheckKind         { return a.kind }
func (a RecallAttempt) Confidence() Confidence  { return a.confidence }
func (a RecallAttempt) Outcome() AttemptOutcome { return a.outcome }
func (a RecallAttempt) GradedBy() GradedBy      { return a.gradedBy }

func (a RecallAttempt) SelectedOption() *string { return copyStringPtr(a.selectedOption) }

// copyStringPtr returns a pointer to a copy of *p, or nil for a nil p, so
// neither the caller nor RecallAttempt can mutate the other's memory
// through a shared pointer.
func copyStringPtr(p *string) *string {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

// IsZero reports whether a was never constructed via NewRecallAttempt.
func (a RecallAttempt) IsZero() bool {
	return a.checkKey.IsZero() || a.kind.IsZero() || a.confidence.IsZero() || a.outcome.IsZero() || a.gradedBy.IsZero()
}
