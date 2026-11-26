package content

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/divijg19/Praxis/internal/challenge"
)

// H1: every Tutorial teaches exactly the primitive its metadata claims.
// The curriculum is the executable specification: metadata, instruction, replay,
// tests, and docs must all agree. Tutorials each teach one primitive,
// so its Concept must appear in the instruction line. Topic labels (utf8,
// hjkl) are whitelisted because they describe a family, not a literal key.
func TestConceptInstructionAgreement(t *testing.T) {
	whitelist := map[string]bool{"utf8": true, "hjkl": true}
	for _, c := range All() {
		if c.Layer != "Tutorial" {
			continue
		}
		m, ok := metadataFor(c.ID)
		if !ok {
			continue
		}
		if whitelist[m.Concept] {
			continue
		}
		if len(c.Content) == 0 {
			t.Errorf("challenge %q has no instruction line", c.ID)
			continue
		}
		if !strings.Contains(c.Content[0], m.Concept) {
			t.Errorf("challenge %q: instruction %q does not teach its Concept %q", c.ID, c.Content[0], m.Concept)
		}
	}
}

// H2: the replay contract check must cover every challenge. Guard against a
// forgotten ID in replay.lua's hand-maintained all_ids; replay verifies data
// contracts, not keystroke solvability.
func TestReplayCoverage(t *testing.T) {
	data, err := os.ReadFile("../../tools/replay/replay.lua")
	if err != nil {
		t.Fatalf("cannot read replay.lua: %v", err)
	}
	raw := string(data)
	start := strings.Index(raw, "local all_ids = {")
	if start < 0 {
		t.Fatal("all_ids table not found in replay.lua")
	}
	block := raw[start:]
	if end := strings.Index(block, "\n}"); end >= 0 {
		block = block[:end]
	}

	replayIDs := map[string]bool{}
	for _, m := range regexp.MustCompile(`"([^"]+)"`).FindAllStringSubmatch(block, -1) {
		replayIDs[m[1]] = true
	}

	for _, c := range All() {
		if !replayIDs[c.ID] {
			t.Errorf("challenge %q is missing from replay.lua all_ids", c.ID)
		}
	}
	for id := range replayIDs {
		if !Exists(id) {
			t.Errorf("replay.lua all_ids contains unknown id %q", id)
		}
	}
}

// H2b: the learner journey harness completes the full curriculum and asserts
// "So far: N challenges completed." where N is derived from `praxis stats`.
// Guard that dynamic derivation so a future edit cannot silently restore a
// stale literal count when the curriculum grows or shrinks.
func TestJourneyCoverage(t *testing.T) {
	data, err := os.ReadFile("../../tools/journey/journey.lua")
	if err != nil {
		t.Fatalf("cannot read journey.lua: %v", err)
	}
	raw := string(data)

	// The harness derives its completion count from `praxis stats`, so there is
	// no literal count to compare. Guard that dynamic path instead.
	if !strings.Contains(raw, `Challenges Completed: (%d+)/(%d+)`) {
		t.Fatal("journey.lua does not derive its completion count from praxis stats")
	}
	if !strings.Contains(raw, `So far: " .. completed_total .. " challenges completed%.`) {
		t.Fatal("journey.lua completion assertion is not derived from praxis stats")
	}
}

// H3: single owner for the verify/result/target shape invariant.
func TestVerifyResultTargetInvariant(t *testing.T) {
	for _, c := range All() {
		switch c.Verify {
		case "cursor":
			if c.Target == "" {
				t.Errorf("cursor challenge %q has empty Target", c.ID)
			}
			if len(c.Result) > 0 {
				t.Errorf("cursor challenge %q has unexpected Result", c.ID)
			}
		case "buffer", "composite":
			if c.Target != "" {
				t.Errorf("buffer-like challenge %q has non-empty Target: %q", c.ID, c.Target)
			}
			if len(c.Result) == 0 {
				t.Errorf("buffer-like challenge %q has empty Result", c.ID)
			}
		}
	}
}

// H4: a challenge name must match its identity. The "Unnamed Register Hunter"
// regression showed names can drift from their lesson; the name must contain a
// recognizable form of the challenge's human-readable purpose, not the bare ID.
func TestChallengeNameNotBareID(t *testing.T) {
	for _, c := range All() {
		if c.Name == "" {
			t.Errorf("challenge %q has empty name", c.ID)
			continue
		}
		// Names like "Find + Delete Word" legitimately differ from the ID;
		// the guard is only against placeholder/empty-looking names.
		if strings.Contains(strings.ToLower(c.Name), "unnamed") ||
			strings.EqualFold(c.Name, c.ID) {
			t.Errorf("challenge %q has a placeholder or ID-like name %q", c.ID, c.Name)
		}
	}
}

// H6: Core is finite and mandatory, not a frozen numeric count. Counts are
// outcomes; the invariants are that every Core id exists, is a Tutorial, and
// teaches a distinct concept.
func TestCoreTutorialInvariants(t *testing.T) {
	ids := CoreTutorialIDs()
	if len(ids) == 0 {
		t.Fatal("Core Tutorial is empty")
	}
	byID := make(map[string]challenge.Challenge, len(All()))
	for _, c := range All() {
		byID[c.ID] = c
	}
	concepts := make(map[string]string, len(ids))
	for _, id := range ids {
		c, ok := byID[id]
		if !ok {
			t.Errorf("Core id %q is not in the catalog", id)
			continue
		}
		if c.Layer != "Tutorial" {
			t.Errorf("Core challenge %q has Layer=%q, must be Tutorial", id, c.Layer)
		}
		m, ok := metadataFor(id)
		if !ok {
			t.Errorf("Core challenge %q has no curriculum metadata", id)
			continue
		}
		if other, duplicate := concepts[m.Concept]; duplicate {
			t.Errorf("Core challenges %q and %q share concept %q", id, other, m.Concept)
		}
		concepts[m.Concept] = id
	}
	if !IsCoreTutorial("delete_inner_word_hunter") {
		t.Error("diw must be Core: Training and Trial compose it")
	}
}

func conceptTaught(concept string, taught map[string]bool) bool {
	if taught[concept] {
		return true
	}
	for prefix := 1; prefix < len(concept); prefix++ {
		operator, object := concept[:prefix], concept[prefix:]
		if !taught[object] {
			continue
		}
		for c := range taught {
			if strings.HasPrefix(c, operator) {
				return true
			}
		}
	}
	return false
}

// H6b: no composite challenge may require a concept absent from Tutorial.
func TestCompositePrerequisitesCovered(t *testing.T) {
	taught := make(map[string]bool)
	for _, c := range All() {
		if c.Layer != "Tutorial" {
			continue
		}
		if m, ok := metadataFor(c.ID); ok {
			taught[m.Concept] = true
		}
	}
	for _, c := range All() {
		if c.Layer != "Training" && c.Layer != "Trial" {
			continue
		}
		m, ok := metadataFor(c.ID)
		if !ok {
			t.Errorf("%q has no curriculum metadata", c.ID)
			continue
		}
		for _, concept := range strings.Split(m.Concept, "/") {
			if !conceptTaught(concept, taught) {
				t.Errorf("%s (%s) requires concept %q not taught in Tutorial", c.ID, c.Layer, concept)
			}
		}
	}
}

// H7: Trial instructions must not prescribe specific keystrokes.
func TestTrialInstructionDoesNotPrescribe(t *testing.T) {
	prescriptiveTerms := []string{
		"use ", "Use ", "press ", "Press ", "type ", "Type ",
		"then ", "hit ", "Hit ",
	}
	for _, c := range All() {
		if c.Layer != "Trial" {
			continue
		}
		if len(c.Content) == 0 {
			continue
		}
		line := c.Content[0]
		for _, term := range prescriptiveTerms {
			if strings.Contains(line, term) {
				t.Errorf("Trial %q instruction %q contains prescriptive term %q; Trials must not prescribe keystrokes", c.ID, line, term)
			}
		}
	}
}

// H5: every instruction line ends with a period, enforcing one editorial
// voice. Trial challenges use imperative goal statements ("Remove the third
// word.") which also end with a period, so the rule is uniform across layers.
func TestInstructionLineTerminates(t *testing.T) {
	for _, c := range All() {
		if len(c.Content) == 0 {
			continue
		}
		line := c.Content[0]
		if !strings.HasSuffix(line, ".") {
			t.Errorf("challenge %q instruction %q does not end with a period", c.ID, line)
		}
	}
}

// H8: Training instructions must follow the consistent "Use X to Y" pattern
// that frames the challenge as fluency practice, not concept teaching.
func TestTrainingInstructionConsistent(t *testing.T) {
	for _, c := range All() {
		if c.Layer != "Training" {
			continue
		}
		if len(c.Content) == 0 {
			continue
		}
		line := c.Content[0]
		if !strings.HasPrefix(line, "Use ") {
			t.Errorf("Training %q instruction %q does not start with 'Use '; Training should prescribe the technique for fluency practice", c.ID, line)
		}
	}
}
