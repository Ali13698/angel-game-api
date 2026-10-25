(function(global){
'use strict';
const NET = global.AngelNet;
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

// رسم تخته
for(let r=1;r<=GRID_COUNT;r++) {
  for(let c=1;c<=GRID_COUNT;c++){
    svg.append("rect")
       .attr("x",getX(c))
       .attr("y",getY(r))
       .attr("width",CELL_SIZE)
       .attr("height",CELL_SIZE)
       .attr("fill",PALETTE.white)
       .attr("stroke",PALETTE.border)
       .attr("stroke-width",1.2);
  }
}

function drawBase(sr,sc,mc,dc){
  const bs=6*CELL_SIZE,bx=getX(sc),by=getY(sr);
  svg.append("rect").attr("x",bx).attr("y",by).attr("width",bs).attr("height",bs).attr("fill",mc).attr("stroke",PALETTE.border).attr("stroke-width",1.2);
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

function paintCell(r,c,col){
  svg.append("rect").attr("x",getX(c)).attr("y",getY(r)).attr("width",CELL_SIZE).attr("height",CELL_SIZE).attr("fill",col).attr("stroke",PALETTE.border).attr("stroke-width",1.2);
}
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

const BOARD_PATH = [
  [7,2],[7,3],[7,4],[7,5],[7,6],
  [6,7],[5,7],[4,7],[3,7],[2,7],[1,7],
  [1,8],[1,9],
  [2,9],[3,9],[4,9],[5,9],[6,9],
  [7,10],[7,11],[7,12],[7,13],[7,14],[7,15],
  [8,15],[9,15],
  [9,14],[9,13],[9,12],[9,11],[9,10],
  [10,9],[11,9],[12,9],[13,9],[14,9],[15,9],
  [15,8],[15,7],
  [14,7],[13,7],[12,7],[11,7],[10,7],
  [9,6],[9,5],[9,4],[9,3],[9,2],[9,1],
  [8,1],[7,1]
];

const COLORS = [
  {name:'شما',key:'red',hex:'#e52521',offset:0, homeLane:[[8,2],[8,3],[8,4],[8,5],[8,6]], base:[[3,3],[3,4],[4,3],[4,4]]},
  {name:'حریف',key:'green',hex:'#1eae3a',offset:13, homeLane:[[2,8],[3,8],[4,8],[5,8],[6,8]], base:[[3,12],[3,13],[4,12],[4,13]]},
  {name:'بازیکن ۳',key:'yellow',hex:'#fdc010',offset:26, homeLane:[[8,14],[8,13],[8,12],[8,11],[8,10]], base:[[12,12],[12,13],[13,12],[13,13]]},
  {name:'بازیکن ۴',key:'blue',hex:'#007cee',offset:39, homeLane:[[14,8],[13,8],[12,8],[11,8],[10,8]], base:[[12,3],[12,4],[13,3],[13,4]]}
];

const tokenLayer = document.getElementById('token-layer');
let boardState = [-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1];
let myTurn = false;

function cellCenter(r,c){
  return [getX(c)+CELL_SIZE/2, getY(r)+CELL_SIZE/2];
}

function cellPosForStep(ci, ti, step){
  const color = COLORS[ci];
  if(step === -1){
    const [r,c] = color.base[ti];
    return cellCenter(r,c);
  }
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
    for(let ti=0; ti<4; ti++){
      let pos = boardState[ci*4 + ti];
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
        
        // ✅ کلیک روی مهره برای حرکت
        el.onclick = () => {
          if(ci === 0 && myTurn){
            NET.send("move", { token: ti });
          }
        };
        tokenLayer.appendChild(el);
      }
      
      const posArr = cellPosForStep(ci, ti, pos);
      el.style.left = ((posArr[0] / 600) * 100) + '%';
      el.style.top  = ((posArr[1] / 600) * 100) + '%';
    }
  });
}
renderTokens();

// ✅ اتصال به وب‌سوکت
window.__rollDice = function(){
  if(!myTurn) return document.getElementById('turnText').textContent = "نوبت شما نیست!";
  NET.send("roll", {});
  document.getElementById('turnText').textContent = "در حال پرتاب تاس...";
};

NET.on('joined', () => {
  document.getElementById('turnText').textContent = "متصل شد. منتظر شروع...";
});

NET.on('game.started', () => {
  myTurn = true;
  document.getElementById('turnText').textContent = "بازی شروع شد! نوبت شماست، تاس بیندازید.";
});

NET.on('game.rolled', (p) => {
  document.getElementById('turnText').textContent = "تاس: " + p.dice + " (مهره را انتخاب کنید)";
});

NET.on('game.moved', (p) => {
  boardState[p.token] = p.pos;
  renderTokens();
});

NET.on('turn_changed', (p) => {
  myTurn = (p.turn === 0);
  document.getElementById('turnText').textContent = myTurn ? "نوبت شماست! تاس بیندازید." : "نوبت حریف...";
});

async function boot() {
  const user = new URLSearchParams(location.search).get('user') || '1001';
  const room = new URLSearchParams(location.search).get('room');
  await NET.connect(user, { roomId: room });
  NET.join(room);
}
boot();

})(window);
