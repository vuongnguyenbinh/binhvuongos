---
title: "Backlog — Việc chưa xử lý"
updated: 2026-04-22
---

# Backlog

Tổng hợp các câu hỏi/đề xuất chưa xử lý từ các plan + report đã xong. Mỗi mục có source, khi nào cân nhắc lại.

## P1 — Ảnh hưởng production

### 1. SMTP mailer integration
- **Source:** brainstorm-260421 triage-convert-user-crud Q1
- **Hiện tại:** Password reset chỉ render link URL cho manager copy + gửi thủ công qua Zalo/Telegram
- **Cần làm:** Tích hợp SMTP (Resend / Postmark / SES) → gửi email reset + welcome + notifications
- **Impact khi có:** Self-service "forgot password" cho user, không phụ thuộc manager; notification email cho comments/triage

### 2. Telegram/Zalo n8n workflows deploy
- **Source:** brainstorm + phase-05-deploy-test webhook plan
- **Hiện tại:** JSON workflow template đã commit (`docs/n8n-flows/`), chưa import vào `auto.binhvuong.vn`
- **Cần làm:** Import telegram-to-inbox.json + zalo-bot-to-inbox.json, set `BVOS_API_KEY` + `ZALO_BOT_SECRET` env vars, activate, smoke test
- **Blocker:** Cần Telegram BOT_TOKEN + Zalo bot credentials

### 3. Drive orphan cleanup reaper
- **Source:** code-review-260420 webhook
- **Hiện tại:** Nếu Drive upload OK nhưng DB insert fail → file Drive orphan (không link tới inbox_item)
- **Cần làm:** Daily cron job list Drive files trong folder, diff với `inbox_items.attachments` + `content.attachments` + ..., xoá orphan >7 days
- **Impact:** Ngăn Drive quota bloat. Ít xảy ra — YAGNI cho hiện tại nhưng nhớ khi scale.

## P2 — Nên làm sớm

### 4. Healthcheck endpoint + Uptime monitor
- **Source:** tester-260421 header-size-fix Q2
- **Hiện tại:** Không có endpoint `/healthz`; user tự phát hiện 5xx
- **Cần làm:** Route `GET /healthz` trả 200 + DB ping; đăng ký BetterStack/UptimeRobot ping 30s
- **Value:** Catch deploy fail sớm, alert Telegram

### 5. `query/*.sql` source of truth
- **Source:** code-review-260420 webhook
- **Hiện tại:** `query/*.sql` có drift vs hand-written `generated/*.sql.go` (project không dùng sqlc dù cấu trúc giống)
- **Cần làm:** Decide — xoá `query/*.sql` hoặc viết tool auto-sync, hoặc treat query/*.sql as design doc only
- **Risk:** Hiện giờ ai đó sửa query/*.sql sẽ nghĩ code tự regenerate → không happen → inconsistency

### 6. Legacy JWT role mismatch window
- **Source:** triage-convert plan Phase 1 risk
- **Hiện tại:** Session cũ với `role=core_staff` hoặc `ctv` trong JWT sẽ fail `RequireRole("owner","manager")` → user gặp 403 cho đến khi login lại
- **Cần làm:** Broadcast notify cho 4 user (2 core_staff + 2 ctv) logout + login lại. Hoặc bump JWT secret để invalidate all tokens.

## P3 — Nice to have / YAGNI

### 7. Multi-tenant manager scope
- **Source:** brainstorm triage-convert Q
- **Hiện tại:** Manager có thể CRUD bất kỳ staff nào
- **Cần làm:** Constrain manager chỉ manage staff trong companies được assign (qua `user_company_assignments`)
- **Khi cần:** Khi scale >3 manager và các cluster không nên xem chung

### 8. Self-service "forgot password" từ login
- **Source:** triage-convert phase-03 risk
- **Hiện tại:** Chỉ admin-initiated reset
- **Cần làm:** Route `/forgot` (public) → nhập email → gửi reset link
- **Blocker:** Cần SMTP (#1)

### 9. Webhook rate-limit riêng
- **Source:** brainstorm webhook Q2
- **Hiện tại:** Global 60 req/phút cho toàn app, áp dụng cả webhook
- **Cần làm:** Tách rate limit webhook 600 req/phút (bot traffic cao)
- **Khi cần:** Khi n8n flow nhận spam hoặc event batch

### 10. Fiber BodyLimit review
- **Source:** tester header-size-fix Q1
- **Hiện tại:** Global 60MB (accommodate webhook multipart 50MB + overhead)
- **Cần làm:** Không có endpoint khác cần upload lớn. Giữ nguyên.

### 11. Webhook Zalo flow smoke test
- **Source:** webhook plan Phase 5
- **Hiện tại:** Chỉ test qua cURL, chưa setup Zalo OA/bot thật
- **Cần làm:** Tạo Zalo personal bot qua bot.zapps.me, set webhook, gửi text + image, verify Drive upload.

### 12. Templ files gitignore side effects
- **Source:** gitignore commit 907bd1b
- **Hiện tại:** `*_templ.go` được gitignore; dev local cần chạy `templ generate` trước khi `go build`
- **Cần làm:** Thêm Makefile target `make dev` để auto-generate, hoặc bump instruction trong CLAUDE.md

### 13. Soft-delete strategy review
- **Source:** triage-convert Q2
- **Hiện tại:** User dùng `deleted_at` (soft); inbox dùng `status='archived'` (vừa có deleted_at vừa có status)
- **Cần làm:** Chuẩn hoá — hoặc tất cả dùng `deleted_at`, hoặc tất cả dùng `status`. Tránh nhầm lẫn.

## Done / Closed

- ✅ n8n deploy location: `auto.binhvuong.vn` (confirmed)
- ✅ Auto-logout soft-deleted user: implemented via `status='active'` check trong AuthRequired
- ✅ Triage race: atomic `WHERE status != 'done'` + 409 response
- ✅ Templ files gitignored + Dockerfile rebuild

## Cách dùng file này

- Bổ sung mục mới với số thứ tự tiếp theo
- Di chuyển lên **Done** khi xử lý xong, kèm ngày
- Review 2 tuần/lần, promote P3 → P2 nếu cần
