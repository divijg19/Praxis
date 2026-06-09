# Validators

Validators determine how a challenge's success condition is evaluated.

Praxis has two validators: **cursor** and **buffer**.

## Validator Registry

The registry is in `internal/validator/validator.go`:

```go
var valid = map[string]bool{
    "cursor": true,
    "buffer": true,
}

func Exists(name string) bool {
    return valid[name]
}
```

`Exists()` is used by `TestValidatorCoverage` to ensure every challenge references a known validator.

## Cursor Validator

**Verify:** `"cursor"`

**Rule:** Navigate the cursor to the target character on screen.

**How it works:**
- Buffer is set to `Content` with `modifiable=false`
- `CursorMoved` autocommand checks if the character under the cursor matches `Target`
- Uses `byte_to_char()` normalization to handle multi-byte content correctly
- On match: sets `state.done = true`, echoes success

**Used by:** 20 challenges (movement, search, navigation, UTF-8 proof)

**Required fields:**
- `Target` — the character the user must navigate to
- `Content` — the buffer content (may include instructions in the same view)
- `Result` — must be nil or empty (cursor challenges don't transform the buffer)

## Buffer Validator

**Verify:** `"buffer"`

**Rule:** Edit the buffer content to exactly match `Result`.

**How it works:**
- Buffer is set to `Content` with `modifiable=true`
- `TextChanged` autocommand fires on Normal mode edits (`x`, `r`, `~`, `dw`, `ciw`, etc.) and increments moves count
- `TextChangedI` autocommand fires on Insert mode keystrokes (typing during `ciw`, `A`, etc.) but does NOT increment moves
- `check_buffer()` compares every line of the current buffer against `Result`
- Comparison is line-by-line exact string equality (no trimming, no fuzzy matching)
- On match: sets `state.done = true`, echoes success

**Used by:** 21 challenges (editing, structural editing, registers)

**Required fields:**
- `Result` — the exact target buffer state after editing
- `Content` — required to have instruction, blank, play area (3+ lines)
- `Target` — must be empty string

## CursorMoved-only dispatching

The Lua frontend only attaches `CursorMoved` for cursor challenges. Buffer challenges use `TextChanged` + `TextChangedI` instead. This is handled by checking `state.verify` when setting up autocommands.

## Future validator candidates

Not yet implemented. Listed for architectural awareness:

- `selection` — validate that a specific text range is visually selected
- `state` — validate a non-buffer state condition
- `register` — validate register contents after an operation
- `operator` — validate completion of an operator-pending mode sequence
