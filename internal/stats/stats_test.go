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

func TestMasteryUnseen(t *testing.T) {
	s := Stats{}
	if got := MasteryTier(s); got != "Unseen" {
		t.Errorf("MasteryTier(%+v) = %q, want Unseen", s, got)
	}
}

func TestMasteryLearning(t *testing.T) {
	s := Stats{Completions: 1}
	if got := MasteryTier(s); got != "Learning" {
		t.Errorf("MasteryTier(%+v) = %q, want Learning", s, got)
	}
	s.Completions = 2
	if got := MasteryTier(s); got != "Learning" {
		t.Errorf("MasteryTier(%+v) = %q, want Learning", s, got)
	}
}

func TestMasteryPracticed(t *testing.T) {
	s := Stats{Completions: 3}
	if got := MasteryTier(s); got != "Practiced" {
		t.Errorf("MasteryTier(%+v) = %q, want Practiced", s, got)
	}
	s.Completions = 7
	if got := MasteryTier(s); got != "Practiced" {
		t.Errorf("MasteryTier(%+v) = %q, want Practiced", s, got)
	}
}

func TestMasteryExperienced(t *testing.T) {
	s := Stats{Completions: 8}
	if got := MasteryTier(s); got != "Experienced" {
		t.Errorf("MasteryTier(%+v) = %q, want Experienced", s, got)
	}
	s.Completions = 20
	if got := MasteryTier(s); got != "Experienced" {
		t.Errorf("MasteryTier(%+v) = %q, want Experienced", s, got)
	}
}

func TestMasteryDistributionEmpty(t *testing.T) {
	m := make(map[string]Stats)
	d := MasteryDistribution(m, 41)
	if d["Unseen"] != 41 {
		t.Errorf("Unseen = %d, want 41", d["Unseen"])
	}
	if d["Learning"] != 0 {
		t.Errorf("Learning = %d, want 0", d["Learning"])
	}
	if d["Practiced"] != 0 {
		t.Errorf("Practiced = %d, want 0", d["Practiced"])
	}
	if d["Experienced"] != 0 {
		t.Errorf("Experienced = %d, want 0", d["Experienced"])
	}
}

func TestMasteryDistributionMixed(t *testing.T) {
	m := map[string]Stats{
		"a": {Completions: 1},
		"b": {Completions: 2},
		"c": {Completions: 5},
		"d": {Completions: 10},
	}
	d := MasteryDistribution(m, 41)
	if d["Unseen"] != 37 {
		t.Errorf("Unseen = %d, want 37", d["Unseen"])
	}
	if d["Learning"] != 2 {
		t.Errorf("Learning = %d, want 2", d["Learning"])
	}
	if d["Practiced"] != 1 {
		t.Errorf("Practiced = %d, want 1", d["Practiced"])
	}
	if d["Experienced"] != 1 {
		t.Errorf("Experienced = %d, want 1", d["Experienced"])
	}
}

func TestNextChallengeEmpty(t *testing.T) {
	m := make(map[string]Stats)
	curric := []string{"a", "b", "c"}
	if got := NextChallenge(m, curric); got != "a" {
		t.Errorf("NextChallenge(empty) = %q, want %q", got, "a")
	}
}

func TestNextChallengePartial(t *testing.T) {
	m := map[string]Stats{
		"a": {Completions: 5},
		"b": {Completions: 3},
		"c": {Completions: 0},
	}
	curric := []string{"a", "b", "c"}
	if got := NextChallenge(m, curric); got != "c" {
		t.Errorf("NextChallenge(partial) = %q, want %q", got, "c")
	}
}

func TestNextChallengeComplete(t *testing.T) {
	m := map[string]Stats{
		"a": {Completions: 5},
		"b": {Completions: 8},
	}
	curric := []string{"a", "b"}
	if got := NextChallenge(m, curric); got != "" {
		t.Errorf("NextChallenge(all done) = %q, want %q", got, "")
	}
}

func TestRecommendedReviewPracticed(t *testing.T) {
	m := map[string]Stats{
		"a": {Completions: 5, LastPlayed: "2026-06-01"},
		"b": {Completions: 4, LastPlayed: "2026-06-10"},
	}
	curric := []string{"a", "b"}
	if got := RecommendedReview(m, curric); got != "a" {
		t.Errorf("RecommendedReview = %q, want %q", got, "a")
	}
}

func TestRecommendedReviewFallbackExperienced(t *testing.T) {
	m := map[string]Stats{
		"a": {Completions: 10, LastPlayed: "2026-05-01"},
		"b": {Completions: 12, LastPlayed: "2026-06-01"},
	}
	curric := []string{"a", "b"}
	if got := RecommendedReview(m, curric); got != "a" {
		t.Errorf("RecommendedReview = %q, want %q", got, "a")
	}
}

func TestRecommendedReviewPrefersPracticedOverExperienced(t *testing.T) {
	m := map[string]Stats{
		"a": {Completions: 5, LastPlayed: "2026-06-10"},
		"b": {Completions: 10, LastPlayed: "2026-06-01"},
	}
	curric := []string{"a", "b"}
	if got := RecommendedReview(m, curric); got != "a" {
		t.Errorf("RecommendedReview = %q (expected Practiced a over Experienced b), want %q", got, "a")
	}
}
