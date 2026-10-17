import { User, GameStat } from './types';
import { StoreItem } from './store';

const ICON_GEAR = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path><circle cx="12" cy="12" r="3"></circle></svg>`;
const ICON_BACK = `<svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>`;
const ICON_PLUS = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>`;
const ICON_PENCIL = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#8b949e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>`;
const ICON_HEART = `<svg width="28" height="28" viewBox="0 0 24 24" fill="{fill}" stroke="{stroke}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`;

const GAME_INFO: Record<string, { name: string, icon: string }> = { 'ludo': { name: 'منچ', icon: '🎯' }, 'backgammon': { name: 'بک‌گمون', icon: '🎲' }, 'chess': { name: 'شطرنج', icon: '♟️' }, 'football': { name: 'فوتبال دستی', icon: '⚽' }, 'billiards': { name: 'بیلیارد', icon: '🎱' } };

export function escapeHtml(s: string): string { return String(s || '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c] as string)); }

function renderProfileCard(user: User, isLiked: boolean, banners: StoreItem[]): string {
  const displayName = user.display_name || user.game_username || ('User' + user.id);
  const avatarContent = user.avatar ? `<img src="${escapeHtml(user.avatar)}">` : (user.name_emoji || '🎮');
  const heartFill = isLiked ? '#ef4444' : 'none';
  const heartStroke = isLiked ? '#ef4444' : '#ffffff';
  const heartSvg = ICON_HEART.replace('{fill}', heartFill).replace('{stroke}', heartStroke);
  
  // Find the banner style based on banner_id
  const defaultBanner = banners.length > 0 ? banners[0].url : '#15181e';
  const bannerStyle = user.banner_id ? user.banner_id : defaultBanner;

  return `
    <div class="profile-card">
      <div class="top-bar">
        <button class="icon-btn" id="btn-open-settings">${ICON_GEAR}</button>
        <button class="icon-btn">${ICON_BACK}</button>
      </div>
      
      <div class="banner" style="background: ${bannerStyle}; background-size: cover; background-position: center;">
        <button class="plus-badge banner-plus" id="btn-edit-banner">${ICON_PLUS}</button>
        <input type="file" id="banner-file-input" accept="image/*" style="display:none;">
      </div>
      
      <div class="avatar-container">
        <div class="avatar" id="btn-edit-avatar">${avatarContent}</div>
        <button class="plus-badge avatar-plus" id="btn-change-avatar">${ICON_PLUS}</button>
        <input type="file" id="avatar-file-input" accept="image/*" style="display:none;">
      </div>
      
      <div class="profile-info">
        <div class="name" style="color:${user.profile_color || '#c084fc'}">${escapeHtml(displayName)}</div>
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

function renderAvatarModal(user: User, avatars: StoreItem[]): string {
  let gridHtml = '';
  avatars.forEach(item => { gridHtml += `<img src="${item.url}" class="avatar-option ${user.avatar === item.url ? 'selected' : ''}" data-url="${item.url}" />`; });
  return `<div id="avatar-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">انتخاب آواتار</h3><div class="avatar-grid" id="avatar-grid">${gridHtml}</div><button class="modal-btn save" id="save-avatar-btn">تایید</button><button class="modal-btn close" id="close-avatar-btn">انصراف</button></div></div>`;
}

function renderBioModal(user: User): string {
  return `<div id="bio-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">ویرایش بیوگرافی</h3><textarea id="bio-input" placeholder="بیوگرافی جدید...">${escapeHtml(user.bio || '')}</textarea><button class="modal-btn save" id="save-bio-btn">ذخیره</button><button class="modal-btn close" id="close-bio-btn">انصراف</button></div></div>`;
}

function renderBannerModal(banners: StoreItem[]): string {
  let gridHtml = '';
  banners.forEach(item => {
    gridHtml += `<div class="banner-option" style="background: ${item.url}; background-size: cover; background-position: center;" data-bg="${item.url}"></div>`;
  });
  return `<div id="banner-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">انتخاب بنر</h3><div class="banner-grid" id="banner-grid">${gridHtml}</div><button class="modal-btn close" id="close-banner-btn">انصراف</button></div></div>`;
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

export function render(
  user: User, 
  stats: GameStat[], 
  banners: StoreItem[], 
  avatars: StoreItem[], 
  appEl: HTMLElement, 
  onUpdate: (patch: any) => void, 
  onUpload: (file: File, type: 'avatar' | 'banner') => void,
  onLike: () => void
) {
  const isLiked = localStorage.getItem('liked_' + user.id) === 'true';
  appEl.innerHTML = `
    <div class="container">
      ${renderProfileCard(user, isLiked, banners)}
      <div class="section-title" style="margin-bottom: 10px; color:#8b949e;">آمار من</div>
      ${renderGameStats(stats)}
    </div>
    ${renderAvatarModal(user, avatars)}
    ${renderBioModal(user)}
    ${renderBannerModal(banners)}
    ${renderSettingsModal()}
  `;
  
  // Avatar Logic (Select from store or Upload)
  let selectedAvatar = user.avatar || '';
  const avatarModal = document.getElementById('avatar-modal');
  document.getElementById('btn-change-avatar')?.addEventListener('click', () => avatarModal?.classList.add('active'));
  document.getElementById('close-avatar-btn')?.addEventListener('click', () => avatarModal?.classList.remove('active'));
  document.querySelectorAll('.avatar-option').forEach(img => img.addEventListener('click', (e) => { document.querySelectorAll('.avatar-option').forEach(i => i.classList.remove('selected')); (e.target as HTMLElement).classList.add('selected'); selectedAvatar = (e.target as HTMLElement).getAttribute('data-url') || ''; }));
  document.getElementById('save-avatar-btn')?.addEventListener('click', () => { onUpdate({ avatar: selectedAvatar }); avatarModal?.classList.remove('active'); });

  // Bio Logic
  const bioModal = document.getElementById('bio-modal');
  document.getElementById('btn-edit-bio')?.addEventListener('click', () => bioModal?.classList.add('active'));
  document.getElementById('close-bio-btn')?.addEventListener('click', () => bioModal?.classList.remove('active'));
  document.getElementById('save-bio-btn')?.addEventListener('click', () => { onUpdate({ bio: (document.getElementById('bio-input') as HTMLTextAreaElement).value }); bioModal?.classList.remove('active'); });

  // Banner Logic (Select from store or Upload)
  const bannerModal = document.getElementById('banner-modal');
  document.getElementById('btn-edit-banner')?.addEventListener('click', () => bannerModal?.classList.add('active'));
  document.getElementById('close-banner-btn')?.addEventListener('click', () => bannerModal?.classList.remove('active'));
  document.querySelectorAll('.banner-option').forEach(bg => bg.addEventListener('click', (e) => {
      const bgStyle = (e.target as HTMLElement).getAttribute('data-bg');
      if (bgStyle) { onUpdate({ banner_id: bgStyle }); bannerModal?.classList.remove('active'); }
  }));

  // Upload Logic (Triggered by hidden file inputs)
  const avatarInput = document.getElementById('avatar-file-input') as HTMLInputElement;
  const bannerInput = document.getElementById('banner-file-input') as HTMLInputElement;

  document.getElementById('btn-change-avatar')?.addEventListener('contextmenu', (e) => { e.preventDefault(); avatarInput.click(); });
  // Double click or special button can trigger upload, for now let's use long press or another UI element. Let's just use a dedicated upload button inside the modal.
  // For simplicity, let's trigger upload if they click a specific button in the modal (to be added) or just use the plus badge for uploading directly, and clicking the avatar for the store.
  
  // Let's make the Plus button open the modal, but add an "Upload" button inside the modal.
  // Actually, the prompt says: "+ فقط برای تغییر آواتار آماده یا انتخاب عکس از گالری گوشی کاربر"
  // Let's add an upload button in the modal itself.
  
  const avatarGrid = document.getElementById('avatar-grid');
  if (avatarGrid) {
      const uploadBtn = document.createElement('div');
      uploadBtn.className = 'avatar-option';
      uploadBtn.style.cssText = 'display:flex;align-items:center;justify-content:center;font-size:24px;background:#2a2e38;color:white;cursor:pointer;';
      uploadBtn.innerHTML = '+';
      uploadBtn.onclick = () => avatarInput.click();
      avatarGrid.appendChild(uploadBtn);
  }

  const bannerGrid = document.getElementById('banner-grid');
  if (bannerGrid) {
      const uploadBtn = document.createElement('div');
      uploadBtn.className = 'banner-option';
      uploadBtn.style.cssText = 'display:flex;align-items:center;justify-content:center;font-size:24px;background:#2a2e38;color:white;cursor:pointer;';
      uploadBtn.innerHTML = '+';
      uploadBtn.onclick = () => bannerInput.click();
      bannerGrid.appendChild(uploadBtn);
  }

  avatarInput?.addEventListener('change', (e) => {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (file) { onUpload(file, 'avatar'); avatarModal?.classList.remove('active'); }
  });

  bannerInput?.addEventListener('change', (e) => {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (file) { onUpload(file, 'banner'); bannerModal?.classList.remove('active'); }
  });

  // Like & Settings
  document.getElementById('like-box')?.addEventListener('click', onLike);
  const settingsModal = document.getElementById('settings-modal');
  document.getElementById('btn-open-settings')?.addEventListener('click', () => settingsModal?.classList.add('active'));
  document.getElementById('close-settings-btn')?.addEventListener('click', () => settingsModal?.classList.remove('active'));
}
