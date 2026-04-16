CREATE TABLE campaigns (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(300) NOT NULL,
    description     TEXT,
    company_id      UUID NOT NULL REFERENCES companies(id),
    owner_id        UUID REFERENCES users(id),
    campaign_type   VARCHAR(50),
    status          VARCHAR(20) NOT NULL DEFAULT 'planning',
    start_date      DATE,
    end_date        DATE,
    target_json     JSONB DEFAULT '{}',
    budget          DECIMAL(15, 2),
    budget_spent    DECIMAL(15, 2) DEFAULT 0,
    notes           TEXT,
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES users(id),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_campaigns_company ON campaigns(company_id, status);
CREATE INDEX idx_campaigns_status ON campaigns(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_campaigns_sync ON campaigns(sync_status) WHERE sync_status IN ('pending', 'error');

CREATE TRIGGER tr_campaigns_updated_at BEFORE UPDATE ON campaigns
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
