package registry

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/angel/angel/internal/core/resource"
	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/room"
)

var (
	ErrUnknownGame = errors.New("registry: unknown game")
	ErrBadSpec     = errors.New("registry: bad spec")
	ErrClosed      = errors.New("registry: closed")
)

type GameFactory interface {
	ID() types.GameID
	New(ctx context.Context) (types.Game, error)
}

type Spec struct {
	GameID         types.GameID
	Players        []types.UserID
	Transport      room.Transport
	StatsRecorder  room.StatsRecorder
	EventPublisher room.EventPublisher
}

type entry struct {
	rm     *room.Room
	cancel context.CancelFunc
	class  resource.Class
}

type Registry struct {
	ctx       context.Context
	gmu       sync.RWMutex
	factories map[types.GameID]GameFactory
	mu        sync.RWMutex
	rooms     map[types.RoomID]*entry
	closed    bool
	resources *resource.Manager
}

func New(ctx context.Context) *Registry {
	return &Registry{
		ctx:       ctx,
		factories: make(map[types.GameID]GameFactory, 8),
		rooms:     make(map[types.RoomID]*entry, 32),
		resources: resource.NewManager(),
	}
}

func (r *Registry) Register(factories ...GameFactory) error {
	r.gmu.Lock()
	defer r.gmu.Unlock()
	for _, f := range factories {
		if f == nil || f.ID() == "" { return ErrBadSpec }
		if _, dup := r.factories[f.ID()]; dup { return ErrBadSpec }
		r.factories[f.ID()] = f
	}
	return nil
}

func (r *Registry) Open(spec Spec) (*room.Room, error) {
	if spec.GameID == "" { return nil, ErrBadSpec }
	factory, ok := r.factories[spec.GameID]
	if !ok { return nil, ErrUnknownGame }

	class := resource.GetClass(string(spec.GameID))
	if err := r.resources.RequestRoom(class); err != nil { return nil, err }

	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		r.resources.ReleaseRoom(class)
		return nil, ErrClosed
	}
	id := randID()
	rmCtx, cancel := context.WithCancel(r.ctx)
	
	gameInstance, err := factory.New(rmCtx)
	if err != nil {
		cancel()
		r.mu.Unlock()
		r.resources.ReleaseRoom(class)
		return nil, err
	}
	
	rm := room.New(room.Config{
		ID: id, GameID: spec.GameID, Game: gameInstance,
		Transport: spec.Transport, StatsRecorder: spec.StatsRecorder,
		EventPublisher: spec.EventPublisher, Players: spec.Players,
	})
	
	r.rooms[id] = &entry{rm: rm, cancel: cancel, class: class}
	r.mu.Unlock()

	go rm.Run(rmCtx)
	return rm, nil
}

func (r *Registry) Retire(ctx context.Context, id types.RoomID) error {
	r.mu.Lock()
	e, ok := r.rooms[id]
	if ok { delete(r.rooms, id) }
	r.mu.Unlock()
	if !ok { return nil }
	
	r.resources.ReleaseRoom(e.class)
	e.cancel()
	return e.rm.Close(ctx)
}

func (r *Registry) Close(ctx context.Context) error {
	r.mu.Lock()
	if r.closed { r.mu.Unlock(); return nil }
	r.closed = true
	ents := make([]*entry, 0, len(r.rooms))
	for _, e := range r.rooms { ents = append(ents, e) }
	r.rooms = make(map[types.RoomID]*entry)
	r.mu.Unlock()

	var wg sync.WaitGroup
	for _, e := range ents {
		wg.Add(1)
		go func(e *entry) {
			defer wg.Done()
			r.resources.ReleaseRoom(e.class)
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
