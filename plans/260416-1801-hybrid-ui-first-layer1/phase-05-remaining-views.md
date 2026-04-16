# Phase 5: Content, Companies, Campaigns & Knowledge Views

## Overview
- **Priority:** P1 — Complete all 8 views for full demo
- **Status:** Pending
- **Effort:** 3h
- Convert 4 remaining views from demo HTML to templ pages

## Context Links
- Demo lines: 1513-1731 (content), 1737-1940 (companies), 1946-2156 (campaigns), 2162-2321 (knowledge)
- Components from Phase 2

## Files to Create

| File | Purpose |
|------|---------|
| `web/templates/pages/content.templ` | Content pipeline page |
| `web/templates/pages/companies.templ` | Companies portfolio page |
| `web/templates/pages/campaigns.templ` | Campaigns progress page |
| `web/templates/pages/knowledge.templ` | Knowledge base page |
| `internal/handler/content.go` | Content handler |
| `internal/handler/companies.go` | Companies handler |
| `internal/handler/campaigns.go` | Campaigns handler |
| `internal/handler/knowledge.go` | Knowledge handler |

## Implementation Steps

### 1. Content Pipeline (`pages/content.templ`)

Port lines 1513-1731. Structure:
```
ContentPage
├── Header ("Pipeline nội dung" + buttons)
├── 7-column stat strip (Ý tưởng 14 | Đang viết 8 | Cần duyệt 3 | Sửa lại 2 | Duyệt 5 | Đã đăng 42 | Tổng 74)
├── Content table (6 rows)
│   └── Columns: tiêu đề, công ty, loại/nền tảng, người viết, trạng thái, ngày đăng, tương tác
```

**Design changes:**
- "Cần duyệt" stat count: was `text-amber2` → `text-ember`
- REVIEW pill: background `#FDF4E2` → `#FDE8E4`, color `#B8741F` → `#D94F30`
- REVIEW pill dot: `background:#B8741F` → `background:#D94F30`

### 2. Companies (`pages/companies.templ`)

Port lines 1737-1940. Structure:
```
CompaniesPage
├── Header ("5 công ty đang đồng hành" + button)
├── 3-column card grid
│   ├── ABC Education (OK, sage stripe)
│   ├── XYZ Finance (CHÚ Ý, flame stripe) ← was ochre
│   ├── EduPlus Academy (OK, sage stripe)
│   ├── FinHub Vietnam (GẤP, rust stripe)
│   ├── Creative Studio (PAUSE, muted, opacity-70)
│   └── "+ THÊM CÔNG TY" dashed card
```

**Design changes:**
- XYZ card: right stripe was `bg-ochre` → `bg-flame`
- XYZ avatar: was `bg-amber2` → `bg-ember`
- "CHÚ Ý" pill: bg `#FDF4E2` → `#FDE8E4`, color `#B8741F` → `#D94F30`, dot `#C89B3C` → `#E8623A`
- XYZ ontime percentage: was `text-ochre` → `text-flame`
- People avatar circles: any `bg-amber2` → `bg-ember`

### 3. Campaigns (`pages/campaigns.templ`)

Port lines 1946-2156. Structure:
```
CampaignsPage
├── Header ("7 chiến dịch đang chạy" + button)
├── Featured campaign card — SEO Q2 2026 (large, detailed)
│   ├── Header (ABC badge + title + 62% stat)
│   ├── 4 progress bars (backlink, viết bài, MXH, video)
│   ├── Team avatars
│   └── Budget (18.4M / 30M VND)
├── 4 smaller campaign cards (2-column grid)
│   ├── Ra mắt SP XYZ (34%)
│   ├── Training Q2 EDU (88%)
│   ├── Full Launch FIN (23%, rust)
│   └── Content Pillar ABC (41%)
```

**Design changes:**
- "CẦN TĂNG TỐC" label (video bar): was `text-ochre` → `text-flame`
- Video progress gradient: was `#C89B3C, #B8741F` → `#E8623A, #D94F30`
- Video percentage: was `text-ochre` → `text-flame`

### 4. Knowledge Base (`pages/knowledge.templ`)

Port lines 2162-2321. Structure:
```
KnowledgePage
├── Header ("Thư viện của Bình Vương" + search + button)
├── Category filter pills (Tất cả 184 | Bài giảng | SOP | Template | Training | Nguyên liệu | Nghiên cứu)
├── 4-column card grid
│   ├── Featured card (2-col span, HeroSection bg) ← DESIGN CHANGE: stock image
│   ├── SOP card
│   ├── Template card (3 stars)
│   ├── Nguyên liệu card (cream bg, 3 stars)
│   ├── Training card
│   ├── Bài giảng card (MỚI badge)
│   ├── SOP card
│   └── Nghiên cứu card (2 stars)
```

**Design changes:**
- Featured card: was `bg-forest` solid → `HeroSection` component (landscape-02.jpg + forest/75 overlay)
- "→ MỞ" link: was `text-amber2` → `text-ember`
- "MỚI" badge: was `text-amber2` → `text-ember`
- Star ratings: was `text-amber2` → `text-ember`
- SOP pill: bg `#FDF4E2` → `#FDE8E4`, color `#B8741F` → `#D94F30`
- WEB source pill: was `background:#B8741F` → `background:#D94F30`

### 5. Handlers (all follow same pattern)

```go
// internal/handler/content.go
func Content(c *fiber.Ctx) error {
    data := getDemoContentData()
    return render(c, pages.ContentPage(data))
}
```

Each handler defines its own data structs + hardcoded demo data function.

## Todo List
- [ ] Create content data structs + demo data
- [ ] Create `pages/content.templ` — stat strip + table
- [ ] Apply ember/flame to content REVIEW pills
- [ ] Create companies data structs + demo data
- [ ] Create `pages/companies.templ` — 3-col card grid
- [ ] Apply ember/flame to XYZ company card + CHÚ Ý pill
- [ ] Create campaigns data structs + demo data
- [ ] Create `pages/campaigns.templ` — featured card + grid
- [ ] Apply flame to "CẦN TĂNG TỐC" + video progress bar
- [ ] Create knowledge data structs + demo data
- [ ] Create `pages/knowledge.templ` — filter pills + card grid
- [ ] Apply HeroSection to featured knowledge card (landscape-02.jpg)
- [ ] Apply ember to star ratings, MỚI badge, SOP pills
- [ ] Create all 4 route handlers
- [ ] Verify all 4 pages render at correct URLs

## Success Criteria
- All 4 pages render matching demo layout proportions
- Content table: 6 rows with correct pills and status colors
- Companies: 5 cards + add card, XYZ shows flame accent
- Campaigns: featured card with 4 progress bars, 4 smaller cards
- Knowledge: featured card uses stock landscape image with blur overlay
- All amber/ochre → ember/flame replacements complete across all 4 pages
- Total: all 8 views of the app now accessible

## Risk Assessment
- **4 pages in 1 phase**: Risk of rushing. Mitigate: these pages are simpler than dashboard/worklogs (mostly tables/cards, less interactive).
- **Knowledge featured card image**: Different image from dashboard campaigns card (use landscape-02 vs landscape-01).

## Next Steps
→ Phase 6: Responsive polish + demo deployment
