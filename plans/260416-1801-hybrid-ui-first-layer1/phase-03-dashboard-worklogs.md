# Phase 3: Dashboard & Work Logs Views

## Overview
- **Priority:** P0 — Owner sees these first, Work Logs is most important module
- **Status:** Pending
- **Effort:** 3h
- Convert dashboard view (lines 306-691) and work logs view (lines 697-1031) from demo HTML to templ pages

## Context Links
- Demo: `Root/binhvuong-os-demo.html` lines 306-691 (dashboard), 697-1031 (worklogs)
- Layout: `web/templates/layout.templ` (Phase 2)
- Components: `web/templates/components/` (Phase 2)

## Files to Create

| File | Purpose |
|------|---------|
| `web/templates/pages/dashboard.templ` | Owner dashboard page |
| `web/templates/pages/worklogs.templ` | Work logs review page |
| `internal/handler/dashboard.go` | Dashboard route handler |
| `internal/handler/worklogs.go` | Work logs route handler |

## Implementation Steps

### 1. Dashboard page (`pages/dashboard.templ`)

Port demo lines 306-691. Structure:

```
DashboardPage
├── Greeting section (eyebrow + display heading + summary text)
├── Stat grid (4 columns, 12-col grid)
│   ├── StatCard: Output tháng này (324 link, spark bars)
│   ├── StatCard: Nội dung đã đăng (42 bài, mini bar chart)
│   ├── StatCard: Công việc (58/72, progress bar)
│   └── HeroSection: Chiến dịch đang chạy (7, stock image bg) ← DESIGN CHANGE
├── 2-column layout (8+4)
│   ├── Left: Tình trạng công ty (5 CompanyRow components)
│   └── Right: Hôm nay tasks + Inbox preview
```

**Key design changes:**
- "Chiến dịch đang chạy" box: Replace `bg-forest` with `HeroSection` component (stock landscape + forest/75 overlay)
- All `text-amber2` notification dots → `text-ember`
- `bg-amber2` avatar circles (XYZ company) → `bg-ember`
- Pill "CAO" background `#FDF4E2` → `#FDE8E4`, color `#B8741F` → `#D94F30`
- Inbox dots `background:#B8741F` → `background:#D94F30`
- Spark bar last bar `background:#B8741F` → `background:#D94F30`

**Hardcoded data approach:**
Define Go structs for demo data, return from handler:

```go
// internal/handler/dashboard.go
type DashboardData struct {
    Greeting     string
    Stats        []StatItem
    Companies    []CompanyItem
    TodayTasks   []TaskItem
    InboxPreview []InboxItem
    Campaigns    []CampaignPreview
}

func Dashboard(c *fiber.Ctx) error {
    data := getDemoDashboardData() // hardcoded
    component := pages.DashboardPage(data)
    c.Set("Content-Type", "text/html")
    return component.Render(c.Context(), c.Response().BodyWriter())
}
```

This approach makes Layer 2 transition easier — just replace `getDemoDashboardData()` with DB query.

### 2. Work Logs page (`pages/worklogs.templ`)

Port demo lines 697-1031. Structure:

```
WorkLogsPage
├── Header (eyebrow + display heading + action buttons)
├── Status tabs (Chờ duyệt 5 | Đã duyệt 142 | Cần sửa 3 | Từ chối 1)
├── Review table
│   ├── Table header (checkbox, ngày/người, công ty, loại, số lượng, nguồn, ghi chú, hành động)
│   └── 5 rows with:
│       ├── Checkbox
│       ├── Date + person name
│       ├── Company badge + campaign
│       ├── Work type (emoji + name)
│       ├── Quantity (stat-hero style)
│       ├── Evidence links (sheet, photos)
│       ├── Notes text
│       └── Action buttons (approve ✓, fix 🔄, reject ✗)
├── Daily summary footer (6 work type cards)
```

**Key design changes:**
- Status tabs: Active tab `bg-forest` (keep). Other tabs border style (keep).
- "⚠ chờ check" text: was `text-ochre` → now `text-flame`
- Highlighted row: was `bg-ochre/5` → now `bg-flame/5`
- Quantity warning: was `text-ochre` → now `text-flame`
- Action button hover (cần sửa): was `hover:bg-ochre` → now `hover:bg-flame`
- Tháng count in footer: keep `text-sage` (green is fine for positive)

**Hardcoded data:**
```go
type WorkLogEntry struct {
    Date        string
    Time        string
    PersonName  string
    CompanyCode string
    CompanyName string
    CompanyColor string
    CampaignName string
    WorkTypeIcon string
    WorkTypeName string
    Quantity     string
    Unit         string
    SheetURL     string
    PhotoCount   int
    Notes        string
    NeedsCheck   bool  // highlighted row
}
```

### 3. Route handlers

Both handlers follow same pattern:
```go
func WorkLogs(c *fiber.Ctx) error {
    data := getDemoWorkLogsData()
    component := pages.WorkLogsPage(data)
    c.Set("Content-Type", "text/html")
    return component.Render(c.Context(), c.Response().BodyWriter())
}
```

Extract common render helper to reduce duplication:
```go
// internal/handler/render.go
func render(c *fiber.Ctx, component templ.Component) error {
    c.Set("Content-Type", "text/html")
    return component.Render(c.Context(), c.Response().BodyWriter())
}
```

## Todo List
- [ ] Create `internal/handler/render.go` — shared render helper
- [ ] Define dashboard data structs in handler
- [ ] Create `pages/dashboard.templ` — greeting section
- [ ] Create `pages/dashboard.templ` — stat grid (4 cards)
- [ ] Apply HeroSection to "Chiến dịch đang chạy" card (stock image + overlay)
- [ ] Create `pages/dashboard.templ` — company health list (5 rows)
- [ ] Create `pages/dashboard.templ` — today tasks + inbox preview sidebar
- [ ] Replace all amber/ochre colors with ember/flame in dashboard
- [ ] Define work logs data structs in handler
- [ ] Create `pages/worklogs.templ` — header + status tabs
- [ ] Create `pages/worklogs.templ` — review table (5 rows)
- [ ] Create `pages/worklogs.templ` — daily summary footer (6 type cards)
- [ ] Replace all amber/ochre colors with ember/flame in worklogs
- [ ] Verify both pages render via `localhost:3000` and `localhost:3000/work-logs`

## Success Criteria
- Dashboard renders all 4 stat cards, company list, today tasks, inbox preview
- "Chiến dịch đang chạy" shows stock landscape with blurred forest overlay
- Work logs table shows 5 rows with all columns
- All amber/ochre instances replaced with ember/flame
- Both pages use Layout component (consistent header/footer)
- No broken layouts, correct spacing matches demo proportions

## Risk Assessment
- **Large templ files**: Dashboard is complex (~400 lines HTML). Split into sub-templates if >150 lines templ code.
- **Hardcoded data structs**: Keep simple — don't over-engineer for Layer 1. Plain structs, no interfaces.

## Next Steps
→ Phase 4: Inbox & Tasks views
