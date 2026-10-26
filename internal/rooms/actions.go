package rooms

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
)

// ProcessRoll: سرور تاس می‌اندازد. مرورگر فقط درخواست می‌فرستد.
func (r *Room) ProcessRoll(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.CurrentTurn != userID {
		return errors.New("not your turn")
	}

	// تولید تاس تصادفی امن در سطح رمزنگاری
	nBig, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		return errors.New("random generation failed")
	}
	r.LastDice = int(nBig.Int64()) + 1

	r.MustRoll = false 

	// ارسال رویداد تاس به همه
	r.broadcastEvent("ludo.rolled", map[string]interface{}{
		"seat": r.GetPlayerSeat(userID),
		"dice": r.LastDice,
	})
	return nil
}

// ProcessMove: سرور حرکت را تایید می‌کند.
func (r *Room) ProcessMove(ctx context.Context, userID string, tokenIndex int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.CurrentTurn != userID {
		return errors.New("not your turn")
	}
	if r.LastDice == 0 {
		return errors.New("must roll first")
	}

	seat := r.GetPlayerSeat(userID)
	currentPos := r.Tokens[seat][tokenIndex]
	var newPos int
	captured := false

	// قوانین سمت سرور
	if currentPos == -1 { // در خانه
		if r.LastDice != 6 {
			return errors.New("need 6 to exit base")
		}
		newPos = 0
	} else {
		newPos = currentPos + r.LastDice
		if newPos > 38 { // 38 = پایان (Finish)
			return errors.New("overshoot finish")
		}
	}

	// اعمال حرکت
	r.Tokens[seat][tokenIndex] = newPos
	r.LastDice = 0 // ریست تاس

	// ارسال رویداد حرکت
	r.broadcastEvent("ludo.moved", map[string]interface{}{
		"seat":  seat,
		"token": tokenIndex,
		"pos":   newPos,
	})

	// بررسی برنده شدن یا نوبت بعدی
	r.checkWinAndNextTurn(userID, newPos)
	return nil
}

// checkWinAndNextTurn بررسی می‌کند که آیا بازی تمام شده یا نوبت عوض می‌شود
func (r *Room) checkWinAndNextTurn(userID string, newPos int) {
	seat := r.GetPlayerSeat(userID)
	
	// 1. بررسی برنده شدن (رسیدن هر ۴ مهره به ۳۸)
	finishedTokens := 0
	for _, pos := range r.Tokens[seat] {
		if pos == 38 {
			finishedTokens++
		}
	}
	if finishedTokens == 4 {
		r.broadcastEvent("game.over", map[string]interface{}{
			"winner": seat,
		})
		return
	}

	// 2. تعیین نوبت بعدی (اگر ۶ آورد یا مهره زد، دوباره خودش)
	// طبق قرارداد کلاینت، سرور باید همیشه turn_changed را بفرستد
	nextSeat := seat
	if r.LastDice != 6 && !captured { // منطق ساده شده
		nextSeat = (seat + 1) % 2
	}

	r.CurrentTurn = r.Players[nextSeat].ID
	r.broadcastEvent("turn_changed", map[string]interface{}{
		"turn": nextSeat,
	})
}
