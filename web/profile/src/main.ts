import { User, GameStat } from './types';
import { updateProfile } from './api';
import { render } from './ui';
const appEl = document.getElementById('app');
let currentUser: User | null = null;
let currentStats: GameStat[] = [];
async function loadProfile() {
  const p = new URLSearchParams(location.search).get('id');
  if (!p) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">No ID</div>'; return; }
  try {
    const r = await fetch('/api/profile/' + p); const d = await r.json();
    if (!d.ok) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">Not Found</div>'; return; }
    currentUser = d.user; currentStats = d.stats || [];
    if (appEl && currentUser) render(currentUser, currentStats, appEl, handleSaveProfile, handleLikeClick);
  } catch(e) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">Error</div>'; }
}
async function handleSaveProfile(avatar: string, bio: string) {
  if (!currentUser) return;
  const ok = await updateProfile({ avatar, bio });
  if (ok && appEl) { currentUser.avatar = avatar; currentUser.bio = bio; render(currentUser, currentStats, appEl, handleSaveProfile, handleLikeClick); document.getElementById('edit-modal')?.classList.remove('active'); } else { alert('خطا در ذخیره اطلاعات'); }
}
function handleLikeClick() {
  if (!currentUser) return;
  const isLiked = localStorage.getItem('liked_' + currentUser.id) === 'true';
  if (!isLiked) { localStorage.setItem('liked_' + currentUser.id, 'true'); document.getElementById('heart-icon')!.style.color = '#ef4444'; document.getElementById('likes-count')!.innerText = (currentUser.likes_count + 1).toString(); }
  else { localStorage.setItem('liked_' + currentUser.id, 'false'); document.getElementById('heart-icon')!.style.color = '#ffffff'; document.getElementById('likes-count')!.innerText = (currentUser.likes_count).toString(); }
}
loadProfile();