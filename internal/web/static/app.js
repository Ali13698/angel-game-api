const appEl = document.getElementById('app');
let currentUser = null;
let currentStats = [];

function escapeHtml(s) {
  return String(s || '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
}

const ICON_GEAR = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path><circle cx="12" cy="12" r="3"></circle></svg>`;
const ICON_BACK = `<svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>`;
const ICON_PLUS = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>`;
const ICON_PENCIL = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#8b949e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>`;
const ICON_HEART = `<svg width="28" height="28" viewBox="0 0 24 24" fill="{fill}" stroke="{stroke}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`;

function renderLayout(content, activeTab) {
    return `
    <div class="app-shell">
        <div class="top-header">
            <div class="header-spacer"></div>
            <div class="header-logo"><span class="logo-text">Angel</span></div>
            <div class="header-avatar" onclick="window.location.hash='#profile'">
                <img src="https://api.dicebear.com/7.x/bottts/svg?seed=Angel" alt="Profile" />
            </div>
        </div>
        <div class="page-content">${content}</div>
        <div class="bottom-nav">
            <div class="nav-item ${activeTab === 'home' ? 'active' : ''}" onclick="window.location.hash='#home'"><span class="nav-icon">🏠</span></div>
            <div class="nav-item"><span class="nav-icon">👥</span></div>
            <div class="nav-center-btn ${activeTab === 'games' ? 'active' : ''}" onclick="window.location.hash='#games'"><span class="nav-icon">🎮</span></div>
            <div class="nav-item"><span class="nav-icon">💬</span></div>
            <div class="nav-item ${activeTab === 'profile' ? 'active' : ''}" onclick="window.location.hash='#profile'"><span class="nav-icon">👤</span></div>
        </div>
    </div>`;
}

function renderHomePage(user) {
    // این آرایه‌ها الان خالی هستند. بعداً می‌تونی دیتای واقعی رو از API بگیری و اینجا جایگزین کنی.
    const onlineUsers = []; 
    const games = []; 

    return `
        <div class="home-tabs"><div class="tab active">آخرین</div><div class="tab">اخبار</div></div>
        
        ${onlineUsers.length > 0 ? `
            <div class="section-title">آنلاین / موارد دلخواه</div>
            <div class="online-slider">
                ${onlineUsers.map(u => `<div class="online-card"><div class="online-avatar"><img src="https://api.dicebear.com/7.x/adventurer/svg?seed=${u}" /></div><div class="online-name">${u}</div></div>`).join('')}
            </div>
        ` : ''}

        <div class="section-title">بازی ها</div>
        <div class="games-slider">
            ${games.length > 0 ? games.map(g => `<div class="game-card" onclick="window.location.hash='#games'"><div class="game-bg" style="background: ${g.image};"></div><div class="game-info"><span class="game-name">${g.name}</span><button class="play-btn">بازی</button></div></div>`).join('') : '<div style="color:#8b949e; text-align:center; width:100%; padding:20px; font-size:14px;">هنوز بازی‌ای اضافه نشده است</div>'}
        </div>

        <div class="daily-reward-section">
            <div class="reward-header"><span>جایزه روزانه</span><span class="reward-time">۲ ساعت و ۱۵ دقیقه</span></div>
            <div class="reward-title">برنده ۲ بازی</div>
            <div class="quest-item"><div class="quest-icon">8</div><div class="quest-info"><div>۲ اوچو را برنده شوید</div><div class="progress-bar"><div class="progress-fill" style="width:20%;"></div></div></div></div>
        </div>`;
}

function renderGamesPage() {
    const games = []; 
    return `
        <div class="page-header" style="padding: 10px 0;"><h2>بازی ها</h2><div class="header-icons"><span>↑↓</span></div></div>
        <div class="games-tabs"><div class="gtab active">همه</div><div class="gtab">جدید</div><div class="gtab">کارت</div><div class="gtab">تخته</div></div>
        <div class="games-grid">
            ${games.length > 0 ? games.map(g => `<div class="game-grid-card"><div class="grid-game-bg" style="background: ${g.image};"></div><div class="grid-game-name">${g.name}</div></div>`).join('') : '<div style="color:#8b949e; text-align:center; grid-column: 1 / -1; padding:20px;">هنوز بازی‌ای اضافه نشده است</div>'}
        </div>`;
}

function renderProfileCard(user, isLiked) {
  const displayName = user.display_name || user.game_username || ('User' + user.id);
  const avatarContent = user.avatar ? `<img src="${escapeHtml(user.avatar)}">` : (user.name_emoji || '🎮');
  const heartFill = isLiked ? '#ef4444' : 'none';
  const heartStroke = isLiked ? '#ef4444' : '#ffffff';
  const heartSvg = ICON_HEART.replace('{fill}', heartFill).replace('{stroke}', heartStroke);

  return `
    <div class="profile-card" style="margin-top: 20px;">
      <div class="top-bar">
        <button class="icon-btn" id="btn-open-settings">${ICON_GEAR}</button>
        <button class="icon-btn" onclick="window.location.hash='#home'">${ICON_BACK}</button>
      </div>
      <div class="banner" style="background: ${user.banner_id || '#7c3aed'}; background-size: cover; background-position: center;">
        <button class="plus-badge banner-plus">${ICON_PLUS}</button>
      </div>
      <div class="avatar-container">
        <div class="avatar">${avatarContent}</div>
        <button class="plus-badge avatar-plus">${ICON_PLUS}</button>
      </div>
      <div class="profile-info">
        <div class="name" style="color:${user.profile_color || '#c084fc'}">${escapeHtml(displayName)}</div>
        <div class="username">@${escapeHtml(user.game_username || 'unknown')}</div>
        <div class="bio-wrapper"><span class="bio">${escapeHtml(user.bio || '---')}</span><button class="pencil-btn">${ICON_PENCIL}</button></div>
      </div>
      <div class="main-stats">
        <div class="stat-box"><div class="stat-value" id="heart-icon">${heartSvg}</div><div class="stat-label" id="likes-count">${user.likes_count || 0}</div></div>
        <div class="stat-box"><div class="stat-value">${user.age_days || 0}</div><div class="stat-label">روز</div></div>
        <div class="stat-box"><div class="stat-value">${user.wins || 0}</div><div class="stat-label">برد</div></div>
      </div>
      <div class="action-buttons">
        <button class="btn-square-small green"><span class="icon-3d">💎</span> ${user.gems || 0}</button>
        <button class="btn-square-small"><span class="icon-3d">🪙</span> ${user.coins || 0}</button>
      </div>
    </div>`;
}

function router() {
    if (!currentUser || !appEl) return;
    const hash = window.location.hash || '#home';
    if (hash === '#home') {
        appEl.innerHTML = renderLayout(renderHomePage(currentUser), 'home');
    } else if (hash === '#games') {
        appEl.innerHTML = renderLayout(renderGamesPage(), 'games');
    } else if (hash === '#profile') {
        appEl.innerHTML = renderLayout(`<div class="container">${renderProfileCard(currentUser, false)}</div>`, 'profile');
    } else {
        window.location.hash = '#home';
    }
}

window.addEventListener('hashchange', router);

async function loadUser() {
    let p = new URLSearchParams(location.search).get('id');
    
    if (!p) {
        try {
            if (window.Telegram && window.Telegram.WebApp && window.Telegram.WebApp.initDataUnsafe && window.Telegram.WebApp.initDataUnsafe.user) {
                p = window.Telegram.WebApp.initDataUnsafe.user.id;
            } else if (window.Bale && window.Bale.WebApp && window.Bale.WebApp.initDataUnsafe && window.Bale.WebApp.initDataUnsafe.user) {
                p = window.Bale.WebApp.initDataUnsafe.user.id;
            }
        } catch (e) { console.error(e); }
    }

    if (!p) { if (appEl) appEl.innerHTML = 'No ID (Open from Bot)'; return false; }
    try {
        const r = await fetch('/api/profile/' + p);
        const d = await r.json();
        if (!d.ok) throw new Error();
        currentUser = d.user;
        currentStats = d.stats || [];
        return true;
    } catch (e) { if (appEl) appEl.innerHTML = 'Error Loading User'; return false; }
}

(async () => {
    const ok = await loadUser();
    if (ok) {
        if (!window.location.hash) window.location.hash = '#home';
        else router();
    }
})();
