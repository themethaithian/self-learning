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
-- deleted or reattached to a different question.
--
-- No lesson_id either: resolving a due check_key back to a lesson (for
-- display) is a read-path concern that belongs to whichever endpoint Q-2d
-- designs — guessing at a denormalized column here, before that endpoint's
-- actual query shape exists, would be a shot in the dark this ticket has no
-- way to validate.
--
-- ease_factor is DECIMAL(4,2) UNSIGNED, not FLOAT/DOUBLE: MySQL's DECIMAL
-- is exact fixed-point storage, so a value read, adjusted by the SM-2
-- recurrence, and written back hundreds of times over a card's life
-- accumulates no binary-fraction rounding error the way a float column
-- would. 1.30 is SM-2's published hard floor (Wozniak) — the CHECK
-- constraint means even a buggy caller cannot persist a lower value.
-- Note the domain's EaseFactor has no matching UPPER bound, while this
-- column's UNSIGNED DECIMAL(4,2) caps at 99.99 — unreachable in practice
-- (roughly 975 consecutive good reviews from the 2.50 starting value) but
-- the two disagree and this is the only place that says so.
--
-- due_at is DATETIME, not TIMESTAMP: TIMESTAMP tops out at 2038-01-19, but
-- SM-2 intervals compound (I(n) := I(n-1)*EF), so a run of good reviews
-- reaches years-long intervals quickly and due_at is a date computed
-- FORWARD from today — proven live: inserting a 2040 date into a TIMESTAMP
-- column under this repo's sql_mode raises ERROR 1292 (22007). DATETIME
-- has no such ceiling, at the cost of 1 extra byte of storage per row (5
-- vs 4) — irrelevant next to a column whose entire purpose breaks past
-- 2038.
--
-- No DEFAULT on due_at, even though DATETIME still accepts
-- "DEFAULT CURRENT_TIMESTAMP" syntactically (dropped deliberately, not an
-- oversight): TIMESTAMP normalizes to UTC on write and converts back on
-- read using the session time zone, which is what made the app's UTC
-- writes (internal/platform/mysql/pool.go's loc=UTC DSN param) agree with
-- a TIMESTAMP column's own default automatically. DATETIME does no such
-- conversion — it stores literal wall-clock digits. A DEFAULT here would
-- be a SECOND writer, using MySQL's system time zone, and it only matches
-- the app's UTC writes because this container happens to run UTC today.
-- If the eventual VPS container's time zone differs, a default-inserted
-- due_at would be silently wrong by that offset while every other
-- due_at (written explicitly by Q-2d's own inserts) stays correct — a
-- new card would not surface in a due-now query for hours, for a column
-- whose whole purpose is "no special case". One writer, one time zone:
-- Q-2d's repository must always supply due_at itself.
--
-- That value must never be Go's zero-value time.Time either, even though
-- (unlike TIMESTAMP) DATETIME's documented range starts at year 1000 but
-- does not actually enforce it — MySQL accepts 0001-01-01 without error
-- under this repo's sql_mode, proven live. The constraint is semantic, not
-- a validation boundary MySQL happens to provide: "never reviewed" must
-- mean "due now", not "due in year 1". The zero value would satisfy
-- "<= NOW()" today, so nothing would look broken while an
-- out-of-supported-range value sits in the table with no guarantee from
-- MySQL about index ordering, time zone conversion, or driver
-- round-tripping. Same rule for due_at_before/due_at_after below.
CREATE TABLE IF NOT EXISTS review_cards (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    check_key CHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    ease_factor DECIMAL(4,2) UNSIGNED NOT NULL DEFAULT 2.50,
    interval_days INT UNSIGNED NOT NULL DEFAULT 0,
    repetition INT UNSIGNED NOT NULL DEFAULT 0,
    due_at DATETIME NOT NULL,
    -- TIMESTAMP, not DATETIME: this holds a past event time, always within
    -- TIMESTAMP's 1970-2038 range, so it keeps the UTC round-trip
    -- guarantee due_at just gave up above — do not "consistently" convert
    -- this one too.
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
    -- DATETIME, not TIMESTAMP — same 2038 ceiling as review_cards.due_at
    -- above, carrying the same forward-computed dates.
    due_at_before DATETIME NOT NULL,
    due_at_after DATETIME NOT NULL,
    -- reviewed_at is the domain event time the SM-2 formula was applied
    -- against (Advance's reviewedAt parameter) — created_at is merely when
    -- this audit row was written, which Q-2d's advance policy may not run
    -- at the exact moment the triggering attempt happened. Both are
    -- TIMESTAMP, not DATETIME, for the same reason as
    -- review_cards.last_reviewed_at: past event times never need
    -- DATETIME's extended range, so they keep TIMESTAMP's UTC round-trip
    -- guarantee.
    reviewed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_review_logs_quality_range CHECK (quality <= 5),
    CONSTRAINT fk_review_logs_recall_attempt FOREIGN KEY (recall_attempt_id) REFERENCES recall_attempts (id),
    -- One log per attempt, not just intent: the advance policy (docs/tickets/
    -- quiz.md) is exactly one SM-2 advance per recall_attempts row, so a
    -- second log for the same attempt is a bug, not a valid state — this
    -- makes double-application structurally impossible instead of merely
    -- intended, and gives Q-2d's retry path an INSERT it can safely repeat.
    UNIQUE KEY uniq_review_logs_recall_attempt (recall_attempt_id),
    INDEX idx_review_logs_check_key_created_at (check_key, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
