# Challenge Catalog

Total: **41 challenges** across **7 curriculum packs**.

Challenge behavior originates from internal/content/content.go.

Challenges are grouped by historical shipping packs. For the pedagogical stage taxonomy, see PROGRESSION.md.

Primary Concept annotations are curriculum documentation
maintained alongside the challenge catalog.

## Movement

### Motion Rush

- **ID:** `motion_rush`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `hjkl`
- **Context:** `basic navigation`
- **Stage:** `Movement`

#### Content

```text
Move your cursor to the star ★
```

---

### Grid Rush

- **ID:** `grid_rush`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `hjkl`
- **Context:** `grid navigation`
- **Stage:** `Movement`

#### Content

```text
. . . . .
. . . ★ .
. . . . .
```

---

## Search

### Find Hunter

- **ID:** `find_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `f`
- **Context:** `character search`
- **Stage:** `Search`

#### Content

```text
find motions are fast

aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa★
```

---

### Word Hunter

- **ID:** `word_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `w`
- **Context:** `word motion`
- **Stage:** `Search`

#### Content

```text
word motions are fast

alpha beta gamma delta epsilon ★
```

---

### Symbol Hunter

- **ID:** `symbol_hunter`  
- **Verify:** `cursor`  
- **Target:** `@`  
- **Primary Concept:** `F`
- **Context:** `backward symbol search`
- **Stage:** `Search`

#### Content

```text
find the target symbol

.......................@
```

---

### Line Hunter

- **ID:** `line_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `j`
- **Context:** `line navigation`
- **Stage:** `Search`

#### Content

```text
Move to the line with the star

line one
line two
line three
line four
line five
★ line six
```

---

## Structural Navigation

### Paren Hunter

- **ID:** `paren_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `%`
- **Context:** `matching delimiters navigation`
- **Stage:** `Structural Navigation`

#### Content

```text
Jump to the matching paren

(                         )★
```

---

### Sentence Hunter

- **ID:** `sentence_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `)`
- **Context:** `sentence navigation`
- **Stage:** `Structural Navigation`

#### Content

```text
Jump between sentences

First sentence ends here.
Second. Third.
★ Fourth. Fifth.
```

---

### Slash Hunter

- **ID:** `slash_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `/`
- **Context:** `forward search`
- **Stage:** `Search`

#### Content

```text
Search forward to find the target

alpha  bravo  charlie  delta  echo  foxtrot  golf  ★
```

---

### Question Hunter

- **ID:** `question_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `?`
- **Context:** `backward search`
- **Stage:** `Search`

#### Content

```text
Search backward to find the target

line one
line two
line three
line four
line five
line six
line seven
line eight
line nine
line ten
line eleven
line twelve
line thirteen
line fourteen
line fifteen
line sixteen
★ line seventeen
line eighteen
```

---

### Repeat Hunter

- **ID:** `repeat_hunter`  
- **Verify:** `cursor`  
- **Target:** `@`  
- **Primary Concept:** `;`
- **Context:** `repeat motion`
- **Stage:** `Search`

#### Content

```text
Search for ★, then repeat to find @

★  ★  ★  ★  ★  ★  ★  ★  @
```

---

### Inner Paren Hunter

- **ID:** `inner_paren_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `i(`
- **Context:** `select inside parentheses`
- **Stage:** `Structural Navigation`

#### Content

```text
Select inside the parentheses

(hello★)
```

---

### Around Paren Hunter

- **ID:** `around_paren_hunter`  
- **Verify:** `cursor`  
- **Target:** `)`  
- **Primary Concept:** `a(`
- **Context:** `select around parentheses`
- **Stage:** `Structural Navigation`

#### Content

```text
Select around the parentheses

(hello)★
```

---

### Inner Bracket Hunter

- **ID:** `inner_bracket_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `i[`
- **Context:** `select inside brackets`
- **Stage:** `Structural Navigation`

#### Content

```text
Select inside the brackets

[hello★]
```

---

### Around Bracket Hunter

- **ID:** `around_bracket_hunter`  
- **Verify:** `cursor`  
- **Target:** `]`  
- **Primary Concept:** `a[`
- **Context:** `select around brackets`
- **Stage:** `Structural Navigation`

#### Content

```text
Select around the brackets

[hello]★
```

---

### Inner Quote Hunter

- **ID:** `inner_quote_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `i"`
- **Context:** `select inside quotes`
- **Stage:** `Structural Navigation`

#### Content

```text
Select inside the quotes

x"magnificent★"
```

---

### Around Quote Hunter

- **ID:** `around_quote_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `a"`
- **Context:** `select around quotes`
- **Stage:** `Structural Navigation`

#### Content

```text
Select around the quotes

x"magnificent"★
```

---

### Paragraph Hunter

- **ID:** `paragraph_hunter`  
- **Verify:** `cursor`  
- **Target:** `Z`  
- **Primary Concept:** `{`
- **Context:** `paragraph navigation`
- **Stage:** `Structural Navigation`

#### Content

```text
Jump to the last paragraph

first paragraph

second paragraph

third Z
```

---

### Match Hunter

- **ID:** `match_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `%`
- **Context:** `nested delimiter matching`
- **Stage:** `Structural Navigation`

#### Content

```text
Jump to the matching delimiter

[         ]★
```

---

## Editing

### Delete Character Hunter

- **ID:** `delete_character_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `x`
- **Context:** `delete character`
- **Stage:** `Editing`

#### Content

```text
Delete the extra letter

helllo
```

#### Result

```text
Delete the extra letter

hello
```

---

### Replace Character Hunter

- **ID:** `replace_character_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `r`
- **Context:** `replace character`
- **Stage:** `Editing`

#### Content

```text
Replace the wrong letter

hallo
```

#### Result

```text
Replace the wrong letter

hello
```

---

### Toggle Case Hunter

- **ID:** `toggle_case_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `~`
- **Context:** `toggle case`
- **Stage:** `Editing`

#### Content

```text
Toggle the case of each letter

HELLO
```

#### Result

```text
Toggle the case of each letter

hello
```

---

### Delete Word Hunter

- **ID:** `delete_word_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `dw`
- **Context:** `delete word`
- **Stage:** `Editing`

#### Content

```text
Delete the middle word

keep lose keep
```

#### Result

```text
Delete the middle word

keep  keep
```

---

### Change Word Hunter

- **ID:** `change_word_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `ciw`
- **Context:** `simple word replacement`
- **Stage:** `Editing`

#### Content

```text
Change the word using ciw

foo
```

#### Result

```text
Change the word using ciw

bar
```

---

## UTF-8 Proof

### UTF-8 Cursor Hunter

- **ID:** `utf8_cursor_hunter`  
- **Verify:** `cursor`  
- **Target:** `★`  
- **Primary Concept:** `f`
- **Context:** `multi-byte search`
- **Stage:** `Movement`

#### Content

```text
Navigate past multi-byte characters to the star

α β γ ★
```

---

## Structural Editing

### Delete Line Hunter

- **ID:** `delete_line_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `dd`
- **Context:** `delete line`
- **Stage:** `Editing`

#### Content

```text
Delete the middle line

keep
remove
keep
```

#### Result

```text
Delete the middle line

keep
keep
```

---

### Delete To End Hunter

- **ID:** `delete_to_end_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `D`
- **Context:** `delete to end of line`
- **Stage:** `Editing`

#### Content

```text
Delete from cursor to end of line

keep this text remove_this_part
```

#### Result

```text
Delete from cursor to end of line

keep this text
```

---

### Delete Inner Word Hunter

- **ID:** `delete_inner_word_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `diw`
- **Context:** `delete inside word`
- **Stage:** `Text Objects`

#### Content

```text
Delete inside the middle word

keep remove keep
```

#### Result

```text
Delete inside the middle word

keep  keep
```

---

### Delete Around Word Hunter

- **ID:** `delete_around_word_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `daw`
- **Context:** `delete around word`
- **Stage:** `Text Objects`

#### Content

```text
Delete around the middle word

keep remove keep
```

#### Result

```text
Delete around the middle word

keep keep
```

---

### Delete Inner Paren Hunter

- **ID:** `delete_inner_paren_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `di(`
- **Context:** `delete inside parentheses`
- **Stage:** `Text Objects`

#### Content

```text
Delete inside parentheses

func(remove)
```

#### Result

```text
Delete inside parentheses

func()
```

---

### Delete Around Paren Hunter

- **ID:** `delete_around_paren_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `da(`
- **Context:** `delete around parentheses`
- **Stage:** `Text Objects`

#### Content

```text
Delete around parentheses

func(remove)
```

#### Result

```text
Delete around parentheses

func
```

---

### Delete Inner Quote Hunter

- **ID:** `delete_inner_quote_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `di"`
- **Context:** `delete inside quotes`
- **Stage:** `Text Objects`

#### Content

```text
Delete inside quotes

"remove"
```

#### Result

```text
Delete inside quotes

""
```

---

### Delete Around Quote Hunter

- **ID:** `delete_around_quote_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `da"`
- **Context:** `delete around quotes`
- **Stage:** `Text Objects`

#### Content

```text
Delete around quotes

"remove"
```

#### Result

```text
Delete around quotes


```

---

### Change Inner Word Hunter

- **ID:** `change_inner_word_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `ciw`
- **Context:** `word replacement within structural editing`
- **Stage:** `Text Objects`

#### Content

```text
Change the first word

goodbye world
```

#### Result

```text
Change the first word

hello world
```

---

### Change Inner Paren Hunter

- **ID:** `change_inner_paren_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `ci(`
- **Context:** `change inside parentheses`
- **Stage:** `Text Objects`

#### Content

```text
Change inside parentheses

func(remove)
```

#### Result

```text
Change inside parentheses

func(hello)
```

---

### Change Inner Quote Hunter

- **ID:** `change_inner_quote_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `ci"`
- **Context:** `change inside quotes`
- **Stage:** `Text Objects`

#### Content

```text
Change inside quotes

"remove"
```

#### Result

```text
Change inside quotes

"hello"
```

---

### Unnamed Register Hunter

- **ID:** `yank_line_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `yy`
- **Context:** `yank line`
- **Stage:** `Registers`

#### Content

```text
Every yank enters a register. Yank and paste to see.

copy me
```

#### Result

```text
Every yank enters a register. Yank and paste to see.

copy me
copy me
```

---

## Registers

### Named Register Hunter

- **ID:** `named_register_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `"a`
- **Context:** `named register`
- **Stage:** `Registers`

#### Content

```text
Store and retrieve with register a

copy me
```

#### Result

```text
Store and retrieve with register a

copy me
copy me
```

---

### Word Register Hunter

- **ID:** `word_register_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `"A`
- **Context:** `append to register`
- **Stage:** `Registers`

#### Content

```text
Store a word in register a, then append it

alpha beta
```

#### Result

```text
Store a word in register a, then append it

alpha beta alpha
```

---

### Register Replace Hunter

- **ID:** `register_replace_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `"ap`
- **Context:** `replace content from named register`
- **Stage:** `Registers`

#### Content

```text
Store 'correct' in register a, then replace 'wrong'

correct
wrong
```

#### Result

```text
Store 'correct' in register a, then replace 'wrong'

correct
correct
```

---

### Register Duplicate Hunter

- **ID:** `register_duplicate_hunter`  
- **Verify:** `buffer`  
- **Primary Concept:** `"ap`
- **Context:** `duplicate content from named register`
- **Stage:** `Registers`

#### Content

```text
Duplicate 'foo' below 'bar' using register a

foo
bar
```

#### Result

```text
Duplicate 'foo' below 'bar' using register a

foo
bar
foo
```

---

