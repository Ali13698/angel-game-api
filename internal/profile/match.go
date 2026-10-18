package profile

import (
	"context"
	"log"
	"math"

	"github.com/angel/angel/internal/db"
)

// MatchResult داده‌ای است که موتور بازی (Room) پس از پایان بازی تولید می‌کند.
type MatchResult struct {
	UserID int64
	GameID string
	IsWin  bool
	IsDraw bool
}

// RecordMatch تابعی داخلی (In-Memory) برای ثبت نتیجه بازی است.
// این تابع به هیچ عنوان از طریق HTTP قابل دسترسی نیست و فقط موتور بازی می‌تواند آن را صدا بزند.
func RecordMatch(ctx context.Context, res MatchResult) error {
	p := db.Pool()
	if p == nil {
		return errStr("db not init")
	}

	tx, err := p.Begin(ctx)
	if err != nil {
		log.Printf("[STATS] Tx begin failed: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	// ۱. اطمینان از وجود رکورد برای این کاربر و بازی
	_, err = tx.Exec(ctx, `
		INSERT INTO user_game_stats (user_id, game_id, wins, losses, draws, games_played, level, xp)
		VALUES ($1, $2, 0, 0, 0, 0, 1, 0)
		ON CONFLICT (user_id, game_id) DO NOTHING
	`, res.UserID, res.GameID)
	if err != nil {
		log.Printf("[STATS] Insert default failed: %v", err)
		return err
	}

	// ۲. آپدیت آمار (برد، باخت یا مساوی)
	var newGamesPlayed int
	query := `UPDATE user_game_stats SET `
	if res.IsWin {
		query += `wins = wins + 1, `
	} else if res.IsDraw {
		query += `draws = draws + 1, `
	} else {
		query += `losses = losses + 1, `
	}
	query += `games_played = games_played + 1 WHERE user_id = $1 AND game_id = $2 RETURNING games_played`

	err = tx.QueryRow(ctx, query, res.UserID, res.GameID).Scan(&newGamesPlayed)
	if err != nil {
		log.Printf("[STATS] Update stats failed: %v", err)
		return err
	}

	// ۳. محاسبه لِوِل با فرمول دقیق پلاتو: Level = floor(sqrt(games_played) * 1.25)
	newLevel := int(math.Floor(math.Sqrt(float64(newGamesPlayed)) * 1.25))
	if newLevel < 1 {
		newLevel = 1
	}

	// ۴. ثبت لِوِل جدید
	_, err = tx.Exec(ctx, `UPDATE user_game_stats SET level = $1 WHERE user_id = $2 AND game_id = $3`, newLevel, res.UserID, res.GameID)
	if err != nil {
		log.Printf("[STATS] Update level failed: %v", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("[STATS] Tx commit failed: %v", err)
		return err
	}

	log.Printf("[STATS] Match recorded successfully: User %d, Game %s, Win: %v, New Level: %d", res.UserID, res.GameID, res.IsWin, newLevel)
	return nil
}
