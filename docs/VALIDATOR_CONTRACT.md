# Validator Contract

This document formally specifies the behavior of each validator. Validators determine how a challenge's success condition is evaluated.

## Cursor Validator

**Verify value:** `"cursor"`

### Rule

The challenge succeeds when the character under the cursor matches the target character.

### Formal specification

```
Given:
  buf  = challenge content (array of N lines)
  col0 = cursor byte column (0-indexed)
  row  = cursor row (1-indexed)
  line = buf[row-1]  (0-indexed line access)

The cursor position is valid when:
  line is not nil
  AND strcharpart(line, rune_offset(line, col0), 1) == target

Where:
  rune_offset(line, col0) = number of runes in line[0:col0]
  strcharpart(s, start, len) = substring of s at rune start, len runes long
```

### PASS examples

| Challenge | Content (play area) | Target | Correct cursor position |
|---|---|---|---|
| motion_rush | `Move your cursor to the star ★` | ★ | On the ★ character |
| find_hunter | `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa★` | ★ | On the ★ (last character) |
| symbol_hunter | `.......................@` | @ | On the @ (last character) |
| repeat_hunter | `★  ★  ★  ★  ★  ★  ★  ★  @` | @ | On the @ (last character) |
| paren_hunter | `(                         )★` | ★ | On the ★ after the space |
| utf8_cursor_hunter | `α β γ ★` | ★ | On the ★ (bytecol=9, charcol=6) |

### FAIL examples

| Scenario | Why it fails |
|---|---|
| Cursor on wrong character | Character under cursor does not equal target |
| Target not in buffer | No matching character exists |
| Empty target | Target field is `""` (violates contract) |
| Cursor on instruction line | Character is the target but wrong location |

## Buffer Validator

**Verify value:** `"buffer"`

### Rule

The challenge succeeds when the current buffer content exactly matches the result content.

### Formal specification

```
Given:
  current = current buffer lines (array of M strings)
  result  = expected result lines (array of N strings)

The buffer is valid when:
  M == N
  AND current[i] == result[i] for all i in [0, N)

Comparison is byte-exact string equality:
  - No trimming
  - No case folding
  - No whitespace normalization
```

### PASS examples

| Challenge | Result (play area) | Editing technique |
|---|---|---|
| delete_character_hunter | `hello` | `x` on extra `l` |
| delete_word_hunter | `keep  keep` | `dw` on middle word |
| delete_inner_paren_hunter | `func()` | `di(` on inner content |
| toggle_case_hunter | `hello` | `~` on each letter |
| register_duplicate_hunter | `foo` + `bar` + `foo` | `"ayy` + `"ap` to duplicate |

### FAIL examples

| Scenario | Why it fails |
|---|---|
| Buffer has wrong content | At least one line differs from result |
| Buffer has extra line | Line count mismatch (M != N) |
| Buffer has missing line | Line count mismatch (M != N) |
| Trailing whitespace differs | `"keep "` vs `"keep"` — byte-exact comparison |
| Result is empty | Result field is `[]string{}` or nil (violates contract) |

## Contract Enforcement

These contracts are enforced by:

- `TestValidatorCoverage` — every challenge references a registered validator
- `TestNoValidatorDrift` — every registered validator is used by at least one challenge
- `TestResultShapeMatchesVerify` — cursor challenges forbid Result; buffer challenges forbid Target
- `TestExistsCursor` / `TestExistsBuffer` — validators are present in the registry
