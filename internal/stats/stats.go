// Package stats owns progress persistence and all derived progress facts.
// v0.4.0 uses one canonical, strict schema: no legacy readers, no migrations.
package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/divijg19/Praxis/internal/score"
)

var (
	ErrCorrupt      = errors.New("corrupt progress file")
	ErrInvalidStats = errors.New("invalid progress values")
	ErrUnverified   = errors.New("completion was not verified")
)

const (
	currentFormat = 1
	maxEvents     = 2000

	lockTimeout   = 2 * time.Second
	staleLockAge  = 10 * time.Second
	lockRetryWait = 25 * time.Millisecond
)

// Stats is the per-challenge aggregate. Completions are verified completions
// only; Starts record effort and are never used as a completion proxy.
type Stats struct {
	Starts      int    `json:"starts"`
	Completions int    `json:"completions"`
	BestKeys    int    `json:"best_keys"`
	BestEdits   int    `json:"best_edits"`
	BestMarks   int    `json:"best_marks"`
	LastPlayed  string `json:"last_played"`
}

// Event is an immutable record of a verified completion. Raw key logs are
// deliberately not persisted.
type Event struct {
	ID       string `json:"id"`
	Keys     int    `json:"keys"`
	Par      int    `json:"par"`
	HintCost int    `json:"hint_cost,omitempty"`
	Marks    int    `json:"marks"`
	Edits    int    `json:"edits"`
	Seed     string `json:"seed,omitempty"`
	Verified bool   `json:"verified"`
}

// Progress is the single on-disk document.
type Progress struct {
	Format        int              `json:"format"`
	Entries       map[string]Stats `json:"entries"`
	Events        []Event          `json:"events,omitempty"`
	EventsTrimmed int              `json:"events_trimmed,omitempty"`
}

func path() string {
	xdg := os.Getenv("XDG_DATA_HOME")
	if xdg == "" {
		home, _ := os.UserHomeDir()
		xdg = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(xdg, "praxis", "progress.json")
}

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

func empty() Progress {
	return Progress{Format: currentFormat, Entries: make(map[string]Stats)}
}

func decodeProgress(data []byte) (Progress, error) {
	var p Progress
	if err := json.Unmarshal(data, &p); err != nil {
		return empty(), fmt.Errorf("%w: %v", ErrCorrupt, err)
	}
	if p.Format != currentFormat {
		return empty(), fmt.Errorf("%w: format %d (supported: %d)", ErrCorrupt, p.Format, currentFormat)
	}
	if p.Entries == nil {
		p.Entries = make(map[string]Stats)
	}
	if err := validate(p); err != nil {
		return empty(), err
	}
	return p, nil
}

func validate(p Progress) error {
	for id, s := range p.Entries {
		if s.Starts < 0 || s.Completions < 0 || s.BestKeys < 0 || s.BestEdits < 0 || s.BestMarks < 0 {
			return fmt.Errorf("%w: challenge %q has negative values", ErrInvalidStats, id)
		}
	}
	for _, e := range p.Events {
		if !e.Verified || e.Keys <= 0 || e.Marks < 1 || e.Marks > 3 {
			return fmt.Errorf("%w: invalid completion event", ErrInvalidStats)
		}
	}
	return nil
}

// Load returns the canonical progress document. A missing file is fresh state.
func Load() (Progress, error) {
	data, err := os.ReadFile(path())
	if err != nil {
		if os.IsNotExist(err) {
			return empty(), nil
		}
		return empty(), err
	}
	p, err := decodeProgress(data)
	if err == nil {
		return p, nil
	}
	if errors.Is(err, ErrCorrupt) || errors.Is(err, ErrInvalidStats) {
		backup := path() + ".corrupt"
		if backupErr := os.WriteFile(backup, data, 0o600); backupErr != nil {
			return empty(), fmt.Errorf("%w; backup could not be saved: %v", err, backupErr)
		}
		return empty(), fmt.Errorf("%w; backup saved to %s", err, backup)
	}
	return empty(), err
}

// Save writes progress atomically and durably.
func Save(p Progress) error {
	if err := validate(p); err != nil {
		return err
	}
	p.Format = currentFormat
	dir := filepath.Dir(path())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "progress.json.*")
	if err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Sync(); err != nil {
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

// Mutate applies fn to progress under the inter-process lock.
func Mutate(fn func(*Progress) error) error {
	unlock, err := acquireLock()
	if err != nil {
		return err
	}
	defer unlock()
	p, err := Load()
	if err != nil {
		return err
	}
	if err := fn(&p); err != nil {
		return err
	}
	return Save(p)
}

func Reset() error {
	unlock, err := acquireLock()
	if err != nil {
		return err
	}
	defer unlock()
	for _, p := range []string{path(), path() + ".corrupt"} {
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

// Start records that a challenge was opened. It never creates a completion.
func Start(p *Progress, id string) {
	s := p.Entries[id]
	s.Starts++
	p.Entries[id] = s
}

// Complete records a verified completion. Marks are recomputed here; clients
// cannot supply them.
func Complete(p *Progress, e Event) (Stats, int, error) {
	if !e.Verified {
		return Stats{}, 0, ErrUnverified
	}
	e.Marks = score.Marks(e.Keys, e.Par, e.HintCost)
	if e.Marks < 1 || e.Marks > 3 {
		return Stats{}, 0, fmt.Errorf("%w: invalid mark", ErrInvalidStats)
	}
	s := p.Entries[e.ID]
	s.Completions++
	s.LastPlayed = time.Now().Format("2006-01-02")
	if s.Completions == 1 || e.Keys < s.BestKeys {
		s.BestKeys = e.Keys
	}
	if s.Completions == 1 || e.Edits < s.BestEdits {
		s.BestEdits = e.Edits
	}
	if e.Marks > s.BestMarks {
		s.BestMarks = e.Marks
	}
	p.Entries[e.ID] = s
	p.Events = append(p.Events, e)
	if len(p.Events) > maxEvents {
		p.EventsTrimmed += len(p.Events) - maxEvents
		p.Events = append([]Event(nil), p.Events[len(p.Events)-maxEvents:]...)
	}
	return s, e.Marks, nil
}

const (
	learningMax    = 2
	experiencedMin = 8
)

func NextChallenge(p Progress, curriculum []string) string {
	for _, id := range curriculum {
		if p.Entries[id].Completions <= learningMax {
			return id
		}
	}
	return ""
}

func RecommendedReview(p Progress, curriculum []string) string {
	var practiced, experienced string
	var practicedDate, experiencedDate string
	for _, id := range curriculum {
		s := p.Entries[id]
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

func MasteryDistribution(p Progress, totalChallenges int) map[string]int {
	unseen := totalChallenges - len(p.Entries)
	if unseen < 0 {
		unseen = 0
	}
	dist := map[string]int{
		"Unseen": unseen, "Learning": 0, "Practiced": 0, "Experienced": 0,
	}
	for _, s := range p.Entries {
		dist[MasteryTier(s)]++
	}
	return dist
}
