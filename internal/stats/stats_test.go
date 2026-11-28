package stats

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func writeProgress(t *testing.T, data string) {
	t.Helper()
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	if err := os.MkdirAll(filepath.Join(d, "praxis"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "praxis", "progress.json"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadMissing(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	p, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if p.Format != currentFormat || len(p.Entries) != 0 {
		t.Errorf("unexpected fresh progress: %+v", p)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	p := empty()
	p.Entries["grid_rush"] = Stats{Starts: 2, Completions: 1, BestKeys: 5, BestEdits: 1, BestMarks: 3}
	if err := Save(p); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := loaded.Entries["grid_rush"].BestKeys; got != 5 {
		t.Errorf("BestKeys = %d, want 5", got)
	}
}

func TestCorruptAndUnsupportedRejected(t *testing.T) {
	for _, tc := range []struct{ data, want string }{
		{`{broken`, "corrupt"},
		{`{"format":99,"entries":{}}`, "format 99"},
		{`{"format":1,"entries":{"a":{"starts":-1}}}`, "negative"},
	} {
		writeProgress(t, tc.data)
		_, err := Load()
		if err == nil {
			t.Fatalf("expected error for %s", tc.data)
		}
		if !strings.Contains(err.Error(), tc.want) {
			t.Errorf("error = %q, want containing %q", err, tc.want)
		}
	}
}

func TestMutateSerializes(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	const writers = 50
	var wg sync.WaitGroup
	errs := make(chan error, writers)
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- Mutate(func(p *Progress) error { Start(p, "grid_rush"); return nil })
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	p, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := p.Entries["grid_rush"].Starts; got != writers {
		t.Errorf("Starts = %d, want %d", got, writers)
	}
}

func TestCompleteRejectsUnverified(t *testing.T) {
	p := empty()
	if _, _, err := Complete(&p, Event{ID: "grid_rush", Keys: 1, Par: 1, Verified: false}); !errors.Is(err, ErrUnverified) {
		t.Fatalf("error = %v, want ErrUnverified", err)
	}
	if len(p.Entries) != 0 {
		t.Fatal("unverified completion wrote progress")
	}
}

func TestCompleteUsesGoMarks(t *testing.T) {
	p := empty()
	s, marks, err := Complete(&p, Event{ID: "grid_rush", Keys: 7, Par: 5, Marks: 99, Edits: 1, Verified: true})
	if err != nil {
		t.Fatal(err)
	}
	if marks != 2 || s.BestMarks != 2 {
		t.Errorf("marks=%d BestMarks=%d, client-supplied 99 must be ignored", marks, s.BestMarks)
	}
	if len(p.Events) != 1 || p.Events[0].Marks != 2 {
		t.Errorf("event marks = %+v", p.Events)
	}
}

func TestResetRemovesBackup(t *testing.T) {
	writeProgress(t, `{broken`)
	if _, err := Load(); err == nil {
		t.Fatal("expected corrupt error")
	}
	if err := Reset(); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"progress.json", "progress.json.corrupt"} {
		p := filepath.Join(os.Getenv("XDG_DATA_HOME"), "praxis", name)
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("%s still exists after reset", name)
		}
	}
}

func TestRevealHintAndStartClears(t *testing.T) {
	p := empty()
	RevealHint(&p, "grid_rush", 2)
	if got := RevealedHintLevel(p, "grid_rush"); got != 2 {
		t.Fatalf("level = %d, want 2", got)
	}
	RevealHint(&p, "grid_rush", 1)
	if got := RevealedHintLevel(p, "grid_rush"); got != 2 {
		t.Fatalf("level regressed to %d", got)
	}
	Start(&p, "grid_rush")
	if got := RevealedHintLevel(p, "grid_rush"); got != 0 {
		t.Fatalf("Start did not clear level: %d", got)
	}
}

func TestRejectsNegativeHintCost(t *testing.T) {
	p := empty()
	p.Events = []Event{{ID: "grid_rush", Keys: 1, Par: 1, Marks: 3, HintCost: -1, Verified: true}}
	if err := validate(p); err == nil {
		t.Fatal("negative hint cost accepted")
	}
}
