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

# ترکیب دستورات در یک لایه برای جلوگیری از کش شدن در داکر
RUN mkdir -p internal/web/static && touch internal/web/static/.gitkeep && \
    go run ./cmd/build-ts && \
    go build -o myapp .

ENTRYPOINT ["/app/myapp"]
