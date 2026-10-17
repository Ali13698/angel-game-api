import { User } from './types';
import { fetchToken, updateProfile } from './api';
import { render, escapeHtml } from './ui';

const appEl = document.getElementById('app');
let currentUser: User | null = null;

async function loadProfile() {
  const p = new URLSearchParams(location.search).get('id');
  if (!p) { if (appEl) appEl.innerHTML = '<div class="status">No ID</div>'; return; }
  try {
    const r = await fetch('/api/profile/' + p);
    const d = await r.json();
    if (!d.ok) { if (appEl) appEl.innerHTML = '<div class="status">Not Found</div>'; return; }
    currentUser = d.user;
    if (appEl) render(currentUser!, appEl, handleBioClick, handleEmojiClick);
  } catch(e) { if (appEl) appEl.innerHTML = '<div class="status">Error</div>'; }
}

function handleBioClick() {
  const newBio = prompt('بیو جدید:', currentUser?.bio || '');
  if (newBio === null) return;
  updateProfile({ bio: newBio }).then(ok => {
    if (ok && currentUser && appEl) { currentUser.bio = newBio; render(currentUser, appEl, handleBioClick, handleEmojiClick); }
  });
}

function handleEmojiClick(e: string) {
  updateProfile({ name_emoji: e }).then(ok => {
    if (ok && currentUser && appEl) { currentUser.name_emoji = e; render(currentUser, appEl, handleBioClick, handleEmojiClick); }
  });
}

loadProfile();
