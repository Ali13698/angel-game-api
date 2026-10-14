FROM node:20-alpine AS web
WORKDIR /web
COPY games/ludo-ts/package.json ./
RUN npm install --no-audit --no-fund
COPY games/ludo-ts/ ./
RUN npm run build

FROM golang:1.23 AS build
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web /web/dist ./games/ludo-ts/dist
RUN go build -o myapp .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /app/myapp ./myapp
COPY --from=build /app/games/ludo-ts/dist ./games/ludo-ts/dist
COPY --from=build /app/web ./web
ENV PORT=3000
EXPOSE 3000
CMD ["/app/myapp"]
