(function(global){
  'use strict';
  const NET = global.AngelNet;
  const svg = d3.select("#ludo");
  const CELL = 600/15;
  const getX = c => (c-1)*CELL, getY = r => (r-1)*CELL;

  // رسم تخته منچ
  for(let r=1;r<=15;r++) for(let c=1;c<=15;c++){
    svg.append("rect").attr("x",getX(c)).attr("y",getY(r)).attr("width",CELL).attr("height",CELL).attr("fill","#fff").attr("stroke","#dbe2ea");
  }
  // (کدهای رنگی پایه و مسیرها اینجا به صورت ساده‌شده رسم میشن)
  
  let myTurn = false;
  let boardState = [-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1];

  window.__rollDice = function() {
    if(!myTurn) return document.getElementById('turnText').textContent = "نوبت شما نیست!";
    NET.send("roll", {});
    document.getElementById('turnText').textContent = "در حال پرتاب تاس...";
  };

  NET.on('joined', (m) => {
    document.getElementById('turnText').textContent = "متصل شد. منتظر شروع...";
  });

  NET.on('game.started', () => {
    myTurn = true;
    document.getElementById('turnText').textContent = "نوبت شماست! تاس بیندازید.";
  });

  NET.on('game.rolled', (p) => {
    document.getElementById('turnText').textContent = "تاس: " + p.dice;
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
