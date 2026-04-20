# Phase 4 — Docs + n8n Flow Templates

**Effort:** 45m | **Priority:** P1 | **Status:** pending | **Depends on:** Phase 2

## Context
- Docs directory: `docs/`
- Existing docs (from CLAUDE.md pattern): `project-overview-pdr.md`, `code-standards.md`, etc.

## Overview
Viết `docs/webhook-api.md` — hướng dẫn đầy đủ cho 6 nguồn: cURL, iOS Shortcut, Telegram via n8n, Zalo via n8n, bookmarklet, Zapier. Kèm JSON template cho n8n flows (import-ready).

## Files to create
- `docs/webhook-api.md` — main guide (~300-400 LOC markdown)
- `docs/n8n-flows/telegram-to-inbox.json` — n8n flow export
- `docs/n8n-flows/zalo-to-inbox.json` — n8n flow export
- `docs/n8n-flows/README.md` — import instructions

## Implementation steps

### 1. Write `docs/webhook-api.md` structure

```markdown
# Unified Inbox Webhook API

## Endpoint
POST https://os.binhvuong.vn/api/v1/inbox
Header: X-API-Key: <master_key>

## Auth — generate & rotate key
openssl rand -hex 32 → update .env on server → docker compose restart app

## Payload formats
### JSON (text/URL/pre-uploaded attachments)
{ ...spec... }

### Multipart (file upload)
{ ...spec... }

## Validation rules
Table: field | required | constraint

## Response codes
- 201 Created
- 200 + duplicate:true (idempotent hit)
- 400 VALIDATION_ERROR / UPLOAD_ERROR
- 401 Invalid API key
- 500 DB_ERROR

## Source recipes
### 1. cURL
### 2. iOS Shortcut
### 3. Telegram Bot (via n8n)
### 4. Zalo OA (via n8n)
### 5. Browser bookmarklet
### 6. Zapier / Make / IFTTT

## Troubleshooting
- 401 → check X-API-Key header
- 400 item_type → whitelist: note|link|image|voice|file
- Duplicate → verify external_ref unique per source

## Rate limits
60 req/min per IP (global limiter)
```

### 2. n8n flow — Telegram to Inbox

`docs/n8n-flows/telegram-to-inbox.json` (n8n export format):

```json
{
  "name": "Telegram → Inbox",
  "nodes": [
    {
      "name": "Telegram Trigger",
      "type": "n8n-nodes-base.telegramTrigger",
      "parameters": { "updates": ["message"] }
    },
    {
      "name": "Transform",
      "type": "n8n-nodes-base.function",
      "parameters": {
        "functionCode": "const m = items[0].json.message; return [{json: { content: m.text || m.caption || '', source: 'telegram', external_ref: `tg:${m.message_id}`, item_type: (m.text||'').startsWith('http') ? 'link' : 'note' }}];"
      }
    },
    {
      "name": "POST to Inbox",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "https://os.binhvuong.vn/api/v1/inbox",
        "headerParameters": { "parameters": [{"name":"X-API-Key","value":"={{$env.BVOS_API_KEY}}"}] },
        "bodyParameters": "={{$json}}"
      }
    }
  ],
  "connections": { "Telegram Trigger": {"main":[[{"node":"Transform","type":"main","index":0}]]}, "Transform": {"main":[[{"node":"POST to Inbox","type":"main","index":0}]]} }
}
```

### 3. n8n flow — Zalo OA to Inbox

`docs/n8n-flows/zalo-to-inbox.json`:
- Webhook Trigger (path `/zalo-webhook`)
- Function: parse Zalo event, handle `user_send_text` + `user_send_image`
- (If image) Zalo Media API node → get URL
- HTTP Request → POST /api/v1/inbox

### 4. `docs/n8n-flows/README.md`
- How to import JSON flow into n8n
- Set env vars: `BVOS_API_KEY`, Telegram BOT_TOKEN, Zalo OA_SECRET
- Activate workflows

## Todo
- [ ] Write `docs/webhook-api.md` with all 6 recipes
- [ ] Include cURL examples (test from local trước khi paste)
- [ ] Include iOS Shortcut step-by-step
- [ ] Export n8n Telegram flow JSON
- [ ] Export n8n Zalo flow JSON
- [ ] Write `docs/n8n-flows/README.md` import guide
- [ ] Verify rate limit + security notes accurate

## Success criteria
- `docs/webhook-api.md` standalone readable — không cần context khác
- Người ngoài setup được 1 source mới <30 phút theo hướng dẫn
- n8n flow JSON import thành công vào n8n instance fresh
- cURL examples copy-paste chạy được ngay (với `API_KEY` env)

## Risks
- Zalo OA API có thể thay đổi — docs cần version/date
- n8n node version conflict — note version tested (n8n 1.x)
