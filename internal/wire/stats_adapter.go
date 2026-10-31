package wire

import (
	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/room"
)

// StatsAdapter یک آداپتور ساده برای اتصال به room.StatsRecorder است
type StatsAdapter struct{}

func NewStatsAdapter() room.StatsRecorder {
	return &StatsAdapter{}
}

func (s *StatsAdapter) RecordGameResult(userID types.UserID, gameID types.GameID, result string) {
	// در آینده به دیتابیس یا کافکا وصل می‌شود
}
