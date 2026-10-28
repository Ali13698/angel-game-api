FROM golang:1.27-alpine

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY . .

RUN go build -o myapp .

EXPOSE 8080

ENTRYPOINT ["/app/myapp"]
