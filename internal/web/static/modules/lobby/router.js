import { renderProfileScreen } from './ui.js';
import { renderGamesScreen } from './games.js';

export class Router {
    static switch(tabName, element) {
        // مخفی کردن تمام صفحات
        document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
        // غیرفعال کردن تمام تب‌ها
        document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));

        const targetScreen = document.getElementById(`screen-${tabName}`);
        
        if (targetScreen) {
            targetScreen.classList.add('active');
        }

        // مدیریت حالت خاص برای پنل اصلی (Home)
        if (tabName === 'home') {
            document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
            document.getElementById('nav-home').classList.add('active');
        } else if (element) {
            element.classList.add('active');
        }

        // بارگذاری محتوای تب‌ها (Modular Rendering)
        if (tabName === 'profile') {
            renderProfileScreen();
        } else if (tabName === 'games') {
            renderGamesScreen();
        } else if (tabName === 'store') {
            // اینجا بعداً ماژول فروشگاه لود می‌شود
            targetScreen.innerHTML = '<div style="text-align:center;padding:50px;color:#888;">فروشگاه به زودی...</div>';
        } else if (tabName === 'friends') {
            targetScreen.innerHTML = '<div style="text-align:center;padding:50px;color:#888;">دوستان و گروه‌ها به زودی...</div>';
        } else if (tabName === 'rooms') {
            targetScreen.innerHTML = '<div style="text-align:center;padding:50px;color:#888;">روم‌های عمومی به زودی...</div>';
        }
    }
}
