---
type: brainstorm
date: 2026-04-20 21:31 +07
slug: unified-inbox-webhook
status: approved
---

# Brainstorm: Unified Inbox Webhook API

## Problem statement

Hiện có nhiều nền tảng (Zalo OA, Telegram Bot, iOS Shortcut, Zapier, browser extension…) muốn push thông tin vào `inbox_items`. Không muốn mỗi nền tảng một endpoint riêng trong Go app — handler Telegram hiện tại đã broken (APIKeyAuth nhưng Telegram không gửi X-API-Key).

Cần: **1 webhook duy nhất**, auth đơn giản, hỗ trợ cả text lẫn file, idempotent.

## Requirements (đã chốt)

- Auth: **1 master API key** duy nhất (`X-API-Key` header)
- `submitted_by`: luôn gán **owner user** (solo founder, chưa multi-tenant)
- Attachments: **hỗ trợ multipart upload** → Google Drive (reuse `drive.UploadFile`)
- Bridge layer: **n8n** xử lý mọi platform-specific (Telegram verify, Zalo signature, payload transform)
- Go app: **bỏ hẳn** `TelegramWebhook` handler cũ, chỉ còn 1 endpoint unified
- `item_type` strict **whitelist**: `note | link | image | voice | file`
- Max **10 attachments** per item
- Idempotency qua `(source, external_ref)` unique index

## Evaluated approaches

### A. Platform-specific handlers trong Go app (rejected)
Mỗi nền tảng 1 route riêng + signature verify trong Go.
- Pros: không phụ thuộc bridge
- Cons: Go app phình ra mỗi khi thêm platform; logic verify duplicate; maintain khó

### B. Unified webhook + external bridge (n8n) ✅ CHOSEN
Go app 1 endpoint. n8n xử lý Telegram/Zalo webhook → transform → POST vào Go.
- Pros: Go code tối giản (YAGNI); thêm platform mới = thêm n8n flow, không touch Go
- Cons: dependency vào n8n uptime; cần self-host n8n
- Mitigation: n8n self-host trên cùng server 103.97.125.186, có queue retry built-in

### C. Cloudflare Worker bridge (alternative)
Worker serverless thay n8n.
- Pros: zero-maintain, auto-scale
- Cons: phải code mỗi platform bằng JS; không có UI flow builder; debug khó
- Decision: n8n thắng vì UX tốt hơn cho non-technical maintenance

## Final recommended solution

### Endpoint

```
POST https://os.binhvuong.vn/api/v1/inbox
Header: X-API-Key: <master_key>
Content-Type: application/json | multipart/form-data
```

### JSON payload

```json
{
  "content": "string (required, max 10KB)",
  "url": "string (optional)",
  "source": "string ≤30 (required: zalo|telegram|zapier|ios_shortcut|extension|manual|api)",
  "item_type": "note|link|image|voice|file (optional, default auto-detect)",
  "attachment_urls": ["string (optional, max 10)"],
  "external_ref": "string ≤200 (optional, idempotency key)"
}
```

### Multipart (with file)

Form fields: `content`, `source`, `item_type`, `external_ref` + `file` (binary, max 50MB)
→ Handler upload lên Google Drive qua `drive.UploadFile` → URL push vào `attachments[]`

### Validation

| Field | Rule |
|---|---|
| `content` | non-empty, ≤10KB |
| `source` | non-empty, ≤30 chars |
| `item_type` | whitelist: `note`, `link`, `image`, `voice`, `file`; default `note` |
| `attachment_urls` | max 10, mỗi URL ≤2KB |
| `external_ref` | optional; nếu có, check unique với source |
| Total body | JSON ≤10MB, multipart ≤50MB |

### Database

Reuse `inbox_items` table. 1 migration nhỏ:

```sql
-- 000022_inbox_external_ref.up.sql
ALTER TABLE inbox_items ADD COLUMN external_ref VARCHAR(200);
CREATE UNIQUE INDEX idx_inbox_external_ref
  ON inbox_items(source, external_ref)
  WHERE external_ref IS NOT NULL;
```

Attachments structure (JSONB):
```json
[{ "url": "...", "filename": "...", "mime": "image/jpeg", "drive_file_id": "..." }]
```

### Auth key generation

```bash
# Generate
openssl rand -hex 32   # 64-char hex

# Deploy
ssh root@103.97.125.186 "echo 'API_KEY=<new_key>' >> /opt/binhvuongos/.env && \
  cd /opt/binhvuongos && docker compose restart app"
```

Rotate = regenerate + update `.env` + restart (~3s downtime).

### Response

**201 Created:**
```json
{ "success": true, "data": { "id": "uuid", "content": "...", "source": "...", "attachments": [...], "status": "raw", "created_at": "..." } }
```

**409 Conflict** (duplicate external_ref):
```json
{ "success": true, "data": { /* existing item */ }, "duplicate": true }
```

**4xx errors:** `{"success": false, "error": "...", "code": "VALIDATION_ERROR"}`

## Implementation scope

| Task | LOC | File |
|---|---|---|
| Extend `APICreateInbox` — JSON + multipart branching | ~40 | `internal/handler/api_handlers.go` |
| Helper: upload file → Drive → attachment struct | ~30 | `internal/handler/api_handlers.go` |
| Validate item_type whitelist + attachment count | ~20 | same |
| Idempotency: query by (source, external_ref) trước INSERT | ~25 | same |
| Migration `000022_inbox_external_ref` | 10 SQL | `internal/db/migrations/` |
| sqlc query `GetInboxByExternalRef` | ~8 | `internal/db/query/inbox_items.sql` |
| Delete `TelegramWebhook` handler + route | −60 | `integrations.go`, `main.go` |
| Owner user ID resolver (env `OWNER_EMAIL` → lookup) | ~15 | startup hoặc handler |
| `docs/webhook-api.md` usage guide | markdown | docs/ |
| n8n flow templates (Telegram, Zalo) | JSON | docs/n8n/ |

**Total:** ~130 LOC Go + 1 migration + docs. **Effort:** ~3h work.

## Usage guide summary

Chi tiết trong `docs/webhook-api.md` (sẽ tạo ở plan):

- **cURL:** JSON + multipart examples
- **iOS Shortcut:** recipe share-sheet → POST
- **Telegram Bot (via n8n):** Telegram Trigger node → Function transform → HTTP POST
- **Zalo OA (via n8n):** Webhook Trigger → parse event → (nếu image) Zalo Media API → HTTP POST
- **Browser extension / bookmarklet:** JS fetch snippet
- **Zapier / Make:** Webhooks action config

## Risks & mitigations

| Risk | Mitigation |
|---|---|
| Master key leak | Rotate dễ; Cloudflare IP allowlist backup; env var không log |
| n8n down → mất Telegram msg | n8n self-host + queue retry; monitoring systemd |
| Duplicate insert khi n8n retry | `external_ref` unique index |
| Body size vượt CF free 100MB | 50MB hard limit handler; error rõ ràng |
| Owner user ID thay đổi | Resolver dùng `OWNER_EMAIL` env → lookup runtime |

## Success metrics

- n8n flow gửi 10 Telegram msg → 10 `inbox_items` (source=`telegram`, no dup)
- Zalo image qua n8n → `inbox_item` có `attachments[0].url` = Drive link, `item_type=image`
- Rotate key → old trả 401, new trả 201 trong <5s sau restart
- Telegram cũ route `/api/v1/telegram/webhook` trả 404 (đã xoá)
- `docs/webhook-api.md` đủ để người ngoài setup được 1 source mới trong <30 phút

## Next steps

Tạo plan `260420-2131-unified-inbox-webhook` với các phase:
1. Migration + sqlc queries
2. Handler refactor (JSON + multipart + idempotency + validation)
3. Xoá Telegram handler cũ
4. Docs + n8n flow templates
5. Deploy + test end-to-end với cURL + n8n sample flow

## Unresolved questions

- n8n deploy ở đâu? Cùng server Docker `103.97.125.186` (có port rảnh) hay VPS khác? → để phase 5 quyết.
- Có cần rate-limit khác cho webhook endpoint (hiện 60/phút) hay để như global? → đề xuất giữ 60/phút, reassess sau.
