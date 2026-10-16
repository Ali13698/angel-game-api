# Runflare-Environment-Variable-Start
ARG PORT=8080
# Runflare-Environment-Variable-End

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# نصب esbuild (برای build TypeScript)
RUN apk add --no-cache esbuild

COPY . .

# build TypeScript → JS → embed
RUN mkdir -p internal/web/static && \
    cp web/profile/index.html internal/web/static/ && \
    cp web/profile/style.css internal/web/static/ && \
    esbuild web/profile/app.ts --bundle --minify --target=es2020 \
      --outfile=internal/web/static/app.js

RUN go build -o myapp .

ENTRYPOINT ["/app/myapp"]
