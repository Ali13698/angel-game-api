import { User, GameStat } from './types';
import { updateProfile, uploadImage } from './api';
import { renderApp } from './ui/profile';
import { renderAvatarModal, renderEditModal, renderBannerModal, renderSettingsModal } from './ui/modals';
import { showToast } from './ui/toast';
import { DEFAULT_BANNERS, DEFAULT_AVATARS } from './store';
import { logInfo, logError } from './logger';

const appEl = document.getElementById('app');
let currentUser: User | null = null;
let currentStats: GameStat[] = [];

window.addEventListener('unhandledrejection', (event) => {
    logError('Unhandled Promise Rejection', event.reason);
    showToast('خطای غیرمنتظره رخ داد. لطفاً کنسول را بررسی کنید.');
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

  const container = document.querySelector('.container');
  if (container) {
    const modalsWrapper = document.createElement('div');
    modalsWrapper.innerHTML = renderAvatarModal(currentUser, DEFAULT_AVATARS) + 
                              renderEditModal(currentUser) + 
                              renderBannerModal(DEFAULT_BANNERS) + 
                              renderSettingsModal();
    document.body.appendChild(modalsWrapper);
    attachEvents();
  }
}

function attachEvents() {
  if (!currentUser) return;
  logInfo('Attaching UI events...');

  const editModal = document.getElementById('edit-modal');
  document.getElementById('btn-edit-bio')?.addEventListener('click', () => editModal?.classList.add('active'));
  document.getElementById('close-edit-btn')?.addEventListener('click', () => editModal?.classList.remove('active'));
  document.getElementById('save-profile-btn')?.addEventListener('click', () => {
    const uname = (document.getElementById('username-input') as HTMLInputElement).value;
    const bio = (document.getElementById('bio-input') as HTMLTextAreaElement).value;
    logInfo('Saving profile changes...', { uname, bio });
    if (!/^[a-zA-Z0-9]{1,10}$/.test(uname)) {
        showToast('نام کاربری فقط حروف و اعداد (حداکثر ۱۰ کاراکتر) باشد');
        return;
    }
    handleUpdate({ game_username: uname, bio: bio });
    editModal?.classList.remove('active');
  });

  let selectedAvatar = currentUser.avatar || '';
  const avatarModal = document.getElementById('avatar-modal');
  document.getElementById('btn-change-avatar')?.addEventListener('click', () => avatarModal?.classList.add('active'));
  document.getElementById('close-avatar-btn')?.addEventListener('click', () => avatarModal?.classList.remove('active'));
  document.querySelectorAll('.avatar-option').forEach(img => img.addEventListener('click', (e) => { 
      if ((e.target as HTMLElement).id === 'btn-upload-avatar') return;
      document.querySelectorAll('.avatar-option').forEach(i => i.classList.remove('selected')); 
      (e.target as HTMLElement).classList.add('selected'); 
      selectedAvatar = (e.target as HTMLElement).getAttribute('data-url') || ''; 
  }));
  document.getElementById('save-avatar-btn')?.addEventListener('click', () => { 
      logInfo('Saving selected avatar:', selectedAvatar);
      handleUpdate({ avatar: selectedAvatar }); 
      avatarModal?.classList.remove('active'); 
  });

  document.getElementById('btn-upload-avatar')?.addEventListener('click', () => document.getElementById('avatar-file-input')?.click());
  document.getElementById('avatar-file-input')?.addEventListener('change', async (e) => { 
      const file = (e.target as HTMLInputElement).files?.[0]; 
      if (file) { 
          logInfo('Starting avatar upload for file:', file.name);
          const res = await uploadImage(file); 
          if (res.ok && res.url) {
              logInfo('Avatar uploaded, updating profile...');
              await handleUpdate({ avatar: res.url }); 
              avatarModal?.classList.remove('active');
              showToast('آواتار با موفقیت تغییر کرد', 'success');
          } else {
              logError('Avatar upload failed in UI', res);
              showToast('خطا در آپلود: ' + (res.error || 'ناشناخته'));
          }
      } 
  });

  const bannerModal = document.getElementById('banner-modal');
  document.getElementById('btn-edit-banner')?.addEventListener('click', () => bannerModal?.classList.add('active'));
  document.getElementById('close-banner-btn')?.addEventListener('click', () => bannerModal?.classList.remove('active'));
  document.querySelectorAll('.banner-option').forEach(bg => bg.addEventListener('click', (e) => {
      const bgStyle = (e.target as HTMLElement).getAttribute('data-bg');
      if (bgStyle) { 
          logInfo('Saving selected banner:', bgStyle);
          handleUpdate({ banner_id: bgStyle }); 
          bannerModal?.classList.remove('active'); 
      }
  }));

  document.getElementById('like-box')?.addEventListener('click', () => {
      const isLiked = localStorage.getItem('liked_' + currentUser!.id) === 'true';
      logInfo('Toggling like status to:', !isLiked);
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
      logInfo('Update successful, rebuilding UI');
      Object.assign(currentUser, patch); 
      buildUI(); 
  } else { 
      logError('Update failed', res.error);
      showToast(res.error || 'خطا در ذخیره اطلاعات'); 
  }
}

loadProfile();