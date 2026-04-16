# Phase 6: Bookmarks + Prompts Tabs

## Overview
- **Priority:** P2
- **Status:** Pending
- **Effort:** 2h
- Two new pages: /bookmarks and /prompts

## Implementation Steps

### 1. Bookmarks Page (`/bookmarks`)
Layout:
```
Header: "Bookmark" + search + "Thêm link" button
Filter: Tags (tất cả | SEO | Content | Tools | Design | Business)
Grid: 3-col cards
  Each card:
    - Favicon (from URL domain)
    - Title
    - URL (truncated, mono)
    - Description (1-2 lines)
    - Tags (pills)
    - Date added
```

Hardcoded 8-10 sample bookmarks (SEO tools, design resources, business articles).

### 2. Prompts Page (`/prompts`)
Layout:
```
Header: "Prompt Templates" + search + "Tạo prompt" button
Filter: Categories (Tất cả | SEO | Content | AI | Copywriting | Other)
List: Cards with expandable content
  Each card:
    - Title
    - Category pill
    - Preview (first 100 chars)
    - "Copy" button (clipboard)
    - Click → expand full prompt text
```

Hardcoded 6-8 sample prompts (SEO meta description, blog outline, social caption, etc.)

### 3. Navigation Update
In `layout.templ` nav tabs, add:
- Bookmark (icon 🔖)
- Prompt (icon 💬)

### 4. Route Registration
```go
app.Get("/bookmarks", handler.Bookmarks)
app.Get("/prompts", handler.Prompts)
```

### 5. Knowledge Enhancement
In `knowledge.templ`:
- Add "Note/Idea" pill type to filter bar
- Add source + author text under each card
- Source/author filter in filter bar

## Files to Create
- `web/templates/pages/bookmarks.templ`
- `web/templates/pages/prompts.templ`
- `internal/handler/bookmarks.go`
- `internal/handler/prompts.go`

## Files to Modify
- `web/templates/layout.templ` — add 2 nav tabs
- `cmd/server/main.go` — add 2 routes
- `web/templates/pages/knowledge.templ` — add Note type, source/author

## Success Criteria
- /bookmarks shows grid of link cards
- /prompts shows expandable prompt cards
- Copy button copies prompt text to clipboard
- Knowledge has Note/Idea filter + source/author info
- Both tabs in dark mode
