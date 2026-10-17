import { User, GameStat } from './types';
const GAME_INFO: Record<string, { name: string, icon: string }> = { 'ludo': { name: 'منچ', icon: '🎯' }, 'backgammon': { name: 'بک‌گمون', icon: '🎲' }, 'chess': { name: 'شطرنج', icon: '♟️' }, 'football': { name: 'فوتبال دستی', icon: '⚽' }, 'billiards': { name: 'بیلیارد', icon: '🎱' } };
const DICEBEAR_STYLES = ['adventurer', 'avataaars', 'bottts', 'fun-emoji', 'lorelei', 'micah', 'miniavs', 'notionists', 'open-peeps', 'personas'];
export const AVATARS = Array.from({ length: 10 }, (_, i) => `https://api.dicebear.com/7.x/${DICEBEAR_STYLES[i % 10]}/svg?seed=Angel${i + 1}`);
export function escapeHtml(s: string): string { return String(s || '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c] as string)); }
function renderProfileCard(user: User, isLiked: boolean): string {
  const displayName = user.display_name || user.game_username || ('User' + user.id);
  const avatarContent = user.avatar ? `<img src="${escapeHtml(user.avatar)}">` : (user.name_emoji || '🎮');
  const heartColor = isLiked ? '#ef4444' : '#ffffff';
  return `
    <div class="profile-card">
      <div class="banner">
        <button class="settings-btn" id="btn-open-settings">⚙️</button>
        <button class="add-btn" id="btn-add-banner">+</button>
      </div>
      <div class="avatar-container"><div class="avatar" id="btn-edit-avatar">${avatarContent}</div></div>
      <div class="profile-info" style="padding-top: 80px;">
        <div class="name" style="color:${user.profile_color || '#c084fc'}">${escapeHtml(displayName)}</div>
        <div class="username">@${escapeHtml(user.game_username || 'unknown')}</div>
        <div class="bio">${escapeHtml(user.bio || '---')}</div>
      </div>
      <div class="main-stats">
        <div class="stat-box" id="like-box" style="cursor:pointer;">
          <div class="stat-value" id="heart-icon" style="color:${heartColor}; font-size:24px; line-height:1;">❤</div>
          <div class="stat-label" id="likes-count" style="margin-top:4px;">${user.likes_count || 0}</div>
        </div>
        <div class="stat-box"><div class="stat-value">${user.age_days || 0}</div><div class="stat-label">روز</div></div>
        <div class="stat-box"><div class="stat-value">${user.wins || 0}</div><div class="stat-label">برد</div></div>
      </div>
      <div class="action-buttons">
        <button class="btn-square blue" id="btn-edit-bio">✏️</button>
        <button class="btn-square green">💎 ${user.gems || 0}</button>
        <button class="btn-square">🪙 ${user.coins || 0}</button>
      </div>
      <div style="padding: 0 20px 20px;"><button class="btn-hype">هایپ 🚀</button></div>
    </div>`;
}
function renderGameStats(stats: GameStat[]): string {
  if (!stats || stats.length === 0) return `<div class="game-stats-card"><div style="text-align:center;padding:20px;color:#8b949e;">هنوز آماری ثبت نشده است.</div></div>`;
  let html = `<div class="game-stats-card">`;
  stats.forEach(g => {
    const info = GAME_INFO[g.game_id] || { name: g.game_id, icon: '🎮' };
    const winRate = g.games_played > 0 ? Math.round((g.wins / g.games_played) * 100) : 0;
    html += `<div class="game-item"><div class="game-icon">${info.icon}</div><div class="game-details"><div class="game-header"><div><span class="game-name">${info.name}</span></div><div class="game-level">مرحله ${g.level}</div></div><div class="progress-bar"><div class="progress-fill" style="width: ${winRate}%"></div></div><div class="game-meta"><span>${g.games_played} بازی کرده</span><span>${g.wins} برنده شده</span></div></div><div class="game-star">${g.level}</div></div>`;
  });
  return html + `</div>`;
}
function renderEditModal(user: User): string {
  let gridHtml = '';
  AVATARS.forEach(url => { gridHtml += `<img src="${url}" class="avatar-option ${user.avatar === url ? 'selected' : ''}" data-url="${url}" />`; });
  return `<div id="edit-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">ویرایش پروفایل</h3><div class="avatar-grid" id="avatar-grid">${gridHtml}</div><textarea id="bio-input" placeholder="بیوگرافی...">${escapeHtml(user.bio || '')}</textarea><button class="modal-btn save" id="save-profile-btn">ذخیره</button><button class="modal-btn close" id="close-modal-btn">انصراف</button></div></div>`;
}
function renderSettingsModal(): string {
  return `
    <div id="settings-modal" class="modal">
      <div class="modal-content" style="padding: 15px;">
        <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:15px;"><h3 style="margin:0;">تنظیمات</h3><button id="close-settings-btn" style="background:none; border:none; color:white; font-size:20px; cursor:pointer;">✖</button></div>
        <ul class="settings-list">
            <li class="settings-item"><span>حساب</span> <span class="settings-icon">👤</span></li>
            <li class="settings-item"><span>دستگاه‌ها</span> <span class="settings-icon">📱</span></li>
            <li class="settings-item"><span>زبان</span> <span class="settings-icon">🌍</span></li>
            <li class="settings-item"><span>بازی‌ها</span> <span class="settings-icon">🎮</span></li>
            <li class="settings-item"><span>اعلان‌ها</span> <span class="settings-icon">🔔</span></li>
            <li class="settings-item"><span>حریم خصوصی</span> <span class="settings-icon">🔒</span></li>
            <li class="settings-item"><span>تم‌ها</span> <span class="settings-icon">🎨</span></li>
            <li class="settings-item"><span>درباره آنجل</span> <span class="settings-icon">ⓘ</span></li>
            <li class="settings-item"><span>دریافت پشتیبانی</span> <span class="settings-icon">❓</span></li>
        </ul>
      </div>
    </div>`;
}
export function render(user: User, stats: GameStat[], appEl: HTMLElement, onSaveProfile: (avatar: string, bio: string) => void, onLike: () => void) {
  const isLiked = localStorage.getItem('liked_' + user.id) === 'true';
  appEl.innerHTML = `<div class="container">${renderProfileCard(user, isLiked)}<div class="section-title" style="margin-bottom: 10px; color:#8b949e;">آمار من</div>${renderGameStats(stats)}</div>${renderEditModal(user)}${renderSettingsModal()}`;
  let selectedAvatar = user.avatar || '';
  const editModal = document.getElementById('edit-modal');
  [document.getElementById('btn-edit-avatar'), document.getElementById('btn-edit-bio')].forEach(btn => btn?.addEventListener('click', () => editModal?.classList.add('active')));
  document.getElementById('close-modal-btn')?.addEventListener('click', () => editModal?.classList.remove('active'));
  document.querySelectorAll('.avatar-option').forEach(img => img.addEventListener('click', (e) => { document.querySelectorAll('.avatar-option').forEach(i => i.classList.remove('selected')); (e.target as HTMLElement).classList.add('selected'); selectedAvatar = (e.target as HTMLElement).getAttribute('data-url') || ''; }));
  document.getElementById('save-profile-btn')?.addEventListener('click', () => onSaveProfile(selectedAvatar, (document.getElementById('bio-input') as HTMLTextAreaElement).value));
  document.getElementById('like-box')?.addEventListener('click', onLike);
  const settingsModal = document.getElementById('settings-modal');
  document.getElementById('btn-open-settings')?.addEventListener('click', () => settingsModal?.classList.add('active'));
  document.getElementById('close-settings-btn')?.addEventListener('click', () => settingsModal?.classList.remove('active'));
  document.getElementById('btn-add-banner')?.addEventListener('click', () => alert('آپلود بنر به زودی...'));
}