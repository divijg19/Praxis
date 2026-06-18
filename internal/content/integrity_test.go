package content

import (
	"strings"
	"testing"
)

func TestCoreConceptCoverage(t *testing.T) {
	seen := make(map[string]bool)
	for _, c := range All() {
		meta, ok := curriculum[c.ID]
		if !ok {
			t.Errorf("challenge %q missing from curriculum map", c.ID)
			continue
		}
		seen[meta.Concept] = true
	}

	families := []struct {
		name   string
		values []string
	}{
		{"fFtT", []string{"f", "F", "t", "T"}},
		{"wWeEeBb", []string{"w", "W", "e", "E", "b", "B"}},
		{"registers", []string{"\"a", "\"A", "\"ap"}},
	}
	for _, f := range families {
		var found bool
		for _, v := range f.values {
			if seen[v] {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("no challenge teaches a %q family concept (values: %v)", f.name, f.values)
		}
	}

	for _, concept := range []string{"dd", "dw", "diw", "daw", "ciw", "yy"} {
		if !seen[concept] {
			t.Errorf("core concept %q is not taught by any challenge", concept)
		}
	}
}

func TestNoDuplicateChallengeContent(t *testing.T) {
	t.Run("content", func(t *testing.T) {
		seen := make(map[string]string)
		for _, c := range All() {
			key := strings.Join(c.Content, "\n")
			if first, ok := seen[key]; ok {
				t.Errorf("duplicate content between %q and %q", first, c.ID)
			}
			seen[key] = c.ID
		}
	})

	t.Run("buffer_results", func(t *testing.T) {
		seen := make(map[string]string)
		for _, c := range All() {
			if len(c.Result) == 0 {
				continue
			}
			key := strings.Join(c.Result, "\n")
			if first, ok := seen[key]; ok {
				t.Errorf("duplicate buffer result between %q and %q", first, c.ID)
			}
			seen[key] = c.ID
		}
	})
}

func TestCurriculumMapComplete(t *testing.T) {
	if len(curriculum) != len(All()) {
		t.Errorf("curriculum map has %d entries, expected %d", len(curriculum), len(All()))
	}

	for _, c := range All() {
		meta, ok := curriculum[c.ID]
		if !ok {
			t.Errorf("challenge %q missing from curriculum map", c.ID)
			continue
		}
		if meta.Concept == "" {
			t.Errorf("challenge %q has empty curriculum concept", c.ID)
		}
	}

	for id := range curriculum {
		var found bool
		for _, c := range All() {
			if c.ID == id {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("curriculum map has unknown challenge %q", id)
		}
	}
}

func TestTrialIntegrity(t *testing.T) {
	allIDs := make(map[string]bool)
	for _, c := range All() {
		allIDs[c.ID] = true
	}

	for _, c := range All() {
		m, ok := MetadataFor(c.ID)
		if !ok {
			t.Errorf("challenge %q missing from curriculum", c.ID)
			continue
		}

		switch c.Layer {
		case "Trial":
			if c.Evaluation == nil {
				t.Errorf("Trial %q: missing Evaluation", c.ID)
			}
			if len(m.DerivedFrom) == 0 {
				t.Errorf("Trial %q: empty DerivedFrom", c.ID)
			}
			for _, dep := range m.DerivedFrom {
				if !allIDs[dep] {
					t.Errorf("Trial %q: DerivedFrom target %q not found in challenges", c.ID, dep)
				}
			}
		default:
			if len(m.DerivedFrom) != 0 {
				t.Errorf("non-Trial %q (layer=%q): unexpected DerivedFrom %v", c.ID, c.Layer, m.DerivedFrom)
			}
		}
	}
}
