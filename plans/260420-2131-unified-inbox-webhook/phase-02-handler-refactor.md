# Phase 2 — Handler Refactor (JSON + Multipart + Idempotency)

**Effort:** 90m | **Priority:** P1 | **Status:** pending | **Depends on:** Phase 1

## Context
- Existing handler: `internal/handler/api_handlers.go::APICreateInbox` (~30 LOC)
- Drive upload: `internal/handler/upload.go::Upload` (reuse `drive.UploadFile`)
- Config: `internal/config/config.go` — add `OwnerEmail` field
- Existing user lookup: `queries.GetUserByEmail`

## Overview
Refactor `APICreateInbox` to:
1. Branch theo `Content-Type` (JSON vs multipart/form-data)
2. Validate input (whitelist item_type, attachments ≤10, sizes)
3. Idempotency check qua `GetInboxByExternalRef`
4. Upload file → Drive (multipart path)
5. Resolve `submitted_by` = owner user ID (cached on startup)
6. Insert với `attachments` JSONB + `external_ref`

## Files to modify
- `internal/handler/api_handlers.go` — main refactor
- `internal/handler/handler.go` — add `ownerUserID` field on Handler struct, resolve on init
- `internal/config/config.go` — add `OwnerEmail string` from env

## Files to create
- `internal/handler/inbox_webhook_helpers.go` — detectItemType, validateItemType, uploadAttachment helpers (để api_handlers.go không phình)

## Implementation steps

### 1. Config — add OWNER_EMAIL
```go
// internal/config/config.go
type Config struct {
    // ...existing
    OwnerEmail string
}

func Load() *Config {
    return &Config{
        // ...
        OwnerEmail: getenvOr("OWNER_EMAIL", "vuongnguyenbinh@gmail.com"),
    }
}
```

### 2. Handler init — resolve owner ID once
```go
// internal/handler/handler.go
type Handler struct {
    queries     *generated.Queries
    config      *config.Config
    ownerUserID uuid.UUID
}

func NewHandler(q *generated.Queries, cfg *config.Config) *Handler {
    h := &Handler{queries: q, config: cfg}
    // Resolve owner ID — fail fast if missing
    user, err := q.GetUserByEmail(context.Background(), cfg.OwnerEmail)
    if err != nil {
        log.Fatalf("OWNER_EMAIL=%s not found in users table: %v", cfg.OwnerEmail, err)
    }
    h.ownerUserID = user.ID
    return h
}
```

### 3. Helpers (new file)
```go
// internal/handler/inbox_webhook_helpers.go
package handler

var allowedItemTypes = map[string]bool{
    "note": true, "link": true, "image": true, "voice": true, "file": true,
}

const (
    maxAttachments     = 10
    maxContentBytes    = 10 * 1024      // 10KB
    maxMultipartBytes  = 50 * 1024 * 1024 // 50MB
    maxAttachmentURL   = 2048
)

// detectItemType auto-detect từ content / url
func detectItemType(content, url string) string {
    if url != "" || isURL(content) {
        return "link"
    }
    return "note"
}

// validateInboxInput checks common validations
func validateInboxInput(content, source, itemType string, attachmentURLs []string) error {
    if content == "" {
        return errors.New("content is required")
    }
    if len(content) > maxContentBytes {
        return fmt.Errorf("content too large (max %d bytes)", maxContentBytes)
    }
    if source == "" || len(source) > 30 {
        return errors.New("source required, max 30 chars")
    }
    if itemType != "" && !allowedItemTypes[itemType] {
        return fmt.Errorf("invalid item_type (allowed: note, link, image, voice, file)")
    }
    if len(attachmentURLs) > maxAttachments {
        return fmt.Errorf("max %d attachments", maxAttachments)
    }
    for _, u := range attachmentURLs {
        if len(u) > maxAttachmentURL {
            return errors.New("attachment URL too long")
        }
    }
    return nil
}

// uploadAttachmentToDrive handles multipart file → Drive → attachment struct
func (h *Handler) uploadAttachmentToDrive(c *fiber.Ctx, file *multipart.FileHeader) (map[string]any, error) {
    if file.Size > maxMultipartBytes {
        return nil, fmt.Errorf("file too large (max %d bytes)", maxMultipartBytes)
    }
    src, err := file.Open()
    if err != nil {
        return nil, err
    }
    defer src.Close()

    cfg := &drive.Config{
        ClientID: h.config.GoogleClientID,
        ClientSecret: h.config.GoogleClientSecret,
        RefreshToken: h.config.GoogleRefreshToken,
        FolderID: h.config.GoogleDriveFolderID,
    }
    if cfg.RefreshToken == "" {
        return nil, errors.New("drive not configured")
    }
    result, err := drive.UploadFile(c.Context(), cfg, file.Filename, file.Header.Get("Content-Type"), src)
    if err != nil {
        return nil, err
    }
    url := result.WebViewLink
    if url == "" {
        url = fmt.Sprintf("https://drive.google.com/file/d/%s/view", result.FileID)
    }
    return map[string]any{
        "url":           url,
        "filename":      result.FileName,
        "mime":          file.Header.Get("Content-Type"),
        "drive_file_id": result.FileID,
    }, nil
}
```

### 4. Refactor APICreateInbox
```go
// internal/handler/api_handlers.go
func (h *Handler) APICreateInbox(c *fiber.Ctx) error {
    contentType := c.Get("Content-Type")
    var (
        content, url, source, itemType, externalRef string
        attachments []map[string]any
    )

    if strings.HasPrefix(contentType, "multipart/form-data") {
        content = c.FormValue("content")
        url = c.FormValue("url")
        source = c.FormValue("source", "api")
        itemType = c.FormValue("item_type")
        externalRef = c.FormValue("external_ref")

        if file, err := c.FormFile("file"); err == nil && file != nil {
            att, err := h.uploadAttachmentToDrive(c, file)
            if err != nil {
                return c.Status(400).JSON(APIError("UPLOAD_ERROR", err.Error()))
            }
            attachments = append(attachments, att)
        }
    } else {
        var input struct {
            Content        string   `json:"content"`
            URL            string   `json:"url"`
            Source         string   `json:"source"`
            ItemType       string   `json:"item_type"`
            AttachmentURLs []string `json:"attachment_urls"`
            ExternalRef    string   `json:"external_ref"`
        }
        if err := c.BodyParser(&input); err != nil {
            return c.Status(400).JSON(APIError("VALIDATION_ERROR", err.Error()))
        }
        content, url, source, itemType, externalRef = input.Content, input.URL, input.Source, input.ItemType, input.ExternalRef
        if source == "" {
            source = "api"
        }
        for _, u := range input.AttachmentURLs {
            attachments = append(attachments, map[string]any{"url": u})
        }
    }

    // Extract attachment URLs for validation
    attachURLs := make([]string, 0, len(attachments))
    for _, a := range attachments {
        if u, ok := a["url"].(string); ok {
            attachURLs = append(attachURLs, u)
        }
    }
    if err := validateInboxInput(content, source, itemType, attachURLs); err != nil {
        return c.Status(400).JSON(APIError("VALIDATION_ERROR", err.Error()))
    }

    // Auto-detect item_type
    if itemType == "" {
        itemType = detectItemType(content, url)
    }

    // Idempotency check
    if externalRef != "" {
        if existing, err := h.queries.GetInboxByExternalRef(c.Context(), generated.GetInboxByExternalRefParams{
            Source: sql.NullString{String: source, Valid: true},
            ExternalRef: sql.NullString{String: externalRef, Valid: true},
        }); err == nil {
            return c.Status(200).JSON(fiber.Map{"success": true, "duplicate": true, "data": existing})
        }
    }

    // Marshal attachments to JSONB
    attachJSON, _ := json.Marshal(attachments)

    item, err := h.queries.CreateInboxItem(c.Context(), generated.CreateInboxItemParams{
        Content:     content,
        URL:         sql.NullString{String: url, Valid: url != ""},
        Source:      sql.NullString{String: source, Valid: true},
        ItemType:    sql.NullString{String: itemType, Valid: true},
        SubmittedBy: uuid.NullUUID{UUID: h.ownerUserID, Valid: true},
        Attachments: attachJSON,
        ExternalRef: sql.NullString{String: externalRef, Valid: externalRef != ""},
    })
    if err != nil {
        return c.Status(500).JSON(APIError("DB_ERROR", err.Error()))
    }
    return c.Status(201).JSON(APISuccess(item))
}
```

## Todo
- [ ] Add `OwnerEmail` to `config.Config` + `.env.example`
- [ ] Modify `NewHandler` to resolve `ownerUserID` (fail-fast)
- [ ] Create `internal/handler/inbox_webhook_helpers.go`
- [ ] Refactor `APICreateInbox` (branch JSON/multipart, validation, idempotency)
- [ ] Update `generated/inbox_items.sql.go` params struct usage (after sqlc regen)
- [ ] Run `go build ./...` — no compile errors
- [ ] Manual test: cURL JSON → 201; cURL multipart → 201 with attachment; duplicate external_ref → 200 with `duplicate: true`

## Success criteria
- JSON request: content + url + source → 201 với item_type auto-detected
- Multipart with file: upload Drive → attachments[0] có url + drive_file_id
- Duplicate external_ref: 200 (không 500), return existing item
- Invalid item_type: 400 với error message
- Attachments >10: 400
- Content >10KB: 400

## Risks
- Owner user không tồn tại khi seed chạy → startup fail (fail-fast là đúng, không silent)
- `BodyParser` với multipart có thể conflict — test kỹ
- sqlc generated `CreateInboxItemParams.Attachments` type: nếu là `[]byte` thì dùng `json.Marshal`; nếu là `pgtype.JSONB` phải dùng constructor khác. Verify sau Phase 1.
