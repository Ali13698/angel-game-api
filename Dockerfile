FROM golang:1.27
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

# 🛠️ تضمین ساخت فایل‌های WASM در زمان بیلد داکر
RUN go run ./cmd/build-wasm

RUN go build -o myapp .

ENTRYPOINT ["/app/myapp"]
