export interface OnlineUser {
    id: string;
    name: string;
    avatar?: string;
}

function avatarOf(u: OnlineUser): string {
    if (u.avatar) return u.avatar;
    return `https://api.dicebear.com/7.x/avataaars/svg?seed=${encodeURIComponent(u.name || u.id)}`;
}

function escapeHtml(s: string): string {
    return s.replace(/[&<>"']/g, c => ({
        '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'
    }[c] as string));
}

export function renderOnlineUsers(users: OnlineUser[] = []): string {
    if (!users.length) return '';
    return `
    <div class="section-title">آنلاین (${users.length})</div>
    <div class="online-list">
        ${users.map(u => `
            <div class="online-item">
                <img src="${escapeHtml(avatarOf(u))}" class="online-avatar" alt="" />
                <span style="font-size: 12px;">${escapeHtml(u.name)}</span>
            </div>
        `).join('')}
    </div>
    `;
}
