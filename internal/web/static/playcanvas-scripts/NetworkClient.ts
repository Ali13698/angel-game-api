import * as protobuf from 'protobufjs';

/**
 * Enterprise-grade Network Client for PlayCanvas.
 * Handles binary Protobuf communication with the Go backend.
 */
export class NetworkClient {
    private ws: WebSocket | null = null;
    private protoRoot: protobuf.Root | null = null;
    private ClientMsg: protobuf.Type | null = null;
    private ServerMsg: protobuf.Type | null = null;
    
    private url: string = '';
    private isConnected: boolean = false;

    // Events (Callbacks) for PlayCanvas to subscribe to
    public onConnected: () => void = () => {};
    public onDisconnected: () => void = () => {};
    public onError: (err: string) => void = () => {};
    public onMessage: (msg: any) => void = () => {};

    /**
     * Initializes the Protobuf schema.
     * @param protoUrl The URL where angel.proto is hosted.
     */
    public async init(protoUrl: string): Promise<void> {
        try {
            this.protoRoot = await protobuf.load(protoUrl);
            this.ClientMsg = this.protoRoot.lookupType('angel.ClientMsg');
            this.ServerMsg = this.protoRoot.lookupType('angel.ServerMsg');
        } catch (error) {
            this.onError(`Failed to load Proto schema: ${error}`);
            throw error;
        }
    }

    /**
     * Connects to the Go backend securely.
     * @param userId The authenticated User ID.
     */
    public connect(userId: string): void {
        if (!this.ClientMsg || !this.ServerMsg) {
            this.onError("Protobuf not initialized. Call init() first.");
            return;
        }
        
        // Root cause fix: Escape userId to prevent URL injection
        const safeUserId = encodeURIComponent(userId);
        this.url = `wss://game.nomi98.com/ws?user=${safeUserId}`;
        
        this.ws = new WebSocket(this.url);
        this.ws.binaryType = 'arraybuffer';

        this.ws.onopen = () => {
            this.isConnected = true;
            this.onConnected();
        };

        this.ws.onclose = () => {
            this.isConnected = false;
            this.onDisconnected();
        };

        this.ws.onerror = (event) => {
            this.onError("WebSocket error occurred.");
        };

        this.ws.onmessage = (event) => {
            if (!(event.data instanceof ArrayBuffer)) return;
            try {
                const buffer = new Uint8Array(event.data);
                // Decode binary Protobuf directly
                const decoded = this.ServerMsg!.decode(buffer);
                const obj = this.ServerMsg!.toObject(decoded, {
                    enums: String,
                    longs: String,
                    bytes: Uint8Array // Keep payloads as Uint8Array for performance
                });
                this.onMessage(obj);
            } catch (err) {
                this.onError(`Failed to decode ServerMsg: ${err}`);
            }
        };
    }

    /**
     * Sends a command to the Go server.
     * @param type The command type (e.g., 'join', 'move').
     * @param payloadBytes The binary Protobuf payload for the specific command.
     * @param roomId The room ID.
     * @param rev The current authoritative revision.
     */
    public send(type: string, payloadBytes: Uint8Array = new Uint8Array(0), roomId: string = '', rev: number = 0): void {
        if (!this.isConnected || !this.ws || !this.ClientMsg) return;

        try {
            const message = this.ClientMsg.create({
                type: type,
                id: crypto.randomUUID(), // Unique ID for tracking
                room_id: roomId,
                rev: rev,
                payload: payloadBytes
            });

            const buffer = this.ClientMsg.encode(message).finish();
            this.ws.send(buffer);
        } catch (err) {
            this.onError(`Failed to encode ClientMsg: ${err}`);
        }
    }

    public disconnect(): void {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
            this.isConnected = false;
        }
    }
}
