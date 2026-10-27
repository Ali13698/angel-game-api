import { Router } from './router.js';

// Telegram Init
if (window.Telegram && Telegram.WebApp) {
    Telegram.WebApp.ready();
    Telegram.WebApp.expand();
}

// Expose Router globally so HTML onclick attributes can use it
window.Router = Router;

// Initialize app in Home state
document.addEventListener('DOMContentLoaded', () => {
    Router.switch('home', document.getElementById('nav-home'));
});
