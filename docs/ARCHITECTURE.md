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
| `cmd/praxis` | CLI entry point: `list`, `challenge`, `target`, `verify`, `result`, `attempt`, `record`, `stats`, `next`, `stage` |
| `internal/challenge` | `Challenge` struct — the core data model |
| `internal/content` | `All()`, `MetadataFor()`, `ValidStages()` — challenge registry + curriculum metadata |
| `internal/stats` | `Stats` struct, `Load`, `Save`, `Update` — persistent progress tracking |
| `internal/validator` | `Exists(name)` — validator dispatch registry |

### Data flow

```
User request (CLI)
    |
    v
cmd/praxis/main.go
    |
    ├── internal/content.All()  →  challenge.Challenge{ID, Name, Verify, Target, Content, Result, Layer}
    ├── internal/content.MetadataFor(id)  →  content.Metadata{Concept, Context, Stage, Layer}
    ├── internal/stats.Load()  →  map[string]Stats
    |
    v
stdout (CLI response)
```

### The `Challenge` struct

```go
type Challenge struct {
    ID      string     // stable public identifier, never changes
    Name    string     // display name, may evolve
    Verify  string     // "cursor" or "buffer"
    Target  string     // for cursor: the character to navigate to
    Content []string   // buffer content lines
    Result  []string   // for buffer: target buffer state
    Layer   string     // "Tutorial", "Training", "Trial", "Boss"
}
```

## Neovim Frontend (Lua)

The Neovim frontend is loaded via `nvim/plugin/praxis.lua` which requires `nvim/lua/praxis/init.lua`.

### Module layout

| Module | Surface | Responsibility |
|---|---|---|
| `init.lua` | — | Command registration, dispatch, first-time detection |
| `challenge.lua` | Practice | Challenge lifecycle: open, verify, autocmds, result, replay |
| `session.lua` | Reflection | Session tracking: start, record, `:PraxisSession` |
| `ui.lua` | — | Scratch buffer creation and content helpers |
| `onboarding.lua` | Arrival | First-time welcome flow |
| `hub.lua` | Journey | Hub surface (reserved for future use) |

### Practice Surface (challenge.lua)

- Fetch challenge data from the Go binary (`/tmp/praxis`)
- Create buffer with challenge content
- Set up autocmds for validation:
  - `CursorMoved` — cursor challenges (target reached check)
  - `TextChanged` — buffer challenges (buffer matches result, Normal mode edits increment moves)
  - `TextChangedI` — buffer challenges (Insert mode keystrokes, no move increment)
- Checks `state.verify` to decide which autocmd behavior to enable:
  - `"cursor"` → modifiable=false, CursorMoved listener, target check
  - `"buffer"` → modifiable=true, TextChanged + TextChangedI listeners, buffer comparison via `check_buffer()`
- Uses `byte_to_char()` normalization for multi-byte content:

```lua
function byte_to_char(line, bytecol)
    return vim.fn.strchars(string.sub(line, 1, bytecol))
end
```

Converts Neovim's 0-indexed byte column to a 0-indexed character column. Critical for multi-byte content (UTF-8 Greek, emoji, etc.).

### Reflection Surface (session.lua)

- Ephemeral session state (moves, time, completion counts)
- `session.start()` — initializes or continues session tracking per challenge
- `session.record()` — persists completion via `praxis record` CLI, aggregates counters
- `session.show()` — renders `:PraxisSession` buffer

### Arrival Surface (onboarding.lua)

- First-time detection via stats file absence
- Welcome buffer with orientation text
- Enter → `challenge.open("motion_rush")`
- No persistence, no state, no wizard

### Hub Surface (hub.lua)

- Reserved for v0.2.2 (Journey release)
- Currently a stub that errors on invocation

## External interface contracts

| Interface | Consumer | Format |
|---|---|---|
| `praxis list` | CLI, Neovim | One challenge name per line |
| `praxis challenge <id>` | Neovim | Content lines to stdout |
| `praxis target <id>` | Neovim | Target character to stdout |
| `praxis verify <id>` | Neovim | "cursor" or "buffer" |
| `praxis result <id>` | Neovim | Result lines to stdout |
| `praxis attempt <id>` | Neovim | Silent (internal) |
| `praxis record <id> <moves> <time_ms>` | Neovim | Silent (internal) |
| `praxis stats [id]` | CLI | Per-challenge or summary |
