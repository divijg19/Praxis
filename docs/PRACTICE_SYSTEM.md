# Practice System

**Release:** v0.3.5
**Series theme:** How should Praxis live?
**Preceding releases:**

- v0.3.3 (Educational Architecture): How should Praxis teach?
- v0.3.4 (Practice Architecture): What promise does each mode make?
- v0.3.5 (Practice System): How should Praxis live as a long-term deliberate-practice environment?

This document is the primary deliverable of v0.3.5. The release is an
architectural investigation. It answers the questions Praxis must settle before
any capability is added. Where the answer is already clear, a small, conservative
implementation follows. Where it is not, the question is recorded with its
conclusion and left for a later release.

---

## Guiding Principle

Every investigation must eliminate uncertainty. No feature should exist merely
because it seems useful. Every implementation must emerge from a clearly
established educational principle. The v0.3.x series removes more uncertainty
than capability at each step; v0.3.5 is the final foundational release in that
tradition.

---

## Question I: What is a practice session?

**Today:** The product presents one challenge at a time. Completion leads to the
next challenge through `praxis next`. There is no explicit session boundary.

**Findings:**

- A session is not a new runtime construct. It is a learner-held intention: "I
  am here to practice for a while." Forcing a session abstraction into the
  product would add ceremony without educational benefit.
- Practice already flows continuously: Challenge to Result to Next Challenge.
  The rhythm is correct. The only gap is that the learner may not *feel* the
  boundary between "practicing" and "not practicing."
- A session is therefore *several challenges pursued under one intention*. It
  begins when the learner opens Praxis and ends when they leave. It does not
  need a start button, an end screen, or a timer.
- Too much practice is real but self-correcting: fatigue is felt, not computed.
  Praxis should never meter practice.

**Answer:** A practice session is the span between opening Praxis and leaving
it, composed of as many challenges as the learner chooses. Praxis must protect
the continuity of that span, not enclose it. No session object is added.

---

## Question II: What is Training?

**Findings:**

- Training exists to build *fluency*, not to teach mechanics. Mechanics belong
  to Tutorial. This distinction is already enforced in wording (REFERENCE.md,
  CURRICULUM.md) and must be enforced in feeling.
- Among speed, confidence, fluency, composition, recognition, and efficiency,
  the *primary* outcome is **fluency**: the ability to act correctly without
  deliberation. The others are consequences:
  - Recognition is a prerequisite fluency assumes is already internalized.
  - Composition is fluency applied to combined primitives.
  - Speed and efficiency are fluency observed from outside.
  - Confidence is fluency felt from inside.
- Therefore Training's contract is singular: **Training refines fluency.** Every
  Training challenge should be experienceable as "I know this; now make it
  automatic."

**Answer:** Training exists to refine fluency. Speed, confidence, composition,
efficiency, and recognition are consequences of fluency, not separate goals.
This is a permanent definition, not a description of current content.

---

## Question III: What is a Trial?

**Findings:**

- A Trial currently proves *transfer*: can the learner select the right
  composition for an unprescribed scenario? That is its core purpose and must
  remain.
- Of judgment, adaptability, efficiency, and confidence, the ones a Trial
  legitimately measures are *judgment* (choosing the composition) and
  *adaptability* (solving the same class of problem in varied forms). Efficiency
  and confidence are side effects the learner reads off their own performance.
- Trials must not become examinations. An examination has a right answer and a
  grade. A Trial has a solved buffer and a feeling of "I used the right tool."
  The difference is emotional, not technical: no score, no pass/fail, no ranking.

**Answer:** A Trial validates transfer, specifically judgment and adaptability,
by presenting authentic editing problems without prescribing keystrokes. It is
not a test. It is a proof to the learner that they can apply what Training
built.

---

## Question IV: What is progress?

**Findings:**

- Current signals (attempts, completions, mastery, confidence) are adequate as
  *records*. The open question is what they should *represent* to the learner.
- Progress is not a single number. It is the gradual shrinking of uncertainty
  across four dimensions:
  - **Breadth:** how many concepts have been touched.
  - **Depth:** how automatic each has become (mastery tiers).
  - **Consistency:** how reliably success arrives (confidence).
  - **Transfer:** how freely the learner applies a concept in Trials.
- Algorithms remain unchanged. The clarification is that "progress" is the
  learner's growing sense that fewer things feel uncertain, not a march toward
  a finish line. Curriculum Complete is an achievement, not the goal (already
  stated in CURRICULUM.md).

**Answer:** Progress is the reduction of uncertainty across breadth, depth,
consistency, and transfer. The stored signals already capture this; the learner
facing should present progress as "what is still uncertain," never as a
percentage toward done.

---

## Question V: What keeps someone practicing?

**Findings:**

- Extrinsic rewards (streaks, achievements, points) are rejected. They train
  return, not practice.
- Intrinsic reasons people return, observed and defensible:
  - **Competence:** "editing feels easier than it used to."
  - **Fewer mistakes:** "I no longer fumble the obvious cases."
  - **Confidence:** "I know which key to press without thinking."
  - **Curiosity:** "what happens if I combine these?"
- These are all *self-evident improvements the learner can feel*. Praxis
  reinforces them by making improvement visible (best moves, completions) without
  ranking it. The product should let the learner notice their own progress, not
  announce it.

**Answer:** Learners return because practice makes editing visibly easier.
Praxis supports this by surfacing quiet evidence of improvement, never by
rewarding return itself.

---

## Question VI: How should Praxis guide?

**Findings:**

- This absorbs every deferred v0.3.4 topic: recommendations, learner agency,
  exploration, optionality, review.
- The permanent principle: **Praxis reduces uncertainty; it never reduces
  autonomy.** Guidance is a suggestion the learner can ignore. Control is a
  restriction the learner cannot.
- Concrete implications already present and now ratified:
  - `NextChallenge` suggests the next unfinished challenge; the Catalog lets the
    learner go anywhere.
  - `RecommendedReview` suggests a review; the Hub shows it as optional (`[r]`).
  - Additional Lessons are browseable, never required.
- No recommendation engine is needed. The curriculum order plus the review
  suggestion already guide without controlling.

**Answer:** Praxis guides by suggesting a sensible next step and always permitting
departure. Every guidance surface is optional. Autonomy is a permanent right of
the learner.

---

## Question VII: What is the Hub?

**Findings:**

- The Hub is currently a dashboard: it reports Current, Progress, Direction,
  Mastery. That is accurate but cold. It answers "where am I?" but not "why am I
  here?"
- Its proper educational role is **home**: the place a practitioner returns to,
  where the next meaningful action is immediately obvious and the weight of
  status is minimal.
- Transition requirements:
  - Reduce informational noise (Mastery distribution is useful context, not a
    daily readout; it can be compact or secondary).
  - Make the next action obvious (`[Enter]` to practice, `[r]` to review).
  - Replace the feeling of "reading a report" with "stepping back into the
    practice space."
- The Hub is not a launcher, not a reflection essay, not a companion. It is the
  threshold of the practice space.

**Answer:** The Hub is the learner's home: the threshold of practice, not a
dashboard. It should make returning feel like stepping back in, and make the
next action unavoidable to miss.

---

## Question VIII: What should practice feel like?

**Findings (audit of `:Praxis` -> `Challenge Complete`):**

- **Pacing:** Currently good. Challenge to Result to Next is continuous; Enter
  never routes through the Hub (correct).
- **Rhythm:** A natural rhythm is Training, Training, Training, Trial, repeating.
  Repeated Training builds fluency; an occasional Trial proves it transferred.
  The product need not enforce this; `NextChallenge` already favors finishing
  what is started, which naturally clusters related practice.
- **Transitions:** The Result screen is the one place momentum could improve. It
  currently shows moves/time/completions plus retry/continue/back. That is
  sufficient. No modal ceremony should be added.
- **Friction:** The only real friction is re-reading instructional text the
  learner already knows. Instruction should remain present (it anchors the
  challenge) but never grow.
- **Repetition:** Repetition is the mechanism of fluency. It should feel
  intentional, not tedious. The completion count ("Completed N times") frames
  repetition as progress, which is correct.
- **Variety:** Variety arrives through the learner's own choices (Catalog,
  review), not through randomization the product imposes.

**Answer:** Practice should feel continuous, quiet, and self-directed: open,
practice, complete, continue, leave. No new rhythm engine; the existing flow
already embodies it. The work is removing small frictions and tonal noise, not
adding structure.

---

## Question IX: What should never change?

**Permanent invariants (architectural laws):**

1. **Deterministic.** Same input, same behavior. No randomness in scoring,
   ordering, or verification.
2. **Local-first.** All state is a JSON file under XDG. No network, no account,
   no telemetry.
3. **Educational before technical.** Every technical decision serves an
   educational contract; never the reverse.
4. **One learner need per mode.** Tutorial (I don't know), Training (I know but
   am slow), Trials (I think I know). Any feature that fits none of these three
   sentences is out of scope.
5. **Simplicity compounds.** A simpler system is permanently preferable to a
   more capable one that costs complexity.
6. **Practice belongs to the learner.** Guidance suggests; the learner decides.
7. **Three pillars only.** Tutorial, Training, Trials. No fourth mode, no
   sub-product presented as a mode.

These are now constitutional. Future releases obey them or justify the change
explicitly.

---

## Question X: What remains unverified?

**Audit of philosophy against tests:**

- **Journey coverage guard.** The learner journey harness hard-codes
  `Progress: 49/49`. If the curriculum grows, the harness passes with a stale
  count. *Fixed in this release*: `TestJourneyCoverage` (mirrors
  `TestReplayCoverage`) asserts the harness total equals `len(All())`.
- **Replay coverage.** Already guarded by `TestReplayCoverage`.
- **Instructional voice consistency.** `TestInstructionLineTerminates` enforces
  period-terminated instructions. The Hub/About/Result copy now shares a
  unified, period-final, second-person voice; this is a documentation/test
  concern, not a runtime contract, and is maintained by review.
- **Contract verification.** `Validate()` + `TestValidateAll` cover the
  verify/result/target shape for every challenge. `DescriptionFor` is the single
  source consumed by CLI, Lua, and replay.
- **Stats-text contract.** `hub.lua` parses `praxis stats` output. This is a
  declared breaking change (REFERENCE.md). No automated guard exists beyond the
  existing `stats` command tests; a future guard could assert the exact output
  schema, but the current tests cover the fields the Hub reads.

**Conclusion:** The only material blind spot was journey coverage. It is now
closed. No other philosophy-dependent contract lacks a test.

---

## Brainstorming (generated, not committed)

Possible future shapes, recorded without evaluation: practice sessions,
maintenance sessions, warm-ups, cooldowns, editing scenarios, challenge
playlists, thematic practice, mixed drills, composition drills, weakness
practice, review practice, daily routines. None are adopted in v0.3.5.

---

## Architectural Synthesis

Each candidate was weighed against educational quality, conceptual simplicity,
implementation cost, determinism, longevity, and maintainability.

- **Session object:** rejected. Adds ceremony; contradicts "practice belongs to
  the learner."
- **Recommendation engine:** rejected. Guidance already exists via `NextChallenge`
  and `RecommendedReview`; an engine would reduce autonomy.
- **Rhythm/enforcement scheduler:** rejected. The natural Training/Trial cadence
  needs no engine.
- **Hub dashboard upgrade:** rejected in favor of *reduction* (home, not more
  report).
- **Unified instructional voice + journey guard + Hub-as-home refinement:**
  accepted. Removes uncertainty, adds negligible complexity, aligns product with
  philosophy.

---

## Permanent Principles (the educational constitution)

1. Every mode answers exactly one learner need.
2. Tutorial creates independence.
3. Training builds fluency.
4. Trials validate transfer.
5. Guidance never becomes restriction.
6. Simplicity compounds.
7. Practice belongs to the learner.
8. Deterministic, local-first, educational before technical.

---

## Implementation Scope (v0.3.5)

Intentionally conservative. Investigation dominates; code follows only where
answers are clear.

### Verification

- Add `TestJourneyCoverage` so both replay and learner journey harnesses are
  protected against curriculum drift.

### Hub

- Complete the conceptual transition from dashboard to practice home through
  wording and layout refinement, not new functionality.

### UI

- Audit and unify instructional copy across onboarding, Hub, challenge screens,
  and completion screens. Keep the unified, period-final, second-person voice.

### Documentation

- Align docs with the finalized philosophy: Practice System, the seven permanent
  principles, and the "how should Praxis live?" framing.

### Explicit Non-Goals

Recommendation engines, adaptive learning, spaced repetition, unlock systems,
mastery algorithms, AI coaching, achievements, streaks, gamification, dynamic
difficulty. None are implemented.

---

## Success Criteria

v0.3.5 is complete when Praxis can answer, with confidence:

- What is a practice session? (Several challenges under one intention; no object.)
- What is Training fundamentally trying to improve? (Fluency.)
- What is a Trial fundamentally validating? (Transfer: judgment and adaptability.)
- What should progress represent? (Reduction of uncertainty across breadth,
  depth, consistency, transfer.)
- Why do learners return voluntarily? (Practice makes editing visibly easier.)
- How should Praxis guide without controlling? (Suggest always; permit departure.)
- What role should the Hub play? (Home: the threshold of practice.)
- What should deliberate practice feel like? (Continuous, quiet, self-directed.)
- Which educational principles are now permanent? (The seven above.)
- Does every implemented behavior reinforce those principles? (Yes, by audit.)
