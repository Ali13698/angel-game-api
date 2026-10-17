import { User, GameStat } from './types';

const PRESET_EMOJIS = ['🐉','🦊','🌟','🐼','🦁','🐸','🐺','🦄','🐯','🐨','👾','🤖','🧑','👩','🧔','👑'];

// داده‌های تستی برای بازی‌ها (شبیه اسکرین‌شات)
const MOCK_GAMES: GameStat[] = [
  { id: 'tank', name: 'تانک', icon: '🚀', level: 25, rank: 'Captain', played: 384, won: 193, progress: 80 },
  { id: 'chess', name: 'شطرنج', icon: '♟️', level: 16, rank: 'Skilled', played: 196, won: 88, progress: 45 },
  { id: 'backgammon', name: 'بک‌گمون', icon: '🎲', level: 15, rank: 'Skilled', played: 142, won: 60, progress: 90 },
  { id: 'ludo', name: 'منچ', icon: '🎯', level: 7, rank: 'Gamer', played: 47, won: 10, progress: 30 },
  { id: 'basketball', name: 'بسکتبال', icon: '🏀', level: 6, rank: 'Beginner', played: 16, won: 9, progress: 50 },
];

export function escapeHtml(s: string): string {
  return String(s || '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c] as string));
}

export function render(user: User, appEl: HTMLElement, onBioClick: () => void, onEmojiClick: (e: string) => void) {
  const displayName = user.display_name || user.game_username || ('User' + user.id);
  const color = user.profile_color || '#c084fc';
  const avatarContent = user.avatar ? `<img src="${escapeHtml(user.avatar)}">` : (user.name_emoji || '🎮');

  // بخش پروفایل
  let html = `
  <div class="container">
    <div class="profile-card">
      <div class="banner"></div>
      <div class="avatar-container">
        <div class="avatar">${avatarContent}</div>
      </div>
      <div class="profile-info">
        <div class="name" style="color:${color}">${escapeHtml(displayName)}</div>
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
      
      <div style="padding: 0 20px 20px;">
         <button class="btn-hype">هایپ 🚀</button>
      </div>
    </div>

    <div class="section-title" style="margin-bottom: 10px; color:#8b949e;">آمار من</div>
    <div class="game-stats-card">
  `;

  // بخش آمار بازی‌ها
  MOCK_GAMES.forEach(g => {
    html += `
      <div class="game-item">
        <div class="game-icon">${g.icon}</div>
        <div class="game-details">
          <div class="game-header">
            <div><span class="game-name">${g.name}</span> <span class="game-rank">${g.rank}</span></div>
            <div class="game-level">مرحله ${g.level}</div>
          </div>
          <div class="progress-bar"><div class="progress-fill" style="width: ${g.progress}%"></div></div>
          <div class="game-meta">
            <span>${g.played} بازی کرده</span>
            <span>${g.won} برنده شده</span>
          </div>
        </div>
        <div class="game-star">${g.level}</div>
      </div>
    `;
  });

  html += `</div></div>`;
  appEl.innerHTML = html;

  // رویدادها
  document.getElementById('btn-edit-bio')?.addEventListener('click', onBioClick);
}
