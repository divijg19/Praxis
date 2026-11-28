package content

// Hints are authored on challenges. Only costs cross the CLI boundary; hint
// text is revealed explicitly by `praxis hint`.

func HintCount(id string) int {
	c, ok := byID[id]
	if !ok {
		return 0
	}
	return len(c.Hints)
}

// HintCost returns the cost of a single 1-based rung, or 0 if invalid.
func HintCost(id string, rung int) int {
	c, ok := byID[id]
	if !ok || rung < 1 || rung > len(c.Hints) {
		return 0
	}
	return c.Hints[rung-1].Cost
}

// HintCostThrough returns the cumulative cost of all rungs up to and including
// the given level. Levels are clamped to the authored ladder so tampered state
// cannot inflate the cost.
func HintCostThrough(id string, level int) int {
	c, ok := byID[id]
	if !ok || level < 0 {
		return 0
	}
	if level > len(c.Hints) {
		level = len(c.Hints)
	}
	total := 0
	for i := 0; i < level; i++ {
		total += c.Hints[i].Cost
	}
	return total
}

// HintText returns an authored rung, or "" if invalid.
func HintText(id string, rung int) string {
	c, ok := byID[id]
	if !ok || rung < 1 || rung > len(c.Hints) {
		return ""
	}
	return c.Hints[rung-1].Text
}
