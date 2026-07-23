package infra

import (
	"path/filepath"
	"testing"
)

var wantDDIAChapters = []wantContentChapter{
	{"reliable-scalable-maintainable", []string{"reliability", "scalability-load", "latency-percentiles", "coping-with-load", "maintainability"}},
	{"data-models-query-languages", []string{"relational-vs-document", "object-relational-mismatch", "schema-flexibility", "declarative-vs-imperative-queries", "graph-data-models"}},
	{"storage-and-retrieval", []string{"log-structured-storage", "sstables-lsm-trees", "b-trees", "oltp-vs-olap", "column-oriented-storage"}},
	{"encoding-and-evolution", []string{"language-specific-formats-pitfalls", "json-xml-binary", "protobuf-thrift-avro", "schema-evolution-compatibility", "dataflow-modes"}},
	{"replication", []string{"single-leader", "replication-lag-problems", "multi-leader", "leaderless-dynamo-style", "quorum-consistency"}},
	{"partitioning", []string{"partitioning-key-range", "partitioning-by-hash", "hot-spots-skew", "partitioning-secondary-indexes", "rebalancing-and-routing"}},
	{"transactions", []string{"acid-meaning", "read-committed", "snapshot-isolation-mvcc", "lost-updates-and-write-skew", "serializability"}},
	{"trouble-with-distributed-systems", []string{"unreliable-networks", "unreliable-clocks", "process-pauses", "knowledge-truth-lies", "byzantine-faults"}},
	{"consistency-and-consensus", []string{"linearizability", "ordering-and-causality", "total-order-broadcast", "two-phase-commit", "consensus-and-zookeeper"}},
	{"batch-processing", []string{"mapreduce", "batch-joins", "dataflow-engines", "graphs-iterative", "batch-output"}},
	{"stream-processing", []string{"event-streams-and-brokers", "partitioned-logs", "change-data-capture", "event-sourcing", "stream-joins-and-time", "exactly-once-fault-tolerance"}},
	{"future-of-data-systems", []string{"data-integration", "unbundling-databases", "lambda-kappa", "designing-for-correctness", "ethics-of-data"}},
}

// TestLoadTopic_DDIAContentFile guards content/curriculum/ddia.json itself: it
// is data that can rot (a typo'd slug, a dropped concept, a reordered
// chapter) with no compiler to catch it.
func TestLoadTopic_DDIAContentFile(t *testing.T) {
	topic, err := LoadTopic(filepath.Join("..", "..", "..", "content", "curriculum", "ddia.json"))
	if err != nil {
		t.Fatalf("LoadTopic(ddia.json) unexpected error: %v", err)
	}
	if topic.Track().String() != "ddia" || topic.Slug().String() != "designing-data-intensive-applications" {
		t.Fatalf("topic = {%q %q}, want {ddia designing-data-intensive-applications}", topic.Track().String(), topic.Slug().String())
	}

	chapters := topic.Chapters()
	if len(chapters) != len(wantDDIAChapters) {
		t.Fatalf("len(Chapters()) = %d, want %d", len(chapters), len(wantDDIAChapters))
	}

	for i, wantCh := range wantDDIAChapters {
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
