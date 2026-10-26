// Main.ts - Entry point for the PlayCanvas game
import { PlayCanvasBridge } from './PlayCanvasBridge.js';
import { TelegramAuth } from './TelegramAuth.js';

// Global instance for PlayCanvas Scripts to access
let bridge: PlayCanvasBridge | null = null;

export async function bootstrapGame() {
    console.log("[Main] Booting Angel PlayCanvas Client...");

    // 1. Authenticate via Telegram (or fallback)
    const initData = TelegramAuth.getInitData();
    if (!initData) {
        console.error("[Main] Fatal: No initData found. Cannot authenticate.");
        return;
    }

    // 2. Initialize Core Bridge
    bridge = new PlayCanvasBridge();

    // 3. Setup UI Event Listeners (Hooking into PlayCanvas UI)
    bridge.game.onRolled = (seat, dice) => {
        console.log(`[View Update] Player ${seat} rolled a ${dice}`);
        // pc.app.fire('game:rolled', seat, dice);
    };
    
    bridge.game.onMoved = (seat, token, pos) => {
        console.log(`[View Update] Player ${seat} moved token ${token} to ${pos}`);
        // pc.app.fire('game:moved', seat, token, pos);
    };

    bridge.game.onTurnChanged = (turn) => {
        console.log(`[View Update] It is now seat ${turn}'s turn.`);
        // pc.app.fire('game:turn', turn);
    };

    // 4. Start Connection
    try {
        await bridge.start(initData, '/static/proto/angel.proto');
        console.log("[Main] Engine started successfully.");
    } catch (e) {
        console.error("[Main] Failed to start engine:", e);
    }
}

// Ensure global access for PlayCanvas script components
(window as any).AngelBridge = {
    get instance() { return bridge; },
    join: (room: string) => bridge?.sendJoin(room),
    roll: () => bridge?.sendRoll(),
    move: (token: number) => bridge?.sendMove(token)
};

// Auto-start if in browser
if (typeof window !== 'undefined') {
    window.addEventListener('load', () => {
        // Small delay to ensure PlayCanvas is loaded if embedded
        setTimeout(bootstrapGame, 500); 
    });
}
