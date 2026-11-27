package verify

import (
	"testing"

	"github.com/divijg19/Praxis/internal/content"
)

func TestCanonicalTracesVerifyInCleanNeovim(t *testing.T) {
	for _, id := range content.IDs() {
		tr, ok := content.TraceFor(id)
		if !ok {
			t.Fatalf("missing trace for %s", id)
		}
		var keys []string
		for _, r := range tr.Keys {
			keys = append(keys, string(r))
		}
		result, err := Run(id, Input{Keys: keys})
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if !result.Verified {
			t.Errorf("%s: canonical trace did not verify: %s; final=%v", id, result.Reason, result.FinalLines)
		}
	}
}

func TestTaintedInputIsRejected(t *testing.T) {
	result, err := Run("grid_rush", Input{Keys: []string{"j"}, Tainted: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.Verified {
		t.Fatal("tainted input verified")
	}
}

func TestRejectsExCommandKeys(t *testing.T) {
	result, err := Run("grid_rush", Input{Keys: []string{":"}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Verified || result.Reason != "invalid key in log" {
		t.Fatalf("ex-command key was not rejected: %+v", result)
	}
}

func TestRejectsInstructionLineGoal(t *testing.T) {
	// The instruction line contains a ★; the product frontend never accepts it.
	result, err := Run("grid_rush", Input{Keys: []string{"k", "f", "★"}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Verified {
		t.Fatal("instruction-line target was accepted")
	}
}

func TestReportsMatchedPrefixLength(t *testing.T) {
	tr, _ := content.TraceFor("grid_rush")
	var keys []string
	for _, r := range tr.Keys {
		keys = append(keys, string(r))
	}
	keys = append(keys, "j") // extra key after the goal
	result, err := Run("grid_rush", Input{Keys: keys})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Verified || result.Keys != len(tr.Keys) {
		t.Fatalf("matched keys = %d, want %d: %+v", result.Keys, len(tr.Keys), result)
	}
}

func TestRejectsAssembledTermcode(t *testing.T) {
	// "<Char-58>" is ":" when termcodes are expanded; assembled from single
	// printable keys it must still be rejected.
	keys := []string{"<", "C", "h", "a", "r", "-", "5", "8", ">", "q", "a", "!", "<CR>"}
	result, err := Run("grid_rush", Input{Keys: keys})
	if err != nil {
		t.Fatal(err)
	}
	if result.Verified || result.Reason != "invalid key in log" {
		t.Fatalf("assembled termcode was not rejected: %+v", result)
	}
}

func TestRejectsExModeEntry(t *testing.T) {
	keys := []string{"g", "Q", "!", "e", "c", "h", "o", "<CR>"}
	result, err := Run("grid_rush", Input{Keys: keys})
	if err != nil {
		t.Fatal(err)
	}
	if result.Verified || result.Reason != "invalid key in log" {
		t.Fatalf("Ex-mode entry was not rejected: %+v", result)
	}
}

func TestStartPositionIsNotGoal(t *testing.T) {
	result, err := Run("repeat_hunter", Input{Keys: []string{"h"}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Verified {
		t.Fatal("a no-op key verified a challenge already solved at its start")
	}
}

func TestReportsShortestPrefix(t *testing.T) {
	tr, _ := content.TraceFor("grid_rush")
	var keys []string
	for _, r := range tr.Keys {
		keys = append(keys, string(r))
	}
	keys = append(keys, "\x1b") // harmless trailing key
	result, err := Run("grid_rush", Input{Keys: keys})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Verified || result.Keys != len(tr.Keys) {
		t.Fatalf("matched keys = %d, want %d: %+v", result.Keys, len(tr.Keys), result)
	}
}
