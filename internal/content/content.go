package content

import "github.com/divijg19/Praxis/internal/challenge"

func All() []challenge.Challenge {
	return []challenge.Challenge{
		{
			ID:   "motion_rush",
			Name: "Motion Rush",
			Content: []string{
				"Move your cursor to the star ★",
			},
		},
		{
			ID:   "grid_rush",
			Name: "Grid Rush",
			Content: []string{
				". . . . .",
				". . . ★ .",
				". . . . .",
			},
		},
		{
			ID:   "find_hunter",
			Name: "Find Hunter",
			Content: []string{
				"find motions are fast",
				"",
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa★",
			},
		},
		{
			ID:   "word_hunter",
			Name: "Word Hunter",
			Content: []string{
				"word motions are fast",
				"",
				"alpha beta gamma delta epsilon ★",
			},
		},
	}
}
