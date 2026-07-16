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

**Checklist**

See [docs/REFERENCE.md → "Adding a New Challenge"](../../docs/REFERENCE.md) for the full step-by-step checklist (registration, metadata, stable-ID test, replay list, and `tools/verify.sh`).

- [ ] Implemented per the checklist above
- [ ] `tools/verify.sh` passes
