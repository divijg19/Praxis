package content

import "github.com/divijg19/Praxis/internal/challenge"

type Description struct {
	challenge.Challenge
	Stage       string   `json:"stage"`
	Concept     string   `json:"concept"`
	Context     string   `json:"context"`
	DerivedFrom []string `json:"derived_from,omitempty"`
}

func DescriptionFor(id string) (Description, bool) {
	for _, c := range All() {
		if c.ID == id {
			d := Description{Challenge: c}
			if m, ok := metadataFor(id); ok {
				d.Stage = m.Stage
				d.Concept = m.Concept
				d.Context = m.Context
				d.DerivedFrom = m.DerivedFrom
			}
			return d, true
		}
	}
	return Description{}, false
}
