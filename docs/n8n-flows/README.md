# n8n Flow Templates — Bình Vương OS Inbox

Import-ready n8n workflows for incoming webhooks → `POST /api/v1/inbox`.

**n8n instance:** <https://auto.binhvuong.vn>

## Import & Setup

### 1. Set env variable trong n8n

Vào **Settings → Variables** (hoặc edit `.env` của n8n container):

```
BVOS_API_KEY=905cd5982517fec4215c6f91450115e96a6e4093f87816a032db47ce6d62ebdd
```

> Thay bằng master key hiện tại. Rotate key → update biến này trong n8n → workflow tự dùng key mới.

### 2. Import workflow

1. n8n UI → **Workflows** → click **+** (top right) → **Import from File**
2. Chọn file JSON trong thư mục này
3. Save workflow

### 3. Credentials

Mỗi workflow cần credential tương ứng (Telegram Bot API, Zalo OA…). Xem chi tiết từng flow bên dưới.

### 4. Activate

Toggle **Active** ở góc trên phải. Workflow bắt đầu nhận sự kiện.

---

## Flows

### `telegram-to-inbox.json`

Telegram Bot → Bình Vương Inbox.

**Nodes:**
1. **Telegram Trigger** — nhận update từ bot (events: `message`)
2. **Transform Payload** (Code node) — classify `item_type` theo media kind + build `external_ref` dedup key
3. **POST to Inbox** — HTTP POST → `os.binhvuong.vn/api/v1/inbox` với `X-API-Key` header

**Setup:**
1. Tạo Bot qua [@BotFather](https://t.me/BotFather) → lấy `BOT_TOKEN`
2. n8n → **Credentials** → New → `Telegram API` → paste BOT_TOKEN → Save
3. Import `telegram-to-inbox.json`
4. Node **Telegram Trigger** → Credentials → chọn credential vừa tạo
5. Activate workflow → n8n tự set webhook với Telegram

**Test:** chat với bot → row mới xuất hiện trong `/inbox` trên os.binhvuong.vn.

**Giới hạn hiện tại:**
- Chỉ lấy `text` / `caption`, **chưa download file** từ Telegram (photo/voice/document sẽ vào inbox dưới dạng text placeholder `[image]` / `[voice]` / `[file]`)
- Muốn attach file thật: thêm node Telegram → `Get File` → HTTP Request với multipart `file=@<binary>`

**Idempotency:** `external_ref = tg:<chat_id>:<message_id>` → retry Telegram không tạo dup.

---

## Thêm flow mới

Khuôn mẫu chung:

```
[Platform Trigger] → [Code: transform] → [HTTP Request: POST /api/v1/inbox]
```

Required payload fields:
- `content` (string, required)
- `source` (string, e.g. `zalo`, `intercom`, `gmail`)
- `external_ref` (recommended, unique per source)

Optional:
- `url`, `item_type`, `attachment_urls[]`

See `docs/webhook-api.md` for full spec.
