package content

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestHintCostsIncrease(t *testing.T) {
	for _, c := range All() {
		n := HintCount(c.ID)
		if n == 0 {
			continue
		}
		prev := 0
		for level := 1; level <= n; level++ {
			got := HintCostThrough(c.ID, level)
			if got <= prev {
				t.Errorf("%s: cumulative hint cost did not increase at rung %d (%d <= %d)", c.ID, level, got, prev)
			}
			prev = got
		}
		if HintCostThrough(c.ID, n+10) != prev {
			t.Errorf("%s: out-of-range hint level must clamp to ladder length", c.ID)
		}
	}
}

func TestHintTextDoesNotLeakSolutions(t *testing.T) {
	banned := []string{"use ", "Use ", "press ", "Press ", "type ", "Type ", "then ", "hit ", "Hit "}
	for _, c := range All() {
		for i, h := range c.Hints {
			tr, hasTrace := TraceFor(c.ID)
			if hasTrace && strings.Contains(h.Text, tr.Keys) {
				t.Errorf("%s hint %d contains the canonical trace", c.ID, i+1)
			}
			if len(c.Result) > 1 {
				for _, line := range c.Result[1:] {
					if line != "" && strings.Contains(h.Text, line) {
						t.Errorf("%s hint %d contains the expected result", c.ID, i+1)
					}
				}
			}
			if c.Layer == "Trial" {
				for _, term := range banned {
					if strings.Contains(h.Text, term) {
						t.Errorf("%s Trial hint %d prescribes %q", c.ID, i+1, term)
					}
				}
			}
		}
	}
}

func TestHintTextIsNotSerialized(t *testing.T) {
	for _, c := range All() {
		if len(c.Hints) == 0 {
			continue
		}
		data, err := json.Marshal(DescriptionForMust(c.ID))
		if err != nil {
			t.Fatal(err)
		}
		for _, h := range c.Hints {
			if strings.Contains(string(data), h.Text) {
				t.Errorf("%s: hint text leaked into description JSON", c.ID)
			}
		}
		if !strings.Contains(string(data), `"hints":`) {
			t.Errorf("%s: hint costs missing from description JSON", c.ID)
		}
	}
}

func DescriptionForMust(id string) Description {
	d, ok := DescriptionFor(id)
	if !ok {
		panic("missing description: " + id)
	}
	return d
}
