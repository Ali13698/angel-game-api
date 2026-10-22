// game.js - مدیریت WebSocket و Protobuf
let ws = null;
let AngelProto = null;
let onGameUpdateCallback = null;

async function initWS(userId, onUpdate) {
    onGameUpdateCallback = onUpdate;
    AngelProto = await protobuf.load("/proto/angel.proto");
    
    const wsUrl = `wss://${location.host}/ws?user=${userId}`;
    ws = new WebSocket(wsUrl);
    ws.binaryType = 'arraybuffer';

    ws.onopen = () => console.log('⚡ [WS] Connected to Go Engine');
    
    ws.onmessage = (event) => {
        const ServerMsg = AngelProto.lookupType("angel.ServerMsg");
        const decoded = ServerMsg.decode(new Uint8Array(event.data));
        if (onGameUpdateCallback) onGameUpdateCallback(decoded);
    };
    
    ws.onclose = () => {
        console.log('🔴 [WS] Disconnected');
        ws = null;
    };
}

function sendIntent(type, payload = {}) {
    if (!ws || ws.readyState !== 1 || !AngelProto) return;
    
    const ClientMsg = AngelProto.lookupType("angel.ClientMsg");
    const msg = ClientMsg.encode({
        type: type,
        payload: new TextEncoder().encode(JSON.stringify(payload))
    }).finish();
    
    ws.send(msg);
}
