CREATE TABLE IF NOT EXISTS lesson_progress (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    lesson_id BIGINT UNSIGNED NOT NULL,
    state ENUM('locked', 'in_progress', 'passed') NOT NULL,
    first_passed_at TIMESTAMP NULL,
    last_read_at TIMESTAMP NULL,
    UNIQUE KEY uniq_lesson_progress_lesson (lesson_id),
    CONSTRAINT fk_lesson_progress_lesson FOREIGN KEY (lesson_id) REFERENCES lessons (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
