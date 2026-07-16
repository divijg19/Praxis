#!/bin/bash
# Praxis Certification Script
# Single-source release certification for v0.3.x.
# One command provides a complete release verdict: build, lint, format,
# test, and end-to-end (replay + journey) verification.
# All docs reference this script; no doc lists individual commands.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_DIR"

echo "=== go build ./... ==="
go build ./...

echo ""
echo "=== go vet ./... ==="
go vet ./...

echo ""
echo "=== gofmt -l . ==="
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    echo "FAIL: unformatted files:"
    echo "$UNFORMATTED"
    exit 1
fi
echo "PASS: clean"

echo ""
echo "=== golangci-lint run ==="
# Use the module-pinned version via `go run` so it builds with the project's
# Go version (a prebuilt golangci-lint binary can fail on newer Go).
go run github.com/golangci/golangci-lint/cmd/golangci-lint run

echo ""
echo "=== go test ./... -count=1 ==="
go test ./... -count=1

echo ""
echo "=== tools/replay/replay.sh ==="
bash "$PROJECT_DIR/tools/replay/replay.sh"

echo ""
echo "=== tools/journey/journey.sh ==="
bash "$PROJECT_DIR/tools/journey/journey.sh"

echo ""
echo "=== PRAXIS CERTIFIED: ALL CHECKS PASS ==="
