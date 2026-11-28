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
    Hints:   []challenge.Hint{ // optional authored ladder; Tutorials declare none
        {Text: "A word can be selected from inside.", Cost: 1},
        {Text: "Move to the word before acting.", Cost: 2},
    },
}

type Hint struct {
    Text string `json:"-"`  // never serialized; revealed only by `praxis hint`
    Cost int    `json:"cost"`
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
| `Hints` | EVOLVABLE | Optional authored ladder of 1–3 rungs; Tutorials declare none. Costs must increase; `Text` is `json:"-"` and never crosses the CLI boundary. |

### Canonical Representation

`content.DescriptionFor(id)` is the canonical way to obtain a challenge's full representation. It returns a `Description` struct combining challenge data, curriculum metadata, and verified start/par metadata:

```go
type Description struct {
    challenge.Challenge                // embeds ID, Name, Verify, Layer, Target, Content, Result
    Stage       string   `json:"stage"`
    Concept     string   `json:"concept"`
    Context     string   `json:"context"`
    DerivedFrom []string `json:"derived_from,omitempty"`
    StartRow    int      `json:"start_row"` // play-area row for the canonical solution
    StartCol    int      `json:"start_col"`
    Par         int      `json:"par"` // verified physical-key count of the canonical solution
    // Description embeds Challenge, so hint costs appear in JSON but Text never does.
}
```

All consumers (CLI, Lua frontend, replay tool) must obtain challenge data
through `DescriptionFor`. The `Description` struct is the single source of
truth; JSON is a transport format.

**Stability:** v0.4.x is a clean-break canonical model. There are no legacy readers or migrations; the progress file has a single strict schema.

Enforced by: `TestDescriptionForCompleteness`, `TestDescriptionForUnknown`.

### Adding a New Challenge

1. Add the challenge to the end of `All()` in `internal/content/content.go`
2. Add metadata entry to `curriculum` map in `internal/content/curriculum.go`
3. Append the ID to `stableChallengeIDs` in `internal/content/content_test.go` (order guard); names are evolvable and only uniqueness is enforced
4. Register a canonical keystroke trace in `internal/content/solutions.go`.
5. Run `go run ./tools/verify`. Harnesses derive IDs from `praxis catalog --ids` / `--core` and verify every trace in a clean Neovim.

## CLI Surface

### Public Commands

| Command | Output | Exit | Enforced by |
|---|---|---|---|
| `praxis help` | CLI reference for learner commands | 0 | TestStartAndComplete |
| `praxis` (bare) | Next challenge ID + `praxis help` hint | 0 | TestStartAndComplete |
| `praxis describe <id>` | JSON: `Description` struct (challenge + metadata + start position + par + hint costs) | 0 / 1 | TestDescribeFields, TestDescribeHintsCostOnly |
| `praxis trace <id>` | Internal (dev-only, `PRAXIS_DEV=1`): canonical keystroke trace for replay | 0 / 1 | TestTraceGate |
| `praxis catalog` | One name per line, curriculum order | 0 | TestCatalogOutputStable |
| `praxis catalog --ids` | One challenge ID per line, curriculum order | 0 | TestCatalogFlags |
| `praxis catalog --core` | One Core Tutorial ID per line, sorted | 0 | TestCatalogFlags |
| `praxis next` | Next challenge ID; empty if all complete | 0 | TestStartAndComplete |
| `praxis stats [id]` | Per-challenge or summary; `--json` for full state | 0 / 1 | TestStartAndComplete |
| `praxis reset [--yes]` | Erase all progress; interactive (type `RESET`), or `--yes` for automation | 0 / 1 | TestReset |

### Internal (Frontend Transport API)

| Command | Output | Exit | Enforced by |
|---|---|---|---|
| `praxis start <id>` | Silent; records effort only | 0 / 1 | TestStartAndComplete |
| `praxis complete <id>` | stdin key-log JSON → verified verdict + marks JSON | 0 / 1 | TestStartAndComplete, TestTaintedCompletionRejected |
| `praxis verify <id>` | same replay, diagnostic only, no write | 0 / 1 | TestStartAndComplete, TestHintCostDerivedAndSpoofProof |
| `praxis stats --json` | full canonical progress JSON | 0 | TestCatalogJSONHasNoSpoilers |
| `praxis catalog --json` | metadata-only catalog JSON | 0 | TestCatalogJSONHasNoSpoilers |
| `praxis dev mark-complete <id> [count]` | dev-only harness seeding (`PRAXIS_DEV=1`) | 0 / 1 | TestTraceGate |
| `praxis hint <id> <level> [--preview]` | Authored hint: cost-only preview or explicit text + cost reveal | 0 / 1 | TestHintPreviewAndReveal, TestHintCostDerivedAndSpoofProof |

On unknown ID, stderr is `unknown challenge: <id>`.


### `praxis stats [id]`

With an ID, shows per-challenge stats:

```
Starts: 10
Completions: 10
Best Keys: 5
Best Edits: 1
Best Marks: 3/3
Mastery: Experienced
```

Without arguments, shows summary with mastery distribution and practice guidance. Exit 0 on success, 1 on unknown challenge ID.

### `praxis describe <id>`

Returns a canonical JSON representation of the challenge including its metadata and verified start position/par.

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
  "concept": "diw",
  "context": "composite deletion",
  "stage": "Text Objects",
  "start_row": 3,
  "start_col": 0,
  "par": 5
}
```

Enforced by: `TestDescribeFields`, `TestCatalogOutputStable`.

### `praxis catalog`

Returns the name of every registered challenge, one per line, in curriculum order.

Enforced by: `TestCatalogOutputStable`.

### `praxis next`

Returns the first challenge in curriculum order that is not Practiced (i.e., completions ≤ `learningMax` (2)). Output is the challenge ID on stdout, or empty if all challenges are Practiced.

Enforced by: `TestStartAndComplete`, internal content integrity tests

## Validators

Validators determine how a challenge's success condition is evaluated. Praxis has three: **cursor**, **buffer**, and **composite**. The set of valid values is defined once and enforced for every challenge by `validate()` (`internal/content/validate.go`). The CLI calls `ValidateAll()` at startup and the gate runs `TestValidateAll` in tests.

### Cursor Validator

**Verify value:** `"cursor"`

**Rule:** Navigate the cursor to the target character on screen.

Buffer is set to `Content` with `modifiable=false`. A `CursorMoved` autocommand checks if the character under the cursor matches `Target`. Uses `byte_to_char()` normalization to handle multi-byte content correctly. On match: sets `state.done = true`, echoes success.

**Used by:** 11 challenges (movement, search, navigation, UTF-8 proof).

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
| grid_rush | `Move your cursor to the star ★` | ★ | On the ★ character |
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

Buffer is set to `Content` with `modifiable=true`. A `TextChanged` autocommand fires on Normal mode edits and a `TextChangedI` autocommand fires on Insert mode edits; both are informational; the verifier owns the result. `check_buffer()` compares every line of the current buffer against `Result` using byte-exact string equality. On match: sets `state.done = true`, echoes success.

**Used by:** 22 challenges (editing, text objects, registers).

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
| register_replace_hunter | `correct` + `wrong` | `"ayy` to store, `"ap` to replace |

**FAIL examples:**

| Scenario | Why it fails |
|---|---|
| Buffer has wrong content | At least one line differs from result |
| Line count mismatch | Extra or missing line |
| Trailing whitespace differs | `"keep "` vs `"keep"`: byte-exact comparison |
| Result is empty | Result field is nil or `[]string{}` (violates contract) |

### Composite Validator

**Verify value:** `"composite"`

**Rule:** Edit the buffer content to exactly match `Result`.

Behaviorally identical to the buffer validator (byte-exact comparison via `check_buffer()`). `praxis complete` replays the key log in a clean Neovim and decides success; the frontend does not determine the result. Verified keys, par, and marks are reported with provenance on the result card.

**Used by:** 15 challenges (10 Training + 5 Trial).

### Autocommand Dispatching

The Lua frontend only attaches `CursorMoved` for cursor challenges. Buffer challenges use `TextChanged` + `TextChangedI` instead. Handled by checking `state.verify` when setting up autocommands.

### Contract Enforcement

Every challenge uses one of the three verify values, and each value carries the correct shape: cursor challenges declare a `Target` and no `Result`; buffer and composite challenges declare a `Result` and no `Target`. These invariants are enforced by runtime `ValidateAll()` and the content test suite (see TESTING.md).

## Stats

### Schema

```go
type Stats struct {
    Starts      int    // times opened
    Completions int    // verified completions only
    BestKeys    int    // fewest physical keys in a verified run
    BestEdits   int    // fewest edit events in a verified run
    BestMarks   int    // best verified mark, 1..3
    LastPlayed  string // "2006-01-02" format
}
```

Progress is stored in `$XDG_DATA_HOME/praxis/progress.json` (or `~/.local/share/praxis/progress.json`) as one strict document:

```json
{
  "format": 1,
  "entries": {
    "grid_rush": {
      "starts": 2,
      "completions": 1,
      "best_keys": 3,
      "best_edits": 1,
      "best_marks": 3,
      "last_played": "2026-11-27"
    }
  },
  "events": [],
  "revealed_hints": { "grid_rush": 2 }
}
```

There is exactly one format. Legacy formats are not read and are not migrated. Bad JSON, an unknown format, or impossible counters are treated as corruption, backed up to `progress.json.corrupt`, and surfaced as recovery. Writes are flushed with `fsync` before an atomic rename.

| File state | Behavior |
|---|---|
| Missing | Fresh, empty progress |
| Format 1 with valid counters | Loaded normally |
| Bad JSON or impossible counters | CLI exits 1, writes `progress.json.corrupt`, Lua shows recovery |
| Unknown format | Treated as corrupt; CLI exit 1, backup, Lua recovery |

State remains a human-readable JSON file under XDG. A database-backed store was evaluated and rejected: it would add dependencies, migration risk, and new failure modes without serving an educational contract. This is a deliberate architecture decision, not an accident.

| Signal | Meaning |
|---|---|
| Starts | Times the challenge was opened |
| Completions | Verified completions only |
| Best Keys | Fewest physical keys in a verified run |
| Best Edits | Verified edit count from clean replay |
| Best Marks | Best verified mark (1-3) |
| Mastery | Practice depth (derived from Completions) |
| Guidance | Next action (derived from Mastery + LastPlayed + curriculum order) |

### Mastery Tiers

Derived from `Completions`:

| Tier | Threshold |
|---|---|
| Unseen | 0 |
| Learning | 1 to 2 |
| Practiced | 3 to 7 |
| Experienced | 8+ |

### Design Decisions

- Only verified completions are recorded. There is no client-reported completion path.
- Best Keys and Best Edits track minima; Best Marks tracks the maximum.
- Marks are computed in Go from the verified key log, par, and hint cost; clients cannot supply them.
- `revealed_hints` (optional) holds the highest revealed rung per challenge. It
  is the sole authority for hint cost; `praxis hint <id> <level> --preview` is
  inert, reveal persists the level, and `praxis start` clears it on retry.
- Starts are effort only and never become completions.
- Events are local, capped, and compacted without dropping aggregates.
- `Event.Seed` is **reserved for future rounds**; the current CLI never reads
  or writes it, so it must not be treated as part of the v0.4.x contract.

### Practice Guidance

`NextChallenge()` returns the first curriculum-ordered challenge whose Completions ≤ 2 (Unseen or Learning). This means: **finish what you started**. A partially-practiced Learning challenge will be recommended before a new Unseen challenge.

`RecommendedReview()` returns the oldest Practiced challenge by LastPlayed date, falling back to the oldest Experienced if no Practiced challenges exist. Practiced challenges are preferred because they are more likely to benefit from review than deeply-ingrained Experienced ones.

## Reflection (Internal)

Ephemeral per-challenge state (edits, elapsed time, captured keys) lives in `challenge.lua` and is sent to `praxis complete`. Verified results are persisted; unverified runs write nothing.

- Raw key logs are never persisted.
- Marks are descriptive of one run, not a score.

## Hub

The Hub is the primary surface for returning users. It answers "where am I and what should I do next?" and is opened automatically on `:praxis` for non-first-time users.

### Layout

```
── Praxis ──────────────────────────────────────

  Current: Tutorial / Search
  So far: 4 challenges completed.

  Direction:
    Next: Find Hunter / Search
    Review: Grid Rush / Movement

  Practice mix:
    Unseen: 44   Learning: 2   Practiced: 1   Experienced: 1

  [Enter] Continue, or [r] Review.
  [q] Back.
```

Three sections:

| Section | Content | Source |
|---|---|---|
| **Current** | Current stage (stage of first un-Practiced challenge) | `praxis next`, `praxis describe` |
| **Progress** | Positive completion statement; no denominator. `So far: N challenges completed.` | `praxis stats` |
| **Direction** | Next challenge (by name) + stage, review recommendation (by name) + stage | `praxis next`, `praxis describe`, `praxis stats` |
| **Practice mix** | Mastery distribution, shown only after five completions | `praxis stats` |

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
- **Zero Go changes.** All data is sourced from `praxis next`, `praxis stats`,
  and `praxis describe`.

## Integrity Guarantees

### Principles

1. **One concept per challenge.** Every challenge has a documented primary Concept, Context, and Stage. The `curriculum` map in `internal/content/curriculum.go` is the single source of truth.
2. **IDs are permanent.** Challenge identifiers must never be renamed or removed.
3. **Verify values are stable.** The set of valid Verify values (`cursor`, `buffer`, `composite`) is fixed and every challenge uses one.
4. **Curriculum growth is intentional.** New challenges require a documented primary concept and must not duplicate existing content.
5. **Replay verification is mandatory.** Every challenge must appear in replay coverage, execute its canonical Vim keystroke trace, reach the goal state.
6. **Duplicate challenge content is prohibited.** Challenges may share concepts when they teach the concept in a different context or composition.

### Test Suite

These principles are enforced by the integrity and content test suites (see TESTING.md): unique and stable challenge IDs, a complete and acyclic curriculum map, no duplicate challenge content, full replay coverage, and single-owner validators for the verify/result/target shape.

### Anti-Goals

- Challenge difficulty ratings
- Mastery scoring
- Skill trees or dependencies
- Progression gates

