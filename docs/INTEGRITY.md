# Curriculum Integrity

## Principles

1. **One concept per challenge.** Every challenge has a documented primary Concept, Context, and Stage. The curriculum metadata in `internal/content/taxonomy_test.go` is the single source of truth, enforced by `TestCurriculumContextsComplete` and `TestCurriculumMapComplete`.

2. **IDs are permanent.** Challenge identifiers must never be renamed or removed. This is enforced by `TestChallengeIDsStable`.

3. **Validators are stable.** Every registered validator must be used by at least one challenge. This is enforced by `TestNoValidatorDrift`.

4. **Curriculum growth is intentional.** New challenges require a documented primary concept and must not duplicate existing content. This is enforced by `TestCurriculumMapComplete` and `TestNoDuplicateChallengeContent`.

5. **Replay verification is mandatory.** Every challenge must be solvable via the replay harness. This is enforced by `tools/replay/replay.sh`.

6. **Duplicate challenge content is prohibited.** Challenges may share concepts when they teach the concept in a different context or composition. For example, both `paren_hunter` and `match_hunter` teach `%`, but one uses parentheses and the other uses brackets — the technique is the same, the context differs.

## Test Suite

| Test | Enforces |
|---|---|
| `TestCoreConceptCoverage` | Core Vim concepts remain represented |
| `TestNoDuplicateChallengeContent` | No unintended duplicate exercises |
| `TestCurriculumMapComplete` | Every challenge mapped, no orphaned entries |
| `TestUniqueChallengeIDs` | No ID collisions |
| `TestChallengeCount` | No accidental addition/removal |
| `TestNoValidatorDrift` | Validator usage stays current |
| `TestChallengeIDsStable` | IDs never renamed |
| `TestCurriculumContextsComplete` | Every challenge has Concept, Context, Stage |
| `TestConceptContextPairsUnique` | No duplicate (Concept, Context) pairs |
| `TestProgressionCoverage` | All progression stages have challenges |
| `TestStageIntroductionOrder` | Stages introduced in pedagogical order |

## Anti-Goals

- Challenge difficulty ratings
- Mastery scoring
- Skill trees or dependencies
- Progression gates
- Recommendation systems
