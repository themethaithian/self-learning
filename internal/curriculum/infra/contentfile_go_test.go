package infra

import (
	"path/filepath"
	"testing"
)

var wantGoChapters = []wantContentChapter{
	{"runtime-scheduler", []string{"gmp-model", "goroutine-stacks-growth", "preemption", "netpoller", "sysmon"}},
	{"memory-gc", []string{"escape-analysis", "allocator-layout", "gc-tricolor-write-barriers", "gc-pacing-gogc", "stack-vs-heap-performance"}},
	{"concurrency-internals", []string{"channel-internals-hchan", "select-implementation", "mutex-internals", "waitgroup-once-cond", "atomics-memory-model", "context-propagation", "common-concurrency-bugs"}},
	{"types-generics", []string{"interface-internals-itab", "nil-interface-pitfalls", "embedding-composition", "generics-implementation", "reflection-cost"}},
	{"performance-tooling", []string{"pprof-cpu-heap", "execution-tracer", "benchmarking-methodology", "race-detector", "compiler-optimizations-pgo"}},
	{"stdlib-internals", []string{"net-http-server", "http-client-transport", "database-sql-pooling", "encoding-json", "errors-wrapping", "slices-maps-internals"}},
}

// TestLoadTopic_GoContentFile guards content/curriculum/go.json itself: it
// is data that can rot (a typo'd slug, a dropped concept, a reordered
// chapter) with no compiler to catch it.
func TestLoadTopic_GoContentFile(t *testing.T) {
	topic, err := LoadTopic(filepath.Join("..", "..", "..", "content", "curriculum", "go.json"))
	if err != nil {
		t.Fatalf("LoadTopic(go.json) unexpected error: %v", err)
	}
	if topic.Track().String() != "go" || topic.Slug().String() != "deep-go" {
		t.Fatalf("topic = {%q %q}, want {go deep-go}", topic.Track().String(), topic.Slug().String())
	}

	chapters := topic.Chapters()
	if len(chapters) != len(wantGoChapters) {
		t.Fatalf("len(Chapters()) = %d, want %d", len(chapters), len(wantGoChapters))
	}

	for i, wantCh := range wantGoChapters {
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
