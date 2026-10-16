// Angel Profile — TypeScript

interface User {
  id: number;
  first_name?: string;
  game_username?: string;
  display_name?: string;
  avatar?: string;
  photo_url?: string;
  profile_color?: string;
  profile_banner?: string;
  bio?: string;
  banner_id?: string;
  frame_id?: string;
  hat_id?: string;
  chat_bubble_id?: string;
  badge_id?: string;
  name_emoji?: string;
  coins: number;
  gems: number;
  wins: number;
  losses: number;
  level: number;
  likes_count: number;
  streak_days: number;
  active_days: number;
  age_days: number;
  created_at: number;
}

interface ProfileResp { ok: boolean; user?: User; error?: string; }
interface TokenResp { ok: boolean; token?: string; user_id?: string; error?: string; }

const API_BASE = '';
const PRESET_EMOJIS = ['🐉','🦊','🌟','🐼','🦁','🐸','🐺','🦄','🐯','🐨','👾','🤖','🧑','👩','🧔','👑'];

let currentUser: User | null = null;
let token: string | null = null;

// استخراج userId از URL (?id=1) یا localStorage
function getUserIdFromURL(): number | null {
  const p = new URLSearchParams(location.search).get('id');
  if (!p) return null;
  const n = parseInt(p, 10);
  return isNaN(n) ? null : n;
}

// گرفتن JWT (موقت: از /api/auth/token)
async function fetchToken(userId: number): Promise<string | null> {
  try {
    const r = await fetch(`${API_BASE}/api/auth/token?id=${userId}`);
    const d: TokenResp = await r.json();
    return d.ok && d.token ? d.token : null;
  } catch { return null; }
}

async function loadProfile(userId: number): Promise<ProfileResp> {
  const r = await fetch(`${API_BASE}/api/profile/${userId}`);
  return r.json();
}

async function patchProfile(patch: Partial<User>): Promise<ProfileResp> {
  if (!token) return { ok: false, error: 'no token' };
  const r = await fetch(`${API_BASE}/api/me`, {
    method: 'PATCH',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(patch),
  });
  return r.json();
}

function render(user: User): void {
  const app = document.getElementById('app');
  if (!app) return;

  const displayName = user.display_name || user.first_name || user.game_username || `User${user.id}`;
  const avatarContent = user.avatar
    ? `<img src="${escapeHtml(user.avatar)}" alt="">`
    : (user.name_emoji || '🎮');

  const color = user.profile_color || '#8b5cf6';

  app.innerHTML = `
    <div class="banner" style="background:${escapeHtml(user.banner_id ? 'linear-gradient(135deg,#7c3aed,#ec4899)' : color)}"></div>
    <div class="avatar-wrap">
      <div class="avatar">${avatarContent}</div>
    </div>
    <div class="name" style="color:${color}">${escapeHtml(displayName)}</div>
    <div class="username">@${escapeHtml(user.game_username || 'unknown')}</div>
    <div class="bio">${escapeHtml(user.bio || 'بدون بیو')}</div>

    <div class="stats">
      <div><div class="stat-value">${user.wins}</div><div class="stat-label">برد</div></div>
      <div><div class="stat-value">${user.likes_count}</div><div class="stat-label">هاله</div></div>
      <div><div class="stat-value">${user.age_days}</div><div class="stat-label">روز</div></div>
    </div>

    <div class="stats">
      <div><div class="stat-value">${user.coins}</div><div class="stat-label">سکه</div></div>
      <div><div class="stat-value">${user.gems}</div><div class="stat-label">الماس</div></div>
      <div><div class="stat-value">${user.level}</div><div class="stat-label">سطح</div></div>
    </div>

    <div class="actions">
      <button class="btn" id="btn-edit-bio">تغییر بیو</button>
      <button class="btn secondary" id="btn-edit-emoji">تغییر ایموجی</button>
    </div>

    <div class="section">
      <div class="section-title">ایموجی نام</div>
      <div class="emoji-row" id="emoji-row">
        ${PRESET_EMOJIS.map(e => `<button class="emoji-btn${user.name_emoji === e ? ' active' : ''}" data-emoji="${e}">${e}</button>`).join('')}
      </div>
    </div>
  `;

  bindEvents();
}

function bindEvents(): void {
  const bioBtn = document.getElementById('btn-edit-bio');
  bioBtn?.addEventListener('click', async () => {
    const v = prompt('بیو جدید:', currentUser?.bio || '');
    if (v === null) return;
    const r = await patchProfile({ bio: v });
    if (r.ok && r.user) { currentUser = r.user; render(r.user); }
    else alert('خطا: ' + (r.error || 'unknown'));
  });

  document.querySelectorAll<HTMLButtonElement>('.emoji-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
      const e = btn.dataset.emoji;
      if (!e) return;
      const r = await patchProfile({ name_emoji: e });
      if (r.ok && r.user) { currentUser = r.user; render(r.user); }
    });
  });
}

function escapeHtml(s: string): string {
  return String(s).replace(/[&<>"']/g, c => ({
    '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'
  })[c] as string);
}

async function main(): Promise<void> {
  const userId = getUserIdFromURL();
  const app = document.getElementById('app');
  if (!userId) {
    if (app) app.innerHTML = '<div class="status error">?id= در URL لازمه — مثلاً /profile/?id=1</div>';
    return;
  }
  token = await fetchToken(userId);
  const r = await loadProfile(userId);
  if (!r.ok || !r.user) {
    if (app) app.innerHTML = `<div class="status error">${escapeHtml(r.error || 'not found')}</div>`;
    return;
  }
  currentUser = r.user;
  render(r.user);
}

main();
