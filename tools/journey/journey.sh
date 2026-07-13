#!/bin/bash
# Praxis Learner Journey Verification
# Builds binary and runs the Neovim journey harness (validates experience, not content).
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Building praxis binary..."
cd "$PROJECT_DIR"
go build -o /tmp/praxis ./cmd/praxis/

echo "Running learner journey..."
PATH="/tmp:$PATH" nvim --headless -l "$SCRIPT_DIR/journey.lua"
