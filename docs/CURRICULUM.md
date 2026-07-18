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

## 2. Learning Model

Praxis is deliberate practice, not a course. The three layers map to a
single arc:

```
Onboard  →  Practice  →  Improve  →  Return forever
```

**Tutorial is the end of onboarding, not the beginning of a curriculum.** Its
only job is to teach the few essentials so a learner can genuinely use Praxis
on their own. Everything after Tutorial *is* the product.

- **Tutorial** teaches the few essential primitives, with a hint.
- **Training** is repetition and composition — the learner practices and
  chooses what to improve.
- **Trial** is open problem-solving with no prescribed keystrokes.

Each challenge is one buffer the learner solves, then continues. The UI does
not separate these into distinct modes; the layers describe the *kind* of
experience, not a retry count or enforced sequence.

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

Challenges are grouped into six stages (Movement, Search, Structural
Navigation, Editing, Text Objects, Registers). The authoritative primitive
breakdown per stage is in [ARCHITECTURE.md](./ARCHITECTURE.md) (Stage
taxonomy).

---

## 5. Layers

Every challenge belongs to one layer: Tutorial, Training, or Trial. The
intended learning loop is Observe → Practice → Apply, but the UI does not
separate these into distinct modes — each challenge is one buffer the
learner solves, then continues. See [ARCHITECTURE.md](./ARCHITECTURE.md)
(Layer taxonomy) for the per-layer purpose and challenge counts.

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

- **Unseen** — 0 completions
- **Learning** — 1–2 completions
- **Practiced** — 3–7 completions
- **Experienced** — 8+ completions

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
