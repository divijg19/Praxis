# Vision

Most developer education optimizes for consumption.

Praxis optimizes for competence.

No courses. No playlists. No passive lessons. Just challenges, feedback, progression, and deliberate practice.

## Why Praxis Exists

Vim mastery is not about knowing key bindings. It is about developing automatic, fluid, intentional text manipulation. The gap between "knowing that `diw` exists" and "instinctively using `diw` in context" is where most learning paths fail.

Praxis bridges that gap through deliberate practice: repeated, focused, measurable execution of composable skills.

## What Makes Praxis Different

- **Challenge-driven.** Every skill is taught through a concrete editing challenge with a verifiable success condition.
- **Composable.** Motions combine with operators combine with text objects. Praxis teaches the composition, not just the parts.
- **Terminal-first.** No GUI. No browser. The learning environment is the real environment.
- **Measurable.** Every challenge has a pass/fail outcome. Progress is tracked. Mastery is earned, not claimed.

## Design Principles

1. **Do not optimize for consumption.** Optimize for execution.
2. **Teach composition, not isolation.** `diw` is meaningful. `d` alone is not.
3. **Verification must be automatic.** If a human needs to judge success, the feedback loop is too slow.
4. **Content and engine are separate.** The Go engine defines challenges. The Lua frontend validates them. Either can evolve independently.
5. **Stable IDs are a public contract.** Challenge IDs must never change once released. Names may evolve; IDs do not.
