package wasm

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/tetratelabs/wazero/api"
)

type State struct {
	Rev  uint64          `json:"rev"`
	Data json.RawMessage `json:"data"`
}

func (s *State) Version() uint64 { return s.Rev }

func (s *State) Clone() types.State {
	copyState := *s
	copyState.Data = append(json.RawMessage(nil), s.Data...)
	return &copyState
}

type WasmHandler struct {
	GameID string
	mod    api.Module
	ctx    context.Context
}

// NewHandler instantiates the module ONCE per room. (V2 Architecture)
func NewHandler(ctx context.Context, gameID string) (*WasmHandler, error) {
	mod, err := Instantiate(ctx, gameID)
	if err != nil {
		return nil, err
	}
	return &WasmHandler{GameID: gameID, mod: mod, ctx: ctx}, nil
}

func (h *WasmHandler) ID() types.GameID { return types.GameID(h.GameID) }

func (h *WasmHandler) NewState(players []types.UserID) types.State {
	data, _ := json.Marshal(map[string]any{
		"players": players,
		"seats":   players,
		"rev":     1,
	})
	return &State{Rev: 1, Data: data}
}

type wasmResult struct {
	State  json.RawMessage `json:"state"`
	Events []types.Event   `json:"events"`
}

func (h *WasmHandler) ApplyCommand(state types.State, from types.UserID, cmd types.Command) (types.State, []types.Event, error) {
	current, ok := state.(*State)
	if !ok || current == nil {
		return state, nil, errors.New("wasm: bad state")
	}
	if err := validateCommandVersion(cmd.Version, current.Rev); err != nil {
		return state, nil, err
	}

	ctx, cancel := context.WithTimeout(h.ctx, 2*time.Second)
	defer cancel()

	intent, err := json.Marshal(map[string]any{"from": from, "cmd": cmd})
	if err != nil { return state, nil, err }

	output, err := CallApply(ctx, h.mod, string(current.Data), string(intent),
		GameContext{GameID: h.GameID, Players: []Player{{ID: string(from)}}, ActiveItems: cmd.ActiveItems},
	)
	if err != nil { return state, nil, err }

	var result wasmResult
	if err := json.Unmarshal([]byte(output), &result); err != nil { return state, nil, err }
	if !json.Valid(result.State) || string(result.State) == "null" { return state, nil, errors.New("wasm: missing state") }
	if len(result.Events) > 128 { return state, nil, errors.New("wasm: too many events") }

	return &State{Rev: current.Rev + 1, Data: result.State}, result.Events, nil
}

func (h *WasmHandler) OnPlayerJoined(s types.State, id types.UserID) (types.State, []types.Event) { return s, nil }
func (h *WasmHandler) OnPlayerLeft(s types.State, id types.UserID) (types.State, []types.Event) { return s, nil }

// Tick is crucial for Physics games (Billiards). It will be called 60 times per second ONLY for physics games.
func (h *WasmHandler) Tick(s types.State, now time.Time) (types.State, []types.Event) {
	// V3/V4 will implement the actual tick logic for physics games here.
	return s, nil
}

func (h *WasmHandler) IsTerminal(s types.State) bool { return false }
func (h *WasmHandler) GetResult(s types.State) types.GameResult { return types.GameResult{IsDraw: true} }
func (h *WasmHandler) ViewFor(s types.State, viewer types.UserID) any { return s.(*State).Data }

func (h *WasmHandler) Close(ctx context.Context) {
	if h.mod != nil { h.mod.Close(ctx) }
}

var _ types.GameHandler = (*WasmHandler)(nil)
var _ types.StateFilter = (*WasmHandler)(nil)
