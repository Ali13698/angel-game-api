// GameController.ts - Blind State Machine. Parses binary events, protects against replay attacks.
import { ServerMsg, RolledPayload, MovedPayload, TurnPayload } from './Types.js';
import * as protobuf from 'protobufjs';

export class GameController {
    public currentRev: number = 0;
    public isReady: boolean = false;
    
    private protoRoot: protobuf.Root | null = null;

    // Callbacks for the View (PlayCanvas)
    public onRolled: (seat: number, dice: number) => void = () => {};
    public onMoved: (seat: number, token: number, pos: number) => void = () => {};
    public onTurnChanged: (turn: number) => void = () => {};
    public onGameOver: (winnerSeat: number) => void = () => {};

    public async init(protoUrl: string) {
        this.protoRoot = await protobuf.load(protoUrl);
    }

    public handleMessage(msg: ServerMsg) {
        // Anti-Cheat: Prevent replay attacks and out-of-order processing
        if (this.isReady && msg.rev <= this.currentRev) {
            console.warn(`[GameController] Dropped stale msg. Rev: ${msg.rev} <= ${this.currentRev}`);
            return;
        }

        if (msg.error) {
            console.error(`[Server Error] ${msg.error.code}: ${msg.error.error}`);
            return;
        }

        if (msg.state && msg.state.length > 0) {
            this.isReady = true;
        }

        if (msg.events && msg.events.length > 0) {
            for (const ev of msg.events) {
                this.processEvent(ev.type, ev.payload);
            }
        }

        this.currentRev = msg.rev;
    }

    private processEvent(type: string, payloadBin?: Uint8Array) {
        if (!this.protoRoot || !payloadBin || payloadBin.length === 0) return;

        try {
            if (type === 'ludo.rolled') {
                const P = this.protoRoot.lookupType('angel.RolledPayload');
                const p = P.decode(payloadBin) as unknown as RolledPayload;
                this.onRolled(p.seat, p.dice);
            } else if (type === 'ludo.moved') {
                const P = this.protoRoot.lookupType('angel.MovedPayload');
                const p = P.decode(payloadBin) as unknown as MovedPayload;
                this.onMoved(p.seat, p.token, p.pos);
            } else if (type === 'turn_changed') {
                const P = this.protoRoot.lookupType('angel.TurnPayload');
                const p = P.decode(payloadBin) as unknown as TurnPayload;
                this.onTurnChanged(p.turn);
            }
        } catch (e) {
            console.error(`[GameController] Failed to parse event ${type}`, e);
        }
    }

    public reset() {
        this.currentRev = 0;
        this.isReady = false;
    }
}
