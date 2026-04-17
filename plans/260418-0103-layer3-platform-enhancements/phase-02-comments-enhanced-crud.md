# Phase 2: Comments System + Enhanced CRUD

## Overview
- **Priority:** P1
- **Effort:** 29.5h
- **Status:** Pending
- Polymorphic comments system + nâng cấp filter/display cho tất cả modules

## DB Migrations

### 000019_comments.up.sql
```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(30) NOT NULL,
    entity_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_comments_entity ON comments(entity_type, entity_id, created_at DESC);
CREATE INDEX idx_comments_user ON comments(user_id);
CREATE TRIGGER tr_comments_updated_at BEFORE UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
```

## Implementation Steps

### 1. Comments table + queries (1h)
- Migration 000019
- `internal/db/generated/comments.sql.go`:
  - ListCommentsByEntity(entity_type, entity_id) → join users for name+avatar
  - CreateComment(entity_type, entity_id, user_id, body)
  - DeleteComment(id)
  - CountCommentsByEntity(entity_type, entity_id)

### 2. Comments handler (2h)
- `internal/handler/comments.go`:
  - POST /comments — body: entity_type, entity_id, body → create + HTMX return new comment HTML
  - DELETE /comments/:id — soft delete, HTMX swap empty
  - GET /comments?entity_type=task&entity_id=xxx — HTMX partial load comment list

### 3. Comments templ component (2h)
- `web/templates/components/comments.templ`:
  - CommentSection(entityType, entityID, comments []CommentItem) — renders list + add form
  - CommentItem — avatar, name, body, time ago, delete button (if owner)
  - Add form: textarea + submit button, hx-post="/comments" hx-target="#comments-list"
- Wire into all 7 detail pages: tasks, content, work_logs, knowledge, campaigns, bookmarks, companies

### 4. Inbox enhancements (3h)
- Bulk select: checkbox per item + "Gán hàng loạt" floating bar
- POST /inbox/batch-triage — accept IDs[] + destination
- Inline triage: dropdown per item row → hx-post triage without leaving page
- Detail page: edit all fields (content, url, source, item_type, company_id)
- Archive button: POST /inbox/:id/archive → set status='archived'

### 5. WorkLogs enhancements (4h)
- JOIN users: update ListWorkLogs query → include u.full_name
- Add UserName field to WorkLogItem, display in table column
- Date range filter: ?from=2026-04-01&to=2026-04-30
- User filter: ?user_id=xxx dropdown
- Work type filter: ?work_type_id=xxx dropdown
- Chart.js: add `<canvas>` element, hx-get="/work-logs/chart?month=2026-04" returns JSON
- New endpoint GET /work-logs/chart → aggregate by work_type for month → JSON for Chart.js

### 6. Tasks enhancements (4h)
- Kanban cards: show assignee name (join users on ListTasksByStatus)
- Table view toggle: query param ?view=table → render table instead of kanban
- Table view template: `web/templates/pages/tasks_table.templ`
- Date range filter: ?from=&to= on due_date
- Detail page: quantity (decimal) + unit fields (add to tasks table if not exist)
- Detail page: show campaign link if campaign_id set
- Task form: campaign dropdown (pass campaigns list)

### 7. Content enhancements (3h)
- Tag filter: query param ?topic=seo → filter by topics array contains
- Status click: stat bar items link to ?status=idea, ?status=drafting, etc.
- Author column: join users for author_id → display name
- Tags column: render topics[] as pill badges in table row

### 8. Companies enhancements (3h)
- Card: query COUNT open tasks per company, display on card
- Card: completion rate = done tasks / total tasks → percentage
- Logo upload: POST /companies/:id/logo → Drive upload → update logo_url
- Detail: display contact fields (primary_contact_name, phone, email, zalo — already in schema)
- Detail: add description/internal_notes display

### 9. Campaigns enhancements (3h)
- Filter bar: date range, company dropdown, campaign_type dropdown
- Card: show company name (join companies)
- Card: open task count (query tasks WHERE campaign_id=X AND status!='done')
- Card: completion % (done/total tasks)
- Detail: show linked work_logs list
- Detail: show linked tasks list

### 10. Knowledge enhancements (2h)
- Tag filter: query param ?topic=seo → filter topics array
- Detail: file upload button (Drive) → store URL in attachments JSONB
- Detail: assign users for read/edit (visible_to_companies + direct user assign)

### 11. Bookmarks enhancements (2h)
- Date range filter: ?from=&to= on created_at
- Tag filter: ?tag=SEO → filter tags array contains
- Domain search: ?domain=ahrefs → filter WHERE url LIKE '%ahrefs%'

## Files to Create
- `internal/db/migrations/000019_comments.up.sql` + down
- `internal/db/generated/comments.sql.go`
- `internal/handler/comments.go`
- `web/templates/components/comments.templ`
- `web/templates/pages/tasks_table.templ`

## Files to Modify
- All list handlers: add filter query params
- All detail handlers: load + pass comments
- All list templates: add filter bars
- All detail templates: add comments component
- `internal/db/generated/work_logs.sql.go` — JOIN users
- `internal/db/generated/tasks.sql.go` — JOIN users for assignee
- `cmd/server/main.go` — new routes

## Success Criteria
- [ ] Comments work on all 7 detail pages
- [ ] Inbox bulk triage + inline triage + edit + archive
- [ ] Work logs show user name + date/user/type filters + chart
- [ ] Tasks show assignee + table view + date filter
- [ ] Content filter by tag + status click + author column
- [ ] Companies show task count + completion + logo + contacts
- [ ] Campaigns filter by date/company/type + show stats
- [ ] Knowledge tag filter + file upload
- [ ] Bookmarks date/tag/domain filter
