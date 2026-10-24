package hub

import (
	"context"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/room"
)

// runBot: ربات تخصصی که بر اساس نوع بازی تصمیم می‌گیرد.
func (h *Hub) runBot(rm *room.Room, botID types.UserID, gameID types.GameID, grace time.Duration) {
	ctx := context.Background()
	time.Sleep(grace) // صبر برای ورود انسان

	var joined bool
	for i := 0; i < 10 && !joined; i++ {
		if err := rm.Join(ctx, botID); err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		joined = true
	}
	if !joined { return }

	tick := time.NewTicker(1000 * time.Millisecond)
	defer tick.Stop()
	for range tick.C {
		err := rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "roll"}})
		if err != nil { continue }
		
		if gameID == "ludo" {
			for i := 0; i < 4; i++ {
				moveErr := rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{
					Type: "move", Payload: map[string]any{"token": i},
				}})
				if moveErr == nil { break }
			}
		} else if gameID == "billiards" {
			rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{
				Type: "shoot", Payload: map[string]any{"force": 50, "angle": 45},
			}})
		}
	}
}
