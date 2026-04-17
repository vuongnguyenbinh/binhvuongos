# Phase 3: Rich Text + Integrations

## Overview
- **Priority:** P1
- **Effort:** 14h
- **Status:** Pending
- Markdown editor, content body, dashboard notes, real Notion sync, multi-file upload

## DB Migrations

### 000020_user_notes.up.sql
```sql
CREATE TABLE user_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_user_notes_user ON user_notes(user_id);
```

### 000021_content_body.up.sql
```sql
ALTER TABLE content ADD COLUMN IF NOT EXISTS body TEXT;
```

## Implementation Steps

### 1. SimpleMDE integration (2h)
- Add SimpleMDE CDN to `layout.templ`:
  - `<link rel="stylesheet" href="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.css">`
  - `<script src="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.js"></script>`
- Create `web/templates/components/markdown_editor.templ`:
  - MarkdownEditor(name, content, placeholder) — textarea + JS init SimpleMDE
- Create `web/static/js/markdown-init.js` — auto-init SimpleMDE on `.markdown-editor` textareas

### 2. Content body + Markdown render (3h)
- Migration 000021: add body TEXT to content table
- Add `goldmark` to go.mod (markdown→HTML lib)
- `internal/render/markdown.go` — RenderMarkdown(md string) → HTML string
- Update content detail template: show rendered body HTML
- Update content edit form: SimpleMDE editor for body field
- Update CreateContent + UpdateContent handlers: accept body field

### 3. Dashboard quick notes (2h)
- Migration 000020: user_notes table
- `internal/db/generated/user_notes.sql.go` — GetNote(user_id), UpsertNote(user_id, content)
- Dashboard template: SimpleMDE editor panel
- HTMX auto-save: hx-post="/dashboard/notes" hx-trigger="keyup changed delay:2s"
- Handler POST /dashboard/notes — upsert note content

### 4. Notion sync worker (6h)
- `internal/notion/client.go` — Notion API client:
  - CreatePage(databaseID, properties)
  - UpdatePage(pageID, properties)
  - Rate limiter: 2 req/s via time.Tick
- `internal/notion/mapper.go` — per-table mappers:
  - MapCompany(Company) → NotionProperties
  - MapTask(Task) → NotionProperties
  - MapContent(Content) → NotionProperties
  - MapWorkLog(WorkLog) → NotionProperties
  - MapCampaign(Campaign) → NotionProperties
  - MapKnowledgeItem(KnowledgeItem) → NotionProperties
- `internal/notion/sync.go` — SyncWorker:
  - For each syncable table: query WHERE sync_status IN ('pending','error') OR updated_at > synced_at
  - Create or update Notion page
  - Update sync_status, synced_at, notion_page_id
  - Log to notion_sync_log
- `cmd/server/main.go` — start goroutine: `go notion.StartCron(pool, cfg, 1*time.Hour)`
- Read Notion API key + database IDs from settings table

### 5. Notion sync button (1h)
- Admin Settings page: "Đồng bộ Notion ngay" button
- POST /api/v1/notion/sync — trigger sync immediately (run in goroutine)
- Return JSON with sync start message
- Admin UI: show last sync time from notion_sync_log

### 6. Multi-file upload (2h)
- Enhance upload handler: accept multiple files
- POST /upload/multi — loop upload to Drive, return array of {file_id, name, url}
- Update attachments JSONB: append new files to existing array
- Detail pages: show attachments list with download links
- Detail pages: multi-file upload form (input type="file" multiple)

## Files to Create
- `internal/db/migrations/000020_user_notes.up.sql` + down
- `internal/db/migrations/000021_content_body.up.sql` + down
- `internal/db/generated/user_notes.sql.go`
- `internal/render/markdown.go`
- `internal/notion/client.go`
- `internal/notion/mapper.go`
- `internal/notion/sync.go`
- `web/templates/components/markdown_editor.templ`
- `web/static/js/markdown-init.js`

## Files to Modify
- `web/templates/layout.templ` — SimpleMDE CDN
- `web/templates/pages/dashboard.templ` — notes panel
- `web/templates/pages/content_detail.templ` — body editor + render
- `internal/handler/content.go` — body field in create/update
- `internal/handler/dashboard.go` — notes handler
- `internal/handler/upload.go` — multi-file support
- `cmd/server/main.go` — Notion cron goroutine + new routes
- `go.mod` — goldmark

## Success Criteria
- [ ] SimpleMDE editor works on content body + dashboard notes
- [ ] Content detail renders Markdown as HTML
- [ ] Dashboard notes auto-save
- [ ] Notion sync runs every hour, syncs pending records
- [ ] Notion sync button triggers immediate sync
- [ ] Multi-file upload stores array in attachments JSONB
- [ ] Detail pages show attachment list with Drive links
