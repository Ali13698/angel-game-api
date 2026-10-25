(function(global) {
  const NET = global.AngelNet;
  const boardEl = document.getElementById('board');
  const statusEl = document.getElementById('status');
  let myTurn = false;
  let gameActive = false;
  let boardState = [0,0,0,0,0,0,0,0,0];
  
  let turnTimer = null;
  let timeLeft = 10;

  for (let i = 0; i < 9; i++) {
    const cell = document.createElement('div');
    cell.className = 'cell';
    cell.dataset.idx = i;
    cell.onclick = () => {
      if (!gameActive) return;
      if (!myTurn) return statusEl.textContent = "نوبت شما نیست!";
      if (boardState[i] !== 0) return;
      
      // ✅ تایمر را اینجا متوقف نمی‌کنیم! منتظر تایید سرور می‌مانیم
      statusEl.textContent = "در حال ارسال به سرور...";
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

  function startTimer() {
    stopTimer();
    if (!gameActive || !myTurn) return;
    timeLeft = 10;
    statusEl.textContent = `نوبت شماست (${timeLeft})`;
    turnTimer = setInterval(() => {
      timeLeft--;
      if (timeLeft >= 0) {
        statusEl.textContent = `نوبت شماست (${timeLeft})`;
      }
      if (timeLeft < 0) stopTimer();
    }, 1000);
  }

  function stopTimer() {
    if (turnTimer) clearInterval(turnTimer);
    turnTimer = null;
  }

  NET.on('joined', () => { statusEl.textContent = "منتظر شروع..."; });
  
  NET.on('state', (st) => { 
      if (st && st.board) { boardState = st.board; render(); } 
  });

  NET.on('game.started', () => {
    gameActive = true;
    myTurn = true;
    startTimer();
  });

  // ✅ وقتی سرور تایید کرد که حرکت انجام شد، تایمر را متوقف کن
  NET.on('game.moved', (p) => {
    boardState[p.cell] = p.player;
    render();
    stopTimer(); 
  });

  NET.on('turn_changed', (p) => {
    myTurn = (p.turn === 0);
    if (myTurn) {
      startTimer();
    } else {
      stopTimer();
      statusEl.textContent = "نوبت حریف...";
    }
  });

  NET.on('game.over', (p) => {
    gameActive = false;
    myTurn = false;
    stopTimer();
    if (p.winner === 1) statusEl.textContent = "🎉 شما برنده شدید!";
    else if (p.winner === 2) statusEl.textContent = "🤖 ربات برنده شد!";
    else statusEl.textContent = "🤝 بازی مساوی شد!";
    if (p.reason === "timeout") statusEl.textContent += " (حریف اوت شد)";
  });

  // ✅ اگر سرور خطا داد، تایمر را دوباره روشن کن تا کاربر وقت داشته باشد
  NET.on('error', (e) => { 
    statusEl.textContent = "خطا: " + (e.message || e.code); 
    if (myTurn && gameActive) startTimer(); 
  });

  async function boot() {
    const user = new URLSearchParams(location.search).get('user') || '1001';
    try {
        let res = await fetch("/rooms", {
            method: "POST", headers: {"Content-Type": "application/json"},
            body: JSON.stringify({game_id: "tic-tac-toe", players: [user, "bot"]})
        });
        let data = await res.json();
        await NET.connect(user, { roomId: data.room_id });
        NET.join(data.room_id);
    } catch(e) { statusEl.textContent = "خطا: " + e.message; }
  }
  boot();
})(window);
