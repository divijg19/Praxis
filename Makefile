.DEFAULT_GOAL := help
GOLANGCI_LINT_VERSION ?= v2.15.0

.PHONY: help run build test race verify certify install-tools clean

help:            ## list available targets
	@awk -F':.*##' '/^[a-z][a-z -]*:.*##/ {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run:             ## run the CLI
	go run ./cmd/praxis

build:           ## build ./praxis
	go build -o praxis ./cmd/praxis

test:            ## unit + integrity tests
	go test ./... -count=1

race:            ## optional race check
	go test -race ./...

verify certify:  ## the full gate (authoritative)
	go run ./tools/verify

install-tools:   ## install the linter at the version CI uses
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

clean:           ## remove the built binary
	rm -f praxis
