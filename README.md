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

### DSA Mastery

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

```text
Campaign
└── Region
    └── Lesson
        └── Challenge
```

Every challenge is scored.

Every lesson contributes to progression.

Every discipline has a mastery path.

```text
Bronze → Silver → Gold → Platinum → Master → Grandmaster
```

---

## Boss Battles

Knowledge is easy.

Execution is difficult.

Boss battles combine multiple skills into practical scenarios that test real competence rather than memorization.

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
                 `Praxis`

             ┌─────────────┐
             │  Go Engine  │
             └──────┬──────┘
                    │
        ┌───────────┴───────────┐
        │                       │
        ▼                       ▼

   Standalone TUI        Neovim Frontend
       (Go)                  (Lua)
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

Current focus:

* Vim Mastery
* DSA Mastery

Everything else can wait.

---

**Less consumption. More competence.**

