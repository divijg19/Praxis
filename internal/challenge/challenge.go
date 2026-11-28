package challenge

// Verify values and Layers are part of the catalog contract. They are defined
// once here so content, validation, and consumers share the same vocabulary.
const (
	VerifyCursor    = "cursor"
	VerifyBuffer    = "buffer"
	VerifyComposite = "composite"

	LayerTutorial = "Tutorial"
	LayerTraining = "Training"
	LayerTrial    = "Trial"
)

// Hint is an authored, opt-in nudge. Text is never serialized: Description
// embeds Challenge, so JSON must expose only the cost of each rung.
type Hint struct {
	Text string `json:"-"`
	Cost int    `json:"cost"`
}

type Challenge struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Verify  string   `json:"verify"`
	Content []string `json:"content"`
	Result  []string `json:"result"`
	Target  string   `json:"target"`
	Layer   string   `json:"layer"`
	Hints   []Hint   `json:"hints,omitempty"`
}
