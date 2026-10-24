package hub

import (
	"context"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/room"
)

func (h *Hub) runBot(rm *room.Room, botID types.UserID, gameID types.GameID, grace time.Duration) {
	ctx := context.Background()
	time.Sleep(grace)

	var joined bool
	for i := 0; i < 10 && !joined; i++ {
		if err := rm.Join(ctx, botID); err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		joined = true
	}
	if !joined { return }
	logger.Info("🤖 Bot joined room", "room_id", rm.ID())

	tick := time.NewTicker(1500 * time.Millisecond)
	defer tick.Stop()
	
	cellToPlay := 0
	for range tick.C {
		if gameID == "ludo" {
			err := rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "roll"}})
			if err == nil {
				for i := 0; i < 4; i++ {
					moveErr := rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "move", Payload: map[string]any{"token": i}}})
					if moveErr == nil { break }
				}
			}
		} else if gameID == "tic-tac-toe" {
			// ربات دوز: تلاش برای پیدا کردن یک خانه خالی و زدن آن
			for i := 0; i < 9; i++ {
				targetCell := (cellToPlay + i) % 9
				err := rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "move", Payload: map[string]any{"cell": targetCell}}})
				if err == nil {
					cellToPlay = (targetCell + 1) % 9
					break
				}
			}
		}
	}
}
