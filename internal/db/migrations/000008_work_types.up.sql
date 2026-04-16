CREATE TABLE work_types (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    slug            VARCHAR(50) UNIQUE NOT NULL,
    unit            VARCHAR(30) NOT NULL,
    icon            VARCHAR(10),
    color           VARCHAR(20),
    description     TEXT,
    active          BOOLEAN DEFAULT TRUE,
    sort_order      INTEGER DEFAULT 0,
    default_target_per_day INTEGER,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_work_types_active ON work_types(active, sort_order);

CREATE TRIGGER tr_work_types_updated_at BEFORE UPDATE ON work_types
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

INSERT INTO work_types (name, slug, unit, icon, color, sort_order) VALUES
    ('Xây backlink',    'backlink', 'link',     '🔗', '#4F46E5', 1),
    ('Viết bài',        'content',  'bài',      '📝', '#10B981', 2),
    ('Set up ads',      'ads',      'campaign', '📢', '#F59E0B', 3),
    ('Đăng mạng XH',   'social',   'post',     '📱', '#EC4899', 4),
    ('Làm video',       'video',    'video',    '🎥', '#EF4444', 5),
    ('Email outreach',  'email',    'email',    '📧', '#6B7280', 6);
