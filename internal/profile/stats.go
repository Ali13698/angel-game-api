package profile

import (
	"context"

	"github.com/angel/angel/internal/db"
)

type GameStat struct {
	GameID string `json:"game_id"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Draws  int    `json:"draws"`
	Played int    `json:"games_played"`
	Level  int    `json:"level"`
	XP     int    `json:"xp"`
}

func GetGameStats(ctx context.Context, userID int64) ([]GameStat, error) {
	p := db.Pool()
	if p == nil {
		return nil, errStr("db not init")
	}
	rows, err := p.Query(ctx, `SELECT game_id, COALESCE(wins,0), COALESCE(losses,0), COALESCE(draws,0), COALESCE(games_played,0), COALESCE(level,1), COALESCE(xp,0) FROM user_game_stats WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []GameStat
	for rows.Next() {
		var s GameStat
		if err := rows.Scan(&s.GameID, &s.Wins, &s.Losses, &s.Draws, &s.Played, &s.Level, &s.XP); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	if stats == nil {
		stats = []GameStat{}
	}
	return stats, nil
}
