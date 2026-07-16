.PHONY: run build test verify certify

run:
	go run ./cmd/praxis

build:
	go build -o praxis ./cmd/praxis

test:
	go test ./...

verify certify:
	bash tools/verify.sh
