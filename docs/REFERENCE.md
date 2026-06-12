# Reference

## Challenge Model

### Struct

Every challenge is defined in `internal/content/content.go` as part of the `All()` function:

```go
challenge.Challenge{
    ID:      "unique_identifier",
    Name:    "Display Name",
    Verify:  "cursor",         // or "buffer"
    Target:  "★",              // cursor challenges only; empty for buffer
    Content: []string{
        "Instruction line",
        "",
        "play area content",
    },
    Result:  []string{         // buffer challenges only; nil for cursor
        "Instruction line",
        "",
        "edited content",
    },
}
```

### Field Rules

| Field | Stability | Rule |
|---|---|---|
| `ID` | STABLE | Globally unique. Never change after first release. `snake_case` with `_hunter` suffix. |
| `Name` | EVOLVABLE | Display name. May change as curriculum framing shifts. |
| `Verify` | STABLE | `"cursor"` or `"buffer"`. Must match a registered validator. |
| `Target` | STABLE | Required for cursor (non-empty). Must be present in buffer. Empty for buffer. |
| `Content` | STABLE | First line is instruction. Buffer: 3+ lines (instruction, blank, play area). Cursor: 1+ lines. |
| `Result` | STABLE | Required for buffer (non-empty slice, exact target state). Nil for cursor. |

### Adding a New Challenge

1. Add the challenge to the end of `All()` in `internal/content/content.go`
2. Update `TestChallengeIDsStable` with the new ID in the correct position
3. Update `TestChallengeNamesStable` with the new Name
4. Add the ID to the appropriate list in `tools/replay/replay.lua` (cursor or buffer)
5. Run `go test ./...` to verify
6. Run `tools/replay/replay.sh` to verify

## CLI Surface

### `praxis list`

Outputs one challenge name per line, in curriculum order. 41 lines, UTF-8 encoded, trailing newline.

Enforced by: `TestListCount`, `TestListOutputStable`

### `praxis challenge <id>`

Outputs the challenge content lines. First line is the instruction. Trailing newline after each line. Exits 0 on success, 1 on unknown ID. Stderr: `unknown challenge: <id>`.

Enforced by: `TestChallengeLookup`, `TestUnknownChallengeFails`

### `praxis target <id>`

Outputs the target character on a single line. Empty for buffer challenges. Exits 0 on success, 1 on unknown ID.

Enforced by: `TestTargetLookup`, `TestTargetOutputStable`, `TestUnknownTargetFails`

### `praxis verify <id>`

Outputs the validator name (`"cursor"` or `"buffer"`) on a single line. Exits 0 on success, 1 on unknown ID.

Enforced by: `TestVerifyLookup`, `TestVerifyOutputStable`, `TestUnknownVerifyFails`

### `praxis result <id>`

Outputs the result lines for buffer challenges. Trailing newline after each line. Empty output for cursor challenges. Exits 0 on success, 1 on unknown ID.

Enforced by: `TestResultLookup`

### `praxis attempt <id>`

Records an attempt on a challenge. Silent on success (no stdout). Exit 0 on success, 1 on unknown ID. Used internally by the Neovim frontend on challenge start and replay.

### `praxis record <id> <moves> <time_ms>`

Records a challenge completion and updates persistent stats. Silent on success (no stdout). Exit 0 on success. Used internally by the Neovim frontend.

### `praxis stats [id]`

With an ID, shows per-challenge stats:

```
Attempts: 10
Completions: 10
Success Rate: 100%
Best Moves: 2
Best Time: 180ms
Mastery: Experienced
```

Without arguments, shows summary:

```
Challenges Completed: 31/41
Total Attempts: 245

Mastery:
  Unseen: 10
  Learning: 5
  Practiced: 8
  Experienced: 18

Highest Tier: Experienced
```

Followed by practice guidance sections:

```
Next Challenge:
  find_hunter

Recommended Review:
  motion_rush
```

Exits 0 on success, 1 on unknown challenge ID.

Enforced by: `TestStatsCommand`, `TestStatsSummary`, `TestRecordStats`

## Validators

Validators determine how a challenge's success condition is evaluated. Praxis has two: **cursor** and **buffer**.

### Validator Registry

Defined in `internal/validator/validator.go`:

```go
var valid = map[string]bool{
    "cursor": true,
    "buffer": true,
}

func Exists(name string) bool {
    return valid[name]
}
```

### Cursor Validator

**Verify value:** `"cursor"`

**Rule:** Navigate the cursor to the target character on screen.

Buffer is set to `Content` with `modifiable=false`. A `CursorMoved` autocommand checks if the character under the cursor matches `Target`. Uses `byte_to_char()` normalization to handle multi-byte content correctly. On match: sets `state.done = true`, echoes success.

**Used by:** 20 challenges (movement, search, navigation, UTF-8 proof).

**Formal specification:**

```
Given:
  buf  = challenge content (array of N lines)
  col0 = cursor byte column (0-indexed)
  row  = cursor row (1-indexed)
  line = buf[row-1]

The cursor position is valid when:
  line is not nil
  AND strcharpart(line, rune_offset(line, col0), 1) == target

Where:
  rune_offset(line, col0) = number of runes in line[0:col0]
  strcharpart(s, start, len) = substring of s at rune start, len runes long
```

**PASS examples:**

| Challenge | Content (play area) | Target | Correct cursor position |
|---|---|---|---|
| motion_rush | `Move your cursor to the star ★` | ★ | On the ★ character |
| find_hunter | `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa★` | ★ | On the ★ (last character) |
| symbol_hunter | `.......................@` | @ | On the @ (last character) |
| utf8_cursor_hunter | `α β γ ★` | ★ | On the ★ (bytecol=9, charcol=6) |

**FAIL examples:**

| Scenario | Why it fails |
|---|---|
| Cursor on wrong character | Character under cursor does not equal target |
| Target not in buffer | No matching character exists |
| Empty target | Target field is `""` (violates contract) |

### Buffer Validator

**Verify value:** `"buffer"`

**Rule:** Edit the buffer content to exactly match `Result`.

Buffer is set to `Content` with `modifiable=true`. A `TextChanged` autocommand fires on Normal mode edits and increments moves count. `TextChangedI` fires on Insert mode keystrokes but does NOT increment moves. `check_buffer()` compares every line of the current buffer against `Result` using byte-exact string equality. On match: sets `state.done = true`, echoes success.

**Used by:** 21 challenges (editing, structural editing, registers).

**Formal specification:**

```
Given:
  current = current buffer lines (array of M strings)
  result  = expected result lines (array of N strings)

The buffer is valid when:
  M == N
  AND current[i] == result[i] for all i in [0, N)

Comparison is byte-exact string equality:
  - No trimming, no case folding, no whitespace normalization
```

**PASS examples:**

| Challenge | Result (play area) | Editing technique |
|---|---|---|
| delete_character_hunter | `hello` | `x` on extra `l` |
| toggle_case_hunter | `hello` | `~` on each letter |
| register_duplicate_hunter | `foo` + `bar` + `foo` | `"ayy` + `"ap` to duplicate |

**FAIL examples:**

| Scenario | Why it fails |
|---|---|
| Buffer has wrong content | At least one line differs from result |
| Line count mismatch | Extra or missing line |
| Trailing whitespace differs | `"keep "` vs `"keep"` — byte-exact comparison |
| Result is empty | Result field is nil or `[]string{}` (violates contract) |

### Autocommand Dispatching

The Lua frontend only attaches `CursorMoved` for cursor challenges. Buffer challenges use `TextChanged` + `TextChangedI` instead. Handled by checking `state.verify` when setting up autocommands.

### Contract Enforcement

| Test | Enforces |
|---|---|
| `TestValidatorCoverage` | Every challenge references a registered validator |
| `TestNoValidatorDrift` | Every registered validator is used by at least one challenge |
| `TestResultShapeMatchesVerify` | Cursor challenges forbid Result; buffer challenges forbid Target |
| `TestExistsCursor` / `TestExistsBuffer` | Validators are present in the registry |

### Future Candidates

Not implemented. Listed for architectural awareness: `selection`, `state`, `register`, `operator`.

## Stats

### Storage

Stats are stored in `~/.local/share/praxis/stats.json` (`$XDG_DATA_HOME/praxis/stats.json`). Simple JSON, human-readable and editable. No database, no migrations. If the file is missing or corrupt, Praxis silently starts fresh.

### Schema

```go
type Stats struct {
    Attempts    int    // total attempts (including replays and abandoned starts)
    Completions int    // total successful completions
    BestMoves   int    // lowest moves across all completions
    BestTimeMs  int    // fastest time across all completions
    LastPlayed  string // "2006-01-02" format
}
```

### Mastery Tiers

Derived from `Completions`:

| Tier | Threshold |
|---|---|
| Unseen | 0 |
| Learning | 1–2 |
| Practiced | 3–7 |
| Experienced | 8+ |

### Design Decisions

- Stats are recorded only on successful completion
- Best values track minima (fewest moves, fastest time)
- First completion always sets the best value
- Subsequent completions only update if better
- No per-attempt history — only the best is preserved
- Attempts and Completions are tracked independently. Attempts increment on challenge start and replay; Completions increment only on success

### Practice Guidance

`NextChallenge()` returns the first curriculum-ordered challenge whose Completions ≤ 2 (Unseen or Learning). This means: **finish what you started** — a partially-practiced Learning challenge will be recommended before a new Unseen challenge.

`RecommendedReview()` returns the oldest Practiced challenge by LastPlayed date, falling back to the oldest Experienced if no Practiced challenges exist. Practiced challenges are preferred because they are more likely to benefit from review than deeply-ingrained Experienced ones.

## Sessions

### Overview

A session represents a contiguous period of deliberate practice within a single Neovim instance. Sessions are ephemeral — they exist only in memory and are discarded when Neovim exits. There is exactly one session per Neovim instance. It begins when the first `:Praxis <id>` is called and accumulates data until Neovim is closed.

### Command

```
:PraxisSession
```

Opens a read-only scratch buffer showing the current session's aggregated metrics.

```
Session

Challenges: 8
Completions: 5

Session Length: 18m42s
Practice Time: 6m12s

Moves: 74

Avg Moves: 14
Avg Time: 74s
```

### Fields

| Field | Meaning |
|---|---|
| **Challenges** | Total `:Praxis <id>` invocations in this session |
| **Completions** | Total successful completions (including replays) |
| **Session Length** | Wall-clock time since first `:Praxis <id>` |
| **Practice Time** | Sum of individual challenge completion times |
| **Moves** | Total moves across all completions |
| **Avg Moves** | Moves / Completions |
| **Avg Time** | Practice Time / Completions |

Completions counts every successful run, including replays via the `r` key. A single challenge may contribute multiple completions within one session, so `Completions` can exceed `Challenges`.

### Design Decisions

- **No persistence.** Sessions are Neovim-scoped ephemeral state.
- **No automatic popup.** `:PraxisSession` is user-invoked only.
- **Lua-only, no Go.** Session lives as a module-level table in `nvim/lua/praxis/init.lua`.
- **Challenges and Completions are separate counters.** The gap shows unfinished attempts.
- **Session Length tracks wall-clock time.** Not idle-adjusted. Honest but approximate.

## Integrity Guarantees

### Principles

1. **One concept per challenge.** Every challenge has a documented primary Concept, Context, and Stage. The curriculum metadata in `internal/content/taxonomy_test.go` is the single source of truth.
2. **IDs are permanent.** Challenge identifiers must never be renamed or removed.
3. **Validators are stable.** Every registered validator must be used by at least one challenge.
4. **Curriculum growth is intentional.** New challenges require a documented primary concept and must not duplicate existing content.
5. **Replay verification is mandatory.** Every challenge must be solvable via the replay harness.
6. **Duplicate challenge content is prohibited.** Challenges may share concepts when they teach the concept in a different context or composition.

### Test Suite

| Test | Enforces |
|---|---|
| `TestCoreConceptCoverage` | Core Vim concepts remain represented |
| `TestNoDuplicateChallengeContent` | No unintended duplicate exercises |
| `TestCurriculumMapComplete` | Every challenge mapped, no orphaned entries |
| `TestUniqueChallengeIDs` | No ID collisions |
| `TestChallengeCount` | No accidental addition/removal |
| `TestNoValidatorDrift` | Validator usage stays current |
| `TestChallengeIDsStable` | IDs never renamed |
| `TestCurriculumContextsComplete` | Every challenge has Concept, Context, Stage |
| `TestConceptContextPairsUnique` | No duplicate (Concept, Context) pairs |
| `TestProgressionCoverage` | All progression stages have challenges |
| `TestStageIntroductionOrder` | Stages introduced in pedagogical order |

### Anti-Goals

- Challenge difficulty ratings
- Mastery scoring
- Skill trees or dependencies
- Progression gates


## Release Procedure

### 10-Step Process

**1. Format**

```bash
gofmt -l .
```

Must produce no output. If it lists files, run `gofmt -w` on them.

**2. Build**

```bash
go build ./...
```

All packages must compile without errors.

**3. Vet**

```bash
go vet ./...
```

All packages must pass static analysis.

**4. Test**

```bash
go test ./...
```

All tests must pass. Final line must report `PASS` or `ok` for every package.

**5. Replay**

```bash
tools/replay/replay.sh
```

Must report `ALL 41/41 REPLAY TESTS PASS`. If replay fails, fix the issue before proceeding.

**6. Documentation**

If challenge content changed, regenerate the catalog:

```bash
go run scripts/generate_catalog.go > docs/CHALLENGES.md
```

Update `docs/RELEASES.md` with the new version row. Verify all documentation references are consistent.

**7. Stage**

```bash
git add -A
git status
```

Review staged files. Confirm no unintended changes are included.

**8. Commit**

Write a descriptive commit message with title (version tag + one-line summary), body (categorized changes), and discipline section (what did NOT change).

**9. Tag**

```bash
git tag v0.1.<N>
```

The tag version must match the release plan.

**10. Push and Release**

```bash
git push origin v0.1.x v0.1.<N>
gh issue create --title "v0.1.<N>: <title>" --body "<scope>"
gh issue create --title "v0.1.<N>: verify <title>" --body "<checklist>"
```

Create two issues per release: a release issue describing the scope and a verification issue with a checklist. Do not skip steps. Every release follows the same process.
