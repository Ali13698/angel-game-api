const statusEl = document.getElementById('status-text');
const cells = document.querySelectorAll('.cell');
let ws = null;
let myTurn = false;

// دریافت آیدی کاربر (از URL یا تلگرام)
function getUserId() {
    const qs = new URLSearchParams(window.location.search).get("id");
    if (qs) return qs;
    const tg = window.Telegram?.WebApp?.initDataUnsafe?.user?.id;
    return tg ? String(tg) : "1001";
}

function init() {
    if (window.Telegram && Telegram.WebApp) { Telegram.WebApp.ready(); Telegram.WebApp.expand(); }
    connectWebSocket();
}

function connectWebSocket() {
    const userId = getUserId();
    // برای تست فعلی به سرور اصلی وصل میشویم
    ws = new WebSocket(`wss://game.nomi98.com/ws?user=${userId}`);
    ws.binaryType = 'arraybuffer';

    ws.onopen = () => {
        statusEl.textContent = "متصل شد؛ در حال ساخت اتاق...";
        // ارسال درخواست ساخت اتاق برای بازی دوز
        fetch('/rooms', {
            method: 'POST', headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ game_id: 'tic-tac-toe', players: [userId, "bot"] })
        }).then(res => res.json()).then(data => {
            if(data.room_id) {
                // اتصال به اتاق ساخته شده
                ws.send(JSON.stringify({ type: 'join', room_id: data.room_id }));
                statusEl.textContent = "در انتظار حرکت شما...";
                myTurn = true; // فرض میکنیم بازیکن اول هستیم
            }
        });
    };

    ws.onmessage = (event) => {
        try {
            const msg = JSON.parse(new TextDecoder().decode(event.data));
            
            if (msg.events) {
                for (const ev of msg.events) {
                    if (ev.type === 'move_made') {
                        const { seat, cell } = ev.payload;
                        updateCell(cell, seat === 0 ? 'X' : 'O');
                        myTurn = (seat !== 0); // اگر ربات حرکت کرد، نوبت ما میشود
                        statusEl.textContent = myTurn ? "نوبت شماست" : "نوبت ربات...";
                    } else if (ev.type === 'game_over') {
                        const { winner } = ev.payload;
                        if (winner === 0) statusEl.textContent = "🎉 شما بردید!";
                        else if (winner === 1) statusEl.textContent = "🤖 ربات برد!";
                        else statusEl.textContent = "🤝 مساوی شد!";
                        myTurn = false;
                    }
                }
            }
        } catch (e) {
            console.error("Decode error", e);
        }
    };

    ws.onclose = () => { statusEl.textContent = "ارتباط قطع شد."; myTurn = false; };
}

function updateCell(index, symbol) {
    const cell = cells[index];
    if (cell) {
        cell.textContent = symbol;
        cell.classList.add(symbol.toLowerCase());
    }
}

cells.forEach(cell => {
    cell.addEventListener('click', () => {
        if (!myTurn || cell.textContent !== '') return;
        
        const index = parseInt(cell.dataset.index);
        myTurn = false;
        statusEl.textContent = "نوبت ربات...";
        
        // ارسال حرکت به سرور
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({
                type: 'move',
                payload: { cell: index }
            }));
        }
    });
});

init();
