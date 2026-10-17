
const GAME_INFO = { 'ludo': { name: 'منچ', icon: '🎯' }, 'backgammon': { name: 'بک‌گمون', icon: '🎲' }, 'chess': { name: 'شطرنج', icon: '♟️' }, 'football': { name: 'فوتبال دستی', icon: '⚽' }, 'billiards': { name: 'بیلیارد', icon: '🎱' } };
const DICEBEAR_STYLES = ['adventurer', 'avataaars', 'bottts', 'fun-emoji', 'lorelei', 'micah', 'miniavs', 'notionists', 'open-peeps', 'personas'];
const AVATARS = Array.from({ length: 10 }, (_, i) => `https://api.dicebear.com/7.x/${DICEBEAR_STYLES[i % 10]}/svg?seed=Angel${i + 1}`);
function escapeHtml(s) { return String(s || '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c])); }

function render(user, stats, appEl, onSaveProfile) {
    const displayName = user.display_name || user.game_username || ('User' + user.id);
    const avatarContent = user.avatar ? `<img src="${escapeHtml(user.avatar)}">` : (user.name_emoji || '🎮');
    
    let statsHtml = '';
    if (!stats || stats.length === 0) {
        statsHtml = `<div class="game-stats-card"><div style="text-align:center;padding:20px;color:#8b949e;">هنوز آماری ثبت نشده است.</div></div>`;
    } else {
        statsHtml = `<div class="game-stats-card">`;
        stats.forEach(g => {
            const info = GAME_INFO[g.game_id] || { name: g.game_id, icon: '🎮' };
            const winRate = g.games_played > 0 ? Math.round((g.wins / g.games_played) * 100) : 0;
            statsHtml += `<div class="game-item"><div class="game-icon">${info.icon}</div><div class="game-details"><div class="game-header"><div><span class="game-name">${info.name}</span></div><div class="game-level">مرحله ${g.level}</div></div><div class="progress-bar"><div class="progress-fill" style="width: ${winRate}%"></div></div><div class="game-meta"><span>${g.games_played} بازی کرده</span><span>${g.wins} برنده شده</span></div></div><div class="game-star">${g.level}</div></div>`;
        });
        statsHtml += `</div>`;
    }

    let gridHtml = '';
    AVATARS.forEach(url => { gridHtml += `<img src="${url}" class="avatar-option ${user.avatar === url ? 'selected' : ''}" data-url="${url}" />`; });

    appEl.innerHTML = `
    <div class="container">
        <div class="profile-card">
            <div class="banner"></div>
            <div class="avatar-container"><div class="avatar" id="btn-edit-avatar">${avatarContent}</div></div>
            <div class="profile-info">
                <div class="name" style="color:${user.profile_color || '#c084fc'}">${escapeHtml(displayName)}</div>
                <div class="username">@${escapeHtml(user.game_username || 'unknown')}</div>
                <div class="bio">${escapeHtml(user.bio || '---')}</div>
            </div>
            <div class="main-stats">
                <div class="stat-box"><div class="stat-value">${user.likes_count || 0}</div><div class="stat-label">هاله</div></div>
                <div class="stat-box"><div class="stat-value">${user.age_days || 0}</div><div class="stat-label">روز</div></div>
                <div class="stat-box"><div class="stat-value">${user.wins || 0}</div><div class="stat-label">برد</div></div>
            </div>
            <div class="action-buttons">
                <button class="btn-square blue" id="btn-edit-bio">✏️</button>
                <button class="btn-square green">💎 ${user.gems || 0}</button>
                <button class="btn-square">🪙 ${user.coins || 0}</button>
            </div>
            <div style="padding: 0 20px 20px;"><button class="btn-hype">هایپ 🚀</button></div>
        </div>
        <div class="section-title" style="margin-bottom: 10px; color:#8b949e;">آمار من</div>
        ${statsHtml}
    </div>
    <div id="edit-modal" class="modal">
        <div class="modal-content">
            <h3 style="text-align:center; margin-bottom:15px;">ویرایش پروفایل</h3>
            <div class="avatar-grid" id="avatar-grid">${gridHtml}</div>
            <textarea id="bio-input" placeholder="بیوگرافی...">${escapeHtml(user.bio || '')}</textarea>
            <button class="modal-btn save" id="save-profile-btn">ذخیره</button>
            <button class="modal-btn close" id="close-modal-btn">انصراف</button>
        </div>
    </div>`;

    let selectedAvatar = user.avatar || '';
    const modal = document.getElementById('edit-modal');
    [document.getElementById('btn-edit-avatar'), document.getElementById('btn-edit-bio')].forEach(btn => btn?.addEventListener('click', () => modal?.classList.add('active')));
    document.getElementById('close-modal-btn')?.addEventListener('click', () => modal?.classList.remove('active'));
    document.querySelectorAll('.avatar-option').forEach(img => img.addEventListener('click', (e) => { document.querySelectorAll('.avatar-option').forEach(i => i.classList.remove('selected')); (e.target).classList.add('selected'); selectedAvatar = (e.target).getAttribute('data-url') || ''; }));
    document.getElementById('save-profile-btn')?.addEventListener('click', () => onSaveProfile(selectedAvatar, (document.getElementById('bio-input')).value));
}

async function fetchToken(userId) { try { const r = await fetch('/api/auth/token?id=' + userId); const d = await r.json(); return d.ok ? d.token : null; } catch { return null; } }
async function updateProfile(patch) {
    const p = new URLSearchParams(location.search).get('id'); if (!p) return false;
    const token = await fetchToken(parseInt(p, 10)); if (!token) return false;
    try { const r = await fetch('/api/me', { method: 'PATCH', headers: { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' }, body: JSON.stringify(patch) }); return r.ok; } catch { return false; }
}

const appEl = document.getElementById('app');
let currentUser = null;
let currentStats = [];
async function loadProfile() {
    const p = new URLSearchParams(location.search).get('id');
    if (!p) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">No ID</div>'; return; }
    try {
        const r = await fetch('/api/profile/' + p);
        const d = await r.json();
        if (!d.ok) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">Not Found</div>'; return; }
        currentUser = d.user; currentStats = d.stats || [];
        if (appEl && currentUser) render(currentUser, currentStats, appEl, handleSaveProfile);
    } catch(e) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">Error</div>'; }
}
async function handleSaveProfile(avatar, bio) {
    if (!currentUser) return;
    const ok = await updateProfile({ avatar, bio });
    if (ok && appEl) { currentUser.avatar = avatar; currentUser.bio = bio; render(currentUser, currentStats, appEl, handleSaveProfile); document.getElementById('edit-modal')?.classList.remove('active'); } else { alert('خطا در ذخیره اطلاعات'); }
}
loadProfile();
