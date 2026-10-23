// app.js - UI and Routing
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
        <div class="profile-card">
            <div class="banner"></div>
            <div class="avatar-container">
                <div class="avatar">
                    <img src="https://api.dicebear.com/7.x/bottts/svg?seed=${currentUser ? currentUser.id : '1'}" />
                </div>
            </div>
            <div class="profile-info">
                <div class="name">User ${currentUser ? currentUser.id : 'Unknown'}</div>
                <div class="username">@user_${currentUser ? currentUser.id : '1'}</div>
            </div>
            <div class="main-stats">
                <div class="stat-box">
                    <div class="stat-value">0</div>
                    <div class="stat-label">لایک</div>
                </div>
                <div class="stat-box">
                    <div class="stat-value">1</div>
                    <div class="stat-label">روز</div>
                </div>
                <div class="stat-box">
                    <div class="stat-value">0</div>
                    <div class="stat-label">برد</div>
                </div>
            </div>
            <div style="padding: 20px; display: flex; gap: 10px; justify-content: center;">
                <div style="background:#1c2028; padding: 8px 16px; border-radius: 8px; font-size:14px;">💎 0</div>
                <div style="background:#1c2028; padding: 8px 16px; border-radius: 8px; font-size:14px;">🪙 0</div>
            </div>
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
    try {
        const urlParams = new URLSearchParams(location.search);
        let uid = urlParams.get('id');

        if (!uid && window.Telegram && window.Telegram.WebApp && window.Telegram.WebApp.initDataUnsafe) {
            uid = window.Telegram.WebApp.initDataUnsafe.user?.id;
        }
        if (!uid) uid = '1001';

        currentUser = { id: parseInt(uid) };
        
        // اتصال در پس‌زمینه
        initWS(currentUser.id, handleGameMsg);
        
        if (!window.location.hash) window.location.hash = '#home';
        router();
    } catch (e) {
        console.error("Init error:", e);
        appEl.innerHTML = `<div style="color:red; padding:20px; text-align:center;">خطا در بارگذاری: ${e.message}</div>`;
    }
}

window.addEventListener('hashchange', router);
window.addEventListener('load', init);
