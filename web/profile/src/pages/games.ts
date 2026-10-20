export function renderGamesPage(): string {
    const games = [
        { id: 'match', name: 'Match It', image: 'linear-gradient(45deg, #ff9a9e, #fecfef)' },
        { id: '8ball', name: 'بیلیارد', image: 'linear-gradient(45deg, #a18cd1, #fbc2eb)' },
        { id: 'bowling', name: 'بولینگ', image: 'linear-gradient(45deg, #84fab0, #8fd3f4)' },
        { id: 'ludo', name: 'منچ', image: 'linear-gradient(45deg, #fccb90, #d57eeb)' },
        { id: 'football', name: 'فوتبال', image: 'linear-gradient(45deg, #a1c4fd, #c2e9fb)' },
        { id: 'monster', name: 'مانستر', image: 'linear-gradient(45deg, #ffecd2, #fcb69f)' }
    ];

    return `
        <div class="page-header" style="padding: 10px 0;">
            <h2>بازی ها</h2>
            <div class="header-icons"><span>↑↓</span></div>
        </div>
        <div class="games-tabs">
            <div class="gtab active">همه</div>
            <div class="gtab">جدید</div>
            <div class="gtab">کارت</div>
            <div class="gtab">تخته</div>
        </div>
        <div class="games-grid">
            ${games.map(g => `
                <div class="game-grid-card">
                    <div class="grid-game-bg" style="background: ${g.image};"></div>
                    <div class="grid-game-name">${g.name}</div>
                </div>
            `).join('')}
        </div>
    `;
}
