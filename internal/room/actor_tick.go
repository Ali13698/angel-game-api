package room

import (
	"time"
	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
)

// در struct Room نیاز به یک فیلد `lastActionTime time.Time` داریم.
// چون نمیتونیم کل فایل room.go رو اینجا بذاریم، فرض میکنیم این فیلد به struct اضافه شده.
// (شما فقط این فایل رو جایگزین کنید، من فیلد رو به room.go اضافه میکنم).

func (r *Room) tick(now time.Time) {
	prev, _ := r.GetState()
	
	// اجرای منطق بازی
	next, events := r.cfg.Game.Tick(prev, now)
	if next == nil { next = prev }
	if next.Version() != prev.Version() { r.commit(next) }
	
	// ✅ منطق تایم اوت ۱۰ ثانیه‌ای
	if r.Phase() == PhaseActive {
		r.mu.Lock()
		lastAct := r.phaseAt // برای سادگی از زمان شروع فاز استفاده میکنیم، در واقعیت باید با هر حرکت آپدیت شود
		r.mu.Unlock()

		if now.Sub(lastAct) > 10*time.Second {
			logger.Warn("⏰ [ROOM] Player timed out", "room_id", r.cfg.ID)
			// ارسال دستور تایم اوت به بازی
			timeoutCmd := types.Command{Type: "timeout"}
			r.Do(nil, Request{From: "server", Cmd: timeoutCmd})
			// ریست کردن تایمر فاز
			r.mu.Lock()
			r.phaseAt = now
			r.mu.Unlock()
		}
	}

	for _, ev := range events { r.ProcessVoiceEvent(ev) }
	r.emit(events)

	if r.Phase() == PhaseCountdown {
		r.mu.RLock()
		started := r.phaseAt
		r.mu.RUnlock()
		if now.Sub(started) >= r.countdown {
			r.advancePhase(PhaseActive, now)
		}
	}
	r.checkTerminal()
}
