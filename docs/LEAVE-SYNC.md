# Leave — تصریح همگام‌سازی سه‌مرحله‌ای

پاسخ به بند ۵.۱ نامه EXE-2026-1006-16.

## مراحل Leave

Leave یک بازیکن در سه مرحلهٔ پیوسته انجام می‌شود:

### مرحلهٔ ۱ — Enqueue اتمیک

    select {
    case r.priority <- env:
        // ارسال فوری
    default:
        // خط اولویت پر است: انتظار محدود (≤ PriorityLaneWaitTimeout)
    }

- کانال `r.priority` با ظرفیت ثابت `PriorityLaneSize = 64`.
- ارسال روی channel در Go اتمیک است (تضمین زبان).
- هیچ mutex اضافه‌ای لازم نیست — channel همگام‌سازی داخلی دارد.

### مرحلهٔ ۲ — Drain اتمیک (انتظار پاسخ)

    select {
    case err := <-reply:
        return err
    case <-r.quit:
        return ErrClosed
    case <-ctx.Done():
        return ctx.Err()
    }

- `reply` یک channel ظرفیت ۱ است که به Envelope ضمیمه شده.
- Actor پس از اعمال Leave در `step()` با `r.reply(env, nil)` پاسخ می‌فرستد.
- دریافت روی channel اتمیک است.

### مرحلهٔ ۳ — اعمال وضعیت (در Actor)

    next, events := r.cfg.Game.OnPlayerLeft(prev, env.From)
    ...
    if next.Version() != prev.Version() {
        r.commit(next)   // commit با mutex r.mu
    }

- `r.commit()` یک mutex ساده `r.mu` می‌گیرد.
- همین mutex توسط `GetState()` هم گرفته می‌شود.
- هیچ نقطه‌ای بین سه مرحله وجود ندارد که وضعیت نامتناقض دیده شود.

## پنجرهٔ رقابت باقی‌مانده (پذیرفته‌شده)

اگر `ctx.Done()` در حین مرحلهٔ ۲ رخ دهد، caller Leave برمی‌گردد اما
Room ممکن است هنوز در حال پردازش Leave باشد (fire-and-forget).
Room تضمین می‌کند که Leave نهایتاً اعمال می‌شود، ولی caller از نتیجهٔ
نهایی آگاه نمی‌شود. این رفتار در `CONTRACTS.md` مستند شده است.

## نتیجه

سه مرحله به ترتیب زیر همگام‌سازی می‌شوند:

| مرحله | نوع همگام‌سازی | منبع |
|---|---|---|
| ۱. Enqueue | channel send | `r.priority` |
| ۲. Drain | channel recv | `reply` |
| ۳. Commit | mutex | `r.mu` |

هیچ پنجرهٔ رقابتی که باعث ناهمگام‌سازی وضعیت نهایی شود باقی نمی‌ماند.
