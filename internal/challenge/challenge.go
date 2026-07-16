package challenge

type Evaluation struct {
	MaxMoves int `json:"max_moves"`
}

type Challenge struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Verify     string      `json:"verify"`
	Content    []string    `json:"content"`
	Result     []string    `json:"result"`
	Target     string      `json:"target"`
	Layer      string      `json:"layer"`
	Evaluation *Evaluation `json:"evaluation,omitempty"`
}
