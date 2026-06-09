package content

import (
	"testing"

	"github.com/divijg19/Praxis/internal/validator"
)

func TestUniqueChallengeIDs(t *testing.T) {
	seen := make(map[string]bool)
	for _, c := range All() {
		if seen[c.ID] {
			t.Errorf("duplicate ID: %s", c.ID)
		}
		seen[c.ID] = true
	}
}

func TestUniqueChallengeNames(t *testing.T) {
	seen := make(map[string]bool)
	for _, c := range All() {
		if seen[c.Name] {
			t.Errorf("duplicate name: %s", c.Name)
		}
		seen[c.Name] = true
	}
}

func TestAllChallengesHaveVerify(t *testing.T) {
	for _, c := range All() {
		if c.Verify == "" {
			t.Errorf("challenge %s has empty Verify", c.ID)
		}
	}
}

func TestBufferChallengesHaveResults(t *testing.T) {
	for _, c := range All() {
		if c.Verify == "buffer" && len(c.Result) == 0 {
			t.Errorf("buffer challenge %s has empty Result", c.ID)
		}
	}
}

func TestCursorChallengesHaveTargets(t *testing.T) {
	for _, c := range All() {
		if c.Verify == "cursor" && c.Target == "" {
			t.Errorf("cursor challenge %s has empty Target", c.ID)
		}
	}
}

func TestBufferChallengesHaveNoTargets(t *testing.T) {
	for _, c := range All() {
		if c.Verify == "buffer" && c.Target != "" {
			t.Errorf("buffer challenge %s has non-empty Target: %q", c.ID, c.Target)
		}
	}
}

func TestNoEmptyContent(t *testing.T) {
	for _, c := range All() {
		if len(c.Content) == 0 {
			t.Errorf("challenge %s has empty Content", c.ID)
		}
	}
}

func TestInstructionLinePresent(t *testing.T) {
	for _, c := range All() {
		if len(c.Content) == 0 || c.Content[0] == "" {
			t.Errorf("challenge %s has empty or missing instruction line", c.ID)
		}
	}
}

func TestValidatorCoverage(t *testing.T) {
	for _, c := range All() {
		if !validator.Exists(c.Verify) {
			t.Errorf("challenge %s uses unknown validator: %s", c.ID, c.Verify)
		}
	}
}

func TestChallengeIDsStable(t *testing.T) {
	expected := []string{
		"motion_rush",
		"grid_rush",
		"find_hunter",
		"word_hunter",
		"symbol_hunter",
		"line_hunter",
		"paren_hunter",
		"sentence_hunter",
		"slash_hunter",
		"question_hunter",
		"repeat_hunter",
		"inner_paren_hunter",
		"around_paren_hunter",
		"inner_bracket_hunter",
		"around_bracket_hunter",
		"inner_quote_hunter",
		"around_quote_hunter",
		"paragraph_hunter",
		"match_hunter",
		"delete_character_hunter",
		"replace_character_hunter",
		"toggle_case_hunter",
		"delete_word_hunter",
		"change_word_hunter",
		"utf8_cursor_hunter",
		"delete_line_hunter",
		"delete_to_end_hunter",
		"delete_inner_word_hunter",
		"delete_around_word_hunter",
		"delete_inner_paren_hunter",
		"delete_around_paren_hunter",
		"delete_inner_quote_hunter",
		"delete_around_quote_hunter",
		"change_inner_word_hunter",
		"change_inner_paren_hunter",
		"change_inner_quote_hunter",
		"yank_line_hunter",
		"named_register_hunter",
		"word_register_hunter",
		"register_replace_hunter",
		"register_duplicate_hunter",
	}
	got := make([]string, len(All()))
	for i, c := range All() {
		got[i] = c.ID
	}
	if len(got) != len(expected) {
		t.Fatalf("got %d challenges, want %d", len(got), len(expected))
	}
	for i, id := range got {
		if id != expected[i] {
			t.Errorf("challenge[%d] = %q, want %q", i, id, expected[i])
		}
	}
}

func TestChallengeNamesStable(t *testing.T) {
	expected := []string{
		"Motion Rush",
		"Grid Rush",
		"Find Hunter",
		"Word Hunter",
		"Symbol Hunter",
		"Line Hunter",
		"Paren Hunter",
		"Sentence Hunter",
		"Slash Hunter",
		"Question Hunter",
		"Repeat Hunter",
		"Inner Paren Hunter",
		"Around Paren Hunter",
		"Inner Bracket Hunter",
		"Around Bracket Hunter",
		"Inner Quote Hunter",
		"Around Quote Hunter",
		"Paragraph Hunter",
		"Match Hunter",
		"Delete Character Hunter",
		"Replace Character Hunter",
		"Toggle Case Hunter",
		"Delete Word Hunter",
		"Change Word Hunter",
		"UTF-8 Cursor Hunter",
		"Delete Line Hunter",
		"Delete To End Hunter",
		"Delete Inner Word Hunter",
		"Delete Around Word Hunter",
		"Delete Inner Paren Hunter",
		"Delete Around Paren Hunter",
		"Delete Inner Quote Hunter",
		"Delete Around Quote Hunter",
		"Change Inner Word Hunter",
		"Change Inner Paren Hunter",
		"Change Inner Quote Hunter",
		"Unnamed Register Hunter",
		"Named Register Hunter",
		"Word Register Hunter",
		"Register Replace Hunter",
		"Register Duplicate Hunter",
	}
	got := make([]string, len(All()))
	for i, c := range All() {
		got[i] = c.Name
	}
	if len(got) != len(expected) {
		t.Fatalf("got %d challenges, want %d", len(got), len(expected))
	}
	for i, name := range got {
		if name != expected[i] {
			t.Errorf("challenge[%d] name = %q, want %q", i, name, expected[i])
		}
	}
}
