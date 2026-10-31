export const Network = {
    connected: false,
    ws: null,
    ClientMsg: null,
    ServerMsg: null,
    currentRev: 0, // 🛡️ ردیابی نسخه برای ضد تقلب
    roomId: 'room_1', // در آینده باید داینامیک شود
    
    async init(onStatusChange) {
        this.onStatus = onStatusChange;
        try {
            const root = await protobuf.load("../../angel.proto");
            this.ClientMsg = root.lookupType("angel.ClientMsg");
            this.ServerMsg = root.lookupType("angel.ServerMsg");
            this.connect();
        } catch (e) {
            console.error("Proto load failed", e);
            this.onStatus("خطا در بارگذاری پروتکل");
        }
    },
    
    connect() {
        const userId = "1"; 
        this.ws = new WebSocket(`wss://game.nomi98.com/ws?user=${userId}`);
        this.ws.binaryType = 'arraybuffer';
        
        this.ws.onopen = () => {
            this.connected = true;
            this.onStatus("متصل شد؛ در حال ورود...");
            // ارسال دستور Join
            this.send('join', { room_id: this.roomId, game_id: 'ludo' });
        };
        
        this.ws.onclose = () => {
            this.connected = false;
            this.onStatus("ارتباط قطع شد.");
        };
        
        this.ws.onmessage = (event) => {
            if (!(event.data instanceof ArrayBuffer)) return;
            try {
                const buffer = new Uint8Array(event.data);
                const decoded = this.ServerMsg.decode(buffer);
                const msg = this.ServerMsg.toObject(decoded, { enums: String, longs: String, bytes: Uint8Array });
                
                // 🛡️ آپدیت نسخه از سرور
                if (msg.rev) {
                    this.currentRev = Number(msg.rev);
                }

                if (msg.error) {
                    this.onStatus("خطا: " + (msg.error.error || 'نامشخص'));
                    return;
                }
                if (msg.events) {
                    for (const ev of msg.events) {
                        if (ev.type === 'ludo.rolled') this.onStatus("تاس ریخته شد");
                    }
                }
            } catch (e) {
                console.error("Decode error", e);
            }
        };
    },
    
    send(type, payload = {}) {
        if (!this.connected || !this.ws || this.ws.readyState !== WebSocket.OPEN) return;
        try {
            const payloadBytes = new TextEncoder().encode(JSON.stringify(payload));
            const message = this.ClientMsg.create({
                type: type,
                id: crypto.randomUUID(),
                roomId: this.roomId,        // 🛡️ ارسال اتاق
                rev: this.currentRev,       // 🛡️ ارسال نسخه برای ضد تقلب
                payload: payloadBytes
            });
            this.ws.send(this.ClientMsg.encode(message).finish());
        } catch (e) {
            console.error("Send error", e);
        }
    }
};
