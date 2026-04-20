DROP INDEX IF EXISTS idx_inbox_source_external_ref;
ALTER TABLE inbox_items DROP COLUMN IF EXISTS external_ref;
