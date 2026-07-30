-- explanation is nullable, not NOT NULL DEFAULT '': every recall_checks row
-- written before this migration (all 356 pre-AWS-1 questions) has no
-- explanation, and the domain's own optionality rule treats "no explanation"
-- and "empty string" as the same thing (see RecallCheck.Explanation) — NULL
-- is what the writer inserts for that case (internal/curriculum/infra/
-- lessonwriter.go's insertRecallCheck), matching the options column's
-- existing "optional value -> SQL NULL" convention on the same table.
ALTER TABLE recall_checks
    ADD COLUMN explanation TEXT NULL AFTER options;
