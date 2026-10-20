import { User, GameStat } from '../types';
import { renderApp } from '../ui/profile';
import { renderAvatarOptionsModal, renderDiceBearModal, renderEditModal, renderBannerModal, renderSettingsModal, renderCropperModal } from '../ui/modals';
import { DEFAULT_BANNERS, DEFAULT_AVATARS } from '../store';
import { showToast } from '../ui/toast';
import { updateProfile, uploadImage } from '../api';
import { logError } from '../logger';

let cropperInstance: any = null;
let currentUser: User | null = null;
let currentStats: GameStat[] = [];
let appElement: HTMLElement | null = null;

export function initProfilePage(user: User, stats: GameStat[], appEl: HTMLElement) {
    currentUser = user;
    currentStats = stats;
    appElement = appEl;
    buildUI();
}

function buildUI() {
    if (!currentUser || !appElement) return;
    try {
        const isLiked = localStorage.getItem('liked_' + currentUser.id) === 'true';
        renderApp(currentUser, currentStats, DEFAULT_BANNERS, appElement, isLiked);
    } catch (e) {
        logError('Main UI render failed', e);
        appElement.innerHTML = '<div style="color:red;padding:20px;">خطا در نمایش پروفایل</div>';
        return;
    }

    try {
        document.querySelectorAll('.modal').forEach(m => m.remove());
        const modalsWrapper = document.createElement('div');
        modalsWrapper.innerHTML = renderAvatarOptionsModal() + 
                                  renderDiceBearModal(currentUser, DEFAULT_AVATARS) + 
                                  renderEditModal(currentUser) + 
                                  renderBannerModal(DEFAULT_BANNERS) + 
                                  renderSettingsModal() +
                                  renderCropperModal();
        document.body.appendChild(modalsWrapper);
    } catch (e) {
        logError('Modals render failed', e);
    }

    try {
        attachEvents();
    } catch (e) {
        logError('Event attachment failed', e);
    }
}

function attachEvents() {
    if (!currentUser) return;

    const safeListen = (id: string, event: string, handler: () => void) => {
        const el = document.getElementById(id);
        if (el) el.addEventListener(event, handler);
    };

    const editModal = document.getElementById('edit-modal');
    safeListen('btn-edit-bio', 'click', () => editModal?.classList.add('active'));
    safeListen('close-edit-btn', 'click', () => editModal?.classList.remove('active'));
    safeListen('save-profile-btn', 'click', () => {
      const uname = (document.getElementById('username-input') as HTMLInputElement)?.value || '';
      const bio = (document.getElementById('bio-input') as HTMLTextAreaElement)?.value || '';
      if (!/^[a-zA-Z0-9]{1,10}$/.test(uname)) { showToast('نام کاربری فقط حروف و اعداد (حداکثر ۱۰ کاراکتر)'); return; }
      handleUpdate({ game_username: uname, bio: bio });
      editModal?.classList.remove('active');
    });

    const avatarOptsModal = document.getElementById('avatar-options-modal');
    const dicebearModal = document.getElementById('dicebear-modal');
    const cropperModal = document.getElementById('cropper-modal');
    
    safeListen('btn-change-avatar', 'click', () => avatarOptsModal?.classList.add('active'));
    safeListen('close-avatar-options', 'click', () => avatarOptsModal?.classList.remove('active'));
    safeListen('opt-gif', 'click', () => showToast('خرید گیف به زودی!', 'error'));
    safeListen('opt-frame', 'click', () => showToast('خرید قاب به زودی!', 'error'));
    safeListen('opt-dicebear', 'click', () => {
        avatarOptsModal?.classList.remove('active');
        dicebearModal?.classList.add('active');
    });

    let selectedAvatar = currentUser.avatar || '';
    safeListen('close-dicebear-btn', 'click', () => dicebearModal?.classList.remove('active'));
    document.querySelectorAll('#dicebear-grid .avatar-option').forEach(img => img.addEventListener('click', (e) => { 
        document.querySelectorAll('#dicebear-grid .avatar-option').forEach(i => i.classList.remove('selected')); 
        (e.target as HTMLElement).classList.add('selected'); 
        selectedAvatar = (e.target as HTMLElement).getAttribute('data-url') || ''; 
    }));
    safeListen('save-dicebear-btn', 'click', () => { 
        handleUpdate({ avatar: selectedAvatar }); 
        dicebearModal?.classList.remove('active'); 
    });

    safeListen('opt-gallery', 'click', () => {
        avatarOptsModal?.classList.remove('active');
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = 'image/*';
        input.onchange = (e: any) => {
            const file = e.target.files?.[0];
            if (!file) return;
            const reader = new FileReader();
            reader.onload = (ev) => {
                const img = document.getElementById('image-to-crop') as HTMLImageElement;
                img.src = ev.target?.result as string;
                cropperModal?.classList.add('active');
                if (cropperInstance) cropperInstance.destroy();
                cropperInstance = new (window as any).Cropper(img, {
                    aspectRatio: 1, viewMode: 1, dragMode: 'move', autoCropArea: 1, background: false, cropBoxResizable: false, cropBoxMovable: false, toggleDragModeOnDblclick: false
                });
            };
            reader.readAsDataURL(file);
        };
        input.click();
    });

    safeListen('close-crop-btn', 'click', () => {
        cropperModal?.classList.remove('active');
        if (cropperInstance) cropperInstance.destroy();
    });

    safeListen('save-crop-btn', 'click', () => {
        if (!cropperInstance) return;
        cropperInstance.getCroppedCanvas({ width: 256, height: 256 }).toBlob(async (blob: Blob | null) => {
            if (!blob) return;
            cropperModal?.classList.remove('active');
            showToast('در حال آپلود...', 'success');
            const res = await uploadImage(blob);
            if (res.ok && res.url) { await handleUpdate({ avatar: res.url }); showToast('آواتار با موفقیت تغییر کرد', 'success'); } 
            else { showToast('خطا در آپلود: ' + (res.error || 'ناشناخته')); }
        }, 'image/jpeg', 0.9);
    });

    const bannerModal = document.getElementById('banner-modal');
    safeListen('btn-edit-banner', 'click', () => bannerModal?.classList.add('active'));
    safeListen('close-banner-btn', 'click', () => bannerModal?.classList.remove('active'));
    document.querySelectorAll('.banner-option').forEach(bg => bg.addEventListener('click', (e) => {
        const bgStyle = (e.target as HTMLElement).getAttribute('data-bg');
        if (bgStyle) { handleUpdate({ banner_id: bgStyle }); bannerModal?.classList.remove('active'); }
    }));

    safeListen('like-box', 'click', () => {
        if (!currentUser) return;
        const isLiked = localStorage.getItem('liked_' + currentUser.id) === 'true';
        const newIsLiked = !isLiked;
        localStorage.setItem('liked_' + currentUser.id, String(newIsLiked));
        if (newIsLiked) {
            currentUser.likes_count = (currentUser.likes_count || 0) + 1;
        } else {
            currentUser.likes_count = Math.max(0, (currentUser.likes_count || 0) - 1);
        }
        const heartIcon = document.getElementById('heart-icon');
        const likesCount = document.getElementById('likes-count');
        if (heartIcon) {
            const path = heartIcon.querySelector('path');
            if (path) {
                path.setAttribute('fill', newIsLiked ? '#ef4444' : 'none');
                path.setAttribute('stroke', newIsLiked ? '#ef4444' : '#ffffff');
            }
        }
        if (likesCount) {
            likesCount.innerText = currentUser.likes_count.toString();
        }
    });

    const settingsModal = document.getElementById('settings-modal');
    safeListen('btn-open-settings', 'click', () => settingsModal?.classList.add('active'));
    safeListen('close-settings-btn', 'click', () => settingsModal?.classList.remove('active'));

    safeListen('btn-back-home', 'click', () => {
        window.location.hash = '#home';
    });
}

async function handleUpdate(patch: any) {
    if (!currentUser) return;
    const res = await updateProfile(patch);
    if (res.ok) { 
        Object.assign(currentUser, patch); 
        buildUI(); 
    } 
    else { showToast(res.error || 'خطا در ذخیره اطلاعات'); }
}
