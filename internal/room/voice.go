package room

import (
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/wire"
)

const VoiceTurnDuration = 10 * time.Second

// ProcessVoiceEvent زمان‌سنج صدا را به رویداد نوبت بازی متصل می‌کند
func (r *Room) ProcessVoiceEvent(ev types.Event) {
	if ev.Type != "turn_changed" {
		return
	}

	payload, ok := ev.Payload.(map[string]any)
	if !ok {
		return
	}

	uidStr, ok := payload["user_id"].(string)
	if !ok {
		return
	}

	// ۱. باز کردن میکروفون بازیکن و اعلام ۱۰ ثانیه زمان
	r.cfg.Transport.Broadcast(wire.OutboundMsg{
		Type: "voice",
		Payload: map[string]any{
			"user_id":  uidStr,
			"status":   "on",
			"duration": VoiceTurnDuration.Seconds(),
		},
	}.Bytes())

	logger.Info("Voice window opened", "room_id", r.cfg.ID, "user_id", uidStr)

	// ۲. بستن میکروفون پس از ۱۰ ثانیه (با مدیریت بسته شدن اتاق)
	go func() {
		select {
		case <-time.After(VoiceTurnDuration):
			r.cfg.Transport.Broadcast(wire.OutboundMsg{
				Type: "voice",
				Payload: map[string]any{
					"user_id": uidStr,
					"status":  "off",
				},
			}.Bytes())
			logger.Info("Voice window closed", "room_id", r.cfg.ID, "user_id", uidStr)
		case <-r.quit:
			// اگر اتاق قبل از اتمام ۱۰ ثانیه بسته شد، تایمر متوقف می‌شود
			logger.Info("Voice window aborted: room closed", "room_id", r.cfg.ID)
		}
	}()
}
