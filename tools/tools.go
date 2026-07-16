//go:build tools
// +build tools

// Package tools pins build/verification tooling as module dependencies so they
// build with the project's Go version via `go run`. See tools/verify.sh.
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
