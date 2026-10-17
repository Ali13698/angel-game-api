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

# کامپایل TypeScript به JavaScript با خود Go (بدون نیاز به ابزار خارجی)
RUN go run ./cmd/build-ts

# بیلد نهایی اپلیکیشن
RUN go build -o myapp .

ENTRYPOINT ["/app/myapp"]
