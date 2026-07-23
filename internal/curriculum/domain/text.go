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
