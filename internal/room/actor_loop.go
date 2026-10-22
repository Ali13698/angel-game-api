package room

import (
	"context"
	"sync"
	"time"
)

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
	
	r.shutdownOnce.Do(func() { close(r.shutdown) })
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
