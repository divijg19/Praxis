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
| `cmd/praxis` | CLI: `describe`, `catalog` (`--ids`/`--core`/`--json`), `next`, `stats` (`--json`), `report` (`--json`), `start`, `hint`, `complete`, `verify`, `reset`, `help`; `trace`/`dev mark-complete` are dev-only |
| `internal/challenge` | `Challenge` struct plus `Layer` and `Verify` constants: the core data model |
| `internal/score` | Only mark computation, from verified keys, par, and hint cost |
| `internal/verify` | Config-free clean-Neovim replay of key logs |
| `internal/content` | `All()`, `DescriptionFor()`, `metadataFor()`, `validStage()`, `IDs()`, `Exists()`, `HintCount`/`HintCost`/`HintCostThrough`/`HintText`, `ValidateAll()`: challenge registry, authored hint ladders, and curriculum metadata |
| `internal/stats` | `Stats`/`Event`, `Load`, `Save`, `Start`, `RevealHint`/`RevealedHintLevel`, `Complete`: sole owner of progress aggregates, revealed hint state, and review selection |

### Ownership notes

- **`internal/stats`** is the only package that computes progress, mastery
  tiers, and review recommendations, and the only owner of the versioned JSON
  file and its validation. It also owns revealed hint state: `RevealHint`
  persists the highest rung, `Start` clears it, and the persisted level is the
  only authority for hint cost when `complete`/`verify` derive marks.
  `cmd/praxis` (`stats` command) is the **sole CLI renderer** of that data.
- **`internal/content`** owns challenge-model validity: `validate()` defines
  it, `ValidateAll()` runs it at CLI startup, and `TestValidateAll` enforces it
  in the gate.
- **Stats-text contract:** `cmd/praxis` emits `Challenges Completed: N/M`
  (where `M` is always the full catalog size). `hub.lua` consumes that line
  directly. Any change to the `stats` command's output text — including
  corruption, invalid-value, and unsupported-version errors — is a breaking
  change for the Neovim frontend and must be coordinated with `hub.lua`.

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
    ├── internal/stats.Load()       →  Progress
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
    Layer      string      // "Tutorial", "Training", or "Trial"
    Hints      []Hint      // authored, opt-in ladder; Text is json:"-"
}

type Hint struct {
    Text string `json:"-"` // never serialized or shown in the challenge buffer
    Cost int    `json:"cost"`
}
```

### Curriculum Metadata

Challenge distribution across stages and layers. Enforced by tests in `internal/content/content_test.go`.

#### Stage taxonomy

| Stage | Purpose | Primitives | Tutorial Challenges |
|---|---|---|---|
| Movement | cursor control | h, j, k, l, UTF-8 navigation | 2 |
| Search | target acquisition | f, w, j, /, ?, ; | 6 |
| Structural Navigation | semantic movement | %, ), i(, a(, i", a", } | 7 |
| Editing | mutation | x, r, ~, dw, ciw, dd, d$ | 7 |
| Text Objects | scoped mutation | diw, daw, di(, da(, di", da", ci(, ci" | 8 |
| Registers | memory | yy, "a, "ap | 3 |

> The "Tutorial Challenges" column is the Tutorial-layer breakdown only (the 33 Tutorial challenges: 11 Core + 22 Additional Lessons). The full per-stage distribution across all layers is the matrix below.

#### Layer taxonomy

| Layer | Purpose (contract) | Scaffolding | Challenges |
|---|---|---|---|
| Tutorial (Core) | teach the essentials; finish onboarding | instruction-led, ordered | 11 |
| Tutorial (Additional Lessons) | teach more; never blocks (internally the `optional` tier) | instruction-led, free | 22 |
| Training | never teaches mechanics; refines fluency | combine primitives | 10 |
| Trial | integrate; solve your own way | select composition | 5 |

#### Distribution matrix

```
                     Tutorial  Training  Trial  Total
Movement             2         0         0      2
Search               6         0         0      6
Structural Nav       7         0         0      7
Editing              7         3         2      12
Text Objects         8         4         3      15
Registers            3         3         0      6
Total                33        10        5      48
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

The Neovim frontend is loaded at startup: `plugin/praxis.lua` calls `require('praxis')`, which registers the entry point. The canonical invocation name is lowercase (`praxis` in the shell, `:praxis` in Nvim through a command-line alias); Vim requires the underlying user command to start uppercase, so `util.command()` derives it from `util.name`. No plugin manager is required, but the repository's `nvim/` directory must be on `runtimepath`.

### Module layout

| Module | Surface | Responsibility |
|---|---|---|
| `init.lua` | none | Command registration, dispatch, binary availability check, recovery, orphan cleanup, first-time detection |
| `challenge.lua` | Practice | Challenge lifecycle: open, verify, autocmds, result, retry, invalid-id recovery |
| `ui.lua` | none | Scratch buffer creation, content helpers, `recovery()` screen |
| `onboarding.lua` | Arrival | First-time welcome flow |
| `hub.lua` | Progress | Hub surface: stats, current location, direction, mastery, return-to-previous-buffer |

### Practice Surface (challenge.lua)

- Fetch challenge data from the `praxis` CLI binary
- Create buffer with challenge content
- Set up autocmds for validation:
  - `CursorMoved`: cursor challenges (target reached check)
  - `TextChanged`: buffer challenges (buffer matches result)
  - `TextChangedI`: buffer challenges (Insert mode edits)
- Checks `state.verify` to decide which autocmd behavior to enable:
  - `"cursor"` → modifiable=false, CursorMoved listener, target check
  - `"buffer"` → modifiable=true, TextChanged + TextChangedI listeners, buffer comparison via `check_buffer()`
   - `"composite"` → same as buffer; success is decided by clean replay
- Uses `byte_to_char()` normalization for multi-byte content:

```lua
function byte_to_char(line, bytecol)
    return vim.fn.strchars(string.sub(line, 1, bytecol))
end
```

Converts Neovim's 0-indexed byte column to a 0-indexed character column. Critical for multi-byte content (UTF-8 Greek, emoji, etc.).

### Reflection (keylog.lua + challenge.lua)

- Ephemeral per-challenge state (`state` table): captured keys, target, verify mode
- On completion, `render_result()` sends the key log to `praxis complete`; Go verifies it in a clean Neovim and persists only verified results
- `internal/verify` owns replay, `internal/score` owns marks, `internal/stats` owns aggregates

See [REFERENCE.md](./REFERENCE.md) for the surface-by-surface walkthrough
(onboarding, hub, result screen) and the full external interface contract
table.

### Runtime Loop

A challenge is a scratch (`nofile`) buffer; exactly one "Praxis" surface
is visible at a time. The loop:

1. `:praxis [<id>]` (or onboarding `s` / hub `<CR>`) → `challenge.open(id)`.
2. `challenge.open` fetches the description from the `praxis` CLI, renders
   `content`, and arms an autocmd: `CursorMoved` for cursor challenges,
   or `TextChanged` / `TextChangedI` for buffer / composite challenges.
3. The learner edits. Each relevant event increments `edits`; `vim.on_key`
   captures physical keys. On a goal match, capture stops.
4. `render_result()` sends keys to `praxis complete`, which replays them in a
   config-free Neovim and returns a verdict. Only verified runs are persisted.
   The buffer is replaced with a result card showing provenance for every value.
5. `r` re-runs the challenge; `<Enter>` opens the next challenge via
   `praxis next` (or returns to the Hub); `q` returns to the Hub.

`ui.create_buffer` removes only Praxis scratch buffers with the same name and
otherwise suffixes the new buffer name ("Praxis 2"). Unrelated user files,
including clean files named "Praxis", are never deleted.
