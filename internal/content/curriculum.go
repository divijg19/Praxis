package content

type Metadata struct {
	Concept     string
	Context     string
	Stage       string
	DerivedFrom []string
}

const (
	stageMovement             = "Movement"
	stageSearch               = "Search"
	stageStructuralNavigation = "Structural Navigation"
	stageEditing              = "Editing"
	stageTextObjects          = "Text Objects"
	stageRegisters            = "Registers"
)

var curriculum = map[string]Metadata{
	// Movement
	"motion_rush":        {"hjkl", "basic navigation", stageMovement, nil},
	"grid_rush":          {"hjkl", "grid navigation", stageMovement, nil},
	"utf8_cursor_hunter": {"utf8", "multibyte navigation", stageMovement, nil},
	// Search
	"find_hunter":     {"f", "character search", stageSearch, nil},
	"word_hunter":     {"w", "word motion", stageSearch, nil},
	"line_hunter":     {"j", "line navigation", stageSearch, nil},
	"slash_hunter":    {"/", "forward search", stageSearch, nil},
	"question_hunter": {"?", "backward search", stageSearch, nil},
	"repeat_hunter":   {";", "repeat motion", stageSearch, nil},
	// Structural Navigation
	"paren_hunter":          {"%", "matching delimiters navigation", stageStructuralNavigation, nil},
	"sentence_hunter":       {")", "sentence navigation", stageStructuralNavigation, nil},
	"inner_paren_hunter":    {"i(", "select inside parentheses", stageStructuralNavigation, nil},
	"around_paren_hunter":   {"a(", "select around parentheses", stageStructuralNavigation, nil},
	"inner_bracket_hunter":  {"i[", "select inside brackets", stageStructuralNavigation, nil},
	"around_bracket_hunter": {"a[", "select around brackets", stageStructuralNavigation, nil},
	"inner_quote_hunter":    {"i\"", "select inside quotes", stageStructuralNavigation, nil},
	"around_quote_hunter":   {"a\"", "select around quotes", stageStructuralNavigation, nil},
	"paragraph_hunter":      {"{", "paragraph navigation", stageStructuralNavigation, nil},
	// Editing
	"delete_character_hunter":  {"x", "delete character", stageEditing, nil},
	"replace_character_hunter": {"r", "replace character", stageEditing, nil},
	"toggle_case_hunter":       {"~", "toggle case", stageEditing, nil},
	"delete_word_hunter":       {"dw", "delete word", stageEditing, nil},
	"change_word_hunter":       {"ciw", "simple word replacement", stageEditing, nil},
	"delete_line_hunter":       {"dd", "delete line", stageEditing, nil},
	"delete_to_end_hunter":     {"d$", "delete to end of line", stageEditing, nil},
	// Text Objects
	"delete_inner_word_hunter":   {"diw", "delete inside word", stageTextObjects, nil},
	"delete_around_word_hunter":  {"daw", "delete around word", stageTextObjects, nil},
	"delete_inner_paren_hunter":  {"di(", "delete inside parentheses", stageTextObjects, nil},
	"delete_around_paren_hunter": {"da(", "delete around parentheses", stageTextObjects, nil},
	"delete_inner_quote_hunter":  {"di\"", "delete inside quotes", stageTextObjects, nil},
	"delete_around_quote_hunter": {"da\"", "delete around quotes", stageTextObjects, nil},
	"change_inner_paren_hunter":  {"ci(", "change inside parentheses", stageTextObjects, nil},
	"change_inner_quote_hunter":  {"ci\"", "change inside quotes", stageTextObjects, nil},
	// Registers
	"yank_line_hunter":          {"yy", "yank line", stageRegisters, nil},
	"word_register_hunter":      {"\"a", "store word in register", stageRegisters, nil},
	"register_replace_hunter":   {"\"ap", "replace content from named register", stageRegisters, nil},
	"register_duplicate_hunter": {"\"ap", "duplicate content from named register", stageRegisters, nil},
	// Training
	"find_diw_combo":        {"diw", "composite deletion", stageTextObjects, nil},
	"find_daw_combo":        {"daw", "composite deletion around", stageTextObjects, nil},
	"find_di_paren_combo":   {"di(", "composite delete inside parens", stageTextObjects, nil},
	"find_ca_quote_combo":   {"ca\"", "composite change around quotes", stageTextObjects, nil},
	"find_ciw_combo":        {"ciw", "composite change word", stageEditing, nil},
	"dw_dot_combo":          {"dw", "composite delete with repeat", stageEditing, nil},
	"ciw_dot_combo":         {"ciw", "composite change with repeat", stageEditing, nil},
	"yank_paste_combo":      {"yy", "composite yank and paste", stageRegisters, nil},
	"dd_paste_combo":        {"dd", "composite cut and paste", stageRegisters, nil},
	"dd_paste_before_combo": {"dd", "composite cut and paste before", stageRegisters, nil},
	// Trial
	"trial_find_delete":   {"diw", "recognition: f+diw on third argument", stageTextObjects, []string{"find_diw_combo"}},
	"trial_find_change":   {"ca\"", "recognition: f+ca quote replacement", stageTextObjects, []string{"find_ca_quote_combo"}},
	"trial_dot_repeat":    {"dw", "recognition: dot repeat of edit", stageEditing, []string{"dw_dot_combo"}},
	"trial_delete_choice": {"diw/daw", "recognition: inner vs around word", stageTextObjects, []string{"find_diw_combo", "find_daw_combo"}},
	"trial_repeat_choice": {"dw/ciw", "recognition: repeat vs re-execute", stageEditing, []string{"dw_dot_combo", "ciw_dot_combo"}},
}

func metadataFor(id string) (Metadata, bool) {
	m, ok := curriculum[id]
	return m, ok
}
