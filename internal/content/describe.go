package content

import "github.com/divijg19/Praxis/internal/challenge"

type Description struct {
	challenge.Challenge
	Stage       string   `json:"stage"`
	Concept     string   `json:"concept"`
	Context     string   `json:"context"`
	DerivedFrom []string `json:"derived_from,omitempty"`
	StartRow    int      `json:"start_row"`
	StartCol    int      `json:"start_col"`
	Par         int      `json:"par"`
}

func DescriptionFor(id string) (Description, bool) {
	c, ok := byID[id]
	if !ok {
		return Description{}, false
	}
	d := Description{Challenge: c, Par: ParFor(id)}
	if tr, ok := TraceFor(id); ok {
		d.StartRow, d.StartCol = tr.StartRow, tr.StartCol
	}
	if m, ok := metadataFor(id); ok {
		d.Stage = m.Stage
		d.Concept = m.Concept
		d.Context = m.Context
		d.DerivedFrom = m.DerivedFrom
	}
	return d, true
}
