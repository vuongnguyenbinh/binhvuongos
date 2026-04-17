# Brainstorm Report — Layer 3: Platform Enhancements

## Bối cảnh

Bình Vương OS Layer 2 production-ready: 101 handlers, 14 DB tables, 12 trang, JWT auth, Google Drive upload, HTMX, pagination, search. Cần nâng cấp UX/features theo yêu cầu mới.

## Phân tích yêu cầu — 5 Subsystems

### Subsystem 1: Core Platform (Ưu tiên 1)
Ảnh hưởng toàn bộ trang — làm trước.

| Feature | Effort | Approach |
|---------|--------|----------|
| Sticky header (bỏ scroll) | 1h | CSS `position: sticky` remove `overflow-x-auto` |
| Favicon + admin avatar → binhvuong.vn/favicon.png | 30m | Layout template update |
| Việt hoá status/priority labels | 2h | Rà soát helper functions, thêm mapping |
| Admin Settings page (SMTP, Notion, n8n, work_types) | 4h | New `settings` table + admin UI |
| SMTP config + forgot password + email verify | 4h | Go `net/smtp` + token table |
| Google OAuth login | 3h | `golang.org/x/oauth2` + Google provider |
| User avatar upload (Drive) | 1h | Reuse existing Drive upload |
| Notification bell — real notifications | 3h | `notifications` table + unread count |
| Dashboard greeting theo giờ + Unsplash background | 3h | Time-based greeting + Unsplash API |
| Chat bubble (n8n webhook) | 3h | JS floating widget + fetch to webhook URL |
| User CRUD nâng cao (edit/delete + permissions) | 3h | Update users handler + form |
| Admin editable categories/types | 2h | CRUD work_types + content_types |

**Subtotal: ~29h**

### Subsystem 2: Comments System (Ưu tiên 2)
1 bảng polymorphic dùng chung 7+ modules.

| Feature | Effort |
|---------|--------|
| Migration: `comments` table | 30m |
| Comment handler + HTMX component | 2h |
| Wire vào: tasks, content, work_logs, knowledge, campaigns, bookmarks, companies | 3h |

**Subtotal: ~5.5h**

### Subsystem 3: Enhanced CRUD (Ưu tiên 3)
Nâng cấp filter, hiển thị, detail cho từng module.

| Module | Key Features | Effort |
|--------|-------------|--------|
| Inbox | Bulk select, inline triage, edit all fields, archive | 3h |
| WorkLogs | User column, date range filter, user/type filter, chart (Chart.js) | 4h |
| Tasks | User name, table view, date filter, quantity/unit, Trello-like | 4h |
| Content | Tag filter, status click filter, user column, tags column | 3h |
| Companies | Task count, completion rate, logo upload, contact/entity | 3h |
| Campaigns | Date/company/type filter, task count, completion, company | 3h |
| Knowledge | Tag filter, file upload | 2h |
| Bookmarks | Date/tag/domain filter | 2h |

**Subtotal: ~24h**

### Subsystem 4: Rich Text + File Attachments (Ưu tiên 4)
Content body cần editor, multi-file upload.

| Feature | Effort | Approach |
|---------|--------|----------|
| SimpleMDE editor integration | 2h | CDN script + textarea enhancement |
| Content detail: render Markdown | 1h | goldmark Go lib |
| Dashboard quick notes (rich text) | 2h | SimpleMDE + `user_notes` table |
| Multi-file upload per entity | 2h | Extend existing Drive upload |

**Subtotal: ~7h**

### Subsystem 5: Integrations (Ưu tiên 5)

| Feature | Effort |
|---------|--------|
| Notion sync worker thật (mapper + cron) | 6h |
| Notion sync button on UI | 1h |

**Subtotal: ~7h**

## Tổng effort ước tính: ~72h

## Thiết kế kỹ thuật

### DB Migrations mới cần

```sql
-- comments (polymorphic)
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(30) NOT NULL, -- 'task','content','work_log'...
    entity_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_comments_entity ON comments(entity_type, entity_id, created_at DESC);

-- notifications
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(300) NOT NULL,
    body TEXT,
    link VARCHAR(500),
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notifications_user ON notifications(user_id, read_at NULLS FIRST, created_at DESC);

-- settings (key-value for admin config)
CREATE TABLE settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- user_notes (dashboard sticky notes)
CREATE TABLE user_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- password_reset_tokens
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    token VARCHAR(100) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Thư viện JS mới (CDN)
- **SimpleMDE** — Markdown editor cho content body + notes
- **Chart.js** — Biểu đồ work logs tháng
- **Unsplash Source** — Background ảnh theo giờ (no API key needed)

### Approach cho Chat Bubble (n8n)
- JS floating button + chat panel
- Config webhook URL qua Admin Settings
- POST to n8n webhook: `{user, message, session_id}`
- n8n response hiển thị trong chat panel
- Session persist via localStorage

## Điểm mù đã phát hiện

1. **Không có activity log** — Cần audit trail khi nhiều user edit
2. **Không có user avatar field** — Schema users đã có `avatar_url` nhưng chưa dùng
3. **Content body trống** — Schema có `body TEXT` trong knowledge nhưng không có trong content → cần migration thêm `body` column cho content
4. **Không có tags table** — Nhiều module dùng `TEXT[]` cho tags, nhưng không có autocomplete/quản lý tags tập trung

## Đề xuất bổ sung

1. **Activity log table** — Track ai sửa gì, khi nào → audit trail
2. **Tags autocomplete** — Aggregate unique tags across modules cho dropdown suggest
3. **Backup cron** — pg_dump daily → Google Drive
4. **Health check endpoint** — `/api/health` cho monitoring

## Thứ tự triển khai đề xuất

```
Phase 1: Core Platform (Subsystem 1) — 29h
  → Deploy + Test

Phase 2: Comments + Enhanced CRUD (Subsystem 2+3) — 29.5h
  → Deploy + Test

Phase 3: Rich Text + Integrations (Subsystem 4+5) — 14h
  → Deploy + Test
```
