import { User } from './types';
import { updateProfile } from './api';
import { render } from './ui';

const appEl = document.getElementById('app');
let currentUser: User | null = null;

async function loadProfile() {
  const p = new URLSearchParams(location.search).get('id');
  if (!p) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">No ID</div>'; return; }
  try {
    const r = await fetch('/api/profile/' + p);
    const d = await r.json();
    if (!d.ok) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">Not Found</div>'; return; }
    currentUser = d.user;
    if (appEl) render(currentUser!, appEl, handleBioClick, () => {});
  } catch(e) { if (appEl) appEl.innerHTML = '<div style="text-align:center;padding:20px;">Error</div>'; }
}

function handleBioClick() {
  const newBio = prompt('بیو جدید:', currentUser?.bio || '');
  if (newBio === null) return;
  updateProfile({ bio: newBio }).then(ok => {
    if (ok && currentUser && appEl) { 
      currentUser.bio = newBio; 
      render(currentUser, appEl, handleBioClick, () => {}); 
    }
  });
}

loadProfile();
