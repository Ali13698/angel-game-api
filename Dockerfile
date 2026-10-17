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

RUN go build -o myapp .

ENTRYPOINT ["/app/myapp"]
