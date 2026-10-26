// BoardMapper.ts - Translates Server Game State (pos) to PlayCanvas 3D Coordinates (X, Y, Z)
export class BoardMapper {
    private readonly CELL_SIZE = 1;
    private readonly BOARD_OFFSET = 7.5; // Center of 15x15 grid

    // Standard Ludo Path (X, Z coordinates for a 15x15 grid)
    // Simplified for the core logic. Replace with exact 52-cell path.
    private trackPositions: { x: number, z: number }[] = [];
    private homePaths: { [seat: number]: { x: number, z: number }[] } = { 0: [], 1: [], 2: [], 3: [] };

    constructor() {
        this.generateTrack();
        this.generateHomePaths();
    }

    private generateTrack() {
        // Generate a dummy circular path for testing
        // In production, hardcode the exact 52 XZ coordinates of the Ludo board
        for (let i = 0; i < 52; i++) {
            const angle = (i / 52) * Math.PI * 2;
            this.trackPositions.push({
                x: Math.cos(angle) * 5.5,
                z: Math.sin(angle) * 5.5
            });
        }
    }

    private generateHomePaths() {
        // Dummy home paths (straight lines to center)
        this.homePaths[0] = [{ x: 0, z: 4 }, { x: 0, z: 3 }, { x: 0, z: 2 }, { x: 0, z: 1 }, { x: 0, z: 0 }];
        this.homePaths[1] = [{ x: 4, z: 0 }, { x: 3, z: 0 }, { x: 2, z: 0 }, { x: 1, z: 0 }, { x: 0, z: 0 }];
        this.homePaths[2] = [{ x: 0, z: -4 }, { x: 0, z: -3 }, { x: 0, z: -2 }, { x: 0, z: -1 }, { x: 0, z: 0 }];
        this.homePaths[3] = [{ x: -4, z: 0 }, { x: -3, z: 0 }, { x: -2, z: 0 }, { x: -1, z: 0 }, { x: 0, z: 0 }];
    }

    /**
     * Returns 3D coordinates for a specific token position.
     * @param seat The player's seat (0-3)
     * @param pos The server's position (-1 = base, 0-51 = track, 52-56 = home, 57 = finish)
     */
    public getCoordinates(seat: number, pos: number): { x: number, y: number, z: number } {
        if (pos === -1) {
            // Base position (outside the board)
            const bases = [{ x: -5, z: -5 }, { x: 5, z: -5 }, { x: 5, z: 5 }, { x: -5, z: 5 }];
            return { x: bases[seat].x, y: 0.5, z: bases[seat].z };
        }
        
        if (pos >= 0 && pos <= 51) {
            const absolutePos = (seat * 13 + pos) % 52; // Standard Ludo offset
            return { x: this.trackPositions[absolutePos].x, y: 0.5, z: this.trackPositions[absolutePos].z };
        }
        
        if (pos >= 52 && pos <= 56) {
            const p = this.homePaths[seat][pos - 52];
            return { x: p.x, y: 0.5, z: p.z };
        }

        return { x: 0, y: 0.5, z: 0 }; // Center (Finish)
    }
}
