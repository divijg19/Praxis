# Challenge Rationale

## motion_rush

**Concept:** `hjkl`
**Why:** First challenge. Introduces basic cursor movement (hjkl) with a single visible target.
**Not redundant:** Foundation challenge. All subsequent movement challenges build on this skill.
**Followed by:** `grid_rush`

## grid_rush

**Concept:** `hjkl`
**Why:** Reinforces hjkl with multi-line grid navigation. Requires planning a path rather than moving to a visible target.
**Not redundant:** Grid navigation introduces directional planning absent from `motion_rush`.
**Followed by:** `find_hunter`

## find_hunter

**Concept:** `f`
**Why:** First search-based motion. Introduces `f` for character-level forward search.
**Not redundant:** Teaches character search as an alternative to grid-based movement.
**Followed by:** `word_hunter`

## word_hunter

**Concept:** `w`
**Why:** Introduces word-level forward motion (`w`). Teaches navigation by word boundaries, faster than character-by-character movement.
**Not redundant:** Word motions enable faster navigation than character-by-character movement.
**Followed by:** `line_hunter`

## line_hunter

**Concept:** `j`
**Why:** Teaches vertical navigation to a specific line. Introduces the concept of targeting by line position.
**Not redundant:** Line-level navigation is distinct from character or word navigation.
**Followed by:** `paren_hunter`

## paren_hunter

**Concept:** `%`
**Why:** Introduces matching delimiter navigation. Teaches `%` for jumping between paired delimiters.
**Not redundant:** Foundation structural navigation challenge. Teaches the `%` motion itself rather than nested matching.
**Followed by:** `sentence_hunter`

## sentence_hunter

**Concept:** `)`
**Why:** Introduces sentence-level navigation. Teaches `)` and `(` for moving between sentences.
**Not redundant:** Sentence navigation is a distinct structural motion unrelated to delimiter matching.
**Followed by:** `slash_hunter`

## slash_hunter

**Concept:** `/`
**Why:** Introduces forward search with `/`. Teaches text-pattern search as a navigation strategy.
**Not redundant:** Text search is more flexible than character search (`f`).
**Followed by:** `question_hunter`

## question_hunter

**Concept:** `?`
**Why:** Introduces backward search with `?`. Complements forward search from `slash_hunter`.
**Not redundant:** Backward search (`?`) completes the bidirectional search model.
**Followed by:** `repeat_hunter`

## repeat_hunter

**Concept:** `;`
**Why:** Introduces search repetition with `;` and `,`. Teaches efficient reuse of the last search.
**Not redundant:** Repetition is a distinct skill from performing the initial search.
**Followed by:** `inner_paren_hunter`

## inner_paren_hunter

**Concept:** `i(`
**Why:** First text-object selection. Introduces `i(` for selecting inside parentheses.
**Not redundant:** Text objects are the foundation for operator + text-object editing patterns.
**Followed by:** `around_paren_hunter`

## around_paren_hunter

**Concept:** `a(`
**Why:** Introduces `a(` for selecting around parentheses (including the delimiters). Complements `i(`.
**Not redundant:** `a(` selects a different range than `i(` (includes delimiters).
**Followed by:** `inner_bracket_hunter`

## inner_bracket_hunter

**Concept:** `i[`
**Why:** Extends text-object selection to brackets. Reinforces the `i` prefix pattern.
**Not redundant:** Different delimiter type (brackets vs parentheses), different visual structure.
**Followed by:** `around_bracket_hunter`

## around_bracket_hunter

**Concept:** `a[`
**Why:** Introduces `a[` for selecting around brackets. Complements `i[`.
**Not redundant:** `a[` selects a different range than `i[` (includes brackets).
**Followed by:** `inner_quote_hunter`

## inner_quote_hunter

**Concept:** `i"`
**Why:** Extends text-object selection to quotes. Reinforces the `i` prefix pattern with string content.
**Not redundant:** Quotes are a new delimiter context requiring different cursor positioning.
**Followed by:** `around_quote_hunter`

## around_quote_hunter

**Concept:** `a"`
**Why:** Introduces `a"` for selecting around quotes. Complements `i"`.
**Not redundant:** `a"` selects quotes + content, distinct from `i"` which selects content only.
**Followed by:** `paragraph_hunter`

## paragraph_hunter

**Concept:** `{`
**Why:** Introduces paragraph-level navigation with `{` and `}`.
**Not redundant:** Paragraph navigation is a block-level motion distinct from line, sentence, or delimiter motions.
**Followed by:** `delete_character_hunter`

## delete_character_hunter

**Concept:** `x`
**Why:** First editing challenge. Introduces `x` for single-character deletion.
**Not redundant:** Foundation editing skill. All subsequent editing challenges build on this.
**Followed by:** `replace_character_hunter`

## replace_character_hunter

**Concept:** `r`
**Why:** Introduces `r` for single-character replacement without entering insert mode.
**Not redundant:** Replace (`r`) is a different operation than delete-then-insert.
**Followed by:** `toggle_case_hunter`

## toggle_case_hunter

**Concept:** `~`
**Why:** Introduces case toggling with `~`. Teaches a characterwise transformation.
**Not redundant:** Case toggling is a specialized editing operation distinct from delete or replace.
**Followed by:** `delete_word_hunter`

## delete_word_hunter

**Concept:** `dw`
**Why:** First operator + motion editing. Introduces `dw` for deleting a word.
**Not redundant:** `dw` combines an operator (`d`) with a motion (`w`), a foundational Vim editing pattern.
**Followed by:** `change_word_hunter`

## change_word_hunter

**Concept:** `ciw`
**Why:** Introduces the change operator with a text object. Teaches `ciw` for replacing a word.
**Not redundant:** `ciw` demonstrates operator + text-object composition, distinct from operator + motion (`dw`).
**Followed by:** `delete_inner_word_hunter`

## utf8_cursor_hunter

**Concept:** `utf8`
**Why:** Validates that cursor positioning works correctly with multi-byte (UTF-8) characters. Uses `h` to move left past Greek letters to a star target.
**Not redundant:** UTF-8 correctness is a distinct guarantee from single-byte cursor movement. This challenge exists to verify encoding handling.
**Followed by:** `delete_line_hunter`

## delete_line_hunter

**Concept:** `dd`
**Why:** Introduces line-level deletion with `dd`.
**Not redundant:** `dd` is a distinct operator from word or character deletion, and the most common delete operation.
**Followed by:** `delete_to_end_hunter`

## delete_to_end_hunter

**Concept:** `d$`
**Why:** Introduces `d$` for deleting from cursor to end of line.
**Not redundant:** `d$` is a distinct deletion range from word or character deletion.
**Followed by:** `delete_inner_word_hunter`

## delete_inner_word_hunter

**Concept:** `diw`
**Why:** First text-object deletion. Introduces `diw` for deleting inside a word.
**Not redundant:** `diw` deletes only the word content, unlike `daw` which also deletes surrounding whitespace.
**Followed by:** `delete_around_word_hunter`

## delete_around_word_hunter

**Concept:** `daw`
**Why:** Introduces `daw` for deleting a word and its surrounding whitespace. Complements `diw`.
**Not redundant:** `daw` deletes more aggressively than `diw`, producing different results.
**Followed by:** `delete_inner_paren_hunter`

## delete_inner_paren_hunter

**Concept:** `di(`
**Why:** Extends text-object deletion to parentheses. Introduces `di(`.
**Not redundant:** Text-object deletion with parentheses is structurally distinct from word-based deletion.
**Followed by:** `delete_around_paren_hunter`

## delete_around_paren_hunter

**Concept:** `da(`
**Why:** Introduces `da(` for deleting around parentheses (including delimiters).
**Not redundant:** `da(` deletes the parentheses themselves, unlike `di(` which preserves them.
**Followed by:** `delete_inner_quote_hunter`

## delete_inner_quote_hunter

**Concept:** `di"`
**Why:** Introduces `di"` for deleting inside quotes.
**Not redundant:** Quote-based deletion is structurally distinct from parentheses-based deletion.
**Followed by:** `delete_around_quote_hunter`

## delete_around_quote_hunter

**Concept:** `da"`
**Why:** Introduces `da"` for deleting around quotes (including the quotes themselves).
**Not redundant:** `da"` deletes the quotes, unlike `di"` which preserves them.
**Followed by:** `change_inner_paren_hunter`

## change_inner_paren_hunter

**Concept:** `ci(`
**Why:** Introduces `ci(` for changing content inside parentheses.
**Not redundant:** Change inside parentheses (`ci(`) is a distinct operation from delete inside parentheses (`di(`).
**Followed by:** `change_inner_quote_hunter`

## change_inner_quote_hunter

**Concept:** `ci"`
**Why:** Introduces `ci"` for changing content inside quotes.
**Not redundant:** Change inside quotes is a distinct operation from delete inside quotes.
**Followed by:** `yank_line_hunter`

## yank_line_hunter

**Concept:** `yy`
**Why:** Introduces yanking (copying) with `yy`. First register-adjacent challenge.
**Not redundant:** Yanking is a distinct operation from deleting or changing.
**Followed by:** `word_register_hunter`

## word_register_hunter

**Concept:** `"a`
**Why:** Introduces storing a word in a named register with `"ayw`. Teaches scoped register storage beyond whole-line yanks.
**Not redundant:** Storing a word (`"ayw`) is distinct from storing a whole line (`"ayy`).
**Followed by:** `register_replace_hunter`

## register_replace_hunter

**Concept:** `"ap`
**Why:** Teaches register-based replacement. Stores content, then uses it to overwrite other text.
**Not redundant:** Register replacement demonstrates a practical use case (replace-with-register) distinct from basic store-and-retrieve.
**Followed by:** `register_duplicate_hunter`

## register_duplicate_hunter

**Concept:** `"ap`
**Why:** Teaches register-based duplication. Stores content, then appends a copy of it elsewhere.
**Not redundant:** Duplication demonstrates a different register use case than replacement.
**Followed by:** (none — final challenge)
