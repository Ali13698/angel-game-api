FROM golang:1.27 AS builder

WORKDIR /app

# این خط کش رانفلر را می‌شکند تا مجبور شود فایل جدید را بیلد کند
RUN echo "Busting cache at $(date)"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/myapp .
COPY --from=builder /app/internal/web/static ./internal/web/static

EXPOSE 8080
ENTRYPOINT ["/app/myapp"]
