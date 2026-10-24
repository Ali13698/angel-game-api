(()=>{
'use strict';
const getX = c => (c-1)*40, getY = r => (r-1)*40;
const svg = d3.select("#ludo");
const BOARD_PATH = [[7,2],[7,3],[7,4],[7,5],[7,6],[6,7],[5,7],[4,7],[3,7],[2,7],[1,7],[1,8],[1,9],[2,9],[3,9],[4,9],[5,9],[6,9],[7,10],[7,11],[7,12],[7,13],[7,14],[7,15],[8,15],[9,15],[9,14],[9,13],[9,12],[9,11],[9,10],[10,9],[11,9],[12,9],[13,9],[14,9],[15,9],[15,8],[15,7],[14,7],[13,7],[12,7],[11,7],[10,7],[9,6],[9,5],[9,4],[9,3],[9,2],[9,1],[8,1],[7,1]];
const COLORS = [{hex:'#e52521',offset:0},{hex:'#1eae3a',offset:13},{hex:'#fdc010',offset:26},{hex:'#007cee',offset:39}];

// رسم تخته (فقط ظاهر)
for(let r=1;r<=15;r++) for(let c=1;c<=15;c++) svg.append("rect").attr("x",getX(c)).attr("y",getY(r)).attr("width",40).attr("height",40).attr("fill","#fff").attr("stroke","#dbe2ea");
// (کدهای رسم خانه‌های رنگی شما در اینجا قرار می‌گیرد تا حجم کم شود)

let myUserId = new URLSearchParams(location.search).get("id") || "1001";
let localState = null;

function renderTokens() {
    if(!localState) return;
    document.getElementById('token-layer').innerHTML = '';
    localState.tokens.forEach((playerTokens, pIdx) => {
        playerTokens.forEach((pos, tIdx) => {
            let x = 0, y = 0;
            if(pos === -1) { x = 60 + (tIdx%2)*40; y = 60 + Math.floor(tIdx/2)*40; } // مختصات فرضی پایه
            else { let idx = (COLORS[pIdx].offset + pos) % 52; let cell = BOARD_PATH[idx]; x = getX(cell[1])+20; y = getY(cell[0])+20; }
            
            let el = document.createElement('div');
            el.className = 'token';
            el.style.left = (x/6)+'%'; el.style.top = (y/6)+'%';
            el.style.background = COLORS[pIdx].hex;
            el.onclick = () => {
                if(localState.turn === pIdx && myUserId === localState.seats[pIdx]) {
                    window.sendIntent("move", {token: tIdx});
                }
            };
            document.getElementById('token-layer').appendChild(el);
        });
    });
}

window.__rollDice = () => {
    if(localState && localState.seats[localState.turn] === myUserId) window.sendIntent("roll", {});
    else alert("نوبت شما نیست!");
};

document.getElementById('startGameBtn').addEventListener('click', async () => {
    document.getElementById('settingsScreen').style.display = 'none';
    let playersArr = [myUserId, "bot_1"]; // ارسال به سرور برای ساخت ربات
    
    let res = await fetch("/rooms", {method: "POST", headers: {"Content-Type": "application/json"}, body: JSON.stringify({game_id: "ludo", players: playersArr})});
    let data = await res.json();
    
    window.initWS(myUserId, (msg) => {
        if(msg.type === 'joined') {
            document.getElementById('turnText').innerText = "بازی شروع شد!";
            // دریافت وضعیت اولیه
            window.sendIntent("get_state", {room_id: data.room_id});
        }
        if(msg.type === 'state') {
            localState = JSON.parse(new TextDecoder().decode(msg.state));
            renderTokens();
            let isMe = localState.seats[localState.turn] === myUserId;
            document.getElementById('turnText').innerText = isMe ? "نوبت شماست! تاس بریزید." : "نوبت حریف...";
        }
        if(msg.type === 'event' && msg.event.type === 'ludo.rolled') {
            let p = JSON.parse(new TextDecoder().decode(msg.event.payload));
            // انیمیشن تاس با عدد p.die
        }
    });
    setTimeout(() => window.sendIntent("join", {room_id: data.room_id}), 500);
});
})();
