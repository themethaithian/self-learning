-- review_cards is SM-2 scheduling state, ONE row per check_key (unlike
-- recall_attempts, which is append-only with many rows per check_key) — a
-- card is the single current position in the SRS schedule, so uniqueness on
-- check_key is the whole point of the table, enforced the same way
-- lesson_progress enforces "one row per lesson" (surrogate id PK + a
-- separate UNIQUE KEY on the natural key, not the natural key as PK).
--
-- No FK to recall_checks, for the same reason recall_attempts has none
-- (check_key is content-addressed — see
-- internal/learning/domain/checkkey.go): a card can outlive the
-- recall_checks row it was scheduled from (question text edited, or the
-- whole lesson re-imported) and keeps its schedule rather than being
-- deleted or reattached to a different question. It also means a due-cards
-- query can return a check_key that no longer resolves to a live question
-- — the read path Q-2d builds must treat that as "skip it", not an error.
--
-- No lesson_id either: resolving a due check_key back to a lesson (for
-- display) is a read-path concern that belongs to whichever endpoint Q-2d
-- designs — guessing at a denormalized column here, before that endpoint's
-- actual query shape exists, would be a shot in the dark this ticket has no
-- way to validate.
--
-- ease_factor is DECIMAL(4,2), not FLOAT/DOUBLE: MySQL's DECIMAL is exact
-- fixed-point storage, so a value read, adjusted by the SM-2 recurrence, and
-- written back hundreds of times over a card's life accumulates no
-- binary-fraction rounding error the way a float column would. 1.30 is
-- SM-2's published hard floor (Wozniak) — the CHECK constraint means even a
-- buggy caller cannot persist a lower value.
--
-- due_at defaults to CURRENT_TIMESTAMP ("due immediately") rather than a
-- NULL/never-reviewed sentinel: a brand-new, never-reviewed card should
-- surface in a "due now" query with no special case. Handoff note for
-- Q-2d: internal/learning/domain.ReviewCard represents this same "never
-- reviewed, always due" state as Go's zero-value time.Time, which cannot be
-- written to a TIMESTAMP column as-is (MySQL's TIMESTAMP range starts at
-- 1970-01-01 00:00:01 UTC) — the repository mapping must substitute
-- time.Now() when persisting a fresh card, not serialize the zero value.
CREATE TABLE IF NOT EXISTS review_cards (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    check_key CHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    ease_factor DECIMAL(4,2) UNSIGNED NOT NULL DEFAULT 2.50,
    interval_days INT UNSIGNED NOT NULL DEFAULT 0,
    repetition INT UNSIGNED NOT NULL DEFAULT 0,
    due_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_reviewed_at TIMESTAMP NULL,
    UNIQUE KEY uniq_review_cards_check_key (check_key),
    CONSTRAINT chk_review_cards_ease_factor_floor CHECK (ease_factor >= 1.30),
    -- Supports Q-2d's "which cards are due now, soonest first"
    -- (WHERE due_at <= NOW() ORDER BY due_at ASC).
    INDEX idx_review_cards_due_at (due_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- review_logs is one append-only row per SM-2 advance: which recall_attempts
-- row caused it, the quality applied, and the card's full before/after
-- state, so a bad advance is auditable and review_cards can be recomputed
-- from this history alone if it is ever corrupted — review_cards is a
-- cache of "replay these logs" answer, this table is the source of truth.
-- That is also why there is no FK to review_cards: a corrupted or deleted
-- review_cards row must not take its own audit trail down with it, and
-- rebuilding a card means inserting a fresh review_cards row from this
-- table's history, never repointing old log rows at a new id.
--
-- recall_attempt_id DOES have a FK, unlike check_key: recall_attempts.id is
-- a stable AUTO_INCREMENT surrogate key that content re-import never
-- touches (only recall_checks.id churns on re-import — see checkkey.go),
-- so referencing it directly is safe the way referencing recall_checks.id
-- is not.
--
-- Per Q-2b's finding (recall_attempts.created_at is TIMESTAMP, second
-- precision, and two rows have tied on it with no race involved): any read
-- of this table's history for one check_key must be
-- ORDER BY created_at DESC, id DESC, never created_at DESC alone.
CREATE TABLE IF NOT EXISTS review_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    check_key CHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    recall_attempt_id BIGINT UNSIGNED NOT NULL,
    quality TINYINT UNSIGNED NOT NULL,
    ease_factor_before DECIMAL(4,2) UNSIGNED NOT NULL,
    ease_factor_after DECIMAL(4,2) UNSIGNED NOT NULL,
    interval_days_before INT UNSIGNED NOT NULL,
    interval_days_after INT UNSIGNED NOT NULL,
    repetition_before INT UNSIGNED NOT NULL,
    repetition_after INT UNSIGNED NOT NULL,
    due_at_before TIMESTAMP NOT NULL,
    due_at_after TIMESTAMP NOT NULL,
    -- reviewed_at is the domain event time the SM-2 formula was applied
    -- against (Advance's reviewedAt parameter) — created_at is merely when
    -- this audit row was written, which Q-2d's advance policy may not run
    -- at the exact moment the triggering attempt happened.
    reviewed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_review_logs_quality_range CHECK (quality <= 5),
    CONSTRAINT fk_review_logs_recall_attempt FOREIGN KEY (recall_attempt_id) REFERENCES recall_attempts (id),
    INDEX idx_review_logs_check_key_created_at (check_key, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
