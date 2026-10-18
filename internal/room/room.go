package room

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/wire"
)

var (
	ErrClosed       = errors.New("room: closed")
	ErrBusy         = errors.New("room: inbox is full")
	ErrWrongPhase   = errors.New("room: command not accepted in this phase")
	ErrStaleVersion = errors.New("room: stale command version")
	ErrPriorityFull = errors.New("room: priority lane full")
)

type Phase int

const (
	PhaseLobby Phase = iota
	PhaseCountdown
	PhaseActive
	PhaseFinished
	PhaseAborted
)

func (p Phase) String() string {
	switch p {
	case PhaseLobby:
		return "LOBBY"
	case PhaseCountdown:
		return "COUNTDOWN"
	case PhaseActive:
		return "ACTIVE"
	case PhaseFinished:
		return "FINISHED"
	case PhaseAborted:
		return "ABORTED"
	}
	return "UNKNOWN"
}

const DefaultCountdownDuration = 3 * time.Second
const PriorityLaneSize = 64
const PriorityLaneWaitTimeout = 100 * time.Millisecond

// TickBufferSize دقیقاً ۱ (نه ۰، نه ۲+).
// مرجع: CON-2026-1006-12 بند «ه» + EXE-2026-1006-16 بند ۲.۲.
const TickBufferSize = 1

// FanoutBufferSize ظرفیت کانال Fanout.
// سرریز → صعود کنترلشده per-audience (EXE-2026-1006-17 بند ۳.۱).
const FanoutBufferSize = 256

// Transport رابط لبهٔ ارسال به کلاینت‌ها.
//
// CloseSlowClient: در سرریز ToPlayer فراخوانی می‌شود.
// SignalBackpressure: در سرریز ToRoom فراخوانی می‌شود.
type Transport interface {
	Broadcast(frame []byte)
	SendTo(id types.UserID, frame []byte)
	CloseSlowClient(id types.UserID, code string, reason string)
	SignalBackpressure(code string, reason string)
}

type StatsRecorder interface {
	RecordMatch(ctx context.Context, userID types.UserID, gameID string, isWin bool, isDraw bool) error
}

type Config struct {
	ID                types.RoomID
	GameID            types.GameID
	Game              types.Game
	Transport         Transport
	Players           []types.UserID
	InboxSize         int
	CountdownDuration time.Duration
	StatsRecorder     StatsRecorder
}

type Request struct {
	From types.UserID
	Cmd  types.Command
}

type Envelope struct {
	From    types.UserID
	Cmd     types.Command
	Leave   bool
	ReplyTo chan<- error
}

type pendingEvent struct {
	seq      uint64
	event    types.Event
	streamID string
	gameID   string
	rev      uint64
}

type Room struct {
	cfg       Config
	inbox     chan Envelope
	priority  chan Envelope
	tickCh    chan struct{}
	fanoutCh  chan []pendingEvent
	quit      chan struct{}
	done      chan struct{}
	shutdown  chan struct{}
	fanoutDone chan struct{}
	closeOnce  sync.Once
	shutdownOnce sync.Once

	gate   sync.RWMutex
	closed bool

	mu  sync.RWMutex
	cur types.State

	phase     Phase
	phaseAt   time.Time
	countdown time.Duration
	joinedCount int

	seq atomic.Uint64

	// شمارنده‌های سرریز fanout (EXE-2026-1006-17 بند ۳.۱-ب).
	dropped          atomic.Uint64 // فقط Internal + ToSpectators
	slowClientClosed atomic.Uint64 // ToPlayer overflow
	backpressure     atomic.Uint64 // ToRoom overflow
}

func New(cfg Config) *Room {
	if cfg.InboxSize <= 0 {
		cfg.InboxSize = 64
	}
	if cfg.CountdownDuration <= 0 {
		cfg.CountdownDuration = DefaultCountdownDuration
	}
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

func (r *Room) ID() types.RoomID      { return r.cfg.ID }
func (r *Room) GameID() types.GameID  { return r.cfg.GameID }
func (r *Room) Done() <-chan struct{} { return r.done }
func (r *Room) LastSeq() uint64       { return r.seq.Load() }

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

// Leave — سهمرحله‌ای. مرجع: docs/LEAVE-SYNC.md
func (r *Room) Join(ctx context.Context, id types.UserID) error {
	r.mu.Lock(); prev := r.cur; r.mu.Unlock()
	if prev == nil { return ErrClosed }
	next, events := r.cfg.Game.OnPlayerJoined(prev, id)
	if next != nil && next.Version() != prev.Version() { r.commit(next) }
	r.emit(events)
	r.mu.Lock(); r.joinedCount++; jc := r.joinedCount; ph := r.phase; r.mu.Unlock()
	if jc >= 2 && ph == PhaseLobby { r.advancePhase(PhaseCountdown, time.Now()) }
	return nil
}

func (r *Room) Leave(ctx context.Context, id types.UserID) error {
	reply := make(chan error, 1)
	env := Envelope{From: id, Leave: true, ReplyTo: reply}

	select {
	case r.priority <- env:
	default:
		waitCtx, cancel := context.WithTimeout(ctx, PriorityLaneWaitTimeout)
		defer cancel()
		select {
		case r.priority <- env:
		case <-r.quit:
			return ErrClosed
		case <-waitCtx.Done():
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return ErrPriorityFull
		}
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
	})
	select {
	case <-r.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Room) enqueueTick() {
	select {
	case r.tickCh <- struct{}{}:
	default:
	}
}

func (r *Room) Run(ctx context.Context) error {
	defer close(r.done)

	go r.tickerLoop(ctx)

	var fanoutWG sync.WaitGroup
	fanoutWG.Add(1)
	go func() {
		defer fanoutWG.Done()
		defer close(r.fanoutDone)
		r.fanoutLoop(ctx)
	}()

	r.advancePhase(PhaseLobby, time.Now())

	err := r.actorLoop(ctx)

	// سیگنال به ticker و fanout که تمام شد.
	r.shutdownOnce.Do(func() { close(r.shutdown) })

	// صبر برای تخلیهٔ کامل fanout (EXE-2026-1006-17 بند ۳.۴).
	fanoutWG.Wait()

	return err
}

func (r *Room) tickerLoop(ctx context.Context) {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-r.quit:
			return
		case <-r.shutdown:
			return
		case <-t.C:
			r.enqueueTick()
		}
	}
}

func (r *Room) actorLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			r.advancePhase(PhaseAborted, time.Now())
			return ctx.Err()
		case <-r.quit:
			r.advancePhase(PhaseAborted, time.Now())
			return nil
		case env := <-r.priority:
			r.step(env)
		case env := <-r.inbox:
			r.step(env)
		case <-r.tickCh:
			r.tick(time.Now())
		}
	}
}

// fanoutLoop کار advisory را روی goroutine جدا انجام می‌دهد.
// در shutdown، قبل از خروج، بافر را تخلیه می‌کند (بند ۳.۴).
func (r *Room) fanoutLoop(ctx context.Context) {
	for {
		select {
		case <-r.shutdown:
			r.drainFanout()
			return
		case <-ctx.Done():
			r.drainFanout()
			return
		case <-r.quit:
			r.drainFanout()
			return
		case pending := <-r.fanoutCh:
			r.deliver(pending)
		}
	}
}

func (r *Room) drainFanout() {
	for {
		select {
		case pending := <-r.fanoutCh:
			r.deliver(pending)
		default:
			return
		}
	}
}

func (r *Room) step(env Envelope) {
	prev, _ := r.GetState()
	if env.Leave {
		next, events := r.cfg.Game.OnPlayerLeft(prev, env.From)
		if next == nil {
			next = prev
		}
		if next.Version() != prev.Version() {
			r.commit(next)
		}
		r.emit(events)
		r.reply(env, nil)
		r.checkTerminal()
		return
	}

	if r.Phase() != PhaseActive {
		r.reply(env, ErrWrongPhase)
		return
	}

	if env.Cmd.Version > 0 && env.Cmd.Version < prev.Version() {
		r.reply(env, ErrStaleVersion)
		return
	}

	next, events, err := r.cfg.Game.ApplyCommand(prev, env.From, env.Cmd)
	if err != nil {
		r.reply(env, err)
		return
	}
	if next == nil {
		next = prev
	}
	if next.Version() != prev.Version() {
		r.commit(next)
	}
	r.emit(events)
	r.reply(env, nil)
	r.checkTerminal()
}

func (r *Room) tick(now time.Time) {
	prev, _ := r.GetState()
	next, events := r.cfg.Game.Tick(prev, now)
	if next == nil {
		next = prev
	}
	if next.Version() != prev.Version() {
		r.commit(next)
	}
	r.emit(events)

	if r.Phase() == PhaseCountdown {
		r.mu.RLock()
		started := r.phaseAt
		r.mu.RUnlock()
		if now.Sub(started) >= r.countdown {
			r.advancePhase(PhaseActive, now)
		}
	}

	r.checkTerminal()
}

func (r *Room) advancePhase(next Phase, now time.Time) {
	r.mu.Lock()
	cur := r.phase
	if cur == next {
		r.mu.Unlock()
		return
	}
	allowed := false
	switch cur {
	case PhaseLobby:
		allowed = next == PhaseCountdown || next == PhaseAborted
	case PhaseCountdown:
		allowed = next == PhaseActive || next == PhaseAborted
	case PhaseActive:
		allowed = next == PhaseFinished || next == PhaseAborted
	case PhaseFinished:
		allowed = next == PhaseAborted
	case PhaseAborted:
		allowed = false
	}
	if !allowed {
		r.mu.Unlock()
		return
	}
	r.phase = next
	r.phaseAt = now
	r.mu.Unlock()
}

func (r *Room) checkTerminal() {
	if r.Phase() != PhaseActive {
		return
	}
	st, _ := r.GetState()
	if r.cfg.Game.IsTerminal(st) {
		r.advancePhase(PhaseFinished, time.Now())
		if r.cfg.StatsRecorder != nil {
			res := r.cfg.Game.GetResult(st)
			go func() {
				ctx := context.Background()
				if res.IsDraw {
					_ = r.cfg.StatsRecorder.RecordMatch(ctx, res.WinnerID, string(r.cfg.GameID), false, true)
					_ = r.cfg.StatsRecorder.RecordMatch(ctx, res.LoserID, string(r.cfg.GameID), false, true)
				} else {
					_ = r.cfg.StatsRecorder.RecordMatch(ctx, res.WinnerID, string(r.cfg.GameID), true, false)
					_ = r.cfg.StatsRecorder.RecordMatch(ctx, res.LoserID, string(r.cfg.GameID), false, false)
				}
			}()
		}
	}
}

func (r *Room) commit(next types.State) {
	r.mu.Lock()
	r.cur = next
	r.mu.Unlock()
}

// emit رویدادها را برای fanout آماده می‌کند (hot path، غیرمسدودکننده).
// در سرریز، صعود کنترلشده per-audience طبق EXE-2026-1006-17 بند ۳.۱.
func (r *Room) emit(events []types.Event) {
	if len(events) == 0 || r.cfg.Transport == nil {
		return
	}
	_, rev := r.GetState()
	streamID := string(r.cfg.ID)
	gameID := string(r.cfg.GameID)

	pending := make([]pendingEvent, 0, len(events))
	for i := range events {
		seq := r.seq.Add(1)
		pending = append(pending, pendingEvent{
			seq:      seq,
			event:    events[i],
			streamID: streamID,
			gameID:   gameID,
			rev:      rev,
		})
	}

	select {
	case r.fanoutCh <- pending:
		return
	default:
	}

	// بافر پر است. صعود کنترلشده.
	r.handleFanoutOverflow(pending)
}

// handleFanoutOverflow سیاست سرریز را per-audience اعمال می‌کند.
//
//   - Internal / ToSpectators → drop + شمارش
//   - ToPlayer → بستن کانال بازیکن با CodeSlowClient
//   - ToRoom → سیگنال backpressure
func (r *Room) handleFanoutOverflow(pending []pendingEvent) {
	for i := range pending {
		p := &pending[i]
		switch p.event.Audience {
		case types.Internal, types.ToSpectators:
			r.dropped.Add(1)
		case types.ToPlayer:
			r.slowClientClosed.Add(1)
			if p.event.Target != "" {
				r.cfg.Transport.CloseSlowClient(
					p.event.Target,
					wire.CodeSlowClient,
					"fanout overflow",
				)
			}
		case types.ToRoom:
			r.backpressure.Add(1)
			r.cfg.Transport.SignalBackpressure(
				wire.CodeBackpressure,
				"fanout overflow",
			)
		}
	}
}

// deliver کار I/O را انجام می‌دهد. فقط روی goroutine fanout.
func (r *Room) deliver(pending []pendingEvent) {
	for i := range pending {
		p := &pending[i]
		ev := wire.NewEvent(p.streamID, p.seq, p.event)
		frame := wire.OutboundMsg{
			Type:   wire.TypeEvent,
			RoomID: p.streamID,
			GameID: p.gameID,
			Rev:    p.rev,
			Event:  &ev,
		}.Bytes()

		switch p.event.Audience {
		case types.Internal:
			continue
		case types.ToRoom, types.ToSpectators:
			r.cfg.Transport.Broadcast(frame)
		case types.ToPlayer:
			if p.event.Target != "" {
				r.cfg.Transport.SendTo(p.event.Target, frame)
			}
		}
	}
}

func (r *Room) reply(env Envelope, err error) {
	if env.ReplyTo == nil {
		return
	}
	select {
	case env.ReplyTo <- err:
	default:
	}
}
