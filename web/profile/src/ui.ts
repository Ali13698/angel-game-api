import { UserProfile } from './types';

const ICON_PENCIL = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#8b949e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>`;
const ICON_PLUS = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>`;
const ICON_HEART = (filled: boolean) => `<svg width="28" height="28" viewBox="0 0 24 24" fill="${filled ? '#ef4444' : 'none'}" stroke="${filled ? '#ef4444' : '#ffffff'}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`;

export function renderProfile(user: UserProfile, isLiked: boolean): string {
    const heart = ICON_HEART(isLiked);
    return `
    <div class="profile-card" style="margin-top: 20px;">
      <div class="banner" style="background: ${user.banner_id || 'linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)'};">
        <button class="plus-badge banner-plus" id="btn-edit-banner">${ICON_PLUS}</button>
      </div>
      
      <div class="avatar-container">
        <div class="avatar"><img src="${user.avatar || 'https://api.dicebear.com/7.x/bottts/svg?seed=1'}" /></div>
        <button class="plus-badge avatar-plus" id="btn-change-avatar">${ICON_PLUS}</button>
      </div>
      
      <div class="profile-info">
        <div class="name" style="color:${user.profile_color || '#c084fc'}">${user.display_name || user.game_username}</div>
        <div class="username">@${user.game_username}</div>
        <div class="bio-wrapper">
            <span class="bio">${user.bio || 'بیوگرافی...'}</span>
            <button class="pencil-btn" id="btn-edit-bio">${ICON_PENCIL}</button>
        </div>
      </div>
      
      <div class="main-stats">
        <div class="stat-box" id="like-box" style="cursor:pointer;">
          <div class="stat-value" id="heart-icon">${heart}</div>
          <div class="stat-label" id="likes-count">${user.likes_count || 0}</div>
        </div>
        <div class="stat-box"><div class="stat-value">${user.age_days || 1}</div><div class="stat-label">روز</div></div>
        <div class="stat-box"><div class="stat-value">${user.wins || 0}</div><div class="stat-label">برد</div></div>
      </div>
      
      <div class="action-buttons">
        <button class="btn-square-small green"><span class="icon-3d">💎</span> ${user.gems || 0}</button>
        <button class="btn-square-small"><span class="icon-3d">🪙</span> ${user.coins || 0}</button>
      </div>
    </div>
  `;
}

export function renderModals(user: UserProfile): string {
    let avatarGrid = '';
    for (let i = 1; i <= 20; i++) {
        const url = `https://api.dicebear.com/7.x/bottts/svg?seed=${i}`;
        avatarGrid += `<img src="${url}" class="avatar-option" data-url="${url}" />`;
    }

    const banners = [
        'linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)',
        'linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)',
        'linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)',
        'linear-gradient(135deg, #10b981 0%, #3b82f6 100%)',
        'linear-gradient(135deg, #1f2937 0%, #111827 100%)'
    ];
    let bannerGrid = banners.map(b => `<div class="banner-option" style="background: ${b};" data-bg="${b}"></div>`).join('');

    return `
    <!-- Modal Avatar -->
    <div id="avatar-modal" class="modal">
      <div class="modal-content">
        <h3 style="text-align:center; margin-bottom:15px;">انتخاب آواتار</h3>
        <div class="avatar-grid" id="avatar-grid">${avatarGrid}</div>
        <button class="modal-btn save" id="save-avatar-btn">تایید</button>
        <button class="modal-btn close" id="close-avatar-btn">انصراف</button>
      </div>
    </div>

    <!-- Modal Banner -->
    <div id="banner-modal" class="modal">
      <div class="modal-content">
        <h3 style="text-align:center; margin-bottom:15px;">انتخاب بنر</h3>
        <div class="banner-grid" id="banner-grid">${bannerGrid}</div>
        <button class="modal-btn close" id="close-banner-btn">انصراف</button>
      </div>
    </div>

    <!-- Modal Edit Profile (Username & Bio) -->
    <div id="edit-modal" class="modal">
      <div class="modal-content">
        <h3 style="text-align:center; margin-bottom:15px;">ویرایش اطلاعات</h3>
        
        <label style="font-size:12px; color:#8b949e; margin-bottom:5px; display:block;">نام کاربری (حداکثر ۱۰ حرف انگلیسی/عدد)</label>
        <input type="text" id="username-input" maxlength="10" value="${user.game_username || ''}" style="width:100%; background:var(--inner-card-bg); border:1px solid #333; border-radius:8px; color:white; padding:10px; margin-bottom:15px; font-family:inherit;">
        
        <label style="font-size:12px; color:#8b949e; margin-bottom:5px; display:block;">بیوگرافی (حداکثر ۱۰۰ کاراکتر - ایموجی مجاز است)</label>
        <textarea id="bio-input" maxlength="100" placeholder="بیوگرافی شما...">${user.bio || ''}</textarea>
        
        <button class="modal-btn save" id="save-profile-btn">ذخیره</button>
        <button class="modal-btn close" id="close-edit-btn">انصراف</button>
      </div>
    </div>
    `;
}
