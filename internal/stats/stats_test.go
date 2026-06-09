package stats

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadFileNotExist(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	m, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 0 {
		t.Errorf("got %d entries, want 0", len(m))
	}
}

func TestSaveAndLoad(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	m := map[string]Stats{
		"motion_rush": {Attempts: 5, Completions: 4, BestMoves: 2, BestTimeMs: 150, LastPlayed: "2026-06-09"},
	}
	if err := Save(m); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	s := loaded["motion_rush"]
	if s.Attempts != 5 || s.Completions != 4 || s.BestMoves != 2 || s.BestTimeMs != 150 || s.LastPlayed != "2026-06-09" {
		t.Errorf("got %+v", s)
	}
}

func TestUpdateIncrements(t *testing.T) {
	m := make(map[string]Stats)
	Update(m, "a", 3, 200)
	Update(m, "a", 5, 300)
	s := m["a"]
	if s.Attempts != 2 {
		t.Errorf("Attempts = %d, want 2", s.Attempts)
	}
	if s.Completions != 2 {
		t.Errorf("Completions = %d, want 2", s.Completions)
	}
}

func TestUpdateBestMoves(t *testing.T) {
	m := make(map[string]Stats)
	Update(m, "a", 10, 100) // first: sets best to 10
	Update(m, "a", 3, 200)  // better: updates to 3
	Update(m, "a", 7, 150)  // worse: keeps 3
	s := m["a"]
	if s.BestMoves != 3 {
		t.Errorf("BestMoves = %d, want 3", s.BestMoves)
	}
}

func TestUpdateBestTime(t *testing.T) {
	m := make(map[string]Stats)
	Update(m, "a", 5, 1000) // first
	Update(m, "a", 3, 500)  // better time
	Update(m, "a", 7, 800)  // worse time
	s := m["a"]
	if s.BestTimeMs != 500 {
		t.Errorf("BestTimeMs = %d, want 500", s.BestTimeMs)
	}
}

func TestUpdateLastPlayed(t *testing.T) {
	m := make(map[string]Stats)
	Update(m, "a", 1, 50)
	s := m["a"]
	today := time.Now().Format("2006-01-02")
	if s.LastPlayed != today {
		t.Errorf("LastPlayed = %q, want %q", s.LastPlayed, today)
	}
}

func TestLoadCorruptFile(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	os.MkdirAll(filepath.Join(d, "praxis"), 0755)
	os.WriteFile(filepath.Join(d, "praxis", "stats.json"), []byte("{broken"), 0644)
	m, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 0 {
		t.Errorf("got %d entries, want 0", len(m))
	}
}

func TestUpdateReturnsUpdated(t *testing.T) {
	m := make(map[string]Stats)
	s := Update(m, "a", 5, 200)
	if s.Attempts != 1 || s.Completions != 1 || s.BestMoves != 5 || s.BestTimeMs != 200 {
		t.Errorf("returned %+v", s)
	}
}
