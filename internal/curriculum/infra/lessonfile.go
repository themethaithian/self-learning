package infra

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

const defaultLessonVersion = 1

type recallCheckFileDTO struct {
	Q              string   `json:"q"`
	ExpectedAnswer string   `json:"expected_answer"`
	Type           string   `json:"type"`
	Options        []string `json:"options,omitempty"`
	// Explanation is optional: absent and "" decode identically for a plain
	// string field, and the domain treats both as "no explanation" (see
	// domain.validateRecallExplanation) — every pre-AWS-1 content file omits
	// this key and must keep importing unchanged.
	Explanation string `json:"explanation,omitempty"`
}

type referenceFileDTO struct {
	Title  string `json:"title"`
	Source string `json:"source"`
	Why    string `json:"why"`
}

type lessonFileDTO struct {
	Topic        string               `json:"topic"`
	ConceptID    string               `json:"concept_id"`
	Version      int                  `json:"version"`
	TitleEn      string               `json:"title_en"`
	EstMinutes   int                  `json:"est_minutes"`
	BodyMd       string               `json:"body_md"`
	RecallChecks []recallCheckFileDTO `json:"recall_checks"`
	References   []referenceFileDTO   `json:"references"`
}

// LoadLesson reads one lesson content file (the shape documented in
// content/lessons/<topic>/<concept>.json) and builds the Lesson aggregate
// through the domain constructors, so a bad file fails the same way a bad DB
// row would. It also returns the topic slug: SaveLesson needs it to resolve
// the concept the lesson belongs to, since Lesson itself carries only the
// concept slug. A file whose folder or filename disagrees with its own
// topic/concept_id fields is rejected here rather than silently imported
// under the wrong concept.
func LoadLesson(path string) (topicSlug string, lesson domain.Lesson, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", domain.Lesson{}, fmt.Errorf("infra: load lesson file %s: %w", path, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	var dto lessonFileDTO
	if err := dec.Decode(&dto); err != nil {
		return "", domain.Lesson{}, fmt.Errorf("infra: load lesson file %s: decode: %w", path, err)
	}
	if dec.More() {
		return "", domain.Lesson{}, fmt.Errorf("infra: load lesson file %s: unexpected trailing content after JSON value", path)
	}

	wantTopic := filepath.Base(filepath.Dir(path))
	wantConcept := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if dto.Topic != wantTopic {
		return "", domain.Lesson{}, fmt.Errorf("infra: load lesson file %s: topic %q does not match folder %q", path, dto.Topic, wantTopic)
	}
	if dto.ConceptID != wantConcept {
		return "", domain.Lesson{}, fmt.Errorf("infra: load lesson file %s: concept_id %q does not match filename %q", path, dto.ConceptID, wantConcept)
	}

	l, err := dto.toDomain()
	if err != nil {
		return "", domain.Lesson{}, fmt.Errorf("infra: load lesson file %s: %w", path, err)
	}
	return dto.Topic, l, nil
}

func (d lessonFileDTO) toDomain() (domain.Lesson, error) {
	slug, err := domain.NewSlug(d.ConceptID)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("concept_id: %w", err)
	}

	version := d.Version
	if version == 0 {
		version = defaultLessonVersion
	}

	estMinutes, err := domain.NewEstMinutes(d.EstMinutes)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("lesson %q: est_minutes: %w", d.ConceptID, err)
	}

	references := make([]domain.Reference, 0, len(d.References))
	for i, r := range d.References {
		ref, err := domain.NewReference(r.Title, r.Source, r.Why)
		if err != nil {
			return domain.Lesson{}, fmt.Errorf("lesson %q: reference %d: %w", d.ConceptID, i, err)
		}
		references = append(references, ref)
	}

	// position is the recall check's 1-based index in the array, never a
	// file field: the content contract has no "position" key, so two array
	// entries can never collide on position the way two chapters' explicit
	// "position" ints (T7's content format) can.
	recallChecks := make([]domain.RecallCheck, 0, len(d.RecallChecks))
	for i, rc := range d.RecallChecks {
		kind, err := domain.NewRecallKind(rc.Type)
		if err != nil {
			return domain.Lesson{}, fmt.Errorf("lesson %q: recall check %d: %w", d.ConceptID, i, err)
		}
		position, err := domain.NewPosition(i + 1)
		if err != nil {
			return domain.Lesson{}, fmt.Errorf("lesson %q: recall check %d: position: %w", d.ConceptID, i, err)
		}
		check, err := domain.NewRecallCheck(position, kind, rc.Q, rc.ExpectedAnswer, rc.Options, rc.Explanation)
		if err != nil {
			return domain.Lesson{}, fmt.Errorf("lesson %q: recall check %d: %w", d.ConceptID, i, err)
		}
		recallChecks = append(recallChecks, check)
	}

	return domain.NewLesson(slug, version, d.TitleEn, estMinutes, d.BodyMd, references, recallChecks)
}
