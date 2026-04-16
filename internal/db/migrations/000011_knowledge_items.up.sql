CREATE TABLE knowledge_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(500) NOT NULL,
    description     TEXT,
    body            TEXT,
    category        VARCHAR(30) NOT NULL,
    topics          TEXT[],
    quality_rating  INTEGER CHECK (quality_rating BETWEEN 1 AND 3),
    scope           VARCHAR(20) NOT NULL DEFAULT 'shared',
    company_id      UUID REFERENCES companies(id),
    visible_to_companies UUID[] DEFAULT '{}',
    format          VARCHAR(30),
    source_url      TEXT,
    attachments     JSONB DEFAULT '[]',
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES users(id),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_knowledge_category ON knowledge_items(category, status);
CREATE INDEX idx_knowledge_scope ON knowledge_items(scope);
CREATE INDEX idx_knowledge_company ON knowledge_items(company_id) WHERE company_id IS NOT NULL;
CREATE INDEX idx_knowledge_topics ON knowledge_items USING GIN(topics);
CREATE INDEX idx_knowledge_sync ON knowledge_items(sync_status) WHERE sync_status IN ('pending', 'error');
CREATE INDEX idx_knowledge_search ON knowledge_items
    USING GIN(to_tsvector('simple', unaccent(title || ' ' || COALESCE(description, ''))));

CREATE TRIGGER tr_knowledge_updated_at BEFORE UPDATE ON knowledge_items
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
