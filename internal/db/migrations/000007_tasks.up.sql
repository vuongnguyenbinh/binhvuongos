CREATE TABLE tasks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(500) NOT NULL,
    description     TEXT,
    category        VARCHAR(30),
    group_name      VARCHAR(200),
    company_id      UUID REFERENCES companies(id),
    assignee_id     UUID REFERENCES users(id),
    objective_id    UUID,
    content_id      UUID,
    campaign_id     UUID REFERENCES campaigns(id),
    status          VARCHAR(20) NOT NULL DEFAULT 'todo',
    priority        VARCHAR(20) NOT NULL DEFAULT 'normal',
    due_date        DATE,
    due_date_end    DATE,
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    attachments     JSONB DEFAULT '[]',
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES users(id),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_tasks_assignee ON tasks(assignee_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_tasks_company ON tasks(company_id, status);
CREATE INDEX idx_tasks_due_date ON tasks(due_date) WHERE status NOT IN ('done', 'cancelled');
CREATE INDEX idx_tasks_group ON tasks(group_name) WHERE group_name IS NOT NULL;
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_sync ON tasks(sync_status) WHERE sync_status IN ('pending', 'error');
CREATE INDEX idx_tasks_campaign ON tasks(campaign_id) WHERE campaign_id IS NOT NULL;

CREATE TRIGGER tr_tasks_updated_at BEFORE UPDATE ON tasks
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
