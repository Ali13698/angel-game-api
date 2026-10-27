export function renderGamesScreen() {
    const container = document.getElementById('screen-games');
    
    const games = [
        // ⬇️ تغییر مهم: مسیر به فایل lobby.html تغییر کرد ⬇️
        { id: 'ludo', name: 'منچ فضایی', icon: '🎲', route: '/games/ludo/lobby.html' }
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
