import { getUserId } from '../../../lobby/src/api.js';

const API_BASE = '/api';

document.addEventListener('DOMContentLoaded', async () => {
    if (window.Telegram && Telegram.WebApp) { Telegram.WebApp.ready(); Telegram.WebApp.expand(); }
    
    const favBtn = document.getElementById('fav-btn');
    const settingsModal = document.getElementById('settings-modal');
    const createPrivateBtn = document.getElementById('create-private-btn');
    const startGameBtn = document.getElementById('start-game-btn');
    const closeModal = document.getElementById('close-modal');

    // ۱. بررسی وضعیت علاقه‌مندی
    try {
        const res = await fetch(`${API_BASE}/favorites`, {
            headers: { 'Authorization': 'Bearer ' + btoa(getUserId()) }
        });
        const data = await res.json();
        if (data.ok && data.favorites.includes('ludo')) {
            favBtn.classList.add('active');
            favBtn.textContent = '♥';
        }
    } catch (e) { console.warn("Favorites fetch failed"); }

    // ۲. کلیک روی قلب
    favBtn.addEventListener('click', async () => {
        try {
            const res = await fetch(`${API_BASE}/favorites/toggle`, {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + btoa(getUserId())
                },
                body: JSON.stringify({ game_id: 'ludo' })
            });
            const data = await res.json();
            if (data.ok) {
                const isAdded = data.status === 'added';
                favBtn.classList.toggle('active', isAdded);
                favBtn.textContent = isAdded ? '♥' : '♡';
            }
        } catch (e) { console.error(e); }
    });

    // ۳. باز و بسته کردن مودال تنظیمات
    createPrivateBtn.addEventListener('click', () => settingsModal.classList.add('active'));
    closeModal.addEventListener('click', () => settingsModal.classList.remove('active'));

    // ۴. تنظیمات کلیک (فقط برای ظاهر)
    document.querySelectorAll('.options span').forEach(opt => {
        opt.addEventListener('click', function() {
            this.parentElement.querySelectorAll('span').forEach(s => s.classList.remove('active'));
            this.classList.add('active');
        });
    });

    // ۵. ساخت اتاق خصوصی و رفتن به صفحه بازی
    startGameBtn.addEventListener('click', () => {
        window.location.href = `index.html?mode=private`;
    });
});

window.startQuickMatch = function() {
    window.location.href = `index.html`;
};
