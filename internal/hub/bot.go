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

type LudoState struct {
	Board  [16]int  `json:"board"`
	Turn   int      `json:"turn"`
	Seats  []string `json:"seats"`
	Dice   int      `json:"dice"`
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

	tick := time.NewTicker(1200 * time.Millisecond)
	defer tick.Stop()

	for range tick.C {
		st, _ := rm.GetState()
		ws, ok := st.(*wasm.State)
		if !ok { continue }

		if gameID == "tic-tac-toe" {
			var state TicTacToeState
			if err := json.Unmarshal(ws.Data, &state); err != nil { continue }
			if state.Winner != 0 { return }

			botSeat := -1
			for i, s := range state.Seats { if s == string(botID) { botSeat = i; break } }

			if botSeat != -1 && state.Turn == botSeat {
				botP := botSeat + 1; humanP := 1; if botP == 1 { humanP = 2 }
				bestMove := getBestMove(state.Board, botP, humanP)
				if bestMove != -1 {
					rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "move", Payload: map[string]any{"cell": bestMove}}})
				}
			}
		} else if gameID == "ludo" || gameID == "ludo-premium" {
			var state LudoState
			if err := json.Unmarshal(ws.Data, &state); err != nil { continue }
			if state.Winner != 0 { return }

			botSeat := -1
			for i, s := range state.Seats { if s == string(botID) { botSeat = i; break } }

			if botSeat != -1 && state.Turn == botSeat {
				if state.Dice == 0 {
					rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "roll"}})
				} else {
					// تلاش برای حرکت دادن هر ۴ مهره
					for i := 0; i < 4; i++ {
						err := rm.Do(ctx, room.Request{From: botID, Cmd: types.Command{Type: "move", Payload: map[string]any{"token": i}}})
						if err == nil { break }
					}
				}
			}
		}
	}
}
