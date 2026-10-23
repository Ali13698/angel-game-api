import { fetchProfile, updateProfile } from './api';
import { renderProfile, renderModals } from './ui';
import { UserProfile } from './types';

let currentUser: UserProfile | null = null;
let isLiked = false;
let selectedAvatar = '';

function getUserId(): number {
  const fromQuery = new URLSearchParams(location.search).get('id');
  if (fromQuery && /^\d+$/.test(fromQuery)) return parseInt(fromQuery, 10);
  const tg = (window as any).Telegram?.WebApp;
  return Number(tg?.initDataUnsafe?.user?.id) || 1001;
}

function showToast(msg: string, color: string = '#3b82f6') {
    const t = document.createElement('div');
    t.style.cssText = `position:fixed; top:20px; left:50%; transform:translateX(-50%); background:${color}; color:white; padding:10px 20px; border-radius:8px; z-index:9999; font-family:inherit; transition:0.3s;`;
    t.innerText = msg;
    document.body.appendChild(t);
    setTimeout(() => t.remove(), 3000);
}

function bindEvents(userId: number) {
    const $ = (id: string) => document.getElementById(id);

    // مودال‌ها
    const avatarModal = $('avatar-modal');
    const bannerModal = $('banner-modal');
    const editModal = $('edit-modal');

    $('btn-change-avatar')?.addEventListener('click', () => avatarModal?.classList.add('active'));
    $('close-avatar-btn')?.addEventListener('click', () => avatarModal?.classList.remove('active'));
    
    document.querySelectorAll('.avatar-option').forEach(el => {
        el.addEventListener('click', (e) => {
            document.querySelectorAll('.avatar-option').forEach(x => x.classList.remove('selected'));
            (e.target as HTMLElement).classList.add('selected');
            selectedAvatar = (e.target as HTMLElement).getAttribute('data-url') || '';
        });
    });

    $('save-avatar-btn')?.addEventListener('click', async () => {
        if (!selectedAvatar) return;
        if (await updateProfile(userId, { avatar: selectedAvatar })) {
            showToast('آواتار با موفقیت تغییر کرد', '#22c55e');
            window.location.reload();
        }
    });

    $('btn-edit-banner')?.addEventListener('click', () => bannerModal?.classList.add('active'));
    $('close-banner-btn')?.addEventListener('click', () => bannerModal?.classList.remove('active'));
    
    document.querySelectorAll('.banner-option').forEach(el => {
        el.addEventListener('click', async (e) => {
            const bg = (e.target as HTMLElement).getAttribute('data-bg');
            if (bg && await updateProfile(userId, { banner_id: bg })) {
                showToast('بنر با موفقیت تغییر کرد', '#22c55e');
                window.location.reload();
            }
        });
    });

    $('btn-edit-bio')?.addEventListener('click', () => editModal?.classList.add('active'));
    $('close-edit-btn')?.addEventListener('click', () => editModal?.classList.remove('active'));

    $('save-profile-btn')?.addEventListener('click', async () => {
        const uname = ($('username-input') as HTMLInputElement).value.trim();
        const bio = ($('bio-input') as HTMLTextAreaElement).value.trim();

        // اعتبارسنجی نام کاربری (فقط حروف و اعداد انگلیسی - حداکثر ۱۰)
        if (!/^[a-zA-Z0-9]{1,10}$/.test(uname)) {
            showToast('نام کاربری فقط حروف و اعداد (حداکثر ۱۰ کاراکتر)', '#ef4444');
            return;
        }
        if (bio.length > 100) {
            showToast('بیوگرافی حداکثر ۱۰۰ کاراکتر', '#ef4444');
            return;
        }

        if (await updateProfile(userId, { game_username: uname, bio: bio })) {
            showToast('اطلاعات با موفقیت ذخیره شد', '#22c55e');
            editModal?.classList.remove('active');
            setTimeout(() => window.location.reload(), 1000);
        }
    });

    // لایک (قلب)
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
        // در اینجا می‌توانید درخواست API برای ثبت لایک در دیتابیس بفرستید
    });
}

async function init() {
    const app = document.getElementById('app');
    if (!app) return;
    app.innerHTML = '<div style="text-align:center; padding:50px;">در حال بارگذاری...</div>';

    const uid = getUserId();
    const user = await fetchProfile(uid);

    if (!user) {
        app.innerHTML = '<div style="text-align:center; color:red; padding:50px;">کاربر پیدا نشد</div>';
        return;
    }

    currentUser = user;
    app.innerHTML = `<div class="container">${renderProfile(user, isLiked)}</div>${renderModals(user)}`;
    bindEvents(uid);
}

document.addEventListener('DOMContentLoaded', init);
