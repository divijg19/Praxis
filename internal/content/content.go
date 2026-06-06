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
	}
}
