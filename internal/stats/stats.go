package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	// ErrCorrupt identifies a progress file that could not be decoded.
	ErrCorrupt = errors.New("corrupt progress file")
	// ErrUnsupportedVersion identifies a file written by an incompatible release.
	ErrUnsupportedVersion = errors.New("unsupported progress version")
	// ErrInvalidStats identifies semantically impossible progress values.
	ErrInvalidStats = errors.New("invalid progress values")
)

const (
	// currentVersion is written by Save. Version 1 is a bare map of challenge
	// IDs to Stats and remains readable for backward compatibility.
	currentVersion = 2

	lockTimeout   = 2 * time.Second
	staleLockAge  = 10 * time.Second
	lockRetryWait = 25 * time.Millisecond
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

// acquireLock serializes read-modify-write cycles across CLI processes. A
// stale lock (left by a crashed process) is reclaimed after staleLockAge.
func acquireLock() (func(), error) {
	lock := path() + ".lock"
	if err := os.MkdirAll(filepath.Dir(lock), 0o755); err != nil {
		return nil, err
	}
	deadline := time.Now().Add(lockTimeout)
	for {
		f, err := os.OpenFile(lock, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			_, _ = fmt.Fprintln(f, os.Getpid())
			if closeErr := f.Close(); closeErr != nil {
				_ = os.Remove(lock)
				return nil, closeErr
			}
			return func() { _ = os.Remove(lock) }, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}
		if info, statErr := os.Stat(lock); statErr == nil && time.Since(info.ModTime()) > staleLockAge {
			_ = os.Remove(lock)
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("progress file is locked by another process")
		}
		time.Sleep(lockRetryWait)
	}
}

// Mutate applies fn to the persisted progress map under a lock. Load and Save
// remain available for read-only callers; only writers need this cycle.
func Mutate(fn func(map[string]Stats) error) error {
	unlock, err := acquireLock()
	if err != nil {
		return err
	}
	defer unlock()

	m, err := Load()
	if err != nil {
		return err
	}
	if err := fn(m); err != nil {
		return err
	}
	return Save(m)
}

// progressFile is the version 2 on-disk envelope. Version 1 files are a bare
// map[string]Stats and remain readable by decodeProgress.
type progressFile struct {
	Version int              `json:"version"`
	Entries map[string]Stats `json:"entries"`
}

func decodeProgress(data []byte) (map[string]Stats, error) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCorrupt, err)
	}
	if top == nil {
		// "null" is valid JSON but not a valid progress map.
		return make(map[string]Stats), nil
	}

	var m map[string]Stats
	if rawVersion, ok := top["version"]; ok {
		var version int
		if err := json.Unmarshal(rawVersion, &version); err != nil {
			return nil, fmt.Errorf("%w: invalid version: %v", ErrCorrupt, err)
		}
		if version != currentVersion {
			return nil, fmt.Errorf("%w: %d (supported: %d)", ErrUnsupportedVersion, version, currentVersion)
		}
		rawEntries, ok := top["entries"]
		if !ok {
			return nil, fmt.Errorf("%w: missing entries", ErrCorrupt)
		}
		if err := json.Unmarshal(rawEntries, &m); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCorrupt, err)
		}
	} else {
		// Version 1: a bare map of challenge ID to Stats.
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCorrupt, err)
		}
	}
	if m == nil {
		m = make(map[string]Stats)
	}
	if err := validateProgress(m); err != nil {
		return nil, err
	}
	return m, nil
}

// validateProgress rejects impossible counters. Unknown challenge IDs remain
// non-fatal: stale keys are ignored by the CLI summary and must survive a
// load/save round trip. Record-only completions are also legitimate.
func validateProgress(m map[string]Stats) error {
	for id, s := range m {
		if s.Attempts < 0 || s.Completions < 0 || s.BestMoves < 0 || s.BestTimeMs < 0 {
			return fmt.Errorf("%w: challenge %q has negative values", ErrInvalidStats, id)
		}
	}
	return nil
}

func Load() (map[string]Stats, error) {
	data, err := os.ReadFile(path())
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]Stats), nil
		}
		return make(map[string]Stats), err
	}
	m, err := decodeProgress(data)
	if err != nil {
		if errors.Is(err, ErrCorrupt) {
			backup := path() + ".corrupt"
			if backupErr := os.WriteFile(backup, data, 0o644); backupErr != nil {
				return make(map[string]Stats), fmt.Errorf("%w; backup could not be saved: %v", err, backupErr)
			}
			return make(map[string]Stats), fmt.Errorf("%w; backup saved to %s", err, backup)
		}
		return make(map[string]Stats), err
	}
	return m, nil
}

func Save(m map[string]Stats) error {
	dir := filepath.Dir(path())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(progressFile{Version: currentVersion, Entries: m}, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "stats.json.*")
	if err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		// Cleanup is best-effort; preserve the write error as the result.
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Sync(); err != nil {
		// Flush before rename so a crash cannot leave an empty progress file.
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), path())
}

func Reset() error {
	unlock, err := acquireLock()
	if err != nil {
		return err
	}
	defer unlock()

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
	r := float64(s.Completions) / float64(s.Attempts)
	if r < 0 {
		return 0
	}
	if r > 1 {
		return 1
	}
	return r
}

func Confidence(s Stats) string {
	if s.Attempts == 0 {
		return "—"
	}
	switch r := SuccessRate(s); {
	case r >= 0.80:
		return "High"
	case r >= 0.60:
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
	var practiced, experienced string
	var practicedDate, experiencedDate string
	for _, id := range curriculum {
		s := m[id]
		if s.LastPlayed == "" {
			continue
		}
		switch MasteryTier(s) {
		case "Practiced":
			if practicedDate == "" || s.LastPlayed < practicedDate {
				practicedDate, practiced = s.LastPlayed, id
			}
		case "Experienced":
			if experiencedDate == "" || s.LastPlayed < experiencedDate {
				experiencedDate, experienced = s.LastPlayed, id
			}
		}
	}
	if practiced != "" {
		return practiced
	}
	return experienced
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
	unseen := totalChallenges - len(m)
	if unseen < 0 {
		unseen = 0
	}
	dist := map[string]int{
		"Unseen":      unseen,
		"Learning":    0,
		"Practiced":   0,
		"Experienced": 0,
	}
	for _, s := range m {
		dist[MasteryTier(s)]++
	}
	return dist
}
