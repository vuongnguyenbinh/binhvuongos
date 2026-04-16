---
title: "Bình Vương OS — Hybrid UI-First Layer 1 (Static UI)"
description: "Convert demo HTML to Go templ components with design changes, serve as static demo for stakeholder feedback"
status: completed
priority: P1
effort: 16h
tags: [frontend, ui, templ, go, htmx]
blockedBy: []
blocks: []
created: 2026-04-16
---

# Bình Vương OS — Layer 1: Static UI

## Overview

Chuyển demo HTML (2377 dòng, 8 views) thành Go templ components. Áp dụng design changes (đỏ/gradient thay amber, ảnh stock Windows-style cho sections đặc). Serve bằng Fiber, hardcoded data. Mục đích: demo cho stakeholders, lấy feedback UX trước khi build backend.

**Layer 2-3 (DB + Business Logic) sẽ plan riêng sau khi Layer 1 được approve.**

## Design Changes (so với demo HTML gốc)

1. **Amber/Ochre → Đỏ/Đỏ cam/Gradient**
   - `amber2: #B8741F` → `ember: #D94F30` (đỏ cam chính)
   - `ochre: #C89B3C` → `flame: #E8623A` (đỏ cam sáng)
   - Gradient: `linear-gradient(135deg, #C0392B, #E74C3C, #D94F30)` cho accents
   - Warning pills, badges, notification dots → dùng ember/flame thay vì amber/ochre
   - Spark bar highlight cuối → ember thay amber

2. **Solid-color sections → Stock image + blur overlay**
   - "Chiến dịch đang chạy" (dashboard stat box bg-forest) → ảnh landscape + dark overlay 70%
   - Knowledge featured card (bg-forest) → ảnh landscape + dark overlay 70%
   - Pattern: `background: url(stock.jpg) center/cover; + overlay div rgba(31,61,46,0.75) + backdrop-blur`
   - Stock images: Unsplash landscape (mountains, lakes, fields) — 3-4 ảnh rotate

3. **Giữ nguyên**
   - Font stack: Phudu + Inter + JetBrains Mono
   - Layout: top tabs (không sidebar), max-width 1440px
   - forest (#1F3D2E), sage (#4A7C59), rust (#A64545), ivory (#FAF8F3)
   - Tất cả component styles: .pill, .eyebrow, .kanban-card, .stat-hero, v.v.

## Tech Stack (Layer 1)

- **Go 1.22+** with Fiber v2
- **templ** for HTML templates (type-safe, compiled)
- **Tailwind CSS v3** (build via CLI, not CDN)
- **HTMX** (CDN for now, tab switching)
- **Air** for hot reload dev
- **Docker** for dev environment

## Project Structure

```
binhvuongos/
├── Root/                          # Existing PRD + schema + demo
├── cmd/
│   └── server/
│       └── main.go                # Fiber server entry
├── internal/
│   └── handler/
│       ├── dashboard.go           # Route handlers
│       ├── worklogs.go
│       ├── inbox.go
│       ├── tasks.go
│       ├── content.go
│       ├── companies.go
│       ├── campaigns.go
│       └── knowledge.go
├── web/
│   ├── templates/
│   │   ├── layout.templ           # Base layout (header, tabs, footer)
│   │   ├── components/
│   │   │   ├── pill.templ         # Reusable pill badge
│   │   │   ├── progress.templ     # Progress bar
│   │   │   ├── stat-card.templ    # Stat hero card
│   │   │   ├── kanban-card.templ  # Kanban card
│   │   │   ├── check.templ        # Checkbox
│   │   │   └── hero-section.templ # Section with stock image bg
│   │   ├── pages/
│   │   │   ├── dashboard.templ
│   │   │   ├── worklogs.templ
│   │   │   ├── inbox.templ
│   │   │   ├── tasks.templ
│   │   │   ├── content.templ
│   │   │   ├── companies.templ
│   │   │   ├── campaigns.templ
│   │   │   └── knowledge.templ
│   │   └── partials/
│   │       ├── header.templ
│   │       ├── tabs.templ
│   │       └── footer.templ
│   ├── static/
│   │   ├── css/
│   │   │   ├── input.css          # Tailwind input
│   │   │   └── output.css         # Generated
│   │   ├── js/
│   │   │   └── app.js             # Tab switching, checkbox toggle
│   │   └── img/
│   │       ├── landscape-01.jpg   # Stock images for hero sections
│   │       ├── landscape-02.jpg
│   │       └── landscape-03.jpg
│   └── tailwind.config.js
├── go.mod
├── go.sum
├── Makefile                       # dev, build, templ generate, tailwind
├── Dockerfile                     # Multi-stage build
├── docker-compose.yml             # Dev environment
├── .air.toml                      # Hot reload config
└── README.md
```

## Phases

| Phase | Name | Status | Effort |
|-------|------|--------|--------|
| 1 | [Project Scaffolding & Dev Tooling](./phase-01-project-scaffolding.md) | Pending | 3h |
| 2 | [Design System & Shared Components](./phase-02-design-system-components.md) | Pending | 3h |
| 3 | [Dashboard & Work Logs Views](./phase-03-dashboard-worklogs.md) | Pending | 3h |
| 4 | [Inbox & Tasks Views](./phase-04-inbox-tasks.md) | Pending | 3h |
| 5 | [Content, Companies, Campaigns & Knowledge Views](./phase-05-remaining-views.md) | Pending | 3h |
| 6 | [Responsive Polish & Demo Deployment](./phase-06-polish-deploy.md) | Pending | 1h |

## Dependencies

- Go 1.22+ installed
- Node.js (for Tailwind CSS CLI)
- Docker & Docker Compose
- Stock landscape images (Unsplash, license-free)

## Out of Scope (Layer 1)

- Database, migrations, any backend logic
- Authentication, JWT, sessions
- Real data, API endpoints
- Telegram bot, Notion sync
- Form submissions (forms render but don't POST)
- File upload functionality
- Mobile-specific optimizations (basic responsive only)

## Next Plan (after Layer 1 approved)

- **Layer 2**: DB + Auth + Core CRUD (plan separately after stakeholder feedback)
- **Layer 3**: Business logic, integrations, Notion sync
