package domain

import "fmt"

// ChunkState is a LessonProgress's position in the locked -> in_progress ->
// passed lifecycle.
type ChunkState struct {
	value string
}

// chunkStates is the single source of truth for which states exist;
// NewChunkState and the package-level singletons below both derive from it
// so the two can never drift apart.
var chunkStates = [...]ChunkState{{value: "locked"}, {value: "in_progress"}, {value: "passed"}}

var (
	chunkStateLocked     = chunkStates[0]
	chunkStateInProgress = chunkStates[1]
	chunkStatePassed     = chunkStates[2]
)

func NewChunkState(raw string) (ChunkState, error) {
	for _, s := range chunkStates {
		if s.value == raw {
			return s, nil
		}
	}
	return ChunkState{}, fmt.Errorf("learning: chunk state %q: %w", raw, ErrInvalidChunkState)
}

func (s ChunkState) String() string { return s.value }

// IsZero reports whether s was never constructed via NewChunkState.
func (s ChunkState) IsZero() bool { return s.value == "" }

func (s ChunkState) IsLocked() bool     { return s.value == "locked" }
func (s ChunkState) IsInProgress() bool { return s.value == "in_progress" }
func (s ChunkState) IsPassed() bool     { return s.value == "passed" }
