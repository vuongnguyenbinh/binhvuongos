---
title: "Triage Convert + User CRUD + Role Simplification"
description: "Fix inbox triage to actually create target rows; consolidate 4 roles → 3; complete user CRUD with perms + password reset"
status: pending
priority: P1
effort: 5h
tags: [inbox, triage, users, roles, auth, crud]
blockedBy: []
blocks: []
created: 2026-04-21
---

# Triage Convert + User CRUD + Role Simplification

## Overview

Gộp 2 concerns sit cùng plan vì overlap (`users.role` + `middleware.RequireRole` + shared UI admin):

1. **Bug fix**: Triage từ inbox → tasks/content/knowledge hiện chỉ mark `status=done`, **không insert** vào target table. Rewrite thành transaction thật.
2. **Feature**: 4 roles cũ (owner/core_staff/ctv/staff) → 3 roles (owner/manager/staff) + full User CRUD với phân quyền + password reset (SMTP-ready).

## Context

- **Brainstorm:** [plans/reports/brainstorm-260421-0120-triage-convert-user-crud.md](../reports/brainstorm-260421-0120-triage-convert-user-crud.md)
- **Reuse:** `password_reset_tokens` (migration 000018), existing `CreateTask/CreateContent/CreateKnowledgeItem` queries, `RequireRole` middleware
- **Không đụng:** auth login flow, dashboard, webhook plan trước

## Architecture

```
Inbox row
  ├─ [→ Task]     ── HTMX → modal(triage_task)    → POST /inbox/:id/convert?target=task
  ├─ [→ Content]  ── HTMX → modal(triage_content) → POST /inbox/:id/convert?target=content
  ├─ [→ Knowledge]── HTMX → modal(triage_kb)      → POST /inbox/:id/convert?target=knowledge
  └─ [🗑 Archive]  ── POST /inbox/:id/archive (existing)

ConvertInbox handler (TX):
  INSERT INTO <target> (...) RETURNING id
  UPDATE inbox_items SET status='done', converted_to_*, processed_at=NOW()
  → redirect /inbox

Users page:
  /users [owner|manager]
    ├─ List + role-gated edit/delete buttons
    ├─ Create (role whitelist theo actor)
    ├─ /users/:id/edit → update
    ├─ /users/:id/delete → soft (status='archived')
    └─ /users/:id/reset-password → INSERT token + render link to share
```

## Phases

| # | Phase | Effort | Status |
|---|---|---|---|
| 1 | [Role migration + i18n labels](phase-01-role-migration.md) | 30m | pending |
| 2 | [User CRUD + permission middleware](phase-02-user-crud.md) | 90m | pending |
| 3 | [Password reset flow (SMTP-ready)](phase-03-password-reset.md) | 60m | pending |
| 4 | [Triage convert handler + queries](phase-04-triage-convert.md) | 60m | pending |
| 5 | [HTMX triage modals + UI wiring](phase-05-triage-modals.md) | 75m | pending |
| 6 | [Deploy + E2E tests](phase-06-deploy-test.md) | 30m | pending |

## Key dependencies

- Tasks tables có `title`, `company_id`, `priority`, `due_date`, `status`, `assignee_id` (check schema)
- Content có `title`, `company_id`, `content_type`, `status`
- Knowledge có `title`, `body`, `category`, `tags`
- `inbox_items.converted_to_type` + `converted_to_id` đã có từ migration 000005 — chưa ai dùng

## Success criteria

- Triage 1 inbox → Task: row mới ở `/tasks`, inbox marked done, `converted_to_id` point tới task ID
- Tương tự Content, Knowledge
- Migration 000023 clean: all users có role ∈ {owner, manager, staff}
- Manager tạo được user role=staff, KHÔNG tạo được role=manager/owner qua form tampering
- Owner tạo được cả 3 role
- Manager không edit/delete owner/manager khác (403)
- Reset password: click → render link URL → mở form → save pass mới
- All existing pages vẫn work (no regression)

## Risks

| Risk | Mitigation |
|---|---|
| Migration UPDATE oversweeps | Backup DB; tight WHERE; dry-run SELECT first |
| Manager privilege escalation via form | Server-side whitelist actor_role → allowed target_role, never trust form |
| Triage race (double-click) | `WHERE status != 'done'` in UPDATE; affected=0 → treat as already-processed success |
| Soft-deleted user keeps session | Check `status='active'` in AuthRequired middleware |
