package infra

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

func fullRow(track, topicSlug, topicTitle string, topicPos int, chapterSlug, chapterTitle string, chapterPos int, conceptSlug, conceptTitle, conceptOutline string, conceptPos int) conceptRow {
	return conceptRow{
		topicTrack:      track,
		topicSlug:       topicSlug,
		topicTitle:      topicTitle,
		topicPosition:   topicPos,
		chapterSlug:     sql.NullString{String: chapterSlug, Valid: true},
		chapterTitle:    sql.NullString{String: chapterTitle, Valid: true},
		chapterPosition: sql.NullInt64{Int64: int64(chapterPos), Valid: true},
		conceptSlug:     sql.NullString{String: conceptSlug, Valid: true},
		conceptTitle:    sql.NullString{String: conceptTitle, Valid: true},
		conceptOutline:  sql.NullString{String: conceptOutline, Valid: true},
		conceptPosition: sql.NullInt64{Int64: int64(conceptPos), Valid: true},
	}
}

func topicOnlyRow(track, slug, title string, pos int) conceptRow {
	return conceptRow{topicTrack: track, topicSlug: slug, topicTitle: title, topicPosition: pos}
}

func chapterOnlyRow(track, topicSlug, topicTitle string, topicPos int, chapterSlug, chapterTitle string, chapterPos int) conceptRow {
	return conceptRow{
		topicTrack:      track,
		topicSlug:       topicSlug,
		topicTitle:      topicTitle,
		topicPosition:   topicPos,
		chapterSlug:     sql.NullString{String: chapterSlug, Valid: true},
		chapterTitle:    sql.NullString{String: chapterTitle, Valid: true},
		chapterPosition: sql.NullInt64{Int64: int64(chapterPos), Valid: true},
	}
}

type wantConcept struct {
	slug     string
	title    string
	position int
}

type wantChapter struct {
	slug     string
	title    string
	position int
	concepts []wantConcept
}

type wantTopic struct {
	track    string
	slug     string
	title    string
	position int
	chapters []wantChapter
}

func assertTopics(t *testing.T, got []domain.Topic, want []wantTopic) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("assembleTree() returned %d topics, want %d", len(got), len(want))
	}
	for i, w := range want {
		tp := got[i]
		if tp.Track().String() != w.track || tp.Slug().String() != w.slug || tp.Title() != w.title || tp.Position().Int() != w.position {
			t.Fatalf("topic[%d] = {%q %q %q %d}, want {%q %q %q %d}",
				i, tp.Track().String(), tp.Slug().String(), tp.Title(), tp.Position().Int(), w.track, w.slug, w.title, w.position)
		}
		chapters := tp.Chapters()
		if len(chapters) != len(w.chapters) {
			t.Fatalf("topic[%d] %q has %d chapters, want %d", i, w.slug, len(chapters), len(w.chapters))
		}
		for j, wc := range w.chapters {
			ch := chapters[j]
			if ch.Slug().String() != wc.slug || ch.Title() != wc.title || ch.Position().Int() != wc.position {
				t.Fatalf("topic[%d] chapter[%d] = {%q %q %d}, want {%q %q %d}",
					i, j, ch.Slug().String(), ch.Title(), ch.Position().Int(), wc.slug, wc.title, wc.position)
			}
			concepts := ch.Concepts()
			if len(concepts) != len(wc.concepts) {
				t.Fatalf("topic[%d] chapter[%d] %q has %d concepts, want %d", i, j, wc.slug, len(concepts), len(wc.concepts))
			}
			for k, wco := range wc.concepts {
				co := concepts[k]
				if co.Slug().String() != wco.slug || co.Title() != wco.title || co.Position().Int() != wco.position {
					t.Fatalf("topic[%d] chapter[%d] concept[%d] = {%q %q %d}, want {%q %q %d}",
						i, j, k, co.Slug().String(), co.Title(), co.Position().Int(), wco.slug, wco.title, wco.position)
				}
			}
		}
	}
}

func TestAssembleTree(t *testing.T) {
	tests := []struct {
		name string
		rows []conceptRow
		want []wantTopic
	}{
		{
			name: "empty rows",
			rows: nil,
			want: []wantTopic{},
		},
		{
			name: "one topic one chapter one concept",
			rows: []conceptRow{
				fullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1),
			},
			want: []wantTopic{
				{track: "go", slug: "go-basics", title: "Go Basics", position: 1, chapters: []wantChapter{
					{slug: "syntax", title: "Syntax", position: 1, concepts: []wantConcept{
						{slug: "variables", title: "Variables", position: 1},
					}},
				}},
			},
		},
		{
			name: "topics sharing a position are ordered by slug, not row order",
			rows: []conceptRow{
				fullRow("go", "zzz-topic", "ZZZ Topic", 1, "ch", "Ch", 1, "co", "Co", "outline", 1),
				fullRow("go", "aaa-topic", "AAA Topic", 1, "ch", "Ch", 1, "co", "Co", "outline", 1),
			},
			want: []wantTopic{
				{track: "go", slug: "aaa-topic", title: "AAA Topic", position: 1, chapters: []wantChapter{
					{slug: "ch", title: "Ch", position: 1, concepts: []wantConcept{
						{slug: "co", title: "Co", position: 1},
					}},
				}},
				{track: "go", slug: "zzz-topic", title: "ZZZ Topic", position: 1, chapters: []wantChapter{
					{slug: "ch", title: "Ch", position: 1, concepts: []wantConcept{
						{slug: "co", title: "Co", position: 1},
					}},
				}},
			},
		},
		{
			name: "multiple topics sorted by position regardless of row order",
			rows: []conceptRow{
				fullRow("aws", "aws-storage", "AWS Storage", 2, "s3", "S3", 1, "buckets", "Buckets", "outline", 1),
				fullRow("ddd", "ddd-tactical", "DDD Tactical", 1, "aggregates", "Aggregates", 1, "entities", "Entities", "outline", 1),
			},
			want: []wantTopic{
				{track: "ddd", slug: "ddd-tactical", title: "DDD Tactical", position: 1, chapters: []wantChapter{
					{slug: "aggregates", title: "Aggregates", position: 1, concepts: []wantConcept{
						{slug: "entities", title: "Entities", position: 1},
					}},
				}},
				{track: "aws", slug: "aws-storage", title: "AWS Storage", position: 2, chapters: []wantChapter{
					{slug: "s3", title: "S3", position: 1, concepts: []wantConcept{
						{slug: "buckets", title: "Buckets", position: 1},
					}},
				}},
			},
		},
		{
			name: "multiple chapters per topic sorted by position",
			rows: []conceptRow{
				fullRow("go", "go-basics", "Go Basics", 1, "concurrency", "Concurrency", 2, "goroutines", "Goroutines", "outline", 1),
				fullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1),
			},
			want: []wantTopic{
				{track: "go", slug: "go-basics", title: "Go Basics", position: 1, chapters: []wantChapter{
					{slug: "syntax", title: "Syntax", position: 1, concepts: []wantConcept{
						{slug: "variables", title: "Variables", position: 1},
					}},
					{slug: "concurrency", title: "Concurrency", position: 2, concepts: []wantConcept{
						{slug: "goroutines", title: "Goroutines", position: 1},
					}},
				}},
			},
		},
		{
			name: "multiple concepts per chapter, scrambled position order",
			rows: []conceptRow{
				fullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "loops", "Loops", "outline", 2),
				fullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1),
			},
			want: []wantTopic{
				{track: "go", slug: "go-basics", title: "Go Basics", position: 1, chapters: []wantChapter{
					{slug: "syntax", title: "Syntax", position: 1, concepts: []wantConcept{
						{slug: "variables", title: "Variables", position: 1},
						{slug: "loops", title: "Loops", position: 2},
					}},
				}},
			},
		},
		{
			name: "rows for two topics fully interleaved",
			rows: []conceptRow{
				fullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1),
				fullRow("aws", "aws-storage", "AWS Storage", 2, "s3", "S3", 1, "buckets", "Buckets", "outline", 1),
				fullRow("go", "go-basics", "Go Basics", 1, "concurrency", "Concurrency", 2, "goroutines", "Goroutines", "outline", 1),
				fullRow("aws", "aws-storage", "AWS Storage", 2, "s3", "S3", 1, "versioning", "Versioning", "outline", 2),
			},
			want: []wantTopic{
				{track: "go", slug: "go-basics", title: "Go Basics", position: 1, chapters: []wantChapter{
					{slug: "syntax", title: "Syntax", position: 1, concepts: []wantConcept{
						{slug: "variables", title: "Variables", position: 1},
					}},
					{slug: "concurrency", title: "Concurrency", position: 2, concepts: []wantConcept{
						{slug: "goroutines", title: "Goroutines", position: 1},
					}},
				}},
				{track: "aws", slug: "aws-storage", title: "AWS Storage", position: 2, chapters: []wantChapter{
					{slug: "s3", title: "S3", position: 1, concepts: []wantConcept{
						{slug: "buckets", title: "Buckets", position: 1},
						{slug: "versioning", title: "Versioning", position: 2},
					}},
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assembleTree(tt.rows)
			if err != nil {
				t.Fatalf("assembleTree() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("assembleTree() returned nil, want a non-nil (possibly empty) slice")
			}
			assertTopics(t, got, tt.want)
		})
	}
}

func TestAssembleTree_InvalidData(t *testing.T) {
	tests := []struct {
		name        string
		rows        []conceptRow
		wantContain []string
	}{
		{
			name: "invalid track names the offending topic",
			rows: []conceptRow{
				fullRow("not-a-track", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1),
			},
			wantContain: []string{"go-basics"},
		},
		{
			name: "empty topic slug fails",
			rows: []conceptRow{
				fullRow("go", "", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1),
			},
			wantContain: []string{"invalid slug"},
		},
		{
			name: "topic position zero fails, names topic",
			rows: []conceptRow{
				fullRow("go", "go-basics", "Go Basics", 0, "syntax", "Syntax", 1, "variables", "Variables", "outline", 1),
			},
			wantContain: []string{"go-basics"},
		},
		{
			name: "concept position zero fails, names topic and chapter",
			rows: []conceptRow{
				fullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "variables", "Variables", "outline", 0),
			},
			wantContain: []string{"go-basics", "syntax"},
		},
		{
			name: "concept empty slug fails, names topic and chapter",
			rows: []conceptRow{
				fullRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1, "", "Variables", "outline", 1),
			},
			wantContain: []string{"go-basics", "syntax"},
		},
		{
			name: "topic with no chapters fails, names topic",
			rows: []conceptRow{
				topicOnlyRow("go", "go-basics", "Go Basics", 1),
			},
			wantContain: []string{"go-basics", "no children"},
		},
		{
			name: "chapter with no concepts fails, names topic and chapter",
			rows: []conceptRow{
				chapterOnlyRow("go", "go-basics", "Go Basics", 1, "syntax", "Syntax", 1),
			},
			wantContain: []string{"go-basics", "syntax", "no children"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := assembleTree(tt.rows)
			if err == nil {
				t.Fatal("assembleTree() expected error, got nil")
			}
			for _, want := range tt.wantContain {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("assembleTree() error = %q, want it to contain %q", err.Error(), want)
				}
			}
		})
	}
}
