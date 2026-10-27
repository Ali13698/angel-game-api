import { drawBoard } from './board.js';
import { Network } from './network.js';

if (window.Telegram && Telegram.WebApp) {
    Telegram.WebApp.ready();
    Telegram.WebApp.expand();
}

document.addEventListener('DOMContentLoaded', () => {
    const statusEl = document.getElementById('status-text');
    const rollBtn = document.getElementById('roll-btn');

    // ۱. رسم تخته
    drawBoard();

    // ۲. راه‌اندازی شبکه
    Network.init((msg) => {
        statusEl.textContent = msg;
    });

    // ۳. رویداد دکمه تاس
    rollBtn.addEventListener('click', () => {
        if (Network.connected) {
            statusEl.textContent = "در حال پرتاب تاس...";
            rollBtn.disabled = true;
            Network.send('roll');
            // (در فاز بعدی، با دریافت پاسخ سرور، دکمه مجدداً فعال می‌شود)
        } else {
            statusEl.textContent = "اتصال به سرور برقرار نیست.";
        }
    });
});
