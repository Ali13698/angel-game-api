// Config.ts - Centralized Configuration
export const Config = {
    // Network
    WS_URL: 'wss://game.nomi98.com/ws',
    // ROOT CAUSE FIX: Go serves static at root, not /static/
    PROTO_URL: '/proto/angel.proto', 
    
    // Game Settings
    GAME_ID: 'ludo',
    MAX_PLAYERS: 4,
    
    // UI Colors
    COLORS: {
        RED: '#e52521',
        GREEN: '#1eae3a',
        BLUE: '#007cee',
        YELLOW: '#fdc010',
        BG_DARK: '#0f1115',
        TEXT_MAIN: '#ffffff'
    },
    DEBUG: true
};
