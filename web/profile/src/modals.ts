import { User } from './types';

export function renderAvatarModal(user: User): string {
    let grid = '';
    for (let i = 1; i <= 20; i++) {
        const url = `https://api.dicebear.com/7.x/bottts/svg?seed=${i}`;
        grid += `<img src="${url}" class="avatar-option" data-url="${url}" />`;
    }
    return `<div id="avatar-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">انتخاب آواتار</h3><div class="avatar-grid" id="avatar-grid">${grid}</div><button class="modal-btn save" id="save-avatar-btn">تایید</button><button class="modal-btn close" id="close-avatar-btn">انصراف</button></div></div>`;
}

export function renderBannerModal(): string {
    const banners = [
        'linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)',
        'linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)',
        'linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)',
        'linear-gradient(135deg, #10b981 0%, #3b82f6 100%)',
        'linear-gradient(135deg, #1f2937 0%, #111827 100%)'
    ];
    let grid = banners.map(b => `<div class="banner-option" style="background: ${b};" data-bg="${b}"></div>`).join('');
    return `<div id="banner-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">انتخاب بنر</h3><div class="banner-grid" id="banner-grid">${grid}</div><button class="modal-btn close" id="close-banner-btn">انصراف</button></div></div>`;
}

export function renderEditModal(user: User): string {
    return `
    <div id="edit-modal" class="modal">
      <div class="modal-content">
        <h3 style="text-align:center; margin-bottom:15px;">ویرایش اطلاعات</h3>
        <label style="font-size:12px; color:#8b949e; margin-bottom:5px; display:block;">نام کاربری (فقط انگلیسی و عدد - حداکثر ۱۰)</label>
        <input type="text" id="username-input" maxlength="10" value="${user.game_username || ''}" style="width:100%; background:var(--inner-card-bg); border:1px solid #333; border-radius:8px; color:white; padding:10px; margin-bottom:15px;">
        <label style="font-size:12px; color:#8b949e; margin-bottom:5px; display:block;">بیوگرافی (حداکثر ۱۰۰ کاراکتر - ایموجی مجاز)</label>
        <textarea id="bio-input" maxlength="100" style="width:100%; height:80px; background:var(--inner-card-bg); border:1px solid #333; border-radius:8px; color:white; padding:10px; margin-bottom:15px;">${user.bio || ''}</textarea>
        <button class="modal-btn save" id="save-profile-btn">ذخیره</button>
        <button class="modal-btn close" id="close-edit-btn">انصراف</button>
      </div>
    </div>`;
}
