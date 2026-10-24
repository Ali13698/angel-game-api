/* net.js — AngelNet: the ONLY place that builds ClientMsg frames. */
(function (global) {
  'use strict';
  const T = Object.freeze({
    JOIN: 'join', ROLL: 'roll', MOVE: 'move',
    JOINED: 'joined', ERROR: 'error',
    ROLLED: 'ludo.rolled', MOVED: 'ludo.moved', CAPTURED: 'ludo.captured',
    TURN: 'ludo.turn_changed', HOME: 'ludo.home', OVER: 'ludo.game_over'
  });
  const enc = new TextEncoder();
  const dec = new TextDecoder();
  const PROTO_URLS = ['/proto/angel.proto', '/static/proto/angel.proto', '/static/angel.proto'];
  const PROTO_SRC = `syntax = "proto3"; package angel;
    message ClientMsg { string type = 1; string id = 2; string room_id = 3; int64 rev = 4; bytes payload = 5; }
    message GameEvent { string event_id = 1; string stream_id = 2; int64 seq = 3; string type = 4; bytes payload = 5; }
    message ErrorMsg { string code = 1; string error = 2; }
    message ServerMsg { string type = 1; string room_id = 2; string game_id = 3; int64 rev = 4;
      repeated GameEvent events = 5; ErrorMsg error = 6; bytes state = 7; string from = 8; bytes payload = 9; }`;
  async function loadProto(root, url) {
    const list = url ? [url] : PROTO_URLS;
    for (const u of list) { try { return await root.load(u); } catch (e) { /* try next */ } }
    return root.parse(PROTO_SRC).root;
  }
  const uid = () => (global.crypto && crypto.randomUUID)
    ? crypto.randomUUID()
    : Date.now().toString(36) + Math.random().toString(36).slice(2, 8);
  const hasField = (t, n) => !!(t && ((t.fields && t.fields[n]) || (t._fields && t._fields[n])));
  const fname = (t, camel, snake) => (hasField(t, snake) ? snake : camel);
  function parseJsonBytes(bytes) {
    if (bytes == null) return null;
    if (typeof bytes === 'string') { try { return JSON.parse(bytes); } catch { return null; } }
    if (!bytes.length) return null;
    try { return JSON.parse(dec.decode(bytes)); } catch { return null; }
  }
  class AngelNet {
    constructor() {
      this.ws = null; this.proto = null; this.Rev = null;
      this.userId = null; this.roomId = '';
      this._q = []; this._h = new Map();
      this._ClientMsg = null; this._ServerMsg = null;
      this._k = { cRoom: 'roomId', sRoom: 'roomId', sGame: 'gameId' };
    }
    async connect(userId, opts = {}) {
      this.userId = userId; this.roomId = opts.roomId || '';
      this.proto = await loadProto(global.protobuf, opts.protoUrl);
      this._ClientMsg = this.proto.lookupType('angel.ClientMsg');
      this._ServerMsg = this.proto.lookupType('angel.ServerMsg');
      this._k = { cRoom: fname(this._ClientMsg, 'roomId', 'room_id'), sRoom: fname(this._ServerMsg, 'roomId', 'room_id'), sGame: fname(this._ServerMsg, 'gameId', 'game_id') };
      const scheme = location.protocol === 'https:' ? 'wss' : 'ws';
      const url = opts.url || `${scheme}://${location.host}/ws?user=${encodeURIComponent(userId)}`;
      this.ws = new WebSocket(url);
      this.ws.binaryType = 'arraybuffer';
      this.ws.onopen = () => this._flush();
      this.ws.onmessage = (e) => this._onFrame(e.data);
      this.ws.onclose = () => this._emit('closed', null);
      this.ws.onerror = () => this._emit('closed', null);
      return this;
    }
    send(type, payload = {}, extra = {}) {
      if (!type) throw new Error('AngelNet.send: type is required');
      const frame = this._ClientMsg.encode(this._ClientMsg.create({
        type: String(type), id: uid(), [this._k.cRoom]: String(extra.roomId ?? this.roomId ?? ''),
        rev: Number(extra.rev ?? this.Rev ?? 0), payload: enc.encode(JSON.stringify(payload || {}))
      })).finish();
      if (this.ws && this.ws.readyState === 1) this.ws.send(frame); else this._q.push(frame);
      return frame;
    }
    _flush() { const q = this._q; this._q = []; for (const f of q) { try { this.ws.send(f); } catch (e) { this._q.push(f); break; } } }
    join(roomId) { if (roomId) this.roomId = roomId; return this.send(T.JOIN, {}, { roomId: this.roomId }); }
    roll() { return this.send(T.ROLL, {}); }
    move(token) { return this.send(T.MOVE, { token: Number(token) }); }
    _onFrame(buf) {
      let msg;
      try { msg = this._ServerMsg.toObject(this._ServerMsg.decode(new Uint8Array(buf)), { bytes: Uint8Array }); }
      catch (err) { this._emit('badframe', err); return; }
      const room = msg[this._k.sRoom], game = msg[this._k.sGame];
      if (msg.rev != null && Number(msg.rev) > 0) this.Rev = Number(msg.rev);
      if (room) this.roomId = room;
      if (msg.error && (msg.error.code || msg.error.error)) { this._emit(T.ERROR, { code: msg.error.code, message: msg.error.error }); return; }
      if (msg.state && msg.state.length) { const st = parseJsonBytes(msg.state); if (st) this._emit('state', st); }
      const events = [];
      if (msg.events && msg.events.length) events.push(...msg.events);
      else if (msg.event && msg.event.type) events.push(msg.event);
      for (const ev of events) { const payload = parseJsonBytes(ev.payload) || {}; this._emit('event', { type: ev.type, seq: ev.seq, payload }); this._emit(ev.type, payload); }
      if (!events.length && msg.type && msg.type !== T.ERROR) this._emit(msg.type, { rev: msg.rev, roomId: room, gameId: game, payload: parseJsonBytes(msg.payload) });
      this._emit('msg', msg);
    }
    on(type, fn) { if (!this._h.has(type)) this._h.set(type, []); this._h.get(type).push(fn); return () => this.off(type, fn); }
    off(type, fn) { const l = this._h.get(type); if (!l) return; const i = l.indexOf(fn); if (i >= 0) l.splice(i, 1); }
    _emit(type, data) { const list = this._h.get(type); if (list) for (const fn of list.slice()) { try { fn(data); } catch (e) { console.error(e); } } if (type !== '*') this._emit('*', { type, data }); }
  }
  global.AngelNet = new AngelNet();
  global.AngelNet.TYPES = T;
})(window);
