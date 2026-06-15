# Releases

## Versioning

Praxis uses `v0.0.x` for the foundation phase, `v0.1.x` for the public contracts phase, and `v0.2.x` for the expand phase. Each release is tagged and pushed to GitHub with a corresponding pair of issues.

## Release Pattern

Every release follows a two-issue pattern:

1. **Release issue** — tracks the content/feature scope
2. **Verification issue** — checklist for testing and regression

## Creating a Release

See `REFERENCE.md` (Release Procedure section) for the mandatory process.

## API Contract

**Challenge IDs are stable once released.** They must never be renamed or removed. The `yank_line_hunter` precedent established in v0.0.19 demonstrates that curriculum framing (Name, Instruction) can evolve while the identifier remains permanent.

This is enforced by `TestChallengeIDsStable` and `TestChallengeNamesStable`.

## Release History

| Tag | Title | Challenges | Tests |
|---|---|---|---|
| v0.0.1 | Initial prototype | 1 | 0 |
| v0.0.2 | Parameterize challenge content | 1 | 0 |
| v0.0.3 | Establish Neovim-to-Go-CLI Bridge | 2 | 0 |
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
| **v0.1.4** | **Curriculum Taxonomy & Progression** | **41** | **52** |
| **v0.1.5** | **Mastery Labels** | **41** | **56** |
| **v0.1.6** | **Documentation Consolidation & Mastery Surfacing** | **41** | **58** |
| **v0.1.7** | **Practice Guidance & Documentation Hygiene** | **41** | **66** |
| **v0.1.8** | **Measurement Completion** | **41** | **71** |
| **v0.1.9** | **Execution Confidence & Freeze** | **41** | **79** |
| **v0.1.10** | **Documentation Audit & Consistency Pass** | **41** | **79** |
| **v0.2.0** | **Curriculum Taxonomy** | **41** | **82** |
| **v0.2.1** | **Arrival** | **41** | **82** |
| **v0.2.2** | **Production Curriculum Metadata** | **41** | **90** |
