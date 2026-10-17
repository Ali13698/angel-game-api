import { User, StoreItem } from '../types';
import { escapeHtml } from './profile';

const ICON_LOCK = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>`;
const ICON_GALLERY = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>`;
const ICON_GIF = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><text x="12" y="16" text-anchor="middle" font-size="10" fill="currentColor">GIF</text></svg>`;
const ICON_FRAME = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16v16H4z"></path><path d="M9 4v16"></path><path d="M15 4v16"></path><path d="M4 9h16"></path><path d="M4 15h16"></path></svg>`;
const ICON_AVATAR = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><circle cx="12" cy="10" r="3"></circle><path d="M7 20.662V19a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v1.662"></path></svg>`;

export function renderAvatarOptionsModal(): string {
  return `
    <div id="avatar-options-modal" class="modal">
      <div class="modal-content" style="padding: 15px; background: #1c2028;">
        <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:15px;">
            <h3 style="margin:0; font-size:16px;">تغییر تصویر</h3>
            <button id="close-avatar-options" style="background:none; border:none; color:white; font-size:20px; cursor:pointer;">✖</button>
        </div>
        
        <div class="avatar-options-grid">
            <div class="option-btn locked" id="opt-gif">
                <div class="opt-icon">${ICON_GIF}</div>
                <div class="opt-label">انتخاب گیف</div>
                <div class="lock-badge">${ICON_LOCK} 2</div>
            </div>
            <div class="option-btn active" id="opt-gallery">
                <div class="opt-icon">${ICON_GALLERY}</div>
                <div class="opt-label">گالری</div>
            </div>
            <div class="option-btn locked" id="opt-frame">
                <div class="opt-icon">${ICON_FRAME}</div>
                <div class="opt-label">قاب</div>
                <div class="lock-badge">${ICON_LOCK}</div>
            </div>
            <div class="option-btn active" id="opt-dicebear">
                <div class="opt-icon">${ICON_AVATAR}</div>
                <div class="opt-label">آواتار</div>
            </div>
        </div>
      </div>
    </div>`;
}

export function renderCropperModal(): string {
    return `
    <div id="cropper-modal" class="modal">
      <div class="modal-content" style="padding: 15px; background: #1c2028; width: 95%; max-width: 500px;">
        <h3 style="text-align:center; margin-bottom:15px;">تنظیم عکس</h3>
        <div style="width: 100%; height: 300px; background: #000; border-radius: 8px; overflow: hidden;">
            <img id="image-to-crop" style="max-width: 100%; display: block;" />
        </div>
        <div style="display:flex; gap:10px; margin-top:15px;">
            <button class="modal-btn save" id="save-crop-btn">تایید و آپلود</button>
            <button class="modal-btn close" id="close-crop-btn" style="width:40%;">انصراف</button>
        </div>
      </div>
    </div>`;
}

export function renderDiceBearModal(user: User, avatars: StoreItem[]): string {
  let gridHtml = '';
  avatars.forEach(item => { gridHtml += `<img src="${item.url}" class="avatar-option ${user.avatar === item.url ? 'selected' : ''}" data-url="${item.url}" />`; });
  return `<div id="dicebear-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">انتخاب آواتار</h3><div class="avatar-grid" id="dicebear-grid">${gridHtml}</div><button class="modal-btn save" id="save-dicebear-btn">تایید</button><button class="modal-btn close" id="close-dicebear-btn">انصراف</button></div></div>`;
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