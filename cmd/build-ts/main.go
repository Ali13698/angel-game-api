package main

import (
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

const htmlContent = `<!doctype html>
<html lang="fa" dir="rtl">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
  <title>پروفایل — Angel</title>
  <link rel="stylesheet" href="style.css">
</head>
<body>
  <div class="card" id="app">
    <div class="status">در حال بارگذاری...</div>
  </div>
  <script type="module" src="app.js"></script>
</body>
</html>`

const cssContent = `* { box-sizing: border-box; margin: 0; padding: 0; }
body {
  font-family: system-ui, -apple-system, 'Segoe UI', Tahoma, sans-serif;
  background: #0a0e1a;
  color: #fff;
  min-height: 100vh;
  padding: 16px;
}
.card {
  max-width: 480px;
  margin: 0 auto;
  background: #131826;
  border-radius: 20px;
  overflow: hidden;
}
.banner {
  height: 140px;
  background: linear-gradient(135deg, #7c3aed, #ec4899);
  position: relative;
}
.avatar-wrap {
  position: relative;
  margin-top: -60px;
  display: flex;
  justify-content: center;
}
.avatar {
  width: 120px; height: 120px;
  border-radius: 50%;
  background: #1a1f3a;
  border: 4px solid #131826;
  display: flex; align-items: center; justify-content: center;
  font-size: 56px;
  overflow: hidden;
}
.avatar img { width: 100%; height: 100%; object-fit: cover; }
.name { text-align: center; font-size: 22px; font-weight: 700; margin-top: 10px; }
.username { text-align: center; color: #7dd3fc; font-size: 13px; margin-top: 4px; }
.bio { text-align: center; color: #94a3b8; font-size: 14px; margin: 12px 24px; min-height: 20px; }
.stats {
  display: grid; grid-template-columns: repeat(3, 1fr);
  gap: 8px; margin: 16px 20px;
  text-align: center;
}
.stat-value { font-size: 20px; font-weight: 700; color: #a855f7; }
.stat-label { font-size: 11px; color: #64748b; margin-top: 4px; }
.actions {
  display: flex; gap: 8px; padding: 16px 20px;
}
.btn {
  flex: 1;
  padding: 12px;
  border: 0;
  border-radius: 12px;
  background: #7c3aed;
  color: #fff;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
}
.btn.secondary { background: #1e293b; }
.btn:active { opacity: 0.85; }
.section {
  padding: 16px 20px;
  border-top: 1px solid #1e293b;
}
.section-title { font-size: 13px; color: #64748b; margin-bottom: 10px; }
.emoji-row { display: flex; gap: 8px; flex-wrap: wrap; }
.emoji-btn {
  width: 44px; height: 44px;
  border-radius: 50%;
  border: 2px solid transparent;
  background: #1a1f3a;
  font-size: 22px;
  cursor: pointer;
  display: flex; align-items: center; justify-content: center;
}
.emoji-btn.active { border-color: #a855f7; }
.status { text-align: center; color: #64748b; font-size: 12px; padding: 20px; }
.error { color: #f87171; }`

const tsContent = `const API_BASE = '';
const PRESET_EMOJIS = ['🐉','🦊','🌟','🐼','🦁','🐸','🐺','🦄','🐯','🐨','👾','🤖','🧑','👩','🧔','👑'];
let currentUser: any = null;
let token: string | null = null;

function getUserIdFromURL(): number | null {
  const p = new URLSearchParams(location.search).get('id');
  if (!p) return null;
  const n = parseInt(p, 10);
  return isNaN(n) ? null : n;
}

async function fetchToken(userId: number): Promise<string | null> {
  try {
    const r = await fetch(API_BASE + '/api/auth/token?id=' + userId);
    const d = await r.json();
    return d.ok && d.token ? d.token : null;
  } catch { return null; }
}

async function loadProfile(userId: number): Promise<any> {
  const r = await fetch(API_BASE + '/api/profile/' + userId);
  return r.json();
}

async function patchProfile(patch: any): Promise<any> {
  if (!token) return { ok: false, error: 'no token' };
  const r = await fetch(API_BASE + '/api/me', {
    method: 'PATCH',
    headers: { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' },
    body: JSON.stringify(patch),
  });
  return r.json();
}

function escapeHtml(s: string): string {
  return String(s).replace(/[&<>"']/g, (c: string) => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c] as string));
}

function render(user: any): void {
  const app = document.getElementById('app');
  if (!app) return;
  const displayName = user.display_name || user.first_name || user.game_username || ('User' + user.id);
  const avatarContent = user.avatar
    ? '<img src="' + escapeHtml(user.avatar) + '" alt="">'
    : (user.name_emoji || '🎮');
  const color = user.profile_color || '#8b5cf6';
  app.innerHTML = '<div class="banner" style="background:linear-gradient(135deg,#7c3aed,#ec4899)"></div>' +
    '<div class="avatar-wrap"><div class="avatar">' + avatarContent + '</div></div>' +
    '<div class="name" style="color:' + color + '">' + escapeHtml(displayName) + '</div>' +
    '<div class="username">@' + escapeHtml(user.game_username || 'unknown') + '</div>' +
    '<div class="bio">' + escapeHtml(user.bio || 'بدون بیو') + '</div>' +
    '<div class="stats"><div><div class="stat-value">' + user.wins + '</div><div class="stat-label">برد</div></div>' +
    '<div><div class="stat-value">' + user.likes_count + '</div><div class="stat-label">هاله</div></div>' +
    '<div><div class="stat-value">' + user.age_days + '</div><div class="stat-label">روز</div></div></div>' +
    '<div class="stats"><div><div class="stat-value">' + user.coins + '</div><div class="stat-label">سکه</div></div>' +
    '<div><div class="stat-value">' + user.gems + '</div><div class="stat-label">الماس</div></div>' +
    '<div><div class="stat-value">' + user.level + '</div><div class="stat-label">سطح</div></div></div>' +
    '<div class="actions"><button class="btn" id="btn-edit-bio">تغییر بیو</button>' +
    '<button class="btn secondary" id="btn-edit-emoji">تغییر ایموجی</button></div>' +
    '<div class="section"><div class="section-title">ایموجی نام</div><div class="emoji-row" id="emoji-row">' +
    PRESET_EMOJIS.map((e: string) => '<button class="emoji-btn' + (user.name_emoji === e ? ' active' : '') + '" data-emoji="' + e + '">' + e + '</button>').join('') +
    '</div></div>';
  bindEvents();
}

function bindEvents(): void {
  document.getElementById('btn-edit-bio')?.addEventListener('click', async () => {
    const v = prompt('بیو جدید:', currentUser?.bio || '');
    if (v === null) return;
    const r = await patchProfile({ bio: v });
    if (r.ok && r.user) { currentUser = r.user; render(r.user); }
    else alert('خطا: ' + (r.error || 'unknown'));
  });
  document.querySelectorAll('.emoji-btn').forEach((btn: any) => {
    btn.addEventListener('click', async () => {
      const e = btn.dataset.emoji;
      if (!e) return;
      const r = await patchProfile({ name_emoji: e });
      if (r.ok && r.user) { currentUser = r.user; render(r.user); }
    });
  });
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
    if (app) app.innerHTML = '<div class="status error">' + escapeHtml(r.error || 'not found') + '</div>';
    return;
  }
  currentUser = r.user;
  render(r.user);
}

main();`

func main() {
	log.Println("[build-ts] generating static assets and compiling TS...")
	
	os.MkdirAll("internal/web/static", 0755)
	os.WriteFile("internal/web/static/index.html", []byte(htmlContent), 0644)
	os.WriteFile("internal/web/static/style.css", []byte(cssContent), 0644)

	result := api.Build(api.BuildOptions{
		Stdin: &api.StdinOptions{
			Contents:   tsContent,
			ResolveDir: ".",
			Sourcefile: "app.ts",
			Loader:     api.LoaderTS,
		},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
		Outfile:           "internal/web/static/app.js",
	})
	if len(result.Errors) > 0 {
		log.Fatalf("[build-ts] errors: %v", result.Errors)
	}
	log.Println("[build-ts] OK: app.js, index.html, style.css generated.")
}
