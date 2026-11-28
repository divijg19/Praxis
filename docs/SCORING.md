# Scoring without a score

Praxis records **marks** for a verified run. Marks are not a score: they are
never summed, ranked, converted to a percentage, or required for completion.

## Verification chain

```
physical keys (vim.on_key) → praxis complete → clean config-free nvim replay
                            → exact goal state → Go computes marks → event
```

- The frontend never sends a mark or completion count.
- `praxis complete` refuses an empty or tainted key log.
- Replay ignores user config, mappings, plugins, and XDG state.
- `verified:false` writes nothing and shows an honest failure card.

## Marks formula

With `K` verified keys, `P` par, and `H` hint key-equivalents:

```
E = K + H
E <= P          → 3 marks
E <= P + margin → 2 marks   (margin = max(2, floor(P/2)))
otherwise       → 1 mark
```

Every verified completion earns at least one mark. Marks are monotone in both
`K` and `H`, property-tested in `internal/score`.

## What marks are not

- No XP, levels, badges, leaderboards, Elo, or aggregate score.
- No gate on progression; three completions still advance a challenge.
- No comparison with other learners.
- No generated, adaptive, or auto-revealed hints, and no hint gates
  progression. Authored, opt-in, costed hints are allowed (see below).

## Hints

Hints are **authored, finite, opt-in ladders** attached to a challenge. They are
never generated, adaptive, auto-revealed (for example on failure), or used to
gate progress. Tutorials declare no authored hints: the instruction line is
their only hint.

- A Training challenge has two rungs (costs 1 and 2); a Trial has three
  (costs 1, 2, and 3). Rung text is hand-written and must never reproduce a
  canonical keystroke trace or the expected result, and Trial rungs must not
  prescribe specific keystrokes.
- Previewing a rung is inert and shows only its cost. Revealing it persists the
  highest rung reached in `progress.json`; that persisted state is the single
  authority for `H`. Clients cannot supply `hint_cost`, `hint_level`, or marks.
- `praxis start` clears revealed hints, so every retry is a clean attempt. A
  revealed hint never changes a previously recorded event or mark.
