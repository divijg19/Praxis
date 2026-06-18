# Curriculum Progression

## Overview

Praxis challenges are organized into six sequential stages. Each stage introduces a new category of Vim technique. Challenges within a stage build on skills from prior stages but may interleave with adjacent stages for pedagogical depth.

> **Note:** Stages represent pedagogical dependencies (what should be learned before what). A challenge belongs to exactly one stage (skill domain) and exactly one layer (pedagogical intent: Tutorial, Training, or Trial). These are intentionally separate taxonomies.

```
Movement → Search → Structural Navigation → Editing → Text Objects → Registers
```

---

## 1. Movement (3 challenges)

Fundamental cursor movement. Character-level navigation, grid traversal, and multi-byte encoding awareness.

| Challenge | Concept |
|---|---|
| motion_rush | `hjkl` |
| grid_rush | `hjkl` |
| utf8_cursor_hunter | `f` |

---

## 2. Search (7 challenges)

Pattern-based navigation. Character search, word motions, line targeting, and search repetition.

| Challenge | Concept |
|---|---|
| find_hunter | `f` |
| word_hunter | `w` |
| symbol_hunter | `F` |
| line_hunter | `j` |
| slash_hunter | `/` |
| question_hunter | `?` |
| repeat_hunter | `;` |

---

## 3. Structural Navigation (10 challenges)

Delimiter matching, sentence and paragraph navigation, and text-object selection. The bridge between navigation and editing.

| Challenge | Concept |
|---|---|
| paren_hunter | `%` |
| sentence_hunter | `)` |
| inner_paren_hunter | `i(` |
| around_paren_hunter | `a(` |
| inner_bracket_hunter | `i[` |
| around_bracket_hunter | `a[` |
| inner_quote_hunter | `i"` |
| around_quote_hunter | `a"` |
| paragraph_hunter | `{` |
| match_hunter | `%` |

---

## 4. Editing (7 challenges)

Character and line editing. Delete, replace, toggle case, word editing, and line operations.

| Challenge | Concept |
|---|---|
| delete_character_hunter | `x` |
| replace_character_hunter | `r` |
| toggle_case_hunter | `~` |
| delete_word_hunter | `dw` |
| change_word_hunter | `ciw` |
| delete_line_hunter | `dd` |
| delete_to_end_hunter | `D` |

---

## 5. Text Objects (9 challenges)

Operator + text-object combinations. Delete inside/around, change inside, applied to words, parentheses, and quotes.

| Challenge | Concept |
|---|---|
| delete_inner_word_hunter | `diw` |
| delete_around_word_hunter | `daw` |
| delete_inner_paren_hunter | `di(` |
| delete_around_paren_hunter | `da(` |
| delete_inner_quote_hunter | `di"` |
| delete_around_quote_hunter | `da"` |
| change_inner_word_hunter | `ciw` |
| change_inner_paren_hunter | `ci(` |
| change_inner_quote_hunter | `ci"` |

---

## 6. Registers (5 challenges)

Yanking, named registers, appending, and register-based replacement and duplication.

| Challenge | Concept |
|---|---|
| yank_line_hunter | `yy` |
| named_register_hunter | `"a` |
| word_register_hunter | `"A` |
| register_replace_hunter | `"ap` |
| register_duplicate_hunter | `"ap` |

---

## Stage Totals

| Stage | Tutorial | Training | Trial | Total |
|---|---|---|---|---|
| Movement | 3 | — | — | 3 |
| Search | 7 | — | — | 7 |
| Structural Navigation | 10 | — | — | 10 |
| Editing | 7 | 3 | 2 | 12 |
| Text Objects | 9 | 4 | 3 | 16 |
| Registers | 5 | 3 | — | 8 |
| **Total** | **41** | **10** | **5** | **56** |

---

## Training Layer

Training challenges teach **cross-family skill compositions**. While Tutorial challenges focus on single concepts (e.g., `diw` alone, `dw` alone), Training challenges combine two or more families:

- **Search + Text Objects** — find a character, then apply a text-object operation
- **Editing + Repeat** — perform an edit, then repeat it with `.`
- **Registers + Editing** — yank, cut, or paste as part of an edit sequence

Training challenges use the `composite` validator, which behaves like `buffer` (byte-exact buffer comparison) but enforces a `MaxMoves` threshold to prevent bruteforce approaches.

All 10 Training challenges are listed in [`CHALLENGES.md`](CHALLENGES.md) under the Training section.

---

## Trial Layer

Trial challenges test **composition recognition** under move constraints. Unlike Training (which teaches *how* to combine), Trial challenges present a goal without naming the technique. The user must recognize which composition applies and execute it within a declared move budget.

Trial instructions describe **goals, never techniques**. Every Trial traces to one or more specific Training challenges it expects the user to have mastered.

All 5 Trial challenges are listed in [`CHALLENGES.md`](CHALLENGES.md) under the Trial section.

| Challenge | Tests recognition of |
|---|---|
| `trial_find_delete` | `find_diw_combo` |
| `trial_find_change` | `find_ca_quote_combo` |
| `trial_dot_repeat` | `dw_dot_combo` |
| `trial_delete_choice` | `find_diw_combo` + `find_daw_combo` |
| `trial_repeat_choice` | `dw_dot_combo` + `ciw_dot_combo` |

---

## Future Mastery System (Not Implemented)

Future releases may map mastery signals onto visible ranks:

    Bronze → Silver → Gold → Platinum → Diamond → Master → Grandmaster

No thresholds, algorithms, or schemas are defined. This is a design anchor, not a commitment.
