// Package score contains the only mark computation in Praxis. Marks describe a
// verified run; they are never summed, ranked, or used to gate completion.
package score

// Marks returns 1, 2, or 3 marks for a verified run. Hint costs are key
// equivalents applied here, never in the frontend.
func Marks(keys, par, hintCost int) int {
	if par <= 0 || keys < 0 || hintCost < 0 {
		return 1
	}
	effective := keys + hintCost
	if effective <= par {
		return 3
	}
	margin := par / 2
	if margin < 2 {
		margin = 2
	}
	if effective <= par+margin {
		return 2
	}
	return 1
}
