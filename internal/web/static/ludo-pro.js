(function() {
  const NET = window.AngelNet;
  let myTurn = false;

  window.__rollDice = () => {
    if(!myTurn) return alert("نوبت شما نیست!");
    NET.send("roll", {});
    document.getElementById('turnText').innerText = "در حال پرتاب تاس...";
  };

  NET.on('joined', () => { document.getElementById('turnText').innerText = "بازی متصل شد. منتظر شروع..."; });
  NET.on('game.rolled', (p) => { 
    document.getElementById('turnText').innerText = "تاس: " + p.dice + " (منتظر حرکت...)";
    // در نسخه کامل، دکمه های حرکت مهره ها اینجا فعال میشن
  });
  NET.on('turn_changed', (p) => {
    myTurn = (p.turn === 0);
    document.getElementById('turnText').innerText = myTurn ? "نوبت شماست! تاس بیندازید." : "نوبت حریف...";
  });
  NET.on('game.moved', (p) => { console.log("Moved:", p); });

  async function boot() {
    const user = new URLSearchParams(location.search).get('user') || '1001';
    const room = new URLSearchParams(location.search).get('room');
    await NET.connect(user, { roomId: room });
    NET.join(room);
  }
  boot();
})();
