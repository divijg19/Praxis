# Testing

## Test Suite Overview

Praxis has **87 tests** across 4 packages:

| Package | Tests | What they verify |
|---|---|---|
| `internal/content` | 33 | Content invariants, layout, stability, contracts, curriculum integrity, taxonomy, describe canonical correctness |
| `internal/stats` | 31 | Stats persistence, attempt/completion tracking, best-value, mastery tiers, distribution, guidance, confidence |
| `internal/validator` | 5 | Validator registry (`cursor`, `buffer`, `composite`), UTF-8 normalization |
| `cmd/praxis` | 18 | CLI subprocess behavior, output format contracts, describe/catalog/help commands, attempt/record commands, confidence levels |

## Running Tests

```bash
# All tests
go test ./... -count=1

# Specific package
go test ./internal/content/... -v

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
| `TestAllChallengesHaveLayer` | Missing Layer field |
| `TestLayerValidValues` | Invalid Layer value (typo, wrong case, stray whitespace) |
| `TestValidatorCoverage` | Unknown validator type |
| `TestNoValidatorDrift` | Registered validator unused by any challenge |
| `TestResultMatchesVerify` | Buffer/composite without Result / cursor with Result |
| `TestResultShapeMatchesVerify` | Cursor challenge with Result / buffer/composite with Target |
| `TestCursorChallengesHaveTargets` | Cursor challenge without Target |
| `TestBufferChallengesHaveNoTargets` | Buffer-like challenge with incorrect Target |
| `TestNoEmptyContent` | Completely empty Content |
| `TestInstructionLinePresent` | Missing first-content-line instruction |
| `TestBufferChallengeLayout` | Buffer-like challenge with <3 lines or missing blank |
| `TestCursorChallengeLayout` | Cursor challenge with <1 line or empty instruction |
| `TestContentResultLineCountReasonable` | Buffer-like challenge with wildly mismatched Content/Result line count |
| `TestCompositeHasEvaluation` | Composite challenge missing Evaluation or invalid MaxMoves |

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
| `TestCurriculumCoverage` | Challenge missing Concept, Context, or Stage |
| `TestConceptContextPairsUnique` | Two challenges with identical (Concept, Context) |
| `TestProgressionCoverage` | Progression stage with zero challenges |
| `TestStageIntroductionOrder` | Stages introduced out of pedagogical order |

## Describe Tests (`internal/content/describe_test.go`)

| Test | What it catches |
|---|---|
| `TestDescriptionForCompleteness` | DescriptionFor returns complete and correct data for all 56 challenges: every field matches source (Challenge + Metadata + Evaluation) |
| `TestDescriptionForUnknown` | Unknown ID returns zero-value Description and false |

## Validator Tests (`internal/validator/validator_test.go`, `utf8_test.go`)

| Test | What it catches |
|---|---|
| `TestExistsCursor` | "cursor" validator unregistered |
| `TestExistsBuffer` | "buffer" validator unregistered |
| `TestExistsComposite` | "composite" validator unregistered |
| `TestExistsUnknown` | Non-existent validator falsely registered |
| `TestUTF8CursorNormalization` | byte_to_char regression with multi-byte content |

## CLI Tests (`cmd/praxis/main_test.go`)

CLI tests build and run the praxis binary as a subprocess, verifying output and exit codes.

| Test | What it catches |
|---|---|
| `TestCatalogOutputStable` | Catalog output names or order drift |
| `TestDescribeCommand` | Describe JSON for a buffer challenge (no Evaluation) |
| `TestDescribeComposite` | Describe JSON for a composite challenge (with Evaluation) |
| `TestDescribeUnknown` | Non-existent describe ID exits 1 |
| `TestHelpCommand` | Help output includes all public commands, excludes removed commands |
| `TestBarePraxis` | Bare invocation shows next challenge + `praxis help` hint |
| `TestRecordStats` | CLI stat recording and best-value persistence |
| `TestStatsCommand` | Per-challenge stats output format |
| `TestStatsSummary` | Summary stats output format |
| `TestStatsUnknownChallenge` | Unknown ID exits 1 for stats command |
| `TestAttemptCommand` | Attempt tracking via CLI |
| `TestAttemptWithRecord` | Full attempt+record workflow and confidence rendering |
| `TestAttemptUnknown` | Non-existent ID exits 1 for attempt |
| `TestStatsCommandConfidenceLevels` | Confidence rendering across all 4 states |

## Adding a New Test

1. Add the test function to the appropriate file
2. If adding a challenge-level invariant, add it to `content_test.go`. If adding a curriculum-integrity test, add it to `integrity_test.go`. If adding a taxonomy test, add it to `taxonomy_test.go`.
3. If adding a CLI-level test, add it to `main_test.go`
4. If adding a new validator, add its test to `validator_test.go`. The `utf8_test.go` file validates byte-to-character normalization for the cursor validator — update it if adding multi-byte content
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
4. Reports per-layer summary and total: `ALL 56/56 REPLAY TESTS PASS`

### Interpreting Results

Each challenge reports `PASS <id>` or `FAIL <id>`.

- A **cursor challenge** passes when navigating to the target character completes successfully
- A **buffer challenge** passes when setting the buffer content to the result matches exactly
- The **UTF-8 challenge** passes when byte-to-character normalization correctly identifies the star at bytecol=9, charcol=6
- A **composite challenge** passes when setting the buffer content matches and moves ≤ MaxMoves

The summary shows per-layer breakout and final count. All 56 must pass.

To add a new challenge: add its ID to `all_ids` in `tools/replay/replay.lua`, run the replay to verify, then commit the updated file.

### Limitations

- Buffer challenges are verified by directly setting the buffer to the result, not by simulating keystrokes
- This means replay tests validate that content+result pairs are consistent, not that the teaching technique works
- Full technique validation requires manual testing within Neovim
