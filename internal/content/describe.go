package content

import (
	"encoding/json"

	"github.com/divijg19/Praxis/internal/challenge"
)

type Description struct {
	ID         string                `json:"id"`
	Name       string                `json:"name"`
	Verify     string                `json:"verify"`
	Layer      string                `json:"layer"`
	Stage      string                `json:"stage"`
	Concept    string                `json:"concept"`
	Context    string                `json:"context"`
	Target     string                `json:"target"`
	Content    []string              `json:"content"`
	Result     []string              `json:"result"`
	Evaluation *challenge.Evaluation `json:"evaluation,omitempty"`
}

func Describe(id string) (string, bool) {
	for _, c := range All() {
		if c.ID == id {
			d := Description{
				ID:         c.ID,
				Name:       c.Name,
				Verify:     c.Verify,
				Layer:      c.Layer,
				Target:     c.Target,
				Content:    c.Content,
				Result:     c.Result,
				Evaluation: c.Evaluation,
			}
			if m, ok := MetadataFor(id); ok {
				d.Stage = m.Stage
				d.Concept = m.Concept
				d.Context = m.Context
			}
			data, err := json.Marshal(d)
			if err != nil {
				return "", false
			}
			return string(data), true
		}
	}
	return "", false
}
