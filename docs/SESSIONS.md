# Practice Sessions

## Overview

A session represents a contiguous period of deliberate practice within a single Neovim instance. Sessions are ephemeral — they exist only in memory and are discarded when Neovim exits.

There is exactly one session per Neovim instance. It begins when the first `:Praxis <id>` is called and accumulates data until Neovim is closed.

## Command

```
:PraxisSession
```

Opens a read-only scratch buffer showing the current session's aggregated metrics.

### Before any practice

```
Session

Challenges: 0
Completions: 0

Session Length: 0s
Practice Time: 0s

Moves: 0

Avg Moves: 0
Avg Time: 0s
```

### During practice

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

## Fields

| Field | Meaning |
|---|---|
| **Challenges** | Total `:Praxis <id>` invocations in this session |
| **Completions** | Total successful completions (including replays) |
| **Session Length** | Wall-clock time since the first `:Praxis <id>` was invoked |
| **Practice Time** | Sum of individual challenge completion times |
| **Moves** | Total moves across all completions |
| **Avg Moves** | `Moves / Completions` |
| **Avg Time** | `Practice Time / Completions` |

## Design Decisions

- **No persistence.** Sessions are Neovim-scoped ephemeral state. If you restart Neovim, the session starts fresh.
- **No automatic popup.** `:PraxisSession` is user-invoked only. No `VimLeavePre` or completion hooks show this view automatically.
- **Lua-only, no Go.** The session lives as a module-level table in `nvim/lua/praxis/init.lua`. No new package, CLI command, persistence layer, or test surface.
- **Challenges and Completions are separate counters.** `Challenges` increments on every `:Praxis <id>` call. `Completions` increments only on success. A user who opens 8 challenges and completes 5 sees both numbers — the gap shows unfinished attempts.
- **Session Length tracks wall-clock time.** Not idle-adjusted. It is simply `now - start_ns` from when the first challenge was opened. This is honest but approximate (includes any idle time in Neovim).

## Future

Session data is intentionally minimal. Potential expansions (each a separate release):

- Per-challenge breakdown within a session
- Idle time detection for accurate Session Length
- Cross-session aggregation (persisting session summaries)
