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
  <link rel="stylesheet" href="/profile/style.css">
</head>
<body>
  <div class="card" id="app">
    <div class="status">در حال بارگذاری...</div>
  </div>
  <script type="module" src="/profile/app.js"></script>
</body>
</html>`

const cssContent = `* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: system-ui, sans-serif; background: #0a0e1a; color: #fff; min-height: 100vh; padding: 16px; }
.card { max-width: 480px; margin: 0 auto; background: #131826; border-radius: 20px; overflow: hidden; }
.banner { height: 140px; background: linear-gradient(135deg, #7c3aed, #ec4899); }
.avatar-wrap { margin-top: -60px; display: flex; justify-content: center; }
.avatar { width: 120px; height: 120px; border-radius: 50%; background: #1a1f3a; border: 4px solid #131826; display: flex; align-items: center; justify-content: center; font-size: 56px; overflow: hidden; }
.avatar img { width: 100%; height: 100%; object-fit: cover; }
.name { text-align: center; font-size: 22px; font-weight: 700; margin-top: 10px; }
.username { text-align: center; color: #7dd3fc; font-size: 13px; margin-top: 4px; }
.bio { text-align: center; color: #94a3b8; font-size: 14px; margin: 12px 24px; min-height: 20px; }
.stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; margin: 16px 20px; text-align: center; }
.stat-value { font-size: 20px; font-weight: 700; color: #a855f7; }
.stat-label { font-size: 11px; color: #64748b; margin-top: 4px; }
.status { text-align: center; color: #64748b; font-size: 14px; padding: 40px 20px; }
.error { color: #f87171; font-weight: bold; }`

const tsContent = `
window.onerror = function(msg, url, line, col, err) {
  const app = document.getElementById('app');
  if (app) app.innerHTML = '<div class="status error">خطای JS: ' + msg + '</div>';
};

const API_BASE = '';
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
  try {
    const r = await fetch(API_BASE + '/api/profile/' + userId);
    if (!r.ok) throw new Error('HTTP ' + r.status);
    return r.json();
  } catch (e: any) {
    throw new Error('خطا در دریافت اطلاعات: ' + e.message);
  }
}

function escapeHtml(s: string): string {
  return String(s).replace(/[&<>"']/g, (c: string) => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c] as string));
}

function render(user: any): void {
  const app = document.getElementById('app');
  if (!app) return;
  const displayName = user.display_name || user.first_name || user.game_username || ('User' + user.id);
  const avatarContent = user.avatar ? '<img src="' + escapeHtml(user.avatar) + '">' : (user.name_emoji || '🎮');
  const color = user.profile_color || '#8b5cf6';
  
  app.innerHTML = '<div class="banner"></div>' +
    '<div class="avatar-wrap"><div class="avatar">' + avatarContent + '</div></div>' +
    '<div class="name" style="color:' + color + '">' + escapeHtml(displayName) + '</div>' +
    '<div class="username">@' + escapeHtml(user.game_username || 'unknown') + '</div>' +
    '<div class="bio">' + escapeHtml(user.bio || 'بدون بیو') + '</div>' +
    '<div class="stats"><div><div class="stat-value">' + user.wins + '</div><div class="stat-label">برد</div></div>' +
    '<div><div class="stat-value">' + user.likes_count + '</div><div class="stat-label">هاله</div></div>' +
    '<div><div class="stat-value">' + user.age_days + '</div><div class="stat-label">روز</div></div></div>' +
    '<div class="stats"><div><div class="stat-value">' + user.coins + '</div><div class="stat-label">سکه</div></div>' +
    '<div><div class="stat-value">' + user.gems + '</div><div class="stat-label">الماس</div></div>' +
    '<div><div class="stat-value">' + user.level + '</div><div class="stat-label">سطح</div></div></div>';
}

async function main() {
  try {
    const userId = getUserIdFromURL();
    const app = document.getElementById('app');
    if (!userId) {
      if (app) app.innerHTML = '<div class="status error">?id= در URL لازمه</div>';
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
  } catch(e: any) {
    const app = document.getElementById('app');
    if (app) app.innerHTML = '<div class="status error">خطا: ' + e.message + '</div>';
  }
}

main();`

func main() {
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
	log.Println("[build-ts] OK.")
}
