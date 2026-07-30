package mcqguess

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type recallCheckFileDTO struct {
	Type           string   `json:"type"`
	ExpectedAnswer string   `json:"expected_answer"`
	Options        []string `json:"options,omitempty"`
}

type lessonFileDTO struct {
	Topic        string               `json:"topic"`
	RecallChecks []recallCheckFileDTO `json:"recall_checks"`
}

// LoadQuestions walks root recursively for *.json lesson files and returns
// every recall_check whose type is "mcq" as a Question, tagged with the
// file's own "topic" field as its track. short_answer checks are skipped
// here — they have no options, so no guessing heuristic applies to them,
// and letting one into the denominator would understate every hit rate
// without anyone noticing.
//
// This does not go through curriculuminfra.LoadLesson on purpose: that
// loader decodes with DisallowUnknownFields and constructs the full domain
// Lesson aggregate, which would reject a structurally broken file (e.g. an
// expected answer missing from its own options) before this package ever
// saw it. This tool needs to detect and report that condition itself (see
// Measure), and it has no concept row to resolve a topic/concept_id pair
// against, so it reads the handful of fields it needs directly and ignores
// the rest.
func LoadQuestions(root string) ([]Question, error) {
	var questions []Question

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		qs, err := loadFile(path)
		if err != nil {
			return err
		}
		questions = append(questions, qs...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func loadFile(path string) ([]Question, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("mcqguess: read %s: %w", path, err)
	}

	var dto lessonFileDTO
	if err := json.Unmarshal(raw, &dto); err != nil {
		return nil, fmt.Errorf("mcqguess: decode %s: %w", path, err)
	}

	questions := make([]Question, 0, len(dto.RecallChecks))
	for _, rc := range dto.RecallChecks {
		if rc.Type != "mcq" {
			continue
		}
		questions = append(questions, Question{
			Track:          dto.Topic,
			Source:         path,
			Options:        rc.Options,
			ExpectedAnswer: rc.ExpectedAnswer,
		})
	}
	return questions, nil
}
