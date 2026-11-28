package content

import (
	"testing"
)

var stableChallengeIDs = []string{
	"grid_rush",
	"find_hunter",
	"word_hunter",
	"line_hunter",
	"paren_hunter",
	"sentence_hunter",
	"slash_hunter",
	"question_hunter",
	"repeat_hunter",
	"inner_paren_hunter",
	"around_paren_hunter",
	"inner_quote_hunter",
	"around_quote_hunter",
	"paragraph_hunter",
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
	"change_inner_paren_hunter",
	"change_inner_quote_hunter",
	"yank_line_hunter",
	"word_register_hunter",
	"register_replace_hunter",
	"find_diw_combo",
	"find_daw_combo",
	"find_di_paren_combo",
	"find_ca_quote_combo",
	"find_ciw_combo",
	"dw_dot_combo",
	"ciw_dot_combo",
	"yank_paste_combo",
	"dd_paste_combo",
	"dd_paste_before_combo",
	"trial_find_delete",
	"trial_find_change",
	"trial_dot_repeat",
	"trial_delete_choice",
	"trial_repeat_choice",
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

func TestInstructionLinePresent(t *testing.T) {
	for _, c := range All() {
		if len(c.Content) == 0 || c.Content[0] == "" {
			t.Errorf("challenge %s has empty or missing instruction line", c.ID)
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

func TestChallengeNamesUnique(t *testing.T) {
	seen := make(map[string]bool)
	for _, c := range All() {
		if seen[c.Name] {
			t.Errorf("duplicate challenge name %q", c.Name)
		}
		seen[c.Name] = true
	}
}

func TestBufferChallengeLayout(t *testing.T) {
	for _, c := range All() {
		if c.Verify != "buffer" && c.Verify != "composite" {
			continue
		}
		if len(c.Content) < 3 {
			t.Errorf("buffer-like challenge %s has fewer than 3 content lines", c.ID)
			continue
		}
		if c.Content[1] != "" {
			t.Errorf("buffer-like challenge %s content[1] is not blank", c.ID)
		}
		if len(c.Result) > 0 && c.Result[0] != c.Content[0] {
			t.Errorf("buffer-like challenge %s result[0] != content[0]", c.ID)
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

func TestLayerFor(t *testing.T) {
	for _, c := range All() {
		if got := LayerFor(c.ID); got != c.Layer {
			t.Errorf("LayerFor(%s) = %q, want %q", c.ID, got, c.Layer)
		}
	}
	if got := LayerFor("does_not_exist"); got != "" {
		t.Errorf("LayerFor(unknown) = %q, want empty", got)
	}
}

func TestLayerValidValues(t *testing.T) {
	allowed := map[string]bool{
		"Tutorial": true,
		"Training": true,
		"Trial":    true,
	}
	for _, c := range All() {
		if !allowed[c.Layer] {
			t.Errorf("challenge %s has invalid Layer %q", c.ID, c.Layer)
		}
	}
}

func TestCurriculumLayerDistribution(t *testing.T) {
	counts := map[string]int{}
	for _, c := range All() {
		counts[c.Layer]++
	}
	if counts["Tutorial"] != 33 {
		t.Errorf("Tutorial = %d, want 33", counts["Tutorial"])
	}
	if counts["Training"] != 10 {
		t.Errorf("Training = %d, want 10", counts["Training"])
	}
	if counts["Trial"] != 5 {
		t.Errorf("Trial = %d, want 5", counts["Trial"])
	}
}

func TestContentResultLineCountReasonable(t *testing.T) {
	for _, c := range All() {
		if c.Verify != "buffer" && c.Verify != "composite" {
			continue
		}
		n := len(c.Result)
		m := len(c.Content)
		if n < m-1 || n > m+1 {
			t.Errorf("buffer-like challenge %s: Result has %d lines, Content has %d (reasonable range: %d to %d)",
				c.ID, n, m, m-1, m+1)
		}
	}
}
