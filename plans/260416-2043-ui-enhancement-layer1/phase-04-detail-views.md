# Phase 4: Detail Views (6 Modules)

## Overview
- **Priority:** P1
- **Status:** Pending
- **Effort:** 4h
- Click any row/card → detail page for content, companies, tasks, work-logs, campaigns, knowledge

## Implementation Steps

### 1. Shared Detail Layout
Create `web/templates/pages/detail_layout.templ`:
```go
templ DetailLayout(module string, title string, backURL string, backLabel string) {
  <div class="mb-8">
    <a href={templ.SafeURL(backURL)} class="text-xs text-muted hover:text-ink mono">← {backLabel}</a>
    <div class="flex items-end justify-between mt-4">
      <div>
        <p class="eyebrow mb-2">{module}</p>
        <h1 class="display text-4xl md:text-5xl font-bold leading-none">{title}</h1>
      </div>
      <div class="flex items-center gap-2">
        <button class="px-3 py-1.5 bg-surface border border-hairline rounded text-xs text-muted mono dark:bg-dark-surface dark:border-dark-border">Chỉnh sửa</button>
        <button class="px-3 py-1.5 bg-rust text-white rounded text-xs mono">Xoá</button>
      </div>
    </div>
  </div>
  { children... }
}
```

### 2. Detail Pages (hardcoded demo data)

**Content Detail** (`/content/1`):
- Title, status pill, author, company badge
- Source file link (Google Docs)
- Publish info, platform, metrics (reach/engagement)
- Review notes section
- Markdown body placeholder

**Company Detail** (`/companies/1`):
- Header with logo, name, status, industry, role
- Tabs: Tasks | Content | Work Logs | Campaigns | Knowledge | People
- Each tab shows filtered list (hardcoded)

**Task Detail** (`/tasks/1`):
- Title, priority, status, assignee avatar
- Description, due date, group name
- Attachments, related campaign
- Activity timeline placeholder

**Work Log Detail** (`/work-logs/1`):
- Date, person, company, campaign
- Work type, quantity, evidence links
- Review status, admin notes
- Screenshots gallery

**Campaign Detail** (`/campaigns/1`):
- Header with progress %, date range, budget
- Progress bars per work type
- Recent work logs table
- Linked tasks and content

**Knowledge Detail** (`/knowledge/1`):
- Title, category pill, topics
- Source/author info, quality rating stars
- Body (markdown rendered or link to external)
- Shared companies list

### 3. Route Registration
In `cmd/server/main.go`, add:
```go
app.Get("/content/:id", handler.ContentDetail)
app.Get("/companies/:id", handler.CompanyDetail)
app.Get("/tasks/:id", handler.TaskDetail)
app.Get("/work-logs/:id", handler.WorkLogDetail)
app.Get("/campaigns/:id", handler.CampaignDetail)
app.Get("/knowledge/:id", handler.KnowledgeDetail)
```

### 4. Link List Items to Detail
Update list templates: wrap rows/cards in `<a href="/{module}/1">` links.

## Files to Create
- `web/templates/pages/detail_layout.templ`
- `web/templates/pages/content_detail.templ`
- `web/templates/pages/company_detail.templ`
- `web/templates/pages/task_detail.templ`
- `web/templates/pages/worklog_detail.templ`
- `web/templates/pages/campaign_detail.templ`
- `web/templates/pages/knowledge_detail.templ`
- `internal/handler/detail.go`

## Files to Modify
- `cmd/server/main.go` — 6 new routes
- All list page templates — add links to rows/cards

## Success Criteria
- Click any item → detail page with back button
- Detail pages show relevant info for each module
- Consistent layout across all 6 detail pages
- Dark mode works on detail pages
