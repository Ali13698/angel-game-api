package main

import (
	"encoding/json"
	"net/http"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/hub"
)

func registerRooms(mux *http.ServeMux, h *hub.Hub) {
	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 64<<10)

		var body struct {
			GameID  string   `json:"game_id"`
			Players []string `json:"players"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if body.GameID == "" || len(body.Players) == 0 {
			http.Error(w, "game_id and players required", http.StatusBadRequest)
			return
		}

		players := make([]types.UserID, len(body.Players))
		for i, player := range body.Players {
			players[i] = types.UserID(player)
		}

		rm, err := h.OpenRoom(types.GameID(body.GameID), players)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"room_id": string(rm.ID()),
		})
	})
}
