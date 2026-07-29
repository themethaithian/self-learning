// Package domain models the learning bounded context: LessonProgress and the
// ChunkState lifecycle it moves through, the pure Gate function that decides
// which lessons in a chapter are currently unlocked, and RecallAttempt, one
// self-graded (today) recall-check attempt keyed by content-addressed
// CheckKey.
package domain
