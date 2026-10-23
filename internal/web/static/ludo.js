(()=>{
'use strict';
const GRID_COUNT = 15;
const CELL_SIZE = 600 / GRID_COUNT;
const PALETTE = {
  red:"#e52521", redDark:"#c41c18",
  green:"#1eae3a", greenDark:"#178e2e",
  blue:"#007cee", blueDark:"#0062c4",
  yellow:"#fdc010", yellowDark:"#dda403",
  border:"#dbe2ea", white:"#ffffff"
};
const svg = d3.select("#ludo");
const getX = c => (c-1)*CELL_SIZE, getY = r => (r-1)*CELL_SIZE;

// ==================== ساخت تخته ====================
for(let r=1;r<=GRID_COUNT;r++) for(let c=1;c<=GRID_COUNT;c++){
  svg.append("rect").attr("x",getX(c)).attr("y",getY(r)).attr("width",CELL_SIZE).attr("height",CELL_SIZE)
    .attr("fill",PALETTE.white).attr("stroke",PALETTE.border).attr("stroke-width",1.2);
}
function drawBase(sr,sc,mc,dc){
  const bs=6*CELL_SIZE,bx=getX(sc),by=getY(sr);
  svg.append("rect").attr("x",bx).attr("y",by).attr("width",bs).attr("height",bs).attr("fill",mc)
    .attr("stroke",PALETTE.border).attr("stroke-width",1.2);
  const m=CELL_SIZE*0.9,ins=bs-2*m,ix=bx+m,iy=by+m;
  svg.append("rect").attr("x",ix).attr("y",iy).attr("width",ins).attr("height",ins).attr("rx",10).attr("fill",dc);
  const o=ins*0.28;
  [[ix+o,iy+o],[ix+ins-o,iy+o],[ix+o,iy+ins-o],[ix+ins-o,iy+ins-o]].forEach(([cx,cy])=>{
    svg.append("circle").attr("cx",cx).attr("cy",cy).attr("r",CELL_SIZE*0.38).attr("fill",PALETTE.white);
    svg.append("circle").attr("cx",cx).attr("cy",cy).attr("r",CELL_SIZE*0.22).attr("fill",dc);
  });
}
drawBase(1,1,PALETTE.red,PALETTE.redDark);
drawBase(1,10,PALETTE.green,PALETTE.greenDark);
drawBase(10,1,PALETTE.blue,PALETTE.blueDark);
drawBase(10,10,PALETTE.yellow,PALETTE.yellowDark);
function paintCell(r,c,col){svg.append("rect").attr("x",getX(c)).attr("y",getY(r)).attr("width",CELL_SIZE).attr("height",CELL_SIZE).attr("fill",col).attr("stroke",PALETTE.border).attr("stroke-width",1.2);}
for(let c=2;c<=6;c++)paintCell(8,c,PALETTE.red);paintCell(7,2,PALETTE.red);
for(let r=2;r<=6;r++)paintCell(r,8,PALETTE.green);paintCell(2,9,PALETTE.green);
for(let c=10;c<=14;c++)paintCell(8,c,PALETTE.yellow);paintCell(9,14,PALETTE.yellow);
for(let r=10;r<=14;r++)paintCell(r,8,PALETTE.blue);paintCell(14,7,PALETTE.blue);
const ccx=getX(7),ccy=getY(7),cs=3*CELL_SIZE,mx=ccx+cs/2,my=ccy+cs/2;
[
  {p:`${ccx},${ccy} ${ccx+cs},${ccy} ${mx},${my}`,c:PALETTE.green},
  {p:`${ccx+cs},${ccy} ${ccx+cs},${ccy+cs} ${mx},${my}`,c:PALETTE.yellow},
  {p:`${ccx},${ccy+cs} ${ccx+cs},${ccy+cs} ${mx},${my}`,c:PALETTE.blue},
  {p:`${ccx},${ccy} ${ccx},${ccy+cs} ${mx},${my}`,c:PALETTE.red}
].forEach(t=>svg.append("polygon").attr("points",t.p).attr("fill",t.c));
svg.append("rect").attr("x",ccx).attr("y",ccy).attr("width",cs).attr("height",cs).attr("fill","none").attr("stroke",PALETTE.border).attr("stroke-width",1.2);

const COLORS = [
  {name:'بازیکن ۱',key:'red',hex:'#e52521',offset:0, homeLane:[[8,2],[8,3],[8,4],[8,5],[8,6]], base:[[3,3],[3,4],[4,3],[4,4]]},
  {name:'بازیکن ۲',key:'green',hex:'#1eae3a',offset:13, homeLane:[[2,8],[3,8],[4,8],[5,8],[6,8]], base:[[3,12],[3,13],[4,12],[4,13]]},
  {name:'بازیکن ۳',key:'yellow',hex:'#fdc010',offset:26, homeLane:[[8,14],[8,13],[8,12],[8,11],[8,10]], base:[[12,12],[12,13],[13,12],[13,13]]},
  {name:'بازیکن ۴',key:'blue',hex:'#007cee',offset:39, homeLane:[[14,8],[13,8],[12,8],[11,8],[10,8]], base:[[12,3],[12,4],[13,3],[13,4]]}
];
const CENTER = [8,8];
const MAX_POS = 56;

// ==================== وضعیت محلی (فقط برای انیمیشن) ====================
let state = {red:[-1,-1,-1,-1],green:[-1,-1,-1,-1],yellow:[-1,-1,-1,-1],blue:[-1,-1,-1,-1]};
let current = 0;
let isMyTurn = false;
let myUserId = new URLSearchParams(location.search).get("id") || "1001";
const tokenLayer = document.getElementById('token-layer');
let rotX=-26, rotY=38;
const cube = document.getElementById('cube');
const diceArea = document.getElementById('diceArea');

function cellCenter(r,c){ return [getX(c)+CELL_SIZE/2, getY(r)+CELL_SIZE/2]; }
function cellPosForStep(ci, ti, step){
  const color = COLORS[ci];
  if(step === -1){ const [r,c] = color.base[ti]; return cellCenter(r,c); }
  if(step >= 0 && step <= 50){
    const idx = (color.offset + step) % 52;
    const [r,c] = BOARD_PATH[idx] || [8,8];
    return cellCenter(r,c);
  }
  if(step >= 51 && step <= 55){
    const [r,c] = color.homeLane[step - 51] || [8,8];
    return cellCenter(r,c);
  }
  return cellCenter(8,8);
}

function renderTokens(){
  COLORS.forEach((color,ci)=>{
    state[color.key].forEach((pos,ti)=>{
      const id = 'tk-' + ci + '-' + ti;
      let el = document.getElementById(id);
      if(!el){
        el = document.createElement('div');
        el.className = 'token';
        el.id = id;
        const body = document.createElement('div');
        body.className = 'token-body';
        body.style.background = color.hex;
        el.appendChild(body);
        el.addEventListener('click', () => {
          if(ci === current && isMyTurn){
             window.sendIntent("move", { token: ti });
          }
        });
        tokenLayer.appendChild(el);
      }
      const posArr = cellPosForStep(ci, ti, pos);
      el.style.left = ((posArr[0] / 600) * 100) + '%';
      el.style.top  = ((posArr[1] / 600) * 100) + '%';
    });
  });
}

function animateDice(val) {
    const faceRot = {1:{x:0,y:0},2:{x:0,y:180},3:{x:0,y:-90},4:{x:0,y:90},5:{x:-90,y:0},6:{x:90,y:0}};
    const target=faceRot[val];
    const fromX=rotX, fromY=rotY;
    cube.style.transition='none';
    cube.style.transform=`rotateX(${fromX}deg) rotateY(${fromY}deg)`;
    void cube.offsetWidth;
    const finalX=Math.round((fromX-target.x)/360)*360+target.x+1080;
    const finalY=Math.round((fromY-target.y)/360)*360+target.y+1080;
    rotX=finalX; rotY=finalY;
    requestAnimationFrame(()=>{
      cube.style.transition='transform 1.7s cubic-bezier(.15,.85,.2,1)';
      cube.style.transform=`rotateX(${finalX}deg) rotateY(${finalY}deg)`;
    });
}

window.__rollDice = () => {
    if(!isMyTurn) return document.getElementById('turnText').innerText = "نوبت شما نیست!";
    window.sendIntent("roll", {});
    document.getElementById('turnText').innerText = "در حال ارسال به سرور...";
};

// ==================== اتصال به سرور ====================
function connectToServer(roomId) {
    window.initWS(myUserId, (msg) => {
        console.log("Server Event:", msg);
        let payload = {};
        if (msg.payload) {
             try { payload = JSON.parse(new TextDecoder().decode(msg.payload)); } catch(e){}
        }

        if (msg.type === 'joined') {
            document.getElementById('turnText').innerText = "اتاق متصل شد.";
        }
        
        if (msg.type === 'event') {
            const evType = msg.event.type;
            let evPayload = {};
            try { evPayload = JSON.parse(new TextDecoder().decode(msg.event.payload)); } catch(e){}
            
            if (evType === 'turn_changed') {
                current = evPayload.seat;
                isMyTurn = (String(current) === String(myUserId));
                document.getElementById('turnText').innerText = isMyTurn ? "نوبت شماست!" : "نوبت حریف...";
                document.querySelectorAll('.player').forEach(e=>e.classList.remove('active'));
                document.getElementById('p'+current)?.classList.add('active');
            }
            
            if (evType === 'ludo.rolled') {
                animateDice(evPayload.die);
            }

            if (evType === 'ludo.moved') {
                const ci = evPayload.seat;
                const ti = evPayload.token;
                state[COLORS[ci].key][ti] = evPayload.position;
                renderTokens();
            }
        }
    });
    setTimeout(() => window.sendIntent("join", { room_id: roomId }), 500);
}

// راه‌اندازی
const roomId = new URLSearchParams(location.search).get("room");
if(roomId) connectToServer(roomId);
else document.getElementById('turnText').innerText = "خطا: کد اتاق یافت نشد";

})();
