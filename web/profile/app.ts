const appEl = document.getElementById('app');
if (appEl) appEl.innerHTML = '<div class="status">در حال دریافت...</div>';
const API_BASE = '';
let currentUser = null;

function getUserIdFromURL() {
  const p = new URLSearchParams(location.search).get('id');
  return p ? parseInt(p, 10) : null;
}

async function fetchToken(userId) {
  try {
    const r = await fetch(API_BASE + '/api/auth/token?id=' + userId);
    const d = await r.json();
    return d.ok ? d.token : null;
  } catch { return null; }
}

async function updateBio(newBio) {
  const uid = getUserIdFromURL();
  const token = await fetchToken(uid);
  if (!token) return false;
  try {
    const r = await fetch(API_BASE + '/api/me', {
      method: 'PATCH',
      headers: { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' },
      body: JSON.stringify({ bio: newBio })
    });
    return r.ok;
  } catch { return false; }
}

function escapeHtml(s) {
  return String(s).replace(/[&<>"']/g, function(c) { return {'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]; });
}

function render(user) {
  if (!appEl) return;
  currentUser = user;
  const displayName = user.display_name || user.first_name || user.game_username || ('User' + user.id);
  const color = user.profile_color || '#8b5cf6';
  appEl.innerHTML = '<div class="banner"></div>' +
    '<div class="avatar-wrap"><div class="avatar">' + (user.avatar ? '<img src="' + escapeHtml(user.avatar) + '">' : '🎮') + '</div></div>' +
    '<div class="name" style="color:' + color + '">' + escapeHtml(displayName) + '</div>' +
    '<div class="username">@' + escapeHtml(user.game_username || 'unknown') + '</div>' +
    '<div class="bio">' + escapeHtml(user.bio || '---') + '</div>' +
    '<div class="stats"><div><div class="stat-value">' + user.coins + '</div><div class="stat-label">سکه</div></div>' +
    '<div><div class="stat-value">' + user.gems + '</div><div class="stat-label">الماس</div></div>' +
    '<div><div class="stat-value">' + user.level + '</div><div class="stat-label">سطح</div></div></div>' +
    '<button class="btn" id="btn-edit-bio">تغییر بیو</button>';
  
  document.getElementById('btn-edit-bio').addEventListener('click', function() {
    document.getElementById('bio-input').value = currentUser.bio || '';
    document.getElementById('bio-modal').classList.remove('hidden');
  });
}

document.getElementById('btn-cancel-bio').addEventListener('click', function() {
  document.getElementById('bio-modal').classList.add('hidden');
});

document.getElementById('btn-confirm-bio').addEventListener('click', async function() {
  const newBio = document.getElementById('bio-input').value;
  const ok = await updateBio(newBio);
  if (ok) {
    document.getElementById('bio-modal').classList.add('hidden');
    currentUser.bio = newBio;
    const bioEl = document.querySelector('.bio');
    if (bioEl) bioEl.textContent = newBio || '---';
  } else {
    alert('خطا در ذخیره');
  }
});

async function main() {
  const uid = getUserIdFromURL();
  if (!uid) { if (appEl) appEl.innerHTML = '<div class="status">No ID</div>'; return; }
  try {
    const r = await fetch(API_BASE + '/api/profile/' + uid);
    const d = await r.json();
    if (!d.ok) { if (appEl) appEl.innerHTML = '<div class="status">Not Found</div>'; return; }
    render(d.user);
  } catch(e) { if (appEl) appEl.innerHTML = '<div class="status">Error</div>'; }
}
main();