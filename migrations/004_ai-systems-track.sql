-- Adds 'ai-systems' to the topics.track ENUM. Appending at the end is safe:
-- MySQL stores ENUM values by their internal index, so existing rows are
-- unaffected. Keep this ordering in sync with domain.Tracks() in
-- internal/curriculum/domain/track.go.
ALTER TABLE topics MODIFY COLUMN track ENUM('ddd', 'distsys', 'aws', 'go', 'dsa', 'ddia', 'ai-systems') NOT NULL;
