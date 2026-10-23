import { logInfo, logError } from '../../profile/src/logger';

function getUserId(): number | null {
  const fromQuery = new URLSearchParams(location.search).get('id');
  if (fromQuery && /^\d+$/.test(fromQuery)) return parseInt(fromQuery, 10);
  const tg = (window as any).Telegram?.WebApp;
  const tgId = tg?.initDataUnsafe?.user?.id;
  if (tgId) return Number(tgId);
  return 1001;
}

function el(id: string): HTMLElement | null {
  return document.getElementById(id);
}

function renderRoomShell(): void {
  const app = el('app');
  if (!app || el('gameRoom')) return;
  
  // 🛠️ تست بیلد: این نوار سبز باید در صفحه ظاهر شود
  app.insertAdjacentHTML('beforeend', `
    <div class="game-room active" id="gameRoom">
      <div style="background: #22c55e; color: white; padding: 10px; text-align: center; font-weight: bold; font-size: 14px; z-index: 9999;">
        ✅ اتصال به هسته Go برقرار شد! (نسخه بیلد خودکار: v1.0.5)
      </div>
      <div class="game-header">
        <button id="closeRoomBtn" style="background:none;border:none;color:white;">خروج</button>
        <div id="roomStatus">در حال ساخت اتاق...</div>
      </div>
      <div class="board-area"><div id="gameStateView">در انتظار شروع...</div></div>
      <div class="chat-area">
        <div class="chat-messages" id="chatMsgs"></div>
        <div class="chat-input-area">
          <input type="text" id="chatInput" placeholder="پیام..." />
          <button id="sendChatBtn">ارسال</button>
        </div>
      </div>
    </div>`);
    
  el('closeRoomBtn')?.addEventListener('click', () => el('gameRoom')?.remove());
  el('sendChatBtn')?.addEventListener('click', sendChat);
  el('chatInput')?.addEventListener('keypress', (e: KeyboardEvent) => {
    if (e.key === 'Enter') sendChat();
  });
}

function sendChat(): void {
  const input = el('chatInput') as HTMLInputElement | null;
  const text = input?.value.trim();
  if (!text) return;
  const send = (window as any).sendIntent;
  if (typeof send === 'function') {
    send('chat', { text });
    input.value = '';
  } else {
    logError('sendIntent در دسترس نیست');
  }
}

async function joinGame(gameId: string): Promise<void> {
  const userId = getUserId();
  if (!userId) { alert('ابتدا وارد شوید'); return; }
  
  renderRoomShell();
  const status = el('roomStatus');
  
  try {
    const res = await fetch('/rooms', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ game_id: gameId, players: [userId.toString()] })
    });
    if (!res.ok) throw new Error(await res.text());
    const data = await res.json();
    if (status) status.innerText = `اتاق: ${data.room_id}`;
    
    setTimeout(() => {
      const send = (window as any).sendIntent;
      if (typeof send === 'function') send('join', { room_id: data.room_id });
    }, 500);
  } catch (e) {
    if (status) status.innerText = `خطا: ${(e as Error).message}`;
  }
}

(window as any).joinGame = joinGame;
(window as any).sendChat = sendChat;
