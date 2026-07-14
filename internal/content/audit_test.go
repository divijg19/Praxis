package content

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// H1 — every Tutorial teaches exactly the primitive its metadata claims.
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

// H2 — the replay harness must exercise every challenge. Guard against a
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

// H3 — single owner for the verify/result/target shape invariant.
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
