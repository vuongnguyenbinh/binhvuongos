CREATE TABLE notion_sync_log (
    id              BIGSERIAL PRIMARY KEY,
    table_name      VARCHAR(50) NOT NULL,
    record_id       UUID NOT NULL,
    action          VARCHAR(20) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    notion_page_id  TEXT,
    error_message   TEXT,
    duration_ms     INTEGER,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sync_log_record ON notion_sync_log(table_name, record_id, created_at DESC);
CREATE INDEX idx_sync_log_status ON notion_sync_log(status, created_at DESC) WHERE status = 'error';
CREATE INDEX idx_sync_log_created ON notion_sync_log(created_at);
