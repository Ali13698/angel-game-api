const appEl = document.getElementById('app');
let currentUser = null;
let ws = null;
let currentRoom = null;

// لود کردن فایل proto برای ارتباط باینری
let AngelProto = null;
async function initProto() {
    AngelProto = await protobuf.load("/proto/angel.proto"); // این فایل باید در پوشه static کپی شود
}

function renderLayout(content, activeTab) {
    return `
    <div class="app-shell">
        <div class="top-header">
            <div class="logo-text">Angel</div>
            <div class="header-avatar" onclick="window.location.hash='#profile'">
                <img src="https://api.dicebear.com/7.x/bottts/svg?seed=${currentUser ? currentUser.id : 'Guest'}" style="width:36px;height:36px;border-radius:50%;border:2px solid var(--accent-blue);" />
            </div>
        </div>
        <div class="page-content">${content}</div>
        <div class="bottom-nav">
            <div class="nav-item ${activeTab === 'home' ? 'active' : ''}" onclick="window.location.hash='#home'"><span class="nav-icon">🏠</span></div>
            <div class="nav-item"><span class="nav-icon">👥</span></div>
            <div class="nav-center-btn" onclick="window.location.hash='#games'"><span class="nav-icon">🎮</span></div>
            <div class="nav-item"><span class="nav-icon">💬</span></div>
            <div class="nav-item ${activeTab === 'profile' ? 'active' : ''}" onclick="window.location.hash='#profile'"><span class="nav-icon">👤</span></div>
        </div>
    </div>`;
}

function renderHomePage() {
    return `
        <div class="section-title">بازی‌های آنلاین</div>
        <div class="games-grid">
            <div class="game-card" onclick="joinGame('ludo')">
                <div class="game-icon">🎲</div>
                <div class="game-name">منچ (Ludo)</div>
            </div>
            <div class="game-card">
                <div class="game-icon">🎱</div>
                <div class="game-name">بیلیارد (به‌زودی)</div>
            </div>
        </div>
    `;
}

function renderProfilePage() {
    return `
        <div style="text-align:center; margin-top: 20px;">
            <img src="https://api.dicebear.com/7.x/bottts/svg?seed=${currentUser ? currentUser.id : '1'}" style="width:100px;height:100px;border-radius:50%;border:4px solid var(--accent-blue);" />
            <h2 style="margin-top: 10px;">User ${currentUser ? currentUser.id : 'Unknown'}</h2>
            <p style="color:var(--text-muted);">سطح ۱</p>
            <div style="display:flex; justify-content:center; gap:20px; margin-top:20px;">
                <div style="background:var(--inner-card-bg); padding:15px; border-radius:10px; width:100px;">
                    <div style="font-size:20px; font-weight:bold;">0</div>
                    <div style="font-size:12px; color:var(--text-muted);">برد</div>
                </div>
                <div style="background:var(--inner-card-bg); padding:15px; border-radius:10px; width:100px;">
                    <div style="font-size:20px; font-weight:bold;">0</div>
                    <div style="font-size:12px; color:var(--text-muted);">باخت</div>
                </div>
            </div>
        </div>
    `;
}

function router() {
    const hash = window.location.hash || '#home';
    if (hash === '#home') appEl.innerHTML = renderLayout(renderHomePage(), 'home');
    else if (hash === '#profile') appEl.innerHTML = renderLayout(renderProfilePage(), 'profile');
    else if (hash === '#games') appEl.innerHTML = renderLayout(renderHomePage(), 'games');
}

// اتصال به هسته Go
async function connectWS() {
    if (ws) return;
    await initProto();
    
    const wsUrl = `wss://${location.host}/ws?user=${currentUser.id}`;
    ws = new WebSocket(wsUrl);
    ws.binaryType = 'arraybuffer';

    ws.onopen = () => {
        console.log('WS Connected');
        // ساخت پیام Join با Protobuf
        const ClientMsg = AngelProto.lookupType("angel.ClientMsg");
        const payload = ClientMsg.encode({ type: 'ping' }).finish();
        ws.send(payload);
    };

    ws.onmessage = (event) => {
        const ServerMsg = AngelProto.lookupType("angel.ServerMsg");
        const decoded = ServerMsg.decode(new Uint8Array(event.data));
        console.log('WS Received:', decoded);
        handleServerMessage(decoded);
    };

    ws.onclose = () => {
        console.log('WS Disconnected');
        ws = null;
    };
}

function handleServerMessage(msg) {
    if (msg.type === 'pong') {
        console.log('Pong received from Go Core!');
    }
}

// ورود به بازی
function joinGame(gameId) {
    if (!currentUser) return alert("ابتدا لاگین کنید");
    
    document.body.innerHTML += `
    <div class="game-room active" id="gameRoom">
        <div class="game-header">
            <button onclick="document.getElementById('gameRoom').remove()" style="background:none;border:none;color:white;">خروج</button>
            <div>منچ - در حال اتصال...</div>
            <div>⏱️ ۱۰s</div>
        </div>
        <div class="board-area">
            <div style="color:var(--text-muted);">صفحه بازی اینجا لود می‌شود (WebGL)</div>
        </div>
        <div class="chat-area">
            <div class="chat-messages" id="chatMsgs"></div>
            <div class="chat-input-area">
                <input type="text" id="chatInput" placeholder="پیام..." onkeypress="if(event.key==='Enter') sendChat()">
                <button onclick="sendChat()">ارسال</button>
            </div>
        </div>
    </div>`;
}

function sendChat() {
    const input = document.getElementById('chatInput');
    const text = input.value.trim();
    if (!text) return;
    
    // برای چت فعلاً از مسیر مستقل استفاده می‌کنیم
    const chatWs = new WebSocket(`wss://${location.host}/ws/chat?room_id=room1`);
    chatWs.onopen = () => {
        chatWs.send(JSON.stringify({ text: text, user: currentUser.id }));
        chatWs.close();
    };
    
    const msgs = document.getElementById('chatMsgs');
    msgs.innerHTML += `<div><b>من:</b> ${text}</div>`;
    input.value = '';
}

// شروع برنامه
async function init() {
    const urlParams = new URLSearchParams(location.search);
    let uid = urlParams.get('id');

    if (!uid && window.Telegram && window.Telegram.WebApp && window.Telegram.WebApp.initDataUnsafe) {
        uid = window.Telegram.WebApp.initDataUnsafe.user?.id;
    }
    
    if (!uid) uid = '1001'; // Guest ID for testing

    currentUser = { id: parseInt(uid) };
    
    await connectWS();
    
    if (!window.location.hash) window.location.hash = '#home';
    router();
}

window.addEventListener('hashchange', router);
window.addEventListener('load', init);
