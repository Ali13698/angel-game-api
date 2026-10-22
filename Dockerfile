FROM golang:1.27

ARG PORT=8080
ARG BRIDGE_URL="ws://10.101.66.192/bridge"

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    PORT=${PORT} \
    BRIDGE_URL=${BRIDGE_URL}

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go run ./cmd/build-ts

RUN go build -o myapp .

ENTRYPOINT ["/app/myapp"]
