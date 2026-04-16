# Solo Expert Bình Vương OS — PostgreSQL Schema

> Schema hoàn chỉnh cho hệ thống quản lý Solo Founder.
> Stack: PostgreSQL 16 + Go + Docker + sync Notion 1h/lần.
> Mọi bảng đều có: `id UUID`, `created_at`, `updated_at`, cột sync.

---

## Tổng quan 12 bảng

| # | Bảng | Mục đích | Ghi chú |
|---|---|---|---|
| 1 | `users` | Tài khoản đăng nhập | Cả bạn + nhân sự + CTV |
| 2 | `companies` | Danh sách công ty | Xương sống |
| 3 | `user_company_assignments` | Phân quyền user ↔ company | Many-to-many |
| 4 | `inbox_items` | Hộp ghi nhanh | Telegram → đây |
| 5 | `tasks` | Công việc | Gộp projects qua trường `group_name` |
| 6 | `content` | Nội dung | Gộp pipeline + library |
| 7 | `campaigns` | Chiến dịch/nhóm output | Cho work_logs |
| 8 | `work_types` | Cấu hình loại công việc hàng ngày | Backlink, Content, Ads... |
| 9 | `work_logs` | Log output hàng ngày | **Module mới** |
| 10 | `knowledge_items` | Kho kiến thức + nguồn liệu | Gộp |
| 11 | `objectives` | Mục tiêu OKR | Nhẹ, chỉ title + quarter |
| 12 | `notion_sync_log` | Nhật ký sync | Debug khi lỗi |

---

## Quy ước chung

### Các cột phổ biến (gần như mọi bảng đều có)

```sql
id              UUID PRIMARY KEY DEFAULT gen_random_uuid()
created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
created_by      UUID REFERENCES users(id)

-- Cột phục vụ sync Notion (chỉ bảng cần sync)
notion_page_id  TEXT UNIQUE
synced_at       TIMESTAMPTZ
sync_status     TEXT DEFAULT 'pending'  -- pending|synced|error|skip
sync_error      TEXT

-- Cột soft delete (cho bảng quan trọng)
deleted_at      TIMESTAMPTZ
```

### Extension cần bật

```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";     -- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pg_trgm";      -- full-text search
CREATE EXTENSION IF NOT EXISTS "unaccent";     -- bỏ dấu tiếng Việt khi search
```

### Trigger auto-update `updated_at`

```sql
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

Áp dụng cho mọi bảng có `updated_at` (ghi ở cuối từng bảng).

---

## BẢNG 1: `users` — Tài khoản

> Cả bạn, staff, CTV, nhân sự client đều ở đây. Khác nhau ở `role`.

```sql
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Thông tin cơ bản
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   TEXT NOT NULL,
    full_name       VARCHAR(200) NOT NULL,
    
    -- Phân loại
    role            VARCHAR(30) NOT NULL DEFAULT 'staff',
    -- Values: 'owner' | 'core_staff' | 'ctv' | 'client_staff' | 'intern'
    
    -- Thông tin bổ sung
    avatar_url      TEXT,
    phone           VARCHAR(30),
    telegram_id     VARCHAR(100),              -- Dùng cho Telegram bot
    zalo_contact    VARCHAR(100),
    specialties     TEXT[],                    -- Mảng kỹ năng: ['SEO', 'Content']
    rate_note       TEXT,                      -- Mức phí/lương — free text
    
    -- Trạng thái
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    -- Values: 'active' | 'on_leave' | 'inactive'
    
    start_date      DATE,
    last_login_at   TIMESTAMPTZ,
    
    -- Notes
    internal_notes  TEXT,                      -- Chỉ bạn thấy (điểm mạnh/yếu)
    
    -- Sync
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    
    -- Timestamps
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
```

**Giải thích các trường chính:**

| Trường | Ghi chú |
|---|---|
| `role` | `owner` = bạn. `core_staff` = staff cốt lõi. `ctv` = cộng tác viên. `client_staff` = nhân sự công ty khách. Dùng để phân quyền. |
| `specialties` | Mảng Postgres — query `WHERE 'SEO' = ANY(specialties)` |
| `telegram_id` | Để xác thực khi user gửi message vào Telegram bot |
| `internal_notes` | Chỉ role `owner` thấy được (enforce ở app layer) |

---

## BẢNG 2: `companies` — Danh sách công ty (Xương sống)

```sql
CREATE TABLE companies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Thông tin cơ bản
    name            VARCHAR(200) NOT NULL,
    short_code      VARCHAR(10) UNIQUE,        -- ABC, FIN1, EDU
    slug            VARCHAR(100) UNIQUE,        -- dùng cho URL
    logo_url        TEXT,
    
    -- Phân loại
    industry        VARCHAR(50),
    -- Values: 'education'|'finance'|'food'|'tech'|'ecommerce'|'realestate'|'other'
    
    my_role         VARCHAR(30) NOT NULL,
    -- Values: 'owner'|'cofounder'|'mentor'
    
    scope           TEXT[],                    -- ['content', 'seo', 'training']
    
    -- Trạng thái
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    -- Values: 'active'|'paused'|'ended'
    
    health          VARCHAR(20) DEFAULT 'ok',
    -- Values: 'ok'|'attention'|'urgent'
    
    -- Thông tin liên hệ
    primary_contact_name    VARCHAR(200),
    primary_contact_phone   VARCHAR(30),
    primary_contact_zalo    VARCHAR(100),
    primary_contact_email   VARCHAR(255),
    
    -- Thời gian
    start_date      DATE,
    end_date        DATE,
    
    -- Ghi chú
    description     TEXT,
    internal_notes  TEXT,
    
    -- Sync
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    
    -- Timestamps
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
```

---

## BẢNG 3: `user_company_assignments` — Phân quyền

> 1 user có thể thuộc nhiều company. Bảng này quyết định ai thấy gì.

```sql
CREATE TABLE user_company_assignments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id      UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    
    -- Vai trò trong company này
    role_in_company VARCHAR(50),
    -- Ví dụ: 'content_writer'|'seo_executive'|'ads_specialist'|'manager'
    
    -- Quyền
    can_view        BOOLEAN DEFAULT TRUE,
    can_edit        BOOLEAN DEFAULT TRUE,
    can_approve     BOOLEAN DEFAULT FALSE,    -- Duyệt work_logs của người khác
    
    -- Thời gian
    start_date      DATE,
    end_date        DATE,                      -- NULL = vẫn đang làm
    
    -- Ghi chú
    notes           TEXT,
    
    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(user_id, company_id)
);

CREATE INDEX idx_uca_user ON user_company_assignments(user_id) WHERE end_date IS NULL;
CREATE INDEX idx_uca_company ON user_company_assignments(company_id) WHERE end_date IS NULL;

CREATE TRIGGER tr_uca_updated_at BEFORE UPDATE ON user_company_assignments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
```

**Query quan trọng: user X được xem task của company nào?**

```sql
SELECT c.*
FROM companies c
JOIN user_company_assignments uca ON c.id = uca.company_id
WHERE uca.user_id = $1
  AND uca.can_view = TRUE
  AND (uca.end_date IS NULL OR uca.end_date > CURRENT_DATE);
```

---

## BẢNG 4: `inbox_items` — Hộp ghi nhanh

> Cổng vào duy nhất. Telegram → đây. Triage sau.

```sql
CREATE TABLE inbox_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Nội dung
    content         TEXT NOT NULL,              -- Text/link/ghi chú
    url             TEXT,                       -- Nếu là link, tách riêng
    
    -- Phân loại sơ bộ
    source          VARCHAR(30) DEFAULT 'manual',
    -- Values: 'telegram'|'zalo'|'tiktok'|'facebook'|'web'|'email'|'manual'
    
    item_type       VARCHAR(30),
    -- Values: 'idea'|'task'|'link'|'note'|'file'|null
    
    -- Trạng thái
    status          VARCHAR(20) NOT NULL DEFAULT 'raw',
    -- Values: 'raw'|'processing'|'done'|'archived'
    
    destination     VARCHAR(30),                -- Đã quyết định chuyển đi đâu
    -- Values: 'tasks'|'content'|'knowledge'|'archive'|'delete'
    
    -- Context
    company_id      UUID REFERENCES companies(id),  -- Nếu biết
    submitted_by    UUID REFERENCES users(id),      -- Ai gửi (thường là bạn)
    
    -- Attachments
    attachments     JSONB DEFAULT '[]',          -- Mảng URL ảnh/file
    -- Ví dụ: [{"url": "...", "type": "image", "name": "screenshot.png"}]
    
    -- Metadata từ Telegram (nếu có)
    telegram_message_id  TEXT,
    telegram_chat_id     TEXT,
    
    -- Ghi chú sau triage
    triage_notes    TEXT,
    processed_at    TIMESTAMPTZ,                -- Thời điểm triage
    
    -- Liên kết sau triage (optional)
    converted_to_type   VARCHAR(30),             -- 'task'|'content'|'knowledge'
    converted_to_id     UUID,                    -- ID record đã tạo
    
    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_inbox_status ON inbox_items(status, created_at DESC);
CREATE INDEX idx_inbox_source ON inbox_items(source);
CREATE INDEX idx_inbox_created_by ON inbox_items(submitted_by, created_at DESC);

-- Auto-archive sau 7 ngày nếu vẫn là 'raw'
-- Chạy bằng cron trong app Go, không dùng pg_cron để tránh phụ thuộc

CREATE TRIGGER tr_inbox_updated_at BEFORE UPDATE ON inbox_items
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
```

**Lưu ý:** Inbox **không sync Notion** (là staging area, không cần nhân đôi).

---

## BẢNG 5: `tasks` — Công việc

> 1 bảng cho mọi task. `group_name` thay thế Projects DB.

```sql
CREATE TABLE tasks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Nội dung
    title           VARCHAR(500) NOT NULL,
    description     TEXT,
    
    -- Phân loại
    category        VARCHAR(30),
    -- Values: 'content'|'seo'|'training'|'admin'|'event'|'design'|'meeting'|'other'
    
    group_name      VARCHAR(200),               -- Thay thế Projects DB
    -- Ví dụ: "Ra mắt SP tháng 5" | "Chiến dịch SEO Q2"
    
    -- Relations
    company_id      UUID REFERENCES companies(id),
    assignee_id     UUID REFERENCES users(id),
    objective_id    UUID REFERENCES objectives(id),
    content_id      UUID,                       -- FK đến content, sẽ add sau vì forward ref
    campaign_id     UUID REFERENCES campaigns(id),
    
    -- Trạng thái
    status          VARCHAR(20) NOT NULL DEFAULT 'todo',
    -- Values: 'todo'|'in_progress'|'waiting'|'review'|'done'|'cancelled'
    
    priority        VARCHAR(20) NOT NULL DEFAULT 'normal',
    -- Values: 'urgent'|'high'|'normal'|'low'
    
    -- Thời gian
    due_date        DATE,
    due_date_end    DATE,                       -- Nếu là khoảng thời gian
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    
    -- Attachments
    attachments     JSONB DEFAULT '[]',
    
    -- Sync
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    
    -- Timestamps
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
```

---

## BẢNG 6: `content` — Nội dung (Pipeline + Library)

```sql
CREATE TABLE content (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Nội dung
    title           VARCHAR(500) NOT NULL,
    
    -- Phân loại
    content_type    VARCHAR(30) NOT NULL,
    -- Values: 'blog'|'social_post'|'video'|'reel'|'email'|'slide'|'infographic'|'podcast'
    
    platforms       TEXT[],                     -- ['facebook', 'tiktok', 'website']
    topics          TEXT[],                     -- ['seo', 'branding']
    
    -- Relations
    company_id      UUID NOT NULL REFERENCES companies(id),
    author_id       UUID NOT NULL REFERENCES users(id),
    campaign_id     UUID REFERENCES campaigns(id),
    
    -- Trạng thái
    status          VARCHAR(20) NOT NULL DEFAULT 'idea',
    -- Values: 'idea'|'drafting'|'review'|'revise'|'approved'|'published'|'killed'
    
    -- Publish info
    publish_date    DATE,
    published_url   TEXT,
    source_file_url TEXT,                       -- Link Google Docs gốc
    attachments     JSONB DEFAULT '[]',
    
    -- Performance metrics (nhân sự tự điền)
    reach           INTEGER DEFAULT 0,
    engagement      INTEGER DEFAULT 0,
    -- Engagement rate tính bằng query hoặc computed column
    
    -- Visibility (các công ty khác được xem)
    visible_to_companies UUID[] DEFAULT '{}',
    
    -- Ghi chú
    notes           TEXT,
    review_notes    TEXT,                       -- Feedback khi review
    
    -- Sync
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    
    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES users(id),
    deleted_at      TIMESTAMPTZ
);

-- Computed column cho engagement rate
ALTER TABLE content ADD COLUMN engagement_rate DECIMAL(5,2)
    GENERATED ALWAYS AS (
        CASE 
            WHEN reach > 0 THEN ROUND((engagement::DECIMAL / reach * 100), 2)
            ELSE 0
        END
    ) STORED;

CREATE INDEX idx_content_status ON content(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_content_company ON content(company_id, status);
CREATE INDEX idx_content_author ON content(author_id);
CREATE INDEX idx_content_publish_date ON content(publish_date DESC) WHERE status = 'published';
CREATE INDEX idx_content_campaign ON content(campaign_id) WHERE campaign_id IS NOT NULL;
CREATE INDEX idx_content_sync ON content(sync_status) WHERE sync_status IN ('pending', 'error');

-- Forward FK cho tasks.content_id (giờ content đã tồn tại)
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_content
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE SET NULL;

CREATE TRIGGER tr_content_updated_at BEFORE UPDATE ON content
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
```

---

## BẢNG 7: `campaigns` — Chiến dịch

```sql
CREATE TABLE campaigns (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Thông tin cơ bản
    name            VARCHAR(300) NOT NULL,
    description     TEXT,
    
    -- Relations
    company_id      UUID NOT NULL REFERENCES companies(id),
    owner_id        UUID REFERENCES users(id),  -- Người phụ trách
    
    -- Phân loại
    campaign_type   VARCHAR(50),
    -- Values: 'seo'|'product_launch'|'brand'|'content_series'|'training'|'event'|'other'
    
    -- Trạng thái
    status          VARCHAR(20) NOT NULL DEFAULT 'planning',
    -- Values: 'planning'|'running'|'paused'|'ended'|'cancelled'
    
    -- Thời gian
    start_date      DATE,
    end_date        DATE,
    
    -- Mục tiêu (linh hoạt theo từng loại work_type)
    target_json     JSONB DEFAULT '{}',
    -- Ví dụ: {"backlink": 300, "content": 50, "ads": 10}
    -- Key phải khớp với work_types.slug
    
    budget          DECIMAL(15, 2),              -- Ngân sách VND
    budget_spent    DECIMAL(15, 2) DEFAULT 0,
    
    -- Ghi chú
    notes           TEXT,
    
    -- Sync
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    
    -- Timestamps
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
```

---

## BẢNG 8: `work_types` — Cấu hình loại công việc

> Bảng config. Thêm loại mới = INSERT 1 dòng, không sửa code.

```sql
CREATE TABLE work_types (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    name            VARCHAR(100) NOT NULL,      -- "Xây backlink"
    slug            VARCHAR(50) UNIQUE NOT NULL, -- "backlink"
    unit            VARCHAR(30) NOT NULL,        -- "link" | "bài" | "campaign"
    
    -- Hiển thị
    icon            VARCHAR(10),                 -- Emoji
    color           VARCHAR(20),                 -- Hex color #4F46E5
    description     TEXT,
    
    -- Cấu hình
    active          BOOLEAN DEFAULT TRUE,
    sort_order      INTEGER DEFAULT 0,
    
    -- Metadata
    default_target_per_day INTEGER,              -- KPI mặc định/người/ngày (optional)
    
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_work_types_active ON work_types(active, sort_order);

CREATE TRIGGER tr_work_types_updated_at BEFORE UPDATE ON work_types
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Seed data ban đầu
INSERT INTO work_types (name, slug, unit, icon, color, sort_order) VALUES
    ('Xây backlink',    'backlink', 'link',     '🔗', '#4F46E5', 1),
    ('Viết bài',        'content',  'bài',      '📝', '#10B981', 2),
    ('Set up ads',      'ads',      'campaign', '📢', '#F59E0B', 3),
    ('Đăng mạng XH',    'social',   'post',     '📱', '#EC4899', 4),
    ('Làm video',       'video',    'video',    '🎥', '#EF4444', 5),
    ('Email outreach',  'email',    'email',    '📧', '#6B7280', 6);
```

---

## BẢNG 9: `work_logs` — **BẢNG CHÍNH** của module mới

> Nhân sự báo cáo output cuối ngày. Module quan trọng nhất.

```sql
CREATE TABLE work_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Thời gian làm việc
    work_date       DATE NOT NULL,              -- Ngày làm thực tế
    
    -- Relations
    user_id         UUID NOT NULL REFERENCES users(id),
    company_id      UUID NOT NULL REFERENCES companies(id),
    work_type_id    UUID NOT NULL REFERENCES work_types(id),
    campaign_id     UUID REFERENCES campaigns(id),     -- Optional
    
    -- Số lượng
    quantity        DECIMAL(10, 2) NOT NULL,    -- Số lượng nhân sự báo
    -- Dùng DECIMAL để hỗ trợ "1.5 bài" nếu cần
    
    -- Evidence
    sheet_url       TEXT,                        -- Link Google Sheet chi tiết
    evidence_url    TEXT,                        -- Link drive/folder ảnh
    screenshots     JSONB DEFAULT '[]',          -- Mảng URL ảnh upload trực tiếp
    
    -- Ghi chú
    notes           TEXT,                        -- Nhân sự note
    admin_notes     TEXT,                        -- Bạn note khi review
    
    -- Workflow
    status          VARCHAR(20) NOT NULL DEFAULT 'submitted',
    -- Values: 'submitted'|'approved'|'rejected'|'needs_fix'
    
    submitted_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_at     TIMESTAMPTZ,
    reviewed_by     UUID REFERENCES users(id),
    
    -- Sync
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    
    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- 1 người không log trùng loại cùng ngày cho cùng campaign
    CONSTRAINT uq_work_log_unique UNIQUE (user_id, work_date, work_type_id, campaign_id)
);

CREATE INDEX idx_work_logs_date ON work_logs(work_date DESC);
CREATE INDEX idx_work_logs_user_date ON work_logs(user_id, work_date DESC);
CREATE INDEX idx_work_logs_company_date ON work_logs(company_id, work_date DESC);
CREATE INDEX idx_work_logs_campaign ON work_logs(campaign_id, work_date DESC) 
    WHERE campaign_id IS NOT NULL;
CREATE INDEX idx_work_logs_work_type ON work_logs(work_type_id, work_date DESC);
CREATE INDEX idx_work_logs_status ON work_logs(status) WHERE status = 'submitted';
CREATE INDEX idx_work_logs_sync ON work_logs(sync_status) WHERE sync_status IN ('pending', 'error');

CREATE TRIGGER tr_work_logs_updated_at BEFORE UPDATE ON work_logs
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
```

### View hữu ích: Tổng hợp theo ngày/tuần/tháng

```sql
-- VIEW 1: Tổng theo tháng, breakdown work_type
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

-- VIEW 2: Tiến độ campaign vs target
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
WHERE c.target_json ? wt.slug    -- Chỉ lấy work_type có trong target
GROUP BY c.id, c.name, c.company_id, c.start_date, c.end_date, c.status,
         wt.slug, wt.name, wt.unit, c.target_json;

-- VIEW 3: Hiệu suất trung bình theo người
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
```

---

## BẢNG 10: `knowledge_items` — Kho kiến thức + Nguồn liệu

> Gộp bài giảng, SOP, template, nguyên liệu thô, curated content.

```sql
CREATE TABLE knowledge_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Nội dung
    title           VARCHAR(500) NOT NULL,
    description     TEXT,
    body            TEXT,                        -- Nội dung markdown (nếu viết trên app)
    
    -- Phân loại
    category        VARCHAR(30) NOT NULL,
    -- Values: 'lecture'|'sop'|'template'|'training'|'raw_material'|'research'|'curated'
    
    topics          TEXT[],                      -- ['seo', 'content']
    
    -- Chất lượng (dùng cho raw_material/curated)
    quality_rating  INTEGER CHECK (quality_rating BETWEEN 1 AND 3),
    -- 1⭐ | 2⭐⭐ | 3⭐⭐⭐
    
    -- Scope
    scope           VARCHAR(20) NOT NULL DEFAULT 'shared',
    -- Values: 'shared' (dùng chung) | 'company_specific' | 'personal'
    
    company_id      UUID REFERENCES companies(id),  -- Nếu company_specific
    visible_to_companies UUID[] DEFAULT '{}',        -- Chia sẻ cho cty nào
    
    -- Source
    format          VARCHAR(30),
    -- Values: 'notion_page'|'google_docs'|'pdf'|'video'|'slide'|'external_link'
    
    source_url      TEXT,
    attachments     JSONB DEFAULT '[]',
    
    -- Trạng thái
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    -- Values: 'active'|'draft'|'archived'
    
    -- Sync
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    
    -- Timestamps
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

-- Full-text search tiếng Việt
CREATE INDEX idx_knowledge_search ON knowledge_items 
    USING GIN(to_tsvector('simple', unaccent(title || ' ' || COALESCE(description, ''))));

CREATE TRIGGER tr_knowledge_updated_at BEFORE UPDATE ON knowledge_items
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
```

---

## BẢNG 11: `objectives` — Mục tiêu (OKR nhẹ)

> Giữ gọn. Key Results ghi trong `key_results_json` dạng JSONB.

```sql
CREATE TABLE objectives (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Nội dung
    title           VARCHAR(500) NOT NULL,
    description     TEXT,
    
    -- Relations
    company_id      UUID NOT NULL REFERENCES companies(id),
    owner_id        UUID REFERENCES users(id),
    
    -- Thời gian
    quarter         VARCHAR(10) NOT NULL,        -- 'Q1/2026'
    year            INTEGER,                      -- 2026 (query dễ hơn)
    
    -- Trạng thái
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    -- Values: 'active'|'achieved'|'at_risk'|'missed'
    
    progress        INTEGER DEFAULT 0,            -- 0-100
    
    -- Key Results (JSONB để linh hoạt)
    key_results     JSONB DEFAULT '[]',
    -- Format: [
    --   {"description": "...", "target": 50000, "current": 32000, "unit": "followers"},
    --   ...
    -- ]
    
    -- Ghi chú
    notes           TEXT,
    
    -- Sync
    notion_page_id  TEXT UNIQUE,
    synced_at       TIMESTAMPTZ,
    sync_status     VARCHAR(20) DEFAULT 'pending',
    sync_error      TEXT,
    
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES users(id),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_objectives_company ON objectives(company_id, quarter);
CREATE INDEX idx_objectives_status ON objectives(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_objectives_quarter ON objectives(quarter);

CREATE TRIGGER tr_objectives_updated_at BEFORE UPDATE ON objectives
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
```

---

## BẢNG 12: `notion_sync_log` — Nhật ký sync

> Debug khi sync lỗi. Giữ 30 ngày rồi auto-delete.

```sql
CREATE TABLE notion_sync_log (
    id              BIGSERIAL PRIMARY KEY,
    
    -- Thông tin sync
    table_name      VARCHAR(50) NOT NULL,       -- 'tasks', 'content', ...
    record_id       UUID NOT NULL,
    action          VARCHAR(20) NOT NULL,        -- 'create'|'update'|'delete'
    
    -- Kết quả
    status          VARCHAR(20) NOT NULL,        -- 'success'|'error'|'skipped'
    notion_page_id  TEXT,
    error_message   TEXT,
    
    -- Performance
    duration_ms     INTEGER,
    
    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sync_log_record ON notion_sync_log(table_name, record_id, created_at DESC);
CREATE INDEX idx_sync_log_status ON notion_sync_log(status, created_at DESC) 
    WHERE status = 'error';
CREATE INDEX idx_sync_log_created ON notion_sync_log(created_at);

-- Auto-delete log cũ hơn 30 ngày (chạy qua cron app Go)
-- DELETE FROM notion_sync_log WHERE created_at < NOW() - INTERVAL '30 days';
```

---

## Sơ đồ quan hệ (ERD đơn giản)

```
                    ┌──────────────┐
                    │    users     │
                    └──────┬───────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
          ▼                ▼                ▼
┌──────────────────┐  ┌─────────┐  ┌──────────────┐
│user_company_     │  │inbox_   │  │work_logs     │
│assignments       │  │items    │  │              │
└─────┬────────────┘  └─────────┘  └──────┬───────┘
      │                                    │
      ▼                                    │
┌──────────────┐                          │
│  companies   │◄─────────────────────────┤
│  (xương sống)│                          │
└──────┬───────┘                          │
       │                                  │
       │       ┌──────────┐               │
       ├──────►│campaigns │◄──────────────┤
       │       └─────┬────┘               │
       │             │                    │
       │             ▼                    │
       │        ┌─────────┐           ┌──────────┐
       ├───────►│  tasks  │           │work_types│
       │        └─────────┘           └──────────┘
       │                                  ▲
       │        ┌─────────┐                │
       ├───────►│ content │────────────────┘
       │        └─────────┘    (work_logs.work_type_id)
       │
       │        ┌─────────────┐
       ├───────►│objectives   │
       │        └─────────────┘
       │
       │        ┌──────────────┐
       └───────►│knowledge_    │
                │items         │
                └──────────────┘
```

---

## Query mẫu cho Dashboard

### Dashboard 1: Tổng quan tháng hiện tại cho bạn

```sql
-- Tổng output tháng này theo loại
SELECT 
    wt.name,
    wt.unit,
    wt.icon,
    SUM(wl.quantity) AS total
FROM work_logs wl
JOIN work_types wt ON wl.work_type_id = wt.id
WHERE wl.work_date >= DATE_TRUNC('month', CURRENT_DATE)
  AND wl.status = 'approved'
GROUP BY wt.id, wt.name, wt.unit, wt.icon, wt.sort_order
ORDER BY wt.sort_order;
```

### Dashboard 2: Tiến độ campaigns đang chạy

```sql
SELECT 
    campaign_name,
    work_type_name,
    unit,
    actual,
    target,
    progress_pct,
    CASE
        WHEN progress_pct >= 90 THEN '🟢'
        WHEN progress_pct >= 60 THEN '🟡'
        ELSE '🔴'
    END AS indicator
FROM v_campaign_progress
WHERE status = 'running'
ORDER BY campaign_name, work_type_name;
```

### Dashboard 3: Top performer tuần này

```sql
SELECT 
    u.full_name,
    wt.name AS work_type,
    SUM(wl.quantity) AS total,
    wt.unit
FROM work_logs wl
JOIN users u ON wl.user_id = u.id
JOIN work_types wt ON wl.work_type_id = wt.id
WHERE wl.work_date >= DATE_TRUNC('week', CURRENT_DATE)
  AND wl.status = 'approved'
GROUP BY u.id, u.full_name, wt.id, wt.name, wt.unit
ORDER BY total DESC
LIMIT 10;
```

### Dashboard 4: Việc cần bạn review

```sql
-- work_logs chờ duyệt
SELECT COUNT(*) AS pending_reviews FROM work_logs WHERE status = 'submitted';

-- Content cần duyệt
SELECT COUNT(*) AS content_review FROM content WHERE status = 'review';

-- Tasks quá hạn
SELECT COUNT(*) AS overdue FROM tasks 
WHERE due_date < CURRENT_DATE AND status NOT IN ('done', 'cancelled');

-- Inbox chưa xử lý
SELECT COUNT(*) AS raw_inbox FROM inbox_items WHERE status = 'raw';
```

---

## Tổng kết

| Bảng | Dòng ước tính sau 1 năm | Mục đích |
|---|---|---|
| users | 20–30 | Tài khoản |
| companies | 10–15 | Công ty |
| user_company_assignments | 30–50 | Phân quyền |
| inbox_items | 3,000+ | Ghi nhanh |
| tasks | 5,000+ | Công việc |
| content | 1,500+ | Nội dung |
| campaigns | 50–100 | Chiến dịch |
| work_types | 6–10 | Config |
| **work_logs** | **15,000+** | **Output hàng ngày** |
| knowledge_items | 500–1,000 | Tài liệu |
| objectives | 50–100 | OKR |
| notion_sync_log | 50,000+ (auto-cleanup) | Debug |

**Tổng dung lượng ước tính sau 1 năm:** 200–500MB. Postgres xử lý thừa sức.

**Hiệu suất queries:**
- Dashboard load: < 100ms (với index đúng)
- Submit work_log: < 50ms
- Sync Notion 100 records: ~60s (do rate limit 2 req/sec)

---

## Checklist triển khai

1. [ ] Tạo database `binhvuong_os` với extension cần thiết
2. [ ] Chạy migration theo thứ tự: users → companies → user_company_assignments → các bảng khác
3. [ ] Seed `work_types` với 6 loại mặc định
4. [ ] Tạo user đầu tiên (bạn) với role='owner'
5. [ ] Tạo 1 company test để verify relations
6. [ ] Test insert work_log + query summary view
7. [ ] Setup migration tool: `golang-migrate` hoặc `goose`

---

*Schema v1.0 — Tối ưu cho 20 user, dưới 10 company, sync Notion 1h/lần.*
