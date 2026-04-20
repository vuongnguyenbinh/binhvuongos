# Unified Inbox Webhook API

> Gửi bất kỳ tin nhắn, link, ảnh, file từ mọi nền tảng (Zalo, Telegram, Zapier, iOS Shortcut, browser extension…) vào `inbox_items` của Bình Vương OS qua **1 endpoint duy nhất**.

## Endpoint

```
POST https://os.binhvuong.vn/api/v1/inbox
```

**Auth header:** `X-API-Key: <master_api_key>`

**Content-Type:** `application/json` hoặc `multipart/form-data`

## Generate & Rotate master API key

```bash
# Tạo key mới (64 ký tự hex, 256-bit entropy)
openssl rand -hex 32

# Update trên server
ssh root@103.97.125.186 << 'EOF'
cd /opt/binhvuongos
sed -i '/^API_KEY=/d' .env
echo "API_KEY=<new_key>" >> .env
grep -q OWNER_EMAIL .env || echo "OWNER_EMAIL=vuongnguyenbinh@gmail.com" >> .env
docker compose restart app
EOF
```

Rotate key = generate mới + update `.env` + restart (~3s downtime). Client cũ trả 401 ngay.

## Payload — JSON

```json
{
  "content": "string (required, max 10KB)",
  "url": "string (optional)",
  "source": "string ≤30 chars (required)",
  "item_type": "note | link | image | voice | file (optional, auto-detect nếu bỏ trống)",
  "attachment_urls": ["https://... (optional, max 10)"],
  "external_ref": "string ≤200 (optional, dùng cho idempotency)"
}
```

## Payload — Multipart (upload file)

| Field | Type | Ghi chú |
|---|---|---|
| `content` | string | Required |
| `source` | string | Required, ≤30 chars |
| `item_type` | string | Optional whitelist |
| `external_ref` | string | Optional |
| `attachment_urls` | repeated string | Optional, max 10 |
| `file` | binary | Optional, max 50MB, auto-upload lên Google Drive |

## Validation rules

| Field | Rule |
|---|---|
| `content` | non-empty, ≤10KB |
| `source` | non-empty, ≤30 chars |
| `item_type` | whitelist: `note`, `link`, `image`, `voice`, `file` |
| `external_ref` | ≤200 chars; unique per `(source, external_ref)` |
| `attachment_urls` | ≤10 URL, mỗi URL ≤2048 chars |
| `file` (multipart) | ≤50MB; cần Google Drive credentials đã cấu hình |

## Response

### 201 Created — inbox_item vừa tạo

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "content": "...",
    "url": "...",
    "source": "zalo",
    "item_type": "image",
    "status": "raw",
    "attachments": [{"url":"...", "filename":"...", "mime":"image/jpeg", "drive_file_id":"..."}],
    "external_ref": "zalo:msg:123",
    "submitted_by": "uuid",
    "created_at": "2026-04-20T16:27:00+07:00",
    "updated_at": "2026-04-20T16:27:00+07:00"
  }
}
```

### 200 OK + `duplicate: true` — hit idempotency

```json
{ "success": true, "duplicate": true, "data": { /* existing row */ } }
```

### 400 Bad Request

```json
{ "success": false, "error": { "code": "VALIDATION_ERROR", "message": "content is required" } }
```

Error codes: `VALIDATION_ERROR`, `UPLOAD_ERROR`

### 401 Unauthorized

```json
{ "success": false, "error": "Invalid or missing API key" }
```

### 500 Internal

```json
{ "success": false, "error": { "code": "DB_ERROR", "message": "..." } }
```

## Recipes

### 1. cURL

```bash
export API_KEY=your_master_key_here
export URL=https://os.binhvuong.vn/api/v1/inbox

# a) Note text
curl -X POST $URL -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"content":"Idea cho landing page","source":"manual"}'

# b) Link (auto-detect item_type)
curl -X POST $URL -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"content":"Ahrefs blog","url":"https://ahrefs.com/blog","source":"manual"}'

# c) Idempotent (retry-safe)
curl -X POST $URL -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"content":"x","source":"zapier","external_ref":"zap_run_12345"}'

# d) Multipart file upload
curl -X POST $URL -H "X-API-Key: $API_KEY" \
  -F "content=Screenshot từ client" \
  -F "source=zalo" \
  -F "item_type=image" \
  -F "file=@/path/to/screenshot.png"

# e) Pre-uploaded attachments
curl -X POST $URL -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "content":"Meeting notes",
    "source":"drive",
    "item_type":"file",
    "attachment_urls":["https://drive.google.com/file/d/abc/view"]
  }'
```

### 2. iOS Shortcut

1. Mở app **Shortcuts** → New Shortcut → Add Action
2. Chọn **Get Contents of URL**
3. URL: `https://os.binhvuong.vn/api/v1/inbox`
4. Method: `POST`
5. Headers:
   - `X-API-Key` = `<paste_master_key>`
   - `Content-Type` = `application/json`
6. Request Body → JSON:
   ```json
   { "content": "Shortcut Input", "source": "ios_shortcut" }
   ```
   Thay `Shortcut Input` bằng token magic variable.
7. **Add to Share Sheet** → Accept types: Text, URLs.
8. Test: chia sẻ text/URL bất kỳ từ Safari → xuất hiện trong `/inbox`.

### 3. Telegram Bot (qua n8n)

**Yêu cầu:** n8n đã chạy (self-host hoặc cloud), Telegram Bot token từ [@BotFather](https://t.me/BotFather).

1. n8n → New workflow → **Telegram Trigger** node
   - Credential: nhập BOT_TOKEN
   - Updates: `message`
2. **Function** node — transform payload:
   ```js
   const m = items[0].json.message;
   const text = m.text || m.caption || '';
   return [{
     json: {
       content: text,
       source: 'telegram',
       external_ref: `tg:${m.message_id}`,
       item_type: text.startsWith('http') ? 'link' : 'note'
     }
   }];
   ```
3. **HTTP Request** node:
   - Method: `POST`
   - URL: `https://os.binhvuong.vn/api/v1/inbox`
   - Headers: `X-API-Key: {{$env.BVOS_API_KEY}}`, `Content-Type: application/json`
   - Body: `={{$json}}`
4. Activate workflow. Tin nhắn Telegram → tự động vào `inbox_items`.

> **File đính kèm Telegram** (photo/voice): thêm Function node parse `m.photo[]` hoặc `m.voice`, gọi Telegram `getFile` API, download về, POST lại endpoint dạng multipart.

### 4. Zalo OA (qua n8n)

**Yêu cầu:** Zalo Official Account đã có, OA secret + app_id, webhook URL trỏ về n8n.

1. n8n → **Webhook Trigger** node
   - Path: `/zalo-oa`
   - HTTP method: POST
   - Đăng ký URL này với Zalo OA webhook config
2. **Function** node:
   ```js
   const ev = items[0].json;
   const isImage = ev.event_name === 'user_send_image';
   return [{
     json: {
       content: ev.message?.text || (isImage ? '[Image]' : ''),
       source: 'zalo',
       external_ref: `zalo:${ev.message?.msg_id}`,
       item_type: isImage ? 'image' : 'note',
       attachment_urls: isImage ? ev.message.attachments?.map(a => a.payload.url) : []
     }
   }];
   ```
3. **HTTP Request** node:
   - Method: `POST`, URL: `https://os.binhvuong.vn/api/v1/inbox`
   - Headers: `X-API-Key: {{$env.BVOS_API_KEY}}`
   - Body: `={{$json}}`
4. Activate.

> **Verify Zalo signature:** thêm Function node check HMAC-SHA256 của body với OA secret trước khi pass.

### 5. Browser bookmarklet

Paste vào bookmark URL, click trên bất kỳ page nào:

```javascript
javascript:(async()=>{const r=await fetch('https://os.binhvuong.vn/api/v1/inbox',{method:'POST',headers:{'X-API-Key':'YOUR_KEY','Content-Type':'application/json'},body:JSON.stringify({content:document.title,url:location.href,source:'bookmarklet',item_type:'link'})});alert(r.ok?'✓ Saved':'✗ Failed: '+r.status)})();
```

Thay `YOUR_KEY` bằng master API key.

### 6. Zapier / Make / IFTTT

**Zapier:**

1. Trigger: event tuỳ chọn (Gmail label, Airtable row, etc.)
2. Action → **Webhooks by Zapier** → POST
3. URL: `https://os.binhvuong.vn/api/v1/inbox`
4. Payload type: `json`
5. Data:
   - `content`: `{{trigger_field}}`
   - `source`: `zapier`
   - `external_ref`: `{{zap_step_id}}` (cho idempotency)
6. Headers:
   - `X-API-Key`: `your_key`
   - `Content-Type`: `application/json`

**Make (Integromat)** và **IFTTT Webhooks** config tương tự.

## Rate limits

- Global: 60 req/phút per IP (Fiber `limiter` middleware)
- Body size: JSON ≤10MB, multipart ≤50MB
- Khuyến nghị: bridge (n8n) batch nếu nguồn gửi spike

## Security notes

- Master key không bao giờ log vào server log (middleware không in header)
- Rotate key nếu nghi lộ → restart <5s
- Cloudflare Tunnel → HTTPS enforced (no plain HTTP)
- Backup: add Cloudflare IP allowlist ở Zero Trust nếu cần hạn chế

## Troubleshooting

| Triệu chứng | Nguyên nhân | Fix |
|---|---|---|
| 401 Invalid API key | Sai hoặc thiếu `X-API-Key` | Check header, verify `.env API_KEY` |
| 400 `invalid item_type` | `item_type` ngoài whitelist | Dùng `note`, `link`, `image`, `voice`, `file` |
| 400 `content is required` | Field `content` rỗng | Ensure client gửi content non-empty |
| 400 `too many attachments` | >10 URL trong `attachment_urls` | Giảm xuống ≤10 |
| 400 UPLOAD_ERROR | File >50MB hoặc Drive chưa config | Giảm size hoặc check `GOOGLE_REFRESH_TOKEN` env |
| 500 DB_ERROR | Unique constraint / connection | Check logs; nếu duplicate `(source, external_ref)` thì client nên nhận 200 |
| Always `source=api` in DB | Client không gửi `source` | Set explicit source mỗi request |

## Schema — `inbox_items` (reference)

| Column | Type | Note |
|---|---|---|
| `id` | UUID | PK |
| `content` | TEXT | message body |
| `url` | TEXT | optional, auto từ `content` nếu là URL |
| `source` | VARCHAR(30) | platform tag |
| `item_type` | VARCHAR(30) | whitelist |
| `status` | VARCHAR(20) | default `raw` |
| `submitted_by` | UUID → users | luôn = owner |
| `attachments` | JSONB | array of `{url, filename, mime, drive_file_id}` |
| `external_ref` | VARCHAR(200) | idempotency key per source |
| `created_at` | TIMESTAMPTZ | |

Unique: `(source, external_ref) WHERE external_ref IS NOT NULL`
