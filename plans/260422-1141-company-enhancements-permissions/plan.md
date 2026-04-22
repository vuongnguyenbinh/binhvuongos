---
title: "Company Enhancements + Permission Matrix"
description: "Archive/logo/filter + deadline warning (badge+dashboard+notifications) cho companies; tài liệu perm matrix 3 role"
status: pending
priority: P1
effort: 5h15m
tags: [companies, notifications, permissions, ui, cron]
blockedBy: []
blocks: []
created: 2026-04-22
---

# Company Enhancements + Permission Matrix

Full context: [../reports/brainstorm-260422-1141-company-enhancements-permissions.md](../reports/brainstorm-260422-1141-company-enhancements-permissions.md)

## Phases

| # | File | Effort |
|---|---|---|
| 1 | Migration 000024: notifications dedup fields | 45m |
| 2 | Archive/unarchive handlers + queries | 30m |
| 3 | Filter (active/all/archived) + UI tabs | 30m |
| 4 | Logo upload (png/jpeg/svg, 5MB, owner+manager) | 45m |
| 5 | Deadline badge helper + render list + detail | 30m |
| 6 | Dashboard widget "Sắp hết hạn" (role-aware) | 45m |
| 7 | Deadline notifier goroutine (startup + 24h) | 45m |
| 8 | `docs/permissions.md` matrix 3-role | 45m |
| 9 | Deploy + E2E | 30m |

Do straight-through, no per-phase file — phases are small + isolated. Report in each Bash/Edit action.
