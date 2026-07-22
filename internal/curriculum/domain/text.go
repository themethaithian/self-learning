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
)

func validateTitle(raw string) (string, error) {
	return validateBounded(raw, maxTitleRunes, ErrInvalidTitle)
}

func validateOutline(raw string) (string, error) {
	return validateBounded(raw, maxOutlineRunes, ErrInvalidOutline)
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
