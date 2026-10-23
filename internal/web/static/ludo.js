(()=>{
'use strict';
const GRID_COUNT = 15;
const CELL_SIZE = 600 / GRID_COUNT;
const PALETTE = { red:"#e52521", green:"#1eae3a", blue:"#007cee", yellow:"#fdc010", border:"#dbe2ea", white:"#ffffff" };
const getX = c => (c-1)*CELL_SIZE, getY = r => (r-1)*CELL_SIZE;
const svg = d3.select("#ludo");

// ساخت تخته (همان کدهای شما - بدون تغییر)
for(let r=1;r<=GRID_COUNT;r++) for(let c=1;c<=GRID_COUNT;c++){
  svg.append("rect").attr("x",getX(c)).attr("y",getY(r)).attr("width",CELL_SIZE).attr("height",CELL_SIZE).attr("fill",PALETTE.white).attr("stroke",PALETTE.border).attr("stroke-width",1.2);
}
// ... (بقیه کدهای ساخت تخته، خانه‌ها و ستاره‌ها را دقیقاً از فایل خودتان اینجا کپی کنید) ...

// State محلی (فقط برای انیمیشن)
let localState = {red:[-1,-1,-1,-1],green:[-1,-1,-1,-1],yellow:[-1,-1,-1,-1],blue:[-1,-1,-1,-1]};
let current = 0;
let isMyTurn = false;
let myUserId = new URLSearchParams(location.search).get("id") || "1001";

// اتصال به سرور Go
function connectToServer(roomId) {
    window.initWS(myUserId, (msg) => {
        console.log("Server Event:", msg);
        
        if (msg.type === 'joined') {
            document.getElementById('turnText').innerText = "اتاق متصل شد. منتظر شروع...";
        }
        
        if (msg.type === 'event' && msg.event) {
            // دریافت رویدادهای بازی از Go
            const evType = msg.event.type;
            let payload = {};
            try { payload = JSON.parse(new TextDecoder().decode(msg.event.payload)); } catch(e){}
            
            if (evType === 'turn_changed') {
                current = payload.seat;
                isMyTurn = (String(current) === String(myUserId));
                document.getElementById('turnText').innerText = isMyTurn ? "نوبت شماست!" : "نوبت حریف...";
                document.querySelectorAll('.player').forEach(e=>e.classList.remove('active'));
                document.getElementById('p'+current)?.classList.add('active');
            }
            
            if (evType === 'ludo.rolled') {
                animateDice(payload.die);
            }

            if (evType === 'ludo.moved') {
                // انیمیشن حرکت مهره (به جای منطق محلی)
                localState[COLORS[current].key][payload.token] = payload.position;
                renderTokens();
            }
        }
    });
    // ارسال پیام Join به سرور
    setTimeout(() => window.sendIntent("join", { room_id: roomId }), 500);
}

// دکمه تاس زدن (ارسال به سرور)
function rollDice(){
    if(!isMyTurn) return alert("نوبت شما نیست!");
    // ارسال قصد به سرور Go
    window.sendIntent("roll", {});
    document.getElementById('turnText').innerText = "در حال ارسال به سرور...";
}
window.__rollDice = rollDice;

// انیمیشن تاس (فقط بصری)
function animateDice(val) {
    const cube = document.getElementById('cube');
    const diceArea = document.getElementById('diceArea');
    // ... کدهای انیمیشن چرخش تاس شما (بدون تغییر در منطق) ...
    // در پایان انیمیشن، عدد val را نمایش بده
}

function renderTokens() {
    // همان کدهای رندر شما، اما فقط بر اساس localState که از سرور آپدیت می‌شود
}

// راه‌اندازی اولیه
const roomId = new URLSearchParams(location.search).get("room");
if(roomId) connectToServer(roomId);
else document.getElementById('turnText').innerText = "خطا: کد اتاق یافت نشد";

})();
