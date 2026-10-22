// app.js - فقط مدیریت UI و مسیریابی
let currentUser = null;
const appEl = document.getElementById('app');

function renderLayout(content, activeTab) {
    return `
    <div class="app-shell">
        <div class="top-header">
            <div class="logo-text">Angel</div>
            <div class="header-avatar" onclick="window.location.hash='#profile'">
                <img src="https://api.dicebear.com/7.x/bottts/svg?seed=${currentUser ? currentUser.id : 'Guest'}" style="width:36px;height:36px;border-radius:50%;border:2px solid var(--accent-blue);" />
            </div>
        </div>
        <div class="page-content">${content}</div>
        <div class="bottom-nav">
            <div class="nav-item ${activeTab === 'home' ? 'active' : ''}" onclick="window.location.hash='#home'"><span class="nav-icon">🏠</span></div>
            <div class="nav-center-btn" onclick="window.location.hash='#games'"><span class="nav-icon">🎮</span></div>
            <div class="nav-item ${activeTab === 'profile' ? 'active' : ''}" onclick="window.location.hash='#profile'"><span class="nav-icon">👤</span></div>
        </div>
    </div>`;
}

function renderHomePage() {
    return `
        <div class="section-title">بازی‌های آنلاین</div>
        <div class="games-grid">
            <div class="game-card" onclick="joinGame('ludo')">
                <div class="game-icon">🎲</div>
                <div class="game-name">منچ (Ludo)</div>
            </div>
            <div class="game-card" onclick="alert('به‌زودی!')">
                <div class="game-icon">🎱</div>
                <div class="game-name">بیلیارد (به‌زودی)</div>
            </div>
        </div>
    `;
}

function renderProfilePage() {
    return `
        <div style="text-align:center; margin-top: 20px;">
            <img src="https://api.dicebear.com/7.x/bottts/svg?seed=${currentUser ? currentUser.id : '1'}" style="width:100px;height:100px;border-radius:50%;border:4px solid var(--accent-blue);" />
            <h2 style="margin-top: 10px;">User ${currentUser ? currentUser.id : 'Unknown'}</h2>
            <p style="color:var(--text-muted);">سطح ۱</p>
        </div>
    `;
}

function router() {
    if (!currentUser) return;
    const hash = window.location.hash || '#home';
    if (hash === '#home') appEl.innerHTML = renderLayout(renderHomePage(), 'home');
    else if (hash === '#profile') appEl.innerHTML = renderLayout(renderProfilePage(), 'profile');
    else if (hash === '#games') appEl.innerHTML = renderLayout(renderHomePage(), 'games');
}

async function init() {
    const urlParams = new URLSearchParams(location.search);
    let uid = urlParams.get('id');

    if (!uid && window.Telegram && window.Telegram.WebApp && window.Telegram.WebApp.initDataUnsafe) {
        uid = window.Telegram.WebApp.initDataUnsafe.user?.id;
    }
    if (!uid) uid = '1001'; // Guest ID for testing

    currentUser = { id: parseInt(uid) };
    
    // اتصال اولیه WebSocket در پس‌زمینه
    initWS(currentUser.id, handleGameMsg);
    
    if (!window.location.hash) window.location.hash = '#home';
    router();
}

window.addEventListener('hashchange', router);
window.addEventListener('load', init);
