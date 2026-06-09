# API Contract

This document defines the public API contracts of Praxis. These contracts are stable: future releases must preserve backward compatibility unless a major version bump documents the break.

## Challenge Contracts

### Challenge.ID

**STABLE** — Must never change after first release.

The ID is the permanent identifier for a challenge. It is used by the CLI, Lua frontend, replay system, and documentation. Renaming an ID silently breaks all consumers.

Enforced by: `TestChallengeIDsStable`

### Challenge.Name

**EVOLVABLE** — May change as curriculum framing shifts.

The display name is a user-facing label. The `yank_line_hunter` → `Unnamed Register Hunter` precedent (v0.0.19) demonstrates that names describe pedagogical intent, not identity.

Enforced by: `TestChallengeNamesStable`

### Challenge.Verify

**STABLE** — Must be `"cursor"` or `"buffer"`.

The verify field selects which validator evaluates success. Adding a new validator requires registering it in `internal/validator/validator.go` and documenting it in `VALIDATOR_CONTRACT.md`.

Enforced by: `TestValidatorCoverage`, `TestNoValidatorDrift`

### Challenge.Target

**STABLE** — Required for cursor challenges, must be empty for buffer challenges.

The target is the character the user must navigate to (cursor) or the empty string (buffer).

Enforced by: `TestCursorChallengesHaveTargets`, `TestBufferChallengesHaveNoTargets`, `TestResultShapeMatchesVerify`

### Challenge.Result

**STABLE** — Required for buffer challenges, must be nil/empty for cursor challenges.

Result is the exact target buffer state after the user completes an editing challenge. Comparison is line-by-line exact string equality.

Enforced by: `TestBufferChallengesHaveResults`, `TestResultMatchesVerify`, `TestResultShapeMatchesVerify`

### Challenge.Content

**STABLE** — First line is the instruction line. Layout is part of the contract.

Content is the initial buffer state shown to the user. For buffer challenges: instruction, blank separator, then play area (3+ lines). For cursor challenges: instruction, optional blank separator, play area (1+ lines).

Enforced by: `TestInstructionLinePresent`, `TestNoEmptyContent`, `TestBufferChallengeLayout`, `TestCursorChallengeLayout`

## CLI Contracts

### `praxis list`

Outputs one challenge name per line, in curriculum order.

```
Motion Rush
Grid Rush
Find Hunter
...
```

- 41 lines
- UTF-8 encoded
- Trailing newline after each name

Enforced by: `TestListCount`, `TestListOutputStable`

### `praxis challenge <id>`

Outputs the challenge content lines.

- First line is the instruction
- Trailing newline after each line
- Exits 0 on success, 1 on unknown ID
- Stderr: `unknown challenge: <id>` on failure

Enforced by: `TestChallengeLookup`, `TestUnknownChallengeFails`

### `praxis target <id>`

Outputs the target character on a single line.

- Single line with trailing newline
- Empty for buffer challenges
- Exits 0 on success, 1 on unknown ID
- Stderr: `unknown challenge: <id>` on failure

Enforced by: `TestTargetLookup`, `TestTargetOutputStable`, `TestUnknownTargetFails`

### `praxis verify <id>`

Outputs the validator name on a single line.

- Single line with trailing newline
- `"cursor"` or `"buffer"`
- Exits 0 on success, 1 on unknown ID
- Stderr: `unknown challenge: <id>` on failure

Enforced by: `TestVerifyLookup`, `TestVerifyOutputStable`, `TestUnknownVerifyFails`

### `praxis result <id>`

Outputs the result lines for buffer challenges.

- Trailing newline after each line
- Empty output for cursor challenges
- Exits 0 on success, 1 on unknown ID
- Stderr: `unknown challenge: <id>` on failure

Enforced by: `TestResultLookup`

### `praxis record <id> <moves> <time_ms>`

Records a challenge completion and updates persistent stats.

- Silent on success (no stdout)
- Exit 0 on success
- Used internally by the Neovim frontend

### `praxis stats [id]`

Outputs challenge statistics.

With an ID, shows per-challenge stats:

```
Attempts: 12
Completions: 10
Best Moves: 2
Best Time: 180ms
```

Without arguments, shows summary:

```
Challenges Completed: 31/41
Total Attempts: 245
```

- Exits 0 on success
- Exits 1 on unknown challenge ID
- Stderr: `unknown challenge: <id>` on failure

## Stats Contract

See `docs/STATS.md` for the full stats system documentation.

Stats are stored in `~/.local/share/praxis/stats.json` (`$XDG_DATA_HOME/praxis/stats.json`).

Each record:

```go
type Stats struct {
    Attempts    int     // total completions
    Completions int     // same as Attempts (every recorded attempt is a success)
    BestMoves   int     // lowest moves across all completions
    BestTimeMs  int     // fastest time across all completions
    LastPlayed  string  // "2006-01-02" format
}
```

## Replay Contract

The replay system is the canonical integration verification for Praxis. See `docs/REPLAY.md` for full documentation.

- 41/41 challenges must pass
- Deterministic for a given binary and content set
- Headless Neovim session (no display required)
- CI-safe with Neovim installed

Enforced by: `tools/replay/replay.sh` + `tools/replay/replay.lua`

## Validator Contracts

See `docs/VALIDATOR_CONTRACT.md` for formal specifications of the cursor and buffer validators, including PASS and FAIL examples.

## Release Process

See `docs/RELEASE_CHECKLIST.md` for the mandatory 10-step release process.
