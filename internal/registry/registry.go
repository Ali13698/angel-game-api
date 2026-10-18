package registry

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/room"
)

var (
	ErrUnknownGame = errors.New("registry: unknown game")
	ErrBadSpec     = errors.New("registry: bad spec")
	ErrClosed      = errors.New("registry: closed")
)

type Spec struct {
	GameID    types.GameID
	Players   []types.UserID
	Transport room.Transport
	StatsRecorder room.StatsRecorder
}

type entry struct {
	rm     *room.Room
	cancel context.CancelFunc
}

type Registry struct {
	ctx context.Context

	gmu   sync.RWMutex
	games map[types.GameID]types.Game

	mu     sync.RWMutex
	rooms  map[types.RoomID]*entry
	closed bool
}

func New(ctx context.Context) *Registry {
	return &Registry{
		ctx:   ctx,
		games: make(map[types.GameID]types.Game, 8),
		rooms: make(map[types.RoomID]*entry, 32),
	}
}

func (r *Registry) Register(games ...types.Game) error {
	r.gmu.Lock()
	defer r.gmu.Unlock()
	for _, g := range games {
		if g == nil {
			return ErrBadSpec
		}
		id := g.ID()
		if id == "" {
			return ErrBadSpec
		}
		if _, dup := r.games[id]; dup {
			return ErrBadSpec
		}
		r.games[id] = g
	}
	return nil
}

func (r *Registry) Game(id types.GameID) (types.Game, bool) {
	r.gmu.RLock()
	defer r.gmu.RUnlock()
	g, ok := r.games[id]
	return g, ok
}

func (r *Registry) Open(spec Spec) (*room.Room, error) {
	if spec.GameID == "" {
		return nil, ErrBadSpec
	}
	game, ok := r.Game(spec.GameID)
	if !ok {
		return nil, ErrUnknownGame
	}

	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil, ErrClosed
	}
	id := randID()
	rmCtx, cancel := context.WithCancel(r.ctx)
	rm := room.New(room.Config{
		ID:        id,
		GameID:    spec.GameID,
		Game:      game,
		Transport: spec.Transport,
		StatsRecorder: spec.StatsRecorder,
		Players:   spec.Players,
	})
	r.rooms[id] = &entry{rm: rm, cancel: cancel}
	r.mu.Unlock()

	go rm.Run(rmCtx)
	return rm, nil
}

func (r *Registry) Get(id types.RoomID) (*room.Room, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.rooms[id]
	if !ok {
		return nil, false
	}
	return e.rm, true
}

func (r *Registry) Retire(ctx context.Context, id types.RoomID) error {
	r.mu.Lock()
	e, ok := r.rooms[id]
	if ok {
		delete(r.rooms, id)
	}
	r.mu.Unlock()
	if !ok {
		return nil
	}
	e.cancel()
	return e.rm.Close(ctx)
}

func (r *Registry) Close(ctx context.Context) error {
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil
	}
	r.closed = true
	ents := make([]*entry, 0, len(r.rooms))
	for _, e := range r.rooms {
		ents = append(ents, e)
	}
	r.rooms = make(map[types.RoomID]*entry)
	r.mu.Unlock()

	var wg sync.WaitGroup
	for _, e := range ents {
		wg.Add(1)
		go func(e *entry) {
			defer wg.Done()
			e.cancel()
			_ = e.rm.Close(ctx)
		}(e)
	}
	wg.Wait()
	return nil
}

func randID() types.RoomID {
	var b [10]byte
	if _, err := rand.Read(b[:]); err != nil {
		return types.RoomID("room" + time.Now().Format("150405"))
	}
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	return types.RoomID(strings.ToLower(enc.EncodeToString(b[:])))
}
