CREATE TABLE inbox_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content         TEXT NOT NULL,
    url             TEXT,
    source          VARCHAR(30) DEFAULT 'manual',
    item_type       VARCHAR(30),
    status          VARCHAR(20) NOT NULL DEFAULT 'raw',
    destination     VARCHAR(30),
    company_id      UUID REFERENCES companies(id),
    submitted_by    UUID REFERENCES users(id),
    attachments     JSONB DEFAULT '[]',
    telegram_message_id  TEXT,
    telegram_chat_id     TEXT,
    triage_notes    TEXT,
    processed_at    TIMESTAMPTZ,
    converted_to_type   VARCHAR(30),
    converted_to_id     UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_inbox_status ON inbox_items(status, created_at DESC);
CREATE INDEX idx_inbox_source ON inbox_items(source);
CREATE INDEX idx_inbox_created_by ON inbox_items(submitted_by, created_at DESC);

CREATE TRIGGER tr_inbox_updated_at BEFORE UPDATE ON inbox_items
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
