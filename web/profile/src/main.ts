import { fetchProfile, updateProfile } from './api';
import { renderProfileCard } from './ui';
import { renderAvatarModal, renderBannerModal, renderEditModal } from './modals';
import { User } from './types';

let currentUser: User | null = null;
let isLiked = false;
let selectedAvatar = '';

function getUserId(): number {
    const urlParams = new URLSearchParams(location.search);
    const id = urlParams.get('id');
    if (id) return parseInt(id, 10);
    const tg = (window as any).Telegram?.WebApp;
    return Number(tg?.initDataUnsafe?.user?.id) || 1001;
}

function showToast(msg: string, color: string = '#3b82f6') {
    const t = document.createElement('div');
    t.style.cssText = `position:fixed; top:20px; left:50%; transform:translateX(-50%); background:${color}; color:white; padding:10px 20px; border-radius:8px; z-index:9999;`;
    t.innerText = msg;
    document.body.appendChild(t);
    setTimeout(() => t.remove(), 3000);
}

function attachEvents(userId: number) {
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

    // مودال آواتار
    $('btn-change-avatar')?.addEventListener('click', () => $('avatar-modal')?.classList.add('active'));
    $('close-avatar-btn')?.addEventListener('click', () => $('avatar-modal')?.classList.remove('active'));
    document.querySelectorAll('.avatar-option').forEach(el => {
        el.addEventListener('click', (e) => {
            document.querySelectorAll('.avatar-option').forEach(x => (x as HTMLElement).style.border = '2px solid transparent');
            (e.target as HTMLElement).style.border = '2px solid #3b82f6';
            selectedAvatar = (e.target as HTMLElement).getAttribute('data-url') || '';
        });
    });
    $('save-avatar-btn')?.addEventListener('click', async () => {
        if (selectedAvatar && await updateProfile(userId, { avatar: selectedAvatar })) {
            showToast('آواتار تغییر کرد', '#22c55e');
            window.location.reload();
        }
    });

    // مودال بنر
    $('btn-edit-banner')?.addEventListener('click', () => $('banner-modal')?.classList.add('active'));
    $('close-banner-btn')?.addEventListener('click', () => $('banner-modal')?.classList.remove('active'));
    document.querySelectorAll('.banner-option').forEach(el => {
        el.addEventListener('click', async (e) => {
            const bg = (e.target as HTMLElement).getAttribute('data-bg');
            if (bg && await updateProfile(userId, { banner_id: bg })) {
                showToast('بنر تغییر کرد', '#22c55e');
                window.location.reload();
            }
        });
    });

    // مودال ویرایش
    $('btn-edit-bio')?.addEventListener('click', () => $('edit-modal')?.classList.add('active'));
    $('close-edit-btn')?.addEventListener('click', () => $('edit-modal')?.classList.remove('active'));
    $('save-profile-btn')?.addEventListener('click', async () => {
        const uname = ($('username-input') as HTMLInputElement).value;
        const bio = ($('bio-input') as HTMLTextAreaElement).value;
        
        if (!/^[a-zA-Z0-9]{1,10}$/.test(uname)) {
            showToast('نام کاربری فقط انگلیسی و عدد (حداکثر ۱۰)', '#ef4444');
            return;
        }
        if (await updateProfile(userId, { game_username: uname, bio: bio })) {
            showToast('اطلاعات ذخیره شد', '#22c55e');
            window.location.reload();
        }
    });
}

async function init() {
    const app = document.getElementById('app');
    if (!app) return;
    
    const uid = getUserId();
    currentUser = await fetchProfile(uid);
    
    if (!currentUser) {
        currentUser = { id: uid, game_username: 'user_'+uid, display_name: 'User '+uid, avatar: '', bio: '', likes_count: 0, coins: 0, gems: 0, wins: 0, age_days: 1 };
    }
    
    app.innerHTML = `
        <div class="app-shell">
            <div style="background:#22c55e; color:white; padding:8px; text-align:center; font-size:13px; font-weight:bold;">✅ کدهای TypeScript با موفقیت اجرا شدند!</div>
            <div class="page-content">
                ${renderProfileCard(currentUser, isLiked)}
            </div>
        </div>
        ${renderAvatarModal(currentUser)}
        ${renderBannerModal()}
        ${renderEditModal(currentUser)}
    `;
    
    attachEvents(uid);
}

document.addEventListener('DOMContentLoaded', init);
