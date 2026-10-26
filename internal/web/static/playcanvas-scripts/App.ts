// App.ts - Master Bootstrapper (Wires everything together)
import { Config } from './Config.js';
import { PlayCanvasStarter } from './PlayCanvasStarter.js';
import { LudoViewScript } from './LudoViewScript.js';
import { UIManager } from './UIManager.js';
import { DiceController } from './DiceController.js';
import { MatchmakerUI } from './MatchmakerUI.js';
import { TelegramAuth } from './TelegramAuth.js';
import { PlayCanvasBridge } from './PlayCanvasBridge.js';

export class App {
    private pcStarter: PlayCanvasStarter;
    private ui: UIManager;
    private bridge: PlayCanvasBridge;

    constructor() {
        this.ui = new UIManager();
        this.pcStarter = new PlayCanvasStarter();
        this.bridge = new PlayCanvasBridge();
    }

    public async boot() {
        console.log("[App] Booting Angel PlayCanvas...");
        
        // 1. Initialize PlayCanvas Engine
        this.pcStarter.init();

        // 2. Expose the bridge globally (CRITICAL FIX)
        (window as any).AngelBridge = {
            instance: this.bridge,
            join: (room: string) => this.bridge.sendJoin(room),
            roll: () => this.bridge.sendRoll(),
            move: (token: number) => this.bridge.sendMove(token)
        };
        
        // 3. Mount 3D Views
        const ludoView = new LudoViewScript(this.pcStarter.app);

        // 4. Initialize UI Components
        const diceCtrl = new DiceController(this.ui);
        const matchmaker = new MatchmakerUI(this.ui);

        // 5. Connect to Backend (Using the correct Bridge start method)
        if (this.bridge) {
            this.bridge.net.onConnected = () => {
                console.log("[App] Connected to server successfully.");
                this.ui.hideLoading();
                matchmaker.show(); // Show menu after connected
            };
            
            this.bridge.net.onError = (err) => {
                this.ui.showError(err);
            };

            // Get user ID and start the full initialization sequence
            const safeId = TelegramAuth.getUserIdForDisplay();
            
            // This calls net.init(), game.init(), and net.connect() properly
            await this.bridge.start(safeId, Config.PROTO_URL); 
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
