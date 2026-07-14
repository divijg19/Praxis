# Testing

**Purpose:** How is Praxis verified?

---

## Categories

Every test belongs to exactly one category.

| Category        | Purpose                                  | Where |
| --------------- | ---------------------------------------- | ----- |
| Correctness     | Individual behavior of a unit or command  | `internal/**/*_test.go` |
| Integrity       | Curriculum invariants (IDs, reachability, acyclicity) | `internal/content/*_test.go` |
| Regression      | Prevent previously-fixed bugs from returning | `internal/**/*_test.go` |
| Replay          | Curriculum correctness, end-to-end content | `tools/replay/replay.lua`, `tools/verify.sh` |
| Learner Journey | End-to-end experience, real product       | `tools/journey/journey.lua`, `tools/journey/journey.sh` |

Replay validates **content correctness**. The learner journey validates the **experience** (navigation, recovery, solving, completion). They are different layers; both must pass.

---

## Correctness

Unit-level behavior: `next`, `stats`, `attempt`, `record`, `describe`, `reset`, mastery calculations, confidence, recommended-review selection.

Run:

```bash
go test ./...
```

## Integrity

Curriculum-wide invariants. If these fail, the product is fundamentally broken:

- unique challenge IDs
- curriculum reachability (`praxis next` walks every challenge)
- `DerivedFrom` acyclicity (no cycles)
- every target and result exists

These live in `internal/content/integrity_test.go` and must never be removed.

## Regression

Behavior that once broke and must stay fixed. Each has a dedicated regression test. Examples:

- no two challenges share an ID
- `describe` on an unknown id returns a clean error (not a crash)
- corrupted or missing `stats.json` degrades gracefully

## Replay

Drives all 56 challenges through the CLI and asserts each opens, solves, and completes. This is curriculum correctness, not the learner experience.

```bash
# Requires the `praxis` binary on PATH (the harness builds it to /tmp/praxis).
# Prefer `tools/verify.sh`, which builds it for you.
nvim --headless -l tools/replay/replay.lua
```

## Learner Journey

Executes the real product — `:Praxis`, `:Praxis <id>`, real keystrokes — from first launch to completion and back. Validates recovery paths, voice, and navigation. Runs inside `tools/verify.sh` as one command covering all five test categories.

```bash
bash tools/journey/journey.sh
```

---

## Verification

```bash
tools/verify.sh
```

Runs: `go test`, `go vet`, `gofmt`, replay (all 56 challenges), and the learner journey. This is the single command to verify all five categories — Correctness, Integrity, Regression, Replay, and Learner Journey — before committing.
