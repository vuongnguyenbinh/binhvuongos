---
title: "Bình Vương OS — UI Enhancement (Dark Mode, Details, New Tabs)"
description: "Fix hero text, add dark mode, responsive, detail views, filters, pagination, bookmarks, prompts, markdown editor"
status: pending
priority: P1
effort: 20h
tags: [frontend, ui, dark-mode, responsive, htmx, milkdown]
blockedBy: []
blocks: []
created: 2026-04-16
---

# UI Enhancement — Layer 1 Polish

## Overview

Nâng cấp UI Layer 1 đã deploy tại https://os.binhvuong.vn. 3 nhóm: (A) fixes + dark mode + responsive, (B) detail views + filters + pagination, (C) new tabs Bookmark/Prompt + markdown editor.

## Technical Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Dark mode | Tailwind `dark:` class-based | Toggle button header, localStorage persist. Standard approach |
| Rich text | Milkdown (Markdown-native) | Output = pure Markdown → Notion/Obsidian compatible. ~40KB lazy loaded |
| Drag-drop | SortableJS + HTMX | ~10KB, touch support, trigger hx-post on drop |
| Pagination | Server-side, shared templ component | 20 items/page, numbered pages |
| Date filter | Native `<input type="date">` + HTMX | No JS library needed |

## Dark Palette

```
Light → Dark
ivory #FAF8F3 → #0F1117
surface #FFFFFF → #1A1D26
ink #1A1918 → #E8E6E3
hairline #E8E4DB → #2D3039
cream #F2EEE4 → #1F2229
muted #6B665E → #8B8880
paper bg → solid dark bg
Accent colors (forest, ember, sage, rust) → unchanged
```

## New Files to Create

```
web/templates/components/
├── date_filter.templ        # Shared date range filter
├── pagination.templ         # Shared pagination component
├── modal.templ              # Shared modal/popup wrapper
└── dark_toggle.templ        # Dark mode toggle button

web/templates/pages/
├── detail_layout.templ      # Shared detail page layout
├── inbox_detail.templ       # Inbox popup content
├── content_detail.templ     # Content detail page
├── company_detail.templ     # Company detail page
├── task_detail.templ        # Task detail page
├── worklog_detail.templ     # Work log detail page
├── campaign_detail.templ    # Campaign detail page
├── knowledge_detail.templ   # Knowledge detail page
├── bookmarks.templ          # NEW: Bookmarks page
└── prompts.templ            # NEW: Prompts page

web/static/js/
├── dark-mode.js             # Dark mode toggle logic
├── sortable-init.js         # SortableJS kanban init
└── milkdown-init.js         # Milkdown editor init (lazy)

internal/handler/
├── bookmarks.go             # Bookmarks handler
├── prompts.go               # Prompts handler
└── detail.go                # Shared detail handlers
```

## Files to Modify

```
web/templates/layout.templ          # Add dark mode toggle, dark classes
web/static/css/input.css            # Add dark: vars, fix hero overlay
tailwind.config.js                  # Add darkMode: 'class', dark colors
web/templates/pages/dashboard.templ # Dark classes, fix hero
web/templates/pages/worklogs.templ  # Dark classes, pagination, date filter
web/templates/pages/inbox.templ     # Checkboxes, create btn, modal, dark
web/templates/pages/tasks.templ     # SortableJS, dark classes
web/templates/pages/content.templ   # Dark, pagination, detail links
web/templates/pages/companies.templ # Dark, detail links
web/templates/pages/campaigns.templ # Dark, detail links
web/templates/pages/knowledge.templ # Dark, source/author, Note type, detail
web/static/js/app.js               # Import sortable, dark toggle
cmd/server/main.go                  # New routes for bookmarks, prompts, details
Dockerfile                          # Add npm deps (sortablejs, milkdown)
```

## Phases

| Phase | Name | Status | Effort |
|-------|------|--------|--------|
| 1 | [Fix Hero + Dark Mode](./phase-01-darkmode-fixes.md) | Pending | 4h |
| 2 | [Mobile Responsive](./phase-02-responsive.md) | Pending | 2h |
| 3 | [Shared Components (Filter, Pagination, Modal)](./phase-03-shared-components.md) | Pending | 3h |
| 4 | [Detail Views (6 modules)](./phase-04-detail-views.md) | Pending | 4h |
| 5 | [Inbox + Kanban Enhancements](./phase-05-inbox-kanban.md) | Pending | 3h |
| 6 | [Bookmarks + Prompts Tabs](./phase-06-new-tabs.md) | Pending | 2h |
| 7 | [Milkdown Editor + Deploy](./phase-07-editor-deploy.md) | Pending | 2h |

## Dependencies

- SortableJS (CDN or npm): kanban drag-drop
- Milkdown (npm): markdown editor
- No backend changes needed (still Layer 1 static)
