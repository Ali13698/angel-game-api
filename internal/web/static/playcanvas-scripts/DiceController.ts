// DiceController.ts - Listens to server dice rolls and triggers visual updates
import { UIManager } from './UIManager.js';

export class DiceController {
    private ui: UIManager;
    private isRolling: boolean = false;
    
    // DOM Elements (Assuming they exist in HTML)
    private diceBtn: HTMLElement | null;
    private diceDisplay: HTMLElement | null;

    constructor(ui: UIManager) {
        this.ui = ui;
        this.diceBtn = document.getElementById('btn-roll');
        this.diceDisplay = document.getElementById('dice-result');

        this.bindEvents();
    }

    private bindEvents() {
        const bridge = (window as any).AngelBridge?.instance;
        
        // 1. Hook server response (Anti-Cheat: server determines the dice)
        if (bridge) {
            bridge.game.onRolled = (seat: number, dice: number) => {
                this.animateDiceResult(dice);
            };
            
            bridge.game.onTurnChanged = (turn: number) => {
                this.updateTurnUI(turn);
            };
        }

        // 2. Hook button click
        if (this.diceBtn) {
            this.diceBtn.addEventListener('click', () => {
                if (this.isRolling) return;
                this.isRolling = true;
                this.diceBtn?.classList.add('disabled');
                
                // Send intent to Go server
                const bridge = (window as any).AngelBridge;
                if (bridge) bridge.roll();
            });
        }
    }

    private animateDiceResult(result: number) {
        this.isRolling = false;
        this.diceBtn?.classList.remove('disabled');

        if (this.diceDisplay) {
            this.diceDisplay.textContent = result.toString();
        }
        
        // TODO: Add PlayCanvas 3D Dice rotation animation here using `pc.Entity.tween`
    }

    private updateTurnUI(turnSeat: number) {
        const turnText = document.getElementById('turn-indicator');
        if (turnText) {
            turnText.textContent = `نوبت بازیکن ${turnSeat + 1}`;
        }
        
        // Enable button only if it's our turn (assuming seat 0 is us for now)
        if (turnSeat === 0 && this.diceBtn) {
            this.diceBtn.classList.remove('disabled');
        } else if (this.diceBtn) {
            this.diceBtn.classList.add('disabled');
        }
    }
}
