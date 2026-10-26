// Config.ts - Centralized Configuration (No magic strings)

export const Config = {
    // Network
    WS_URL: 'wss://game.nomi98.com/ws',
    PROTO_URL: '/static/proto/angel.proto',
    
    // Game Settings
    GAME_ID: 'ludo',
    MAX_PLAYERS: 4,
    
    // UI Colors (Matching Plato Dark Theme)
    COLORS: {
        RED: '#e52521',
        GREEN: '#1eae3a',
        BLUE: '#007cee',
        YELLOW: '#fdc010',
        BG_DARK: '#0f1115',
        TEXT_MAIN: '#ffffff'
    },
    
    // Debug Mode (Set to false in production)
    DEBUG: false
};
