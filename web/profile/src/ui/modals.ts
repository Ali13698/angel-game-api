import { User, StoreItem } from '../types';
import { escapeHtml } from './profile';

const ICON_PLUS = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>`;

export function renderAvatarModal(user: User, avatars: StoreItem[]): string {
  let gridHtml = '';
  avatars.forEach(item => { gridHtml += `<img src="${item.url}" class="avatar-option ${user.avatar === item.url ? 'selected' : ''}" data-url="${item.url}" />`; });
  gridHtml += `<div class="avatar-option" id="btn-upload-avatar" style="display:flex;align-items:center;justify-content:center;font-size:24px;background:#2a2e38;color:white;cursor:pointer;">${ICON_PLUS}</div>`;
  return `<div id="avatar-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">انتخاب آواتار</h3><div class="avatar-grid" id="avatar-grid">${gridHtml}</div><input type="file" id="avatar-file-input" accept="image/*" style="display:none;"><button class="modal-btn save" id="save-avatar-btn">تایید</button><button class="modal-btn close" id="close-avatar-btn">انصراف</button></div></div>`;
}

export function renderEditModal(user: User): string {
  return `
    <div id="edit-modal" class="modal">
      <div class="modal-content">
        <h3 style="text-align:center; margin-bottom:15px;">ویرایش اطلاعات</h3>
        <label style="font-size:12px; color:#8b949e; margin-bottom:5px; display:block;">نام کاربری (فقط حروف و اعداد - حداکثر ۱۰ کاراکتر)</label>
        <input type="text" id="username-input" maxlength="10" pattern="[a-zA-Z0-9]+" value="${escapeHtml(user.game_username || '')}" style="width:100%; background:var(--inner-card-bg); border:1px solid #333; border-radius:8px; color:white; padding:10px; margin-bottom:15px; font-family:inherit;">
        <label style="font-size:12px; color:#8b949e; margin-bottom:5px; display:block;">بیوگرافی (حداکثر ۱۰۰ کاراکتر)</label>
        <textarea id="bio-input" maxlength="100" placeholder="بیوگرافی جدید...">${escapeHtml(user.bio || '')}</textarea>
        <button class="modal-btn save" id="save-profile-btn">ذخیره</button>
        <button class="modal-btn close" id="close-edit-btn">انصراف</button>
      </div>
    </div>`;
}

export function renderBannerModal(banners: StoreItem[]): string {
  let gridHtml = '';
  banners.forEach(item => {
    if (item.owned) {
        gridHtml += `<div class="banner-option" style="background: ${item.url}; background-size: cover; background-position: center;" data-bg="${item.url}"></div>`;
    }
  });
  return `<div id="banner-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">انتخاب بنر</h3><div class="banner-grid" id="banner-grid">${gridHtml}</div><button class="modal-btn close" id="close-banner-btn">انصراف</button></div></div>`;
}

export function renderSettingsModal(): string {
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