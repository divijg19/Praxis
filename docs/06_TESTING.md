# Testing

## Test Suite Overview

Praxis has **27 tests** across 3 packages:

| Package | Tests | What they verify |
|---|---|---|
| `internal/content` | 15 | Content invariants, layout, stability |
| `internal/validator` | 4 | Validator registry, UTF-8 normalization |
| `cmd/praxis` | 8 | CLI subprocess behavior |

## Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/content/... -v

# Force non-cached run
go clean -testcache && go test ./...

# CLI tests (requires compilation)
go test ./cmd/praxis/... -v
```

## Content Tests (`internal/content/content_test.go`)

These are the most important tests. They protect the challenge registry from accidental corruption.

### Identity protection

| Test | What it catches |
|---|---|
| `TestUniqueChallengeIDs` | Duplicate challenge IDs (collision on insert) |
| `TestUniqueChallengeNames` | Duplicate display names |
| `TestChallengeIDsStable` | Accidental ID rename or reorder |
| `TestChallengeNamesStable` | Accidental name drift |
| `TestChallengeCount` | Challenge count mismatch (from stability list) |

### Field validation

| Test | What it catches |
|---|---|
| `TestAllChallengesHaveVerify` | Missing Verify field |
| `TestValidatorCoverage` | Unknown validator type |
| `TestResultMatchesVerify` | Buffer without Result / cursor with Result |
| `TestCursorChallengesHaveTargets` | Cursor challenge without Target |
| `TestBufferChallengesHaveNoTargets` | Buffer challenge with incorrect Target |
| `TestNoEmptyContent` | Completely empty Content |
| `TestInstructionLinePresent` | Missing first-content-line instruction |
| `TestBufferChallengeLayout` | Buffer challenge with <3 lines or missing blank |
| `TestCursorChallengeLayout` | Cursor challenge with <1 line or empty instruction |

## Validator Tests (`internal/validator/validator_test.go`, `utf8_test.go`)

| Test | What it catches |
|---|---|
| `TestExistsCursor` | "cursor" validator unregistered |
| `TestExistsBuffer` | "buffer" validator unregistered |
| `TestExistsUnknown` | Non-existent validator falsely registered |
| `TestUTF8CursorNormalization` | byte_to_char regression with multi-byte content |

## CLI Tests (`cmd/praxis/main_test.go`)

CLI tests build and run the praxis binary as a subprocess, verifying output and exit codes.

| Test | What it catches |
|---|---|
| `TestListCount` | List output not matching challenge count |
| `TestChallengeLookup` | Challenge content mismatch |
| `TestTargetLookup` | Target output mismatch |
| `TestVerifyLookup` | Verify output mismatch |
| `TestResultLookup` | Result output mismatch |
| `TestUnknownChallengeFails` | Non-existent ID exits 1 |
| `TestUnknownTargetFails` | Non-existent target exits 1 |
| `TestUnknownVerifyFails` | Non-existent verify exits 1 |

## Adding a New Test

1. Add the test function to the appropriate file
2. If adding a challenge-level invariant, add it to `content_test.go`
3. If adding a CLI-level test, add it to `main_test.go`
4. If adding a new validator, add its test to `validator_test.go` and update `utf8_test.go` if relevant
5. Run `go test ./...` to verify

## Replay Verification

See `docs/REPLAY.md` and `tools/replay/replay.sh` for the end-to-end Neovim replay harness that validates all 41 challenges in a real Neovim session.
