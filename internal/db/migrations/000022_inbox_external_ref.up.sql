ALTER TABLE inbox_items ADD COLUMN external_ref VARCHAR(200);

CREATE UNIQUE INDEX idx_inbox_source_external_ref
    ON inbox_items(source, external_ref)
    WHERE external_ref IS NOT NULL;
