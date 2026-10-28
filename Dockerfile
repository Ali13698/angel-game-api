# --- Stage 1: Build ---
FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /angel-server ./cmd/server

# --- Stage 2: Final (Production) ---
FROM alpine:latest
# نصب ca-certificates برای ارتباط امن و curl برای دیباگ و ارسال لاگ‌ها
RUN apk --no-cache add ca-certificates curl

WORKDIR /app
COPY --from=builder /angel-server .

EXPOSE 8000
CMD ["./angel-server"]
