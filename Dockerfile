FROM golang:1.27 AS builder
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

WORKDIR /app
COPY go.mo[d] go.su[m] ./
RUN go mod download
COPY . .

# ساخت فرانت‌اند
RUN go run ./cmd/build-ts

# بیلد سرور Go (فایل‌های WASM از قبل در گیت هستند)
RUN go build -o myapp .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/myapp .
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
EXPOSE 8080
ENTRYPOINT ["/app/myapp"]
