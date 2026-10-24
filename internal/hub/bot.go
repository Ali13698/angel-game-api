package hub

import (
	"context"
	"encoding/json"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/room"
	"github.com/angel/angel/internal/wasm"
)

type TicTacToeState struct {
	Board  [9]int   `json:"board"`
	Turn   int      `json:"turn"`
	Seats  []string `json:"seats"`
	Winner int      `json:"winner"`
}

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

	tick := time.NewTicker(1500 * time.Millisecond)
	defer tick.Stop()

	for range tick.C {
		if gameID == "ludo" {
			// منطق منچ (ساده)
			err := rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "roll"}})
			if err == nil {
				for i := 0; i < 4; i++ {
					moveErr := rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "move", Payload: map[string]any{"token": i}}})
					if moveErr == nil { break }
				}
			}
		} else if gameID == "tic-tac-toe" {
			// ✅ منطق هوشمند دوز
			st, _ := rm.GetState()
			ws, ok := st.(*wasm.State)
			if !ok { continue }

			var state TicTacToeState
			if err := json.Unmarshal(ws.Data, &state); err != nil { continue }
			
			if state.Winner != 0 { return } // بازی تمام شده

			// پیدا کردن سیت ربات
			botSeat := -1
			for i, s := range state.Seats {
				if s == string(botID) { botSeat = i; break }
			}

			// اگر نوبت ربات است، حرکت کن
			if botSeat != -1 && state.Turn == botSeat {
				botP := botSeat + 1
				humanP := 1
				if botP == 1 { humanP = 2 }

				bestMove := getBestMove(state.Board, botP, humanP)
				if bestMove != -1 {
					logger.Info("🤖 Bot making smart move", "cell", bestMove)
					rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "move", Payload: map[string]any{"cell": bestMove}}})
				}
			}
		}
	}
}
