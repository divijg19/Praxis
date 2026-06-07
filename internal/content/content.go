package content

import "github.com/divijg19/Praxis/internal/challenge"

func All() []challenge.Challenge {
	return []challenge.Challenge{
		{
			ID:     "motion_rush",
			Name:   "Motion Rush",
			Target: "★",
			Content: []string{
				"Move your cursor to the star ★",
			},
		},
		{
			ID:     "grid_rush",
			Name:   "Grid Rush",
			Target: "★",
			Content: []string{
				". . . . .",
				". . . ★ .",
				". . . . .",
			},
		},
		{
			ID:     "find_hunter",
			Name:   "Find Hunter",
			Target: "★",
			Content: []string{
				"find motions are fast",
				"",
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa★",
			},
		},
		{
			ID:     "word_hunter",
			Name:   "Word Hunter",
			Target: "★",
			Content: []string{
				"word motions are fast",
				"",
				"alpha beta gamma delta epsilon ★",
			},
		},
		{
			ID:     "symbol_hunter",
			Name:   "Symbol Hunter",
			Target: "@",
			Content: []string{
				"find the target symbol",
				"",
				".......................@",
			},
		},
	}
}
