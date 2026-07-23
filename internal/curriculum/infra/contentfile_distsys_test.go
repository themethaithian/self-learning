package infra

import (
	"path/filepath"
	"testing"
)

var wantDistsysChapters = []wantContentChapter{
	{"foundations", []string{"why-distributed", "transparency-goals", "scalability-dimensions", "fallacies-of-distributed-computing"}},
	{"architectures", []string{"client-server", "multi-tier-layered", "peer-to-peer", "microservices-vs-monolith", "event-driven-architecture"}},
	{"communication", []string{"rpc-fundamentals", "message-queues", "publish-subscribe", "rest-vs-grpc", "serialization-formats"}},
	{"naming-discovery", []string{"flat-naming", "structured-naming-dns", "service-discovery"}},
	{"coordination-time", []string{"physical-clocks-ntp", "lamport-clocks", "vector-clocks", "distributed-mutex", "leader-election", "gossip-protocols"}},
	{"consistency-replication", []string{"replication-motivation", "linearizability-sequential", "causal-consistency", "eventual-consistency", "client-centric-models", "quorum-protocols", "crdt-intro"}},
	{"fault-tolerance", []string{"failure-models", "failure-detection", "process-resilience", "consensus-raft", "two-phase-commit", "sagas-compensation", "recovery-checkpointing"}},
	{"interview-patterns", []string{"cap-pacelc", "partitioning-sharding", "idempotency-retries", "backpressure-rate-limiting", "outbox-pattern", "distributed-caching"}},
}

// TestLoadTopic_DistsysContentFile guards content/curriculum/distsys.json
// itself: it is data that can rot (a typo'd slug, a dropped concept, a
// reordered chapter) with no compiler to catch it.
func TestLoadTopic_DistsysContentFile(t *testing.T) {
	topic, err := LoadTopic(filepath.Join("..", "..", "..", "content", "curriculum", "distsys.json"))
	if err != nil {
		t.Fatalf("LoadTopic(distsys.json) unexpected error: %v", err)
	}
	if topic.Track().String() != "distsys" || topic.Slug().String() != "distributed-systems" {
		t.Fatalf("topic = {%q %q}, want {distsys distributed-systems}", topic.Track().String(), topic.Slug().String())
	}

	chapters := topic.Chapters()
	if len(chapters) != len(wantDistsysChapters) {
		t.Fatalf("len(Chapters()) = %d, want %d", len(chapters), len(wantDistsysChapters))
	}

	for i, wantCh := range wantDistsysChapters {
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
