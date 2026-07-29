-- check_key is CHAR(64) ascii_bin, not the table's default utf8mb4_0900_ai_ci:
-- it always holds lowercase hex we produce ourselves (never user-facing
-- text needing accent/case folding), so byte-exact comparison is what a hash
-- column needs, and ascii (1 byte/char) is a QUARTER the byte width of
-- utf8mb4 (up to 4 bytes/char) for an index this will be queried through on
-- every SRS review (Q-2c).
--
-- No UNIQUE constraint on check_key: recall_attempts is append-only, many
-- rows per check_key over time.
CREATE TABLE IF NOT EXISTS recall_attempts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    lesson_id BIGINT UNSIGNED NOT NULL,
    check_key CHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    type ENUM('short_answer', 'mcq') NOT NULL,
    confidence ENUM('guessed', 'unsure', 'confident') NOT NULL,
    outcome ENUM('correct', 'incorrect') NOT NULL,
    selected_option TEXT NULL,
    graded_by ENUM('self', 'llm') NOT NULL DEFAULT 'self',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Supports Q-2c's "this check's attempts, newest first"
    -- (WHERE check_key = ? ORDER BY created_at DESC).
    INDEX idx_recall_attempts_check_key_created_at (check_key, created_at DESC),
    CONSTRAINT fk_recall_attempts_lesson FOREIGN KEY (lesson_id) REFERENCES lessons (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
