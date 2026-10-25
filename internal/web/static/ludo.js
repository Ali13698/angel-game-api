(function(global){
  'use strict';
  const NET = global.AngelNet;
  let myTurn = false;

  window.__rollDice = function() {
    if(!myTurn) return document.getElementById('turnText').textContent = "نوبت شما نیست!";
    NET.send("roll", {});
    document.getElementById('turnText').textContent = "در حال پرتاب تاس...";
  };

  NET.on('joined', (m) => {
    document.getElementById('turnText').textContent = "بازی متصل شد. منتظر شروع...";
  });

  NET.on('game.started', () => {
    myTurn = true;
    document.getElementById('turnText').textContent = "بازی شروع شد! نوبت شماست.";
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
