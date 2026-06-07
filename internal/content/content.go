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
		{
			ID:     "line_hunter",
			Name:   "Line Hunter",
			Target: "★",
			Content: []string{
				"Move to the line with the star",
				"",
				"line one",
				"line two",
				"line three",
				"line four",
				"line five",
				"★ line six",
			},
		},
		{
			ID:     "paren_hunter",
			Name:   "Paren Hunter",
			Target: "★",
			Content: []string{
				"Jump to the matching paren",
				"",
				"(                         )★",
			},
		},
		{
			ID:     "sentence_hunter",
			Name:   "Sentence Hunter",
			Target: "★",
			Content: []string{
				"Jump between sentences",
				"",
				"First sentence ends here.",
				"Second. Third.",
				"★ Fourth. Fifth.",
			},
		},
		{
			ID:     "slash_hunter",
			Name:   "Slash Hunter",
			Target: "★",
			Content: []string{
				"Search forward to find the target",
				"",
				"alpha  bravo  charlie  delta  echo  foxtrot  golf  ★",
			},
		},
		{
			ID:     "question_hunter",
			Name:   "Question Hunter",
			Target: "★",
			Content: []string{
				"Search backward to find the target",
				"",
				"line one",
				"line two",
				"line three",
				"line four",
				"line five",
				"line six",
				"line seven",
				"line eight",
				"line nine",
				"line ten",
				"line eleven",
				"line twelve",
				"line thirteen",
				"line fourteen",
				"line fifteen",
				"line sixteen",
				"★ line seventeen",
				"line eighteen",
			},
		},
		{
			ID:     "repeat_hunter",
			Name:   "Repeat Hunter",
			Target: "@",
			Content: []string{
				"Search for ★, then repeat to find @",
				"",
				"★  ★  ★  ★  ★  ★  ★  ★  @",
			},
		},
	}
}
