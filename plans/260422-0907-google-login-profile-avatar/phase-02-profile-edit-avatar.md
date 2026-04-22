# Phase 2 — Self-Profile Edit + Avatar Upload

**Effort:** 60m

## Files

### New
- `internal/handler/profile_edit.go` — ProfileUpdate, ProfileAvatar

### Modify
- `web/templates/pages/profile.templ` — form edit + avatar upload
- `cmd/server/main.go` — 2 routes
- `internal/db/generated/users.sql.go` — `UpdateOwnProfile` query

## Queries

```go
type UpdateOwnProfileParams struct {
    ID       pgtype.UUID
    FullName string
    Phone    string
}

func (q *Queries) UpdateOwnProfile(ctx, arg) error {
    _, err := q.pool.Exec(ctx,
        "UPDATE users SET full_name=$2, phone=$3 WHERE id=$1",
        arg.ID, arg.FullName, arg.Phone)
    return err
}
```

## Handler

```go
// profile_edit.go

func (h *Handler) ProfileUpdate(c *fiber.Ctx) error {
    user := GetUser(c)
    fullName := strings.TrimSpace(c.FormValue("full_name"))
    if fullName == "" {
        return c.Status(400).SendString("Thiếu họ tên")
    }
    if err := h.queries.UpdateOwnProfile(c.Context(), generated.UpdateOwnProfileParams{
        ID: user.ID, FullName: fullName, Phone: c.FormValue("phone"),
    }); err != nil {
        return c.Status(500).SendString("Lỗi cập nhật")
    }
    return c.Redirect("/profile")
}

func (h *Handler) ProfileAvatar(c *fiber.Ctx) error {
    user := GetUser(c)
    file, err := c.FormFile("avatar")
    if err != nil {
        return c.Status(400).SendString("Thiếu file")
    }
    if file.Size > 10*1024*1024 {
        return c.Status(400).SendString("Ảnh quá lớn (max 10MB)")
    }
    // Reuse drive.UploadFile via inbox_webhook_helpers pattern
    src, _ := file.Open()
    defer src.Close()
    cfg := &drive.Config{
        ClientID: h.config.GoogleClientID, ClientSecret: h.config.GoogleClientSecret,
        RefreshToken: h.config.GoogleRefreshToken, FolderID: h.config.GoogleDriveFolderID,
    }
    result, err := drive.UploadFile(c.Context(), cfg, file.Filename, file.Header.Get("Content-Type"), src)
    if err != nil {
        return c.Status(500).SendString("Upload fail: " + err.Error())
    }
    avatarURL := result.WebViewLink
    if avatarURL == "" {
        avatarURL = fmt.Sprintf("https://drive.google.com/file/d/%s/view", result.FileID)
    }
    if err := h.queries.UpdateUserAvatar(c.Context(), user.ID, avatarURL); err != nil {
        return c.Status(500).SendString("Lỗi lưu avatar")
    }
    return c.Redirect("/profile")
}
```

## Profile templ

Extend `profile.templ`:
```templ
<div class="bg-surface border border-hairline rounded-lg p-6 mb-6">
    <h3 class="display text-lg font-bold mb-4">Thông tin cá nhân</h3>
    <form method="POST" action="/profile/update" class="space-y-4">
        <div>
            <label class="text-sm">Email</label>
            <input type="email" value={ user.Email } disabled class="...bg-cream/50 text-muted"/>
        </div>
        <div>
            <label class="text-sm">Họ tên *</label>
            <input name="full_name" value={ user.FullName } required class="..."/>
        </div>
        <div>
            <label class="text-sm">SĐT</label>
            <input name="phone" value={ user.Phone } class="..."/>
        </div>
        <button type="submit" class="px-6 py-2 bg-forest text-white rounded mono text-sm">LƯU</button>
    </form>
</div>

<div class="bg-surface border border-hairline rounded-lg p-6 mb-6">
    <h3 class="display text-lg font-bold mb-4">Ảnh đại diện</h3>
    if user.AvatarURL != "" {
        <img src={ user.AvatarURL } class="w-20 h-20 rounded-full mb-3"/>
    }
    <form method="POST" action="/profile/avatar" enctype="multipart/form-data" class="flex gap-2">
        <input type="file" name="avatar" accept="image/*" required class="text-sm"/>
        <button type="submit" class="px-4 py-2 bg-forest text-white rounded text-xs mono">UPLOAD</button>
    </form>
</div>

<!-- existing password change form below -->
```

## Routes

```go
app.Post("/profile/update", h.ProfileUpdate)
app.Post("/profile/avatar", h.ProfileAvatar)
// existing: /profile, /profile/password
```

## Todo
- [ ] Query `UpdateOwnProfile`
- [ ] `profile_edit.go` 2 handlers
- [ ] `profile.templ` form edit + avatar upload
- [ ] Routes
- [ ] `go build` pass

## Success criteria
- POST /profile/update → user row updated in DB
- POST /profile/avatar image → Drive URL saved, image renders on profile page
- Avatar >10MB → 400 error
