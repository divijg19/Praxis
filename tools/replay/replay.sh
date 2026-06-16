#!/bin/bash
# Praxis Replay Verification
# Builds binary and runs Neovim replay against all 51 challenges.
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Building praxis binary..."
cd "$PROJECT_DIR"
go build -o /tmp/praxis ./cmd/praxis/

echo "Running Neovim replay..."
nvim --headless -l "$SCRIPT_DIR/replay.lua"
