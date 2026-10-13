# مستند قراردادها — v1.0.0+pentad

این سند قراردادهای رسمی فاز ۰ را تعریف می‌کند. هر تغییری در این
قراردادها نیازمند تأیید کتبی مشاور ارشد است.

---

## ۱. Command

    Command {
        ID        string
        RoomID    RoomID
        Type      string
        Payload   any
        Version   uint64
    }

قواعد:
- `ID` شناسهٔ یکتای درخواست است؛ برای تطبیق پاسخ.
- `RoomID` اتاقی که دستور برای آن ارسال شده.
- `Type` نوع دستور (مثلاً `roll`, `move`, `join`, `ping`).
- `Payload` دادهٔ دستور (اختیاری).
- `Version` نسخهٔ وضعیتی که کلاینت دیده است.

اعتبارسنجی نسخه (فاز ۰، قرارداد):
- `Version == 0` یعنی «بدون اطلاعات نسخه» — پذیرفته می‌شود.
- `Version > 0` و `Version < state.Version()` یعنی نسخهٔ کهنه — رد می‌شود.
- استفادهٔ منطقی از `Version` در Game Handler → فاز ۱.

---

## ۲. Event

    Event {
        Type      string
        Payload   any
        Audience  Audience
        Target    UserID
    }

مقادیر مجاز `Audience`:

| مقدار | معنا |
|---|---|
| `ToRoom` | برای همهٔ اعضای اتاق |
| `ToPlayer` | فقط برای بازیکن مشخص‌شده در `Target` |
| `ToSpectators` | برای تماشاگران |
| `Internal` | فقط داخلی؛ به بیرون ارسال نمی‌شود |

قواعد:
- `Target` فقط برای `ToPlayer` معنا دارد.
- `Broadcast()` برای `ToRoom` و `ToSpectators` مقدار `true` برمی‌گرداند.

---

## ۳. رابط GameHandler

    GameID() string
    NewState(players []UserID) State
    ApplyCommand(state State, from UserID, cmd Command) (State, []Event, error)
    OnPlayerLeft(state State, id UserID) (State, []Event)
    Tick(state State, now time.Time) (State, []Event)
    IsTerminal(state State) bool

قواعد:
- هر بازی باید این رابط را کامل پیاده‌سازی کند.
- `Game` نام قدیمی این رابط است و صرفاً برای سازگاری موقت نگه داشته شده؛
  در کد جدید از `GameHandler` استفاده شود.

---

## ۴. رابط GameState

    Version() uint64
    Clone() State

قواعد:
- `Version()` نسخهٔ یکنوای صعودی وضعیت را برمی‌گرداند.
- `Clone()` یک کپی مستقل از وضعیت می‌سازد (deep copy برای فیلدهای تغییرپذیر).

---

## ۵. چرخهٔ حیات اتاق (Room Actor)

    LOBBY → COUNTDOWN → ACTIVE → FINISHED | ABORTED

قواعد:
- گذارها فقط از داخل Actor انجام می‌شوند.
- دستورات بازی فقط در فاز `ACTIVE` پذیرفته می‌شوند؛ در غیر این صورت
  با `ErrWrongPhase` رد می‌شوند.
- `LOBBY → COUNTDOWN` هنگام شروع Actor.
- `COUNTDOWN → ACTIVE` پس از پایان زمان شمارش معکوس.
- `ACTIVE → FINISHED` وقتی `GameHandler.IsTerminal(state)` مقدار `true` بدهد.
- هر فاز → `ABORTED` در صورت بسته شدن یا کنسل شدن.

---

## ۶. صف ورودی Room Actor

- صف اصلی: bounded (پیش‌فرض ۶۴). سرریز → `ErrBusy` (بدون بلاک).
- خط اولویت: bounded (۶۴) برای `PlayerLeft` و `Shutdown`؛ هرگز رد نمی‌شود.
- هر دو صف در حلقهٔ `select` Actor بررسی می‌شوند.

---

## ۷. خطاهای رسمی

| خطا | کاربرد |
|---|---|
| `ErrClosed` | اتاق بسته شده است |
| `ErrBusy` | صف ورودی پر است |
| `ErrWrongPhase` | دستور در فاز جاری پذیرفته نمی‌شود |
| `ErrStaleVersion` | نسخهٔ دستور قدیمی‌تر از وضعیت فعلی است |
