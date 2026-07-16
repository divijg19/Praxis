name: Challenge proposal
about: Propose a new challenge for the Praxis curriculum (extension model)
title: "[challenge] "
labels: challenge
assignees: ''
---

**Proposed ID**
A short `snake_case` id, e.g. `trial_search_jump`. Must be unique.

**Layer / Category**
Tutorial, Training, or Trial? Which category (Cursor / Buffer / Composite)?

**What it teaches**
One sentence on the Vim skill being drilled.

**Validation brief**
- Verify mode (`cursor` / `buffer` / `composite`)
- For Trial: which challenges it is `derived_from`

**Checklist (see docs/REFERENCE.md "Adding a New Challenge")**
- [ ] Added to `All()` in `internal/content/content.go`
- [ ] Added metadata to `curriculum` in `internal/content/curriculum.go`
- [ ] Added to `stableChallengeIDs` / `stableChallengeNames` in `internal/content/content_test.go`
- [ ] Added to `all_ids` in `tools/replay/replay.lua`
- [ ] `tools/verify.sh` passes
