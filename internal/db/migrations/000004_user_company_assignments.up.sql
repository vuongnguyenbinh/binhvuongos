CREATE TABLE user_company_assignments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id      UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    role_in_company VARCHAR(50),
    can_view        BOOLEAN DEFAULT TRUE,
    can_edit        BOOLEAN DEFAULT TRUE,
    can_approve     BOOLEAN DEFAULT FALSE,
    start_date      DATE,
    end_date        DATE,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, company_id)
);

CREATE INDEX idx_uca_user ON user_company_assignments(user_id) WHERE end_date IS NULL;
CREATE INDEX idx_uca_company ON user_company_assignments(company_id) WHERE end_date IS NULL;

CREATE TRIGGER tr_uca_updated_at BEFORE UPDATE ON user_company_assignments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
