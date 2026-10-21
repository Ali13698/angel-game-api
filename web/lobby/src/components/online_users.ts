export function renderOnlineUsers(): string {
    // لیست کاربران آنلاین (نگاروف حذف شد)
    const users = ['cev', 'aree0x', 'nilo', 'mari__a'];
    
    return `
    <div class="section-title">آنلاین / موارد دلخواه</div>
    <div class="online-list">
        ${users.map(u => `
            <div class="online-item">
                <img src="https://api.dicebear.com/7.x/avataaars/svg?seed=${u}" class="online-avatar" />
                <span style="font-size: 12px;">${u}</span>
            </div>
        `).join('')}
    </div>
    `;
}
