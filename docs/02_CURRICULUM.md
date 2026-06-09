# Curriculum

Praxis organizes challenges into **curriculum packs**. Each pack builds on previous ones. The ordering is pedagogical: skills introduced in earlier packs are prerequisites for later ones.

## v0.0.x Curriculum Path

### Movement (2 challenges)
Navigate to a target on screen using basic cursor motions.

- `motion_rush` — Move your cursor to the star
- `grid_rush` — Navigate a grid to reach the star

**Verify:** cursor (navigate to target)

---

### Search (4 challenges)
Use word motions, find/till, and search to reach targets efficiently.

- `find_hunter` — Find motion to reach the star
- `word_hunter` — Word motions across words
- `symbol_hunter` — Find to a specific symbol
- `line_hunter` — Line-wise navigation

**Verify:** cursor

---

### Structural Navigation (13 challenges)
Navigate by structure: parentheses, sentences, search, text objects.

- `paren_hunter`, `sentence_hunter` — Structural jumps
- `slash_hunter`, `question_hunter`, `repeat_hunter` — Search-based navigation
- `inner_paren_hunter`, `around_paren_hunter` — Inside/around text objects
- `inner_bracket_hunter`, `around_bracket_hunter`
- `inner_quote_hunter`, `around_quote_hunter`
- `paragraph_hunter` — Paragraph navigation
- `match_hunter` — Matching delimiter jump

**Verify:** cursor

---

### Editing (5 challenges)
Modify text at the character and word level.

- `delete_character_hunter` — `x`
- `replace_character_hunter` — `r`
- `toggle_case_hunter` — `~`
- `delete_word_hunter` — `dw`
- `change_word_hunter` — `ciw`

**Verify:** buffer (edit buffer to match result)

---

### UTF-8 Proof (1 challenge)
Validate byte-to-character normalization with multi-byte content.

- `utf8_cursor_hunter` — Navigate past Greek characters to ★

**Verify:** cursor

---

### Structural Editing (12 challenges)
Apply operators to text objects: delete, change, yank.

- `delete_line_hunter` — `dd`
- `delete_to_end_hunter` — `D`
- `delete_inner_word_hunter` — `diw`
- `delete_around_word_hunter` — `daw`
- `delete_inner_paren_hunter` — `di(`
- `delete_around_paren_hunter` — `da(`
- `delete_inner_quote_hunter` — `di"`
- `delete_around_quote_hunter` — `da"`
- `change_inner_word_hunter` — `ciw → hello`
- `change_inner_paren_hunter` — `ci( → hello`
- `change_inner_quote_hunter` — `ci" → hello`
- `yank_line_hunter` — `yyp` (reframed as unnamed register lesson)

**Verify:** buffer

---

### Registers (4 challenges)
Store, retrieve, and reuse text across registers.

- `named_register_hunter` — `"ayy` + `"ap`
- `word_register_hunter` — `"ayiw` + `A`
- `register_replace_hunter` — `"ayiw` + `diw` + `"ap`
- `register_duplicate_hunter` — `"ayy` + `"ap` across lines

**Verify:** buffer

---

## Total

**41 challenges** across **7 curriculum packs**, 20 cursor-verified + 21 buffer-verified.

## v0.1.x Direction

The curriculum focus shifts from breadth to depth:
- Precision mode: stricter validation, timing, and accuracy tracking
- Challenge quality: better content, feedback, and progression
- Mastery systems: bronze through grandmaster tiers
