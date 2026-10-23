import { logInfo, logError } from './logger'; // فرض بر وجود لاگر

function getUserId(): number {
  const fromQuery = new URLSearchParams(location.search).get('id');
  if (fromQuery && /^\d+$/.test(fromQuery)) return parseInt(fromQuery, 10);
  const tg = (window as any).Telegram?.WebApp;
  return Number(tg?.initDataUnsafe?.user?.id) || 1001;
}

function renderProfile() {
  const app = document.getElementById('app');
  if (!app) return;
  const uid = getUserId();

  app.innerHTML = `
    <div class="app-shell">
      <!-- 🛠️ تست بیلد: این بخش باید در صفحه پروفایل دیده شود -->
      <div style="background: #ef4444; color: white; padding: 12px; text-align: center; font-weight: bold; font-size: 15px; border-bottom: 2px solid #b91c1c;">
        🚨 تست بیلد CI/CD: نسخه جدید پروفایل فعال شد!
      </div>
      
      <div class="top-header">
        <div class="logo-text">Angel</div>
        <div class="header-avatar" onclick="window.location.hash='#home'">
          <img src="https://api.dicebear.com/7.x/bottts/svg?seed=${uid}" style="width:36px;height:36px;border-radius:50%;border:2px solid var(--accent-blue);" />
        </div>
      </div>
      <div class="page-content" style="text-align:center; margin-top: 20px;">
        <img src="https://api.dicebear.com/7.x/bottts/svg?seed=${uid}" style="width:100px;height:100px;border-radius:50%;border:4px solid var(--accent-blue);" />
        <h2 style="margin-top: 10px;">User ${uid}</h2>
        <p style="color:var(--text-muted);">آیدی شما با موفقیت شناسایی شد</p>
        
        <button onclick="window.location.hash='#home'" style="margin-top: 20px; background: var(--accent-blue); color: white; border: none; padding: 10px 20px; border-radius: 8px; font-family: inherit;">بازگشت به خانه</button>
      </div>
    </div>
  `;
}

document.addEventListener('DOMContentLoaded', renderProfile);
