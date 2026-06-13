# `Praxis`
> **Mastery through practice.**

Most developer education optimizes for consumption.
`Praxis` optimizes for competence.

No courses.
No playlists.
No passive lessons.

Just challenges, feedback, progression, and deliberate practice.

---

## What is `Praxis`?

`Praxis` is a terminal-first mastery game for developers.

It transforms technical disciplines into campaigns, lessons, challenges, trials, and boss battles designed to build real skill.

Learn by doing.

Measure improvement.

Earn mastery.

---

## Initial Campaigns

### Vim Mastery

```text
Movement
Word Motions
Find & Till
Text Objects
Registers
Macros
Marks
Buffers
Windows
Real Editing Missions
```

### Future: DSA Mastery

```text
Arrays
Strings
Hash Maps
Linked Lists
Trees
Graphs
Heaps
Greedy
Dynamic Programming
```

Future disciplines may include Git, Linux, Shell, Systems Programming, Networking, and Infrastructure.

---

## Gameplay

Every challenge is scored on moves and time. Progress is tracked persistently. Mastery tiers (Unseen → Learning → Practiced → Experienced) are earned through repetition.

*Future plans include campaign, region, and lesson hierarchies for organizing disciplines.*

---

## Future: Boss Battles

Boss battles combine multiple skills into practical scenarios that test real competence rather than memorization.

*Not yet implemented. Planned for future releases.*

---

## Run Anywhere

Launch from the terminal:

```bash
praxis
```

Or directly from Neovim:

```vim
:Praxis
```

---

## Architecture

```text
                 Praxis

             ┌─────────────┐
             │  Go Engine  │
             └──────┬──────┘
                    │
              ┌─────┬─────┐
              │           │
              ▼           ▼
        ┌──────────┐ ┌──────────┐
        │   CLI    │ │  Neovim  │
        │ (Go cmd) │ │  (Lua)   │
        └──────────┘ └──────────┘
```

One engine.

Multiple frontends.

Shared progression, content, scoring, and persistence.

---

## Philosophy

Every feature in `Praxis` must answer one question:

> Does this help the user develop real skill through deliberate practice?

If not, it doesn't belong.

---

## Status

Early development.

**Implemented:** Vim challenge engine, persistent stats, mastery tiers, practice guidance, Neovim frontend, CLI.

**Future:** DSA Mastery campaign, Boss Battles, rank progression, standalone TUI, additional disciplines.

---

**Less consumption. More competence.**

