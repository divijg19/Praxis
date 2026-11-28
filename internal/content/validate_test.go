package content

import (
	"strings"
	"testing"

	"github.com/divijg19/Praxis/internal/challenge"
)

// TestValidateAll is the single place the challenge model is validated: every
// challenge must satisfy Validate. Misconfigurations fail here, fast.
func TestValidateAll(t *testing.T) {
	for _, c := range All() {
		m, ok := metadataFor(c.ID)
		if !ok {
			t.Errorf("challenge %s: missing metadata", c.ID)
			continue
		}
		if problems := validate(c, m); len(problems) > 0 {
			t.Errorf("challenge %s invalid:\n  %v", c.ID, problems)
		}
	}
}

// TestRejectsInvalidChallenges proves Validate detects each malformed field.
func TestRejectsInvalidChallenges(t *testing.T) {
	baseChallenge := challenge.Challenge{
		ID:      "x",
		Name:    "Use x",
		Verify:  "cursor",
		Content: []string{"Use x."},
		Target:  "★",
		Layer:   "Tutorial",
	}
	baseMetadata := Metadata{Concept: "x", Context: "test", Stage: stageMovement}

	modify := func(fn func(*challenge.Challenge)) challenge.Challenge {
		c := baseChallenge
		fn(&c)
		return c
	}
	withMetadata := func(fn func(*Metadata)) Metadata {
		m := baseMetadata
		fn(&m)
		return m
	}

	tests := []struct {
		name string
		c    challenge.Challenge
		m    Metadata
		want string
	}{
		{"empty id", modify(func(c *challenge.Challenge) { c.ID = "" }), baseMetadata, "empty id"},
		{"empty name", modify(func(c *challenge.Challenge) { c.Name = "" }), baseMetadata, "empty name"},
		{"empty content", modify(func(c *challenge.Challenge) { c.Content = nil }), baseMetadata, "empty content"},
		{"bad stage", baseChallenge, withMetadata(func(m *Metadata) { m.Stage = "Bogus" }), "invalid or empty stage"},
		{"empty concept", baseChallenge, withMetadata(func(m *Metadata) { m.Concept = "" }), "empty concept"},
		{"empty verify", modify(func(c *challenge.Challenge) { c.Verify = "" }), baseMetadata, "empty verify"},
		{"cursor without target", modify(func(c *challenge.Challenge) { c.Target = "" }), baseMetadata, "cursor challenge missing target"},
		{"buffer without result", modify(func(c *challenge.Challenge) { c.Verify = "buffer" }), baseMetadata, "missing result"},
		{"unknown verify", modify(func(c *challenge.Challenge) { c.Verify = "bogus" }), baseMetadata, "unknown verify"},
		{"trial without derived_from", modify(func(c *challenge.Challenge) { c.Layer = "Trial" }), baseMetadata, "trial missing derived_from"},
		{"tutorial with hints", modify(func(c *challenge.Challenge) {
			c.Hints = []challenge.Hint{{Text: "Look closer.", Cost: 1}}
		}), baseMetadata, "tutorial must not declare hints"},
		{"empty hint text", modify(func(c *challenge.Challenge) {
			c.Layer = "Training"
			c.Hints = []challenge.Hint{{Text: "", Cost: 1}}
		}), baseMetadata, "empty text"},
		{"invalid hint cost", modify(func(c *challenge.Challenge) {
			c.Layer = "Training"
			c.Hints = []challenge.Hint{{Text: "Look closer.", Cost: 0}}
		}), baseMetadata, "invalid cost"},
		{"too many hint rungs", modify(func(c *challenge.Challenge) {
			c.Layer = "Training"
			c.Hints = []challenge.Hint{
				{Text: "a", Cost: 1}, {Text: "b", Cost: 1},
				{Text: "c", Cost: 1}, {Text: "d", Cost: 1},
			}
		}), baseMetadata, "1..3 rungs"},
		{"non-increasing hint costs", modify(func(c *challenge.Challenge) {
			c.Layer = "Training"
			c.Hints = []challenge.Hint{
				{Text: "a", Cost: 2}, {Text: "b", Cost: 1},
			}
		}), baseMetadata, "costs must increase"},
		{"duplicate hint text", modify(func(c *challenge.Challenge) {
			c.Layer = "Training"
			c.Hints = []challenge.Hint{
				{Text: "Look closer.", Cost: 1}, {Text: "Look closer.", Cost: 2},
			}
		}), baseMetadata, "duplicate hint text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problems := validate(tt.c, tt.m)
			if !strings.Contains(strings.Join(problems, "\n"), tt.want) {
				t.Fatalf("problems = %v, want one containing %q", problems, tt.want)
			}
		})
	}
}

func TestValidateAcceptsBaseline(t *testing.T) {
	c := challenge.Challenge{
		ID:      "x",
		Name:    "Use x",
		Verify:  "cursor",
		Content: []string{"Use x."},
		Target:  "★",
		Layer:   "Tutorial",
	}
	m := Metadata{Concept: "x", Context: "test", Stage: stageMovement}
	if problems := validate(c, m); len(problems) != 0 {
		t.Errorf("valid baseline reported problems: %v", problems)
	}
}
