#!/bin/bash
# Praxis Verification Script
# Single-source verification procedure for v0.2.x.
# All docs reference this script; no doc lists individual commands.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=== go test ./... -count=1 ==="
cd "$PROJECT_DIR"
go test ./... -count=1

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
echo "=== tools/replay/replay.sh ==="
bash "$PROJECT_DIR/tools/replay/replay.sh"

echo ""
echo "=== ALL CHECKS PASS ==="
