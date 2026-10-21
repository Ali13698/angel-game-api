FROM golang:1.27
ARG PORT=8080
ARG BRIDGE_URL="ws://10.101.66.192/bridge"

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mo[d] go.su[m] ./
RUN go mod download

COPY . .

# ساخت فایل های استاتیک (فرانت اند) قبل از بیلد نهایی
RUN go run ./cmd/build-ts

# بیلد نهایی
RUN go build -o myapp .

ENTRYPOINT ["/app/myapp"]
