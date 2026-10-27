// Network Layer - Anti-Cheat Bridge for Space Ludo
const Network = {
    mySeat: 1, 
    connected: false,
    ws: null,
    ClientMsg: null,
    ServerMsg: null,

    async getFingerprint() {
        try {
            const canvas = document.createElement('canvas');
            const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl');
            const debugInfo = gl ? gl.getExtension('WEBGL_debug_renderer_info') : null;
            const gpu = debugInfo ? gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL) : 'unknown';
            const data = [navigator.userAgent, navigator.language, navigator.hardwareConcurrency || 'unknown', screen.width, screen.height, screen.colorDepth, Intl.DateTimeFormat().resolvedOptions().timeZone, gpu].join('|');
            const msgBuffer = new TextEncoder().encode(data);
            const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer);
            const hashArray = Array.from(new Uint8Array(hashBuffer));
            return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
        } catch (e) {
            return "fallback_device";
        }
    },

    async init() {
        try {
            window.NetworkEvents.onConnection(false, 'در حال اتصال به سرور...');
            const root = await protobuf.load("../../assets/angel.proto");
            this.ClientMsg = root.lookupType("angel.ClientMsg");
            this.ServerMsg = root.lookupType("angel.ServerMsg");
            
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

    async connect() {
        const roomId = "room_1";
        const userId = "1"; 
        
        this.ws = new WebSocket(`wss://game.nomi98.com/ws?user=${userId}`);
        this.ws.binaryType = 'arraybuffer';

        this.ws.onopen = async () => {
            this.connected = true;
            window.NetworkEvents.onConnection(true, 'متصل شد؛ در حال ورود به اتاق...');
            
            const fingerprint = await this.getFingerprint();
            this.send('join', { room_id: roomId, game_id: 'ludo', fingerprint: fingerprint });
        };

        this.ws.onclose = () => {
            this.connected = false;
            window.NetworkEvents.onConnection(false, 'ارتباط با سرور قطع شد.');
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
