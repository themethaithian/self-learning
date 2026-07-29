package domain

import "fmt"

// Confidence is the user's self-rated recall confidence, captured at
// Commit time — before Reveal — so it reflects genuine recall strength
// rather than being biased by having just seen correct/incorrect.
type Confidence struct {
	value string
}

// confidences is the single source of truth for which values exist;
// NewConfidence is the only constructor, so no other value can ever exist.
var confidences = [...]Confidence{{value: "guessed"}, {value: "unsure"}, {value: "confident"}}

func NewConfidence(raw string) (Confidence, error) {
	for _, c := range confidences {
		if c.value == raw {
			return c, nil
		}
	}
	return Confidence{}, fmt.Errorf("learning: confidence %q: %w", raw, ErrInvalidConfidence)
}

func (c Confidence) String() string { return c.value }

// IsZero reports whether c was never constructed via NewConfidence.
func (c Confidence) IsZero() bool { return c.value == "" }
