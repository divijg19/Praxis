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
- No generation of hints or adaptive difficulty.
