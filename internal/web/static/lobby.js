function getUserId() {
    let t = new URLSearchParams(location.search).get("id");
    if (t && /^\d+$/.test(t)) return parseInt(t, 10);
    let e = window.Telegram?.WebApp?.initDataUnsafe?.user?.id;
    return e ? Number(e) : 1001;
}

function renderRoom() {
    let t = document.getElementById("app");
    if (!t || document.getElementById("gameRoom")) return;
    
    t.insertAdjacentHTML("beforeend", `
    <div class="game-room active" id="gameRoom" style="padding:20px; color:white;">
      <div style="background: #22c55e; color: white; padding: 10px; text-align: center; font-weight: bold; border-radius:8px; margin-bottom:15px;">
        ✅ اتصال به هسته Go برقرار شد!
      </div>
      <div class="game-header" style="display:flex; justify-content:space-between; margin-bottom:15px;">
        <button id="closeRoomBtn" style="background:#ef4444; border:none; color:white; padding:5px 10px; border-radius:5px;">خروج</button>
        <div id="roomStatus">در حال ساخت اتاق...</div>
      </div>
      <div class="board-area" style="background:#1c2028; height:200px; display:flex; align-items:center; justify-content:center; border-radius:8px;">
        <div id="gameStateView">در انتظار شروع...</div>
      </div>
    </div>`);
    
    document.getElementById("closeRoomBtn")?.addEventListener("click", () => {
        document.getElementById("gameRoom")?.remove();
    });
}

async function joinGame(gameId) {
    let userId = getUserId();
    if (!userId) return alert("ابتدا وارد شوید");
    
    renderRoom();
    let statusEl = document.getElementById("roomStatus");
    
    try {
        let res = await fetch("/rooms", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({game_id: gameId, players: [userId.toString()]})
        });
        
        if (!res.ok) throw new Error(await res.text());
        let data = await res.json();
        
        if (statusEl) statusEl.innerText = `اتاق: ${data.room_id}`;
        
        window.initWS(userId, (msg) => {
            console.log("Server Msg:", msg);
            if (msg.type === 'joined') {
                document.getElementById("gameStateView").innerText = "بازی شروع شد!";
            }
        });

        setTimeout(() => {
            if (window.sendIntent) window.sendIntent("join", { room_id: data.room_id });
        }, 500);

    } catch (err) {
        if (statusEl) statusEl.innerText = `خطا: ${err.message}`;
    }
}

window.joinGame = joinGame;
