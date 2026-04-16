CREATE TABLE companies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(200) NOT NULL,
    short_code      VARCHAR(10) UNIQUE,
    slug            VARCHAR(100) UNIQUE,
    logo_url        TEXT,
    industry        VARCHAR(50),
    my_role         VARCHAR(30) NOT NULL,
    scope           TEXT[],
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    health          VARCHAR(20) DEFAULT 'ok',
    primary_contact_name    VARCHAR(200),
    primary_contact_phone   VARCHAR(30),
    primary_contact_zalo    VARCHAR(100),
    primary_contact_email   VARCHAR(255),
    start_date      DATE,
    end_date        DATE,
    description     TEXT,
    internal_notes  TEXT,
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES users(id),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_companies_status ON companies(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_companies_health ON companies(health) WHERE status = 'active';
CREATE INDEX idx_companies_short_code ON companies(short_code);
CREATE INDEX idx_companies_sync ON companies(sync_status) WHERE sync_status IN ('pending', 'error');

CREATE TRIGGER tr_companies_updated_at BEFORE UPDATE ON companies
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
