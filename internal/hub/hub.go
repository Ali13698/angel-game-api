package hub

import (
	"context"
	"net/http"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/registry"
	"github.com/angel/angel/internal/room"
	"github.com/coder/websocket"
)

func (h *Hub) OpenRoom(gameID types.GameID, players []types.UserID) (*room.Room, error) {
	w := &wireRoom{h: h}
	rm, err := h.reg.Open(registry.Spec{
		GameID: gameID, Players: players, Transport: w, StatsRecorder: h.StatsRecorder, EventPublisher: h.EventPublisher,
	})
	if err != nil { return nil, err }

	if len(players) == 1 {
		go h.startBot(rm)
	}

	h.mu.Lock(); h.rooms[rm.ID()] = rm; h.mu.Unlock()
	logger.Info("Room opened", "room_id", rm.ID(), "game_id", gameID)
	return rm, nil
}

// راننده ربات: هم وارد می‌شود هم تاس می‌ریزد
func (h *Hub) startBot(rm *room.Room) {
	time.Sleep(3 * time.Second)
	botID := types.UserID("bot_1")
	
	if err := rm.Join(context.Background(), botID); err != nil {
		logger.Warn("Bot failed to join", "error", err)
		return
	}
	
	// حلقه ساده ربات برای بازی
	for {
		time.Sleep(4 * time.Second)
		st, _ := rm.GetState()
		if st == nil { return }
		
		// اگر نوبت ربات بود، تاس بریز (کد ساده شده)
		// در یک پروژه واقعی، State را پارس می‌کنیم و اگر Turn برابر با botID بود، roll می‌کنیم
		cmd := types.Command{Type: "roll"}
		_ = rm.Do(context.Background(), room.Request{From: botID, Cmd: cmd})
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := types.UserID(r.URL.Query().Get("user"))
	if uid == "" { http.Error(w, "user required", http.StatusUnauthorized); return }
	
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil { return }
	conn.SetReadLimit(maxMessageSize)
	c := &client{hub: h, conn: conn, out: make(chan []byte, 64), id: uid}
	
	h.mu.Lock(); h.clients[c] = struct{}{}; h.mu.Unlock()
	go c.writeLoop()
	c.readLoop()
}
