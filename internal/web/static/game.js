// game.js - مدیریت WebSocket
let ws = null;
let AngelProto = null;
let onGameUpdateCallback = null;

async function initWS(userId, onUpdate) {
    onGameUpdateCallback = onUpdate;
    if (!AngelProto) {
        AngelProto = await protobuf.load("/proto/angel.proto");
    }
    
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
    if (!ws || ws.readyState !== 1 || !AngelProto) {
        console.log("WS not ready");
        return;
    }
    
    const ClientMsg = AngelProto.lookupType("angel.ClientMsg");
    const msg = ClientMsg.encode({
        type: type,
        payload: new TextEncoder().encode(JSON.stringify(payload))
    }).finish();
    
    ws.send(msg);
    console.log(`[WS] Sent ${type}`, payload);
}

function handleGameMsg(msg) {
    console.log("Received from server:", msg);
    if (msg.type === 'joined') {
        document.getElementById('gameStateView').innerText = "بازی شروع شد!";
    }
    if (msg.type === 'event') {
        const payload = new TextDecoder().decode(msg.event.payload);
        console.log("Game Event:", msg.event.type, payload);
    }
}

function sendChat() {
    const input = document.getElementById('chatInput');
    const text = input.value.trim();
    if (!text) return;
    sendIntent("chat", { text: text });
    document.getElementById('chatMsgs').innerHTML += `<div><b>من:</b> ${text}</div>`;
    input.value = '';
}
