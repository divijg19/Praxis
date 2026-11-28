package content

import (
	"fmt"
	"slices"
	"strings"

	"github.com/divijg19/Praxis/internal/challenge"
)

// validStage reports whether s is one of the curriculum stages defined in
// curriculum.go, so the taxonomy has a single owner.
func validStage(s string) bool {
	return slices.Contains(stages, s)
}

// ValidateAll checks the compiled-in catalog at runtime. Content is static, so
// a failure is a build defect; the CLI reports it on stderr and exits rather
// than proceeding with an invalid catalog.
func ValidateAll() error {
	for _, c := range All() {
		m, ok := metadataFor(c.ID)
		if !ok {
			return fmt.Errorf("challenge %s: missing metadata", c.ID)
		}
		if problems := validate(c, m); len(problems) > 0 {
			return fmt.Errorf("challenge %s: %s", c.ID, strings.Join(problems, "; "))
		}
	}
	return nil
}

// validateHints enforces authored, finite, costed ladders. Tutorial lessons
// rely on the instruction line and must not declare hints.
func validateHints(c challenge.Challenge, problems *[]string) {
	if c.Layer == challenge.LayerTutorial {
		if len(c.Hints) > 0 {
			*problems = append(*problems, fmt.Sprintf("%s: tutorial must not declare hints", c.ID))
		}
		return
	}
	if len(c.Hints) == 0 || len(c.Hints) > 3 {
		*problems = append(*problems, fmt.Sprintf("%s: hint ladder must have 1..3 rungs", c.ID))
		return
	}
	seen := make(map[string]bool, len(c.Hints))
	prevCost := 0
	for i, h := range c.Hints {
		if h.Text == "" {
			*problems = append(*problems, fmt.Sprintf("%s: hint %d has empty text", c.ID, i+1))
		}
		if seen[h.Text] {
			*problems = append(*problems, fmt.Sprintf("%s: duplicate hint text", c.ID))
		}
		seen[h.Text] = true
		if h.Cost <= 0 || h.Cost > 3 {
			*problems = append(*problems, fmt.Sprintf("%s: hint %d has invalid cost %d", c.ID, i+1, h.Cost))
		} else if h.Cost <= prevCost {
			*problems = append(*problems, fmt.Sprintf("%s: hint costs must increase", c.ID))
		}
		prevCost = h.Cost
	}
}

// validate reports well-formedness problems for a single challenge. It is the
// single source of challenge-model validation; the gate runs it over every
// challenge in tests and the CLI runs ValidateAll at startup so a
// misconfiguration fails fast and in one place.
func validate(c challenge.Challenge, m Metadata) []string {
	var problems []string

	if c.ID == "" {
		problems = append(problems, "empty id")
	}
	if c.Name == "" {
		problems = append(problems, fmt.Sprintf("%s: empty name", c.ID))
	}
	if len(c.Content) == 0 {
		problems = append(problems, fmt.Sprintf("%s: empty content", c.ID))
	}
	if !validStage(m.Stage) {
		problems = append(problems, fmt.Sprintf("%s: invalid or empty stage %q", c.ID, m.Stage))
	}
	if m.Concept == "" {
		problems = append(problems, fmt.Sprintf("%s: empty concept", c.ID))
	}

	switch c.Verify {
	case "":
		problems = append(problems, fmt.Sprintf("%s: empty verify", c.ID))
	case challenge.VerifyCursor:
		if c.Target == "" {
			problems = append(problems, fmt.Sprintf("%s: cursor challenge missing target", c.ID))
		}
	case challenge.VerifyBuffer, challenge.VerifyComposite:
		if len(c.Result) == 0 {
			problems = append(problems, fmt.Sprintf("%s: %s challenge missing result", c.ID, c.Verify))
		}
	default:
		problems = append(problems, fmt.Sprintf("%s: unknown verify %q", c.ID, c.Verify))
	}

	if c.Layer == challenge.LayerTrial && len(m.DerivedFrom) == 0 {
		problems = append(problems, fmt.Sprintf("%s: trial missing derived_from", c.ID))
	}

	validateHints(c, &problems)
	return problems
}
