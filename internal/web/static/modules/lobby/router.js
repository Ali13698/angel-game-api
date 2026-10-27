export class Router {
    static switch(tabName, element) {
        // Hide all screens
        document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
        // Deactivate all nav items
        document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));

        const targetScreen = document.getElementById(`screen-${tabName}`);
        
        if (targetScreen) {
            targetScreen.classList.add('active');
        }

        // Special case for home (empty screen)
        if (tabName === 'home') {
            document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
            document.getElementById('nav-home').classList.add('active');
        } else if (element) {
            element.classList.add('active');
        }

        // Trigger specific module loads if needed
        if (tabName === 'profile') {
            import('./ui.js').then(m => m.loadProfile());
        } else if (tabName === 'games') {
            import('./games.js').then(m => m.renderGames());
        }
    }
}
