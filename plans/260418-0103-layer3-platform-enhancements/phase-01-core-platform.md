# Phase 1: Core Platform

## Overview
- **Priority:** P0
- **Effort:** 29h
- **Status:** Pending
- Nâng cấp nền tảng chung: header, favicon, Việt hoá, settings, auth, notifications, dashboard greeting, chat bubble

## DB Migrations

### 000016_settings.up.sql
```sql
CREATE TABLE settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Seed defaults
INSERT INTO settings (key, value) VALUES
    ('smtp_host', ''),
    ('smtp_port', '587'),
    ('smtp_user', ''),
    ('smtp_pass', ''),
    ('smtp_from', ''),
    ('notion_api_key', ''),
    ('notion_database_ids', '{}'),
    ('n8n_webhook_url', ''),
    ('google_oauth_client_id', ''),
    ('google_oauth_client_secret', ''),
    ('unsplash_keywords', 'vietnam,hanoi,nature');
```

### 000017_notifications.up.sql
```sql
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
```

### 000018_password_reset_tokens.up.sql
```sql
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    token VARCHAR(100) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_reset_tokens_token ON password_reset_tokens(token) WHERE used_at IS NULL;
```

## Implementation Steps

### 1. Sticky header + Favicon (1h)
- Edit `web/templates/layout.templ`: remove `overflow-x-auto` on nav, ensure `sticky top-0`
- Add `<link rel="icon" href="https://binhvuong.vn/favicon.png">` to head
- Header avatar: replace hardcoded "BV" with `<img>` from user.avatar_url or favicon fallback

### 2. Việt hoá labels (2h)
- Create `internal/handler/i18n.go` — centralized Vietnamese label maps
- Status: todo→Cần làm, in_progress→Đang làm, waiting→Chờ, review→Cần duyệt, done→Hoàn thành, cancelled→Đã huỷ
- Priority: urgent→Gấp, high→Cao, normal→Trung bình, low→Thấp
- Content status: idea→Ý tưởng, drafting→Đang viết, review→Cần duyệt, revise→Sửa lại, approved→Đã duyệt, published→Đã đăng
- Work log status: submitted→Chờ duyệt, approved→Đã duyệt, rejected→Từ chối, needs_fix→Cần sửa
- Campaign status: planning→Lên kế hoạch, running→Đang chạy, paused→Tạm dừng, ended→Kết thúc
- Update all templ files to use centralized labels

### 3. Settings table + Admin page (4h)
- Migration 000016
- `internal/db/generated/settings.sql.go` — GetSetting, SetSetting, ListSettings
- `internal/handler/admin.go` — GET/POST /admin/settings
- `web/templates/pages/admin_settings.templ` — form groups: SMTP, Notion, n8n, Google OAuth
- Route: owner-only via RequireRole middleware

### 4. SMTP + Forgot password (4h)
- Migration 000018
- `internal/email/email.go` — SendEmail function using settings SMTP config
- `internal/handler/auth.go` — add ForgotPassword, ResetPassword handlers
- `web/templates/pages/forgot_password.templ` + `reset_password.templ`
- Routes: GET/POST /auth/forgot-password, GET/POST /auth/reset-password?token=xxx
- Token: random 64-char, expires 1h

### 5. Google OAuth login (3h)
- Add `golang.org/x/oauth2` to go.mod
- `internal/handler/oauth.go` — GoogleLogin (redirect), GoogleCallback (exchange code, find/create user)
- Read client_id/secret from settings table
- Routes: GET /auth/google, GET /auth/google/callback
- Login page: add "Đăng nhập với Google" button

### 6. User avatar upload (1h)
- POST /users/:id/avatar — accept file, upload to Drive, update users.avatar_url
- Header template: show `<img src={avatar_url}>` if set, else initial letter
- Users list: show avatar thumbnails

### 7. User CRUD nâng cao (3h)
- GET /users/:id — user detail/edit page
- POST /users/:id — update full fields (full_name, role, phone, status, telegram_id)
- POST /users/:id/delete — soft delete
- Show user_company_assignments with can_view/can_edit/can_approve

### 8. Notifications (3h)
- Migration 000017
- `internal/db/generated/notifications.sql.go` — ListByUser, CountUnread, MarkRead, Create
- `internal/handler/notifications.go` — GET /notifications, POST /notifications/:id/read
- Header bell: HTMX hx-get="/notifications/count" → show unread badge
- Create notifications on: work_log approved/rejected, task assigned, content reviewed
- Notification dropdown panel (HTMX load)

### 9. Dashboard greeting + background (3h)
- `internal/handler/dashboard.go` — time-based greeting using Asia/Ho_Chi_Minh timezone
  - 5:00-10:59 → "Chào buổi sáng" + morning keywords
  - 11:00-12:59 → "Chào buổi trưa" + noon keywords  
  - 13:00-17:59 → "Chào buổi chiều" + afternoon keywords
  - 18:00-21:59 → "Chào buổi tối" + evening keywords
  - 22:00-4:59 → "Chào đêm muộn" + night keywords
- Background: `https://source.unsplash.com/1600x400/?{keywords},{time_keyword}`
- Pass greeting + image URL to dashboard template
- Hero section with background image + overlay text

### 10. Chat bubble n8n (3h)
- `web/static/js/chat-bubble.js` — floating widget
  - Button bottom-right: avatar icon https://binhvuong.vn/favicon.png
  - Click → expand chat panel (300x400px)
  - Input field + send button
  - POST to n8n webhook URL (from inline config injected by templ)
  - Body: `{user_name, message, session_id, timestamp}`
  - Display response in chat panel
  - Session: `localStorage.getItem('bvos_chat_session')` or generate UUID
- Layout template: inject chat bubble JS + config

### 11. Admin work_types CRUD (2h)
- GET /admin/work-types — list with edit/delete
- POST /admin/work-types — create new
- POST /admin/work-types/:id — update (name, slug, unit, icon, color, sort_order)
- POST /admin/work-types/:id/delete — soft delete (set active=false)
- `web/templates/pages/admin_work_types.templ`

## Files to Create
- `internal/db/migrations/000016-000018` (6 SQL files)
- `internal/db/generated/settings.sql.go`
- `internal/db/generated/notifications.sql.go`
- `internal/db/generated/password_reset_tokens.sql.go`
- `internal/handler/admin.go`
- `internal/handler/oauth.go`
- `internal/handler/notifications.go`
- `internal/handler/i18n.go`
- `internal/email/email.go`
- `web/templates/pages/admin_settings.templ`
- `web/templates/pages/admin_work_types.templ`
- `web/templates/pages/forgot_password.templ`
- `web/templates/pages/reset_password.templ`
- `web/static/js/chat-bubble.js`

## Files to Modify
- `web/templates/layout.templ` — favicon, sticky header, notification bell, chat bubble
- `web/templates/pages/login.templ` — Google OAuth button, forgot password link
- `web/templates/pages/dashboard.templ` — greeting + background image
- `cmd/server/main.go` — new routes
- `go.mod` — add golang.org/x/oauth2
- All templ files — Việt hoá labels

## Success Criteria
- [ ] Header sticky, no horizontal scroll
- [ ] Favicon = binhvuong.vn/favicon.png
- [ ] All status/priority labels in Vietnamese
- [ ] Admin Settings page saves/loads SMTP, Notion, n8n, OAuth config
- [ ] Forgot password email flow works
- [ ] Google OAuth login works
- [ ] User avatars show in header + lists
- [ ] Notification bell shows unread count
- [ ] Dashboard greeting changes by time of day
- [ ] Chat bubble sends to n8n webhook
- [ ] Work types editable in admin
