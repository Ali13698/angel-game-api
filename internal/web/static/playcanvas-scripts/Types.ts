// Types.ts - Strict typing for Enterprise architecture (No 'any')

export interface ServerMsg {
    type: string;
    room_id?: string;
    game_id?: string;
    rev: number;
    events: GameEvent[];
    error?: { code: string; error: string };
    state?: Uint8Array;
    from?: string;
    payload?: Uint8Array;
}

export interface GameEvent {
    event_id: string;
    stream_id: string;
    seq: number;
    type: string;
    payload?: Uint8Array;
}

// Payload interfaces matching Go's WASM events
export interface RolledPayload { dice: number; seat: number; }
export interface MovedPayload { token: number; pos: number; seat: number; }
export interface TurnPayload { turn: number; }
