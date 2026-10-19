package wasm

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"github.com/angel/angel/internal/core/types"
)

type State struct {
	Rev  uint64          `json:"rev"`
	Data json.RawMessage `json:"data"`
}

func (s *State) Version() uint64 { return s.Rev }
func (s *State) Clone() types.State {
	cp := *s
	cp.Data = append(json.RawMessage(nil), s.Data...)
	return &cp
}

type WasmHandler struct { GameID string }

func NewHandler(gameID string) *WasmHandler { return &WasmHandler{GameID: gameID} }
func (h *WasmHandler) ID() types.GameID { return types.GameID(h.GameID) }

func (h *WasmHandler) NewState(players []types.UserID) types.State {
	initialData := map[string]any{"players": players, "rev": 1}
	b, _ := json.Marshal(initialData)
	return &State{Rev: 1, Data: b}
}

type wasmResult struct {
	State  json.RawMessage `json:"state"`
	Events []types.Event   `json:"events"`
}

func (h *WasmHandler) ApplyCommand(s types.State, from types.UserID, cmd types.Command) (types.State, []types.Event, error) {
	ws, ok := s.(*State)
	if !ok { return s, nil, errors.New("wasm: bad state") }

	ctx := context.Background()
	mod, err := Instantiate(ctx, h.GameID)
	if err != nil { return s, nil, err }
	defer mod.Close(ctx)

	intentData := map[string]any{"from": from, "cmd": cmd}
	intentB, _ := json.Marshal(intentData)
	ctxData := GameContext{GameID: h.GameID, Players: []Player{{ID: string(from)}}}
	
	outStr, err := CallApply(ctx, mod, string(ws.Data), string(intentB), ctxData)
	if err != nil { return s, nil, err }

	var res wasmResult
	if err := json.Unmarshal([]byte(outStr), &res); err != nil { return s, nil, err }

	next := &State{Rev: ws.Rev + 1, Data: res.State}
	return next, res.Events, nil
}

func (h *WasmHandler) OnPlayerJoined(s types.State, id types.UserID) (types.State, []types.Event) { return s, nil }
func (h *WasmHandler) OnPlayerLeft(s types.State, id types.UserID) (types.State, []types.Event) { return s, nil }
func (h *WasmHandler) Tick(s types.State, now time.Time) (types.State, []types.Event) { return s, nil }
func (h *WasmHandler) IsTerminal(s types.State) bool { return false }
func (h *WasmHandler) GetResult(s types.State) types.GameResult { return types.GameResult{IsDraw: true} }
func (h *WasmHandler) ViewFor(s types.State, viewer types.UserID) any { return s.(*State).Data }

var _ types.GameHandler = (*WasmHandler)(nil)
var _ types.StateFilter = (*WasmHandler)(nil)
