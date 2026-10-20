import { User } from '../types';

export function renderHomePage(user: User): string {
    const onlineUsers = ['mari__a', 'nilo', 'aree0x', 'Ngarof', 'cev'];
    const games = [
        { id: 'match', name: 'Match It', image: 'linear-gradient(45deg, #ff9a9e, #fecfef)' },
        { id: '8ball', name: 'بیلیارد', image: 'linear-gradient(45deg, #a18cd1, #fbc2eb)' },
        { id: 'bowling', name: 'بولینگ', image: 'linear-gradient(45deg, #84fab0, #8fd3f4)' },
        { id: 'ludo', name: 'منچ', image: 'linear-gradient(45deg, #fccb90, #d57eeb)' }
    ];

    return `
        <div class="home-tabs">
            <div class="tab active">آخرین</div>
            <div class="tab">اخبار</div>
        </div>
        
        <div class="section-title">آنلاین / موارد دلخواه</div>
        <div class="online-slider">
            ${onlineUsers.map(u => `
                <div class="online-card">
                    <div class="online-avatar">
                        <img src="https://api.dicebear.com/7.x/adventurer/svg?seed=${u}" />
                    </div>
                    <div class="online-name">${u}</div>
                </div>
            `).join('')}
        </div>

        <div class="section-title">بازی ها</div>
        <div class="games-slider">
            ${games.map(g => `
                <div class="game-card" onclick="window.location.hash='#games'">
                    <div class="game-bg" style="background: ${g.image};"></div>
                    <div class="game-info">
                        <span class="game-name">${g.name}</span>
                        <button class="play-btn">بازی</button>
                    </div>
                </div>
            `).join('')}
        </div>

        <div class="daily-reward-section">
            <div class="reward-header">
                <span>جایزه روزانه</span>
                <span class="reward-time">۲ ساعت و ۱۵ دقیقه</span>
            </div>
            <div class="reward-title">برنده ۲ بازی</div>
            <div class="quest-item">
                <div class="quest-icon">8</div>
                <div class="quest-info">
                    <div>۲ اوچو را برنده شوید</div>
                    <div class="progress-bar"><div class="progress-fill" style="width:20%;"></div></div>
                </div>
            </div>
        </div>
    `;
}
