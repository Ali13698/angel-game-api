.PHONY: test race lint build clean

test:
	go test ./...

race:
	go test -race ./...

lint:
	go vet ./...

build:
	go build ./...

clean:
	go clean -testcache
