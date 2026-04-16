CREATE TABLE work_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_date       DATE NOT NULL,
    user_id         UUID NOT NULL REFERENCES users(id),
    company_id      UUID NOT NULL REFERENCES companies(id),
    work_type_id    UUID NOT NULL REFERENCES work_types(id),
    campaign_id     UUID REFERENCES campaigns(id),
    quantity        DECIMAL(10, 2) NOT NULL,
    sheet_url       TEXT,
    evidence_url    TEXT,
    screenshots     JSONB DEFAULT '[]',
    notes           TEXT,
    admin_notes     TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'submitted',
    submitted_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at     TIMESTAMPTZ,
    reviewed_by     UUID REFERENCES users(id),
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_work_log_unique UNIQUE (user_id, work_date, work_type_id, campaign_id)
);

CREATE INDEX idx_work_logs_date ON work_logs(work_date DESC);
CREATE INDEX idx_work_logs_user_date ON work_logs(user_id, work_date DESC);
CREATE INDEX idx_work_logs_company_date ON work_logs(company_id, work_date DESC);
CREATE INDEX idx_work_logs_campaign ON work_logs(campaign_id, work_date DESC) WHERE campaign_id IS NOT NULL;
CREATE INDEX idx_work_logs_work_type ON work_logs(work_type_id, work_date DESC);
CREATE INDEX idx_work_logs_status ON work_logs(status) WHERE status = 'submitted';
CREATE INDEX idx_work_logs_sync ON work_logs(sync_status) WHERE sync_status IN ('pending', 'error');

CREATE TRIGGER tr_work_logs_updated_at BEFORE UPDATE ON work_logs
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Views
CREATE OR REPLACE VIEW v_work_summary_monthly AS
SELECT
    DATE_TRUNC('month', wl.work_date)::DATE AS month,
    wl.company_id,
    c.name AS company_name,
    wl.work_type_id,
    wt.name AS work_type_name,
    wt.slug AS work_type_slug,
    wt.unit,
    SUM(wl.quantity) AS total_quantity,
    COUNT(DISTINCT wl.user_id) AS people_count,
    COUNT(*) AS log_count,
    AVG(wl.quantity) AS avg_per_log
FROM work_logs wl
JOIN work_types wt ON wl.work_type_id = wt.id
JOIN companies c ON wl.company_id = c.id
WHERE wl.status = 'approved'
GROUP BY month, wl.company_id, c.name, wl.work_type_id, wt.name, wt.slug, wt.unit;

CREATE OR REPLACE VIEW v_campaign_progress AS
SELECT
    c.id AS campaign_id,
    c.name AS campaign_name,
    c.company_id,
    c.start_date,
    c.end_date,
    c.status,
    wt.slug AS work_type_slug,
    wt.name AS work_type_name,
    wt.unit,
    COALESCE(SUM(wl.quantity), 0) AS actual,
    (c.target_json->>wt.slug)::DECIMAL AS target,
    CASE
        WHEN (c.target_json->>wt.slug)::DECIMAL > 0
        THEN ROUND(COALESCE(SUM(wl.quantity), 0) * 100.0 / (c.target_json->>wt.slug)::DECIMAL, 1)
        ELSE NULL
    END AS progress_pct
FROM campaigns c
CROSS JOIN work_types wt
LEFT JOIN work_logs wl ON wl.campaign_id = c.id
    AND wl.work_type_id = wt.id
    AND wl.status = 'approved'
WHERE c.target_json ? wt.slug
GROUP BY c.id, c.name, c.company_id, c.start_date, c.end_date, c.status,
         wt.slug, wt.name, wt.unit, c.target_json;

CREATE OR REPLACE VIEW v_user_performance AS
SELECT
    wl.user_id,
    u.full_name,
    wl.work_type_id,
    wt.name AS work_type_name,
    wt.unit,
    DATE_TRUNC('month', wl.work_date)::DATE AS month,
    AVG(wl.quantity) AS daily_avg,
    SUM(wl.quantity) AS month_total,
    COUNT(DISTINCT wl.work_date) AS active_days,
    MIN(wl.quantity) AS min_daily,
    MAX(wl.quantity) AS max_daily
FROM work_logs wl
JOIN users u ON wl.user_id = u.id
JOIN work_types wt ON wl.work_type_id = wt.id
WHERE wl.status = 'approved'
GROUP BY wl.user_id, u.full_name, wl.work_type_id, wt.name, wt.unit,
         DATE_TRUNC('month', wl.work_date);
