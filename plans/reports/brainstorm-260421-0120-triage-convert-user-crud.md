---
type: brainstorm
date: 2026-04-21 01:20 +07
slug: triage-convert-user-crud
status: approved
---

# Brainstorm: Triage Convert + User CRUD Role Simplification

## Problem statement

**Bug:** Triage từ inbox sang tasks/content/knowledge "không hoạt động" — handler hiện chỉ UPDATE `inbox_items.status='done'`, **không insert** vào target table. Fields `converted_to_type`/`converted_to_id` để trống. User thấy item biến mất khỏi inbox nhưng không xuất hiện ở đích.

**Feature:** Consolidate 4 roles cũ (owner/core_staff/ctv/staff) xuống 3 roles rõ ràng (owner/manager/staff), bổ sung đầy đủ User CRUD với phân quyền.

## Decisions (đã chốt qua AskUserQuestion)

| Q | A |
|---|---|
| Triage UX | Modal nhập quick-fields (HTMX) |
| Targets | Tasks + Content + Knowledge (bỏ Bookmarks) |
| Role model | 3 roles: owner + manager + staff |
| Perm matrix | Owner CRUD all / Manager CRUD staff only |
| Post-convert | Redirect về `/inbox` |
| Migration mapping | core_staff → manager, ctv → staff (plan ban đầu) |
| Password reset | Implement đầy đủ, SMTP thêm sau |

## Architecture

### Triage Convert flow

```
Inbox list/detail
  ↓ click "→ Task" button
HTMX modal (partial template)
  ↓ user điền quick-fields
POST /inbox/:id/convert?target=task|content|knowledge
  ↓
Handler ConvertInbox:
  BEGIN TX
    INSERT INTO <target> (title, content, company_id, submitted_by, ...)
    UPDATE inbox_items SET status='done',
      converted_to_type=$target, converted_to_id=$new_id,
      processed_at=NOW(), triage_notes=$notes
    WHERE id=$inbox_id AND status != 'done'
  COMMIT
  ↓ return 302 /inbox (with HX-Redirect for HTMX)
```

### Role + User CRUD

```
Actor role    → Can CRUD user with target role:
──────────────────────────────────────────────
owner         → owner | manager | staff
manager       → staff (only)
staff         → none (self-profile only qua /profile)
```

**Migration 000023:**
```sql
UPDATE users SET role='manager' WHERE role='core_staff';
UPDATE users SET role='staff'   WHERE role IN ('ctv', 'staff');
ALTER TABLE users
  ADD CONSTRAINT chk_user_role
  CHECK (role IN ('owner', 'manager', 'staff'));
```

### Password reset flow (SMTP-ready)

```
Manager/owner clicks "Reset password" trên user row
  ↓
Backend:
  1. Generate random token (32 hex)
  2. INSERT INTO password_reset_tokens (user_id, token, expires_at = NOW() + 1h)
  3. Build URL: https://os.binhvuong.vn/reset/<token>
  4. Render flash message: "Link đặt lại mật khẩu: <URL>" — manager copy + share manual
  ↓ (later, khi có SMTP)
  5. Gửi email tự động via internal/mailer/
  ↓
User clicks link → GET /reset/:token → form new password
  → POST /reset/:token → UPDATE users.password_hash + DELETE token row
```

**Bonus:** Cấu trúc này zero-refactor khi add SMTP — chỉ thêm call `mailer.SendResetEmail()` sau INSERT token.

## Files

### New
- `internal/db/migrations/000023_role_consolidation.up.sql` + `.down.sql`
- `internal/handler/inbox_convert.go` — ConvertInbox handler
- `internal/handler/user_crud.go` — UpdateUser, DeleteUser, ResetPassword handlers
- `internal/middleware/user_perm.go` — CanManageUser helper
- `web/templates/components/triage_modal.templ` — 3 modal variants
- `web/templates/pages/user_edit.templ` — edit form page
- `web/templates/pages/password_reset.templ` — reset form (token URL)

### Modified
- `cmd/server/main.go` — register new routes
- `internal/handler/users.go` — enforce role whitelist theo actor
- `internal/handler/inbox.go` — deprecate old `/inbox/:id/triage` (keep archive path only)
- `internal/db/query/users.sql` + `inbox_items.sql` + `password_reset_tokens.sql` — add queries
- `web/templates/pages/inbox.templ` + `inbox_detail.templ` — change quick-action buttons from form submit → HTMX modal trigger
- `web/templates/pages/users.templ` — role-gated edit/delete buttons
- `internal/handler/i18n.go` — role labels VN

## Effort estimate

| Phase | Scope | Effort |
|---|---|---|
| 1 | Migration 000023 + role constraint + update existing i18n labels | 20m |
| 2 | CanManageUser middleware + role whitelist enforcement in CreateUser | 30m |
| 3 | User CRUD handlers (update/delete/reset-password) + edit page templ | 60m |
| 4 | Triage convert handler (transaction) + query additions | 60m |
| 5 | 3 modal templ + HTMX wiring + inbox row button updates | 75m |
| 6 | Password reset public flow (`/reset/:token` GET+POST) | 45m |
| 7 | Deploy + E2E test (triage 3 targets, user CRUD, pw reset) | 30m |
| **Total** | | **~5h 20m** |

## Risks

| Risk | Severity | Mitigation |
|---|---|---|
| Migration UPDATE affects wrong rows | HIGH | Backup DB trước; WHERE clause tight; dry-run SELECT trước |
| Manager escalation (tạo user role=owner/manager qua form tampering) | HIGH | Server-side enforce role whitelist based on actor — not trust form value |
| Triage race: user nhấn nhiều lần cùng lúc | MED | WHERE status != 'done' trong UPDATE; nếu affected rows=0, coi như đã xử lý, trả success |
| Soft-delete user vẫn có active session → vẫn truy cập được | MED | Check `status='active'` trong AuthRequired middleware |
| Password reset token enumeration | LOW | Token 32-hex = 128 bits; always return generic response |

## Success criteria

- [ ] Triage 1 inbox item → Task: row xuất hiện `/tasks` + inbox marked done với `converted_to_id` đúng
- [ ] Triage → Content: tương tự
- [ ] Triage → Knowledge: tương tự
- [ ] Migration 000023 chạy clean, all existing users có role ∈ {owner, manager, staff}
- [ ] Manager login → `/users` chỉ thấy create form with role=staff; không tạo được manager qua form tampering
- [ ] Owner login → `/users` tạo được cả manager + staff
- [ ] Manager không edit/delete được owner hoặc manager khác (403)
- [ ] Reset password → generate token, render URL; token click mở form đặt pass mới, save thành công
- [ ] Old `/inbox/:id/triage` route archive still works; old convert logic removed/redirected

## Unresolved

1. Khi nào bổ sung SMTP mailer? Sau plan này hay sau nhiều plan nữa?
2. Soft delete vs hard delete cho user — hiện đang plan soft (`status='archived'`). OK không?
3. Auto-logout user khi bị soft-delete? (check status trong AuthRequired middleware)
