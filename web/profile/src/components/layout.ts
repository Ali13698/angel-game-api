export function renderLayout(content: string, activeTab: string): string {
    return `
    <div class="app-shell">
        <div class="top-header">
            <div class="header-spacer"></div>
            <div class="header-logo">
                <span class="logo-text">Angel</span>
            </div>
            <div class="header-avatar" onclick="window.location.hash='#profile'">
                <img src="https://api.dicebear.com/7.x/bottts/svg?seed=Angel" alt="Profile" />
            </div>
        </div>
        
        <div class="page-content">
            ${content}
        </div>

        <div class="bottom-nav">
            <div class="nav-item ${activeTab === 'home' ? 'active' : ''}" onclick="window.location.hash='#home'">
                <span class="nav-icon">🏠</span>
            </div>
            <div class="nav-item">
                <span class="nav-icon">👥</span>
            </div>
            <div class="nav-center-btn ${activeTab === 'games' ? 'active' : ''}" onclick="window.location.hash='#games'">
                <span class="nav-icon">🎮</span>
            </div>
            <div class="nav-item">
                <span class="nav-icon">💬</span>
            </div>
            <div class="nav-item ${activeTab === 'profile' ? 'active' : ''}" onclick="window.location.hash='#profile'">
                <span class="nav-icon">👤</span>
            </div>
        </div>
    </div>
    `;
}
