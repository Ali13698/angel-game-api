package room

import (
	"context"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
)

// handleGameEnd وظیفه ثبت آمار و ارسال رویداد به کافکا را بر عهده دارد
func (r *Room) handleGameEnd(st types.State) {
	r.advancePhase(PhaseFinished, time.Now())
	res := r.cfg.Game.GetResult(st)

	// ۱. ثبت آمار در دیتابیس
	if r.cfg.StatsRecorder != nil {
		go func() {
			ctx := context.Background()
			if res.IsDraw {
				_ = r.cfg.StatsRecorder.RecordMatch(ctx, res.WinnerID, string(r.cfg.GameID), false, true)
				_ = r.cfg.StatsRecorder.RecordMatch(ctx, res.LoserID, string(r.cfg.GameID), false, true)
			} else {
				_ = r.cfg.StatsRecorder.RecordMatch(ctx, res.WinnerID, string(r.cfg.GameID), true, false)
				_ = r.cfg.StatsRecorder.RecordMatch(ctx, res.LoserID, string(r.cfg.GameID), false, false)
			}
		}()
	}

	// ۲. ارسال رویداد به کافکا برای تحلیل داده‌ها و Replay
	if r.cfg.EventPublisher != nil {
		go func() {
			payload := map[string]any{
				"room_id":   string(r.cfg.ID),
				"game_id":   string(r.cfg.GameID),
				"winner_id": string(res.WinnerID),
				"loser_id":  string(res.LoserID),
				"is_draw":   res.IsDraw,
				"timestamp": time.Now().Unix(),
			}
			err := r.cfg.EventPublisher.Publish(context.Background(), "game.ended", string(r.cfg.ID), payload)
			if err != nil {
				logger.Error("Failed to publish game.ended to Kafka", "room_id", r.cfg.ID, "error", err)
			} else {
				logger.Info("Event published to Kafka", "topic", "game.ended", "room_id", r.cfg.ID)
			}
		}()
	}
}
