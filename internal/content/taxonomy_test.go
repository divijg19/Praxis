package content

import "testing"

type CurriculumMetadata struct {
	Concept string
	Context string
	Stage   string
}

const (
	StageMovement             = "Movement"
	StageSearch               = "Search"
	StageStructuralNavigation = "Structural Navigation"
	StageEditing              = "Editing"
	StageTextObjects          = "Text Objects"
	StageRegisters            = "Registers"
)

var curriculum = map[string]CurriculumMetadata{
	// Movement
	"motion_rush":              {"hjkl", "basic navigation", StageMovement},
	"grid_rush":                {"hjkl", "grid navigation", StageMovement},
	"utf8_cursor_hunter":       {"f", "multi-byte search", StageMovement},
	// Search
	"find_hunter":              {"f", "character search", StageSearch},
	"word_hunter":              {"w", "word motion", StageSearch},
	"symbol_hunter":            {"F", "backward symbol search", StageSearch},
	"line_hunter":              {"j", "line navigation", StageSearch},
	"slash_hunter":             {"/", "forward search", StageSearch},
	"question_hunter":          {"?", "backward search", StageSearch},
	"repeat_hunter":            {";", "repeat motion", StageSearch},
	// Structural Navigation
	"paren_hunter":             {"%", "matching delimiters navigation", StageStructuralNavigation},
	"sentence_hunter":          {")", "sentence navigation", StageStructuralNavigation},
	"inner_paren_hunter":       {"i(", "select inside parentheses", StageStructuralNavigation},
	"around_paren_hunter":      {"a(", "select around parentheses", StageStructuralNavigation},
	"inner_bracket_hunter":     {"i[", "select inside brackets", StageStructuralNavigation},
	"around_bracket_hunter":    {"a[", "select around brackets", StageStructuralNavigation},
	"inner_quote_hunter":       {"i\"", "select inside quotes", StageStructuralNavigation},
	"around_quote_hunter":      {"a\"", "select around quotes", StageStructuralNavigation},
	"paragraph_hunter":         {"{", "paragraph navigation", StageStructuralNavigation},
	"match_hunter":             {"%", "nested delimiter matching", StageStructuralNavigation},
	// Editing
	"delete_character_hunter":  {"x", "delete character", StageEditing},
	"replace_character_hunter": {"r", "replace character", StageEditing},
	"toggle_case_hunter":       {"~", "toggle case", StageEditing},
	"delete_word_hunter":       {"dw", "delete word", StageEditing},
	"change_word_hunter":       {"ciw", "simple word replacement", StageEditing},
	"delete_line_hunter":       {"dd", "delete line", StageEditing},
	"delete_to_end_hunter":     {"D", "delete to end of line", StageEditing},
	// Text Objects
	"delete_inner_word_hunter":  {"diw", "delete inside word", StageTextObjects},
	"delete_around_word_hunter": {"daw", "delete around word", StageTextObjects},
	"delete_inner_paren_hunter":  {"di(", "delete inside parentheses", StageTextObjects},
	"delete_around_paren_hunter": {"da(", "delete around parentheses", StageTextObjects},
	"delete_inner_quote_hunter":  {"di\"", "delete inside quotes", StageTextObjects},
	"delete_around_quote_hunter": {"da\"", "delete around quotes", StageTextObjects},
	"change_inner_word_hunter":   {"ciw", "word replacement within structural editing", StageTextObjects},
	"change_inner_paren_hunter":  {"ci(", "change inside parentheses", StageTextObjects},
	"change_inner_quote_hunter":  {"ci\"", "change inside quotes", StageTextObjects},
	// Registers
	"yank_line_hunter":            {"yy", "yank line", StageRegisters},
	"named_register_hunter":       {"\"a", "named register", StageRegisters},
	"word_register_hunter":        {"\"A", "append to register", StageRegisters},
	"register_replace_hunter":     {"\"ap", "replace content from named register", StageRegisters},
	"register_duplicate_hunter":   {"\"ap", "duplicate content from named register", StageRegisters},
}

func TestCurriculumContextsComplete(t *testing.T) {
	validStages := map[string]bool{
		StageMovement:             true,
		StageSearch:               true,
		StageStructuralNavigation: true,
		StageEditing:              true,
		StageTextObjects:          true,
		StageRegisters:            true,
	}

	for _, c := range All() {
		meta, ok := curriculum[c.ID]
		if !ok {
			t.Errorf("challenge %q missing from curriculum metadata", c.ID)
			continue
		}
		if meta.Concept == "" {
			t.Errorf("challenge %q has empty Concept", c.ID)
		}
		if meta.Context == "" {
			t.Errorf("challenge %q has empty Context", c.ID)
		}
		if !validStages[meta.Stage] {
			t.Errorf("challenge %q has invalid Stage %q", c.ID, meta.Stage)
		}
	}
}

func TestConceptContextPairsUnique(t *testing.T) {
	seen := make(map[string]string)
	for _, c := range All() {
		meta := curriculum[c.ID]
		key := meta.Concept + "|" + meta.Context
		if first, ok := seen[key]; ok {
			t.Errorf("duplicate (Concept, Context) pair between %q and %q: (%q, %q)",
				first, c.ID, meta.Concept, meta.Context)
		}
		seen[key] = c.ID
	}
}

func TestProgressionCoverage(t *testing.T) {
	stages := make(map[string]bool)
	for _, c := range All() {
		meta := curriculum[c.ID]
		stages[meta.Stage] = true
	}

	expected := []string{StageMovement, StageSearch, StageStructuralNavigation,
		StageEditing, StageTextObjects, StageRegisters}
	for _, s := range expected {
		if !stages[s] {
			t.Errorf("progression stage %q has no challenges", s)
		}
	}
}

func TestStageIntroductionOrder(t *testing.T) {
	stageOrder := []string{
		StageMovement,
		StageSearch,
		StageStructuralNavigation,
		StageEditing,
		StageTextObjects,
		StageRegisters,
	}

	first := make(map[string]int)
	for i, c := range All() {
		meta := curriculum[c.ID]
		if _, ok := first[meta.Stage]; !ok {
			first[meta.Stage] = i
		}
	}

	for j := 1; j < len(stageOrder); j++ {
		prev := stageOrder[j-1]
		curr := stageOrder[j]
		if first[curr] < first[prev] {
			t.Errorf("stage %q (first at index %d) appears before %q (first at index %d)",
				curr, first[curr], prev, first[prev])
		}
	}
}
