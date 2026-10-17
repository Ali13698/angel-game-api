# --- Stage 1: Build Frontend (TypeScript -> JavaScript) ---
FROM node:20-alpine AS frontend-builder
WORKDIR /build
COPY web/profile/ ./
RUN npm install -g esbuild && \
    esbuild app.ts --bundle --minify --target=es2020 --outfile=app.js

# --- Stage 2: Build Backend (Go) ---
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# کپی کردن فایل‌های استاتیک و فایل JS کامپایل شده از مرحله قبل
RUN mkdir -p internal/web/static && \
    cp web/profile/index.html internal/web/static/ && \
    cp web/profile/style.css internal/web/static/ && \
    cp /build/app.js internal/web/static/app.js

RUN go build -o myapp .

ENTRYPOINT ["/app/myapp"]
