package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// CheckKey content-addresses one recall_checks question:
// SHA256("topic/concept/canonical-question"), hex-encoded lowercase. It
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
type CheckKey struct {
	value string
}

// NewCheckKey validates topic itself (there is no dedicated VO for a topic
// slug in this package yet, so its shape is checked directly), but trusts
// concept's shape as already enforced by LessonRef and question's as already
// enforced by CanonicalQuestion — only their zero values are checked here,
// never re-deriving rules those constructors already guarantee.
func NewCheckKey(topic string, concept LessonRef, question CanonicalQuestion) (CheckKey, error) {
	if !IsValidSlugShape(topic) {
		return CheckKey{}, fmt.Errorf("learning: check key: topic %q: %w", topic, ErrInvalidCheckKey)
	}
	if concept.IsZero() {
		return CheckKey{}, fmt.Errorf("learning: check key: concept is zero: %w", ErrInvalidCheckKey)
	}
	if question.IsZero() {
		return CheckKey{}, fmt.Errorf("learning: check key: question is zero: %w", ErrInvalidCheckKey)
	}

	sum := sha256.Sum256([]byte(topic + "/" + concept.String() + "/" + question.String()))
	return CheckKey{value: hex.EncodeToString(sum[:])}, nil
}

func (k CheckKey) String() string { return k.value }

// IsZero reports whether k was never constructed via NewCheckKey.
func (k CheckKey) IsZero() bool { return k.value == "" }
