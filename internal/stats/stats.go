package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Stats struct {
	Attempts    int    `json:"attempts"`
	Completions int    `json:"completions"`
	BestMoves   int    `json:"best_moves"`
	BestTimeMs  int    `json:"best_time_ms"`
	LastPlayed  string `json:"last_played"`
}

func Path() string {
	xdg := os.Getenv("XDG_DATA_HOME")
	if xdg == "" {
		home, _ := os.UserHomeDir()
		xdg = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(xdg, "praxis", "stats.json")
}

func Load() (map[string]Stats, error) {
	data, err := os.ReadFile(Path())
	if err != nil {
		return make(map[string]Stats), nil
	}
	var m map[string]Stats
	if err := json.Unmarshal(data, &m); err != nil {
		return make(map[string]Stats), nil
	}
	if m == nil {
		return make(map[string]Stats), nil
	}
	return m, nil
}

func Save(m map[string]Stats) error {
	dir := filepath.Dir(Path())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "stats.json.*")
	if err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), Path())
}

func Update(m map[string]Stats, id string, moves, timeMs int) Stats {
	s := m[id]
	s.Attempts++
	s.Completions++
	s.LastPlayed = time.Now().Format("2006-01-02")
	if s.Completions == 1 {
		s.BestMoves = moves
		s.BestTimeMs = timeMs
	} else {
		if moves < s.BestMoves {
			s.BestMoves = moves
		}
		if timeMs < s.BestTimeMs {
			s.BestTimeMs = timeMs
		}
	}
	m[id] = s
	return s
}

const (
	LearningMax    = 2
	ExperiencedMin = 8
)

func NextChallenge(m map[string]Stats, curriculum []string) string {
	for _, id := range curriculum {
		if m[id].Completions <= LearningMax {
			return id
		}
	}
	return ""
}

func RecommendedReview(m map[string]Stats, curriculum []string) string {
	var oldestPracticed, oldestExperienced string
	var practicedDate, experiencedDate string
	for _, id := range curriculum {
		s := m[id]
		if s.LastPlayed == "" {
			continue
		}
		switch MasteryTier(s) {
		case "Practiced":
			if practicedDate == "" || s.LastPlayed < practicedDate {
				practicedDate = s.LastPlayed
				oldestPracticed = id
			}
		case "Experienced":
			if experiencedDate == "" || s.LastPlayed < experiencedDate {
				experiencedDate = s.LastPlayed
				oldestExperienced = id
			}
		}
	}
	if oldestPracticed != "" {
		return oldestPracticed
	}
	return oldestExperienced
}

func MasteryTier(s Stats) string {
	switch {
	case s.Completions >= ExperiencedMin:
		return "Experienced"
	case s.Completions > LearningMax:
		return "Practiced"
	case s.Completions >= 1:
		return "Learning"
	default:
		return "Unseen"
	}
}

func MasteryDistribution(m map[string]Stats, totalChallenges int) map[string]int {
	dist := map[string]int{
		"Unseen":      totalChallenges - len(m),
		"Learning":    0,
		"Practiced":   0,
		"Experienced": 0,
	}
	for _, s := range m {
		dist[MasteryTier(s)]++
	}
	return dist
}
