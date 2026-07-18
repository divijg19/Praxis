# Reference

## Challenge Model

### Struct

Every challenge is defined in `internal/content/content.go` as part of the `All()` function:

```go
challenge.Challenge{
    ID:      "unique_identifier",
    Name:    "Display Name",
    Verify:  "cursor",         // "cursor", "buffer", or "composite"
    Target:  "★",              // cursor challenges only; empty for buffer/composite
    Content: []string{
        "Instruction line",
        "",
        "play area content",
    },
    Result:  []string{         // buffer/composite challenges only; nil for cursor
        "Instruction line",
        "",
        "edited content",
    },
    Layer:   "Tutorial",       // "Tutorial", "Training", or "Trial"
    Evaluation: &challenge.Evaluation{MaxMoves: 10}, // composite only; nil otherwise
}
```

### Field Rules

| Field | Stability | Rule |
|---|---|---|
| `ID` | STABLE | Globally unique. Never change after first release. `snake_case`: Tutorial/Training use `_hunter` or `_combo` suffix; Trials use a `trial_` prefix. |
| `Name` | EVOLVABLE | Display name. May change as curriculum framing shifts. |
| `Verify` | STABLE | `"cursor"`, `"buffer"`, or `"composite"`. Must be a valid Verify value. |
| `Target` | STABLE | Required for cursor (non-empty). Empty for buffer/composite. |
| `Content` | STABLE | First line is instruction. Buffer/composite: 3+ lines (instruction, blank, play area). Cursor: 1+ lines. |
| `Result` | STABLE | Required for buffer/composite (non-empty slice, exact target state). Nil for cursor. |
| `Layer` | STABLE | `"Tutorial"`, `"Training"`, or `"Trial"`. |
| `Evaluation` | STABLE | Non-nil only for `"composite"`. Contains `MaxMoves` threshold. |

### Canonical Representation

`content.DescriptionFor(id)` is the canonical way to obtain a challenge's full representation. It returns a `Description` struct combining challenge data, curriculum metadata, and (for composite challenges) evaluation:

```go
type Description struct {
    challenge.Challenge                // embeds ID, Name, Verify, Layer, Target, Content, Result, Evaluation
    Stage       string   `json:"stage"`
    Concept     string   `json:"concept"`
    Context     string   `json:"context"`
    DerivedFrom []string `json:"derived_from,omitempty"`
}
```

All consumers (CLI, Lua frontend, replay tool) must obtain challenge data
through `DescriptionFor`. The `Description` struct is the single source of
truth; JSON is a transport format.

**Stability:** The JSON field names and types are stable for the v0.3.x series (frozen). New fields may be added but existing fields will not be renamed, removed, or have their types changed.

Enforced by: `TestDescriptionForCompleteness`, `TestDescriptionForUnknown`.

### Adding a New Challenge

1. Add the challenge to the end of `All()` in `internal/content/content.go`
2. Add metadata entry to `curriculum` map in `internal/content/curriculum.go`
3. Append the ID to `stableChallengeIDs` and the name to `stableChallengeNames` in `internal/content/content_test.go` (rename guard)
4. Append the ID to the `all_ids` list in `tools/replay/replay.lua` (guarded by `TestReplayCoverage`)
5. Run `tools/verify.sh`. The total challenge counts in the unit tests and the journey harness are derived from `content.All()`, so no manual count edits are required.

## CLI Surface

### Public Commands

| Command | Output | Exit | Enforced by |
|---|---|---|---|
| `praxis help` | CLI reference, all commands listed with descriptions | 0 | TestHelpCommand |
| `praxis` (bare) | Next challenge ID + `praxis help` hint | 0 | TestBarePraxis |
| `praxis describe <id>` | JSON: `Description` struct (challenge + metadata + evaluation) | 0 / 1 | TestDescribeCommand, TestDescribeComposite, TestDescribeUnknown |
| `praxis catalog` | One name per line, curriculum order | 0 | TestCatalogOutputStable |
| `praxis next` | Next challenge ID; empty if all complete | 0 | TestNextCommand |
| `praxis stats [id]` | Per-challenge or summary | 0 / 1 | TestStatsCommand, TestStatsSummary, TestRecordStats, TestStatsCommandConfidenceLevels |
| `praxis reset [--yes]` | Erase all progress; interactive (type `RESET`), or `--yes` for automation | 0 / 1 | TestResetCommand |

### Internal (Frontend Transport API)

| Command | Output | Exit | Enforced by |
|---|---|---|---|
| `praxis attempt <id>` | Silent (no stdout) | 0 / 1 | TestAttemptCommand, TestAttemptUnknown |
| `praxis record <id> <moves> <time_ms>` | Silent (no stdout) | 0 | TestRecordStats |

On unknown ID, stderr is `unknown challenge: <id>`.


### `praxis stats [id]`

With an ID, shows per-challenge stats:

```
Attempts: 10
Completions: 10
Success Rate: 100%
Best Moves: 2
Best Time: 180 ms
Mastery: Experienced
Confidence: High
```

Without arguments, shows summary with mastery distribution and practice guidance. Exit 0 on success, 1 on unknown challenge ID.

### `praxis describe <id>`

Returns a canonical JSON representation of the challenge including its metadata and (for composite challenges) evaluation.

**Output schema:**
```json
{
  "id": "find_diw_combo",
  "name": "Find + Delete Word",
  "verify": "composite",
  "layer": "Training",
  "target": "",
  "content": ["Delete the word using f<char> and diw", "", "keep delete keep"],
  "result": ["Delete the word using f<char> and diw", "", "keep  keep"],
  "concept": "f",
  "context": "composite deletion",
  "stage": "Text Objects",
  "evaluation": {
    "max_moves": 8
  }
}
```

The `evaluation` field is emitted only for `"composite"` challenges (via `json:"evaluation,omitempty"`). Non-composite challenges omit this field entirely.

Enforced by: `TestDescribeCommand`, `TestDescribeComposite`, `TestDescribeUnknown`.

### `praxis catalog`

Returns the name of every registered challenge, one per line, in curriculum order.

Enforced by: `TestCatalogOutputStable`.

### `praxis next`

Returns the first challenge in curriculum order that is not Practiced (i.e., completions ≤ LearningMax=2). Output is the challenge ID on stdout, or empty if all challenges are Practiced.

Enforced by: `TestNextCommand`, `TestNextCommandAfterCompletion`, `TestNextCommandComplete`

## Validators

Validators determine how a challenge's success condition is evaluated. Praxis has three: **cursor**, **buffer**, and **composite**. The set of valid values is defined once and enforced for every challenge by `Validate()` (`internal/content/validate.go`), which runs over all challenges via `TestValidateAll`.

### Cursor Validator

**Verify value:** `"cursor"`

**Rule:** Navigate the cursor to the target character on screen.

Buffer is set to `Content` with `modifiable=false`. A `CursorMoved` autocommand checks if the character under the cursor matches `Target`. Uses `byte_to_char()` normalization to handle multi-byte content correctly. On match: sets `state.done = true`, echoes success.

**Used by:** 18 challenges (movement, search, navigation, UTF-8 proof).

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

**Used by:** 19 challenges (editing, text objects, registers).

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
| Trailing whitespace differs | `"keep "` vs `"keep"`: byte-exact comparison |
| Result is empty | Result field is nil or `[]string{}` (violates contract) |

### Composite Validator

**Verify value:** `"composite"`

**Rule:** Edit the buffer content to exactly match `Result`, within a maximum number of moves.

Behaviorally identical to the buffer validator (TextChanged + TextChangedI listeners, byte-exact comparison via `check_buffer()`). Additionally, the Lua frontend enforces a `MaxMoves` threshold:

- On each Normal-mode edit (TextChanged), moves counter is incremented
- If moves exceed `MaxMoves`, challenge is failed with "Over the move limit. Press [r] to retry."
- Result screen shows `Moves: N / MaxMoves` instead of just `Moves: N`

The `MaxMoves` value is transmitted via the `describe` JSON endpoint as `evaluation.max_moves`.

**Used by:** 15 challenges (10 Training + 5 Trial).

### Autocommand Dispatching

The Lua frontend only attaches `CursorMoved` for cursor challenges. Buffer challenges use `TextChanged` + `TextChangedI` instead. Handled by checking `state.verify` when setting up autocommands.

### Contract Enforcement

Every challenge uses one of the three verify values, and each value carries the correct shape: cursor challenges declare a `Target` and no `Result`; buffer and composite challenges declare a `Result` and no `Target`. Composite challenges always declare a positive `MaxMoves`. These invariants are enforced by the content test suite (see TESTING.md).

## Stats

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

Stored in `~/.local/share/praxis/stats.json` (`$XDG_DATA_HOME/praxis/stats.json`). Simple JSON, human-readable and editable. No database, no migrations.

| Signal | Meaning |
|---|---|
| Attempts | Effort invested |
| Completions | Raw success count |
| Success Rate | Reliability metric (Completions / Attempts) |
| Mastery | Practice depth (derived from Completions) |
| Confidence | Execution reliability (derived from Success Rate) |
| Guidance | Next action (derived from Mastery + LastPlayed + curriculum order) |

### Confidence ≠ Mastery

These are orthogonal dimensions. High Confidence does not imply high Mastery, and vice versa.

| | High Confidence | Low Confidence |
|---|---|---|
| High Mastery | 20 attempts, 16 completions (80%) | 20 attempts, 10 completions (50%) |
| Low Mastery | 1 attempt, 1 completion (100%) | 1 attempt, 0 completions (0%) |

### Mastery Tiers

Derived from `Completions`:

| Tier | Threshold |
|---|---|
| Unseen | 0 |
| Learning | 1 to 2 |
| Practiced | 3 to 7 |
| Experienced | 8+ |

### Design Decisions

- Stats are recorded only on successful completion
- Best values track minima (fewest moves, fastest time)
- First completion always sets the best value
- Subsequent completions only update if better
- No per-attempt history. Only the best is preserved
- Attempts and Completions are tracked independently. Attempts increment on challenge start and replay; Completions increment only on success
- Confidence is derived from SuccessRate at display time. Not stored in JSON. Thresholds: 80 percent or higher is High, 60 percent or higher is Medium, and below 60 percent is Low. The display shows an em dash when there are no attempts, which is a no-data signal and not the same as Low.
- Confidence and Mastery are orthogonal. Mastery answers "how much have I practiced?"; Confidence answers "how reliably am I executing?"

### Practice Guidance

`NextChallenge()` returns the first curriculum-ordered challenge whose Completions ≤ 2 (Unseen or Learning). This means: **finish what you started**. A partially-practiced Learning challenge will be recommended before a new Unseen challenge.

`RecommendedReview()` returns the oldest Practiced challenge by LastPlayed date, falling back to the oldest Experienced if no Practiced challenges exist. Practiced challenges are preferred because they are more likely to benefit from review than deeply-ingrained Experienced ones.

## Reflection (Internal)

Per-challenge tracking is an internal implementation detail. There is no public command or surface for it. Ephemeral per-challenge counters (moves, elapsed time) live in the `state` table inside `challenge.lua` and are aggregated via the `praxis record` CLI on completion.

- **No persistence.** Discarded when the challenge buffer closes.
- **No public command.** No `:PraxisSession`.

## Hub

The Hub is the primary surface for returning users. It answers "where am I and what should I do next?" and is opened automatically on `:Praxis` for non-first-time users.

### Layout

```
── Praxis ──────────────────────────────────────

  Current: Tutorial / Search
  Progress: 4/52

  Direction:
    Next: Find Hunter / Search
    Review: Motion Rush / Movement

  Mastery:
    Unseen: 48   Learning: 2   Practiced: 1   Experienced: 1

  [Enter] Continue, or [r] Review.
  [q] Back.
```

Three sections:

| Section | Content | Source |
|---|---|---|
| **Current** | Current stage (stage of first un-Practiced challenge), progress count | `praxis next`, `praxis describe`, `praxis stats` |
| **Direction** | Next challenge (by name) + stage, review recommendation (by name) + stage | `praxis next`, `praxis describe`, `praxis stats` |
| **Mastery** | Mastery distribution (compact one-line) | `praxis stats` |

### Actions

| Key | Action |
|---|---|
| `<CR>` | Opens the next challenge directly (never opens Hub as intermediate step) |
| `r` | Opens the recommended review challenge (only shown when a review is recommended) |
| `q` | Returns to the previous buffer |

### Design Principles

- **Hub is for orientation.** You visit it when you need context (starting Praxis, returning after a break, wondering where you are).
- **Result screen is for momentum.** After completing a challenge, Enter goes directly to the next challenge. It never goes through the Hub.
- **Practice flow:** Challenge → Result → Next Challenge. Navigation only when you explicitly choose it.
- **Zero Go changes.** All data is sourced from `praxis next`, `praxis stats`.

## Integrity Guarantees

### Principles

1. **One concept per challenge.** Every challenge has a documented primary Concept, Context, and Stage. The `curriculum` map in `internal/content/curriculum.go` is the single source of truth.
2. **IDs are permanent.** Challenge identifiers must never be renamed or removed.
3. **Verify values are stable.** The set of valid Verify values (`cursor`, `buffer`, `composite`) is fixed and every challenge uses one.
4. **Curriculum growth is intentional.** New challenges require a documented primary concept and must not duplicate existing content.
5. **Replay verification is mandatory.** Every challenge must be solvable via the replay harness.
6. **Duplicate challenge content is prohibited.** Challenges may share concepts when they teach the concept in a different context or composition.

### Test Suite

These principles are enforced by the integrity and content test suites (see TESTING.md): unique and stable challenge IDs, a complete and acyclic curriculum map, no duplicate challenge content, full replay coverage, and single-owner validators for the verify/result/target shape.

### Anti-Goals

- Challenge difficulty ratings
- Mastery scoring
- Skill trees or dependencies
- Progression gates


## Release Procedure

1. **Verify**: run `tools/verify.sh`. It runs build, lint (`go run github.com/golangci/golangci-lint/cmd/golangci-lint run`), format, vet, tests, replay, and journey. All checks must pass.
2. **Build**: run `go build ./...` so all packages compile.
3. **Documentation**: if content changed, update the relevant doc under `docs/` (the challenge catalog is available at runtime via `praxis catalog`).
4. **Stage**: run `git add -A && git status` to verify staged files.
5. **Commit**: write a descriptive message with a title (version plus summary), a body (categorized changes), and a discipline section (what did NOT change).
6. **Tag**: run `git tag <version>` so it matches the release plan.
7. **Push and Release**: run `git push origin <branch> <version>`, then create release and verification issues on GitHub.

Every release follows the same process. Do not skip steps.
