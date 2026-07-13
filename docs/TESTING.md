# Testing

**Purpose:** Testing philosophy for Praxis.

---

## Tier 1 — Structural Invariants

If these fail, Praxis is fundamentally broken. Examples: unique challenge IDs, curriculum reachability, DerivedFrom acyclicity, replay pass/fail. Never remove. Protect permanently.

## Tier 2 — Learning Invariants

If these fail, a learner-visible behavior changed. Examples: NextChallenge progression, mastery calculations, stats updates, instruction lines. Preserve across releases; test for regressions.

## Tier 3 — Release-Specific

If these fail, a release-phase constraint was violated. Examples: layer distribution (41/10/5), ID stability, challenge count (56). Candidates for removal when the release constraint is no longer meaningful.

---

## Verification

```bash
tools/verify.sh
```

Runs: `go test`, `go vet`, `gofmt`, replay (56 challenges). This is the single command to verify all tiers before committing.
