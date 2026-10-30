package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/angel/angel/internal/core/types"
)

type State struct {
	Rev  uint64          `json:"rev"`
	Data json.RawMessage `json:"data"`
}

func (s State) Version() uint64 { return s.Rev }
func (s State) Clone() types.State {
	d := make(json.RawMessage, len(s.Data))
	copy(d, s.Data)
	return State{Rev: s.Rev, Data: d}
}

type WasmHandler struct{ GameID string }

func NewHandler(gameID string) *WasmHandler { return &WasmHandler{GameID: gameID} }
func (h *WasmHandler) ID() string           { return h.GameID }
func (h *WasmHandler) NewState(p []types.User) types.State {
	return State{Rev: 1, Data: json.RawMessage(`{}`)}
}

func (h *WasmHandler) ApplyCommand(state types.State, from types.UserID, cmd types.Command) (types.State, []types.Event, error) {
	current, ok := state.(State)
	if !ok {
		return state, nil, fmt.Errorf("wasm: invalid state type")
	}
	if err := validateCommandVersion(cmd.Version, current.Rev); err != nil {
		return state, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	mod, err := Instantiate(ctx, h.GameID)
	if err != nil {
		return state, nil, err
	}
	defer mod.Close(ctx)

	intent, err := json.Marshal(map[string]any{"from": from, "cmd": cmd})
	if err != nil {
		return state, nil, err
	}

	out, err := CallApply(ctx, mod, string(current.Data), string(intent), GameContext{
		GameID: h.GameID,
		Players: []Player{{ID: string(from)}},
	})
	if err != nil {
		return state, nil, err
	}

	var res struct {
		State  json.RawMessage `json:"state"`
		Events []types.Event   `json:"events"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		return state, nil, err
	}

	return State{Rev: current.Rev + 1, Data: res.State}, res.Events, nil
}

func (h *WasmHandler) OnPlayerJoined(s types.State, id types.UserID) (types.State, []types.Event, error) {
	return s, nil, nil
}
func (h *WasmHandler) OnPlayerLeft(s types.State, id types.UserID) (types.State, []types.Event, error) {
	return s, nil, nil
}
func (h *WasmHandler) Tick(s types.State, now time.Time) (types.State, []types.Event, bool, error) {
	return s, nil, false, nil
}
func (h *WasmHandler) IsTerminal(s types.State) bool { return false }
