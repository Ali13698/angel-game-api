function getUserId() {
    let t = new URLSearchParams(location.search).get("id");
    if (t && /^\d+$/.test(t)) return parseInt(t, 10);
    let e = window.Telegram?.WebApp?.initDataUnsafe?.user?.id;
    return e ? Number(e) : 1001;
}

function renderGameRoomUI(roomCode) {
    // اگر قبلا باز شده پاک کن
    document.getElementById("gameRoomOverlay")?.remove();

    const html = `
    <div id="gameRoomOverlay" style="position:fixed; top:0; left:0; width:100%; height:100%; background:#0c0f14; z-index:9999; display:flex; flex-direction:column;">
        
        <!-- هدر اتاق -->
        <div style="display:flex; justify-content:space-between; align-items:center; padding:15px; background:#15181e; border-bottom:1px solid #1c2028;">
            <button id="btnLeaveRoom" style="background:#ef4444; color:white; border:none; padding:8px 16px; border-radius:8px; font-family:inherit; cursor:pointer;">خروج</button>
            <div style="color:#8b5cf6; font-weight:bold; font-size:16px;">کد اتاق: ${roomCode}</div>
            <div style="width:40px;"></div> <!-- برای توازن -->
        </div>

        <!-- بازیکنان -->
        <div style="display:flex; justify-content:space-between; padding:20px;">
            <div style="text-align:center;">
                <img src="https://api.dicebear.com/7.x/bottts/svg?seed=sim1" style="width:60px; height:60px; border-radius:50%; border:3px solid #22c55e; background:#15181e;">
                <div style="color:#fff; margin-top:5px; font-size:14px;">شما</div>
            </div>
            <div style="color:#8b949e; font-size:24px; font-weight:bold; align-self:center;">VS</div>
            <div style="text-align:center;">
                <img src="https://api.dicebear.com/7.x/bottts/svg?seed=sim2" style="width:60px; height:60px; border-radius:50%; border:3px solid #ef4444; background:#15181e;">
                <div style="color:#fff; margin-top:5px; font-size:14px;">حریف</div>
            </div>
        </div>

        <!-- صفحه بازی (Placeholder) -->
        <div style="flex:1; display:flex; align-items:center; justify-content:center; padding:20px;">
            <div style="width:100%; max-width:400px; aspect-ratio:1/1; background:#1c2028; border-radius:16px; border:2px dashed #3b82f6; display:flex; flex-direction:column; align-items:center; justify-content:center;">
                <div style="font-size:50px; margin-bottom:10px;">🎲</div>
                <div style="color:#fff; font-weight:bold;" id="roomStatusText">در حال آماده‌سازی تخته...</div>
            </div>
        </div>

        <!-- دکمه پرتاب تاس (برای تست) -->
        <div style="padding:20px; text-align:center;">
            <button id="btnRollDice" style="background:#3b82f6; color:white; border:none; padding:15px 40px; border-radius:30px; font-size:18px; font-weight:bold; font-family:inherit; cursor:pointer; box-shadow:0 4px 15px rgba(59,130,246,0.4);">🎲 پرتاب تاس</button>
        </div>
    </div>`;

    document.body.insertAdjacentHTML('beforeend', html);

    document.getElementById('btnLeaveRoom').addEventListener('click', () => {
        document.getElementById("gameRoomOverlay").remove();
    });

    document.getElementById('btnRollDice').addEventListener('click', () => {
        if (window.sendIntent) {
            window.sendIntent("roll", {});
            document.getElementById("roomStatusText").innerText = "منتظر پاسخ سرور...";
        }
    });
}

async function joinGame(gameId) {
    let userId = getUserId();
    if (!userId) return alert("ابتدا وارد شوید");
    
    // نمایش لودینگ ساده
    let statusEl = document.getElementById("roomStatusText");
    
    try {
        // ۱. ساخت اتاق در Go
        let res = await fetch("/rooms", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({game_id: gameId, players: [userId.toString(), "bot"]})
        });
        
        if (!res.ok) throw new Error(await res.text());
        let data = await res.json();
        
        // ۲. نمایش UI اتاق
        renderGameRoomUI(data.room_id);
        
        // ۳. اتصال WebSocket
        window.initWS(userId, (msg) => {
            console.log("Server Msg:", msg);
            let st = document.getElementById("roomStatusText");
            if (msg.type === 'joined' && st) {
                st.innerText = "اتاق متصل شد، منتظر شروع...";
            }
            if (msg.type === 'event' && st) {
                // اگر رویداد تاس آمد
                if (msg.event && msg.event.type.includes('roll')) {
                    st.innerText = "تاس پرتاب شد!";
                }
            }
        });

        // ۴. ارسال پیام Join
        setTimeout(() => {
            if (window.sendIntent) window.sendIntent("join", { room_id: data.room_id });
        }, 500);

    } catch (err) {
        alert(`خطا در ساخت اتاق: ${err.message}`);
    }
}

window.joinGame = joinGame;
