package hub

import (
	"context"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/room"
)

// runBot: هسته کور. هیچ دانشی از تاس یا مهره ندارد.
// فقط تلاش می‌کند roll کند و اگر شد، مهره‌ها را ۰ تا ۳ امتحان می‌کند.
func (h *Hub) runBot(rm *room.Room, botID types.UserID, grace time.Duration) {
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

	tick := time.NewTicker(1000 * time.Millisecond)
	defer tick.Stop()
	for range tick.C {
		// تلاش برای ریختن تاس
		err := rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "roll"}})
		if err != nil {
			continue // نوبت ربات نیست یا خطای دیگر
		}
		
		// اگر تاس ریخت، تلاش برای حرکت دادن مهره‌های ۰ تا ۳
		for i := 0; i < 4; i++ {
			moveErr := rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{
				Type: "move", Payload: map[string]any{"token": i},
			}})
			if moveErr == nil { break } // حرکت موفق بود
		}
	}
}
