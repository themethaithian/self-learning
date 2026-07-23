package domain

import "testing"

func statesOf(t *testing.T, raw ...string) []ChunkState {
	t.Helper()
	out := make([]ChunkState, len(raw))
	for i, r := range raw {
		out[i] = mustChunkState(t, r)
	}
	return out
}

func TestGate(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "empty chapter",
			in:   nil,
			want: nil,
		},
		{
			name: "single lesson not started stays unlocked",
			in:   []string{"locked"},
			want: []string{"in_progress"},
		},
		{
			name: "single lesson already passed stays passed",
			in:   []string{"passed"},
			want: []string{"passed"},
		},
		{
			name: "first lesson always unlocked regardless of state",
			in:   []string{"in_progress", "locked", "locked"},
			want: []string{"in_progress", "locked", "locked"},
		},
		{
			name: "sequential unlock only when prior passed",
			in:   []string{"locked", "locked", "locked"},
			want: []string{"in_progress", "locked", "locked"},
		},
		{
			name: "second unlocks once first is passed",
			in:   []string{"passed", "locked", "locked"},
			want: []string{"passed", "in_progress", "locked"},
		},
		{
			name: "a gap: prior in_progress keeps the rest locked",
			in:   []string{"passed", "in_progress", "locked"},
			want: []string{"passed", "in_progress", "locked"},
		},
		{
			name: "all passed unlocks all",
			in:   []string{"passed", "passed", "passed"},
			want: []string{"passed", "passed", "passed"},
		},
		{
			name: "you can never skip: third stays locked even if marked in_progress without second passed",
			in:   []string{"passed", "in_progress", "in_progress"},
			want: []string{"passed", "in_progress", "locked"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Gate(statesOf(t, tt.in...))
			want := statesOf(t, tt.want...)
			if len(got) != len(want) {
				t.Fatalf("Gate() returned %d states, want %d", len(got), len(want))
			}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("Gate()[%d] = %v, want %v", i, got[i], want[i])
				}
			}
		})
	}
}
