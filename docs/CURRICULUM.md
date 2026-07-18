# Curriculum

**Purpose:** What Praxis teaches, including vocabulary, the learning loop, navigation, stages, layers, progression, and what happens after the trial.

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

## 2. Learning Model

Praxis is deliberate practice, not a course. Each mode makes a different
educational contract with the learner:

```
Tutorial ──┬── Core      teaches, mandatory, finite: "I will make you independent."
           └── Additional Lessons  teaches, never blocks: "Learn this when you're curious."
Training                improves: "You know it; I will help you master it."
Trials                  integrates: "Solve the problem your own way."
```

**Tutorial is the end of onboarding, not the beginning of a curriculum.** Its
only job is to teach the few essentials so a learner can genuinely use Praxis
on their own. When the Core is complete, Tutorial is *finished*, not paused,
not 50% done. Everything after Tutorial *is* the product.

- **Tutorial (Core)** is the mandatory, carefully ordered, confidence-building
  onboarding. Finish it and you are independent.
- **Tutorial (Additional Lessons)** is the rest of the lesson library. It still
  teaches, but never blocks Training, Trials, or exploration. It is the same
  Tutorial experience. The learner is simply browsing more Tutorial material,
  not entering a different product.
- **Training** never teaches mechanics. It refines fluency.
- **Trial** is open problem-solving with no prescribed keystrokes.

Each challenge is one buffer the learner solves, then continues. The UI does
not separate these into distinct modes; the layers describe the *kind* of
experience, not a retry count or enforced sequence.

### The one-need principle

Every educational mode exists to answer exactly one learner need:

| Mode | The learner's need |
|---|---|
| Tutorial | "I don't know." |
| Training | "I know, but I'm slow." |
| Trials | "I think I know." |

If a proposed feature can't be placed under one of these three sentences, it
does not belong in the educational model. This is the litmus test for future
design debates.

### What "finished" means

Praxis distinguishes three milestones so the learner is never confused about
where they stand:

- **Tutorial Complete**: Core onboarding is finished. You are now independent.
  This is the only milestone the system announces as "complete."
- **Practice Ongoing**: Training and Trials never truly finish. Fluency is a
  process, not a checkbox.
- **Curriculum Complete**: every exercise has been experienced at least once.
  This is an achievement, not the educational goal.

### Why Additional Lessons exist

Optional Tutorial material exists to **remove unnecessary prerequisites
without removing instruction**. Someone can ignore it, return to it later,
study it before Training, or consult it after failing a Trial. It adapts to the
learner instead of forcing them. This is exactly why it is part of Tutorial
and never a separate mode.

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
  not bypass progression. It is a reference view.

---

## 4. Stages

Challenges are grouped into six stages (Movement, Search, Structural
Navigation, Editing, Text Objects, Registers). The authoritative primitive
breakdown per stage is in [ARCHITECTURE.md](./ARCHITECTURE.md) (Stage
taxonomy).

---

## 5. Layers

Every challenge belongs to one layer: Tutorial, Training, or Trial. Each
layer makes a distinct educational contract (see §2), and the UI does not
separate them into distinct modes. Each challenge is one buffer the learner
solves, then continues. See [ARCHITECTURE.md](./ARCHITECTURE.md) (Layer
taxonomy) for the per-layer purpose and challenge counts.

### Tutorial: Core and Additional Lessons

Tutorial is split by *purpose*, not difficulty. Internally this is the
`core` / `optional` tier on each Tutorial challenge; externally it is presented
as Core and Additional Lessons so the learner never feels they have entered a
separate product.

- **Core**: the mandatory onboarding (currently ten exercises: movement,
  search, line movement, character delete, word delete, character replace,
  word change, line yank, plus two foundational extras). Completing Core
  finishes Tutorial.
- **Additional Lessons**: every other Tutorial challenge (internally
  `optional`). It teaches, but never blocks progression. The learner explores
  it from the Catalog when curious. See §2 for why it exists.

A challenge's tier is part of its metadata (`Tier: core | optional`) and is
assigned in `internal/content/curriculum.go`.

---

## 6. Progression

### Stage & layer progression

Progression order is fixed by the curriculum definition:

```
Movement → Search → Structural Navigation → Editing → Text Objects → Registers
Tutorial → Training → Trial
```

Stages exist as metadata (the Hub shows the next challenge's stage as
`Current`) but do not gate progression.

### When to advance

Progression is **count-based**, not stage-gated. A challenge stays current
until it has been **completed three times**, after which the next unpracticed
challenge becomes current. Mastery tiers are also count-based (see
`internal/stats`):

- **Unseen**: 0 completions
- **Learning**: 1 to 2 completions
- **Practiced**: 3 to 7 completions
- **Experienced**: 8 or more completions

The full curriculum is complete when every challenge has been completed three times. At that point the Hub shows `Current: Complete`.

### When to revisit

A learner revisits earlier material when:

1. **Stuck on a composition.** Repeatedly failing a Training or Trial that depends on earlier primitives.
2. **Low confidence.** A primitive has not been practiced enough to feel automatic.
3. **Returning after absence.** The learner returns after time away from Praxis.

Revisiting is a pedagogical tool, not a penalty. Progress is preserved when returning to current material.

---

## 7. After the Curriculum

The curriculum ends when every challenge has been completed three times
(count-based; see §6). At that point the Hub shows `Current: Complete`
and the learner can press `[r]` to review the recommended challenge or
`[q]` to finish. Any challenge may be revisited at any time for practice
or improvement.
