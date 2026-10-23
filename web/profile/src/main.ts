import { fetchProfile } from './api';
import { renderProfileCard, renderGamesPage } from './ui';
import { User } from './types';

let currentUser: User | null = null;
let isLiked = false;

function getUserId(): number {
    const urlParams = new URLSearchParams(location.search);
    const id = urlParams.get('id');
    if (id) return parseInt(id, 10);
    const tg = (window as any).Telegram?.WebApp;
    return Number(tg?.initDataUnsafe?.user?.id) || 1001;
}

function renderAppShell(content: string, activeTab: 'profile' | 'games'): string {
    return `
    <div class="app-shell">
        <div class="top-header" style="background: #0c0f14; padding: 10px; display: flex; justify-content: space-between; align-items: center;">
            <div style="font-weight: bold; font-size: 20px;">Angel</div>
            <img src="${currentUser?.avatar || 'https://api.dicebear.com/7.x/bottts/svg?seed=' + (currentUser?.id || 1)}" style="width: 36px; height: 36px; border-radius: 50%; border: 2px solid var(--accent-blue);" />
        </div>
        <div class="page-content" style="padding-bottom: 80px;">
            ${content}
        </div>
        <div class="bottom-nav">
            <div class="nav-item ${activeTab === 'profile' ? 'active' : ''}" onclick="switchTab('profile')">
                <span class="nav-icon">👤</span>
            </div>
            <div class="nav-center-btn" onclick="switchTab('games')">
                <span class="nav-icon">🎮</span>
            </div>
            <div class="nav-item" onclick="window.location.href='https://t.me'">
                <span class="nav-icon">💬</span>
            </div>
        </div>
    </div>`;
}

(window as any).switchTab = (tab: string) => {
    const app = document.getElementById('app');
    if (!app || !currentUser) return;

    if (tab === 'profile') {
        app.innerHTML = renderAppShell(`<div class="container">${renderProfileCard(currentUser, isLiked)}</div>`, 'profile');
        attachProfileEvents();
    } else if (tab === 'games') {
        app.innerHTML = renderAppShell(renderGamesPage(), 'games');
    }
};

function attachProfileEvents() {
    const $ = (id: string) => document.getElementById(id);
    
    // لایک
    $('like-box')?.addEventListener('click', () => {
        isLiked = !isLiked;
        const heartSvg = $('heart-icon')?.querySelector('svg');
        const likesLabel = $('likes-count');
        if (heartSvg) {
            heartSvg.setAttribute('fill', isLiked ? '#ef4444' : 'none');
            heartSvg.setAttribute('stroke', isLiked ? '#ef4444' : '#ffffff');
        }
        if (likesLabel && currentUser) {
            currentUser.likes_count += isLiked ? 1 : -1;
            likesLabel.innerText = currentUser.likes_count.toString();
        }
    });
}

async function init() {
    const app = document.getElementById('app');
    if (!app) return;
    app.innerHTML = '<div style="text-align:center; padding:50px; color:white;">در حال بارگذاری...</div>';

    const uid = getUserId();
    let user = await fetchProfile(uid);
    
    if (!user) {
        user = { id: uid, game_username: 'user_'+uid, display_name: 'User '+uid, avatar: '', bio: '', likes_count: 0, coins: 0, gems: 0, wins: 0, age_days: 1 };
    }
    currentUser = user;

    // نمایش اولیه صفحه پروفایل
    (window as any).switchTab('profile');
}

document.addEventListener('DOMContentLoaded', init);
