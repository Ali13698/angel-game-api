// lobby.js - مدیریت ورود به بازی و صفحه بازی
function joinGame(gameId) {
    if (!currentUser) return alert("ابتدا وارد شوید");
    
    document.getElementById('app').innerHTML += `
    <div class="game-room active" id="gameRoom">
        <div class="game-header">
            <button onclick="document.getElementById('gameRoom').remove()" style="background:none;border:none;color:white;">خروج</button>
            <div>در حال اتصال به سرور...</div>
        </div>
        <div class="board-area"><div id="gameStateView" style="color:var(--text-muted);">در انتظار شروع...</div></div>
        <div class="chat-area">
            <div class="chat-messages" id="chatMsgs"></div>
            <div class="chat-input-area">
                <input type="text" id="chatInput" placeholder="پیام..." onkeypress="if(event.key==='Enter') sendChat()">
                <button onclick="sendChat()">ارسال</button>
            </div>
        </div>
    </div>`;

    // ایجاد اتاق در سرور Go
    fetch('/rooms', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ game_id: gameId, players: [currentUser.id.toString()] })
    }).then(r => r.json()).then(data => {
        if (data.room_id) {
            initWS(currentUser.id, handleGameMsg);
            setTimeout(() => sendIntent("join", { room_id: data.room_id }), 500);
        }
    });
}

function handleGameMsg(msg) {
    console.log("Received from server:", msg);
    if (msg.type === 'joined') {
        document.querySelector('.game-header div:nth-child(2)').innerText = `اتاق: ${msg.room_id}`;
        document.getElementById('gameStateView').innerText = "بازی شروع شد!";
    }
    if (msg.type === 'event') {
        // آپدیت وضعیت بازی (WASM)
        const payload = new TextDecoder().decode(msg.event.payload);
        console.log("Game Event:", msg.event.type, payload);
    }
}

function sendChat() {
    const input = document.getElementById('chatInput');
    const text = input.value.trim();
    if (!text) return;
    
    // ارسال از طریق WebSocket
    sendIntent("chat", { text: text });
    
    const msgs = document.getElementById('chatMsgs');
    msgs.innerHTML += `<div><b>من:</b> ${text}</div>`;
    input.value = '';
}
