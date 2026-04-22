---
title: "Google Login + Self-Profile Edit + Avatar + User Detail"
description: "OAuth Google login (admin-gated), user self-edit profile + avatar upload, admin user detail page, header avatar dropdown"
status: completed
completed: 2026-04-22
priority: P1
effort: 3h30m
tags: [auth, oauth, google, profile, avatar, users]
blockedBy: []
blocks: []
created: 2026-04-22
---

# Google Login + Profile + Avatar + User Detail

## Overview

4 features liên quan, gộp 1 plan:
1. **Google OAuth login** coexist với email/password; chỉ user đã tồn tại (admin-gated) mới login được
2. **Self-profile edit** — user sửa full_name, phone, avatar của chính mình
3. **Avatar upload** — upload lên Drive, lưu URL vào `users.avatar_url`
4. **User detail admin view** — `/users/:id` xem profile + work logs + tasks assigned
5. **Header avatar dropdown** — click avatar → Profile / Logout

## Context

- **Brainstorm:** [../reports/brainstorm-260422-0907-google-login-profile-avatar.md](../reports/brainstorm-260422-0907-google-login-profile-avatar.md)
- **Reuse:** `GOOGLE_CLIENT_ID` + `SECRET` (add redirect URI Google Cloud Console), `drive.UploadFile`, `CanManageUser`
- **Không cần migration:** `users.avatar_url` đã có

## Phases

| # | Phase | Effort |
|---|---|---|
| 1 | [Google OAuth infra + login flow](phase-01-google-oauth.md) | 75m |
| 2 | [Self-profile edit + avatar upload](phase-02-profile-edit-avatar.md) | 60m |
| 3 | [User detail admin page](phase-03-user-detail.md) | 45m |
| 4 | [Header avatar dropdown](phase-04-header-dropdown.md) | 20m |
| 5 | [Deploy + E2E test](phase-05-deploy-test.md) | 30m |

## Google Cloud Console prerequisite

User phải tự cấu hình trước (không thể automate):
1. Google Cloud Console → API & Services → Credentials
2. Open OAuth 2.0 Client ID (cái đang dùng cho Drive)
3. Authorized redirect URIs → Add `https://os.binhvuong.vn/auth/google/callback`
4. Save

Nếu chưa có OAuth consent screen → khai báo tên app "Bình Vương OS", scopes: openid + email + profile.

## Success criteria

Xem brainstorm `## Success criteria` — 10 tiêu chí E2E.

## Risks

Xem brainstorm `## Risks` — CSRF state, email case, avatar quá lớn.
