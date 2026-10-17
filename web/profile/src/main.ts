import { User, GameStat } from './types';
import { updateProfile, uploadImage } from './api';
import { renderApp } from './ui/profile';
import { renderAvatarOptionsModal, renderDiceBearModal, renderEditModal, renderBannerModal, renderSettingsModal, renderCropperModal } from './ui/modals';
import { showToast } from './ui/toast';
import { DEFAULT_BANNERS, DEFAULT_AVATARS } from './store';
import { logInfo, logError } from './logger';

const appEl = document.getElementById('app');
let currentUser: User | null = null;
let currentStats: GameStat[] = [];
let cropperInstance: any = null; // For Cropper.js

window.addEventListener('unhandledrejection', (event) => {
    logError('Unhandled Promise Rejection', event.reason);
    showToast('خطای غیرمنتظره رخ داد.');
});

async function loadProfile() {
  const p = new URLSearchParams(location.search).get('id');
  if (!p) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">No ID provided</div>'; return; }
  logInfo('Starting profile load for ID:', p);
  try {
    const r = await fetch('/api/profile/' + p);
    const d = await r.json();
    if (!d.ok) throw new Error(d.error || 'User not found');
    currentUser = d.user; currentStats = d.stats || [];
    logInfo('Profile loaded successfully', currentUser);
    if (appEl && currentUser) {
      buildUI();
    }
  } catch(e: any) { 
    logError('Load profile failed', e);
    if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;color:red;direction:rtl;">Error: ' + e.message + '</div>'; 
  }
}

function buildUI() {
  if (!currentUser || !appEl) return;
  logInfo('Building UI...');
  const isLiked = localStorage.getItem('liked_' + currentUser.id) === 'true';
  renderApp(currentUser, currentStats, DEFAULT_BANNERS, appEl, isLiked);

  // Remove old modals if exist
  document.querySelectorAll('.modal').forEach(m => m.remove());

  const modalsWrapper = document.createElement('div');
  modalsWrapper.innerHTML = renderAvatarOptionsModal() + 
                            renderDiceBearModal(currentUser, DEFAULT_AVATARS) + 
                            renderEditModal(currentUser) + 
                            renderBannerModal(DEFAULT_BANNERS) + 
                            renderSettingsModal() +
                            renderCropperModal();
  document.body.appendChild(modalsWrapper);
  attachEvents();
}

function attachEvents() {
  if (!currentUser) return;
  logInfo('Attaching UI events...');

  // --- Profile Edit Modal ---
  const editModal = document.getElementById('edit-modal');
  document.getElementById('btn-edit-bio')?.addEventListener('click', () => editModal?.classList.add('active'));
  document.getElementById('close-edit-btn')?.addEventListener('click', () => editModal?.classList.remove('active'));
  document.getElementById('save-profile-btn')?.addEventListener('click', () => {
    const uname = (document.getElementById('username-input') as HTMLInputElement).value;
    const bio = (document.getElementById('bio-input') as HTMLTextAreaElement).value;
    if (!/^[a-zA-Z0-9]{1,10}$/.test(uname)) { showToast('نام کاربری فقط حروف و اعداد (حداکثر ۱۰ کاراکتر) باشد'); return; }
    handleUpdate({ game_username: uname, bio: bio });
    editModal?.classList.remove('active');
  });

  // --- Avatar Options Modal ---
  const avatarOptsModal = document.getElementById('avatar-options-modal');
  const dicebearModal = document.getElementById('dicebear-modal');
  const cropperModal = document.getElementById('cropper-modal');
  
  document.getElementById('btn-change-avatar')?.addEventListener('click', () => avatarOptsModal?.classList.add('active'));
  document.getElementById('close-avatar-options')?.addEventListener('click', () => avatarOptsModal?.classList.remove('active'));
  
  document.getElementById('opt-gif')?.addEventListener('click', () => showToast('خرید گیف به زودی!', 'error'));
  document.getElementById('opt-frame')?.addEventListener('click', () => showToast('خرید قاب به زودی!', 'error'));
  
  document.getElementById('opt-dicebear')?.addEventListener('click', () => {
      avatarOptsModal?.classList.remove('active');
      dicebearModal?.classList.add('active');
  });

  // --- DiceBear Modal ---
  let selectedAvatar = currentUser.avatar || '';
  document.getElementById('close-dicebear-btn')?.addEventListener('click', () => dicebearModal?.classList.remove('active'));
  document.querySelectorAll('#dicebear-grid .avatar-option').forEach(img => img.addEventListener('click', (e) => { 
      document.querySelectorAll('#dicebear-grid .avatar-option').forEach(i => i.classList.remove('selected')); 
      (e.target as HTMLElement).classList.add('selected'); 
      selectedAvatar = (e.target as HTMLElement).getAttribute('data-url') || ''; 
  }));
  document.getElementById('save-dicebear-btn')?.addEventListener('click', () => { 
      handleUpdate({ avatar: selectedAvatar }); 
      dicebearModal?.classList.remove('active'); 
  });

  // --- Gallery & Cropper Logic ---
  document.getElementById('opt-gallery')?.addEventListener('click', () => {
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
                  aspectRatio: 1,
                  viewMode: 1,
                  dragMode: 'move',
                  autoCropArea: 1,
                  background: false,
                  cropBoxResizable: false, // Fixed size circle
                  cropBoxMovable: false,   // Fixed position circle
                  toggleDragModeOnDblclick: false
              });
          };
          reader.readAsDataURL(file);
      };
      input.click();
  });

  document.getElementById('close-crop-btn')?.addEventListener('click', () => {
      cropperModal?.classList.remove('active');
      if (cropperInstance) cropperInstance.destroy();
  });

  document.getElementById('save-crop-btn')?.addEventListener('click', () => {
      if (!cropperInstance) return;
      cropperInstance.getCroppedCanvas({ width: 256, height: 256 }).toBlob(async (blob: Blob | null) => {
          if (!blob) return;
          cropperModal?.classList.remove('active');
          showToast('در حال آپلود...', 'success');
          const res = await uploadImage(blob);
          if (res.ok && res.url) {
              await handleUpdate({ avatar: res.url }); 
              showToast('آواتار با موفقیت تغییر کرد', 'success');
          } else {
              showToast('خطا در آپلود: ' + (res.error || 'ناشناخته'));
          }
      }, 'image/jpeg', 0.9);
  });

  // --- Banner Modal ---
  const bannerModal = document.getElementById('banner-modal');
  document.getElementById('btn-edit-banner')?.addEventListener('click', () => bannerModal?.classList.add('active'));
  document.getElementById('close-banner-btn')?.addEventListener('click', () => bannerModal?.classList.remove('active'));
  document.querySelectorAll('.banner-option').forEach(bg => bg.addEventListener('click', (e) => {
      const bgStyle = (e.target as HTMLElement).getAttribute('data-bg');
      if (bgStyle) { handleUpdate({ banner_id: bgStyle }); bannerModal?.classList.remove('active'); }
  }));

  // --- Like & Settings ---
  document.getElementById('like-box')?.addEventListener('click', () => {
      const isLiked = localStorage.getItem('liked_' + currentUser!.id) === 'true';
      localStorage.setItem('liked_' + currentUser!.id, String(!isLiked));
      buildUI();
  });

  const settingsModal = document.getElementById('settings-modal');
  document.getElementById('btn-open-settings')?.addEventListener('click', () => settingsModal?.classList.add('active'));
  document.getElementById('close-settings-btn')?.addEventListener('click', () => settingsModal?.classList.remove('active'));
}

async function handleUpdate(patch: any) {
  if (!currentUser) return;
  const res = await updateProfile(patch);
  if (res.ok) { 
      Object.assign(currentUser, patch); 
      buildUI(); 
  } else { 
      showToast(res.error || 'خطا در ذخیره اطلاعات'); 
  }
}

loadProfile();