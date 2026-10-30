package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/tetratelabs/wazero/api"
)

var executionSlots = make(chan struct{}, 4)

type State struct {
	Rev  uint64          `json:"rev"`
	Data json.RawMessage `json:"data"`
}

func (s State) Version() uint64 {
	return s.Rev
}

func (s State) Clone() types.State {
	dataCopy := make(json.RawMessage, len(s.Data))
	copy(dataCopy, s.Data)
	return State{
		Rev:  s.Rev,
		Data: dataCopy,
	}
}

type WasmHandler struct {
	GameID string
}

func NewHandler(gameID string) *WasmHandler {
	return &WasmHandler{GameID: gameID}
}

func (h *WasmHandler) ID() string {
	return h.GameID
}

func (h *WasmHandler) NewState(players []types.User) types.State {
	return State{
		Rev:  1,
		Data: json.RawMessage(`{}`),
	}
}

func (h *WasmHandler) ApplyCommand(state types.State, from types.UserID, cmd types.Command) (types.State, []types.Event, error) {
	current, ok := state.(State)
	if !ok {
		return state, nil, fmt.Errorf("wasm: invalid state type")
	}

	// Strict Phase 1 revision validation before instantiation
	if err := validateCommandVersion(cmd.Version, current.Rev); err != nil {
		return state, nil, err
	}

	select {
	case executionSlots <- struct{}{}:
		defer func() { <-executionSlots }()
	default:
		return state, nil, fmt.Errorf("wasm: execution capacity reached")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	mod, err := Instantiate(ctx, h.GameID)
	if err != nil {
		return state, nil, fmt.Errorf("failed to instantiate WASM for %s: %w", h.GameID, err)
	}
	defer mod.Close(ctx)

	intentBytes, err := json.Marshal(map[string]any{
		"from": from,
		"cmd":  cmd,
	})
	if err != nil {
		return state, nil, fmt.Errorf("failed to marshal intent: %w", err)
	}

	resStr, err := CallApply(ctx, mod, string(current.Data), string(intentBytes), GameContext{
		GameID:  h.GameID,
		Players: []Player{{ID: string(from)}},
	})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return state, nil, fmt.Errorf("wasm: execution timeout exceeded (1s)")
		}
		return state, nil, fmt.Errorf("wasm apply error: %w", err)
	}

	var result struct {
		State  json.RawMessage `json:"state"`
		Events []types.Event   `json:"events"`
	}
	if err := json.Unmarshal([]byte(resStr), &result); err != nil {
		return state, nil, fmt.Errorf("failed to unmarshal WASM result: %w", err}
	}

	if current.Rev == ^uint64(0) {
		return state, nil, fmt.Errorf("wasm: state revision overflow")
	}

	newState := State{
		Rev:  current.Rev + 1,
		Data: result.State,
	}

	return newState, result.Events, nil
}

func (h *WasmHandler) OnPlayerJoined(state types.State, userID types.UserID) (types.State, []types.Event, error) {
	return state, nil, nil
}

func (h *WasmHandler) OnPlayerLeft(state types.State, userID types.UserID) (types.State, []types.Event, error) {
	return state, nil, nil
}

func (h *WasmHandler) Tick(state types.State, now time.Time) (types.State, []types.Event, bool, error) {
	return state, nil, false, nil
}

func (h *WasmHandler) IsTerminal(state types.State) bool {
	return false
}
