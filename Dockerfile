FROM registry-de.runflare.com/public/golang:1.23
# Runflare-Environment-Variable-Start
ARG PORT=8080
ARG GHT="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"
# Runflare-Environment-Variable-End

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mo[d] go.su[m] ./
RUN go mod download

COPY . .

RUN go build -o myapp .

ENTRYPOINT ["/app/myapp"]
