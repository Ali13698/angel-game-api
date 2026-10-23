FROM golang:1.27 AS builder
ARG PORT=8080
ARG BRIDGE_URL="ws://10.101.66.192/bridge"
ARG KAFKA_ADDR=kafka-kub-service:9092
ARG KAFKA_USER=admin

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

WORKDIR /app
COPY go.mo[d] go.su[m] ./
RUN go mod download
COPY . .

# 🛠️ اجرای همزمان اسکریپت و بیلد در یک لایه (کش داکر را کاملاً بی‌اثر می‌کند)
RUN go run ./cmd/build-wasm && \
    go build -a -o myapp .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/myapp .

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
EXPOSE 8080
ENTRYPOINT ["/app/myapp"]
