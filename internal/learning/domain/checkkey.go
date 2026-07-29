package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// CheckKey content-addresses one recall_checks question:
// SHA256("topic/concept/trimmed-question"), hex-encoded lowercase. It
// exists instead of a foreign key to recall_checks.id because content
// import deletes and re-inserts every recall_checks row on each run (see
// internal/curriculum/infra/lessonwriter.go's deleteRecallChecksSQL /
// insertRecallCheckSQL), so those ids are never stable across a re-import —
// unlike lessons.id, which import preserves via
// "ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id)". (lesson_id, position)
// was rejected too: reordering a lesson's checks on re-import would
// silently reattach an existing attempt history to the wrong question.
// Content-addressing fails safe instead — editing a question's text mints a
// new key and starts fresh history, rather than misattaching the old one.
//
// Only the question is trimmed before hashing — no case-folding or other
// unicode normalization is applied: any such rule would itself be a second
// source of truth that could drift from what curriculum actually stored.
// recall_checks.question is always already-trimmed by
// curriculum.NewRecallCheck's validateRecallQuestion, so hashing the trimmed
// form reproduces what MySQL stores byte for byte — PROVIDED question is
// that canonical, already-stored text. NewCheckKey itself cannot verify
// that: recall_checks.question is matched under MySQL's case/accent-
// insensitive collation (utf8mb4_0900_ai_ci), so a client's case-variant
// phrasing can pass a "this question belongs to this lesson" check yet
// still be byte-different from the stored text. Callers MUST resolve and
// pass the DB's own canonical question, never raw client input — see
// app.Service.RecordAttempt, which only calls NewCheckKey with the
// canonical text the repository resolved, and
// internal/learning/infra/repository.go's selectRecallCheckExistsSQL, which
// returns that canonical text alongside the match.
type CheckKey struct {
	value string
}

// NewCheckKey validates topic and question itself (there is no dedicated
// VO for a topic slug in this package yet, so its shape is checked
// directly), but trusts concept's shape as already enforced by LessonRef —
// only its zero value is checked here, never re-deriving a rule LessonRef's
// own constructor already guarantees.
func NewCheckKey(topic string, concept LessonRef, question string) (CheckKey, error) {
	if !IsValidSlugShape(topic) {
		return CheckKey{}, fmt.Errorf("learning: check key: topic %q: %w", topic, ErrInvalidCheckKey)
	}
	if concept.IsZero() {
		return CheckKey{}, fmt.Errorf("learning: check key: concept is zero: %w", ErrInvalidCheckKey)
	}
	trimmed := strings.TrimSpace(question)
	if trimmed == "" {
		return CheckKey{}, fmt.Errorf("learning: check key: question is empty: %w", ErrInvalidCheckKey)
	}

	sum := sha256.Sum256([]byte(topic + "/" + concept.String() + "/" + trimmed))
	return CheckKey{value: hex.EncodeToString(sum[:])}, nil
}

func (k CheckKey) String() string { return k.value }

// IsZero reports whether k was never constructed via NewCheckKey.
func (k CheckKey) IsZero() bool { return k.value == "" }
