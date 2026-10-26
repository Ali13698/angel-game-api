// App.ts - Master Bootstrapper (Wires everything together)
import { Config } from './Config.js';
import { PlayCanvasStarter } from './PlayCanvasStarter.js';
import { LudoViewScript } from './LudoViewScript.js';
import { UIManager } from './UIManager.js';
import { DiceController } from './DiceController.js';
import { MatchmakerUI } from './MatchmakerUI.js';
import { TelegramAuth } from './TelegramAuth.js';
import { PlayCanvasBridge } from './PlayCanvasBridge.js'; // IMPORT ADDED

export class App {
    private pcStarter: PlayCanvasStarter;
    private ui: UIManager;
    private bridge: PlayCanvasBridge; // Reference to the bridge

    constructor() {
        this.ui = new UIManager();
        this.pcStarter = new PlayCanvasStarter();
        this.bridge = new PlayCanvasBridge(); // Instantiate the bridge
    }

    public async boot() {
        console.log("[App] Booting Angel PlayCanvas...");
        
        // 1. Initialize PlayCanvas Engine
        this.pcStarter.init();

        // 2. Expose the bridge globally (CRITICAL FIX)
        // This allows PlayCanvas scripts and UI components to access the network layer
        (window as any).AngelBridge = {
            instance: this.bridge,
            join: (room: string) => this.bridge.sendJoin(room),
            roll: () => this.bridge.sendRoll(),
            move: (token: number) => this.bridge.sendMove(token)
        };
        
        // 3. Mount 3D Views (Now they can find window.AngelBridge)
        const ludoView = new LudoViewScript(this.pcStarter.app);

        // 4. Initialize UI Components
        const diceCtrl = new DiceController(this.ui);
        const matchmaker = new MatchmakerUI(this.ui);

        // 5. Connect to Backend
        if (this.bridge) {
            await this.bridge.game.init(Config.PROTO_URL);
            
            this.bridge.net.onConnected = () => {
                this.ui.hideLoading();
                matchmaker.show(); // Show menu after connected
            };
            
            this.bridge.net.onError = (err) => {
                this.ui.showError(err);
            };

            // For now, use fallback Auth. Production requires sending initData in binary frame.
            const safeId = TelegramAuth.getUserIdForDisplay();
            this.bridge.net.connect(safeId);
        } else {
            this.ui.showError("Fatal: Could not initialize Bridge.");
        }
    }
}

// Start the App
window.addEventListener('load', () => {
    setTimeout(() => {
        const app = new App();
        app.boot();
    }, 500);
});
