// App.ts - Master Bootstrapper (Wires everything together)
import { Config } from './Config.js';
import { PlayCanvasStarter } from './PlayCanvasStarter.js';
import { LudoViewScript } from './LudoViewScript.js';
import { UIManager } from './UIManager.js';
import { DiceController } from './DiceController.js';
import { MatchmakerUI } from './MatchmakerUI.js';
import { TelegramAuth } from './TelegramAuth.js';

export class App {
    private pcStarter: PlayCanvasStarter;
    private ui: UIManager;

    constructor() {
        this.ui = new UIManager();
        this.pcStarter = new PlayCanvasStarter();
    }

    public async boot() {
        console.log("[App] Booting Angel PlayCanvas...");
        
        // 1. Initialize PlayCanvas Engine
        this.pcStarter.init();
        
        // 2. Mount 3D Views
        const ludoView = new LudoViewScript(this.pcStarter.app);

        // 3. Initialize UI Components
        const diceCtrl = new DiceController(this.ui);
        const matchmaker = new MatchmakerUI(this.ui);

        // 4. Connect to Backend (Assuming AngelBridge is injected globally by NetworkClient)
        const bridge = (window as any).AngelBridge?.instance;
        if (bridge) {
            await bridge.game.init(Config.PROTO_URL);
            
            bridge.net.onConnected = () => {
                this.ui.hideLoading();
                matchmaker.show(); // Show menu after connected
            };
            
            bridge.net.onError = (err) => {
                this.ui.showError(err);
            };

            // For now, use fallback Auth. Production requires sending initData in binary frame.
            const safeId = TelegramAuth.getUserIdForDisplay();
            bridge.net.connect(safeId);
        } else {
            this.ui.showError("Fatal: Network layer missing.");
        }
    }
}

// Start the App
window.addEventListener('load', () => {
    // Small delay to ensure PlayCanvas CDN is ready
    setTimeout(() => {
        const app = new App();
        app.boot();
    }, 500);
});
