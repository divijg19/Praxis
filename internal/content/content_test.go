package content

import (
	"testing"

	"github.com/divijg19/Praxis/internal/validator"
)

var stableChallengeIDs = []string{
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
	got := make([]string, len(All()))
	for i, c := range All() {
		got[i] = c.ID
	}
	if len(got) != len(stableChallengeIDs) {
		t.Fatalf("got %d challenges, want %d", len(got), len(stableChallengeIDs))
	}
	for i, id := range got {
		if id != stableChallengeIDs[i] {
			t.Errorf("challenge[%d] = %q, want %q", i, id, stableChallengeIDs[i])
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

func TestChallengeCount(t *testing.T) {
	if got := len(All()); got != len(stableChallengeIDs) {
		t.Errorf("got %d challenges, want %d", got, len(stableChallengeIDs))
	}
}

func TestResultMatchesVerify(t *testing.T) {
	for _, c := range All() {
		switch c.Verify {
		case "buffer":
			if len(c.Result) == 0 {
				t.Errorf("buffer challenge %s has empty Result", c.ID)
			}
		case "cursor":
			if len(c.Result) > 0 {
				t.Errorf("cursor challenge %s has unexpected Result", c.ID)
			}
		}
	}
}

func TestBufferChallengeLayout(t *testing.T) {
	for _, c := range All() {
		if c.Verify != "buffer" {
			continue
		}
		if len(c.Content) < 3 {
			t.Errorf("buffer challenge %s has fewer than 3 content lines", c.ID)
			continue
		}
		if c.Content[1] != "" {
			t.Errorf("buffer challenge %s content[1] is not blank", c.ID)
		}
		if len(c.Result) > 0 && c.Result[0] != c.Content[0] {
			t.Errorf("buffer challenge %s result[0] != content[0]", c.ID)
		}
	}
}

func TestCursorChallengeLayout(t *testing.T) {
	for _, c := range All() {
		if c.Verify != "cursor" {
			continue
		}
		if len(c.Content) < 1 {
			t.Errorf("cursor challenge %s has empty Content", c.ID)
			continue
		}
		if c.Content[0] == "" {
			t.Errorf("cursor challenge %s has empty instruction line", c.ID)
		}
		if c.Target == "" {
			t.Errorf("cursor challenge %s has empty Target", c.ID)
		}
	}
}

func TestNoValidatorDrift(t *testing.T) {
	used := make(map[string]bool)
	for _, c := range All() {
		used[c.Verify] = true
	}
	for _, name := range []string{"cursor", "buffer"} {
		if !used[name] {
			t.Errorf("validator %q is registered but unused by any challenge", name)
		}
	}
}

func TestResultShapeMatchesVerify(t *testing.T) {
	for _, c := range All() {
		switch c.Verify {
		case "cursor":
			if c.Target == "" {
				t.Errorf("cursor challenge %s has empty Target", c.ID)
			}
			if len(c.Result) > 0 {
				t.Errorf("cursor challenge %s has unexpected Result", c.ID)
			}
		case "buffer":
			if c.Target != "" {
				t.Errorf("buffer challenge %s has non-empty Target: %q", c.ID, c.Target)
			}
			if len(c.Result) == 0 {
				t.Errorf("buffer challenge %s has empty Result", c.ID)
			}
		}
	}
}

func TestAllChallengesHaveLayer(t *testing.T) {
	for _, c := range All() {
		if c.Layer == "" {
			t.Errorf("challenge %s has empty Layer", c.ID)
		}
	}
}

func TestLayerValidValues(t *testing.T) {
	allowed := map[string]bool{
		"Tutorial": true,
		"Training": true,
		"Trial":    true,
		"Boss":     true,
	}
	for _, c := range All() {
		if !allowed[c.Layer] {
			t.Errorf("challenge %s has invalid Layer %q", c.ID, c.Layer)
		}
	}
}

func TestAllCurrentChallengesAreTutorials(t *testing.T) {
	for _, c := range All() {
		if c.Layer != "Tutorial" {
			t.Errorf("challenge %s has Layer %q, expected %q", c.ID, c.Layer, "Tutorial")
		}
	}
}

func TestContentResultLineCountReasonable(t *testing.T) {
	for _, c := range All() {
		if c.Verify != "buffer" {
			continue
		}
		n := len(c.Result)
		m := len(c.Content)
		if n < m-1 || n > m+1 {
			t.Errorf("buffer challenge %s: Result has %d lines, Content has %d (reasonable range: %d–%d)",
				c.ID, n, m, m-1, m+1)
		}
	}
}
