# Traces and par

Every challenge has exactly one canonical solution trace in
`internal/content/solutions.go`. A trace is **internal harness data**: it is
never shown to a learner. `praxis trace <id>` exposes it only when
`PRAXIS_DEV=1`.

A trace contains:

- `start_row` / `start_col`: where the learner's cursor begins;
- `keys`: the exact Vim keys of the canonical solution.

Par is `utf8.RuneCountInString(keys)` — the number of physical keys in that
solution. It is recomputed by Go and included in `praxis describe` as `par`.

## Rules

1. Every catalog challenge must have a trace (`TestEveryChallengeHasCanonicalTrace`).
2. A trace must start inside the challenge content and within the line's length.
3. A trace must actually solve its challenge in a clean, config-free Neovim
   (`internal/verify` and `tools/replay/replay.lua`).
4. Par must equal the replayed key count; `describe.par` must equal computed par.
5. A trace must use the concept it claims to teach where that concept is a Vim
   command (e.g. a text-object lesson must use `di(` / `ci"`, not `f` alone).

## Units

Keys and edits are different and must never be compared directly:

- **Keys** are physical keypresses counted by `vim.on_key`, verified by clean
  replay, and used for par and marks.
- **Verified keys** are the shortest key prefix that reached the goal in clean
  Neovim. No separate edit budget exists.

The result card always labels both.
