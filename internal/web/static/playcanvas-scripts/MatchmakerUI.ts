// MatchmakerUI.ts - Manages the HTML Overlay for Matchmaking (Plato Style)
import { UIManager } from './UIManager.js';

export class MatchmakerUI {
    private ui: UIManager;
    private container: HTMLElement | null;

    constructor(ui: UIManager) {
        this.ui = ui;
        this.container = document.getElementById('matchmaker-ui');
        this.bindEvents();
    }

    private bindEvents() {
        const quickMatchBtn = document.getElementById('btn-quick-match');
        const createRoomBtn = document.getElementById('btn-create-room');

        if (quickMatchBtn) {
            quickMatchBtn.addEventListener('click', () => {
                this.ui.showLoading("در حال یافتن حریف...");
                const bridge = (window as any).AngelBridge;
                if (bridge) bridge.instance.net.send('quick_match');
            });
        }

        if (createRoomBtn) {
            createRoomBtn.addEventListener('click', () => {
                this.ui.showLoading("در حال ساخت اتاق...");
                const bridge = (window as any).AngelBridge;
                if (bridge) bridge.instance.net.send('create_room');
            });
        }
    }

    public show() {
        if (this.container) this.container.style.display = 'flex';
    }

    public hide() {
        if (this.container) this.container.style.display = 'none';
    }
}
