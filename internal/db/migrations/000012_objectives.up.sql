CREATE TABLE objectives (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(500) NOT NULL,
    description     TEXT,
    company_id      UUID NOT NULL REFERENCES companies(id),
    owner_id        UUID REFERENCES users(id),
    quarter         VARCHAR(10) NOT NULL,
    year            INTEGER,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    progress        INTEGER DEFAULT 0,
    key_results     JSONB DEFAULT '[]',
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

-- Forward FK for tasks.objective_id
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_objective
    FOREIGN KEY (objective_id) REFERENCES objectives(id) ON DELETE SET NULL;

CREATE INDEX idx_objectives_company ON objectives(company_id, quarter);
CREATE INDEX idx_objectives_status ON objectives(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_objectives_quarter ON objectives(quarter);

CREATE TRIGGER tr_objectives_updated_at BEFORE UPDATE ON objectives
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
