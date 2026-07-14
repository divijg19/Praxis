# Architecture

## Overview

```
              Praxis

          ┌─────────────┐
          │  Go Engine  │
          └──────┬──────┘
                 │
            ┌─────┬─────┐
            │           │
            ▼           ▼
      ┌──────────┐ ┌──────────┐
      │   CLI    │ │  Neovim  │
      │ (Go cmd) │ │  (Lua)   │
      └──────────┘ └──────────┘
```

One engine. Multiple frontends. Shared progression, content, scoring, and persistence.

## Go Engine

The engine is a single-module Go project under `github.com/divijg19/Praxis`.

### Package layout

| Package | Responsibility |
|---|---|---|
| `cmd/praxis` | CLI entry point: `describe`, `catalog`, `next`, `stats`, `reset` (public); `attempt`, `record` (internal) |
| `internal/challenge` | `Challenge` struct — the core data model, `Evaluation` for composite challenges |
| `internal/content` | `All()`, `DescriptionFor()`, `metadataFor()`, `validStages()`, `IDs()`, `Exists()` — challenge registry + curriculum metadata |
| `internal/stats` | `Stats` struct, `Load`, `Save`, `Update` — persistent progress tracking |

### Data flow

```
User request (CLI)
    |
    v
cmd/praxis/main.go
    |
    ├── internal/content.All()      →  challenge.Challenge[]
    ├── internal/content.DescriptionFor(id)  →  (Description, bool)
     ├── internal/content.metadataFor(id)  →  content.Metadata{Concept, Context, Stage}
    ├── internal/stats.Load()       →  map[string]Stats
    |
    v
stdout (CLI response)
```

### The `Challenge` struct

```go
type Challenge struct {
    ID         string      // stable public identifier, never changes
    Name       string      // display name, may evolve
    Verify     string      // "cursor", "buffer", or "composite"
    Target     string      // for cursor: the character to navigate to
    Content    []string    // buffer content lines
    Result     []string    // for buffer/composite: target buffer state
    Layer      string      // "Tutorial", "Training", "Trial", "Boss"
    Evaluation *Evaluation // non-nil only for "composite" challenges
}

type Evaluation struct {
    MaxMoves int // anti-bruteforce threshold
}
```

### Curriculum Metadata

Challenge distribution across stages and layers. Enforced by tests in `internal/content/content_test.go`.

#### Stage taxonomy

| Stage | Purpose | Primitives | Tutorial Challenges |
|---|---|---|---|
| Movement | cursor control | h, j, k, l, UTF-8 navigation | 3 |
| Search | target acquisition | f, F, w, W, /, ?, ; | 7 |
| Structural Navigation | semantic movement | %, ), (, {, }, i(, a(, i[, a[, i", a" | 10 |
| Editing | mutation | x, r, ~, dw, ciw, dd, D | 7 |
| Text Objects | scoped mutation | diw, daw, di(, da(, di", da", ciw, ci(, ci" | 9 |
| Registers | memory | yy, "a, "A, "ap | 5 |

> The "Tutorial Challenges" column is the Tutorial-layer breakdown only (the 41 Tutorial challenges). The full per-stage distribution across all layers is the matrix below.

#### Layer taxonomy

| Layer | Purpose | Scaffolding | Challenges |
|---|---|---|---|
| Tutorial | primitive introduction | Observe → Practice → Apply | 41 |
| Training | composition formation | combine primitives, MaxMoves constraint | 10 |
| Trial | recognition under pressure | select correct composition, budget enforcement | 5 |
| Boss | mastery | (deferred to v0.3) | 0 |

#### Distribution matrix

```
                     Tutorial  Training  Trial  Total
Movement             3         0         0      3
Search               7         0         0      7
Structural Nav       10        0         0      10
Editing              7         3         2      12
Text Objects         9         4         3      16
Registers            5         3         0      8
Total                41        10        5      56
```

#### DerivedFrom lineage

Trials declare the Training challenges they derive from:

| Trial | DerivedFrom | Composition tested |
|---|---|---|
| trial_find_delete | find_diw_combo | f + diw |
| trial_find_change | find_ca_quote_combo | f + ca" |
| trial_dot_repeat | dw_dot_combo | dw + . |
| trial_delete_choice | find_diw_combo, find_daw_combo | diw vs daw |
| trial_repeat_choice | dw_dot_combo, ciw_dot_combo | . vs re-execute |

Enforced by `TestTrialIntegrity` (all targets exist) and `TestDerivedFromAcyclic` (no cycles).

## Neovim Frontend (Lua)

The Neovim frontend is loaded on demand — `:Praxis` triggers `require('praxis')` via the plugin manager's `cmd` lazy-loading.

### Module layout

| Module | Surface | Responsibility |
|---|---|---|
| `init.lua` | — | Command registration, dispatch, binary availability check, recovery, orphan cleanup, first-time detection |
| `challenge.lua` | Practice | Challenge lifecycle: open, verify, autocmds, result, retry, invalid-id recovery |
| `ui.lua` | — | Scratch buffer creation, content helpers, `recovery()` screen |
| `onboarding.lua` | Arrival | First-time welcome flow |
| `hub.lua` | Progress | Hub surface — stats, current location, direction, mastery, return-to-previous-buffer |

### Practice Surface (challenge.lua)

- Fetch challenge data from the `praxis` CLI binary
- Create buffer with challenge content
- Set up autocmds for validation:
  - `CursorMoved` — cursor challenges (target reached check)
  - `TextChanged` — buffer challenges (buffer matches result, Normal mode edits increment moves)
  - `TextChangedI` — buffer challenges (Insert mode keystrokes, no move increment)
- Checks `state.verify` to decide which autocmd behavior to enable:
  - `"cursor"` → modifiable=false, CursorMoved listener, target check
  - `"buffer"` → modifiable=true, TextChanged + TextChangedI listeners, buffer comparison via `check_buffer()`
   - `"composite"` → same as buffer, plus MaxMoves enforcement; echoes "Over the move limit — press [r] to retry." on exceed
- Uses `byte_to_char()` normalization for multi-byte content:

```lua
function byte_to_char(line, bytecol)
    return vim.fn.strchars(string.sub(line, 1, bytecol))
end
```

Converts Neovim's 0-indexed byte column to a 0-indexed character column. Critical for multi-byte content (UTF-8 Greek, emoji, etc.).

### Reflection (inline in challenge.lua)

- Ephemeral per-challenge state (`state` table): moves, elapsed time, target, verify mode
- On completion, `render_result()` persists via `praxis record` CLI and aggregates counters
- No separate session module — tracking lives in the challenge lifecycle

### Arrival Surface (onboarding.lua)

- First-time detection via stats file absence
- Welcome buffer with orientation text
- `[s]` → `challenge.open("motion_rush")`
- No persistence, no state, no wizard

### Hub Surface (hub.lua)

- Returning users see Hub on `:Praxis` instead of raw CLI output
- Layout: Current → Direction → Mastery
  - **Current:** current stage (stage of first un-Practiced challenge), progress count
  - **Direction:** next challenge + its stage, review recommendation + its stage
  - **Mastery:** mastery distribution (compact one-line format)
- `<CR>` opens the next challenge directly via `challenge.open()`
- Data sourced from `praxis next`, `praxis describe`, `praxis stats`

### Journey Surface (challenge.lua — result screen)

- Result screen displays stats with action options:
  - `[r] Retry.` — reset and retry the challenge
  - `[Enter] Continue.` — opens the next curriculum challenge
  - `[q] Back.` — returns to the Hub
- Enter on result screen skips Hub entirely (Challenge → Result → Next Challenge)

### Runtime Loop

A challenge is a scratch (`nofile`) buffer; exactly one "Praxis" surface
is visible at a time. The loop:

1. `:Praxis [<id>]` (or onboarding `s` / hub `<CR>`) → `challenge.open(id)`.
2. `challenge.open` fetches the description from the `praxis` CLI, renders
   `content`, and arms an autocmd: `CursorMoved` for cursor challenges,
   or `TextChanged` / `TextChangedI` for buffer / composite challenges.
3. The learner edits. Each relevant event increments `moves` and, on a
   successful verify, calls `render_result()`.
4. `render_result()` records the attempt (`praxis record`), refreshes stats
   (`praxis stats`), and replaces the buffer with the result view
    (`[r] Retry.`, `[Enter] Continue.`, `[q] Back.`).
5. `r` re-runs the challenge; `<Enter>` opens the next challenge via
   `praxis next` (or returns to the Hub); `q` returns to the Hub.

`ui.create_buffer` reuses or uniquifies the buffer name so surfaces that
share the "Praxis" label never collide.

## External interface contracts

| Interface | Consumer | Format |
|---|---|---|
| `praxis help` | CLI | CLI reference, commands listed with descriptions |
| `praxis` (bare) | CLI | Next challenge ID + `praxis help` hint |
| `praxis describe <id>` | Neovim, replay | JSON: `Description` struct (challenge + metadata + evaluation) |
| `praxis catalog` | CLI | One challenge name per line (flat, no grouping) |
| `praxis next` | Neovim, hub | Challenge ID or empty |
| `praxis attempt <id>` | Neovim | Silent (internal) |
| `praxis record <id> <moves> <time_ms>` | Neovim | Silent (internal) |
| `praxis stats [id]` | CLI | Per-challenge or summary |
| `praxis reset [--yes]` | CLI | Wipe all progress (confirm or pass `--yes`) |
