package content

import (
	"encoding/json"
	"testing"
)

func TestDescribeRoundTrip(t *testing.T) {
	for _, c := range All() {
		jsonStr, ok := Describe(c.ID)
		if !ok {
			t.Errorf("Describe(%q) returned false", c.ID)
			continue
		}
		var d Description
		if err := json.Unmarshal([]byte(jsonStr), &d); err != nil {
			t.Errorf("Describe(%q) returned invalid JSON: %v", c.ID, err)
			continue
		}
		if d.ID != c.ID {
			t.Errorf("Describe(%q) returned Description.ID=%q", c.ID, d.ID)
		}
		if d.Name != c.Name {
			t.Errorf("Describe(%q) returned Description.Name=%q, want %q", c.ID, d.Name, c.Name)
		}
		if d.Verify != c.Verify {
			t.Errorf("Describe(%q) returned Description.Verify=%q, want %q", c.ID, d.Verify, c.Verify)
		}
		if d.Layer != c.Layer {
			t.Errorf("Describe(%q) returned Description.Layer=%q, want %q", c.ID, d.Layer, c.Layer)
		}
	}
}

func TestDescribeEvaluation(t *testing.T) {
	for _, c := range All() {
		jsonStr, ok := Describe(c.ID)
		if !ok {
			continue
		}
		var d Description
		json.Unmarshal([]byte(jsonStr), &d)
		if c.Verify == "composite" {
			if d.Evaluation == nil {
				t.Errorf("Describe(%q): composite challenge has nil Evaluation", c.ID)
			} else if d.Evaluation.MaxMoves <= 0 {
				t.Errorf("Describe(%q): composite challenge has invalid MaxMoves=%d", c.ID, d.Evaluation.MaxMoves)
			}
		} else if d.Evaluation != nil {
			t.Errorf("Describe(%q): non-composite challenge has non-nil Evaluation", c.ID)
		}
	}
}

func TestDescribeUnknown(t *testing.T) {
	_, ok := Describe("nonexistent_challenge_id")
	if ok {
		t.Error("Describe(\"nonexistent_challenge_id\") returned true, expected false")
	}
}
