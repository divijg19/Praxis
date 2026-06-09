# Content Format

This document describes how to add a new challenge to Praxis.

## The Challenge Struct

Every challenge is defined in `internal/content/content.go` as part of the `All()` function.

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

## Field Rules

### ID

- **Must be globally unique** across all challenges
- **Must never change** once released (stable public contract)
- Use `snake_case` with a descriptive prefix and `_hunter` suffix
- Examples: `motion_rush`, `delete_character_hunter`, `named_register_hunter`
- IDs from the same curriculum pack should be adjacent in the slice order

### Name

- Display name shown to the user
- May evolve as curriculum framing shifts
- Example: `yank_line_hunter` was reframed to `Unnamed Register Hunter` in v0.0.19

### Verify

- `"cursor"` for navigation-to-target challenges
- `"buffer"` for edit-the-buffer-to-match challenges
- Must match a registered validator name (see `internal/validator/validator.go`)

### Target

- Required for cursor challenges (non-empty string)
- Must be present somewhere in the buffer content
- For buffer challenges: must be empty string `""`

### Content

- First line must always be the **instruction line**
- Second line is typically blank (`""`) for challenges with multi-line content
- The remaining lines are the **play area** the user interacts with
- Cursor challenges may have 1+ lines; buffer challenges must have 3+ lines (instruction, blank, play area)

### Result

- Required for buffer challenges (non-empty slice)
- Must be the exact target state of the buffer after the user edits it
- Instruction line must match `Content[0]`
- NOT included for cursor challenges (leave as nil)

## Curriculum Pack Placement

Challenges should be added at the end of the `All()` slice, after the last challenge of the most recent release. The slice order is the pedagogical order: skills build on earlier skills.

When adding a new pack, ensure:
1. `TestChallengeIDsStable` is updated with the new ID in the correct position
2. `TestChallengeNamesStable` is updated with the new Name
3. `TestChallengeCount` automatically passes (derives from the stability list length)
4. `tools/replay/replay.lua` is updated with the new ID in the correct category

## API Contract

The full API contract is documented in `docs/API_CONTRACT.md`.

```go
Challenge.ID       // STABLE — never change after first release
Challenge.Verify   // STABLE — "cursor" or "buffer"
Challenge.Target   // STABLE — required for cursor, empty for buffer
Challenge.Result   // STABLE — required for buffer, nil for cursor
Challenge.Content  // STABLE — first line is instruction
Challenge.Name     // EVOLVABLE — curriculum framing may shift
```

The `yank_line_hunter` precedent: IDs are permanent identifiers. Names describe the current pedagogical framing. Both are protected by stability tests.
