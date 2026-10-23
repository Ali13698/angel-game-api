# مرحله ۱: ساخت (Builder)
FROM golang:1.27 AS builder
ARG PORT=8080
ARG BRIDGE_URL="ws://10.101.66.192/bridge"
ARG KAFKA_ADDR=kafka-kub-service:9092
ARG KAFKA_USER=admin

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mo[d] go.su[m] ./
RUN go mod download

COPY . .

# 🛠️ کامپایل بازی‌ها (تولید فایل‌های game.wasm)
RUN go run ./cmd/build-wasm

# بیلد سرور Go
RUN go build -o myapp .

# مرحله ۲: محیط نهایی (Runtime)
FROM debian:bookworm-slim
WORKDIR /app

# انتقال فایل اجرایی و پوشه بازی‌ها از مرحله قبل به محیط نهایی
COPY --from=builder /app/myapp .
COPY --from=builder /app/games ./games

# اطمینان از وجود گواهی‌های SSL
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

EXPOSE 8080
ENTRYPOINT ["/app/myapp"]
