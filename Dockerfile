FROM golang:1.27-alpine
# Runflare-Environment-Variable-Start
ARG PORT=8080
ARG KAFKA_ADDR=kafka-kub-service:9092
ARG KAFKA_USER=admin
# Runflare-Environment-Variable-End

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 🔴 اصلاح: بیلد از مسیر cmd/server
RUN go build -o angel-server ./cmd/server

EXPOSE 8080
ENTRYPOINT ["/app/angel-server"]
