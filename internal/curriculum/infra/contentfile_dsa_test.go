package infra

import (
	"path/filepath"
	"testing"
)

var wantDsaChapters = []wantContentChapter{
	{"arrays-hashing", []string{"hash-frequency", "prefix-sums", "two-sum-family"}},
	{"two-pointers", []string{"converging", "in-place-partition"}},
	{"sliding-window", []string{"fixed", "variable-shrink"}},
	{"stack", []string{"monotonic-stack", "matching-pairs"}},
	{"binary-search", []string{"on-index", "on-answer-space", "rotated-arrays"}},
	{"linked-list", []string{"reversal", "fast-slow-cycle", "merge"}},
	{"trees", []string{"dfs-patterns", "bfs-level-order", "bst-properties", "lca", "serialization"}},
	{"tries", []string{"build-search", "word-search-with-trie"}},
	{"heap", []string{"top-k", "two-heaps-median", "k-way-merge"}},
	{"backtracking", []string{"subsets-permutations", "constraint-pruning"}},
	{"graphs", []string{"representation-traversal", "islands-components", "topological-sort", "union-find"}},
	{"advanced-graphs", []string{"dijkstra", "mst", "bellman-ford"}},
	{"dp-1d", []string{"memo-vs-tabulation", "house-robber-family", "coin-change", "lis"}},
	{"dp-2d", []string{"grid-paths", "lcs-edit-distance", "knapsack-01"}},
	{"greedy", []string{"exchange-argument", "jump-gas"}},
	{"intervals", []string{"sort-merge", "sweep-line-rooms"}},
	{"math-geometry", []string{"matrix-ops", "number-theory-basics"}},
	{"bit-manipulation", []string{"bit-tricks", "xor-patterns"}},
}

// TestLoadTopic_DsaContentFile guards content/curriculum/dsa.json itself: it
// is data that can rot (a typo'd slug, a dropped concept, a reordered
// chapter) with no compiler to catch it.
func TestLoadTopic_DsaContentFile(t *testing.T) {
	topic, err := LoadTopic(filepath.Join("..", "..", "..", "content", "curriculum", "dsa.json"))
	if err != nil {
		t.Fatalf("LoadTopic(dsa.json) unexpected error: %v", err)
	}
	if topic.Track().String() != "dsa" || topic.Slug().String() != "dsa-neetcode-150" {
		t.Fatalf("topic = {%q %q}, want {dsa dsa-neetcode-150}", topic.Track().String(), topic.Slug().String())
	}

	chapters := topic.Chapters()
	if len(chapters) != len(wantDsaChapters) {
		t.Fatalf("len(Chapters()) = %d, want %d", len(chapters), len(wantDsaChapters))
	}

	for i, wantCh := range wantDsaChapters {
		ch := chapters[i]
		if ch.Slug().String() != wantCh.slug {
			t.Errorf("Chapters()[%d].Slug() = %q, want %q", i, ch.Slug().String(), wantCh.slug)
		}
		if ch.Position().Int() != i+1 {
			t.Errorf("Chapters()[%d].Position() = %d, want %d", i, ch.Position().Int(), i+1)
		}

		concepts := ch.Concepts()
		if len(concepts) != len(wantCh.concepts) {
			t.Fatalf("Chapters()[%d] %q has %d concepts, want %d", i, wantCh.slug, len(concepts), len(wantCh.concepts))
		}
		for j, wantSlug := range wantCh.concepts {
			c := concepts[j]
			if c.Slug().String() != wantSlug {
				t.Errorf("Chapters()[%d].Concepts()[%d].Slug() = %q, want %q", i, j, c.Slug().String(), wantSlug)
			}
			if c.Position().Int() != j+1 {
				t.Errorf("Chapters()[%d].Concepts()[%d].Position() = %d, want %d", i, j, c.Position().Int(), j+1)
			}
		}
	}
}
