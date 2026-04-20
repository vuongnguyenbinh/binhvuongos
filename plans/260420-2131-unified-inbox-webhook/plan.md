---
title: "Unified Inbox Webhook API"
description: "Single /api/v1/inbox endpoint — JSON + multipart, idempotency, Drive upload, replace Telegram handler"
status: in_progress
priority: P1
effort: 3h
completed_effort: 2h
tags: [api, webhook, inbox, integrations, n8n]
blockedBy: []
blocks: []
created: 2026-04-20
---

# Unified Inbox Webhook API

## Overview

Hợp nhất tất cả incoming integrations (Zalo OA, Telegram, Zapier, iOS Shortcut, browser extension…) qua **1 endpoint** `POST /api/v1/inbox`. Platform-specific logic (signature verify, payload transform) nằm ở n8n bridge layer, Go app chỉ còn logic core: validate → dedupe → upload attachment → insert.

Xoá hẳn `TelegramWebhook` handler cũ (bị broken do APIKeyAuth wrap).

## Context

- **Brainstorm:** [plans/reports/brainstorm-260420-2131-unified-inbox-webhook.md](../reports/brainstorm-260420-2131-unified-inbox-webhook.md)
- **Reuse:** `drive.UploadFile()`, `APIKeyAuth` middleware, `inbox_items` table
- **Bridge:** n8n self-host (Telegram Trigger + Zalo Webhook → HTTP Request to Go)

## Architecture

```
[Zalo OA][Telegram][Zapier][iOS][Ext] → n8n flows → POST /api/v1/inbox (X-API-Key)
                                                              ↓
                                                 JSON or multipart
                                                              ↓
                                         ┌──────────────────────────────┐
                                         │  Validate (item_type whitelist,│
                                         │  source, attachments ≤10)      │
                                         │  Dedupe (source, external_ref) │
                                         │  Upload file → Drive (if mp)   │
                                         │  Insert inbox_items            │
                                         └──────────────────────────────┘
                                                              ↓
                                              submitted_by = OWNER user
```

## Phases

| # | Phase | Effort | Status |
|---|---|---|---|
| 1 | [Migration + sqlc query](phase-01-migration-query.md) | 20m | ✅ done |
| 2 | [Handler refactor — JSON + multipart + idempotency](phase-02-handler-refactor.md) | 90m | ✅ done |
| 3 | [Cleanup old Telegram handler](phase-03-cleanup-telegram.md) | 20m | ✅ done |
| 4 | [Docs + n8n flow templates](phase-04-docs-n8n-flows.md) | 45m | ✅ done (main docs; n8n JSON templates deferred pending n8n URL) |
| 5 | [Deploy + E2E test](phase-05-deploy-test.md) | 30m | ⏳ pending (waiting for deploy approval) |

## Key dependencies

- Existing: `internal/handler/api_handlers.go::APICreateInbox`, `internal/drive/`, `internal/middleware/api_key.go`
- Env: `API_KEY`, `OWNER_EMAIL`, Google Drive credentials (đã có)
- External: n8n (phase 5 setup)

## Success criteria

- 1 endpoint duy nhất accept JSON + multipart
- Idempotent qua `(source, external_ref)` unique
- `submitted_by` luôn = owner user (lookup qua `OWNER_EMAIL` env)
- Telegram handler cũ bị xoá, `/api/v1/telegram/webhook` trả 404
- `docs/webhook-api.md` đầy đủ examples cho 6 nguồn (cURL, iOS, Telegram-n8n, Zalo-n8n, bookmarklet, Zapier)
- E2E: n8n flow mẫu gửi 10 messages → 10 inbox_items (no dup)

## Risks

- Owner user resolver lỗi nếu `OWNER_EMAIL` không match DB → fail startup (fail-fast OK)
- File upload 50MB timeout qua Cloudflare → 50MB hard cap handler
- Master key leak → rotate procedure trong docs
