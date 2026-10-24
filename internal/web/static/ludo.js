(function (global) {
  'use strict';
  const NET = global.AngelNet, T = NET.TYPES;
  const $ = (id) => document.getElementById(id);
  let mySeat = null, myUser = null;

  // ⚠️ در اینجا فقط استایل و لاگیک ساده برای تست گذاشته شده تا از کدهای قدیمی دور شویم
  global.__rollDice = function () { NET.roll(); };

  function setText(t) { const el = $('turnText'); if (el) el.textContent = t; }

  NET.on('state', (st) => {
      setText("اتاق متصل شد!");
  });
  NET.on(T.ROLLED, (p) => { console.log("Dice rolled:", p.die); setText("تاس: " + p.die); });
  NET.on(T.JOINED, (m) => { if (m && m.roomId) NET.roomId = m.roomId; setText("شما وارد شدید."); });
  NET.on(T.ERROR, (e) => { setText('خطا: ' + (e.message || e.code)); });
  NET.on('closed', () => setText('اتصال قطع شد...'));

  async function boot() {
    const params = new URLSearchParams(location.search);
    const user = params.get('user') || '1001';
    const room = params.get('room') || 'main';
    myUser = user;
    
    setText("در حال اتصال...");
    await NET.connect(user, { roomId: room });
    NET.join(room);
  }

  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', boot);
  else boot();
})(window);
