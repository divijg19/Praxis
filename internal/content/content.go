package content

import "github.com/divijg19/Praxis/internal/challenge"

func All() []challenge.Challenge {
	return []challenge.Challenge{
		{
			ID:     "grid_rush",
			Name:   "Grid Rush",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use h, j, k and l to move to the star ★.",
				". . . . .",
				". . . ★ .",
				". . . . .",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "find_hunter",
			Name:   "Find Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use f to jump to the next character.",
				"",
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa★",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "word_hunter",
			Name:   "Word Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use w to jump to the next word.",
				"",
				"alpha beta gamma delta epsilon ★",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "line_hunter",
			Name:   "Line Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use j and k to move to the line with the star.",
				"",
				"line one",
				"line two",
				"line three",
				"line four",
				"line five",
				"★ line six",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "paren_hunter",
			Name:   "Paren Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use % to jump to the matching parenthesis.",
				"",
				"(                         )★",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "sentence_hunter",
			Name:   "Sentence Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use ) to jump forward to the next sentence.",
				"",
				"First sentence ends here.",
				"Second. Third.",
				"★ Fourth. Fifth.",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "slash_hunter",
			Name:   "Slash Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use / to search forward for the target.",
				"",
				"alpha  bravo  charlie  delta  echo  foxtrot  golf  ★",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "question_hunter",
			Name:   "Question Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use ? to search backward for the target.",
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
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "repeat_hunter",
			Name:   "Repeat Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use ; to repeat your last find and reach the next star.",
				"",
				"★  ★  ★  ★  ★  ★  ★  ★  ★",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "inner_paren_hunter",
			Name:   "Inner Paren Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use i( to target inside the parentheses.",
				"",
				"(hello★)",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "around_paren_hunter",
			Name:   "Around Paren Hunter",
			Verify: challenge.VerifyCursor,
			Target: ")",
			Content: []string{
				"Use a( to target around the parentheses.",
				"",
				"(hello)★",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "inner_quote_hunter",
			Name:   "Inner Quote Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use i\" to target inside the quotes.",
				"",
				"x\"magnificent★\"",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "around_quote_hunter",
			Name:   "Around Quote Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use a\" to target around the quotes.",
				"",
				"x\"magnificent\"★",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "paragraph_hunter",
			Name:   "Paragraph Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use } to jump forward to the next paragraph.",
				"",
				"first paragraph",
				"",
				"second paragraph",
				"",
				"★ third paragraph",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "delete_character_hunter",
			Name:   "Delete Character Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use x to delete the extra letter.",
				"",
				"helllo",
			},
			Result: []string{
				"Use x to delete the extra letter.",
				"",
				"hello",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "replace_character_hunter",
			Name:   "Replace Character Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use r to replace the wrong letter.",
				"",
				"hallo",
			},
			Result: []string{
				"Use r to replace the wrong letter.",
				"",
				"hello",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "toggle_case_hunter",
			Name:   "Toggle Case Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use ~ to toggle the case of each letter.",
				"",
				"HELLO",
			},
			Result: []string{
				"Use ~ to toggle the case of each letter.",
				"",
				"hello",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "delete_word_hunter",
			Name:   "Delete Word Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use dw to delete the middle word.",
				"",
				"alpha bravo charlie",
			},
			Result: []string{
				"Use dw to delete the middle word.",
				"",
				"alpha charlie",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "change_word_hunter",
			Name:   "Change Word Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use ciw to change the word.",
				"",
				"foo",
			},
			Result: []string{
				"Use ciw to change the word.",
				"",
				"bar",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "utf8_cursor_hunter",
			Name:   "UTF-8 Cursor Hunter",
			Verify: challenge.VerifyCursor,
			Target: "★",
			Content: []string{
				"Use h to move past multi-byte characters to the star ★.",
				"",
				"α β γ ★",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "delete_line_hunter",
			Name:   "Delete Line Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use dd to delete the middle line.",
				"",
				"keep",
				"remove",
				"keep",
			},
			Result: []string{
				"Use dd to delete the middle line.",
				"",
				"keep",
				"keep",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "delete_to_end_hunter",
			Name:   "Delete To End Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use d$ to delete from cursor to end of line.",
				"",
				"keep this text remove_this_part",
			},
			Result: []string{
				"Use d$ to delete from cursor to end of line.",
				"",
				"keep this text",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "delete_inner_word_hunter",
			Name:   "Delete Inner Word Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use diw to delete inside the middle word.",
				"",
				"keep remove keep",
			},
			Result: []string{
				"Use diw to delete inside the middle word.",
				"",
				"keep  keep",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "delete_around_word_hunter",
			Name:   "Delete Around Word Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use daw to delete around the middle word.",
				"",
				"keep remove keep",
			},
			Result: []string{
				"Use daw to delete around the middle word.",
				"",
				"keep keep",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "delete_inner_paren_hunter",
			Name:   "Delete Inner Paren Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use di( to delete inside parentheses.",
				"",
				"func(remove)",
			},
			Result: []string{
				"Use di( to delete inside parentheses.",
				"",
				"func()",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "delete_around_paren_hunter",
			Name:   "Delete Around Paren Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use da( to delete around parentheses.",
				"",
				"func(remove)",
			},
			Result: []string{
				"Use da( to delete around parentheses.",
				"",
				"func",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "delete_inner_quote_hunter",
			Name:   "Delete Inner Quote Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use di\" to delete inside quotes.",
				"",
				"\"remove\"",
			},
			Result: []string{
				"Use di\" to delete inside quotes.",
				"",
				"\"\"",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "delete_around_quote_hunter",
			Name:   "Delete Around Quote Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use da\" to delete around quotes.",
				"",
				"\"remove\"",
			},
			Result: []string{
				"Use da\" to delete around quotes.",
				"",
				"",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "change_inner_paren_hunter",
			Name:   "Change Inner Paren Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use ci( to change inside parentheses.",
				"",
				"func(remove)",
			},
			Result: []string{
				"Use ci( to change inside parentheses.",
				"",
				"func(hello)",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "change_inner_quote_hunter",
			Name:   "Change Inner Quote Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use ci\" to change inside quotes.",
				"",
				"\"remove\"",
			},
			Result: []string{
				"Use ci\" to change inside quotes.",
				"",
				"\"hello\"",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "yank_line_hunter",
			Name:   "Yank Line Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use yy to yank, then p to paste.",
				"",
				"copy me",
			},
			Result: []string{
				"Use yy to yank, then p to paste.",
				"",
				"copy me",
				"copy me",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "word_register_hunter",
			Name:   "Word Register Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use \"ayw to store a word, then \"ap to paste it.",
				"",
				"alpha beta",
			},
			Result: []string{
				"Use \"ayw to store a word, then \"ap to paste it.",
				"",
				"alpha beta alpha",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "register_replace_hunter",
			Name:   "Register Replace Hunter",
			Verify: challenge.VerifyBuffer,
			Target: "",
			Content: []string{
				"Use \"ayy to store 'correct', then \"ap to replace 'wrong'.",
				"",
				"correct",
				"wrong",
			},
			Result: []string{
				"Use \"ayy to store 'correct', then \"ap to replace 'wrong'.",
				"",
				"correct",
				"correct",
			},
			Layer: challenge.LayerTutorial,
		},
		{
			ID:     "find_diw_combo",
			Name:   "Find + Delete Word",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Use f{char} then diw to delete the word.",
				"",
				"keep delete keep",
			},
			Result: []string{
				"Use f{char} then diw to delete the word.",
				"",
				"keep  keep",
			},
			Layer:      challenge.LayerTraining,
			Evaluation: &challenge.Evaluation{MaxMoves: 8},
		},
		{
			ID:     "find_daw_combo",
			Name:   "Find + Delete Around",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Use f{char} then daw to delete around the word.",
				"",
				"keep delete keep",
			},
			Result: []string{
				"Use f{char} then daw to delete around the word.",
				"",
				"keep keep",
			},
			Layer:      challenge.LayerTraining,
			Evaluation: &challenge.Evaluation{MaxMoves: 8},
		},
		{
			ID:     "find_di_paren_combo",
			Name:   "Find + Delete Inside",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Use f{char} then di( to delete inside the parentheses.",
				"",
				"func(remove)",
			},
			Result: []string{
				"Use f{char} then di( to delete inside the parentheses.",
				"",
				"func()",
			},
			Layer:      challenge.LayerTraining,
			Evaluation: &challenge.Evaluation{MaxMoves: 8},
		},
		{
			ID:     "find_ca_quote_combo",
			Name:   "Find + Change Around",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Use f{char} then ca\" to change around the quotes.",
				"",
				"\"remove\"",
			},
			Result: []string{
				"Use f{char} then ca\" to change around the quotes.",
				"",
				"hello",
			},
			Layer:      challenge.LayerTraining,
			Evaluation: &challenge.Evaluation{MaxMoves: 10},
		},
		{
			ID:     "find_ciw_combo",
			Name:   "Find + Change Word",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Use f{char} then ciw to change the word.",
				"",
				"keep delete keep",
			},
			Result: []string{
				"Use f{char} then ciw to change the word.",
				"",
				"keep hello keep",
			},
			Layer:      challenge.LayerTraining,
			Evaluation: &challenge.Evaluation{MaxMoves: 10},
		},
		{
			ID:     "dw_dot_combo",
			Name:   "Delete + Repeat",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Use dw then . to delete the remaining words.",
				"",
				"one two three four",
			},
			Result: []string{
				"Use dw then . to delete the remaining words.",
				"",
				"one four",
			},
			Layer:      challenge.LayerTraining,
			Evaluation: &challenge.Evaluation{MaxMoves: 8},
		},
		{
			ID:     "ciw_dot_combo",
			Name:   "Change + Repeat",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Use ciw then . to change the remaining words.",
				"",
				"one two three",
			},
			Result: []string{
				"Use ciw then . to change the remaining words.",
				"",
				"hello hello three",
			},
			Layer:      challenge.LayerTraining,
			Evaluation: &challenge.Evaluation{MaxMoves: 12},
		},
		{
			ID:     "yank_paste_combo",
			Name:   "Yank + Paste",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Use yy then p to duplicate the line below.",
				"",
				"alpha",
				"beta",
			},
			Result: []string{
				"Use yy then p to duplicate the line below.",
				"",
				"alpha",
				"beta",
				"alpha",
			},
			Layer:      challenge.LayerTraining,
			Evaluation: &challenge.Evaluation{MaxMoves: 8},
		},
		{
			ID:     "dd_paste_combo",
			Name:   "Cut + Paste",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Use dd then p to move the line to the end.",
				"",
				"one",
				"TWO",
				"three",
			},
			Result: []string{
				"Use dd then p to move the line to the end.",
				"",
				"one",
				"three",
				"TWO",
			},
			Layer:      challenge.LayerTraining,
			Evaluation: &challenge.Evaluation{MaxMoves: 8},
		},
		{
			ID:     "dd_paste_before_combo",
			Name:   "Cut + Paste Before",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Use dd then P to swap the two lines.",
				"",
				"first",
				"second",
			},
			Result: []string{
				"Use dd then P to swap the two lines.",
				"",
				"second",
				"first",
			},
			Layer:      challenge.LayerTraining,
			Evaluation: &challenge.Evaluation{MaxMoves: 8},
		},
		// Trial
		{
			ID:     "trial_find_delete",
			Name:   "Find + Delete Trial",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Remove the third word.",
				"",
				"apple banana cherry date",
			},
			Result: []string{
				"Remove the third word.",
				"",
				"apple banana date",
			},
			Layer:      challenge.LayerTrial,
			Evaluation: &challenge.Evaluation{MaxMoves: 6},
		},
		{
			ID:     "trial_find_change",
			Name:   "Find + Change Trial",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Replace the value with DONE.",
				"",
				`status = "pending"`,
			},
			Result: []string{
				"Replace the value with DONE.",
				"",
				`status = "DONE"`,
			},
			Layer:      challenge.LayerTrial,
			Evaluation: &challenge.Evaluation{MaxMoves: 8},
		},
		{
			ID:     "trial_dot_repeat",
			Name:   "Dot Repeat Trial",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Remove every occurrence of the.",
				"",
				"the fox and the dog and the bird",
			},
			Result: []string{
				"Remove every occurrence of the.",
				"",
				"fox and dog and bird",
			},
			Layer:      challenge.LayerTrial,
			Evaluation: &challenge.Evaluation{MaxMoves: 12},
		},
		{
			ID:     "trial_delete_choice",
			Name:   "Delete Choice Trial",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Remove the bracketed content.",
				"",
				"keep [remove] keep",
			},
			Result: []string{
				"Remove the bracketed content.",
				"",
				"keep  keep",
			},
			Layer:      challenge.LayerTrial,
			Evaluation: &challenge.Evaluation{MaxMoves: 6},
		},
		{
			ID:     "trial_repeat_choice",
			Name:   "Repeat Choice Trial",
			Verify: challenge.VerifyComposite,
			Content: []string{
				"Normalize all values to done.",
				"",
				"item = alpha",
				"item = beta",
				"item = gamma",
			},
			Result: []string{
				"Normalize all values to done.",
				"",
				"item = done",
				"item = done",
				"item = done",
			},
			Layer:      challenge.LayerTrial,
			Evaluation: &challenge.Evaluation{MaxMoves: 12},
		},
	}
}

// allIDs and byID are derived once from the compiled-in catalog so lookups do
// not rebuild the 48 challenge literals on every call.
var (
	allIDs = func() []string {
		all := All()
		ids := make([]string, 0, len(all))
		for _, c := range all {
			ids = append(ids, c.ID)
		}
		return ids
	}()
	byID = func() map[string]challenge.Challenge {
		m := make(map[string]challenge.Challenge, len(allIDs))
		for _, c := range All() {
			m[c.ID] = c
		}
		return m
	}()
)

func IDs() []string {
	return append([]string(nil), allIDs...)
}

func Exists(id string) bool {
	_, ok := byID[id]
	return ok
}
