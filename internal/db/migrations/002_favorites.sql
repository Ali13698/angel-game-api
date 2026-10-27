-- 002_favorites.sql
CREATE TABLE IF NOT EXISTS user_favorites (
    user_id BIGINT NOT NULL,
    game_id TEXT NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    PRIMARY KEY (user_id, game_id)
);
