package room

import (
	"context"
	"fmt"
	"time"

	"github.com/angel/angel/internal/cache"
)

const (
	LockTTL = 15 * time.Minute // حداکثر زمان یک بازی فوتبال
)

// AcquireLock: تلاش برای گرفتن قفل برای یک دستگاه
func AcquireLock(ctx context.Context, fingerprint string, userID string) error {
	rdb := cache.Client()
	if rdb == nil {
		return nil // اگر ردیس نبود، اجازه بده (Fail-safe)
	}

	key := fmt.Sprintf("active_device:%s", fingerprint)
	
	// SetNX: فقط اگر کلید وجود نداشت، مقدار را ست کن (اتمیک)
	success, err := rdb.SetNX(ctx, key, userID, LockTTL).Result()
	if err != nil {
		return err
	}
	
	if !success {
		// کلید وجود دارد. ببینیم برای چه کسی است
		existingUser, _ := rdb.Get(ctx, key).Result()
		if existingUser != userID {
			return fmt.Errorf("این دستگاه در حال حاضر با حساب دیگری در بازی است")
		}
		// اگر برای خودش بود، یعنی اتصال مجدد (Reconnect) است، مشکلی نیست
	}
	
	return nil
}

// ReleaseLock: آزاد کردن قفل بعد از پایان بازی
func ReleaseLock(ctx context.Context, fingerprint string, userID string) {
	rdb := cache.Client()
	if rdb == nil { return }
	
	key := fmt.Sprintf("active_device:%s", fingerprint)
	// فقط اگر صاحب قفل خودش باشد، پاکش کن
	existingUser, _ := rdb.Get(ctx, key).Result()
	if existingUser == userID {
		rdb.Del(ctx, key)
	}
}
