package wasm

import (
	"context"
	"fmt"
	"time"

	"github.com/Ali13698/angel-game-api/internal/wire"
)

type Adapter struct {
	loader *Loader
	cache  *InstanceCache
}

func NewAdapter(loader *Loader) *Adapter {
	return &Adapter{
		loader: loader,
		cache:  NewInstanceCache(4),
	}
}

func (a *Adapter) ApplyCommand(ctx context.Context, gameID string, state []byte, cmd wire.Command) ([]byte, []byte, error) {
	inst, err := a.cache.GetOrInstantiate(ctx, a.loader, gameID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get/instantiate game %s: %w", gameID, err)
	}

	current, err := UnmarshalState(state)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	if err := validateCommandVersion(cmd.Version, current.Rev); err != nil {
		return state, nil, err
	}

	// 1-second timeout for game logic execution
	execCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	newState, events, err := inst.Apply(execCtx, state, cmd)
	if err != nil {
		if execCtx.Err() == context.DeadlineExceeded {
			inst.Close() // Evict broken instance
			return nil, nil, fmt.Errorf("game execution timeout (1s exceeded): %w", err)
		}
		return nil, nil, fmt.Errorf("game apply failed: %w", err)
	}

	return newState, events, nil
}

func (a *Adapter) ResetState(ctx context.Context, gameID string, players []string) ([]byte, error) {
	inst, err := a.cache.GetOrInstantiate(ctx, a.loader, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get/instantiate game %s: %w", gameID, err)
	}

	execCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	state, err := inst.Reset(execCtx, players)
	if err != nil {
		if execCtx.Err() == context.DeadlineExceeded {
			inst.Close()
			return nil, fmt.Errorf("game reset timeout: %w", err)
		}
		return nil, fmt.Errorf("game reset failed: %w", err)
	}

	return state, nil
}
