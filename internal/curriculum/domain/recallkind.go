package domain

import "fmt"

// RecallKind is the answer format a RecallCheck expects.
type RecallKind struct {
	value string
}

var recallKinds = [...]RecallKind{{value: "short_answer"}, {value: "mcq"}}

func NewRecallKind(raw string) (RecallKind, error) {
	for _, k := range recallKinds {
		if k.value == raw {
			return k, nil
		}
	}
	return RecallKind{}, fmt.Errorf("curriculum: recall check kind %q: %w", raw, ErrInvalidRecallKind)
}

func (k RecallKind) String() string { return k.value }

// IsZero reports whether k was never constructed via NewRecallKind.
func (k RecallKind) IsZero() bool { return k.value == "" }

func (k RecallKind) IsMCQ() bool { return k.value == "mcq" }
