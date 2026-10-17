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

# ساخت فایل‌های استاتیک با Go و سپس بیلد
RUN mkdir -p internal/web/static && \
    go run ./cmd/build-ts && \
    if [ ! -f "internal/web/static/index.html" ]; then echo "ERROR: index.html not generated!" && exit 1; fi && \
    go build -o myapp .

ENTRYPOINT ["/app/myapp"]
