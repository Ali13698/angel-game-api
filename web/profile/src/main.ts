import { User, GameStat } from './types';
import { renderLayout } from './components/layout';
import { renderHomePage } from './pages/home';
import { renderGamesPage } from './pages/games';
import { initProfilePage } from './pages/profile';
import { logError } from './logger';

const appEl = document.getElementById('app');
let currentUser: User | null = null;
let currentStats: GameStat[] = [];

async function loadUser() {
    const p = new URLSearchParams(location.search).get('id');
    if (!p) {
        if (appEl) appEl.innerHTML = '<div style="padding:20px;text-align:center;">No ID provided. Append ?id=YOUR_ID</div>';
        return false;
    }
    try {
        const r = await fetch('/api/profile/' + p);
        const d = await r.json();
        if (!d.ok) throw new Error(d.error || 'User not found');
        currentUser = d.user;
        currentStats = d.stats || [];
        return true;
    } catch (e: any) {
        logError('Load user failed', e);
        if (appEl) appEl.innerHTML = '<div style="padding:20px;text-align:center;color:red;">Error loading user</div>';
        return false;
    }
}

function router() {
    if (!currentUser || !appEl) return;
    const hash = window.location.hash || '#home';

    if (hash === '#home') {
        appEl.innerHTML = renderLayout(renderHomePage(currentUser), 'home');
    } else if (hash === '#games') {
        appEl.innerHTML = renderLayout(renderGamesPage(), 'games');
    } else if (hash === '#profile') {
        initProfilePage(currentUser, currentStats, appEl);
    } else {
        window.location.hash = '#home';
    }
}

window.addEventListener('hashchange', router);

(async () => {
    const ok = await loadUser();
    if (ok) {
        if (!window.location.hash) window.location.hash = '#home';
        else router();
    }
})();
