import { User, GameStat } from './types';
const GAME_INFO: Record<string, { name: string, icon: string }> = { 'ludo': { name: 'منچ', icon: '🎯' }, 'backgammon': { name: 'بک‌گمون', icon: '🎲' }, 'chess': { name: 'شطرنج', icon: '♟️' }, 'football': { name: 'فوتبال دستی', icon: '⚽' }, 'billiards': { name: 'بیلیارد', icon: '🎱' } };
const DICEBEAR_STYLES = ['adventurer', 'avataaars', 'bottts', 'fun-emoji', 'lorelei', 'micah', 'miniavs', 'notionists', 'open-peeps', 'personas'];
export const AVATARS = Array.from({ length: 10 }, (_, i) => `https://api.dicebear.com/7.x/${DICEBEAR_STYLES[i % 10]}/svg?seed=Angel${i + 1}`);
export function escapeHtml(s: string): string { return String(s || '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c] as string)); }
function renderProfileCard(user: User): string {
  const displayName = user.display_name || user.game_username || ('User' + user.id);
  const avatarContent = user.avatar ? `<img src="${escapeHtml(user.avatar)}">` : (user.name_emoji || '🎮');
  return `<div class="profile-card"><div class="banner"></div><div class="avatar-container"><div class="avatar" id="btn-edit-avatar">${avatarContent}</div></div><div class="profile-info"><div class="name" style="color:${user.profile_color || '#c084fc'}">${escapeHtml(displayName)}</div><div class="username">@${escapeHtml(user.game_username || 'unknown')}</div><div class="bio">${escapeHtml(user.bio || '---')}</div></div><div class="main-stats"><div class="stat-box"><div class="stat-value">${user.likes_count || 0}</div><div class="stat-label">هاله</div></div><div class="stat-box"><div class="stat-value">${user.age_days || 0}</div><div class="stat-label">روز</div></div><div class="stat-box"><div class="stat-value">${user.wins || 0}</div><div class="stat-label">برد</div></div></div><div class="action-buttons"><button class="btn-square blue" id="btn-edit-bio">✏️</button><button class="btn-square green">💎 ${user.gems || 0}</button><button class="btn-square">🪙 ${user.coins || 0}</button></div><div style="padding: 0 20px 20px;"><button class="btn-hype">هایپ 🚀</button></div></div>`;
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
export function render(user: User, stats: GameStat[], appEl: HTMLElement, onSaveProfile: (avatar: string, bio: string) => void) {
  appEl.innerHTML = `<div class="container">${renderProfileCard(user)}<div class="section-title" style="margin-bottom: 10px; color:#8b949e;">آمار من</div>${renderGameStats(stats)}</div>${renderEditModal(user)}`;
  let selectedAvatar = user.avatar || '';
  const modal = document.getElementById('edit-modal');
  [document.getElementById('btn-edit-avatar'), document.getElementById('btn-edit-bio')].forEach(btn => btn?.addEventListener('click', () => modal?.classList.add('active')));
  document.getElementById('close-modal-btn')?.addEventListener('click', () => modal?.classList.remove('active'));
  document.querySelectorAll('.avatar-option').forEach(img => img.addEventListener('click', (e) => { document.querySelectorAll('.avatar-option').forEach(i => i.classList.remove('selected')); (e.target as HTMLElement).classList.add('selected'); selectedAvatar = (e.target as HTMLElement).getAttribute('data-url') || ''; }));
  document.getElementById('save-profile-btn')?.addEventListener('click', () => onSaveProfile(selectedAvatar, (document.getElementById('bio-input') as HTMLTextAreaElement).value));
}
