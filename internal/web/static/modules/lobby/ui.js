import { fetchProfileData } from './api.js';

export async function renderProfileScreen() {
    const container = document.getElementById('screen-profile');
    container.innerHTML = '<div style="text-align:center;padding:50px;color:#888;">در حال دریافت اطلاعات...</div>';

    const user = await fetchProfileData();
    
    if (!user) {
        container.innerHTML = '<div style="text-align:center;padding:50px;color:#ef4444;">خطا در دریافت اطلاعات پروفایل</div>';
        return;
    }

    const wins = user.wins || 0;
    const coins = user.coins || 0;
    const gems = user.gems || 0;
    const displayName = user.display_name || user.game_username || "کاربر مهمان";

    // ساختار HTML پروفایل (توسط ui.js تولید می‌شود)
    container.innerHTML = `
        <div class="profile-card">
            <div class="banner"></div>
            <div class="profile-avatar">
                <svg viewBox="0 0 160 160"><rect width="160" height="160" fill="#475346"/><path d="M80 13 C27 10 13 43 25 84 C32 108 54 137 80 147 C104 135 127 109 136 80 C149 37 126 10 80 13Z" fill="#8b9472"/><path d="M34 62 Q62 54 71 88 Q40 102 34 62Z" fill="#07100b"/><path d="M127 62 Q99 54 89 88 Q120 102 127 62Z" fill="#07100b"/></svg>
            </div>
            <div class="profile-info"><h2 class="username">${displayName}</h2></div>
            <div class="metrics">
                <div class="metric"><span>بردها</span><strong>${wins}</strong></div>
                <div class="metric"><span>سکه</span><strong>${coins}</strong></div>
                <div class="metric"><span>الماس</span><strong>${gems}</strong></div>
            </div>
        </div>
    `;
}
