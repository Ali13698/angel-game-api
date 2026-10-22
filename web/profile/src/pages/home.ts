import { User } from '../types';

export function renderHomePage(user: User): string {
    const onlineUsers: string[] = []; 
    const games: { name: string, image: string }[] = []; 

    return `
        <div class="home-tabs"><div class="tab active">آخرین</div><div class="tab">اخبار</div></div>
        
        ${onlineUsers.length > 0 ? `
            <div class="section-title">آنلاین / موارد دلخواه</div>
            <div class="online-slider">
                ${onlineUsers.map(u => `<div class="online-card"><div class="online-avatar"><img src="https://api.dicebear.com/7.x/adventurer/svg?seed=${u}" /></div><div class="online-name">${u}</div></div>`).join('')}
            </div>
        ` : ''}

        <div class="section-title">بازی ها</div>
        <div class="games-slider">
            ${games.length > 0 ? games.map(g => `<div class="game-card" onclick="window.location.hash='#games'"><div class="game-bg" style="background: ${g.image};"></div><div class="game-info"><span class="game-name">${g.name}</span><button class="play-btn">بازی</button></div></div>`).join('') : '<div style="color:#8b949e; text-align:center; width:100%; padding:20px; font-size:14px;">هنوز بازی‌ای اضافه نشده است</div>'}
        </div>
    `;
}
