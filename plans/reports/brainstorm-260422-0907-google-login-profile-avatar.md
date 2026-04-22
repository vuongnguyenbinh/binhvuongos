---
type: brainstorm
date: 2026-04-22 09:07 +07
slug: google-login-profile-avatar
status: approved
---

# Brainstorm: Google Login + Self-Profile Edit + Avatar + User Detail

## Problem statement

1. **Login Google**: Hiện chỉ có email/password. Muốn thêm Google OAuth để nhanh + không cần nhớ pass. Áp dụng cho mọi role (owner bao gồm).
2. **Self-profile edit**: User hiện tại chỉ đổi được password qua `/profile`, không sửa được full_name/phone/avatar. Cần hoàn thiện.
3. **Avatar upload**: Chưa có cơ chế upload ảnh đại diện; field `users.avatar_url` có sẵn chưa dùng.
4. **User detail admin view**: Admin (owner/manager) hiện chỉ thấy list, chưa có page xem profile + work history của 1 user.

## Decisions (đã chốt)

| Q | A |
|---|---|
| User CRUD extension | Self-profile edit + User detail page (admin) + Avatar upload |
| Google login scope | Chỉ user đã đăng ký (admin-gated) |
| Auth coexist | Cả email/password + Google |
| OAuth creds | Reuse `GOOGLE_CLIENT_ID` với Drive (same app, add redirect URI) |
| Header avatar click | Dropdown Profile + Logout |
| Self user detail | Không cần (staff dùng `/profile`) |
| Owner login Google | Có |

## Architecture

### Auth flow

```
[Login page]
  ├─ <form method=POST /auth/login>  (existing)
  └─ <a href="/auth/google">Đăng nhập với Google</a>
         │
         ▼
GET /auth/google
  1. Generate state = rand_hex(16)
  2. Set cookie gstate=<state> (HttpOnly, 10min)
  3. Build Google authorize URL:
     https://accounts.google.com/o/oauth2/v2/auth?
       client_id=<ID>&redirect_uri=<BASE>/auth/google/callback
       &response_type=code&scope=openid%20email%20profile
       &state=<state>&access_type=online&prompt=select_account
  4. 302 redirect
         │
         ▼ (Google consent)
         │
GET /auth/google/callback?code=X&state=Y
  1. Verify state cookie matches
  2. POST https://oauth2.googleapis.com/token
     → exchange code for access_token
  3. GET https://openidconnect.googleapis.com/v1/userinfo
     → {email, name, picture, email_verified}
  4. Lookup: SELECT * FROM users
             WHERE LOWER(email)=LOWER($1)
               AND status='active' AND deleted_at IS NULL
  5a. Found → generate JWT cookie → 302 /inbox
  5b. Not found → 403 "Tài khoản <email> chưa được cấp quyền"
```

### Profile flow

```
GET /profile                       → existing view
POST /profile/update               → edit full_name, phone (email+role locked)
POST /profile/avatar               → multipart file → drive.UploadFile → users.avatar_url
POST /profile/password             → existing
POST /auth/logout                  → existing
```

### User detail (admin)

```
GET /users/:id  [owner|manager only, perm-gated via CanManageUser]
  Shows:
    - Profile card (avatar, full_name, email, role, phone, status, created_at)
    - Assigned companies (via user_company_assignments)
    - Recent 10 work logs (WHERE user_id=$1)
    - Recent 10 tasks assigned (WHERE assignee_id=$1)
    - Edit button → /users/:id/edit (existing)
```

### Header avatar dropdown

```
┌─────────────────────────────────┐
│  Logo   Nav   ...   [avatar ▼]  │
│                      ├ Profile   │ (click → /profile)
│                      └ Logout    │ (click → POST /auth/logout)
└─────────────────────────────────┘
```

Vanilla JS `onclick` toggles dropdown visibility. No HTMX needed.

## Database — zero migration

- `users.avatar_url` đã có (confirmed in GetUserByID scan)
- `users.full_name`, `phone` đã có

## Config

Thêm 1 env:
```go
// config.go
GoogleRedirectURI string  // e.g. https://os.binhvuong.vn/auth/google/callback
```

Reuse `GoogleClientID` + `GoogleClientSecret`.

**Google Cloud Console setup:**
- OAuth 2.0 Client ID → Authorized redirect URIs → add `https://os.binhvuong.vn/auth/google/callback`
- Authorized JavaScript origins → `https://os.binhvuong.vn`
- OAuth consent screen scopes: `openid`, `email`, `profile` (đã mặc định cho external users)

## Files

### New
- `internal/oauth/google.go` — pure OAuth client (ExchangeCode, FetchUserinfo)
- `internal/handler/google_auth.go` — GoogleLoginRedirect, GoogleCallback handlers
- `internal/handler/profile_edit.go` — ProfileUpdate, ProfileAvatar handlers
- `web/templates/pages/user_detail.templ` — admin user detail page

### Modify
- `internal/handler/users.go` — add UserDetail handler
- `web/templates/pages/profile.templ` — add full edit form + avatar upload
- `web/templates/pages/login.templ` — add Google button
- `web/templates/layout.templ` + `header.templ` (nếu có) — avatar dropdown
- `cmd/server/main.go` — 5 routes: /auth/google, /auth/google/callback, /profile/update, /profile/avatar, /users/:id
- `internal/config/config.go` — `GoogleRedirectURI`

## Scope breakdown (5 phases)

| # | Phase | Files | Effort |
|---|---|---|---|
| 1 | Google OAuth infra + login flow | oauth/google.go, google_auth.go, login.templ, config, main.go | 75m |
| 2 | Self-profile edit + avatar upload | profile_edit.go, profile.templ | 60m |
| 3 | User detail admin page | user_detail.templ, users.go::UserDetail | 45m |
| 4 | Header avatar dropdown + nav wiring | layout.templ, app.js | 20m |
| 5 | Deploy + E2E test | — | 30m |
| **Total** | | | **~3h30m, ~440 LOC** |

## Risks

| Risk | Mitigation |
|---|---|
| CSRF on OAuth callback | `state` param random hex + cookie, verify strict |
| Email case mismatch | `LOWER(email) = LOWER($1)` in user lookup |
| Google email_verified=false (possible for some providers) | Reject with "Email chưa verify ở Google" |
| Avatar upload quá lớn | Reuse 50MB cap, Drive folder same |
| Drive OAuth refresh_token bị invalidate khi user consent app lại cho login | Login flow dùng `access_type=online` (no refresh_token issue for login) → không conflict; Drive vẫn dùng refresh_token riêng vì nó dùng service account-style flow |
| Owner lỡ xoá email Google primary → không login Google được | Email/password vẫn còn (dual auth) |
| Redirect URI mismatch | Cấu hình production = `os.binhvuong.vn`; local dev không login Google (chấp nhận) |

## Success criteria

- [ ] Click "Đăng nhập với Google" → Google consent → callback → JWT cookie → /inbox
- [ ] Email Google không match user DB → 403 với message
- [ ] Email Google match soft-deleted user → 403
- [ ] /profile form edit full_name + phone → save → reflect immediately
- [ ] Avatar upload <10MB image → Drive URL lưu trong users.avatar_url → hiển thị header + profile
- [ ] /users/:id (admin view) → thấy profile + recent logs + tasks
- [ ] Manager xem /users/:owner_id → 403
- [ ] Header click avatar → dropdown Profile + Logout
- [ ] Logout button functional
- [ ] Email/password login vẫn work (regression)
- [ ] Webhook regression OK

## Unresolved

1. Avatar default khi user chưa upload — hiện icon `👤` hay initial letters? (KISS: initials với background forest)
2. Google Cloud Console redirect URI — user (Bình Vương) sẽ tự cấu hình hay cần tôi guide step-by-step?
3. Sau Google login lần đầu, có tự populate `avatar_url` từ Google picture không? (Đề xuất: có, chỉ khi `avatar_url` đang null — không overwrite ảnh user đã upload)
