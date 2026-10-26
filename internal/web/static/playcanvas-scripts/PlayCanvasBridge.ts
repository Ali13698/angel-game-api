// PlayCanvasBridge.ts - Connects the game logic to PlayCanvas Engine safely
import { NetworkClient } from './NetworkClient.js';
import { GameController } from './GameController.js';

export class PlayCanvasBridge {
    public net: NetworkClient;
    public game: GameController;

    constructor() {
        this.net = new NetworkClient();
        this.game = new GameController();
    }

    public async start(initData: string, protoUrl: string) {
        // 1. Load Proto schemas
        await this.net.init(protoUrl);
        await this.game.init(protoUrl);

        // 2. Wire Network to GameController
        this.net.onMessage = (msg) => this.game.handleMessage(msg);
        
        // 3. Connect securely
        // Note: In production, we send initData in the first binary frame for auth.
        // For now, we pass the safely extracted ID to establish connection.
        const safeUserId = "TEMP_ID"; // Will be replaced by Go's verifyInitData response
        this.net.connect(safeUserId);
    }

    // Helper methods for PlayCanvas UI buttons
    public sendJoin(roomId: string) {
        this.net.send('join', new Uint8Array(0), roomId);
    }

    public sendRoll() {
        this.net.send('roll');
    }

    public sendMove(token: number) {
        // PlayCanvas must encode this to Protobuf before sending in strict mode.
        // For now, we pass the intent.
        this.net.send('move', new TextEncoder().encode(JSON.stringify({ token }))); 
    }
}
