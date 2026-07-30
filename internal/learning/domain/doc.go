// Package domain models the learning bounded context: LessonProgress and the
// ChunkState lifecycle it moves through, the pure Gate function that decides
// which lessons in a chapter are currently unlocked, RecallAttempt, one
// self-graded (today) recall-check attempt keyed by content-addressed
// CheckKey, and ReviewCard, the SM-2 scheduling state (also keyed by
// CheckKey) that Advance moves forward one review at a time.
package domain
