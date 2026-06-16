package challenge

type Evaluation struct {
	MaxMoves int `json:"max_moves"`
}

type Challenge struct {
	ID         string
	Name       string
	Verify     string
	Content    []string
	Result     []string
	Target     string
	Layer      string
	Evaluation *Evaluation
}
