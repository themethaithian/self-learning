package infra

import (
	"path/filepath"
	"testing"
)

var wantAISystemsChapters = []wantContentChapter{
	{"llm-foundations", []string{"what-is-an-llm", "tokens-and-context-window", "sampling-and-temperature", "capabilities-and-limits", "cost-latency-tradeoffs"}},
	{"prompting-and-context", []string{"message-roles", "prompt-engineering-basics", "few-shot-and-structured-output", "context-management", "prompt-injection-intro"}},
	{"tool-and-function-calling", []string{"how-function-calling-works", "tool-schemas", "agentic-tool-loop", "parallel-and-sequential-tools", "tool-error-handling"}},
	{"retrieval-augmented-generation", []string{"embeddings-and-semantic-search", "vector-databases", "chunking-strategies", "rag-pipeline-and-grounding", "rag-vs-finetune-vs-long-context"}},
	{"streaming-and-realtime", []string{"token-streaming-sse", "streaming-ux-partial-output", "cancellation-and-backpressure", "latency-optimization"}},
	{"ai-agents-and-orchestration", []string{"agent-loop-react", "planning-and-decomposition", "agent-memory-and-state", "single-vs-multi-agent", "agent-frameworks-landscape"}},
	{"conversational-and-agent-assist", []string{"session-state-machines", "agent-assist-copilot-patterns", "human-in-the-loop", "handoff-orchestration", "contact-center-flows"}},
	{"guardrails-and-safety", []string{"input-output-validation", "hallucination-mitigation", "pii-and-moderation", "prompt-injection-defense", "confidence-and-fallback"}},
	{"llm-evaluation", []string{"why-llm-eval-is-hard", "golden-sets-and-offline-eval", "llm-as-judge", "online-eval-and-ab-testing", "eval-metrics-groundedness"}},
	{"production-llm-systems", []string{"llm-gateway-and-routing", "caching-and-cost-control", "retries-fallback-multi-provider", "prompt-versioning", "llm-observability"}},
	{"crm-and-support-ai", []string{"ai-in-crm-and-ticketing", "triage-routing-deflection", "measuring-support-productivity", "integrating-ai-with-support-stack"}},
}

// TestLoadTopic_AISystemsContentFile guards content/curriculum/ai-systems.json
// itself: it is data that can rot (a typo'd slug, a dropped concept, a
// reordered chapter) with no compiler to catch it.
func TestLoadTopic_AISystemsContentFile(t *testing.T) {
	topic, err := LoadTopic(filepath.Join("..", "..", "..", "content", "curriculum", "ai-systems.json"))
	if err != nil {
		t.Fatalf("LoadTopic(ai-systems.json) unexpected error: %v", err)
	}
	if topic.Track().String() != "ai-systems" || topic.Slug().String() != "ai-and-llm-systems" {
		t.Fatalf("topic = {%q %q}, want {ai-systems ai-and-llm-systems}", topic.Track().String(), topic.Slug().String())
	}

	chapters := topic.Chapters()
	if len(chapters) != len(wantAISystemsChapters) {
		t.Fatalf("len(Chapters()) = %d, want %d", len(chapters), len(wantAISystemsChapters))
	}

	for i, wantCh := range wantAISystemsChapters {
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
