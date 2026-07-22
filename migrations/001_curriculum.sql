CREATE TABLE IF NOT EXISTS topics (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    track ENUM('ddd', 'distsys', 'aws', 'go', 'dsa') NOT NULL,
    slug VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    position INT NOT NULL,
    UNIQUE KEY uniq_topics_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS chapters (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    topic_id BIGINT UNSIGNED NOT NULL,
    slug VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    position INT NOT NULL,
    UNIQUE KEY uniq_chapters_topic_slug (topic_id, slug),
    CONSTRAINT fk_chapters_topic FOREIGN KEY (topic_id) REFERENCES topics (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS concepts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    chapter_id BIGINT UNSIGNED NOT NULL,
    slug VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    outline TEXT NOT NULL,
    position INT NOT NULL,
    UNIQUE KEY uniq_concepts_chapter_slug (chapter_id, slug),
    CONSTRAINT fk_concepts_chapter FOREIGN KEY (chapter_id) REFERENCES chapters (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS lessons (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    concept_id BIGINT UNSIGNED NOT NULL,
    version INT NOT NULL,
    title_en VARCHAR(255) NOT NULL,
    est_minutes INT NOT NULL,
    body_md MEDIUMTEXT NOT NULL,
    refs JSON NOT NULL,
    imported_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uniq_lessons_concept (concept_id),
    CONSTRAINT fk_lessons_concept FOREIGN KEY (concept_id) REFERENCES concepts (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS recall_checks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    lesson_id BIGINT UNSIGNED NOT NULL,
    position INT NOT NULL,
    type ENUM('short_answer', 'mcq') NOT NULL,
    question TEXT NOT NULL,
    expected_answer TEXT NOT NULL,
    options JSON NULL,
    CONSTRAINT fk_recall_checks_lesson FOREIGN KEY (lesson_id) REFERENCES lessons (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
