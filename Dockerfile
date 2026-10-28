# مرحله اول: بیلد کردن با استفاده از ایمیج سبک Alpine
FROM golang:1.27-alpine AS builder

# نصب git برای دانلود پکیج‌های Go
RUN apk add --no-cache git

WORKDIR /app

# کپی کردن فایل‌های ماژول (برای استفاده از کش)
COPY go.mod go.sum ./
RUN go mod download

# کپی کردن بقیه کدها
COPY . .

# بیلد کردن به صورت استاتیک (بدون وابستگی به C)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o myapp .

# مرحله دوم: ایمیج نهایی و بسیار سبک برای اجرا
FROM alpine:latest

WORKDIR /app

# کپی کردن فایل اجرایی و فایل‌های استاتیک
COPY --from=builder /app/myapp .
COPY --from=builder /app/internal/web/static ./internal/web/static

EXPOSE 8080
ENTRYPOINT ["/app/myapp"]
