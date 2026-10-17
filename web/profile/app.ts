const appEl = document.getElementById('app');
if (appEl) appEl.innerHTML = '<div class="status">در حال دریافت...</div>';
const API_BASE = '';
const PRESET_EMOJIS = ['🐉','🦊','🌟','🐼','🦁','🐸','🐺','🦄','🐯','🐨','👾','🤖','🧑','👩','🧔','👑'];
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

async function updateProfile(patch) {
  const uid = getUserIdFromURL();
  const token = await fetchToken(uid);
  if (!token) return false;
  try {
    const r = await fetch(API_BASE + '/api/me', {
      method: 'PATCH',
      headers: { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' },
      body: JSON.stringify(patch)
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
  const avatarContent = user.avatar ? '<img src="' + escapeHtml(user.avatar) + '">' : (user.name_emoji || '🎮');
  
  appEl.innerHTML = '<div class="banner"></div>' +
    '<div class="avatar-wrap"><div class="avatar">' + avatarContent + '</div></div>' +
    '<div class="name" style="color:' + color + '">' + escapeHtml(displayName) + '</div>' +
    '<div class="username">@' + escapeHtml(user.game_username || 'unknown') + '</div>' +
    '<div class="bio">' + escapeHtml(user.bio || '---') + '</div>' +
    '<div class="stats"><div><div class="stat-value">' + user.coins + '</div><div class="stat-label">سکه</div></div>' +
    '<div><div class="stat-value">' + user.gems + '</div><div class="stat-label">الماس</div></div>' +
    '<div><div class="stat-value">' + user.level + '</div><div class="stat-label">سطح</div></div></div>' +
    '<button class="btn" id="btn-edit-bio">تغییر بیو</button>' +
    '<div class="section"><div class="section-title">ایموجی نام</div><div class="emoji-row" id="emoji-row">' +
    PRESET_EMOJIS.map(function(e) { return '<button class="emoji-btn' + (user.name_emoji === e ? ' active' : '') + '" data-emoji="' + e + '">' + e + '</button>'; }).join('') +
    '</div></div>';
  
  document.getElementById('btn-edit-bio').addEventListener('click', function() {
    document.getElementById('bio-input').value = currentUser.bio || '';
    document.getElementById('bio-modal').classList.remove('hidden');
  });
  
  document.querySelectorAll('.emoji-btn').forEach(function(btn) {
    btn.addEventListener('click', async function() {
      const e = btn.dataset.emoji;
      const ok = await updateProfile({ name_emoji: e });
      if (ok) {
        currentUser.name_emoji = e;
        render(currentUser); // Re-render to show new emoji and active state
      }
    });
  });
}

document.getElementById('btn-cancel-bio').addEventListener('click', function() {
  document.getElementById('bio-modal').classList.add('hidden');
});

document.getElementById('btn-confirm-bio').addEventListener('click', async function() {
  const newBio = document.getElementById('bio-input').value;
  const ok = await updateProfile({ bio: newBio });
  if (ok) {
    document.getElementById('bio-modal').classList.add('hidden');
    currentUser.bio = newBio;
    const bioEl = document.querySelector('.bio');
    if (bioEl) bioEl.textContent = newBio || '---';
  } else { alert('خطا در ذخیره'); }
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