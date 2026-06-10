# Releases

## Versioning

Praxis uses `v0.0.x` versioning during the foundation phase and `v0.1.x` for the public contracts phase. Each release is tagged and pushed to GitHub with a corresponding pair of issues.

## Release Pattern

Every release follows a two-issue (or three-issue) pattern:

1. **Release issue** — tracks the content/feature scope
2. **Verification issue** — checklist for testing and regression
3. **Archive issue** (occasional) — documents major decisions or migrations

## Creating a Release

See `docs/RELEASE_CHECKLIST.md` for the mandatory 10-step release process.

1. Implement all changes
2. Verify: `gofmt -l .` clean, `go build ./...`, `go vet ./...`, `go test ./...`
3. Run replay: `tools/replay/replay.sh`
4. Stage changes: `git add <files>`
5. Commit with descriptive message
6. Tag: `git tag v0.<major>.<N>`
7. Push: `git push origin v0.<major>.x v0.<major>.<N>`
8. Create GitHub issues for the release

## API Contract

**Challenge IDs are stable once released.** They must never be renamed or removed. The `yank_line_hunter` precedent established in v0.0.19 demonstrates that curriculum framing (Name, Instruction) can evolve while the identifier remains permanent.

This is enforced by `TestChallengeIDsStable` and `TestChallengeNamesStable`.

## Release History

| Tag | Title | Challenges | Tests |
|---|---|---|---|---|---|
| v0.0.1 | Initial prototype | 1 | 0 |
| v0.0.2 | Parameterize challenge content | 1 | 0 |
| v0.0.3 | Add grid_rush challenge | 2 | 0 |
| v0.0.4 | Extract Lua plugin from inline | 2 | 0 |
| v0.0.5 | Add find_hunter challenge | 3 | 0 |
| v0.0.6 | Export praxis target for Lua | 3 | 0 |
| v0.0.7 | Add word_hunter challenge | 4 | 0 |
| v0.0.8 | Lua frontend reads Go engine | 4 | 0 |
| v0.0.9 | Add grid_rush challenge | 5 | 0 |
| v0.0.10 | Add symbol_hunter challenge | 6 | 0 |
| v0.0.11 | Validate content-only extensibility | 7 | 0 |
| v0.0.12 | Parameterize completion target | 8 | 0 |
| v0.0.13 | Add line, paren, sentence hunters | 11 | 0 |
| v0.0.14 | Search Foundations pack | 14 | 0 |
| v0.0.15 | Structural Navigation Foundations | 19 | 0 |
| v0.0.16 | Validation Architecture | 19 | 3 |
| v0.0.17 | Editing Foundations | 25 | 5 |
| v0.0.18 | Structural Editing Foundations | 37 | 5 |
| v0.0.19 | Register Foundations | 41 | 15 |
| v0.0.20 | Foundation Freeze | 41 | 27 |
| v0.1.0 | Public Contracts & Stability | 41 | 33 |
| **v0.1.1** | **Progress Tracking** | **41** | **45** |
| **v0.1.2** | **Practice Sessions** | **41** | **45** |
| **v0.1.3** | **Curriculum Integrity & Audit** | **41** | **48** |
