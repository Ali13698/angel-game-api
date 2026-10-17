import { logWarn } from '../logger';
export function showToast(message: string, type: 'error' | 'success' = 'error') {
    logWarn('Showing Toast:', message);
    const existing = document.getElementById('app-toast');
    if (existing) existing.remove();
    const toast = document.createElement('div');
    toast.id = 'app-toast';
    toast.className = 'toast ' + type;
    toast.innerText = message;
    document.body.appendChild(toast);
    setTimeout(() => toast.remove(), 3000);
}