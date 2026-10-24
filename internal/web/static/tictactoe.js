(function(global) {
  const NET = global.AngelNet;
  const boardEl = document.getElementById('board');
  const statusEl = document.getElementById('status');
  let myTurn = false;
  let boardState = [0,0,0,0,0,0,0,0,0];

  // ساخت خانه‌ها
  for (let i = 0; i < 9; i++) {
    const cell = document.createElement('div');
    cell.className = 'cell';
    cell.dataset.idx = i;
    cell.onclick = () => {
      if (!myTurn || boardState[i] !== 0) return;
      NET.move(i);
    };
    boardEl.appendChild(cell);
  }

  function render() {
    for (let i = 0; i < 9; i++) {
      const cell = boardEl.children[i];
      if (boardState[i] === 1) { cell.textContent = 'X'; cell.className = 'cell x'; }
      else if (boardState[i] === 2) { cell.textContent = 'O'; cell.className = 'cell o'; }
      else { cell.textContent = ''; cell.className = 'cell'; }
    }
  }

  NET.on('joined', (m) => {
    statusEl.textContent = "بازی شروع شد!";
    myTurn = true;
  });

  NET.on('state', (st) => {
    if (st && st.board) { boardState = st.board; render(); }
  });

  NET.on('game.moved', (p) => {
    boardState[p.cell] = p.player;
    render();
  });

  NET.on('turn_changed', (p) => {
    myTurn = (p.turn === 0);
    statusEl.textContent = myTurn ? "نوبت شماست" : "نوبت حریف...";
  });

  NET.on('error', (e) => { statusEl.textContent = "خطا: " + (e.message || e.code); });
  NET.on('closed', () => { statusEl.textContent = "اتصال قطع شد"; });

  // اتصال با ساخت اتاق
  async function boot() {
    const user = new URLSearchParams(location.search).get('user') || '1001';
    statusEl.textContent = "در حال ساخت اتاق در سرور...";
    
    try {
        // ساخت اتاق در سرور Go
        let res = await fetch("/rooms", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({game_id: "tic-tac-toe", players: [user, "bot"]})
        });
        let data = await res.json();
        
        await NET.connect(user, { roomId: data.room_id });
        NET.join(data.room_id);
    } catch(e) {
        statusEl.textContent = "خطا در ساخت اتاق: " + e.message;
    }
  }
  
  boot();
})(window);
