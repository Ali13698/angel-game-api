package room

import (
	"context"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/wire"
)

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
			seq: seq, event: events[i], streamID: streamID, gameID: gameID, rev: rev,
		})
	}
	select {
	case r.fanoutCh <- pending:
		return
	default:
		r.handleFanoutOverflow(pending)
	}
}

func (r *Room) handleFanoutOverflow(pending []pendingEvent) {
	for i := range pending {
		p := &pending[i]
		switch p.event.Audience {
		case types.Internal, types.ToSpectators:
			r.dropped.Add(1)
		case types.ToPlayer:
			r.slowClientClosed.Add(1)
			if p.event.Target != "" {
				r.cfg.Transport.CloseSlowClient(p.event.Target, wire.CodeSlowClient, "fanout overflow")
			}
		case types.ToRoom:
			r.backpressure.Add(1)
			r.cfg.Transport.SignalBackpressure(wire.CodeBackpressure, "fanout overflow")
		}
	}
}

func (r *Room) deliver(pending []pendingEvent) {
	for i := range pending {
		p := &pending[i]

		payloadBytes := wire.EncodeValue(p.event.Payload)

		ev := &wire.GameEvent{
			EventId:  wire.EventID(p.streamID, p.seq),
			StreamId: p.streamID,
			Seq:      int64(p.seq),
			Type:     p.event.Type,
			Payload:  payloadBytes,
		}

		frame := wire.OutboundMsg{
			Type:   wire.TypeEvent,
			RoomID: p.streamID,
			GameID: p.gameID,
			Rev:    p.rev,
			Events: []*wire.GameEvent{ev},
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
