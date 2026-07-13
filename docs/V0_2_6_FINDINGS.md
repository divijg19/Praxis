# v0.2.6 Findings

**Question:** Can a learner progress through Praxis without external guidance?

---

## Findings

### Tutorials were exercises, not tutorials

The audit found 31/41 Tutorials lacked explicit instruction. They assessed rather than taught. The most impactful change in v0.2.6 was rewriting instruction lines from commands ("Delete the word") into teaching prompts ("Use daw to delete the middle word").

### Instruction text has more impact than move budgets

Rewriting 41 Content[0] strings changed the learner experience more than any validator, budget, or metadata change could. The text the learner reads is the curriculum.

### Six instruction archetypes were discovered

Every Tutorial maps to exactly one of: Reach, Locate, Select, Remove, Replace, Transfer. These form a complete taxonomy of the instruction patterns in Praxis. No Tutorial maps to more than one archetype; no archetype exists without at least one Tutorial.

### Onboarding is introduction, not navigation

The original onboarding was a gate ("Press Enter to begin"). Real onboarding establishes a mental model: what Praxis is, how the curriculum works, what the learner should do first. Onboarding should answer three questions: What is Praxis? What will I be doing? What should I do first?

### Reset capability is required for learner testing

A learning environment must support starting over without filesystem surgery. `praxis reset` is a product requirement, not a developer convenience.

### Replay proves correctness, not teaching quality

56/56 replay tests pass after instruction revisions. This proves the edits are structurally correct. It does not prove the tutorials teach well. Teaching quality must be validated by learner observation, not automated tests.

### Architecture is the learning model, not the content engine

The learning model is Observe → Practice → Apply. Content, validators, layers, and metadata serve that model — not the reverse. The UI presents each challenge as one buffer the learner solves and continues; it does not separate the three stages into distinct modes, but the model still shapes content and ordering.
