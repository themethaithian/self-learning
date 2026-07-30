package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Limits are rune-counted, not byte-counted — Thai text is multi-byte in UTF-8.
const (
	maxTitleRunes   = 255
	maxOutlineRunes = 4000

	maxBodyMdRunes = 50000 // generous ceiling for a mermaid-bearing 5-10 min chunk

	maxReferenceTitleRunes  = 255
	maxReferenceSourceRunes = 255
	maxReferenceWhyRunes    = 500

	maxRecallQuestionRunes = 1000
	maxRecallAnswerRunes   = 2000
	maxRecallOptionRunes   = 255

	// maxRecallExplanationRunes matches maxOutlineRunes: an explanation is
	// structured prose (why the answer is right, why each distractor is
	// wrong, a reusable decision rule — see docs/tickets/aws-cert.md), the
	// same order of magnitude as a concept outline, and well short of a
	// full body_md chunk.
	maxRecallExplanationRunes = 4000
)

func validateTitle(raw string) (string, error) {
	return validateBounded(raw, maxTitleRunes, ErrInvalidTitle)
}

func validateOutline(raw string) (string, error) {
	return validateBounded(raw, maxOutlineRunes, ErrInvalidOutline)
}

func validateBodyMd(raw string) (string, error) {
	return validateBounded(raw, maxBodyMdRunes, ErrInvalidBodyMd)
}

func validateReferenceField(raw string, maxRunes int) (string, error) {
	return validateBounded(raw, maxRunes, ErrInvalidReference)
}

func validateRecallQuestion(raw string) (string, error) {
	return validateBounded(raw, maxRecallQuestionRunes, ErrInvalidRecallQuestion)
}

func validateRecallAnswer(raw string) (string, error) {
	return validateBounded(raw, maxRecallAnswerRunes, ErrInvalidRecallAnswer)
}

func validateRecallOption(raw string) (string, error) {
	return validateBounded(raw, maxRecallOptionRunes, ErrInvalidRecallOptions)
}

// validateRecallExplanation trims raw and, unlike validateBounded, treats an
// empty result as "no explanation" rather than an error — explanation is
// optional (118 pre-AWS-S1 lessons have none), and a plain string field can't
// distinguish an omitted JSON key from an explicit "" on decode anyway, so
// the two must mean the same thing here.
func validateRecallExplanation(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}
	if n := utf8.RuneCountInString(trimmed); n > maxRecallExplanationRunes {
		return "", fmt.Errorf("%d runes exceeds max %d: %w", n, maxRecallExplanationRunes, ErrInvalidRecallExplanation)
	}
	return trimmed, nil
}

func validateBounded(raw string, maxRunes int, sentinel error) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("empty: %w", sentinel)
	}
	if n := utf8.RuneCountInString(trimmed); n > maxRunes {
		return "", fmt.Errorf("%d runes exceeds max %d: %w", n, maxRunes, sentinel)
	}
	return trimmed, nil
}
