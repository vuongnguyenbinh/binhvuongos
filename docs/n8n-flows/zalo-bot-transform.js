// n8n Code node — Zalo Personal Bot → Bình Vương Inbox (step 1 of 6)
// Source spec: https://bot.zapps.me/docs/webhook/
//
// Output shape — consumed by downstream IF → Download → POST (multipart | json):
//   {
//     content, source:"zalo", item_type, external_ref,
//     _has_file: boolean,     // IF node checks this to branch
//     _download_url: string,  // Zalo signed URL for binary (image/sticker)
//     _filename:     string   // filename hint for Drive
//   }
//
// Paste into an n8n Code node (Mode: Run Once for All Items).

const ZALO_BOT_SECRET = $env.ZALO_BOT_SECRET;

const req = $input.first().json;

// 1. Secret verify — header names are lowercased by n8n Webhook
const incomingSecret =
  req.headers?.['x-bot-api-secret-token'] ??
  req.headers?.['X-Bot-Api-Secret-Token'];

if (ZALO_BOT_SECRET && incomingSecret !== ZALO_BOT_SECRET) {
  throw new Error(`Invalid X-Bot-Api-Secret-Token (got=${incomingSecret ?? 'none'})`);
}

// 2. Drill payload
const payload = req.body ?? req;
const result = payload.result ?? {};
const msg = result.message ?? {};
const eventName = result.event_name ?? '';

if (!msg.message_id) {
  throw new Error(`Zalo payload missing message_id; event=${eventName}`);
}

const externalRef = `zalo:${msg.chat?.id ?? 'unknown'}:${msg.message_id}`;

// 3. Classify + extract download URL when media is present
let itemType = 'note';
let content = msg.text || msg.caption || '';
let downloadUrl = '';
let filename = `zalo-${msg.message_id}`;

switch (eventName) {
  case 'message.text.received': {
    const t = (msg.text || '').trim();
    const firstToken = t.split(/\s+/)[0] || '';
    itemType = /^https?:\/\//i.test(firstToken) ? 'link' : 'note';
    content = msg.text || '';
    break;
  }
  case 'message.image.received': {
    itemType = 'image';
    if (typeof msg.photo === 'string') {
      downloadUrl = msg.photo;
    } else if (Array.isArray(msg.photo) && msg.photo.length) {
      const first = msg.photo[0];
      downloadUrl = typeof first === 'string' ? first : (first?.url || first?.file_url || '');
    }
    content = msg.caption || '[Zalo image]';
    filename += '.jpg';
    break;
  }
  case 'message.sticker.received': {
    itemType = 'image';
    const sticker = msg.sticker;
    downloadUrl = typeof sticker === 'string'
      ? sticker
      : (sticker?.url || sticker?.thumbnail_url || '');
    content = '[Zalo sticker]';
    filename += '.webp';
    break;
  }
  default: {
    itemType = 'file';
    content = msg.text || msg.caption || `[Zalo ${eventName || 'unsupported'}]`;
  }
}

return [{
  json: {
    content: content || `[Zalo ${eventName}]`,
    source: 'zalo',
    item_type: itemType,
    external_ref: externalRef,
    _has_file: Boolean(downloadUrl),
    _download_url: downloadUrl,
    _filename: filename
  }
}];
