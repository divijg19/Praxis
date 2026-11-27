package content

import "unicode/utf8"

// Trace is a canonical, internal-only solution. It is used by the certification
// harness to prove a challenge is solvable with real Vim keys, and to derive a
// verified par. StartRow/StartCol are 1-based row and 0-based column in the
// challenge buffer. Keys are fed verbatim with nvim_feedkeys.
type Trace struct {
	StartRow int    `json:"start_row"`
	StartCol int    `json:"start_col"`
	Keys     string `json:"keys"`
}

// traces is deliberately separate from Challenge so solution sequences are not
// part of the learner-facing describe payload.
var traces = map[string]Trace{
	// Tutorial — cursor
	"grid_rush":           {StartRow: 2, StartCol: 0, Keys: "jllllll"},
	"find_hunter":         {StartRow: 3, StartCol: 0, Keys: "f★"},
	"word_hunter":         {StartRow: 3, StartCol: 0, Keys: "5w"},
	"line_hunter":         {StartRow: 3, StartCol: 0, Keys: "5j"},
	"paren_hunter":        {StartRow: 3, StartCol: 0, Keys: "%l"},
	"sentence_hunter":     {StartRow: 3, StartCol: 0, Keys: ")hh"},
	"slash_hunter":        {StartRow: 3, StartCol: 0, Keys: "/★\r"},
	"question_hunter":     {StartRow: 20, StartCol: 0, Keys: "?★\r"},
	"repeat_hunter":       {StartRow: 3, StartCol: 4, Keys: "f★;"},
	"inner_paren_hunter":  {StartRow: 3, StartCol: 6, Keys: "di("},
	"around_paren_hunter": {StartRow: 3, StartCol: 6, Keys: "da("},
	"inner_quote_hunter":  {StartRow: 3, StartCol: 1, Keys: "di\""},
	"around_quote_hunter": {StartRow: 3, StartCol: 1, Keys: "da\""},
	"paragraph_hunter":    {StartRow: 3, StartCol: 0, Keys: "3}0"},
	"utf8_cursor_hunter":  {StartRow: 3, StartCol: 0, Keys: "6l"},

	// Tutorial — buffer
	"delete_character_hunter":    {StartRow: 3, StartCol: 0, Keys: "3lx"},
	"replace_character_hunter":   {StartRow: 3, StartCol: 0, Keys: "lre"},
	"toggle_case_hunter":         {StartRow: 3, StartCol: 0, Keys: "5~"},
	"delete_word_hunter":         {StartRow: 3, StartCol: 0, Keys: "6ldw"},
	"change_word_hunter":         {StartRow: 3, StartCol: 0, Keys: "ciwbar\x1b"},
	"delete_line_hunter":         {StartRow: 3, StartCol: 0, Keys: "jdd"},
	"delete_to_end_hunter":       {StartRow: 3, StartCol: 0, Keys: "14ld$"},
	"delete_inner_word_hunter":   {StartRow: 3, StartCol: 0, Keys: "5ldiw"},
	"delete_around_word_hunter":  {StartRow: 3, StartCol: 0, Keys: "5ldaw"},
	"delete_inner_paren_hunter":  {StartRow: 3, StartCol: 0, Keys: "6ldi("},
	"delete_around_paren_hunter": {StartRow: 3, StartCol: 0, Keys: "6lda("},
	"delete_inner_quote_hunter":  {StartRow: 3, StartCol: 0, Keys: "2ldi\""},
	"delete_around_quote_hunter": {StartRow: 3, StartCol: 0, Keys: "2lda\""},
	"change_inner_paren_hunter":  {StartRow: 3, StartCol: 0, Keys: "6lci(hello\x1b"},
	"change_inner_quote_hunter":  {StartRow: 3, StartCol: 0, Keys: "2lci\"hello\x1b"},
	"yank_line_hunter":           {StartRow: 3, StartCol: 0, Keys: "yyp"},
	"word_register_hunter":       {StartRow: 3, StartCol: 0, Keys: "\"aywA \x1b\"apx"},
	"register_replace_hunter":    {StartRow: 3, StartCol: 0, Keys: "\"ayyjdd\"ap"},

	// Training
	"find_diw_combo":        {StartRow: 3, StartCol: 0, Keys: "fddiw"},
	"find_daw_combo":        {StartRow: 3, StartCol: 0, Keys: "fddaw"},
	"find_di_paren_combo":   {StartRow: 3, StartCol: 0, Keys: "f(di("},
	"find_ca_quote_combo":   {StartRow: 3, StartCol: 0, Keys: "ca\"hello\x1b"},
	"find_ciw_combo":        {StartRow: 3, StartCol: 0, Keys: "fdciwhello\x1b"},
	"dw_dot_combo":          {StartRow: 3, StartCol: 4, Keys: "dw."},
	"ciw_dot_combo":         {StartRow: 3, StartCol: 0, Keys: "ciwhello\x1bw."},
	"yank_paste_combo":      {StartRow: 3, StartCol: 0, Keys: "yyjp"},
	"dd_paste_combo":        {StartRow: 3, StartCol: 0, Keys: "jddp"},
	"dd_paste_before_combo": {StartRow: 3, StartCol: 0, Keys: "jddP"},

	// Trial
	"trial_find_delete":   {StartRow: 3, StartCol: 0, Keys: "2wdw"},
	"trial_find_change":   {StartRow: 3, StartCol: 0, Keys: "ci\"DONE\x1b"},
	"trial_dot_repeat":    {StartRow: 3, StartCol: 0, Keys: "dw.."},
	"trial_delete_choice": {StartRow: 3, StartCol: 0, Keys: "5ldaw"},
	"trial_repeat_choice": {StartRow: 3, StartCol: 0, Keys: "f=wciwdone\x1bj0f=w.j0f=w."},
}

// TraceFor returns the canonical solution for id, if one is registered.
func TraceFor(id string) (Trace, bool) {
	t, ok := traces[id]
	return t, ok
}

// ParFor returns the verified physical-key count of the canonical solution.
func ParFor(id string) int {
	t, ok := TraceFor(id)
	if !ok {
		return 0
	}
	return utf8.RuneCountInString(t.Keys)
}
