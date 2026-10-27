import { renderProfileScreen } from './ui.js';
import { renderGamesScreen } from './games.js';

export class Router {
    static switch(tabName, element) {
        document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
        document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));

        const targetScreen = document.getElementById(`screen-${tabName}`);
        
        if (targetScreen) targetScreen.classList.add('active');

        if (tabName === 'home') {
            document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
            document.getElementById('nav-home').classList.add('active');
        } else if (element) {
            element.classList.add('active');
        }

        // صدا زدن لایه نمایش مربوطه
        if (tabName === 'profile') {
            renderProfileScreen();
        } else if (tabName === 'games') {
            renderGamesScreen();
        }
    }
}
