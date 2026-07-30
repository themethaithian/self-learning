-- explanation is nullable, not NOT NULL DEFAULT '': every recall_checks row
-- written before this migration (all 586 pre-AWS-S1 rows, mcq and
-- short_answer alike) has no explanation, and the domain's own optionality
-- rule treats "no explanation" and "empty string" as the same thing (see
-- RecallCheck.Explanation) — NULL is what the writer inserts for that case
-- (internal/curriculum/infra/lessonwriter.go's insertRecallCheck), matching
-- the options column's existing "optional value -> SQL NULL" convention on
-- the same table.
--
-- The ADD COLUMN itself is guarded, unlike every CREATE TABLE IF NOT EXISTS
-- elsewhere in this directory (see TestAllMigrationsAlterTableAddColumnIsGuarded):
-- MySQL 8 has no `ADD COLUMN IF NOT EXISTS`, and ALTER TABLE is
-- non-transactional DDL that implicitly commits before Migrate's own INSERT
-- into schema_migrations runs (see migrate.go's doc comment) — if the
-- process crashes in exactly that gap, a bare ALTER re-run on the next
-- start fails with "Duplicate column name 'explanation'" (error 1060)
-- forever, and this API's `restart: unless-stopped` turns that into a
-- permanently dead container until someone SSHes in and fixes
-- schema_migrations by hand. The guard below makes the ADD COLUMN a no-op
-- once the column already exists, so a restart after that crash converges
-- instead of looping.
SET @has_explanation_col = (
    SELECT COUNT(*) FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'recall_checks' AND COLUMN_NAME = 'explanation'
);
SET @add_explanation_ddl = IF(@has_explanation_col = 0,
    'ALTER TABLE recall_checks ADD COLUMN explanation TEXT NULL AFTER options',
    'DO 0');
PREPARE add_explanation_stmt FROM @add_explanation_ddl;
EXECUTE add_explanation_stmt;
DEALLOCATE PREPARE add_explanation_stmt;
