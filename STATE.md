# 📌 Angel — وضعیت پروژه
تاریخ: ۱۴۰۵/۰۷/۲۴ (2026-10-16)

## 🟢 سرویس‌های فعال

| سرویس | آدرس | وضعیت |
|---|---|---|
| game-api (Go 1.27) | game.nomi98.com | ✅ |
| angel-bot (Node 20) | nomi98.com | ✅ |
| PostgreSQL 18.1 | angel-db-brp-service:5432 | ✅ |
| Redis | internal | ✅ |
| MinIO | files.nomi98.com | ✅ |
| Bridge (Go↔Node) | internal WS | ✅ |

## 🗄️ PostgreSQL

- DB: `angel-dbhsl_db`
- ۱۷ جدول (auto-migrate via internal/db/migrations/)
- users: ۳۶ ستون، شامل ۷ فیلد پروفایل

## 🎯 معماری قفل‌شده

- **Go**: سرور + موتور بازی
- **TypeScript**: UI مرورگر (پروفایل، فروشگاه، SDK)
- **Unity**: بازی‌های 2D/3D (بعداً)
- **Node (angel-bot)**: ربات‌ها + API موقت (مهاجرت تدریجی به Go)

## 🚧 کارهای مانده

1. API پروفایل در Go
2. صفحه TypeScript پروفایل
3. اتصال دکمه‌های تغییر ظاهر
4. کلاینت Ludo در Unity
5. بازی دوم

## ⚠️ نکات

- `PG_HOST_OVERRIDE` توی env `game-api` → باید پاک شه
- DNS پایدار: `angel-db-brp-service`
- D1 دیگه استفاده نمی‌شه
