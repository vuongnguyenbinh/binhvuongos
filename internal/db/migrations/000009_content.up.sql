CREATE TABLE content (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(500) NOT NULL,
    content_type    VARCHAR(30) NOT NULL,
    platforms       TEXT[],
    topics          TEXT[],
    company_id      UUID NOT NULL REFERENCES companies(id),
    author_id       UUID NOT NULL REFERENCES users(id),
    campaign_id     UUID REFERENCES campaigns(id),
    status          VARCHAR(20) NOT NULL DEFAULT 'idea',
    publish_date    DATE,
    published_url   TEXT,
    source_file_url TEXT,
    attachments     JSONB DEFAULT '[]',
    reach           INTEGER DEFAULT 0,
    engagement      INTEGER DEFAULT 0,
    visible_to_companies UUID[] DEFAULT '{}',
    notes           TEXT,
    review_notes    TEXT,
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES users(id),
    deleted_at      TIMESTAMPTZ
);

ALTER TABLE content ADD COLUMN engagement_rate DECIMAL(5,2)
    GENERATED ALWAYS AS (
        CASE
            WHEN reach > 0 THEN ROUND((engagement::DECIMAL / reach * 100), 2)
            ELSE 0
        END
    ) STORED;

-- Forward FK for tasks.content_id
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_content
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE SET NULL;

-- Forward FK for tasks.objective_id (objectives not yet created, will add in 000012)

CREATE INDEX idx_content_status ON content(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_content_company ON content(company_id, status);
CREATE INDEX idx_content_author ON content(author_id);
CREATE INDEX idx_content_publish_date ON content(publish_date DESC) WHERE status = 'published';
CREATE INDEX idx_content_campaign ON content(campaign_id) WHERE campaign_id IS NOT NULL;
CREATE INDEX idx_content_sync ON content(sync_status) WHERE sync_status IN ('pending', 'error');

CREATE TRIGGER tr_content_updated_at BEFORE UPDATE ON content
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
