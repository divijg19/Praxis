# Testing

## Test Suite Overview

Praxis has **71 tests** across 4 packages:

| Package | Tests | What they verify |
|---|---|---|
| `internal/content` | 25 | Content invariants, layout, stability, contracts, curriculum integrity, taxonomy |
| `internal/stats` | 25 | Stats persistence, attempt/completion tracking, best-value, mastery tiers, distribution, guidance |
| `internal/validator` | 4 | Validator registry, UTF-8 normalization |
| `cmd/praxis` | 17 | CLI subprocess behavior, output format contracts, attempt command |

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
| `TestNoValidatorDrift` | Registered validator unused by any challenge |
| `TestResultMatchesVerify` | Buffer without Result / cursor with Result |
| `TestResultShapeMatchesVerify` | Cursor challenge with Result / buffer challenge with Target |
| `TestCursorChallengesHaveTargets` | Cursor challenge without Target |
| `TestBufferChallengesHaveNoTargets` | Buffer challenge with incorrect Target |
| `TestNoEmptyContent` | Completely empty Content |
| `TestInstructionLinePresent` | Missing first-content-line instruction |
| `TestBufferChallengeLayout` | Buffer challenge with <3 lines or missing blank |
| `TestCursorChallengeLayout` | Cursor challenge with <1 line or empty instruction |
| `TestContentResultLineCountReasonable` | Buffer challenge with wildly mismatched Content/Result line count |

## Integrity Tests (`internal/content/integrity_test.go`)

These tests audit the curriculum itself rather than structural invariants.

| Test | What it catches |
|---|---|
| `TestCoreConceptCoverage` | Core Vim concept accidentally removed from curriculum |
| `TestNoDuplicateChallengeContent` | Content or buffer result duplicated across challenges |
| `TestCurriculumMapComplete` | Challenge added without curriculum mapping, or orphaned map entry |

## Taxonomy Tests (`internal/content/taxonomy_test.go`)

These tests enforce the Concept–Context–Stage metadata model, which is the single source of truth for curriculum documentation.

| Test | What it catches |
|---|---|
| `TestCurriculumContextsComplete` | Challenge missing Concept, Context, or Stage |
| `TestConceptContextPairsUnique` | Two challenges with identical (Concept, Context) |
| `TestProgressionCoverage` | Progression stage with zero challenges |
| `TestStageIntroductionOrder` | Stages introduced out of pedagogical order |

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
| `TestListOutputStable` | List output names or order drift |
| `TestChallengeLookup` | Challenge content mismatch |
| `TestTargetLookup` | Target output mismatch |
| `TestTargetOutputStable` | Target output format drift |
| `TestVerifyLookup` | Verify output mismatch |
| `TestVerifyOutputStable` | Verify output format drift |
| `TestResultLookup` | Result output mismatch |
| `TestUnknownChallengeFails` | Non-existent ID exits 1 |
| `TestUnknownTargetFails` | Non-existent target exits 1 |
| `TestUnknownVerifyFails` | Non-existent verify exits 1 |
| `TestRecordStats` | CLI stat recording and best-value persistence |
| `TestStatsCommand` | Per-challenge stats output format |
| `TestStatsSummary` | Summary stats output format |
| `TestStatsUnknownChallenge` | Unknown ID exits 1 for stats command |

## Adding a New Test

1. Add the test function to the appropriate file
2. If adding a challenge-level invariant, add it to `content_test.go`. If adding a curriculum-integrity test, add it to `integrity_test.go`. If adding a taxonomy test, add it to `taxonomy_test.go`.
3. If adding a CLI-level test, add it to `main_test.go`
4. If adding a new validator, add its test to `validator_test.go` and update `utf8_test.go` if relevant
5. Run `go test ./...` to verify

## Replay Verification

Replay verification runs every challenge through a real Neovim session, ensuring that challenge content, results, and validator behavior are consistent. This is the canonical integration verification — it exercises the entire pipeline from Go binary to Lua frontend to Neovim buffer state.

### Running

```bash
tools/replay/replay.sh
```

The script:
1. Builds the Go binary to `/tmp/praxis`
2. Runs Neovim headless with `tools/replay/replay.lua`
3. Prints PASS/FAIL for each challenge
4. Reports summary: `ALL 41/41 REPLAY TESTS PASS`

### Interpreting Results

Each challenge reports `PASS <id>` or `FAIL <id>`.

- A **cursor challenge** passes when navigating to the target character completes successfully
- A **buffer challenge** passes when setting the buffer content to the result matches exactly
- The **UTF-8 challenge** passes when byte-to-character normalization correctly identifies the star at bytecol=9, charcol=6

The summary line shows the final count. All 41 must pass.

### Adding a New Challenge

1. Add its ID to the appropriate list in `tools/replay/replay.lua` (cursor challenges in `cursor_ids`, buffer challenges in `buffer_ids`)
2. Run the replay to verify
3. Commit the updated `replay.lua`

### Limitations

- Buffer challenges are verified by directly setting the buffer to the result, not by simulating keystrokes
- This means replay tests validate that content+result pairs are consistent, not that the teaching technique works
- Full technique validation requires manual testing within Neovim
