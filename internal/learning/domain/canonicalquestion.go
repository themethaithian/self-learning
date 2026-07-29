package domain

import (
	"fmt"
	"strings"
)

// CanonicalQuestion is a recall_checks.question value exactly as MySQL
// stores it. NewCheckKey requires one instead of a raw string so a case- or
// accent-variant submission can never be hashed directly: recall_checks.question
// matches under MySQL's case/accent-insensitive collation, so a client's
// variant phrasing can pass a "this question belongs to this lesson" check
// yet still be byte-different from the stored text. The repository is the
// only caller with a canonical row in hand (see
// internal/learning/infra/repository.go's selectRecallCheckExistsSQL); no
// other code in this codebase constructs one from raw client input.
type CanonicalQuestion struct {
	value string
}

// NewCanonicalQuestion trims raw the same way curriculum.NewRecallCheck's
// validateRecallQuestion trims before storing, so the trimmed form matches
// what MySQL holds byte for byte.
func NewCanonicalQuestion(raw string) (CanonicalQuestion, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return CanonicalQuestion{}, fmt.Errorf("learning: canonical question: empty: %w", ErrInvalidCheckKey)
	}
	return CanonicalQuestion{value: trimmed}, nil
}

func (q CanonicalQuestion) String() string { return q.value }

// IsZero reports whether q was never constructed via NewCanonicalQuestion.
func (q CanonicalQuestion) IsZero() bool { return q.value == "" }
