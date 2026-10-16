-- Angel PostgreSQL Schema
-- All timestamps in Unix epoch (seconds)

CREATE TABLE IF NOT EXISTS users (
    id                    BIGINT PRIMARY KEY,
    username              TEXT,
    first_name            TEXT,
    last_name             TEXT,
    photo_url             TEXT,
    avatar                TEXT,
    display_name          TEXT,
    game_username         TEXT UNIQUE,
    game_username_locked  INTEGER DEFAULT 0,
    coins                 INTEGER DEFAULT 0,
    gems                  INTEGER DEFAULT 0,
    wins                  INTEGER DEFAULT 0,
    losses                INTEGER DEFAULT 0,
    level                 INTEGER DEFAULT 1,
    xp                    INTEGER DEFAULT 0,
    profile_color         TEXT DEFAULT '#8b5cf6',
    profile_banner        TEXT DEFAULT 'default',
    name_effect           TEXT,
    name_emoji            TEXT,
    bio                   TEXT,
    banner_id             TEXT,
    frame_id              TEXT,
    hat_id                TEXT,
    chat_bubble_id        TEXT,
    badge_id              TEXT,
    referral_code         TEXT UNIQUE,
    referred_by           BIGINT,
    referrals_count       INTEGER DEFAULT 0,
    referral_rewarded     INTEGER DEFAULT 0,
    streak_days           INTEGER DEFAULT 0,
    last_active_date      TEXT,
    likes_count           INTEGER DEFAULT 0,
    active_days           INTEGER DEFAULT 0,
    banned_until          BIGINT DEFAULT 0,
    created_at            BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    last_seen             BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

CREATE INDEX IF NOT EXISTS idx_users_game_username ON users(game_username);
CREATE INDEX IF NOT EXISTS idx_users_referral_code ON users(referral_code);
CREATE INDEX IF NOT EXISTS idx_users_last_seen ON users(last_seen DESC);

CREATE TABLE IF NOT EXISTS user_game_stats (
    user_id        BIGINT NOT NULL,
    game_id        TEXT NOT NULL,
    wins           INTEGER DEFAULT 0,
    losses         INTEGER DEFAULT 0,
    draws          INTEGER DEFAULT 0,
    games_played   INTEGER DEFAULT 0,
    level          INTEGER DEFAULT 1,
    xp             INTEGER DEFAULT 0,
    PRIMARY KEY (user_id, game_id)
);

CREATE TABLE IF NOT EXISTS game_config (
    game_id      TEXT PRIMARY KEY,
    name         TEXT,
    name_fa      TEXT,
    min_players  INTEGER,
    max_players  INTEGER,
    has_bot      INTEGER DEFAULT 1,
    is_active    INTEGER DEFAULT 1,
    icon         TEXT,
    color        TEXT
);

INSERT INTO game_config (game_id,name,name_fa,min_players,max_players,icon,color) VALUES
    ('ludo','Ludo','منچ',2,4,'🎲','#ef4444'),
    ('backgammon','Backgammon','تخته نرد',2,2,'🎯','#3b82f6'),
    ('dots','Dots','نقطه و خط',2,2,'📐','#10b981'),
    ('football','Football','فوتبال سوکر',2,2,'⚽','#22c55e')
ON CONFLICT (game_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS shop_items (
    item_id      TEXT PRIMARY KEY,
    item_type    TEXT NOT NULL,
    name_fa      TEXT NOT NULL,
    price_coins  INTEGER NOT NULL,
    preview      TEXT,
    meta         TEXT,
    is_active    INTEGER DEFAULT 1,
    is_premium   INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS inventory (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL,
    item_id       TEXT NOT NULL,
    item_type     TEXT NOT NULL,
    equipped      INTEGER DEFAULT 0,
    acquired_at   BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    UNIQUE(user_id, item_id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL,
    amount      INTEGER NOT NULL,
    reason      TEXT,
    meta        TEXT,
    created_at  BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

CREATE TABLE IF NOT EXISTS referrals (
    id           BIGSERIAL PRIMARY KEY,
    referrer_id  BIGINT NOT NULL,
    referred_id  BIGINT NOT NULL UNIQUE,
    created_at   BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

CREATE TABLE IF NOT EXISTS daily_progress (
    user_id         BIGINT NOT NULL,
    date            TEXT NOT NULL,
    games_played    INTEGER DEFAULT 0,
    games_rewarded  INTEGER DEFAULT 0,
    spin_at         BIGINT DEFAULT 0,
    PRIMARY KEY (user_id, date)
);

CREATE TABLE IF NOT EXISTS discount_codes (
    code         TEXT PRIMARY KEY,
    coins        INTEGER NOT NULL,
    max_uses     INTEGER DEFAULT 0,
    used_count   INTEGER DEFAULT 0,
    expires_at   BIGINT,
    created_at   BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

CREATE TABLE IF NOT EXISTS code_redemptions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL,
    code        TEXT NOT NULL,
    created_at  BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    UNIQUE(user_id, code)
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id            BIGSERIAL PRIMARY KEY,
    room_code     TEXT NOT NULL,
    from_user_id  BIGINT NOT NULL,
    to_user_id    BIGINT,
    kind          TEXT DEFAULT 'text',
    cipher        TEXT,
    iv            TEXT,
    created_at    BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

CREATE TABLE IF NOT EXISTS likes (
    id            BIGSERIAL PRIMARY KEY,
    from_user_id  BIGINT NOT NULL,
    to_user_id    BIGINT NOT NULL,
    created_at    BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    UNIQUE(from_user_id, to_user_id)
);

CREATE TABLE IF NOT EXISTS reports (
    id           BIGSERIAL PRIMARY KEY,
    reporter_id  BIGINT NOT NULL,
    target_id    BIGINT NOT NULL,
    room_code    TEXT,
    reason       TEXT,
    created_at   BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

CREATE TABLE IF NOT EXISTS bans (
    user_id             BIGINT PRIMARY KEY,
    level               INTEGER DEFAULT 1,
    reason              TEXT,
    banned_at           BIGINT,
    expires_at          BIGINT,
    is_permanent        INTEGER DEFAULT 0,
    unban_games_needed  INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS blocks (
    from_user_id  BIGINT NOT NULL,
    to_user_id    BIGINT NOT NULL,
    created_at    BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    PRIMARY KEY (from_user_id, to_user_id)
);

CREATE TABLE IF NOT EXISTS friends (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL,
    friend_id   BIGINT NOT NULL,
    status      TEXT DEFAULT 'pending',
    created_at  BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    UNIQUE(user_id, friend_id)
);

CREATE TABLE IF NOT EXISTS matches (
    id          BIGSERIAL PRIMARY KEY,
    room_code   TEXT NOT NULL,
    game_id     TEXT NOT NULL,
    player1     BIGINT NOT NULL,
    player2     BIGINT NOT NULL,
    winner      BIGINT,
    started_at  BIGINT,
    ended_at    BIGINT
);
