# مرحله ۱: بیلد کردن برنامه (این مرحله سنگین است اما فقط وقتی go.mod تغییر کند کش می‌شود)
FROM golang:1.27-alpine AS builder
WORKDIR /app

# کپی کردن فایل‌های ماژول برای استفاده از کش
COPY go.mod go.sum ./
RUN go mod download

# کپی کردن بقیه کدها
COPY . .

# بیلد کردن به صورت استاتیک (بدون نیاز به کتابخانه‌های سیستمی)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o myapp .

# مرحله ۲: ایمیج نهایی و فوق‌العاده سبک (فقط برای اجرا)
FROM alpine:latest
WORKDIR /app

# کپی کردن برنامه کامپایل شده از مرحله قبل
COPY --from=builder /app/myapp .

# کپی کردن فایل‌های استاتیک (HTML/CSS/JS)
COPY --from=builder /app/internal/web/static ./internal/web/static

EXPOSE 8080
ENTRYPOINT ["/app/myapp"]
