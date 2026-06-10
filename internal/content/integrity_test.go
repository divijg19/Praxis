package content

import (
	"strings"
	"testing"
)

var curriculumConcepts = map[string]string{
	"motion_rush":               "hjkl",
	"grid_rush":                 "hjkl",
	"find_hunter":               "f",
	"word_hunter":               "w",
	"symbol_hunter":             "F",
	"line_hunter":               "j",
	"paren_hunter":              "%",
	"sentence_hunter":           ")",
	"slash_hunter":              "/",
	"question_hunter":           "?",
	"repeat_hunter":             ";",
	"inner_paren_hunter":        "i(",
	"around_paren_hunter":       "a(",
	"inner_bracket_hunter":      "i[",
	"around_bracket_hunter":     "a[",
	"inner_quote_hunter":        "i\"",
	"around_quote_hunter":       "a\"",
	"paragraph_hunter":          "{",
	"match_hunter":              "%",
	"delete_character_hunter":   "x",
	"replace_character_hunter":  "r",
	"toggle_case_hunter":        "~",
	"delete_word_hunter":        "dw",
	"change_word_hunter":        "ciw",
	"utf8_cursor_hunter":        "f",
	"delete_line_hunter":        "dd",
	"delete_to_end_hunter":      "D",
	"delete_inner_word_hunter":  "diw",
	"delete_around_word_hunter": "daw",
	"delete_inner_paren_hunter":  "di(",
	"delete_around_paren_hunter": "da(",
	"delete_inner_quote_hunter":  "di\"",
	"delete_around_quote_hunter": "da\"",
	"change_inner_word_hunter":   "ciw",
	"change_inner_paren_hunter":  "ci(",
	"change_inner_quote_hunter":  "ci\"",
	"yank_line_hunter":           "yy",
	"named_register_hunter":      "\"a",
	"word_register_hunter":       "\"A",
	"register_replace_hunter":    "\"ap",
	"register_duplicate_hunter":  "\"ap",
}

func TestCoreConceptCoverage(t *testing.T) {
	seen := make(map[string]bool)
	for _, c := range All() {
		concept, ok := curriculumConcepts[c.ID]
		if !ok {
			t.Errorf("challenge %q missing from curriculum map", c.ID)
			continue
		}
		seen[concept] = true
	}

	families := []struct {
		name   string
		values []string
	}{
		{"fFtT", []string{"f", "F", "t", "T"}},
		{"wWeEeBb", []string{"w", "W", "e", "E", "b", "B"}},
		{"registers", []string{"\"a", "\"A", "\"ap"}},
	}
	for _, f := range families {
		var found bool
		for _, v := range f.values {
			if seen[v] {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("no challenge teaches a %q family concept (values: %v)", f.name, f.values)
		}
	}

	for _, concept := range []string{"dd", "dw", "diw", "daw", "ciw", "yy"} {
		if !seen[concept] {
			t.Errorf("core concept %q is not taught by any challenge", concept)
		}
	}
}

func TestNoDuplicateChallengeContent(t *testing.T) {
	t.Run("content", func(t *testing.T) {
		seen := make(map[string]string)
		for _, c := range All() {
			key := strings.Join(c.Content, "\n")
			if first, ok := seen[key]; ok {
				t.Errorf("duplicate content between %q and %q", first, c.ID)
			}
			seen[key] = c.ID
		}
	})

	t.Run("buffer_results", func(t *testing.T) {
		seen := make(map[string]string)
		for _, c := range All() {
			if len(c.Result) == 0 {
				continue
			}
			key := strings.Join(c.Result, "\n")
			if first, ok := seen[key]; ok {
				t.Errorf("duplicate buffer result between %q and %q", first, c.ID)
			}
			seen[key] = c.ID
		}
	})
}

func TestCurriculumMapComplete(t *testing.T) {
	if len(curriculumConcepts) != len(All()) {
		t.Errorf("curriculum map has %d entries, expected %d", len(curriculumConcepts), len(All()))
	}

	for _, c := range All() {
		concept, ok := curriculumConcepts[c.ID]
		if !ok {
			t.Errorf("challenge %q missing from curriculum map", c.ID)
			continue
		}
		if concept == "" {
			t.Errorf("challenge %q has empty curriculum concept", c.ID)
		}
	}

	for id := range curriculumConcepts {
		var found bool
		for _, c := range All() {
			if c.ID == id {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("curriculum map has unknown challenge %q", id)
		}
	}
}
