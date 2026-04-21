# Phase 3 — Password Reset (SMTP-ready)

**Effort:** 60m | **Priority:** P1 | **Depends on:** Phase 2

## Context
- Table `password_reset_tokens` exists (migration 000018)
- SMTP mailer chưa có — flow phải work với manager copy-paste link, add email sending sau này

## Files

### Create
- `internal/handler/password_reset.go` — 3 handlers
- `web/templates/pages/password_reset.templ` — public reset form

### Modify
- `internal/db/query/password_reset_tokens.sql` + generated code — add `CreateResetToken`, `GetValidResetToken`, `DeleteResetToken`
- `cmd/server/main.go` — register public routes (outside AuthRequired)
- `web/templates/pages/users.templ` — add "Reset password" button per row (perm-gated)

## Routes

```go
// Protected (admin)
admin.Post("/users/:id/reset-password", h.GenerateResetLink)   // returns flash + URL

// Public
app.Get("/reset/:token",  h.ResetPasswordPage)
app.Post("/reset/:token", h.ResetPassword)
```

## Queries

```sql
-- name: CreateResetToken :one
INSERT INTO password_reset_tokens (user_id, token, expires_at)
VALUES ($1, $2, NOW() + INTERVAL '1 hour')
RETURNING *;

-- name: GetValidResetToken :one
SELECT * FROM password_reset_tokens
WHERE token = $1 AND used_at IS NULL AND expires_at > NOW()
LIMIT 1;

-- name: MarkResetTokenUsed :exec
UPDATE password_reset_tokens SET used_at = NOW() WHERE token = $1;
```

## Handler flow

```go
func (h *Handler) GenerateResetLink(c *fiber.Ctx) error {
    actor := GetUser(c)
    target, err := h.queries.GetUserByID(...)
    if err != nil || !middleware.CanManageUser(actor, target) {
        return c.Status(403).SendString("Không có quyền")
    }
    // 32 hex = 128-bit entropy
    b := make([]byte, 16)
    rand.Read(b)
    token := hex.EncodeToString(b)
    _, _ = h.queries.CreateResetToken(c.Context(), generated.CreateResetTokenParams{
        UserID: target.ID,
        Token:  token,
    })
    resetURL := fmt.Sprintf("https://os.binhvuong.vn/reset/%s", token)
    // TODO: when SMTP ready, call h.mailer.SendReset(target.Email, resetURL)
    // For now, render flash with URL for manager to copy
    return render(c, pages.UserResetFlash(target.Email, resetURL))
}

func (h *Handler) ResetPasswordPage(c *fiber.Ctx) error {
    token := c.Params("token")
    _, err := h.queries.GetValidResetToken(c.Context(), token)
    if err != nil { return c.Status(404).SendString("Link hết hạn hoặc không hợp lệ") }
    return render(c, pages.PasswordResetForm(token))
}

func (h *Handler) ResetPassword(c *fiber.Ctx) error {
    token := c.Params("token")
    row, err := h.queries.GetValidResetToken(c.Context(), token)
    if err != nil { return c.Status(404).SendString("Link hết hạn") }
    password := c.FormValue("password")
    if len(password) < 8 { return c.Status(400).SendString("Mật khẩu ≥ 8 ký tự") }
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
    _, _ = h.queries.UpdateUserPassword(c.Context(), generated.UpdateUserPasswordParams{
        ID: row.UserID, PasswordHash: string(hash),
    })
    _ = h.queries.MarkResetTokenUsed(c.Context(), token)
    return c.Redirect("/login?reset=1")
}
```

## Todo
- [ ] Add queries: `CreateResetToken`, `GetValidResetToken`, `MarkResetTokenUsed`
- [ ] Check `UpdateUserPassword` exists (from earlier profile pw change); reuse or add
- [ ] Create `handler/password_reset.go`
- [ ] Create `pages/password_reset.templ` (form + submit)
- [ ] Create partial `pages/UserResetFlash` (render URL to copy)
- [ ] Register routes in `main.go` — public `/reset/:token` OUTSIDE `AuthRequired`
- [ ] Add "Reset pass" button in users.templ (perm-gated)
- [ ] `go build ./...` pass

## Success criteria
- Manager click Reset trên staff → URL link hiện ra flash
- URL expired (>1h) → 404 "hết hạn"
- URL dùng 2 lần → lần 2 fail (used_at set)
- User đặt pass mới (≥8 ký tự) → redirect login → login pass mới ok
- Manager click Reset trên manager khác → 403

## Risks
- Concurrent reset requests → 2 tokens cùng valid; chấp nhận (cả 2 đều dùng 1 lần)
- Token in URL logged by Cloudflare/nginx → acceptable risk, short TTL
- Self-reset via admin page nếu actor=staff → không cần support (staff dùng login forgot-password flow riêng khi có SMTP)
