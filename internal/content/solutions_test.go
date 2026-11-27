package content

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/divijg19/Praxis/internal/challenge"
)

func TestEveryChallengeHasCanonicalTrace(t *testing.T) {
	ids := IDs()
	known := make(map[string]bool, len(ids))
	for _, id := range ids {
		known[id] = true
		tr, ok := TraceFor(id)
		if !ok {
			t.Errorf("challenge %q has no canonical solution trace", id)
			continue
		}
		if tr.StartRow < 1 || tr.StartCol < 0 {
			t.Errorf("challenge %q has invalid start %d/%d", id, tr.StartRow, tr.StartCol)
			continue
		}
		c, ok := byID[id]
		if !ok {
			t.Errorf("trace for unknown challenge %q", id)
			continue
		}
		if tr.StartRow > len(c.Content) {
			t.Errorf("challenge %q starts at row %d but has %d lines", id, tr.StartRow, len(c.Content))
		} else if tr.StartCol >= utf8.RuneCountInString(c.Content[tr.StartRow-1]) {
			t.Errorf("challenge %q starts at column %d beyond its line", id, tr.StartCol)
		}
		if tr.Keys == "" {
			t.Errorf("challenge %q has an empty solution trace", id)
		}
		if got := ParFor(id); got != utf8.RuneCountInString(tr.Keys) {
			t.Errorf("ParFor(%q) = %d, want %d", id, got, utf8.RuneCountInString(tr.Keys))
		}
		if d, ok := DescriptionFor(id); !ok || d.Par != ParFor(id) {
			t.Errorf("DescriptionFor(%q) par = %d, want %d", id, d.Par, ParFor(id))
		}
	}
	for id := range traces {
		if !known[id] {
			t.Errorf("orphan solution trace for unknown challenge %q", id)
		}
	}
}

// TestTracesUseTaughtConcept enforces the fidelity rule in TRACES_AND_PAR.md:
// a lesson that names a Vim command must actually use that command.
func TestTracesUseTaughtConcept(t *testing.T) {
	required := map[string]string{
		"di(": "di(", "da(": "da(", "di\"": "di\"", "da\"": "da\"",
		"diw": "diw", "daw": "daw", "ciw": "ciw", "dw": "dw", "dd": "dd",
		"d$": "d$", "yy": "yy", "\"a": "\"a", "\"ap": "\"ap",
		"ci(": "ci(", "ci\"": "ci\"", "ca\"": "ca\"",
		"x": "x", "r": "r", "~": "~", "%": "%", ")": ")", "}": "}",
		"f": "f", "w": "w", "j": "j", "/": "/", "?": "?", ";": ";",
	}
	whitelist := map[string]bool{"hjkl": true, "utf8": true}
	for _, c := range All() {
		// Trials test transfer: they deliberately do not prescribe commands.
		if c.Layer == challenge.LayerTrial {
			continue
		}
		m, ok := metadataFor(c.ID)
		if !ok {
			continue
		}
		concept := strings.Split(m.Concept, "/")[0]
		want, ok := required[concept]
		if !ok || whitelist[concept] {
			continue
		}
		tr, ok := TraceFor(c.ID)
		if !ok {
			continue
		}
		if !strings.Contains(tr.Keys, want) {
			t.Errorf("trace for %s (%s) does not use taught command %q: %q", c.ID, concept, want, tr.Keys)
		}
	}
}

func TestCursorTracesDoNotStartOnTarget(t *testing.T) {
	for _, c := range All() {
		if c.Verify != challenge.VerifyCursor {
			continue
		}
		tr, ok := TraceFor(c.ID)
		if !ok || tr.StartRow < 1 || tr.StartRow > len(c.Content) {
			continue
		}
		line := c.Content[tr.StartRow-1]
		if tr.StartCol > len(line) {
			t.Errorf("%s: start column beyond line", c.ID)
			continue
		}
		idx := utf8.RuneCountInString(line[:tr.StartCol])
		runes := []rune(line)
		if idx < len(runes) && string(runes[idx]) == c.Target {
			t.Errorf("%s: starts on target %q; any key would verify", c.ID, c.Target)
		}
	}
}
