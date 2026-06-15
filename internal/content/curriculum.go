package content

type Metadata struct {
	Concept string
	Context string
	Stage   string
	Layer   string
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
	"motion_rush":        {"hjkl", "basic navigation", StageMovement, "Tutorial"},
	"grid_rush":          {"hjkl", "grid navigation", StageMovement, "Tutorial"},
	"utf8_cursor_hunter": {"f", "multi-byte search", StageMovement, "Tutorial"},
	// Search
	"find_hunter":     {"f", "character search", StageSearch, "Tutorial"},
	"word_hunter":     {"w", "word motion", StageSearch, "Tutorial"},
	"symbol_hunter":   {"F", "backward symbol search", StageSearch, "Tutorial"},
	"line_hunter":     {"j", "line navigation", StageSearch, "Tutorial"},
	"slash_hunter":    {"/", "forward search", StageSearch, "Tutorial"},
	"question_hunter": {"?", "backward search", StageSearch, "Tutorial"},
	"repeat_hunter":   {";", "repeat motion", StageSearch, "Tutorial"},
	// Structural Navigation
	"paren_hunter":          {"%", "matching delimiters navigation", StageStructuralNavigation, "Tutorial"},
	"sentence_hunter":       {")", "sentence navigation", StageStructuralNavigation, "Tutorial"},
	"inner_paren_hunter":    {"i(", "select inside parentheses", StageStructuralNavigation, "Tutorial"},
	"around_paren_hunter":   {"a(", "select around parentheses", StageStructuralNavigation, "Tutorial"},
	"inner_bracket_hunter":  {"i[", "select inside brackets", StageStructuralNavigation, "Tutorial"},
	"around_bracket_hunter": {"a[", "select around brackets", StageStructuralNavigation, "Tutorial"},
	"inner_quote_hunter":    {"i\"", "select inside quotes", StageStructuralNavigation, "Tutorial"},
	"around_quote_hunter":   {"a\"", "select around quotes", StageStructuralNavigation, "Tutorial"},
	"paragraph_hunter":      {"{", "paragraph navigation", StageStructuralNavigation, "Tutorial"},
	"match_hunter":          {"%", "nested delimiter matching", StageStructuralNavigation, "Tutorial"},
	// Editing
	"delete_character_hunter":  {"x", "delete character", StageEditing, "Tutorial"},
	"replace_character_hunter": {"r", "replace character", StageEditing, "Tutorial"},
	"toggle_case_hunter":       {"~", "toggle case", StageEditing, "Tutorial"},
	"delete_word_hunter":       {"dw", "delete word", StageEditing, "Tutorial"},
	"change_word_hunter":       {"ciw", "simple word replacement", StageEditing, "Tutorial"},
	"delete_line_hunter":       {"dd", "delete line", StageEditing, "Tutorial"},
	"delete_to_end_hunter":     {"D", "delete to end of line", StageEditing, "Tutorial"},
	// Text Objects
	"delete_inner_word_hunter":   {"diw", "delete inside word", StageTextObjects, "Tutorial"},
	"delete_around_word_hunter":  {"daw", "delete around word", StageTextObjects, "Tutorial"},
	"delete_inner_paren_hunter":  {"di(", "delete inside parentheses", StageTextObjects, "Tutorial"},
	"delete_around_paren_hunter": {"da(", "delete around parentheses", StageTextObjects, "Tutorial"},
	"delete_inner_quote_hunter":  {"di\"", "delete inside quotes", StageTextObjects, "Tutorial"},
	"delete_around_quote_hunter": {"da\"", "delete around quotes", StageTextObjects, "Tutorial"},
	"change_inner_word_hunter":   {"ciw", "word replacement within structural editing", StageTextObjects, "Tutorial"},
	"change_inner_paren_hunter":  {"ci(", "change inside parentheses", StageTextObjects, "Tutorial"},
	"change_inner_quote_hunter":  {"ci\"", "change inside quotes", StageTextObjects, "Tutorial"},
	// Registers
	"yank_line_hunter":          {"yy", "yank line", StageRegisters, "Tutorial"},
	"named_register_hunter":     {"\"a", "named register", StageRegisters, "Tutorial"},
	"word_register_hunter":      {"\"A", "append to register", StageRegisters, "Tutorial"},
	"register_replace_hunter":   {"\"ap", "replace content from named register", StageRegisters, "Tutorial"},
	"register_duplicate_hunter": {"\"ap", "duplicate content from named register", StageRegisters, "Tutorial"},
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
