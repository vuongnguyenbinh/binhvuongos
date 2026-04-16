CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   TEXT NOT NULL,
    full_name       VARCHAR(200) NOT NULL,
    role            VARCHAR(30) NOT NULL DEFAULT 'staff',
    avatar_url      TEXT,
    phone           VARCHAR(30),
    telegram_id     VARCHAR(100),
    zalo_contact    VARCHAR(100),
    specialties     TEXT[],
    rate_note       TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    start_date      DATE,
    last_login_at   TIMESTAMPTZ,
    internal_notes  TEXT,
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE status = 'active';
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_sync ON users(sync_status) WHERE sync_status IN ('pending', 'error');

CREATE TRIGGER tr_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
