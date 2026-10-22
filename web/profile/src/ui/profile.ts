import { User, GameStat, StoreItem } from '../types';

const ICON_GEAR = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path><circle cx="12" cy="12" r="3"></circle></svg>`;
const ICON_BACK = `<svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>`;
const ICON_PLUS = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>`;
const ICON_PENCIL = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#8b949e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>`;
const ICON_HEART = `<svg width="28" height="28" viewBox="0 0 24 24" fill="{fill}" stroke="{stroke}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`;
const ICON_VERIFIED = `<svg width="18" height="18" viewBox="0 0 24 24" fill="#3b82f6" style="margin-right: 6px;"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg>`;

const GAME_INFO: Record<string, { name: string, icon: string }> = {
  'ludo': { name: 'منچ', icon: '🎯' },
  'football': { name: 'فوتبال دستی', icon: '⚽' },
  'billiards': { name: 'بیلیارد', icon: '🎱' }
};

export function escapeHtml(s: string): string {
  return String(s || '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c] as string));
}

function renderProfileCard(user: User, isLiked: boolean, banners: StoreItem[]): string {
  const displayName = user.display_name || user.game_username || ('User' + user.id);
  const avatarContent = user.avatar ? `<img src="${escapeHtml(user.avatar)}">` : (user.name_emoji || '🎮');
  const heartFill = isLiked ? '#ef4444' : 'none';
  const heartStroke = isLiked ? '#ef4444' : '#ffffff';
  const heartSvg = ICON_HEART.replace('{fill}', heartFill).replace('{stroke}', heartStroke);
  const defaultBanner = banners.length > 0 ? banners[0].url : '#15181e';
  const bannerStyle = user.banner_id ? user.banner_id : defaultBanner;

  return `
    <div class="profile-card">
      <div class="top-bar">
        <button class="icon-btn" id="btn-open-settings">${ICON_GEAR}</button>
        <button class="icon-btn" onclick="window.location.hash='#home'">${ICON_BACK}</button>
      </div>
      <div class="banner" style="background: ${bannerStyle}; background-size: cover; background-position: center;">
        <button class="plus-badge banner-plus" id="btn-edit-banner">${ICON_PLUS}</button>
      </div>
      <div class="avatar-container">
        <div class="avatar" id="btn-edit-avatar">${avatarContent}</div>
        <button class="plus-badge avatar-plus" id="btn-change-avatar">${ICON_PLUS}</button>
      </div>
      <div class="profile-info">
        <div class="name" style="color:${user.profile_color || '#c084fc'}; display: flex; align-items: center; justify-content: center;">
            ${escapeHtml(displayName)} <span title="تایید شده">${ICON_VERIFIED}</span>
        </div>
        <div class="username">@${escapeHtml(user.game_username || 'unknown')}</div>
        <div class="bio-wrapper">
            <span class="bio">${escapeHtml(user.bio || '---')}</span>
            <button class="pencil-btn" id="btn-edit-bio">${ICON_PENCIL}</button>
        </div>
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
        <button class="btn-square-small green"><span class="icon-3d">💎</span> ${user.gems || 0}</button>
        <button class="btn-square-small"><span class="icon-3d">🪙</span> ${user.coins || 0}</button>
      </div>
    </div>`;
}

export function renderGameStats(stats: GameStat[]): string {
  if (!stats || stats.length === 0) return `<div class="game-stats-card"><div style="text-align:center;padding:20px;color:#8b949e;">هنوز آماری ثبت نشده است.</div></div>`;
  
  let html = `<div class="game-stats-card" style="padding: 0; background: transparent;">`;
  
  stats.forEach(g => {
    const info = GAME_INFO[g.game_id] || { name: g.game_id, icon: '🎮' };
    
    let rank = 'Beginner';
    if (g.level >= 80) rank = 'Prime Gamer';
    else if (g.level >= 60) rank = 'Master';
    else if (g.level >= 40) rank = 'Skilled';
    else if (g.level >= 20) rank = 'Captain';
    else if (g.level >= 10) rank = 'Adept';
    else if (g.level >= 5) rank = 'Gamer';
    
    const winRate = g.games_played > 0 ? Math.round((g.wins / g.games_played) * 100) : 0;
    
    html += `
      <div style="background: #1c2028; border-radius: 12px; padding: 15px; margin-bottom: 12px; display: flex; align-items: center; justify-content: space-between; box-shadow: 0 4px 10px rgba(0,0,0,0.2);">
        <div style="display: flex; align-items: center; gap: 15px; flex: 1;">
          <div style="position: relative; width: 60px; height: 60px; flex-shrink: 0;">
            <div style="width: 100%; height: 100%; background: linear-gradient(135deg, #2a2e38, #1c2028); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 35px; border: 1px solid rgba(255,255,255,0.1);">
              ${info.icon}
            </div>
            <div style="position: absolute; bottom: -8px; right: -8px; width: 32px; height: 32px; background: #3b82f6; clip-path: polygon(50% 0%, 61% 35%, 98% 35%, 68% 57%, 79% 91%, 50% 70%, 21% 91%, 32% 57%, 2% 35%, 39% 35%); display: flex; align-items: center; justify-content: center; color: white; font-weight: bold; font-size: 13px; filter: drop-shadow(0 2px 3px rgba(0,0,0,0.5));">
              ${g.level}
            </div>
          </div>
          <div style="flex: 1;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px;">
              <span style="font-size: 16px; font-weight: 800; color: #fff;">${info.name}</span>
              <span style="font-size: 13px; color: #8b949e; font-weight: 600;">${rank}</span>
            </div>
            <div style="height: 8px; background: #15181e; border-radius: 4px; overflow: hidden; margin-bottom: 8px; border: 1px solid rgba(255,255,255,0.05);">
              <div style="width: ${winRate}%; height: 100%; background: linear-gradient(90deg, #22c55e, #4ade80); border-radius: 4px; box-shadow: 0 0 8px rgba(74, 222, 128, 0.4);"></div>
            </div>
            <div style="display: flex; justify-content: space-between; font-size: 12px; color: #8b949e;">
              <span>${g.games_played.toLocaleString()} بازی کرده</span>
              <span>${g.wins.toLocaleString()} برنده شده</span>
            </div>
          </div>
        </div>
      </div>
    `;
  });
  return html + `</div>`;
}

export function renderApp(user: User, stats: GameStat[], banners: StoreItem[], appEl: HTMLElement, isLiked: boolean) {
  appEl.innerHTML = `<div class="container">${renderProfileCard(user, isLiked, banners)}<div class="section-title" style="margin-bottom: 10px; color:#8b949e;">آمار من</div>${renderGameStats(stats)}</div>`;
}
