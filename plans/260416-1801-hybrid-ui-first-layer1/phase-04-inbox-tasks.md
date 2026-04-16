# Phase 4: Inbox & Tasks Views

## Overview
- **Priority:** P0 — Core workflow: capture → triage → assign
- **Status:** Pending
- **Effort:** 3h
- Convert inbox view (lines 1037-1235) and tasks kanban view (lines 1241-1507) from demo HTML to templ pages

## Context Links
- Demo: `Root/binhvuong-os-demo.html` lines 1037-1235 (inbox), 1241-1507 (tasks)
- Components from Phase 2

## Files to Create

| File | Purpose |
|------|---------|
| `web/templates/pages/inbox.templ` | Inbox triage page |
| `web/templates/pages/tasks.templ` | Tasks kanban page |
| `internal/handler/inbox.go` | Inbox route handler |
| `internal/handler/tasks.go` | Tasks route handler |

## Implementation Steps

### 1. Inbox page (`pages/inbox.templ`)

Port demo lines 1037-1235. Structure:

```
InboxPage
├── Header (eyebrow + "Phân loại 12 mục chưa xử lý")
├── 2-column layout (7+5)
│   ├── Left: Inbox item list (6 items)
│   │   ├── Active item (bg-cream/30, border-l-2 forest) — Telegram
│   │   ├── Regular items (row-hover) — Manual, Zalo, Telegram, Web
│   │   └── Fading item (opacity-60) — Email, sắp archive
│   └── Right: Triage panel (sticky)
│       ├── Source badge + ID
│       ├── "Triage vào đâu?" heading
│       ├── Content preview (bg-cream)
│       ├── Destination radios (Kho kiến thức [GỢI Ý], Nội dung, Công việc, Lưu trữ)
│       ├── Company dropdown
│       └── Action buttons (CHUYỂN & TIẾP TỤC, Bỏ qua)
```

**Key design changes:**
- Source pills: TELEGRAM `bg-forest` (keep), WEB was `bg-amber2` → `bg-ember`
- Inbox item dots: was `background:#B8741F` → `background:#D94F30`
- "GỢI Ý" label: was `text-amber2` → `text-ember`
- Suggestion text: was `text-amber2 mono` → `text-ember mono`
- Radio accent: was `accent-forest` (keep for primary), `accent-rust` for delete

### 2. Tasks kanban page (`pages/tasks.templ`)

Port demo lines 1241-1507. Structure:

```
TasksPage
├── Header ("72 đầu việc" + "+ TẠO CÔNG VIỆC" button)
├── Filter bar (company, person, group, priority dropdowns)
├── View switcher (KANBAN active, BẢNG, LỊCH)
├── 5-column Kanban grid
│   ├── Cần làm (14) — 5 cards
│   ├── Đang làm (8) — 3 cards with progress bars, border-left accent
│   ├── Chờ (4) — 2 cards, opacity-80
│   ├── Cần duyệt (5) — 2 cards, border-left sage
│   └── Hoàn thành (41) — 3 cards + "...và 38 việc khác"
```

**Key design changes:**
- "Đang làm" column dot: was `background:#B8741F` → `background:#D94F30`
- Card border-left "Đang làm": was `border-left: 3px solid #B8741F` → `border-left: 3px solid #D94F30`
- Priority pill "CAO": background was `#FDF4E2` → `#FDE8E4`, color was `#B8741F` → `#D94F30`
- Avatar circles: `bg-amber2` → `bg-ember` (for certain people)
- "Chờ" column dot: was `background:#C89B3C` → `background:#E8623A`
- Time warning text: was `text-ochre` → `text-flame`

**Kanban card component usage:**
```go
// In tasks.templ, use KanbanCard component from Phase 2
@components.KanbanCard(components.KanbanCardProps{
    Priority:    "CAO",
    Company:     "ABC",
    Title:       "Xây 300 backlink cho SEO Q2",
    Subtitle:    "SEO Q2 · 187/300 link",
    Progress:    62,
    DueDate:     "30/06",
    AvatarLabel: "Hg",
    AvatarColor: "ember",
    BorderColor: "ember",
})
```

### 3. Hardcoded data

```go
// internal/handler/inbox.go
type InboxItem struct {
    Source      string // telegram, manual, zalo, web, email
    TimeAgo     string
    URL         string
    Content     string
    Attachments string
    Suggestion  string
    IsActive    bool
    IsFading    bool
}

// internal/handler/tasks.go
type KanbanColumn struct {
    Name     string
    DotColor string
    Count    int
    Cards    []TaskCard
}

type TaskCard struct {
    Priority     string
    Company      string
    Title        string
    Subtitle     string
    Progress     int    // -1 if no progress bar
    DueDate      string
    Avatar       string
    AvatarColor  string
    BorderColor  string
    IsDone       bool
    IsWaiting    bool
}
```

## Todo List
- [ ] Define inbox data structs + hardcoded demo data
- [ ] Create `pages/inbox.templ` — inbox item list (6 items)
- [ ] Create `pages/inbox.templ` — triage panel (sticky sidebar)
- [ ] Replace amber/ochre colors with ember/flame in inbox
- [ ] Create inbox handler
- [ ] Define tasks data structs + hardcoded demo data
- [ ] Create `pages/tasks.templ` — header + filter bar + view switcher
- [ ] Create `pages/tasks.templ` — 5-column kanban grid
- [ ] Replace amber/ochre colors with ember/flame in tasks
- [ ] Create tasks handler
- [ ] Verify both pages at `/inbox` and `/tasks`

## Success Criteria
- Inbox shows 6 items with correct source badges + triage panel
- Tasks kanban shows 5 columns with correct card styles
- Active inbox item highlighted with border-left forest
- Kanban "Đang làm" cards have ember border-left + progress bars
- All amber → ember, ochre → flame replacements applied
- Filter buttons and view switcher render (non-functional in Layer 1)

## Risk Assessment
- **Kanban 5-column layout**: On smaller screens, columns may overflow. Use `overflow-x-auto` wrapper.
- **Triage panel sticky**: `sticky top-32` may conflict with header height on some screens. Test.

## Next Steps
→ Phase 5: Content, Companies, Campaigns, Knowledge views
