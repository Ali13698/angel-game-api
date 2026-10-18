package wire

import (
	"context"
	"strconv"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/profile"
)

type StatsAdapter struct{}

func (s *StatsAdapter) RecordMatch(ctx context.Context, userID types.UserID, gameID string, isWin bool, isDraw bool) error {
	dbID, err := strconv.ParseInt(string(userID), 10, 64)
	if err != nil { return err }
	return profile.RecordMatch(ctx, profile.MatchResult{
		UserID: dbID, GameID: gameID, IsWin: isWin, IsDraw: isDraw,
	})
}

func NewStatsAdapter() *StatsAdapter {
	return &StatsAdapter{}
}