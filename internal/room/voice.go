package room

import (
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/wire"
)

const VoiceTurnDuration = 10 * time.Second

func voiceFrame(body map[string]any) []byte {
	return wire.OutboundMsg{Type: "voice", Payload: wire.EncodeValue(body)}.Bytes()
}

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

	r.cfg.Transport.Broadcast(voiceFrame(map[string]any{
		"user_id":  uidStr,
		"status":   "on",
		"duration": VoiceTurnDuration.Seconds(),
	}))

	logger.Info("Voice window opened", "room_id", r.cfg.ID, "user_id", uidStr)

	go func() {
		select {
		case <-time.After(VoiceTurnDuration):
			r.cfg.Transport.Broadcast(voiceFrame(map[string]any{
				"user_id": uidStr,
				"status":  "off",
			}))
			logger.Info("Voice window closed", "room_id", r.cfg.ID, "user_id", uidStr)
		case <-r.quit:
			logger.Info("Voice window aborted: room closed", "room_id", r.cfg.ID)
		}
	}()
}
