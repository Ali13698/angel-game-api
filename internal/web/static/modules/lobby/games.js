export function renderGamesScreen() {
    const container = document.getElementById('screen-games');
    
    // در آینده این لیست می‌تواند از API خوانده شود
    const games = [
        { id: 'ludo', name: 'منچ فضایی', icon: '🎲', route: '/games/ludo/' }
    ];

    container.innerHTML = `
        <div class="games-grid">
            ${games.map(game => `
                <div class="game-card" onclick="window.location.href='${game.route}'">
                    <div class="game-icon">${game.icon}</div>
                    <div class="game-name">${game.name}</div>
                </div>
            `).join('')}
        </div>
    `;
}
