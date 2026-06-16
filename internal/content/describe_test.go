package content

import (
	"reflect"
	"testing"
)

func TestDescriptionForCompleteness(t *testing.T) {
	for _, c := range All() {
		t.Run(c.ID, func(t *testing.T) {
			d, ok := DescriptionFor(c.ID)
			if !ok {
				t.Fatal("DescriptionFor returned false")
			}
			if d.ID != c.ID {
				t.Errorf("ID = %q, want %q", d.ID, c.ID)
			}
			if d.Name != c.Name {
				t.Errorf("Name = %q, want %q", d.Name, c.Name)
			}
			if d.Verify != c.Verify {
				t.Errorf("Verify = %q, want %q", d.Verify, c.Verify)
			}
			if d.Layer != c.Layer {
				t.Errorf("Layer = %q, want %q", d.Layer, c.Layer)
			}
			if d.Target != c.Target {
				t.Errorf("Target = %q, want %q", d.Target, c.Target)
			}
			if !reflect.DeepEqual(d.Content, c.Content) {
				t.Errorf("Content mismatch")
			}
			if !reflect.DeepEqual(d.Result, c.Result) {
				t.Errorf("Result mismatch")
			}
			m, ok := MetadataFor(c.ID)
			if !ok {
				t.Fatal("no metadata")
			}
			if d.Stage != m.Stage {
				t.Errorf("Stage = %q, want %q", d.Stage, m.Stage)
			}
			if d.Concept != m.Concept {
				t.Errorf("Concept = %q, want %q", d.Concept, m.Concept)
			}
			if d.Context != m.Context {
				t.Errorf("Context = %q, want %q", d.Context, m.Context)
			}
			if c.Verify == "composite" {
				if d.Evaluation == nil {
					t.Fatal("expected non-nil Evaluation for composite challenge")
				}
				if d.Evaluation.MaxMoves <= 0 {
					t.Errorf("MaxMoves = %d, want > 0", d.Evaluation.MaxMoves)
				}
				if d.Evaluation.MaxMoves != c.Evaluation.MaxMoves {
					t.Errorf("MaxMoves = %d, want %d", d.Evaluation.MaxMoves, c.Evaluation.MaxMoves)
				}
			} else if d.Evaluation != nil {
				t.Errorf("unexpected non-nil Evaluation for non-composite challenge")
			}
		})
	}
	for id := range curriculum {
		_, ok := DescriptionFor(id)
		if !ok {
			t.Errorf("DescriptionFor(%q) returned false for curriculum entry", id)
		}
	}
}

func TestDescriptionForUnknown(t *testing.T) {
	d, ok := DescriptionFor("nonexistent_challenge_id")
	if ok {
		t.Error("DescriptionFor returned true for unknown ID")
	}
	if d.ID != "" {
		t.Errorf("expected zero-value Description, got ID=%q", d.ID)
	}
}
