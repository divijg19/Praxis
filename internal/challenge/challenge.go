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

type Challenge struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Verify  string   `json:"verify"`
	Content []string `json:"content"`
	Result  []string `json:"result"`
	Target  string   `json:"target"`
	Layer   string   `json:"layer"`
}
