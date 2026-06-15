package content

import "testing"

func TestCurriculumCoverage(t *testing.T) {
	for _, c := range All() {
		_, ok := MetadataFor(c.ID)
		if !ok {
			t.Errorf("challenge %q missing from curriculum metadata", c.ID)
		}
	}
}

func TestValidStages(t *testing.T) {
	valid := ValidStages()
	for _, c := range All() {
		m, ok := MetadataFor(c.ID)
		if !ok {
			continue
		}
		if !valid[m.Stage] {
			t.Errorf("challenge %q has invalid Stage %q", c.ID, m.Stage)
		}
	}
}

func TestConceptContextPairsUnique(t *testing.T) {
	seen := make(map[string]string)
	for _, c := range All() {
		m, ok := MetadataFor(c.ID)
		if !ok {
			continue
		}
		key := m.Concept + "|" + m.Context
		if first, ok := seen[key]; ok {
			t.Errorf("duplicate (Concept, Context) pair between %q and %q: (%q, %q)",
				first, c.ID, m.Concept, m.Context)
		}
		seen[key] = c.ID
	}
}

func TestProgressionCoverage(t *testing.T) {
	stages := make(map[string]bool)
	for _, c := range All() {
		m, ok := MetadataFor(c.ID)
		if !ok {
			continue
		}
		stages[m.Stage] = true
	}

	expected := []string{StageMovement, StageSearch, StageStructuralNavigation,
		StageEditing, StageTextObjects, StageRegisters}
	for _, s := range expected {
		if !stages[s] {
			t.Errorf("progression stage %q has no challenges", s)
		}
	}
}

func TestStageIntroductionOrder(t *testing.T) {
	stageOrder := []string{
		StageMovement,
		StageSearch,
		StageStructuralNavigation,
		StageEditing,
		StageTextObjects,
		StageRegisters,
	}

	first := make(map[string]int)
	for i, c := range All() {
		m, ok := MetadataFor(c.ID)
		if !ok {
			continue
		}
		if _, ok := first[m.Stage]; !ok {
			first[m.Stage] = i
		}
	}

	for j := 1; j < len(stageOrder); j++ {
		prev := stageOrder[j-1]
		curr := stageOrder[j]
		if first[curr] < first[prev] {
			t.Errorf("stage %q (first at index %d) appears before %q (first at index %d)",
				curr, first[curr], prev, first[prev])
		}
	}
}

func TestMetadataMatchesChallengeLayer(t *testing.T) {
	for _, c := range All() {
		m, ok := MetadataFor(c.ID)
		if !ok {
			t.Errorf("challenge %q missing from curriculum metadata", c.ID)
			continue
		}
		if m.Layer != c.Layer {
			t.Errorf("challenge %q: Metadata.Layer=%q != Challenge.Layer=%q",
				c.ID, m.Layer, c.Layer)
		}
	}
}
