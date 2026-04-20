// n8n Code node — Zalo Personal Bot → Bình Vương Inbox
// Source spec: https://bot.zapps.me/docs/webhook/
//
// Upstream trigger: n8n Webhook Trigger (method POST) at path "/zalo-bot"
// Expected headers: X-Bot-Api-Secret-Token (set it equal to env ZALO_BOT_SECRET)
// Expected body shape:
//   { ok: true, result: { message: {...}, event_name: "message.text.received" | ... } }
//
// Output: single item ready for the HTTP Request node → POST /api/v1/inbox
//
// Drop this entire file into a Code node (Mode: Run Once for All Items).

const ZALO_BOT_SECRET = $env.ZALO_BOT_SECRET;

const req = $input.first().json;

// ------------------------------------------------------------
// 1. Verify shared secret — n8n stores incoming headers under .headers
// ------------------------------------------------------------
const incomingSecret =
  req.headers?.['x-bot-api-secret-token'] ??
  req.headers?.['X-Bot-Api-Secret-Token'];

if (ZALO_BOT_SECRET && incomingSecret !== ZALO_BOT_SECRET) {
  throw new Error(`Invalid X-Bot-Api-Secret-Token (got=${incomingSecret ?? 'none'})`);
}

// ------------------------------------------------------------
// 2. Drill into the Zalo payload; webhook wraps real data in .body
// ------------------------------------------------------------
const payload = req.body ?? req;
const result = payload.result ?? {};
const msg = result.message ?? {};
const eventName = result.event_name ?? '';

if (!msg.message_id) {
  throw new Error(`Zalo payload missing message_id; event=${eventName}`);
}

// ------------------------------------------------------------
// 3. Classify item_type + extract attachments
// ------------------------------------------------------------
let itemType = 'note';
let content = msg.text || msg.caption || '';
let url = '';
const attachmentURLs = [];

switch (eventName) {
  case 'message.text.received': {
    const t = (msg.text || '').trim();
    if (/^https?:\/\//i.test(t.split(/\s+/)[0] || '')) {
      itemType = 'link';
      url = t.split(/\s+/)[0];
    } else {
      itemType = 'note';
    }
    content = msg.text || '';
    break;
  }
  case 'message.image.received': {
    itemType = 'image';
    // `photo` may be a string URL or an array — handle both
    if (typeof msg.photo === 'string') attachmentURLs.push(msg.photo);
    else if (Array.isArray(msg.photo)) {
      for (const p of msg.photo) {
        const pu = typeof p === 'string' ? p : (p?.url || p?.file_url);
        if (pu) attachmentURLs.push(pu);
      }
    }
    content = msg.caption || '[Zalo image]';
    break;
  }
  case 'message.sticker.received': {
    itemType = 'image';
    const sticker = msg.sticker?.url || msg.sticker?.thumbnail_url || msg.sticker;
    if (typeof sticker === 'string') attachmentURLs.push(sticker);
    content = '[Zalo sticker]';
    break;
  }
  default: {
    // unsupported / file / voice / video fallbacks
    itemType = 'file';
    content = msg.text || msg.caption || `[Zalo ${eventName || 'unsupported'}]`;
  }
}

// Cap attachment count to match Bình Vương API limit (max 10)
const attachments = attachmentURLs.slice(0, 10);

// ------------------------------------------------------------
// 4. Build final payload for POST /api/v1/inbox
// ------------------------------------------------------------
const out = {
  content: content || `[Zalo ${eventName}]`,
  source: 'zalo',
  item_type: itemType,
  external_ref: `zalo:${msg.chat?.id ?? 'unknown'}:${msg.message_id}`,
};

if (url) out.url = url;
if (attachments.length) out.attachment_urls = attachments;

return [{ json: out }];
