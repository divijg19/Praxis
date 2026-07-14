# Curriculum

**Purpose:** What Praxis teaches — vocabulary, learning loop, navigation, stages, layers, progression, and what happens after the trial.

---

## 1. Vocabulary

### Primitive

A single Vim operation.

### Composition

Two or more primitives combined into a repeatable technique.

### Recognition

Selecting the correct composition for a scenario.

### Mastery

Consistently recognizing and executing the correct composition.

---

## 2. Learning Loop

Every challenge follows a four-phase learning loop:

```
Observe → Practice → Apply → Reflect
```

### Observe

The technique is introduced, explained, and contextualized. The learner sees the technique demonstrated and understands why it exists.

### Practice

The technique is reproduced with support. The learner executes with scaffolding (hints, guidance, guardrails) that prevents frustration while building muscle memory.

### Apply

The technique is executed independently. No scaffolding. The learner demonstrates understanding by completing the task without support.

### Reflect

Reflection is the period after execution where the learner compares intent, execution, and outcome. It follows each active phase. This is a contract, not an implementation — UI mechanisms (result screen, move count, mastery update) are deferred to UX design.

### Attempt states

Every challenge has three attempts. Each attempt is a pedagogical state, not a retry count:

| State | Learner Experience |
|---|---|
| Observe | The technique is introduced and contextualized. The learner sees what to do and why. |
| Practice | The technique is reproduced with support. Hints or guardrails are available. |
| Apply | The technique is executed independently. No hints. The learner demonstrates understanding. |

Attempts progress in order: Observe → Practice → Apply. A learner who succeeds at Apply in a given attempt is not required to repeat the state.

### Layer × O/P/A matrix

| Layer | Observe | Practice | Apply |
|---|---|---|---|
| Tutorial | Learn a single primitive | Guided primitive execution | Use the primitive independently |
| Training | Learn a composition | Guided composition execution | Use the composition within budget |
| Trial | Understand the scenario | Recognize the correct composition | Execute the composition within budget |

---

## 3. Navigation

Praxis has one entry point and a flat set of surfaces. There is no nested
menu hierarchy.

- **`:Praxis`** is the universal (re)entry point. On first launch it opens
  the Onboarding; on later launches it opens the Hub. `:Praxis <id>`
  opens a specific challenge directly.
- **Onboarding** (first launch only) offers four actions:
   - `[s]` Start → opens the next uncompleted challenge.
   - `[e]` Explore → the Catalog (a flat list of all challenges).
   - `[h]` About → what Praxis teaches and how progression works.
   - `[p]` View progress → the Hub.
- **Hub** (returning users) shows Current, Progress, Direction, and Mastery,
   then: `[Enter]` Continue to the next challenge, `[r]` Review the
   recommended challenge, `[q]` Back.
- **Challenge** buffer: solve it to see the result, then `[r]` Retry,
   `[Enter]` Continue, `[q]` Back.
- **Catalog** is a flat, unordered list. It does not group by stage and does
  not bypass progression — it is a reference view.

---

## 4. Stages

Challenges are grouped into six stages. Each stage teaches techniques that the next stage builds on.

- **Movement** — cursor control (h, j, k, l)
- **Search** — target acquisition (f, w, /, ?)
- **Structural Navigation** — semantic movement (%, (), {}, i(, a")
- **Editing** — mutation (x, r, ~, dw, ciw, dd)
- **Text Objects** — scoped mutation (diw, daw, di(, ciw)
- **Registers** — memory (yy, "a, "A, "ap)

---

## 5. Layers

Every challenge belongs to one layer. Layers describe the depth of understanding expected.

- **Tutorial** — learn a single technique in isolation, with a teaching hint in the instruction line.
- **Training** — combine techniques into compositions under a move budget.
- **Trial** — recognize the correct composition for an unfamiliar scenario, with no hint.

The intended learning loop is Observe → Practice → Apply, but the UI does not
separate these into distinct modes: each challenge is one buffer the learner
solves, then continues.

---

## 6. Progression

### Stage progression

```
Movement → Search → Structural Navigation → Editing → Text Objects → Registers
```

This order is fixed by the curriculum definition. Rationale:

1. **Movement first.** Everything requires cursor control.
2. **Search second.** After basic movement, learn to jump directly.
3. **Structural Navigation third.** After search, learn semantic navigation.
4. **Editing fourth.** After navigation, learn to modify text.
5. **Text Objects fifth.** After basic editing, learn scoped operations.
6. **Registers sixth.** After all navigation and editing, learn memory.

### Layer progression

```
Tutorial → Training → Trial
```

Enforced by curriculum order. Rationale:

1. **Tutorial** — learn each primitive in isolation.
2. **Training** — combine primitives into compositions.
3. **Trial** — recognize the correct composition under pressure.

### Full progression path

```
Onboarding
    ↓
Movement Tutorials
    ↓
Search Tutorials
    ↓
Structural Navigation Tutorials
    ↓
Editing Tutorials → Training → Trials
    ↓
Text Objects Tutorials → Training → Trials
    ↓
Registers Tutorials → Training
```

### When to advance

Progression is **count-based**, not stage-gated. A challenge stays current
until it has been **completed three times**, after which the next unpracticed
challenge becomes current. Stages exist as metadata (the Hub shows the next
challenge's stage as `Current`) but do not gate progression.

Mastery tiers are also count-based (see `internal/stats`):

- **Unseen** — 0 completions
- **Learning** — 1–2 completions
- **Practiced** — 3–7 completions
- **Experienced** — 8+ completions

The full curriculum is complete when every challenge has been completed three times. At that point the Hub shows `Current: Complete`.

### When to revisit

A learner revisits an earlier stage when:

1. **Stuck on a composition.** Repeatedly failing a Training or Trial that depends on earlier primitives.
2. **Low confidence.** A primitive has been used primarily in Apply mode without sufficient Observe or Practice grounding.
3. **Returning after absence.** The learner returns after time away from Praxis.

Revisiting is a pedagogical tool, not a penalty. Progress is preserved when returning to current material.

---

## 7. After the Curriculum

The curriculum ends when every challenge has been completed three times
(count-based; see §6). At that point the Hub shows `Current: Complete`
and the learner can press `[r]` to review the recommended challenge or
`[q]` to finish. Any challenge may be revisited at any time for practice
or improvement.
