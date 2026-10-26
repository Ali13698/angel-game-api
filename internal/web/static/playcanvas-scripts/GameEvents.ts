// GameEvents.ts - Centralized Event Emitter for PlayCanvas (Decouples logic from UI)
export class GameEvents {
    private static listeners: Map<string, Function[]> = new Map();

    public static on(event: string, callback: Function) {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, []);
        }
        this.listeners.get(event)!.push(callback);
    }

    public static emit(event: string, ...args: any[]) {
        const callbacks = this.listeners.get(event);
        if (callbacks) {
            callbacks.forEach(cb => {
                try {
                    cb(...args);
                } catch (e) {
                    console.error(`[GameEvents] Error in listener for ${event}`, e);
                }
            });
        }
    }

    public static clear() {
        this.listeners.clear();
    }
}
