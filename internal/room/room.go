package room

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
)

type Room struct {
	cfg         Config
	inbox       chan Envelope
	priority    chan Envelope
	tickCh      chan struct{}
	fanoutCh    chan []pendingEvent
	quit        chan struct{}
	done        chan struct{}
	shutdown    chan struct{}
	fanoutDone  chan struct{}
	closeOnce   sync.Once
	shutdownOnce sync.Once
	gate        sync.RWMutex
	closed      bool
	mu          sync.RWMutex
	cur         types.State
	phase       Phase
	phaseAt     time.Time
	countdown   time.Duration
	joinedCount int
	seq         atomic.Uint64
	dropped     atomic.Uint64
	slowClientClosed atomic.Uint64
	backpressure atomic.Uint64
}

func New(cfg Config) *Room {
	if cfg.InboxSize <= 0 {
		cfg.InboxSize = 64
	}
	if cfg.CountdownDuration <= 0 {
		cfg.CountdownDuration = DefaultCountdownDuration
	}
	logger.Info("Room created", "room_id", cfg.ID, "game_id", cfg.GameID)
	return &Room{
		cfg:        cfg,
		inbox:      make(chan Envelope, cfg.InboxSize),
		priority:   make(chan Envelope, PriorityLaneSize),
		tickCh:     make(chan struct{}, TickBufferSize),
		fanoutCh:   make(chan []pendingEvent, FanoutBufferSize),
		quit:       make(chan struct{}),
		done:       make(chan struct{}),
		shutdown:   make(chan struct{}),
		fanoutDone: make(chan struct{}),
		cur:        cfg.Game.NewState(cfg.Players),
		phase:      PhaseLobby,
		phaseAt:    time.Now(),
		countdown:  cfg.CountdownDuration,
	}
}

func (r *Room) ID() types.RoomID     { return r.cfg.ID }
func (r *Room) GameID() types.GameID { return r.cfg.GameID }
func (r *Room) Done() <-chan struct{} { return r.done }
func (r *Room) LastSeq() uint64      { return r.seq.Load() }

func (r *Room) Dropped() uint64          { return r.dropped.Load() }
func (r *Room) SlowClientClosed() uint64 { return r.slowClientClosed.Load() }
func (r *Room) Backpressure() uint64     { return r.backpressure.Load() }

func (r *Room) Phase() Phase {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.phase
}

func (r *Room) GetState() (types.State, uint64) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.cur, r.cur.Version()
}

func (r *Room) Submit(req Request) error {
	return r.enqueueInbox(context.Background(), Envelope{From: req.From, Cmd: req.Cmd}, false)
}

func (r *Room) Do(ctx context.Context, req Request) error {
	reply := make(chan error, 1)
	env := Envelope{From: req.From, Cmd: req.Cmd, ReplyTo: reply}
	if err := r.enqueueInbox(ctx, env, true); err != nil {
		return err
	}
	select {
	case err := <-reply:
		return err
	case <-r.quit:
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Room) enqueueInbox(ctx context.Context, env Envelope, block bool) error {
	r.gate.RLock()
	defer r.gate.RUnlock()
	if r.closed {
		return ErrClosed
	}
	if block {
		select {
		case r.inbox <- env:
			return nil
		case <-r.quit:
			return ErrClosed
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	select {
	case r.inbox <- env:
		return nil
	default:
		return ErrBusy
	}
}

func (r *Room) Close(ctx context.Context) error {
	r.closeOnce.Do(func() {
		r.gate.Lock()
		r.closed = true
		close(r.quit)
		r.gate.Unlock()
		logger.Info("Room closed", "room_id", r.cfg.ID)
	})
	select {
	case <-r.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
