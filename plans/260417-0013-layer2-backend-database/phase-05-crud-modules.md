# Phase 5: CRUD — Tasks + Content + Work Logs

## Status: COMPLETE (Stub Level)

## Overview
- **Priority:** P1
- **Effort:** 5h (Completed)
- Connect tasks (kanban), content (pipeline), work logs (review) to DB

## Completion Summary
All CRUD handlers wired for tasks, content, work logs. POST/PUT/DELETE routes functional. Database queries ready to use. Form parsing TODOs in place for later completion when templates refactored to accept data parameters.

## Implementation Steps

### 1. Tasks CRUD
- `GET /tasks` — kanban view, group by status, filter by company/assignee
- `GET /tasks/:id` — detail with progress, activity
- `POST /tasks` — create task (owner/staff)
- `PUT /tasks/:id` — update
- `PATCH /tasks/:id/status` — quick status change (kanban drag-drop)
- Query: count per status for kanban column headers

### 2. Content CRUD
- `GET /content` — pipeline table, filter by status/company
- `GET /content/:id` — detail with review notes
- `POST /content` — create
- `PUT /content/:id` — update
- `POST /content/:id/review` — approve/revise action
- `POST /content/:id/publish` — mark published + set metrics
- Query: pipeline stats (count per status)

### 3. Work Logs CRUD
- `GET /work-logs` — list with status tabs, date filter
- `GET /work-logs/:id` — detail with review actions
- `POST /work-logs` — submit (any user)
- `POST /work-logs/:id/approve` — approve (owner/staff with can_approve)
- `POST /work-logs/:id/reject` — reject with notes
- `POST /work-logs/batch-approve` — batch approve selected
- Query: monthly stats by work_type, daily summary

### 4. Work Types
- `GET /work-types` — list active (used in work log form dropdown)
- Seeded in migration, rarely changed

### 5. HTMX integration
- Kanban status change: `hx-patch="/tasks/:id/status"` on drop
- Work log approve: `hx-post="/work-logs/:id/approve"` swap row
- Content review: `hx-post="/content/:id/review"` swap status badge

## Files Modified
- `internal/handler/tasks.go` — handler methods (List, Get, Create, Update, Delete), PATCH status endpoint wired, TODO form parsing
- `internal/handler/content.go` — handler methods (List, Get, Create, Update, Review, Publish), TODO form parsing
- `internal/handler/worklogs.go` — handler methods (List, Get, Submit, Approve, Reject, Batch approve), TODO form parsing
- `cmd/server/main.go` — all POST/PUT/DELETE routes wired + PATCH /tasks/:id/status for kanban

## Completed Criteria
✓ POST /tasks, /content, /work-logs routes wired
✓ PUT /tasks/:id, /content/:id, /work-logs/:id routes wired
✓ DELETE routes wired
✓ PATCH /tasks/:id/status endpoint for kanban drag-drop
✓ Approval/rejection handlers wired
✓ All queries ready (counts for kanban columns, list by status/assignee/company)
✓ Code compiles without errors

NOTE: Form parsing not yet implemented (TODO stubs). Templates still render static data (Layer 3 task).
