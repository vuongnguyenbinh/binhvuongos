# Journal — 2026-04-20 → 2026-04-22: Major Platform Sprint

Comprehensive session journal so any future restart can rebuild context.

## Overview

3 days, **18 commits**, **6 plans completed**, ~2,500 LOC.

Built: unified inbox webhook + n8n bridges (Telegram, Zalo), Triage convert flow, User CRUD with role consolidation, Google OAuth login, self-profile/avatar, header admin dropdown, company management (archive/logo/filter/deadline), permission matrix doc, Drive proxy endpoint.

## Production infra (locked)

| Component | Value |
|---|---|
| App | `https://os.binhvuong.vn` (Cloudflare Tunnel → app server `127.0.0.1:3000`) |
| App server | `103.97.125.186` root / `<REDACTED — see local notes>` (Docker, image `binhvuongos-app`) |
| Database | `103.97.125.131` PostgreSQL, db `binhvuong_os`, user `bvos` / `<REDACTED>` |
| n8n bridge | `https://auto.binhvuong.vn` (workflow automation) |
| API key (master) | `<REDACTED — see server .env: API_KEY>` |
| Cloudflare tunnel ID | `a4048ebd-132d-4135-a1d1-76998f8ee4a7` (name `binhvuong-os`) |
| GitHub | `https://github.com/vuongnguyenbinh/binhvuongos` (branch `main`) |

## Stack

- Go 1.22 + Fiber v2 + templ + Tailwind v3 + HTMX
- Hand-written DB layer (NOT sqlc) — `internal/db/generated/*.go` is manually maintained
- pgx/v5 + pgxpool
- Dockerfile multi-stage; `templ generate` runs on every build
- Entrypoint auto-runs migrations on startup

## Migrations applied (sequential)

| # | Migration | Purpose |
|---|---|---|
| 22 | `inbox_external_ref` | unique `(source, external_ref)` for webhook idempotency |
| 23 | `role_consolidation` | `core_staff→manager`, `ctv→staff`, CHECK constraint `owner|manager|staff` |
| 24 | `notifications_dedup` | `ref_type/ref_id/notif_date` + unique partial index |

## Plans completed

| Plan dir | Outcome |
|---|---|
| `260420-2131-unified-inbox-webhook` | POST `/api/v1/inbox` (JSON + multipart), atomic upsert idempotent, Drive upload, deleted broken Telegram handler, docs |
| `260421-0120-triage-convert-user-crud` | Fixed broken triage (now real pgx tx INSERT into target + UPDATE inbox), 3 HTMX modals (Task/Content/Knowledge), full user CRUD with perm matrix, password reset SMTP-ready |
| `260422-0907-google-login-profile-avatar` | Google OAuth login (admin-gated email whitelist), self-profile edit, avatar upload (Drive), `/users/:id` admin detail page, header avatar dropdown |
| `260422-1141-company-enhancements-permissions` | Archive/unarchive companies, logo upload, filter tabs (active/all/archived), deadline badge + dashboard widget + notifier goroutine, `docs/permissions.md` |
| (inline ad-hoc) | Drive proxy `/drive/:file_id` — stream files through app auth; logos render as `<img>` via internal path |
| (inline ad-hoc) | Admin dropdown menu in avatar (Profile/Notifications/Quản trị section/Logout), mobile bottom-sheet |

## Critical bug fixes

| Fix | Root cause | Solution |
|---|---|---|
| 431 "Request Header Fields Too Large" on triage/comments | Fiber default `ReadBufferSize=4KB` overflowed by JWT + CF + HTMX cookies | Set `ReadBufferSize: 32*1024`, `BodyLimit: 60MB` in `fiber.Config` |
| Triage marked done but no target row created | Old handler only `UPDATE inbox_items SET status='done'`, never INSERT into `tasks/content/knowledge_items` | New `ConvertInbox` handler with pgx transaction (`/inbox/:id/convert?target=...`) |
| `($1 \|\| ' days')::interval` failed pgx encoding | pgx can't encode int32 into text-typed placeholder | Use `$1::int * INTERVAL '1 day'` |
| OAuth `unauthorized_client` on Drive | refresh_token bound to old OAuth app `429460581867-...` but `.env` had new app `294942838735-...` | User regenerated refresh_token via OAuth Playground with new app credentials |
| Drive WebViewLink can't render as `<img>` | URL is HTML preview page, not raw image | Added `/drive/:file_id` proxy that streams Drive content through our auth |
| JWT cookies invalid after role migration | Legacy `role=core_staff` in JWT doesn't match new RequireRole | Rotated `JWT_SECRET` on server — forced all sessions to re-login |

## Architecture decisions (locked)

### Auth model
- 3 roles: `owner` / `manager` / `staff` (legacy core_staff/ctv migrated)
- `RequireRole("owner","manager")` for admin routes
- `CanManageUser(actor, target)` — owner manages all, manager manages staff only
- `AllowedTargetRoles(actor)` — server-side whitelist (never trust form role)
- `AuthRequired` rejects `status != 'active'` (soft-delete takes immediate effect)
- API key: `subtle.ConstantTimeCompare` (timing-safe)

### Webhook idempotency
- `POST /api/v1/inbox` accepts JSON + multipart
- Atomic upsert: `INSERT ... ON CONFLICT (source, external_ref) WHERE external_ref IS NOT NULL DO NOTHING`
- `submitted_by` always = owner user (resolved on startup via `OWNER_EMAIL` env, fail-fast)

### Google OAuth
- Single app `294942838735-d6bqd94e1670b14p41aiqs2l7j5ofqqt.apps.googleusercontent.com` for both Login + Drive
- Login: `access_type=online`, scopes `openid email profile`
- Drive: `access_type=offline` with stored `GOOGLE_REFRESH_TOKEN`
- CSRF state cookie `gstate` (10min TTL, HttpOnly)
- Login flow gates by email whitelist (must exist in `users` + `status=active`)
- Auto-populate `avatar_url` from Google picture on first login if null

### Drive proxy pattern
- Logo + avatar storage: `/drive/<fileID>` (NOT Drive WebViewLink)
- `GET /drive/:file_id` (auth-required, file_id regex whitelist)
- Streams via `drive.DownloadFile()` + `Cache-Control: private, max-age=3600`
- Inbox webhook attachments still store Drive direct URL (different use case)

### n8n bridge architecture
- Go app = 1 endpoint (`/api/v1/inbox`), no platform-specific code
- n8n at `auto.binhvuong.vn` hosts platform workflows:
  - `docs/n8n-flows/telegram-to-inbox.json`
  - `docs/n8n-flows/zalo-bot-to-inbox.json` (7 nodes, auto-download image → Drive)
- Idempotency via `external_ref = tg:<chat>:<msg_id>` / `zalo:<chat>:<msg_id>`

### Templ files
- Gitignored `web/templates/**/*_templ.go` (Dockerfile regenerates on build)
- Local dev must run `templ generate` before `go build`

### Deadline notifier (in-process)
- Background goroutine in `internal/handler/deadline_notifier.go`
- Sleep 30s after boot, then run every 24h
- Idempotent via `(user_id, ref_type, ref_id, notif_date)` unique partial index
- Panic-safe with `defer recover()`
- Owner always notified; assignees per `user_company_assignments`

## Env vars on server `.env` (final state)

```
DATABASE_URL=postgres://bvos:<REDACTED>@103.97.125.131:5432/binhvuong_os?sslmode=disable
JWT_SECRET=<rotated 4/22 12:50 ICT — 96 hex chars>
API_KEY=905cd5982517fec4215c6f91450115e96a6e4093f87816a032db47ce6d62ebdd
PORT=3000
OWNER_EMAIL=vuongnguyenbinh@gmail.com
GOOGLE_CLIENT_ID=294942838735-d6bqd94e1670b14p41aiqs2l7j5ofqqt.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-<REDACTED>  ← SHOULD ROTATE (exposed in chat)
GOOGLE_REFRESH_TOKEN=<REDACTED — see server .env>
GOOGLE_DRIVE_FOLDER_ID=<unchanged>
GOOGLE_REDIRECT_URI=https://os.binhvuong.vn/auth/google/callback  ← default fallback in code
```

`.env` permissions: `600` (both server + local).

## Critical files (mental model)

### Backend Go
- `cmd/server/main.go` — routes + `StartDeadlineNotifier(queries)`
- `internal/middleware/auth.go` — JWT + status check
- `internal/middleware/role.go` — `RequireRole`
- `internal/middleware/user_perm.go` — `CanManageUser`, `AllowedTargetRoles`, `IsAllowedRole`
- `internal/middleware/api_key.go` — constant-time compare
- `internal/oauth/google.go` — OAuth client (AuthorizeURL/ExchangeCode/FetchUserinfo)
- `internal/handler/google_auth.go` — Login flow with CSRF state
- `internal/handler/profile_edit.go` — Self-edit name/phone + avatar upload
- `internal/handler/user_crud.go` — Edit/SoftDelete/UserDetail
- `internal/handler/password_reset.go` — Token + public reset
- `internal/handler/inbox_convert.go` — Transactional triage convert
- `internal/handler/inbox_webhook_helpers.go` — Webhook validation
- `internal/handler/api_handlers.go::APICreateInbox` — Unified webhook
- `internal/handler/deadline_notifier.go` — Background goroutine
- `internal/handler/deadline_badge.go` — Badge helper
- `internal/handler/drive_proxy.go` — `/drive/:file_id` stream
- `internal/drive/drive.go` — UploadFile + DownloadFile

### DB layer
- `internal/db/generated/db.go` — pool wrapper, `Pool()` exposed for ad-hoc tx
- `internal/db/generated/companies.sql.go` — `companyColsAliased(alias)`, ListCompaniesDueSoon[ForUser]
- `internal/db/generated/users.sql.go` — UpdateUser/SoftDeleteUser/UpdateOwnProfile/UpdateUserAvatar/ListUsersByRole
- `internal/db/generated/inbox_items.sql.go` — UpsertInboxItemByExternalRef (atomic ON CONFLICT)
- `internal/db/generated/notifications.sql.go` — CreateNotificationRef (dedupe)
- `internal/db/generated/password_reset_tokens.sql.go`

### Templ
- `web/templates/layout.templ` — `#modal-root`, avatar dropdown shell
- `web/templates/components/triage_modals.templ` — 3 HTMX modals + shell
- `web/templates/pages/profile.templ` — rewritten with avatar + edit
- `web/templates/pages/user_detail.templ`, `user_edit.templ`, `password_reset.templ` — new
- `web/templates/pages/companies.templ` — filter tabs + logo + deadline badge
- `web/templates/pages/company_detail.templ` — admin actions (logo upload + archive)
- `web/templates/pages/dashboard.templ` — DueSoonCompanies widget

### Docs
- `docs/webhook-api.md` — full guide cURL/iOS/Telegram/Zalo/bookmarklet/Zapier
- `docs/permissions.md` — 3-role matrix + enforcement map
- `docs/n8n-flows/{telegram,zalo-bot}-to-inbox.json` + `README.md`
- `plans/backlog.md` — P1/P2/P3 backlog tracker

## Open backlog (P1)

1. SMTP mailer integration — password reset emails (currently render link for manual share)
2. Telegram/Zalo n8n workflows import + activate on `auto.binhvuong.vn` (templates ready)
3. Drive orphan cleanup reaper (daily cron — currently any failed insert leaves Drive file)
4. Healthcheck endpoint `/healthz` + uptime monitor

## Open user-side action items

1. ⚠️ **Rotate `GOOGLE_CLIENT_SECRET`** — exposed in chat at least 2 times
2. Add `https://os.binhvuong.vn/auth/google/callback` to OAuth Console redirect URIs (if not already done)
3. Broadcast to 4 staff users (Hương/Minh/Hà/Linh) re-login required (JWT_SECRET rotated 4/22)
4. Verify Login Google E2E on real browser
5. Verify company logo rendering with real-sized PNG (test was 1x1)

## Last known state (4/22 evening)

- Docker image `8ac8b92e` running on server, 148+ handlers loaded
- Migration 24/u applied
- Drive proxy verified working (PNG bytes streamed, magic bytes correct)
- Notifier goroutine logged "first run starting" + "1 companies scanned, 1 notifications attempted"
- All test data cleaned from DB
- 5 active users: 1 owner, 2 manager, 2 staff

## Memory snapshots (key context I should never lose)

- This is a single-tenant solo founder OS (Bình Vương) for managing <10 companies + 20 staff
- Vietnamese language preference for assistant responses
- File uploads go to Google Drive (not local), via Drive proxy for image rendering
- Templ files (`*_templ.go`) gitignored — regenerated by Dockerfile
- "no sqlc" — generated/*.sql.go is hand-written; modifying schema requires editing BOTH `query/*.sql` (doc) AND `generated/*.sql.go` (code)
- Deploy = git push → SSH server pull + `docker build --no-cache` + `docker compose up -d`
- DB cleanup pattern: after E2E tests, always DELETE test rows + revert touched live rows

## Restore instructions (if context lost)

To resume work:
1. `cd /home/binhvuong-ws/projects/binhvuongos`
2. Read this journal + `plans/backlog.md` + recent `plans/<latest>/plan.md`
3. `git log --oneline | head -20` — see latest commits
4. `cat .env` to recall env structure (don't commit)
5. SSH `103.97.125.186` → `docker logs binhvuongos-app-1 --tail 50` to see runtime state
6. `psql` to `103.97.125.131` if data inspection needed

## Unresolved (this session)

1. Did user actually add `auth/google/callback` redirect URI to OAuth Console? (not confirmed in session)
2. Login Google E2E browser test never executed by user
3. Logo render quality with real PNGs unconfirmed (only 1x1 test)
