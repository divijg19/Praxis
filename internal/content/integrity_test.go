package content

import (
	"sort"
	"strings"
	"testing"

	"github.com/divijg19/Praxis/internal/stats"
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
		m, ok := metadataFor(c.ID)
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

func equalStringSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sorted := make([]string, len(a))
	copy(sorted, a)
	sort.Strings(sorted)

	other := make([]string, len(b))
	copy(other, b)
	sort.Strings(other)

	for i := range sorted {
		if sorted[i] != other[i] {
			return false
		}
	}
	return true
}

func TestDerivedFromAcyclic(t *testing.T) {
	adj := make(map[string][]string)
	for _, c := range All() {
		m, ok := metadataFor(c.ID)
		if !ok {
			continue
		}
		if len(m.DerivedFrom) > 0 {
			adj[c.ID] = m.DerivedFrom
		}
	}

	white := make(map[string]bool)
	gray := make(map[string]bool)
	black := make(map[string]bool)

	for id := range adj {
		white[id] = true
	}

	var dfs func(id string) bool
	dfs = func(id string) bool {
		delete(white, id)
		gray[id] = true

		for _, dep := range adj[id] {
			if gray[dep] {
				return true
			}
			if white[dep] {
				if dfs(dep) {
					return true
				}
			}
		}

		delete(gray, id)
		black[id] = true
		return false
	}

	for id := range adj {
		if white[id] {
			if dfs(id) {
				t.Errorf("DerivedFrom cycle detected involving %q", id)
			}
		}
	}
}

func TestCurriculumReachable(t *testing.T) {
	challenges := All()

	if len(challenges) == 0 {
		t.Fatal("All() returned empty")
	}

	seen := make(map[string]bool)
	ids := make([]string, len(challenges))
	for i, c := range challenges {
		if seen[c.ID] {
			t.Fatalf("duplicate ID %q at index %d", c.ID, i)
		}
		seen[c.ID] = true
		ids[i] = c.ID
	}

	m := make(map[string]stats.Stats)
	next := stats.NextChallenge(m, ids)
	if next != ids[0] {
		t.Fatalf("NextChallenge with empty stats: got %q, want %q", next, ids[0])
	}

	visited := make(map[string]bool)
	for i, id := range ids {
		m[id] = stats.Stats{Completions: 3}
		visited[id] = true

		if i < len(ids)-1 {
			next = stats.NextChallenge(m, ids)
			if next != ids[i+1] {
				t.Fatalf("after completing %q (index %d): NextChallenge = %q, want %q", id, i, next, ids[i+1])
			}
		} else {
			next = stats.NextChallenge(m, ids)
			if next != "" {
				t.Fatalf("after final challenge %q: NextChallenge = %q, want \"\"", id, next)
			}
		}
	}

	for _, id := range ids {
		if !visited[id] {
			t.Errorf("challenge %q was never visited during traversal", id)
		}
	}
}

func TestTrialObjectivesDistinct(t *testing.T) {
	type trialMeta struct {
		id   string
		meta Metadata
	}
	var trials []trialMeta
	for _, c := range All() {
		m, ok := metadataFor(c.ID)
		if !ok {
			continue
		}
		if c.Layer == "Trial" {
			trials = append(trials, trialMeta{c.ID, m})
		}
	}

	for i := 0; i < len(trials); i++ {
		for j := i + 1; j < len(trials); j++ {
			a, b := trials[i], trials[j]
			if equalStringSet(a.meta.DerivedFrom, b.meta.DerivedFrom) &&
				a.meta.Concept == b.meta.Concept &&
				a.meta.Context == b.meta.Context {
				t.Errorf("Trials %q and %q have identical DerivedFrom, Concept, and Context:\n  DerivedFrom: %v\n  Concept: %q\n  Context: %q",
					a.id, b.id, a.meta.DerivedFrom, a.meta.Concept, a.meta.Context)
			}
		}
	}
}

func TestTutorialTiersAssigned(t *testing.T) {
	var core, optional int
	for _, c := range All() {
		if c.Layer != "Tutorial" {
			if TutorialTier(c.ID) != "" {
				t.Errorf("non-Tutorial %q has a Tutorial tier %q", c.ID, TutorialTier(c.ID))
			}
			continue
		}
		switch TutorialTier(c.ID) {
		case "core":
			core++
		case "optional":
			optional++
		default:
			t.Errorf("Tutorial %q has no Tier assigned", c.ID)
		}
	}
	// Core is the intentionally small mandatory onboarding: a handful of
	// exercises, not the bulk of Tutorial. Optional should be larger.
	if core < 5 || core > 15 {
		t.Errorf("Core Tutorial count %d is outside the intended 5 to 15 range", core)
	}
	if optional < core {
		t.Errorf("Optional Tutorial count %d should be at least the Core count %d", optional, core)
	}
}

func TestCoreTutorialSetConsistent(t *testing.T) {
	ids := CoreTutorialIDs()
	if len(ids) == 0 {
		t.Fatal("Core Tutorial set is empty")
	}
	for _, id := range ids {
		if !IsCoreTutorial(id) {
			t.Errorf("CoreTutorialIDs contains %q but IsCoreTutorial disagrees", id)
		}
		if _, ok := metadataFor(id); !ok {
			t.Errorf("Core Tutorial %q has no metadata", id)
		}
		if TutorialTier(id) != "core" {
			t.Errorf("Core Tutorial %q is not marked Tier=core", id)
		}
	}
}
