package content

import "testing"

// TestValidateAll is the single place the challenge model is validated: every
// challenge must satisfy Validate. Misconfigurations fail here, fast.
func TestValidateAll(t *testing.T) {
	for _, c := range All() {
		m, ok := metadataFor(c.ID)
		if !ok {
			t.Errorf("challenge %s: missing metadata", c.ID)
			continue
		}
		if problems := Validate(c, m); len(problems) > 0 {
			t.Errorf("challenge %s invalid:\n  %v", c.ID, problems)
		}
	}
}
