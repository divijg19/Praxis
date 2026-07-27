package content

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
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

// H2: the replay harness must exercise every challenge. Guard against a
// forgotten ID in replay.lua's hand-maintained all_ids (integrity
// principle: every challenge must be solvable via the replay harness).
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
// "Progress: N/N" where N is the challenge count. That total is hand-maintained
// in journey.lua and must stay in lockstep with content.All(). Guard against
// the harness silently passing with a stale count when the curriculum grows or
// shrinks (mirrors TestReplayCoverage for the journey harness).
func TestJourneyCoverage(t *testing.T) {
	data, err := os.ReadFile("../../tools/journey/journey.lua")
	if err != nil {
		t.Fatalf("cannot read journey.lua: %v", err)
	}
	raw := string(data)

	// Match the "Progress: 49/49" assertion the harness makes on completion.
	re := regexp.MustCompile(`Progress: (\d+)/(\d+)`)
	m := re.FindStringSubmatch(raw)
	if m == nil {
		t.Fatal("journey.lua completion count assertion (Progress: N/N) not found")
	}
	// Both halves must agree (it asserts a fully-completed curriculum).
	if m[1] != m[2] {
		t.Fatalf("journey.lua completion count is inconsistent: %s/%s", m[1], m[2])
	}
	harnessTotal, _ := strconv.Atoi(m[1])

	want := len(All())
	if harnessTotal != want {
		t.Errorf("journey.lua completion count %d does not match curriculum size %d; update the assertion in journey.lua", harnessTotal, want)
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

// H6: the Core Tutorial set is frozen at exactly 10 challenges.
func TestCoreTutorialCountFrozen(t *testing.T) {
	const want = 10
	got := len(CoreTutorialIDs())
	if got != want {
		t.Errorf("Core Tutorial has %d challenges, want %d; Core is frozen", got, want)
	}
	for _, id := range CoreTutorialIDs() {
		for _, c := range All() {
			if c.ID != id {
				continue
			}
			if c.Layer != "Tutorial" {
				t.Errorf("Core challenge %q has Layer=%q, must be Tutorial", id, c.Layer)
			}
			break
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
