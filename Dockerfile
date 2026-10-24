FROM golang:1.27
ARG PORT=8080
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o myapp .
ENTRYPOINT ["/app/myapp"]
