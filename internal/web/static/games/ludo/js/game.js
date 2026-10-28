import { drawBoard } from './board.js';
import { Network } from './network.js';

if (window.Telegram && Telegram.WebApp) {
    Telegram.WebApp.ready();
    Telegram.WebApp.expand();
}

document.addEventListener('DOMContentLoaded', () => {
    const statusEl = document.getElementById('status-text');
    const rollBtn = document.getElementById('roll-btn');

    drawBoard();

    Network.init((msg) => {
        statusEl.textContent = msg;
    });

    rollBtn.addEventListener('click', () => {
        if (Network.connected) {
            statusEl.textContent = "در حال پرتاب تاس...";
            rollBtn.disabled = true;
            Network.send('roll');
        } else {
            statusEl.textContent = "اتصال به سرور برقرار نیست.";
        }
    });
});
