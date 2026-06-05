.PHONY: run build test

run:
	go run ./cmd/praxis

build:
	go build -o praxis ./cmd/praxis

test:
	go test ./...
