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
| Replay          | Curriculum correctness, end-to-end content | `tools/replay/replay.lua`, `tools/verify/main.go` |
| Learner Journey | End-to-end experience, real product       | `tools/journey/journey.lua`, `tools/verify/main.go` |

Replay validates **content correctness**. The learner journey validates the **experience** (navigation, recovery, solving, completion). They are different layers; both must pass.

---

## Correctness

Unit-level behavior: `start`, `complete`, `verify`, `next`, `stats`, `describe`, `reset`, mastery and review selection, and mark computation.

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
- corrupted, missing, or unsupported-version `progress.json` degrades gracefully

## Replay

Checks all 48 challenge definitions by rendering each in a scratch buffer and asserting its target/result contract. It executes a canonical Vim keystroke trace for every challenge and asserts the resulting goal state and verified key count. This is curriculum correctness, not the learner experience.

```bash
# Requires PRAXIS_BIN, PRAXIS_DEV=1, and a clean Neovim.
# Prefer `go run ./tools/verify`, which sets all of this up.
PRAXIS_DEV=1 nvim --headless -u NONE -l tools/replay/replay.lua
```

## Learner Journey

Executes the real product, namely `:Praxis`, `:Praxis <id>`, and real keystrokes, from first launch to completion and back. It validates recovery paths, voice, and navigation. It runs inside `go run ./tools/verify` as one command covering all five test categories.

```bash
PRAXIS_DEV=1 nvim --headless -u NONE -l tools/journey/journey.lua
```

---

## Verification

```bash
go run ./tools/verify
```

Runs: `go build`, `go vet`, `gofmt`, `golangci-lint`, `go test`, replay (all 48 challenges), and the learner journey. Any failed check stops the gate with a non-zero exit. This is the single command to verify all five categories, namely Correctness, Integrity, Regression, Replay, and Learner Journey, before committing.

External binaries (`golangci-lint`, `nvim`) are resolved from `PATH`; they are not pinned as module dependencies. Install the linter with `make install-tools` (CI does the same). The gate builds one `praxis` binary into a temporary directory and passes it to both harnesses.

---

## Release Procedure

1. **Verify**: run `go run ./tools/verify`. It runs build, lint, format, vet, tests, replay, and journey. All checks must pass.
2. **Build**: run `go build ./...` so all packages compile.
3. **Documentation**: if content changed, update the relevant doc under `docs/` (the challenge catalog is available at runtime via `praxis catalog`).
4. **Stage**: run `git add -A && git status` to verify staged files.
5. **Commit**: write a descriptive message with a title (version plus summary), a body (categorized changes), and a discipline section (what did NOT change).
6. **Tag**: run `git tag <version>` so it matches the release plan.
7. **Push and Release**: run `git push origin <branch> <version>`; only **after** the push, create the pull request and the release/verification issues on GitHub. Create issues in the open state and do not close them.

Every release follows the same process. Do not skip steps.
