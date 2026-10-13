package room

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
)

// room_contract_test.go — GRK-INT-23 · نسخهٔ نهاییِ قطعی.
// صفر ارجاع به TickInterval · صفر انتظار برای تیکر · seed قطعیِ فاز.
// بیلد روی room.go دست‌نخورده. تستِ مثبتِ اصلی: Do و Close هم‌زمان، هر دو کران‌دار.

const (
	contractDoBudget   = 300 * time.Millisecond
	contractHarnessCap = 3 * time.Second
)

var contractEvents = []types.Event{{Audience: types.ToRoom}}

func contractAllowedErr(err error) bool {
	switch {
	case err == nil,
		errors.Is(err, ErrBusy),
		errors.Is(err, ErrClosed),
		errors.Is(err, ErrWrongPhase),
		errors.Is(err, ErrPriorityFull),
		errors.Is(err, context.Canceled),
		errors.Is(err, context.DeadlineExceeded):
		return true
	}
	return false
}

func contractNoop() Request {
	return Request{From: "alice", Cmd: types.Command{Type: "noop"}}
}

func contractRoom(t *testing.T, inbox int) *Room {
	t.Helper()
	return newTestRoom(t, inbox, nullTransport{}, contractEvents)
}

func seedActive(r *Room) {
	r.advancePhase(PhaseCountdown, time.Now())
	r.advancePhase(PhaseActive, time.Now())
}

func TestRoom_Contract_DoConcurrentClose_BoundedReturn(t *testing.T) {
	const workers = 50

	r := contractRoom(t, 1)
	seedActive(r)

	roomCtx, cancelRoom := context.WithCancel(context.Background())
	defer cancelRoom()
	go r.Run(roomCtx)

	errs := make(chan error, workers)
	for i := 0; i < workers; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), contractDoBudget)
			defer cancel()
			errs <- r.Do(ctx, contractNoop())
		}()
	}

	closeCtx, cancelClose := context.WithTimeout(context.Background(), contractDoBudget)
	defer cancelClose()
	closed := make(chan error, 1)
	go func() { closed <- r.Close(closeCtx) }()

	limit := time.NewTimer(contractHarnessCap)
	defer limit.Stop()

	for i := 0; i < workers; i++ {
		select {
		case err := <-errs:
			if !contractAllowedErr(err) {
				t.Fatalf("Do[%d] = %v؛ خارج از مجموعهٔ مجازِ قراردادی", i, err)
			}
		case <-limit.C:
			t.Fatalf("Do در مهلتِ آشکارساز برنگشت؛ %d فراخوانیِ بی‌پاسخ", workers-i)
		}
	}

	select {
	case err := <-closed:
		if !contractAllowedErr(err) {
			t.Fatalf("Close = %v؛ خارج از مجموعهٔ مجازِ قراردادی", err)
		}
	case <-limit.C:
		t.Fatal("Close در مهلتِ آشکارساز برنگشت")
	}
}

func TestRoom_Contract_HealthyActor_DoSucceeds_CloseIsClean(t *testing.T) {
	r := contractRoom(t, 8)
	seedActive(r)

	roomCtx, cancelRoom := context.WithCancel(context.Background())
	defer cancelRoom()
	go r.Run(roomCtx)

	doCtx, cancelDo := context.WithTimeout(context.Background(), contractDoBudget)
	defer cancelDo()
	if err := r.Do(doCtx, contractNoop()); err != nil {
		t.Fatalf("Do در فازِ ACTIVE = %v؛ انتظار nil", err)
	}

	t.Run("twin_close_stays_bounded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), contractDoBudget)
		defer cancel()
		if err := r.Close(ctx); err != nil && !errors.Is(err, ErrClosed) {
			t.Fatalf("Close دوم = %v", err)
		}
	})

	closeCtx, cancelClose := context.WithTimeout(context.Background(), contractHarnessCap)
	defer cancelClose()
	if err := r.Close(closeCtx); err != nil && !errors.Is(err, ErrClosed) {
		t.Fatalf("Close = %v", err)
	}
}
