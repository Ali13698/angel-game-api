# --- Stage 1: Build ---
FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# بیلد از ریشه پروژه (نقطه) تا main.go اصلی و فایل‌های استاتیک لود شوند
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /angel-server .

# --- Stage 2: Final (Production) ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /app
COPY --from=builder /angel-server .
EXPOSE 8080
CMD ["./angel-server"]
