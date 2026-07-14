package stats

import (
	"encoding/json"
	"fmt"
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

func path() string {
	xdg := os.Getenv("XDG_DATA_HOME")
	if xdg == "" {
		home, _ := os.UserHomeDir()
		xdg = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(xdg, "praxis", "stats.json")
}

func Load() (map[string]Stats, error) {
	data, err := os.ReadFile(path())
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]Stats), nil
		}
		return make(map[string]Stats), err
	}
	var m map[string]Stats
	if err := json.Unmarshal(data, &m); err != nil {
		backup := path() + ".corrupt"
		_ = os.WriteFile(backup, data, 0644)
		return make(map[string]Stats), fmt.Errorf("corrupt progress file; backup saved to %s", backup)
	}
	if m == nil {
		m = make(map[string]Stats)
	}
	return m, nil
}

func Save(m map[string]Stats) error {
	dir := filepath.Dir(path())
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
	return os.Rename(tmp.Name(), path())
}

func Reset() error {
	p := path()
	if err := os.Remove(p); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func Attempt(m map[string]Stats, id string) Stats {
	s := m[id]
	s.Attempts++
	m[id] = s
	return s
}

func SuccessRate(s Stats) float64 {
	if s.Attempts == 0 {
		return 0
	}
	return float64(s.Completions) / float64(s.Attempts)
}

func Confidence(s Stats) string {
	if s.Attempts == 0 {
		return "—"
	}
	switch {
	case SuccessRate(s) >= 0.80:
		return "High"
	case SuccessRate(s) >= 0.60:
		return "Medium"
	default:
		return "Low"
	}
}

func Update(m map[string]Stats, id string, moves, timeMs int) Stats {
	s := m[id]
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
	learningMax    = 2
	experiencedMin = 8
)

func NextChallenge(m map[string]Stats, curriculum []string) string {
	for _, id := range curriculum {
		if m[id].Completions <= learningMax {
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
	case s.Completions >= experiencedMin:
		return "Experienced"
	case s.Completions > learningMax:
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
