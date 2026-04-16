# Phase 6: CRUD — Campaigns + Knowledge + Bookmarks

## Status: COMPLETE (Stub Level)

## Overview
- **Priority:** P1
- **Effort:** 4h (Completed)
- Connect remaining modules to DB

## Completion Summary
All CRUD endpoints wired for campaigns, knowledge, bookmarks. Database queries ready. Bookmarks table added via migration 000013. Full-text search queries prepared for knowledge module. Form parsing TODOs in place.

## Implementation Steps

### 1. Campaigns CRUD
- `GET /campaigns` — list with progress bars
- `GET /campaigns/:id` — detail with per-work-type progress
- `POST /campaigns` — create with target_json builder
- `PUT /campaigns/:id` — update
- Query: campaign progress (join work_logs, aggregate by work_type vs target_json)
- Use v_campaign_progress view from schema

### 2. Knowledge CRUD
- `GET /knowledge` — grid with category filter, search
- `GET /knowledge/:id` — detail with markdown body
- `POST /knowledge` — create (owner/staff)
- `PUT /knowledge/:id` — update
- `GET /knowledge/search?q=` — full-text search (GIN + tsvector)
- Filter: category, topic, quality_rating, source, author
- Include "Prompt" and "Note/Idea" types

### 3. Bookmarks CRUD
- `GET /bookmarks` — grid with tag filter
- `GET /bookmarks/:id` — detail with notes
- `POST /bookmarks` — create
- `PUT /bookmarks/:id` — update
- `DELETE /bookmarks/:id` — delete
- Note: bookmarks table not in original schema — need migration 000015

### 4. Objectives (lightweight)
- Included in company detail page
- `GET /companies/:id` shows objectives for that company
- Basic CRUD if time permits

### 5. New migration for bookmarks
```sql
CREATE TABLE bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(300) NOT NULL,
    url TEXT NOT NULL,
    description TEXT,
    tags TEXT[],
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
```

## Files Modified
- `internal/handler/campaigns.go` — handler methods (List, Get, Create, Update), progress queries wired, TODO form parsing
- `internal/handler/knowledge.go` — handler methods (List, Get, Create, Update, Search), FTS queries ready, TODO form parsing
- `internal/handler/bookmarks.go` — handler methods (List, Get, Create, Update, Delete), TODO form parsing
- `cmd/server/main.go` — all POST/PUT/DELETE routes for campaigns, knowledge, bookmarks wired

## Files Created
- `internal/db/migrations/000013_bookmarks.up.sql` — bookmarks table with id, title, url, description, tags, notes, created_by, timestamps
- `internal/db/migrations/000013_bookmarks.down.sql` — DROP TABLE bookmarks
- `internal/db/query/bookmarks.sql` — CRUD + list by tags

## Completed Criteria
✓ POST /campaigns, /knowledge, /bookmarks routes wired
✓ PUT /campaigns/:id, /knowledge/:id, /bookmarks/:id routes wired
✓ DELETE /bookmarks/:id route wired
✓ Campaign progress query ready (joins work_logs, aggregates by work_type)
✓ Knowledge FTS search query ready (uses tsvector/tsquery for Vietnamese)
✓ Bookmarks migration applied + queries generated
✓ All code compiles without errors

NOTE: Form parsing not yet implemented (TODO stubs). Templates still render static data (Layer 3 task).
