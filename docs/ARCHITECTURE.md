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
|---|---|
| `cmd/praxis` | CLI entry point: `list`, `challenge`, `target`, `verify`, `result`, `attempt`, `record`, `stats` |
| `internal/challenge` | `Challenge` struct — the core data model |
| `internal/content` | `All()` — the complete challenge registry |
| `internal/stats` | `Stats` struct, `Load`, `Save`, `Update` — persistent progress tracking |
| `internal/validator` | `Exists(name)` — validator dispatch registry |

### Data flow

```
User request (CLI)
    |
    v
cmd/praxis/main.go
    |
    v
internal/content.All()  →  challenge.Challenge{ID, Name, Verify, Target, Content, Result, Layer}
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

The Neovim frontend lives in `nvim/lua/praxis/init.lua`. It is loaded via `nvim/plugin/praxis.lua`.

### Responsibilities

- Fetch challenge data from the Go binary (`/tmp/praxis`)
- Create buffer with challenge content
- Set up autocmds for validation:
  - `CursorMoved` — cursor challenges (target reached check)
  - `TextChanged` — buffer challenges (buffer matches result, Normal mode edits increment moves)
  - `TextChangedI` — buffer challenges (Insert mode keystrokes, no move increment)

### Validator dispatch

The Lua frontend checks `state.verify` to decide which autocmd behavior to enable:
- `"cursor"` → modifiable=false, CursorMoved listener, target check
- `"buffer"` → modifiable=true, TextChanged + TextChangedI listeners, buffer comparison via `check_buffer()`

### byte_to_char normalization

```lua
function byte_to_char(line, bytecol)
    return vim.fn.strchars(string.sub(line, 1, bytecol))
end
```

Converts Neovim's 0-indexed byte column to a 0-indexed character column. Critical for multi-byte content (UTF-8 Greek, emoji, etc.).

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
