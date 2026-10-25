(function(global) {
  const NET = global.AngelNet;
  const boardEl = document.getElementById('board');
  const statusEl = document.getElementById('status');
  let myTurn = false;
  let gameActive = false;
  let boardState = [0,0,0,0,0,0,0,0,0];

  for (let i = 0; i < 9; i++) {
    const cell = document.createElement('div');
    cell.className = 'cell';
    cell.dataset.idx = i;
    cell.onclick = () => {
      if (!gameActive) return;
      if (!myTurn) return statusEl.textContent = "نوبت شما نیست!";
      if (boardState[i] !== 0) return;
      statusEl.textContent = "در حال ارسال...";
      NET.send("move", { cell: i });
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

  NET.on('joined', () => { statusEl.textContent = "منتظر شروع..."; });
  NET.on('state', (st) => { if (st && st.board) { boardState = st.board; render(); } });
  
  NET.on('game.started', () => {
    gameActive = true;
    myTurn = true;
    statusEl.textContent = "بازی شروع شد! نوبت شماست.";
  });

  NET.on('game.moved', (p) => { boardState[p.cell] = p.player; render(); });
  
  NET.on('turn_changed', (p) => {
    myTurn = (p.turn === 0);
    statusEl.textContent = myTurn ? "نوبت شماست!" : "نوبت حریف...";
  });

  NET.on('game.over', (p) => {
    gameActive = false;
    myTurn = false;
    if (p.winner === 1) statusEl.textContent = "🎉 شما برنده شدید!";
    else if (p.winner === 2) statusEl.textContent = "🤖 ربات برنده شد!";
    else statusEl.textContent = "🤝 بازی مساوی شد!";
  });

  async function boot() {
    const user = new URLSearchParams(location.search).get('user') || '1001';
    const room = new URLSearchParams(location.search).get('room');
    await NET.connect(user, { roomId: room });
    NET.join(room);
  }
  boot();
})(window);
