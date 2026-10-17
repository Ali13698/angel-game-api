import { User, GameStat } from './types';
import { updateProfile, uploadImage } from './api';
import { render } from './ui';
import { DEFAULT_BANNERS, DEFAULT_AVATARS } from './store';

const appEl = document.getElementById('app');
let currentUser: User | null = null;
let currentStats: GameStat[] = [];

async function loadProfile() {
  const p = new URLSearchParams(location.search).get('id');
  if (!p) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">No ID provided</div>'; return; }
  try {
    const r = await fetch('/api/profile/' + p);
    const d = await r.json();
    if (!d.ok) throw new Error(d.error || 'User not found');
    currentUser = d.user; 
    currentStats = d.stats || [];
    if (appEl && currentUser) {
      render(currentUser, currentStats, DEFAULT_BANNERS, DEFAULT_AVATARS, appEl, handleUpdate, handleUpload, handleLikeClick);
    }
  } catch(e: any) { 
    if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;color:red;direction:rtl;">Error: ' + e.message + '</div>'; 
  }
}

async function handleUpdate(patch: any) {
  if (!currentUser) return;
  const ok = await updateProfile(patch);
  if (ok && appEl) { Object.assign(currentUser, patch); render(currentUser, currentStats, DEFAULT_BANNERS, DEFAULT_AVATARS, appEl, handleUpdate, handleUpload, handleLikeClick); } 
  else { alert('خطا در ذخیره اطلاعات'); }
}

async function handleUpload(file: File, type: 'avatar' | 'banner') {
  if (!currentUser || !appEl) return;
  const url = await uploadImage(file, type);
  if (url) { const patch = type === 'avatar' ? { avatar: url } : { banner_id: `url("${url}")` }; await handleUpdate(patch); } 
  else { alert('خطا در آپلود عکس'); }
}

function handleLikeClick() {
  if (!currentUser) return;
  const isLiked = localStorage.getItem('liked_' + currentUser.id) === 'true';
  if (!isLiked) { localStorage.setItem('liked_' + currentUser.id, 'true'); render(currentUser, currentStats, DEFAULT_BANNERS, DEFAULT_AVATARS, appEl!, handleUpdate, handleUpload, handleLikeClick); } 
  else { localStorage.setItem('liked_' + currentUser.id, 'false'); render(currentUser, currentStats, DEFAULT_BANNERS, DEFAULT_AVATARS, appEl!, handleUpdate, handleUpload, handleLikeClick); }
}

loadProfile();