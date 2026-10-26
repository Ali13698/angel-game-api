// Network Layer - Anti-Cheat Bridge for Space Ludo
const Network = {
    mySeat: 1, // Default to Green (1)
    connected: false,
    ws: null,
    ClientMsg: null,
    ServerMsg: null,

    async init() {
        try {
            window.NetworkEvents.onConnection(false, 'در حال اتصال به سرور...');
            const root = await protobuf.load("/proto/angel.proto");
            this.ClientMsg = root.lookupType("angel.ClientMsg");
            this.ServerMsg = root.lookupType("angel.ServerMsg");
            
            // Save payload types for later decoding
            this.payloadTypes = {
                'ludo.rolled': root.lookupType('angel.RolledPayload'),
                'ludo.moved': root.lookupType('angel.MovedPayload'),
                'turn_changed': root.lookupType('angel.TurnPayload')
            };
            
            this.connect();
        } catch (e) {
            console.error("Proto init failed", e);
            window.NetworkEvents.onConnection(false, 'خطا در بارگذاری پروتکل');
        }
    },

    connect() {
        // Hardcoded room and user for testing (In production, use Telegram InitData)
        const roomId = "room_1";
        const userId = "1"; 
        
        this.ws = new WebSocket(`wss://game.nomi98.com/ws?user=${userId}`);
        this.ws.binaryType = 'arraybuffer';

        this.ws.onopen = () => {
            this.connected = true;
            window.NetworkEvents.onConnection(true, 'متصل شد؛ در حال ورود به اتاق...');
            // Send JOIN command immediately
            this.send('join', { room_id: roomId, game_id: 'ludo' });
        };

        this.ws.onclose = () => {
            this.connected = false;
            window.NetworkEvents.onConnection(false, 'ارتباط با سرور قطع شد.');
        };

        this.ws.onerror = () => {
            window.NetworkEvents.onConnection(false, 'خطای شبکه.');
        };

        this.ws.onmessage = (event) => {
            if (!(event.data instanceof ArrayBuffer)) return;
            try {
                const buffer = new Uint8Array(event.data);
                const decoded = this.ServerMsg.decode(buffer);
                const msg = this.ServerMsg.toObject(decoded, { enums: String, longs: String, bytes: Uint8Array });

                if (msg.error) {
                    window.NetworkEvents.onError(msg.error.error || 'خطای سرور');
                    return;
                }

                if (msg.events && msg.events.length > 0) {
                    for (const ev of msg.events) {
                        this.processEvent(ev.type, ev.payload);
                    }
                }
            } catch (e) {
                console.error("Decode error", e);
            }
        };
    },

    processEvent(type, payloadBytes) {
        if (!payloadBytes || payloadBytes.length === 0) return;
        const Type = this.payloadTypes[type];
        if (!Type) return;

        try {
            const p = Type.decode(payloadBytes);
            if (type === 'ludo.rolled') window.NetworkEvents.onRolled(p.seat, p.dice);
            else if (type === 'ludo.moved') window.NetworkEvents.onMoved(p.seat, p.token, p.pos);
            else if (type === 'turn_changed') window.NetworkEvents.onTurnChanged(p.turn);
        } catch (e) {
            console.error("Event decode error", e);
        }
    },

    send(type, payload = {}) {
        if (!this.connected || !this.ws || this.ws.readyState !== WebSocket.OPEN) return;
        try {
            const payloadBytes = new TextEncoder().encode(JSON.stringify(payload));
            const message = this.ClientMsg.create({
                type: type,
                id: crypto.randomUUID(),
                payload: payloadBytes
            });
            const buffer = this.ClientMsg.encode(message).finish();
            this.ws.send(buffer);
        } catch (e) {
            console.error("Send error", e);
        }
    }
};
