package domain

import (
	"errors"
	"testing"
)

func makeReferences(t *testing.T, n int) []Reference {
	t.Helper()
	out := make([]Reference, n)
	for i := range out {
		out[i] = mustReference(t, "Source Title", "https://example.com", "why this helps")
	}
	return out
}

func makeRecallChecks(t *testing.T, n int) []RecallCheck {
	t.Helper()
	out := make([]RecallCheck, n)
	for i := range out {
		out[i] = mustRecallCheck(t, i+1, "short_answer", "question", "answer", nil)
	}
	return out
}

func TestNewLesson(t *testing.T) {
	validRefs := validReferences(t)
	validChecks := validRecallChecks(t)

	tests := []struct {
		name       string
		slug       Slug
		version    int
		titleEn    string
		estMinutes EstMinutes
		bodyMd     string
		references []Reference
		checks     []RecallCheck
		wantErr    error
	}{
		{
			name: "valid", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "# Aggregate Root\n\nเนื้อหา", references: validRefs, checks: validChecks,
		},
		{
			name: "zero-value slug", slug: Slug{}, version: 1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: validRefs, checks: validChecks,
			wantErr: ErrInvalidSlug,
		},
		{
			name: "version zero", slug: mustSlug(t, "aggregate"), version: 0, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: validRefs, checks: validChecks,
			wantErr: ErrInvalidLessonVersion,
		},
		{
			name: "negative version", slug: mustSlug(t, "aggregate"), version: -1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: validRefs, checks: validChecks,
			wantErr: ErrInvalidLessonVersion,
		},
		{
			name: "empty title", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: validRefs, checks: validChecks,
			wantErr: ErrInvalidTitle,
		},
		{
			name: "zero-value est minutes", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "Aggregate Root",
			estMinutes: EstMinutes{}, bodyMd: "body", references: validRefs, checks: validChecks,
			wantErr: ErrInvalidEstMinutes,
		},
		{
			name: "empty body", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "", references: validRefs, checks: validChecks,
			wantErr: ErrInvalidBodyMd,
		},
		{
			name: "one reference", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: makeReferences(t, 1), checks: validChecks,
			wantErr: ErrInvalidReferenceCount,
		},
		{
			name: "five references", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: makeReferences(t, 5), checks: validChecks,
			wantErr: ErrInvalidReferenceCount,
		},
		{
			name: "zero-value reference in slice", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: []Reference{validRefs[0], {}}, checks: validChecks,
			wantErr: ErrZeroChild,
		},
		{
			name: "two recall checks", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: validRefs, checks: makeRecallChecks(t, 2),
			wantErr: ErrInvalidRecallCheckCount,
		},
		{
			name: "six recall checks", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: validRefs, checks: makeRecallChecks(t, 6),
			wantErr: ErrInvalidRecallCheckCount,
		},
		{
			name: "zero-value recall check in slice", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: validRefs,
			checks:  []RecallCheck{validChecks[0], validChecks[1], {}},
			wantErr: ErrZeroChild,
		},
		{
			name: "duplicate recall check positions", slug: mustSlug(t, "aggregate"), version: 1, titleEn: "Aggregate Root",
			estMinutes: mustEstMinutes(t, 7), bodyMd: "body", references: validRefs,
			checks: []RecallCheck{
				mustRecallCheck(t, 1, "short_answer", "q1", "a1", nil),
				mustRecallCheck(t, 1, "short_answer", "q2", "a2", nil),
				mustRecallCheck(t, 2, "short_answer", "q3", "a3", nil),
			},
			wantErr: ErrDuplicatePosition,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLesson(tt.slug, tt.version, tt.titleEn, tt.estMinutes, tt.bodyMd, tt.references, tt.checks)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewLesson() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewLesson() unexpected error: %v", err)
			}
			if got.Slug() != tt.slug {
				t.Errorf("Slug() = %v, want %v", got.Slug(), tt.slug)
			}
			if got.Version() != tt.version {
				t.Errorf("Version() = %d, want %d", got.Version(), tt.version)
			}
			if len(got.References()) != len(tt.references) {
				t.Errorf("len(References()) = %d, want %d", len(got.References()), len(tt.references))
			}
			if len(got.RecallChecks()) != len(tt.checks) {
				t.Errorf("len(RecallChecks()) = %d, want %d", len(got.RecallChecks()), len(tt.checks))
			}
			if got.IsZero() {
				t.Errorf("IsZero() = true, want false for a validly constructed Lesson")
			}
		})
	}
}

func TestNewLessonSortsRecallChecksByPosition(t *testing.T) {
	l := mustLesson(t, "aggregate", 1, "Aggregate Root", 7, "body", validReferences(t), []RecallCheck{
		mustRecallCheck(t, 3, "short_answer", "q3", "a3", nil),
		mustRecallCheck(t, 1, "short_answer", "q1", "a1", nil),
		mustRecallCheck(t, 2, "short_answer", "q2", "a2", nil),
	})

	checks := l.RecallChecks()
	if len(checks) != 3 {
		t.Fatalf("len(RecallChecks()) = %d, want 3", len(checks))
	}
	for i, want := range []int{1, 2, 3} {
		if checks[i].Position().Int() != want {
			t.Errorf("RecallChecks()[%d].Position() = %d, want %d (not sorted)", i, checks[i].Position().Int(), want)
		}
	}
}

func TestNewLessonCopiesInputSlices(t *testing.T) {
	references := validReferences(t)
	checks := validRecallChecks(t)
	l := mustLesson(t, "aggregate", 1, "Aggregate Root", 7, "body", references, checks)

	references[0] = mustReference(t, "Mutated", "https://mutated.example", "mutated why")
	checks[0] = mustRecallCheck(t, 1, "short_answer", "mutated q", "mutated a", nil)

	if l.References()[0].Title() == "Mutated" {
		t.Errorf("Lesson mutated after caller modified original references slice")
	}
	if l.RecallChecks()[0].Question() == "mutated q" {
		t.Errorf("Lesson mutated after caller modified original recall checks slice")
	}
}

func TestLessonAccessorsDefensiveCopy(t *testing.T) {
	l := mustLesson(t, "aggregate", 1, "Aggregate Root", 7, "body", validReferences(t), validRecallChecks(t))

	refs := l.References()
	refs[0] = mustReference(t, "Mutated", "https://mutated.example", "mutated why")
	if l.References()[0].Title() == "Mutated" {
		t.Errorf("References() mutated externally: got title %q", l.References()[0].Title())
	}

	checks := l.RecallChecks()
	checks[0] = mustRecallCheck(t, 1, "short_answer", "mutated q", "mutated a", nil)
	if l.RecallChecks()[0].Question() == "mutated q" {
		t.Errorf("RecallChecks() mutated externally: got question %q", l.RecallChecks()[0].Question())
	}
}

func TestLessonZeroValue(t *testing.T) {
	var l Lesson
	if !l.IsZero() {
		t.Errorf("zero-value Lesson.IsZero() = false, want true")
	}
}
