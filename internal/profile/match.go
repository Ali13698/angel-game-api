package profile

import (
	"context"
	"log"
	"math"

	"github.com/angel/angel/internal/db"
)

type MatchResult struct {
	UserID int64
	GameID string
	IsWin  bool
	IsDraw bool
}

func RecordMatch(ctx context.Context, res MatchResult) error {
	p := db.Pool()
	if p == nil {
		return errStr("db not init")
	}
	tx, err := p.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	// ۱. اطمینان از وجود رکورد در user_game_stats
	_, err = tx.Exec(ctx, "INSERT INTO user_game_stats (user_id, game_id, wins, losses, draws, games_played, level, xp) VALUES ($1, $2, 0, 0, 0, 0, 1, 0) ON CONFLICT (user_id, game_id) DO NOTHING", res.UserID, res.GameID)
	if err != nil { return err }

	// ۲. آپدیت آمار اختصاصی بازی
	var newGamesPlayed int
	query := "UPDATE user_game_stats SET "
	if res.IsWin { 
		query += "wins = wins + 1, " 
	} else if res.IsDraw { 
		query += "draws = draws + 1, " 
	} else { 
		query += "losses = losses + 1, " 
	}
	query += "games_played = games_played + 1 WHERE user_id = $1 AND game_id = $2 RETURNING games_played"

	err = tx.QueryRow(ctx, query, res.UserID, res.GameID).Scan(&newGamesPlayed)
	if err != nil { return err }

	newLevel := int(math.Floor(math.Sqrt(float64(newGamesPlayed)) * 1.25))
	if newLevel < 1 { newLevel = 1 }

	_, err = tx.Exec(ctx, "UPDATE user_game_stats SET level = $1 WHERE user_id = $2 AND game_id = $3", newLevel, res.UserID, res.GameID)
	if err != nil { return err }

	// ⬇️ ۳. آپدیت جدول اصلی کاربران (این بخش جا افتاده بود) ⬇️
	userQuery := "UPDATE users SET "
	if res.IsWin { 
		userQuery += "wins = wins + 1, coins = coins + 10 " // ۱۰ سکه جایزه برد
	} else if res.IsDraw { 
		userQuery += "coins = coins + 5 "
	} else { 
		userQuery += "losses = losses + 1 " 
	}
	userQuery += "WHERE id = $1"
	
	_, err = tx.Exec(ctx, userQuery, res.UserID)
	if err != nil { return err }

	if err := tx.Commit(ctx); err != nil { return err }
	log.Printf("[STATS] Match recorded: User %d, Game %s, Win: %v, Level: %d", res.UserID, res.GameID, res.IsWin, newLevel)
	return nil
}
