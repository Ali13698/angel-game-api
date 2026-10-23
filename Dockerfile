FROM golang:1.27
# Runflare-Environment-Variable-Start
ARG PORT=8080
ARG BRIDGE_URL="ws://10.101.66.192/bridge"
ARG KAFKA_ADDR=kafka-kub-service:9092
ARG KAFKA_USER=admin
# Runflare-Environment-Variable-End

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

WORKDIR /app
COPY go.mo[d] go.su[m] ./
RUN go mod download
COPY . .

# 🛠️ ساخت فایل‌های WASM قبل از بیلد نهایی
RUN go run ./cmd/build-wasm
# 🛠️ ساخت فرانت‌اند
RUN go run ./cmd/build-ts

# 🛠️ بیلد سرور Go با اجبار به خواندن فایل‌های جدید
RUN go build -a -o myapp .

ENTRYPOINT ["/app/myapp"]
