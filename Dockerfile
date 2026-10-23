FROM golang:1.27 AS builder
ARG PORT=8080
ARG BRIDGE_URL="ws://10.101.66.192/bridge"
ARG KAFKA_ADDR=kafka-kub-service:9092
ARG KAFKA_USER=admin

# 🛠️ این متغیر کش داکر را می‌شکند تا از این به بعد مجبور شود از صفر بیلد کند
ARG CACHEBUSTER=2026-10-23-v2

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

WORKDIR /app
COPY go.mo[d] go.su[m] ./
RUN go mod download
COPY . .

# ۱. کامپایل بازی‌ها به WASM
RUN go run ./cmd/build-wasm
# 🛠️ این خط لاگ را چاپ می‌کند تا در لاگ‌های Runflare ببینیم فایل‌ها ساخته شده‌اند یا نه
RUN ls -la games/ludo/ && ls -la games/test/

# ۲. کامپایل فرانت‌اند TypeScript
RUN go run ./cmd/build-ts

# ۳. بیلد نهایی سرور Go (با فایل‌های WASM جاسازی شده)
RUN go build -o myapp .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/myapp .
COPY --from=builder /app/games ./games
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
EXPOSE 8080
ENTRYPOINT ["/app/myapp"]
