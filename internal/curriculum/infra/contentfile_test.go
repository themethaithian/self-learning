package infra

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	curriculumapp "github.com/themethaithian/self-learning/internal/curriculum/app"
	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

var _ curriculumapp.Loader = FileLoader{}

func TestLoadTopic_Valid(t *testing.T) {
	topic, err := LoadTopic(filepath.Join("testdata", "valid.json"))
	if err != nil {
		t.Fatalf("LoadTopic() unexpected error: %v", err)
	}

	if topic.Track().String() != "go" || topic.Slug().String() != "go-basics" ||
		topic.Title() != "Go Basics" || topic.Position().Int() != 1 {
		t.Fatalf("topic = {%q %q %q %d}, want {go go-basics \"Go Basics\" 1}",
			topic.Track().String(), topic.Slug().String(), topic.Title(), topic.Position().Int())
	}

	chapters := topic.Chapters()
	if len(chapters) != 1 {
		t.Fatalf("len(Chapters()) = %d, want 1", len(chapters))
	}
	if chapters[0].Slug().String() != "syntax" {
		t.Fatalf("Chapters()[0].Slug() = %q, want syntax", chapters[0].Slug().String())
	}

	concepts := chapters[0].Concepts()
	if len(concepts) != 2 {
		t.Fatalf("len(Concepts()) = %d, want 2", len(concepts))
	}
	if concepts[0].Slug().String() != "variables" || concepts[1].Slug().String() != "loops" {
		t.Fatalf("concept slugs = %q, %q, want variables, loops", concepts[0].Slug().String(), concepts[1].Slug().String())
	}
	if want := "ตัวแปรใน Go ประกาศแบบไหนได้บ้าง"; concepts[0].Outline() != want {
		t.Errorf("Concepts()[0].Outline() = %q, want %q", concepts[0].Outline(), want)
	}
	if want := "for loop รูปแบบต่าง ๆ ใน Go"; concepts[1].Outline() != want {
		t.Errorf("Concepts()[1].Outline() = %q, want %q", concepts[1].Outline(), want)
	}
}

func TestLoadTopic_Errors(t *testing.T) {
	tests := []struct {
		name        string
		file        string
		wantErr     error
		wantContain []string
	}{
		{name: "malformed JSON", file: "malformed.json", wantContain: []string{"malformed.json"}},
		{name: "unknown top-level field", file: "unknown_field.json", wantContain: []string{"unknown_field.json"}},
		{name: "unknown field nested in a concept", file: "unknown_field_nested.json", wantContain: []string{"unknown_field_nested.json"}},
		{name: "trailing content after the JSON value", file: "trailing_content.json", wantContain: []string{"trailing_content.json"}},
		{name: "unknown track", file: "unknown_track.json", wantErr: domain.ErrInvalidTrack, wantContain: []string{"unknown_track.json"}},
		{name: "duplicate concept slug across chapters", file: "duplicate_concept_slug.json", wantErr: domain.ErrDuplicateSlug, wantContain: []string{"duplicate_concept_slug.json", "idempotency"}},
		{name: "empty chapters", file: "empty_chapters.json", wantErr: domain.ErrNoChildren, wantContain: []string{"empty_chapters.json"}},
		{name: "position zero", file: "position_zero.json", wantErr: domain.ErrInvalidPosition, wantContain: []string{"position_zero.json"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadTopic(filepath.Join("testdata", tt.file))
			if err == nil {
				t.Fatal("LoadTopic() expected error, got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("LoadTopic() error = %v, want it to wrap %v", err, tt.wantErr)
			}
			for _, want := range tt.wantContain {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("LoadTopic() error = %q, want it to contain %q", err.Error(), want)
				}
			}
		})
	}
}

func TestLoadTopic_MissingFile(t *testing.T) {
	_, err := LoadTopic(filepath.Join("testdata", "does-not-exist.json"))
	if err == nil {
		t.Fatal("LoadTopic() expected an error for a missing file, got nil")
	}
}

// wantContentChapter is one chapter's contract: its slug and the exact,
// in-order concept slugs design.md §7 lists under it. Positions are
// implied by list order (1-based). Shared across all track content-file
// tests (ddd, distsys, aws, go, dsa). Named distinctly from repository_test.go's
// wantChapter, which asserts DB row assembly and carries extra fields.
type wantContentChapter struct {
	slug     string
	concepts []string
}

var wantDDDChapters = []wantContentChapter{
	{"model-driven-foundations", []string{"ubiquitous-language", "model-driven-design", "knowledge-crunching", "hands-on-modelers"}},
	{"building-blocks", []string{"layered-architecture", "entities", "value-objects", "domain-services", "modules", "aggregates", "aggregate-design-rules", "factories", "repositories", "domain-events"}},
	{"supple-design-refactoring", []string{"intention-revealing-interfaces", "side-effect-free-functions", "assertions", "specification-pattern", "making-implicit-concepts-explicit", "refactoring-toward-deeper-insight"}},
	{"strategic-design", []string{"bounded-context", "context-mapping", "shared-kernel", "customer-supplier-conformist", "anticorruption-layer", "open-host-service-published-language", "core-domain-distillation", "generic-subdomains", "large-scale-structure"}},
	{"ddd-in-go", []string{"ddd-go-project-layout", "persistence-without-orm", "in-process-domain-events", "testing-the-domain-layer"}},
}

// TestLoadTopic_DDDContentFile guards content/curriculum/ddd.json itself: it
// is data that can rot (a typo'd slug, a dropped concept, a reordered
// chapter) with no compiler to catch it.
func TestLoadTopic_DDDContentFile(t *testing.T) {
	topic, err := LoadTopic(filepath.Join("..", "..", "..", "content", "curriculum", "ddd.json"))
	if err != nil {
		t.Fatalf("LoadTopic(ddd.json) unexpected error: %v", err)
	}
	if topic.Track().String() != "ddd" || topic.Slug().String() != "domain-driven-design" {
		t.Fatalf("topic = {%q %q}, want {ddd domain-driven-design}", topic.Track().String(), topic.Slug().String())
	}

	chapters := topic.Chapters()
	if len(chapters) != len(wantDDDChapters) {
		t.Fatalf("len(Chapters()) = %d, want %d", len(chapters), len(wantDDDChapters))
	}

	for i, wantCh := range wantDDDChapters {
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
