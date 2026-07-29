package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
)

func expectedCheckKeyHex(t *testing.T, topic, concept, question string) string {
	t.Helper()
	sum := sha256.Sum256([]byte(topic + "/" + concept + "/" + question))
	return hex.EncodeToString(sum[:])
}

func TestNewCheckKey(t *testing.T) {
	concept := mustLessonRef(t, "b-trees")
	question := mustCanonicalQuestion(t, "What is a B-tree?")

	tests := []struct {
		name     string
		topic    string
		concept  LessonRef
		question CanonicalQuestion
		wantErr  error
	}{
		{name: "valid", topic: "ddia", concept: concept, question: question},
		{name: "empty topic", topic: "", concept: concept, question: question, wantErr: ErrInvalidCheckKey},
		{name: "uppercase topic is invalid shape", topic: "DDIA", concept: concept, question: question, wantErr: ErrInvalidCheckKey},
		{name: "zero concept", topic: "ddia", concept: LessonRef{}, question: question, wantErr: ErrInvalidCheckKey},
		{name: "zero question", topic: "ddia", concept: concept, question: CanonicalQuestion{}, wantErr: ErrInvalidCheckKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCheckKey(tt.topic, tt.concept, tt.question)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewCheckKey(%q, %q, %q) error = %v, want wrapping %v", tt.topic, tt.concept, tt.question, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewCheckKey(%q, %q, %q) unexpected error: %v", tt.topic, tt.concept, tt.question, err)
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for a valid key")
			}

			want := expectedCheckKeyHex(t, tt.topic, tt.concept.String(), tt.question.String())
			if got.String() != want {
				t.Errorf("String() = %q, want %q (topic/concept/question order)", got.String(), want)
			}
		})
	}
}

// TestNewCheckKey_HashOrderMatters pins the exact concatenation order
// topic + "/" + concept + "/" + question — swapping topic and concept
// must produce a different key, not the same one.
func TestNewCheckKey_HashOrderMatters(t *testing.T) {
	concept := mustLessonRef(t, "concept-a")
	question := mustCanonicalQuestion(t, "question")

	got, err := NewCheckKey("topic-a", concept, question)
	if err != nil {
		t.Fatalf("NewCheckKey() unexpected error: %v", err)
	}

	correctOrder := expectedCheckKeyHex(t, "topic-a", "concept-a", "question")
	if got.String() != correctOrder {
		t.Fatalf("String() = %q, want %q (topic/concept/question order)", got.String(), correctOrder)
	}

	swappedOrder := expectedCheckKeyHex(t, "concept-a", "topic-a", "question")
	if got.String() == swappedOrder {
		t.Fatalf("String() matched the concept/topic/question order — hash input order is not pinned")
	}
}

// TestNewCheckKey_UsesQuestionVerbatim pins that NewCheckKey applies no
// normalization of its own on top of CanonicalQuestion's value — a
// case-different question is a genuinely different CanonicalQuestion (see
// canonicalquestion_test.go) and must hash to a different key.
func TestNewCheckKey_UsesQuestionVerbatim(t *testing.T) {
	concept := mustLessonRef(t, "concept-a")

	base, err := NewCheckKey("topic-a", concept, mustCanonicalQuestion(t, "What is a B-tree?"))
	if err != nil {
		t.Fatalf("NewCheckKey() unexpected error: %v", err)
	}
	differentCase, err := NewCheckKey("topic-a", concept, mustCanonicalQuestion(t, "what is a b-tree?"))
	if err != nil {
		t.Fatalf("NewCheckKey() unexpected error: %v", err)
	}
	if base.String() == differentCase.String() {
		t.Errorf("a case-different question produced the same key — CanonicalQuestion values must be used verbatim, not re-normalized")
	}
}

func TestCheckKeyZeroValue(t *testing.T) {
	var k CheckKey
	if !k.IsZero() {
		t.Errorf("zero-value CheckKey.IsZero() = false, want true")
	}
}
