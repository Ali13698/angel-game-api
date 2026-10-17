# Runflare-Environment-Variable-Start
ARG PORT=8080
# Runflare-Environment-Variable-End

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ساخت پوشه استاتیک برای جلوگیری از خطای go:embed
RUN mkdir -p internal/web/static && touch internal/web/static/.gitkeep

# تبدیل TS به JS با ابزار Go (تست شده و کار می‌کند)
RUN go run ./cmd/build-ts

# بیلد نهایی
RUN go build -o myapp .

ENTRYPOINT ["/app/myapp"]
