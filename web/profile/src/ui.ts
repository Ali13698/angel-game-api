import { User } from './types';

const ICON_HEART = (filled: boolean) => `<svg width="28" height="28" viewBox="0 0 24 24" fill="${filled ? '#ef4444' : 'none'}" stroke="${filled ? '#ef4444' : '#ffffff'}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`;

export function renderProfile(user: User, isLiked: boolean): string {
    let avatarGrid = '';
    for (let i = 1; i <= 20; i++) {
        avatarGrid += `<img src="https://api.dicebear.com/7.x/bottts/svg?seed=${i}" class="avatar-option" data-url="https://api.dicebear.com/7.x/bottts/svg?seed=${i}" />`;
    }

    const banners = [
        'linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)',
        'linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)',
        'linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)'
    ];
    let bannerGrid = banners.map(b => `<div class="banner-option" style="background: ${b};" data-bg="${b}"></div>`).join('');

    return `
    <div class="app-shell">
        <div style="background:#22c55e; color:white; padding:8px; text-align:center; font-size:13px; font-weight:bold;">✅ نسخه TypeScript فعال شد!</div>
        <div class="page-content">
            <div class="profile-card">
                <div class="banner" id="profile-banner" style="background: ${user.banner_id || banners[0]}; background-size: cover;">
                    <button class="plus-badge banner-plus" id="btn-edit-banner">+</button>
                </div>
                <div class="avatar-container">
                    <div class="avatar"><img id="profile-avatar-img" src="${user.avatar || 'https://api.dicebear.com/7.x/bottts/svg?seed=1'}" /></div>
                    <button class="plus-badge avatar-plus" id="btn-change-avatar">+</button>
                </div>
                <div class="profile-info">
                    <div class="name" style="color:${user.profile_color || '#c084fc'}">${user.display_name || user.game_username}</div>
                    <div class="username">@${user.game_username}</div>
                    <div class="bio-wrapper">
                        <span class="bio" id="profile-bio">${user.bio || 'بیوگرافی...'}</span>
                        <button class="pencil-btn" id="btn-edit-bio">✏️</button>
                    </div>
                </div>
                <div class="main-stats">
                    <div class="stat-box" id="like-box" style="cursor:pointer;">
                        <div class="stat-value" id="heart-icon">${ICON_HEART(isLiked)}</div>
                        <div class="stat-label" id="likes-count">${user.likes_count || 0}</div>
                    </div>
                    <div class="stat-box"><div class="stat-value">${user.age_days || 1}</div><div class="stat-label">روز</div></div>
                    <div class="stat-box"><div class="stat-value">${user.wins || 0}</div><div class="stat-label">برد</div></div>
                </div>
            </div>
        </div>
    </div>

    <!-- Modals -->
    <div id="avatar-modal" class="modal"><div class="modal-content"><h3>انتخاب آواتار</h3><div class="avatar-grid" id="avatar-grid">${avatarGrid}</div><button class="modal-btn close" id="close-avatar-btn">انصراف</button></div></div>
    <div id="banner-modal" class="modal"><div class="modal-content"><h3>انتخاب بنر</h3><div class="banner-grid" id="banner-grid">${bannerGrid}</div><button class="modal-btn close" id="close-banner-btn">انصراف</button></div></div>
    `;
}
