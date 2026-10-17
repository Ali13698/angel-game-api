import { User, GameStat } from './types';

const ICON_GEAR = `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path><circle cx="12" cy="12" r="3"></circle></svg>`;
const ICON_BACK = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>`;
const ICON_PLUS = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>`;
const ICON_PENCIL = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>`;
const ICON_HEART = `<svg width="26" height="26" viewBox="0 0 24 24" fill="{fill}" stroke="{stroke}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`;

const GAME_INFO: Record<string, { name: string, icon: string }> = { 'ludo': { name: 'منچ', icon: '🎯' }, 'backgammon': { name: 'بک‌گمون', icon: '🎲' }, 'chess': { name: 'شطرنج', icon: '♟️' }, 'football': { name: 'فوتبال دستی', icon: '⚽' }, 'billiards': { name: 'بیلیارد', icon: '🎱' } };
const DICEBEAR_STYLES = ['adventurer', 'avataaars', 'bottts', 'fun-emoji', 'lorelei', 'micah', 'miniavs', 'notionists', 'open-peeps', 'personas'];
export const AVATARS = Array.from({ length: 10 }, (_, i) => `https://api.dicebear.com/7.x/${DICEBEAR_STYLES[i % 10]}/svg?seed=Angel${i + 1}`);

export function escapeHtml(s: string): string { return String(s || '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c] as string)); }

function renderProfileCard(user: User, isLiked: boolean): string {
  const displayName = user.display_name || user.game_username || ('User' + user.id);
  const avatarContent = user.avatar ? `<img src="${escapeHtml(user.avatar)}">` : (user.name_emoji || '🎮');
  const heartFill = isLiked ? '#ef4444' : 'none';
  const heartStroke = isLiked ? '#ef4444' : '#ffffff';
  const heartSvg = ICON_HEART.replace('{fill}', heartFill).replace('{stroke}', heartStroke);

  return `
    <div class="profile-card">
      <div class="banner">
        <div class="top-bar">
            <button class="icon-btn settings-btn" id="btn-open-settings">${ICON_GEAR}</button>
            <button class="icon-btn back-btn">${ICON_BACK}</button>
        </div>
        <button class="plus-badge banner-plus" id="btn-add-banner">${ICON_PLUS}</button>
      </div>
      <div class="avatar-container">
        <div class="avatar" id="btn-edit-avatar">${avatarContent}</div>
        <button class="plus-badge avatar-plus" id="btn-add-avatar">${ICON_PLUS}</button>
      </div>
      <div class="profile-info" style="padding-top: 80px;">
        <div class="name" style="color:${user.profile_color || '#c084fc'}">${escapeHtml(displayName)}</div>
        <div class="username">@${escapeHtml(user.game_username || 'unknown')}</div>
        <div class="bio">${escapeHtml(user.bio || '---')}</div>
      </div>
      <div class="main-stats">
        <div class="stat-box" id="like-box" style="cursor:pointer;">
          <div class="stat-value" id="heart-icon">${heartSvg}</div>
          <div class="stat-label" id="likes-count">${user.likes_count || 0}</div>
        </div>
        <div class="stat-box"><div class="stat-value">${user.age_days || 0}</div><div class="stat-label">روز</div></div>
        <div class="stat-box"><div class="stat-value">${user.wins || 0}</div><div class="stat-label">برد</div></div>
      </div>
      <div class="action-buttons">
        <button class="btn-square blue" id="btn-edit-bio">${ICON_PENCIL}</button>
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
  document.getElementById('btn-add-avatar')?.addEventListener('click', () => editModal?.classList.add('active'));
}
