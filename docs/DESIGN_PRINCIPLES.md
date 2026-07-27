# Design Principles

**Purpose:** Structural decisions that govern Praxis's educational architecture. These are distinct from learner-facing policies (in CURRICULUM.md). They guide curriculum design and are enforced (where possible) by the test suite.

**Normative:** Changes to the curriculum should preserve these principles unless the principles themselves are intentionally revised through a documented design decision. This document is the architectural reference for future curriculum evolution.

---

## 1. Three-Layer Architecture

Every challenge belongs to exactly one layer. Each layer answers one learner need:

| Layer | Learner need | What it does |
|---|---|---|
| Tutorial | "I don't know." | Teaches one concept through guided instruction. |
| Training | "I know, but I'm slow." | Builds fluency through deliberate repetition of known compositions. |
| Trial | "I think I know." | Validates transfer — judgment and adaptability in an unprescribed context. |

No fourth mode exists. Any proposed feature that cannot be placed under one of these three needs is rejected.

---

## 2. Educational Dependency Invariant

Dependencies propagate in one direction only:

```
Tutorial → Training → Trial
```

- A Training challenge may depend only on concepts taught in Tutorial.
- A Trial may depend on concepts taught in Tutorial or practiced in Training.
- Tutorial challenges never depend on Training or Trial content.

Every composite challenge declares its dependencies via a `derived_from` field. This invariant is checked by the test suite to ensure the curriculum remains coherent.

---

## 3. Core Tutorial

Core Tutorial is the minimal set of challenges required for learner independence. It is:

- **Mandatory** — must be completed before Training or Trial progression is meaningful.
- **Frozen** — the set of 10 Core challenges is stable. No new challenge is added to Core without removing one.
- **Dependency-complete** — every concept referenced by Training or Trial is taught in Core.

The contrast between Core and Optional Tutorial is a tier, not a separate mode. Optional challenges teach additional concepts without blocking progression.

---

## 4. One Concept Per Tutorial

Each Tutorial challenge teaches exactly one concept (its primary concept). This principle keeps Tutorial focused and prevents compound lessons from confusing the learner.

The primary concept and context are documented in the curriculum metadata (`curriculum.go`). Concepts may reappear across challenges when they are taught in a different context or composed with other concepts.

---

## 5. Training as Composition

Training challenges do not teach new mechanics. They require the learner to compose two or more previously-taught primitives into a repeatable technique.

Training uses the composite validator, which enforces:
- Buffer-editing correctness (byte-exact result matching).
- A `MaxMoves` threshold that prevents brute-force solutions.

A Training challenge that teaches a new primitive instead of composing known ones violates the layer contract.

---

## 6. Trial as Transfer

Trials validate transfer, not recall. The learner must:
- **Judge** which composition fits the scenario (no keystrokes are prescribed).
- **Adapt** by applying the composition in an unprescribed context.

Like Training, Trials use the composite validator with a `MaxMoves` threshold and a `derived_from` field documenting the prerequisite concepts.

---

## 7. Stability Before Growth

IDs are permanent. Replay verification is mandatory. The curriculum grows only when a documented primary concept cannot be taught through any existing challenge. Duplicate concepts are prohibited unless the new context or composition justifies the addition.

These invariants are enforced by the content and integrity test suites.

---

## 8. Guidance Never Restricts

Every suggestion (NextChallenge, RecommendedReview) is optional. The learner can open any challenge by ID at any time. Nothing is gated. Progression is count-based, not stage-gated. The catalog is flat and unordered. Control mechanisms (forced orders, skill trees, difficulty ratings) are anti-goals.
