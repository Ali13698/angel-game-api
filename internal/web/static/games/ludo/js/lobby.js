// تابع مستقل برای گرفتن آیدی کاربر
function getUserId() {
    const qs = new URLSearchParams(window.location.search).get("id");
    if (qs) return qs;
    const tg = window.Telegram?.WebApp?.initDataUnsafe?.user?.id;
    return tg ? String(tg) : "1001"; // Fallback for testing
}

const API_BASE = '/api';

document.addEventListener('DOMContentLoaded', async () => {
    if (window.Telegram && Telegram.WebApp) { Telegram.WebApp.ready(); Telegram.WebApp.expand(); }
    
    const userId = getUserId();
    const favBtn = document.getElementById('fav-btn');
    const settingsModal = document.getElementById('settings-modal');
    const createPrivateBtn = document.getElementById('create-private-btn');
    const startGameBtn = document.getElementById('start-game-btn');
    const closeModal = document.getElementById('close-modal');

    // --- ۱. دریافت آمار زنده ---
    try {
        const res = await fetch(`${API_BASE}/profile/${userId}`);
        const data = await res.json();
        if (data.ok && data.user) {
            document.getElementById('live-score').textContent = `امتیاز: ${data.user.xp || 0}`;
            
            const ludoStats = data.stats ? data.stats.find(s => s.game_id === 'ludo') : null;
            if (ludoStats) {
                document.getElementById('live-games').textContent = `${ludoStats.games_played} بازی`;
                document.getElementById('live-wins').textContent = `${ludoStats.wins} برد`;
                document.getElementById('live-rank').textContent = `رتبه: ${ludoStats.level > 5 ? 'Pro' : 'Gamer'}`;
            } else {
                document.getElementById('live-games').textContent = `0 بازی`;
                document.getElementById('live-wins').textContent = `0 برد`;
                document.getElementById('live-rank').textContent = `رتبه: Newbie`;
            }
        }
    } catch (e) { console.warn("Stats fetch failed", e); }

    // --- ۲. بررسی وضعیت علاقه‌مندی (قلب) ---
    try {
        const res = await fetch(`${API_BASE}/favorites`, {
            headers: { 'Authorization': 'Bearer ' + btoa(userId) }
        });
        const data = await res.json();
        if (data.ok && data.favorites.includes('ludo')) {
            favBtn.classList.add('active');
            favBtn.textContent = '♥';
        }
    } catch (e) { console.warn("Favorites fetch failed", e); }

    // --- ۳. کلیک روی قلب ---
    favBtn.addEventListener('click', async () => {
        try {
            const res = await fetch(`${API_BASE}/favorites/toggle`, {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + btoa(userId)
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

    // --- ۴. باز و بسته کردن مودال ---
    createPrivateBtn.addEventListener('click', () => settingsModal.classList.add('active'));
    closeModal.addEventListener('click', () => settingsModal.classList.remove('active'));

    document.querySelectorAll('.options span').forEach(opt => {
        opt.addEventListener('click', function() {
            this.parentElement.querySelectorAll('span').forEach(s => s.classList.remove('active'));
            this.classList.add('active');
        });
    });

    // --- ۵. شروع بازی ---
    startGameBtn.addEventListener('click', () => {
        window.location.href = `index.html?mode=private`;
    });
});

window.startQuickMatch = function() {
    window.location.href = `index.html`;
};
