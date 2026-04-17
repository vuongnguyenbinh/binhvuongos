// Chat bubble widget — sends messages to n8n webhook
(function() {
  const config = window.BVOS_CHAT || {};
  const webhookUrl = config.webhookUrl || '';
  const userName = config.userName || 'Khách';
  if (!webhookUrl) return; // No webhook configured, don't show

  let sessionId = localStorage.getItem('bvos_chat_session');
  if (!sessionId) {
    sessionId = 'sess_' + Math.random().toString(36).substr(2, 12);
    localStorage.setItem('bvos_chat_session', sessionId);
  }

  // Create bubble HTML
  const bubble = document.createElement('div');
  bubble.innerHTML = `
    <div id="chat-bubble-btn" style="position:fixed;bottom:24px;right:24px;z-index:9999;cursor:pointer;width:56px;height:56px;border-radius:50%;background:#1F3D2E;display:flex;align-items:center;justify-content:center;box-shadow:0 4px 12px rgba(0,0,0,.25);transition:transform .2s" onmouseover="this.style.transform='scale(1.1)'" onmouseout="this.style.transform='scale(1)'">
      <img src="https://binhvuong.vn/favicon.png" alt="Chat" style="width:32px;height:32px;border-radius:50%">
    </div>
    <div id="chat-panel" style="display:none;position:fixed;bottom:90px;right:24px;z-index:9999;width:340px;max-height:440px;background:var(--surface,#fff);border:1px solid var(--hairline,#e5e5e5);border-radius:12px;box-shadow:0 8px 30px rgba(0,0,0,.15);overflow:hidden;font-family:Inter,sans-serif">
      <div style="background:#1F3D2E;color:#fff;padding:12px 16px;display:flex;align-items:center;gap:8px">
        <img src="https://binhvuong.vn/favicon.png" style="width:24px;height:24px;border-radius:50%">
        <span style="font-weight:600;font-size:14px">Bình Vương OS</span>
        <span id="chat-close" style="margin-left:auto;cursor:pointer;opacity:.7;font-size:18px">&times;</span>
      </div>
      <div id="chat-messages" style="height:280px;overflow-y:auto;padding:12px;font-size:13px"></div>
      <div style="padding:8px 12px;border-top:1px solid var(--hairline,#e5e5e5);display:flex;gap:8px">
        <input id="chat-input" type="text" placeholder="Nhập tin nhắn..." style="flex:1;padding:8px 12px;border:1px solid #ddd;border-radius:8px;font-size:13px;outline:none">
        <button id="chat-send" style="padding:8px 16px;background:#1F3D2E;color:#fff;border:none;border-radius:8px;font-size:12px;cursor:pointer;font-weight:600">Gửi</button>
      </div>
    </div>`;
  document.body.appendChild(bubble);

  const btn = document.getElementById('chat-bubble-btn');
  const panel = document.getElementById('chat-panel');
  const close = document.getElementById('chat-close');
  const input = document.getElementById('chat-input');
  const send = document.getElementById('chat-send');
  const msgs = document.getElementById('chat-messages');

  btn.onclick = () => { panel.style.display = panel.style.display === 'none' ? 'block' : 'none'; input.focus(); };
  close.onclick = () => { panel.style.display = 'none'; };

  function addMsg(text, isUser) {
    const div = document.createElement('div');
    div.style.cssText = `margin-bottom:8px;padding:8px 12px;border-radius:8px;max-width:85%;font-size:13px;line-height:1.4;${isUser ? 'background:#1F3D2E;color:#fff;margin-left:auto' : 'background:#f3f3f3;color:#333'}`;
    div.textContent = text;
    msgs.appendChild(div);
    msgs.scrollTop = msgs.scrollHeight;
  }

  async function sendMessage() {
    const text = input.value.trim();
    if (!text) return;
    input.value = '';
    addMsg(text, true);
    try {
      const res = await fetch(webhookUrl, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_name: userName, message: text, session_id: sessionId, timestamp: new Date().toISOString() })
      });
      if (res.ok) {
        const data = await res.json().catch(() => null);
        if (data && data.output) addMsg(data.output, false);
        else addMsg('Đã nhận tin nhắn ✓', false);
      }
    } catch (e) { addMsg('Lỗi kết nối', false); }
  }

  send.onclick = sendMessage;
  input.onkeydown = (e) => { if (e.key === 'Enter') sendMessage(); };
})();
