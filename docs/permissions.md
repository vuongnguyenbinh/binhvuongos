# Permissions Matrix — Bình Vương OS

Tài liệu chính thức về phân quyền 3 role: `owner`, `manager`, `staff`. Bảng này phản ánh trạng thái code production; mọi thay đổi quyền phải update đây đồng thời.

## Legend

| Ký hiệu | Ý nghĩa |
|---|---|
| ✅ | Cho phép đầy đủ |
| ⚠️ | Có điều kiện (chi tiết ở cột ghi chú) |
| ❌ | Không cho phép, backend trả 403 / redirect login |

## Role definitions

| Role | Mô tả | Số lượng điển hình |
|---|---|---|
| `owner` | Super admin, chủ sở hữu hệ thống | 1 |
| `manager` | Quản lý, điều hành multi-company | 1–5 |
| `staff` | Nhân sự, CTV, thực thi công việc | 5–20 |

Legacy roles `core_staff`, `ctv` đã được migration 000023 ánh xạ sang `manager`, `staff`.

## Resource matrix

### Authentication & Profile

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| Đăng nhập email/password | ✅ | ✅ | ✅ | `POST /auth/login` |
| Đăng nhập Google OAuth | ✅ | ✅ | ✅ | `GET /auth/google` — email phải có trong `users` |
| Logout | ✅ | ✅ | ✅ | `POST /auth/logout` |
| Đổi password cá nhân | ✅ | ✅ | ✅ | `POST /profile/password` |
| Sửa profile (name, phone) | ✅ | ✅ | ✅ | `POST /profile/update` |
| Upload avatar | ✅ | ✅ | ✅ | `POST /profile/avatar` |
| Reset password user khác | ✅ | ⚠️ | ❌ | Manager chỉ reset được staff — `CanManageUser` |

### Inbox (`/inbox`)

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| View list | ✅ | ✅ | ✅ | Đều login-protected |
| Create qua form | ✅ | ✅ | ✅ | |
| Create qua webhook `/api/v1/inbox` | ✅ | ✅ | ✅ | Xác thực qua `X-API-Key` master key, không phân role |
| Triage → Task/Content/Knowledge | ✅ | ✅ | ✅ | `POST /inbox/:id/convert?target=…` |
| Archive | ✅ | ✅ | ✅ | `POST /inbox/:id/archive` |
| Edit detail | ✅ | ✅ | ✅ | `POST /inbox/:id/update` |

### Tasks (`/tasks`)

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| View list / detail | ✅ | ✅ | ✅ | |
| Create | ✅ | ✅ | ✅ | Có thể tạo task cho chính mình hoặc người khác |
| Edit bất kỳ task | ✅ | ✅ | ⚠️ | Staff chỉ edit task mình được assign (check ở handler) |
| Delete | ✅ | ✅ | ❌ | |

### Content (`/content`)

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| View list / detail | ✅ | ✅ | ✅ | |
| Create | ✅ | ✅ | ✅ | Content cần `company_id` bắt buộc |
| Edit / Publish | ✅ | ✅ | ⚠️ | Staff edit content mình là author |
| Delete | ✅ | ✅ | ❌ | |

### Companies (`/companies`)

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| View list | ✅ | ✅ | ✅ | |
| View detail | ✅ | ✅ | ⚠️ | Staff thấy tất cả nhưng chỉ có data tasks/content/campaigns của assignments |
| Create | ✅ | ✅ | ❌ | `POST /companies` |
| Edit info | ✅ | ✅ | ❌ | `POST /companies/:id` |
| Upload logo | ✅ | ✅ | ❌ | `POST /companies/:id/logo` |
| Archive / Unarchive | ✅ | ✅ | ❌ | `POST /companies/:id/archive` \| `/unarchive` |
| Assign user to company | ✅ | ✅ | ❌ | |

### Campaigns (`/campaigns`)

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| View | ✅ | ✅ | ✅ | |
| Create | ✅ | ✅ | ❌ | |
| Edit / Delete | ✅ | ✅ | ❌ | |

### Work logs (`/work-logs`)

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| Create work log của mình | ✅ | ✅ | ✅ | |
| View own logs | ✅ | ✅ | ✅ | |
| View all logs | ✅ | ✅ | ❌ | |
| Approve / Reject | ✅ | ✅ | ❌ | `POST /work-logs/:id/approve` |
| Batch approve | ✅ | ✅ | ❌ | |

### Knowledge (`/knowledge`)

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| View / Search | ✅ | ✅ | ✅ | |
| Create | ✅ | ✅ | ✅ | |
| Edit / Delete | ✅ | ✅ | ⚠️ | Staff chỉ sửa mục mình tạo |

### Bookmarks (`/bookmarks`)

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| View / Search | ✅ | ✅ | ✅ | Dùng chung |
| Create / Edit / Delete | ✅ | ✅ | ✅ | |

### User CRUD (`/users`)

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| List | ✅ | ✅ | ❌ | Admin route |
| View detail `/users/:id` | ✅ | ✅ | ❌ | |
| Create user role=owner | ✅ | ❌ | ❌ | |
| Create user role=manager | ✅ | ❌ | ❌ | Server whitelist `AllowedTargetRoles` |
| Create user role=staff | ✅ | ✅ | ❌ | |
| Edit user | ✅ | ⚠️ | ❌ | Manager chỉ edit staff (`CanManageUser`) |
| Delete (soft) user | ✅ | ⚠️ | ❌ | Tương tự; không tự xoá mình |
| Reset user password | ✅ | ⚠️ | ❌ | Tạo link `/reset/<token>`, TTL 1h, single-use |

### Admin settings

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| `/admin/settings` view/edit | ✅ | ✅ | ❌ | `google_oauth_*`, `n8n_webhook_url`, ... |
| `/admin/work-types` CRUD | ✅ | ✅ | ❌ | |

### Notifications

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| View own | ✅ | ✅ | ✅ | Tất cả thấy notifications của mình |
| Mark read | ✅ | ✅ | ✅ | |
| Mark all read | ✅ | ✅ | ✅ | |
| Nhận auto deadline cảnh báo | ✅ | ⚠️ | ⚠️ | Owner luôn nhận; manager/staff phải là assignee của company |

### Dashboard

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| Owner dashboard view | ✅ | ✅ | ❌ | Stats toàn org |
| Staff dashboard view | — | — | ✅ | Personal view |
| Widget "Công ty sắp hết hạn" | ✅ (all) | ✅ (all) | ✅ (assigned) | |

### Webhook API

| Action | Owner | Manager | Staff | Ghi chú |
|---|:---:|:---:|:---:|---|
| `POST /api/v1/inbox` | ✅ | ✅ | ✅ | Auth theo `X-API-Key` (shared secret), KHÔNG check role |

## Enforcement map (code location)

| Check | File | Mechanism |
|---|---|---|
| JWT validate + status=active | `internal/middleware/auth.go` | Cookie → `ParseWithClaims` → `GetUserByID` → check status |
| Role restriction | `internal/middleware/role.go` | `RequireRole("owner","manager")` wraps admin group |
| Manager↔staff scope | `internal/middleware/user_perm.go` | `CanManageUser(actor, target) bool` |
| Role whitelist on create/update | `internal/middleware/user_perm.go` | `AllowedTargetRoles(actor) []string`, `IsAllowedRole(actor, role)` |
| API key | `internal/middleware/api_key.go` | `subtle.ConstantTimeCompare` |
| Self-delete guard | `internal/handler/user_crud.go::DeleteUser` | So sánh `actor.ID == target.ID` |

## Route group registration

File: `cmd/server/main.go`

```go
// Public (no auth)
/login, /auth/login, /auth/logout, /auth/google, /auth/google/callback,
/reset/:token

// API (X-API-Key only)
/api/v1/* — auth: middleware.APIKeyAuth(cfg.APIKey)

// Protected (AuthRequired, all logged-in users)
/, /inbox, /tasks, /content, /companies, /campaigns, /knowledge, /bookmarks,
/work-logs, /profile, /notifications, /search, /dashboard

// Admin (owner + manager)
/users, /users/:id/..., /admin/settings, /admin/work-types
```

## Audit log

| Date | Change | Author |
|---|---|---|
| 2026-04-21 | Roles consolidated: `owner | manager | staff` (migration 000023) | — |
| 2026-04-22 | Company archive + logo + deadline perms added | — |
| 2026-04-22 | Matrix document created | — |
