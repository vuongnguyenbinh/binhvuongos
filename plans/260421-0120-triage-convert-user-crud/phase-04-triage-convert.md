# Phase 4 — Triage Convert Handler + Transactional Queries

**Effort:** 60m | **Priority:** P1

## Context

Target tables require these NOT NULL fields:
- **tasks**: `title` (≤500)
- **content**: `title`, `content_type` (≤30), `company_id` (UUID NOT NULL), `author_id` (UUID NOT NULL)
- **knowledge_items**: need verify — likely `title`, `body`

Current bug: `TriageInbox` only updates `status=done`. No insert into target.

## Files

### Create
- `internal/handler/inbox_convert.go` — `ConvertInbox(c)` handler

### Modify
- `internal/db/generated/inbox_items.sql.go` — add tx-aware `MarkInboxConverted` query
- `internal/db/query/inbox_items.sql` — add query def
- `internal/handler/inbox.go::TriageInbox` — replace broken logic with archive-only or redirect to convert; OR leave deprecated

## Transaction pattern

Use `pgx.Tx` via `pool.BeginTx`. Wrap:
1. INSERT target row (reuse existing `CreateTask`/`CreateContent`/`CreateKnowledgeItem` — but need tx-aware variant or raw SQL)
2. UPDATE inbox_items

```go
func (h *Handler) ConvertInbox(c *fiber.Ctx) error {
    actor := GetUser(c)
    inboxID := middleware.StringToUUID(c.Params("id"))
    target := c.Query("target") // "task" | "content" | "knowledge"

    tx, err := h.pool.BeginTx(c.Context(), pgx.TxOptions{})
    if err != nil { return c.Status(500).SendString("Lỗi DB") }
    defer tx.Rollback(c.Context())

    var newID pgtype.UUID
    switch target {
    case "task":
        // extract form: title, company_id, priority, due_date
        row := tx.QueryRow(c.Context(),
            `INSERT INTO tasks (title, description, company_id, priority, due_date, created_by, status)
             VALUES ($1, $2, $3, $4, $5, $6, 'todo') RETURNING id`,
            c.FormValue("title"), c.FormValue("description"),
            parseUUID(c.FormValue("company_id")),
            c.FormValue("priority"),
            parseDate(c.FormValue("due_date")),
            actor.ID)
        if err := row.Scan(&newID); err != nil {
            return c.Status(400).SendString("Tạo task fail: " + err.Error())
        }
    case "content":
        row := tx.QueryRow(c.Context(),
            `INSERT INTO content (title, content_type, company_id, author_id, status)
             VALUES ($1, $2, $3, $4, 'idea') RETURNING id`,
            c.FormValue("title"), c.FormValue("content_type"),
            parseUUID(c.FormValue("company_id")),
            actor.ID)
        if err := row.Scan(&newID); err != nil {
            return c.Status(400).SendString("Tạo content fail: " + err.Error())
        }
    case "knowledge":
        row := tx.QueryRow(c.Context(),
            `INSERT INTO knowledge_items (title, body, category, created_by)
             VALUES ($1, $2, $3, $4) RETURNING id`,
            c.FormValue("title"), c.FormValue("body"),
            c.FormValue("category"), actor.ID)
        if err := row.Scan(&newID); err != nil {
            return c.Status(400).SendString("Tạo knowledge fail: " + err.Error())
        }
    default:
        return c.Status(400).SendString("Target không hợp lệ")
    }

    // Mark inbox converted — idempotent via status guard
    cmd, err := tx.Exec(c.Context(),
        `UPDATE inbox_items
         SET status='done', converted_to_type=$2, converted_to_id=$3,
             processed_at=NOW(), triage_notes=$4
         WHERE id=$1 AND status != 'done'`,
        inboxID, target, newID, c.FormValue("triage_notes"))
    if err != nil { return c.Status(500).SendString("Update inbox fail") }
    if cmd.RowsAffected() == 0 {
        // already processed by another request — rollback to avoid duplicate target row
        return c.Status(409).SendString("Item đã được xử lý")
    }

    if err := tx.Commit(c.Context()); err != nil {
        return c.Status(500).SendString("Commit fail")
    }
    return c.Redirect("/inbox")
}
```

## Handler struct update

`Handler` needs `pool` access:
```go
type Handler struct {
    queries     *generated.Queries
    pool        *pgxpool.Pool   // NEW — for ad-hoc transactions
    config      *config.Config
    ownerUserID pgtype.UUID
}
```

Alternative: add `BeginTx` method on `generated.Queries`. Simpler = expose pool.

## Route

```go
app.Post("/inbox/:id/convert", h.ConvertInbox)
// Deprecate but keep for compat:
app.Post("/inbox/:id/triage", h.ArchiveInbox) // re-alias — old route now archive-only
```

## Todo
- [ ] Add `pool *pgxpool.Pool` to `Handler` struct + pass from `main.go::NewHandler`
- [ ] Create `handler/inbox_convert.go` with 3 target branches
- [ ] Add `parseUUID`, `parseDate` helpers (or reuse if exist)
- [ ] Register `/inbox/:id/convert` route
- [ ] Deprecate old `/inbox/:id/triage` → alias to archive
- [ ] `go build ./...` pass

## Success criteria
- Convert inbox to task: row in `/tasks`, inbox `converted_to_type='task'`, `converted_to_id=<task_id>`
- Convert to content: `/content` has row
- Convert to knowledge: `/knowledge` has row
- Second convert call on same item → 409 (idempotent)
- INSERT fail (missing company_id for content) → rollback; inbox item stays raw; flash error

## Risks
- Company dropdown required for content (company_id NOT NULL) — modal must prevent submit if empty
- `due_date` parsing: user leaves empty → `NULL` OK
- Rollback safe: pgx `defer Rollback` no-op after Commit
