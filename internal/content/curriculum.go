package content

type Metadata struct {
	Concept     string
	Context     string
	Stage       string
	Layer       string
	DerivedFrom []string
}

const (
	StageMovement             = "Movement"
	StageSearch               = "Search"
	StageStructuralNavigation = "Structural Navigation"
	StageEditing              = "Editing"
	StageTextObjects          = "Text Objects"
	StageRegisters            = "Registers"
)

var curriculum = map[string]Metadata{
	// Movement
	"motion_rush":        {"hjkl", "basic navigation", StageMovement, "Tutorial", nil},
	"grid_rush":          {"hjkl", "grid navigation", StageMovement, "Tutorial", nil},
	"utf8_cursor_hunter": {"utf8", "multibyte navigation", StageMovement, "Tutorial", nil},
	// Search
	"find_hunter":     {"f", "character search", StageSearch, "Tutorial", nil},
	"word_hunter":     {"w", "word motion", StageSearch, "Tutorial", nil},
	"symbol_hunter":   {"F", "backward symbol search", StageSearch, "Tutorial", nil},
	"line_hunter":     {"j", "line navigation", StageSearch, "Tutorial", nil},
	"slash_hunter":    {"/", "forward search", StageSearch, "Tutorial", nil},
	"question_hunter": {"?", "backward search", StageSearch, "Tutorial", nil},
	"repeat_hunter":   {";", "repeat motion", StageSearch, "Tutorial", nil},
	// Structural Navigation
	"paren_hunter":          {"%", "matching delimiters navigation", StageStructuralNavigation, "Tutorial", nil},
	"sentence_hunter":       {")", "sentence navigation", StageStructuralNavigation, "Tutorial", nil},
	"inner_paren_hunter":    {"i(", "select inside parentheses", StageStructuralNavigation, "Tutorial", nil},
	"around_paren_hunter":   {"a(", "select around parentheses", StageStructuralNavigation, "Tutorial", nil},
	"inner_bracket_hunter":  {"i[", "select inside brackets", StageStructuralNavigation, "Tutorial", nil},
	"around_bracket_hunter": {"a[", "select around brackets", StageStructuralNavigation, "Tutorial", nil},
	"inner_quote_hunter":    {"i\"", "select inside quotes", StageStructuralNavigation, "Tutorial", nil},
	"around_quote_hunter":   {"a\"", "select around quotes", StageStructuralNavigation, "Tutorial", nil},
	"paragraph_hunter":      {"{", "paragraph navigation", StageStructuralNavigation, "Tutorial", nil},
	"match_hunter":          {"%", "nested delimiter matching", StageStructuralNavigation, "Tutorial", nil},
	// Editing
	"delete_character_hunter":  {"x", "delete character", StageEditing, "Tutorial", nil},
	"replace_character_hunter": {"r", "replace character", StageEditing, "Tutorial", nil},
	"toggle_case_hunter":       {"~", "toggle case", StageEditing, "Tutorial", nil},
	"delete_word_hunter":       {"dw", "delete word", StageEditing, "Tutorial", nil},
	"change_word_hunter":       {"ciw", "simple word replacement", StageEditing, "Tutorial", nil},
	"delete_line_hunter":       {"dd", "delete line", StageEditing, "Tutorial", nil},
	"delete_to_end_hunter":     {"D", "delete to end of line", StageEditing, "Tutorial", nil},
	// Text Objects
	"delete_inner_word_hunter":   {"diw", "delete inside word", StageTextObjects, "Tutorial", nil},
	"delete_around_word_hunter":  {"daw", "delete around word", StageTextObjects, "Tutorial", nil},
	"delete_inner_paren_hunter":  {"di(", "delete inside parentheses", StageTextObjects, "Tutorial", nil},
	"delete_around_paren_hunter": {"da(", "delete around parentheses", StageTextObjects, "Tutorial", nil},
	"delete_inner_quote_hunter":  {"di\"", "delete inside quotes", StageTextObjects, "Tutorial", nil},
	"delete_around_quote_hunter": {"da\"", "delete around quotes", StageTextObjects, "Tutorial", nil},
	"change_inner_word_hunter":   {"ciw", "word replacement within structural editing", StageTextObjects, "Tutorial", nil},
	"change_inner_paren_hunter":  {"ci(", "change inside parentheses", StageTextObjects, "Tutorial", nil},
	"change_inner_quote_hunter":  {"ci\"", "change inside quotes", StageTextObjects, "Tutorial", nil},
	// Registers
	"yank_line_hunter":          {"yy", "yank line", StageRegisters, "Tutorial", nil},
	"named_register_hunter":     {"\"a", "named register", StageRegisters, "Tutorial", nil},
	"word_register_hunter":      {"\"A", "append to register", StageRegisters, "Tutorial", nil},
	"register_replace_hunter":   {"\"ap", "replace content from named register", StageRegisters, "Tutorial", nil},
	"register_duplicate_hunter": {"\"ap", "duplicate content from named register", StageRegisters, "Tutorial", nil},
	// Training
	"find_diw_combo":        {"f", "composite deletion", StageTextObjects, "Training", nil},
	"find_daw_combo":        {"daw", "composite deletion around", StageTextObjects, "Training", nil},
	"find_di_paren_combo":   {"di(", "composite delete inside parens", StageTextObjects, "Training", nil},
	"find_ca_quote_combo":   {"ca\"", "composite change around quotes", StageTextObjects, "Training", nil},
	"find_ciw_combo":        {"ciw", "composite change word", StageEditing, "Training", nil},
	"dw_dot_combo":          {"dw", "composite delete with repeat", StageEditing, "Training", nil},
	"ciw_dot_combo":         {"ciw", "composite change with repeat", StageEditing, "Training", nil},
	"yank_paste_combo":      {"yy", "composite yank and paste", StageRegisters, "Training", nil},
	"dd_paste_combo":        {"dd", "composite cut and paste", StageRegisters, "Training", nil},
	"dd_paste_before_combo": {"dd", "composite cut and paste before", StageRegisters, "Training", nil},
	// Trial
	"trial_find_delete":   {"f", "recognition: f+diw on third argument", StageTextObjects, "Trial", []string{"find_diw_combo"}},
	"trial_find_change":   {"f", "recognition: f+ca quote replacement", StageTextObjects, "Trial", []string{"find_ca_quote_combo"}},
	"trial_dot_repeat":    {"dw", "recognition: dot repeat of edit", StageEditing, "Trial", []string{"dw_dot_combo"}},
	"trial_delete_choice": {"diw/daw", "recognition: inner vs around word", StageTextObjects, "Trial", []string{"find_diw_combo", "find_daw_combo"}},
	"trial_repeat_choice": {"dw/ciw", "recognition: repeat vs re-execute", StageEditing, "Trial", []string{"dw_dot_combo", "ciw_dot_combo"}},
}

func MetadataFor(id string) (Metadata, bool) {
	m, ok := curriculum[id]
	return m, ok
}

func ValidStages() map[string]bool {
	return map[string]bool{
		StageMovement:             true,
		StageSearch:               true,
		StageStructuralNavigation: true,
		StageEditing:              true,
		StageTextObjects:          true,
		StageRegisters:            true,
	}
}
