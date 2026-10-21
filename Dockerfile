FROM golang:1.27
# Runflare-Environment-Variable-Start
ARG PORT=8080
ARG BRIDGE_URL="ws://10.101.66.192/bridge"
# Runflare-Environment-Variable-End

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mo[d] go.su[m] ./
RUN go mod download

COPY . .

# کامپایل کردن فرانت اند (لابی و پروفایل) قبل از بیلد نهایی
RUN go run ./cmd/build-ts

# بیلد نهایی
RUN go build -o myapp .

ENTRYPOINT ["/app/myapp"]
