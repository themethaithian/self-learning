package domain

// Gate computes the gated ChunkState of every lesson in one chapter, given
// their current states in concept-position order.
//
// Invariant: states[0] is always eligible; states[i] (i>0) is eligible only
// when states[i-1] is Passed — you can never skip ahead. A Passed lesson
// always stays Passed. An eligible lesson that is not yet Passed reads as
// InProgress (there is no separate "unlocked but unread" state); anything
// past the first non-Passed lesson stays Locked.
//
// Gate is pure and total: same input always produces the same output, and
// it neither mutates states nor needs a database round trip, so callers can
// re-derive gating on every read instead of caching it.
func Gate(states []ChunkState) []ChunkState {
	out := make([]ChunkState, len(states))
	for i, s := range states {
		switch {
		case s.IsPassed():
			out[i] = chunkStatePassed
		case i == 0 || states[i-1].IsPassed():
			out[i] = chunkStateInProgress
		default:
			out[i] = chunkStateLocked
		}
	}
	return out
}
