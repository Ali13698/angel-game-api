// lobby.js - ورود به بازی
async function joinGame(gameId) {
    if (!currentUser) return alert("ابتدا وارد شوید");
    
    document.getElementById('app').innerHTML += `
    <div class="game-room active" id="gameRoom">
        <div class="game-header">
            <button onclick="document.getElementById('gameRoom').remove()" style="background:none;border:none;color:white;">خروج</button>
            <div id="roomStatus">در حال ساخت اتاق...</div>
        </div>
        <div class="board-area"><div id="gameStateView">در انتظار شروع...</div></div>
        <div class="chat-area">
            <div class="chat-messages" id="chatMsgs"></div>
            <div class="chat-input-area">
                <input type="text" id="chatInput" placeholder="پیام..." onkeypress="if(event.key==='Enter') sendChat()">
                <button onclick="sendChat()">ارسال</button>
            </div>
        </div>
    </div>`;

    try {
        console.log("Requesting room for", gameId);
        const res = await fetch('/rooms', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ game_id: gameId, players: [currentUser.id.toString()] })
        });
        
        if (!res.ok) {
            const errText = await res.text();
            throw new Error(`Server responded ${res.status}: ${errText}`);
        }

        const data = await res.json();
        
        if (data.room_id) {
            document.getElementById('roomStatus').innerText = `اتاق: ${data.room_id}`;
            setTimeout(() => sendIntent("join", { room_id: data.room_id }), 500);
        }
    } catch (e) {
        console.error("Room creation error:", e);
        document.getElementById('roomStatus').innerText = `خطا: ${e.message}`;
    }
}
