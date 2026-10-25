(function() {
  const NET = window.AngelNet;
  const svg = d3.select("#ludo");
  let myTurn = false;
  let boardState = [-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1];
  
  // رسم یک تخته ساده برای تست (در نسخه نهایی با کدهای D3 خودت جایگزین میشه)
  svg.append("rect").attr("width", 600).attr("height", 600).attr("fill", "#1E1E38");
  
  document.getElementById('turnText').onclick = () => {
    if(!myTurn) return alert("نوبت شما نیست!");
    NET.send("roll", {});
    document.getElementById('turnText').innerText = "در حال پرتاب تاس...";
  };

  NET.on('joined', () => { document.getElementById('turnText').innerText = "منتظر شروع..."; });
  NET.on('game.rolled', (p) => { 
    document.getElementById('turnText').innerText = "تاس: " + p.dice + " (برای حرکت کلیک کنید)";
  });
  NET.on('turn_changed', (p) => {
    myTurn = (p.turn === 0);
    document.getElementById('turnText').innerText = myTurn ? "نوبت شماست (کلیک کنید)" : "نوبت حریف...";
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
