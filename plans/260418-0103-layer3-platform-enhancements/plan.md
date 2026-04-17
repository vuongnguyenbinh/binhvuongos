---
title: "Bình Vương OS — Layer 3: Platform Enhancements"
description: "Core platform upgrades, comments system, enhanced CRUD, rich text, Notion sync"
status: pending
priority: P1
effort: 72h
tags: [platform, ux, comments, crud, richtext, notion, integrations]
blockedBy: []
blocks: []
created: 2026-04-18
---

# Layer 3: Platform Enhancements

## Overview

Nâng cấp Bình Vương OS từ MVP functional lên production-grade UX: comments system xuyên suốt, enhanced filtering/display cho tất cả modules, rich text editor, real Notion sync, admin settings, notifications, Google OAuth, chat bubble.

## Architecture Changes

```
New DB tables: comments, notifications, settings, user_notes, password_reset_tokens
New column: content.body TEXT
New JS libs: SimpleMDE (markdown), Chart.js (charts)
New Go libs: goldmark (markdown render), golang.org/x/oauth2
New endpoints: ~30 routes
```

## Phases

| Phase | Name | Status | Effort |
|-------|------|--------|--------|
| 1 | [Core Platform](./phase-01-core-platform.md) | Pending | 29h |
| 2 | [Comments + Enhanced CRUD](./phase-02-comments-enhanced-crud.md) | Pending | 29.5h |
| 3 | [Rich Text + Integrations](./phase-03-richtext-integrations.md) | Pending | 14h |

## Dependencies

- Layer 2 backend (completed functionally)
- Google OAuth credentials (from Admin Settings)
- SMTP credentials (from Admin Settings)
- Notion API key + database IDs (from Admin Settings)
- n8n webhook URL (from Admin Settings)
