// Matchmaker.ts - Handles quick matchmaking and private rooms via Protobuf
export class Matchmaker {
    private bridge: any;

    constructor(bridge: any) {
        this.bridge = bridge;
    }

    public startQuickMatch() {
        if (!this.bridge) return;
        console.log("[Matchmaker] Requesting quick match...");
        
        // Send 'quick_match' intent to Go server
        // No payload needed, server uses Redis to pair players
        this.bridge.net.send('quick_match', new Uint8Array(0), '', 0);
    }

    public createPrivateRoom() {
        if (!this.bridge) return;
        console.log("[Matchmaker] Creating private room...");
        
        // Send 'create_room' intent
        this.bridge.net.send('create_room', new Uint8Array(0), '', 0);
    }

    public joinPrivateRoom(roomCode: string) {
        if (!this.bridge) return;
        
        // Construct intent (In a real scenario, this would be encoded via Protobuf)
        const encoder = new TextEncoder();
        const payload = encoder.encode(JSON.stringify({ code: roomCode }));
        
        this.bridge.net.send('join_room', payload, '', 0);
    }
}
