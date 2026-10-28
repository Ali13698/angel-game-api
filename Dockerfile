FROM golang:1.27-alpine AS builder

# این خط برای تست است تا در لاگ رانفلر ببینیم فایل جدید اجرا شده یا نه
RUN echo "======================================="
RUN echo "=== BUILDING WITH NEW ALPINE DOCKERFILE ==="
RUN echo "======================================="

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o myapp .

FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/myapp .
COPY --from=builder /app/internal/web/static ./internal/web/static

EXPOSE 8080
ENTRYPOINT ["/app/myapp"]
