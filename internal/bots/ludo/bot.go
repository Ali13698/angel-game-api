package ludo

import (
	"context"
	"math/rand"
	"time"

	"github.com/angel/angel/internal/core/types"
	ludogame "github.com/angel/angel/internal/games/ludo"
	"github.com/angel/angel/internal/room"
)

type Bot struct {
	UserID types.UserID
	Room   *room.Room
}

func NewBot(userID types.UserID, rm *room.Room) *Bot {
	return &Bot{UserID: userID, Room: rm}
}

func (b *Bot) Run(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.step(ctx)
		}
	}
}

func (b *Bot) step(ctx context.Context) {
	st, _ := b.Room.GetState()
	ls, ok := st.(*ludogame.State)
	if !ok || ls.Done {
		return
	}
	botSeat := -1
	for i, u := range ls.Seats {
		if u == string(b.UserID) {
			botSeat = i
			break
		}
	}
	if botSeat < 0 || ls.Turn != botSeat {
		return
	}
	if ls.Die == 0 {
		_ = b.Room.Submit(room.Request{
			From: b.UserID,
			Cmd:  types.Command{Type: "roll"},
		})
		return
	}
	tok := chooseToken(ls, botSeat)
	if tok < 0 {
		return
	}
	_ = b.Room.Submit(room.Request{
		From: b.UserID,
		Cmd:  types.Command{Type: "move", Payload: map[string]any{"token": tok}},
	})
}

func chooseToken(s *ludogame.State, seat int) int {
	base := seat * ludogame.TokensPerSeat
	die := s.Die
	var cand []int
	for i := 0; i < ludogame.TokensPerSeat; i++ {
		pos := s.Tokens[base+i]
		if pos == ludogame.PosBase {
			if die == 6 {
				cand = append(cand, i)
			}
		} else if pos < ludogame.Finish && pos+die <= ludogame.Finish {
			cand = append(cand, i)
		}
	}
	if len(cand) == 0 {
		return -1
	}
	return cand[rand.Intn(len(cand))]
}
