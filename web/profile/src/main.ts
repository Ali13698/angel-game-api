import { fetchProfile } from './api';
import { renderProfile } from './ui';
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

    app.innerHTML = renderProfile(currentUser, isLiked);
    
    // رویداد لایک
    document.getElementById('like-box')?.addEventListener('click', () => {
        isLiked = !isLiked;
        const heartSvg = document.getElementById('heart-icon')?.querySelector('svg');
        const likesLabel = document.getElementById('likes-count');
        if (heartSvg) {
            heartSvg.setAttribute('fill', isLiked ? '#ef4444' : 'none');
            heartSvg.setAttribute('stroke', isLiked ? '#ef4444' : '#ffffff');
        }
        if (likesLabel && currentUser) {
            currentUser.likes_count += isLiked ? 1 : -1;
            likesLabel.innerText = currentUser.likes_count.toString();
        }
    });

    // باز و بسته کردن مودال‌ها
    document.getElementById('btn-change-avatar')?.addEventListener('click', () => document.getElementById('avatar-modal')?.classList.add('active'));
    document.getElementById('close-avatar-btn')?.addEventListener('click', () => document.getElementById('avatar-modal')?.classList.remove('active'));
    
    document.getElementById('btn-edit-banner')?.addEventListener('click', () => document.getElementById('banner-modal')?.classList.add('active'));
    document.getElementById('close-banner-btn')?.addEventListener('click', () => document.getElementById('banner-modal')?.classList.remove('active'));

    // انتخاب آواتار و بنر (فقط آپدیت ظاهری برای تست)
    document.querySelectorAll('.avatar-option').forEach(el => {
        el.addEventListener('click', (e) => {
            const url = (e.target as HTMLElement).getAttribute('data-url');
            if(url) {
                (document.getElementById('profile-avatar-img') as HTMLImageElement).src = url;
                document.getElementById('avatar-modal')?.classList.remove('active');
            }
        });
    });

    document.querySelectorAll('.banner-option').forEach(el => {
        el.addEventListener('click', (e) => {
            const bg = (e.target as HTMLElement).getAttribute('data-bg');
            if(bg) {
                (document.getElementById('profile-banner') as HTMLElement).style.background = bg;
                document.getElementById('banner-modal')?.classList.remove('active');
            }
        });
    });
}

document.addEventListener('DOMContentLoaded', init);
