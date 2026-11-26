package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/divijg19/Praxis/internal/content"
	"github.com/divijg19/Praxis/internal/stats"
)

func usage(msg string) {
	fmt.Fprintln(os.Stderr, "praxis: usage:", msg)
	os.Exit(1)
}

func unknown(id string) {
	fmt.Fprintln(os.Stderr, "unknown challenge:", id)
	os.Exit(1)
}

func positiveInt(s, what string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		fmt.Fprintf(os.Stderr, "praxis: usage: praxis record <id> <moves> <time> (%s must be a positive integer)\n", what)
		return 0, false
	}
	return n, true
}

func banner() {
	fmt.Println("Praxis")
	fmt.Println("Mastery through practice.")
}

func corruptProgress() {
	fmt.Fprintln(os.Stderr, "praxis: progress file is corrupt")
	fmt.Println("Your Praxis progress is corrupted.")
	fmt.Println("If a backup was saved, it is next to the progress file. Run `praxis reset` to start fresh.")
}

func loadOrFail() map[string]stats.Stats {
	m, err := stats.Load()
	if err != nil {
		if errors.Is(err, stats.ErrCorrupt) {
			corruptProgress()
		} else {
			fmt.Fprintln(os.Stderr, "praxis: could not read progress:", err)
		}
		os.Exit(1)
	}
	return m
}

func mutateOrFail(fn func(map[string]stats.Stats) error) {
	if err := stats.Mutate(fn); err != nil {
		if errors.Is(err, stats.ErrCorrupt) {
			corruptProgress()
		} else {
			fmt.Fprintln(os.Stderr, "praxis: could not update progress:", err)
		}
		os.Exit(1)
	}
}

func main() {
	if err := content.ValidateAll(); err != nil {
		fmt.Fprintln(os.Stderr, "praxis: invalid challenge catalog:", err)
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		entrypoint()
		return
	}
	switch os.Args[1] {
	case "help":
		help()
	case "next":
		if len(os.Args) != 2 {
			usage("praxis next")
		}
		next()
	case "catalog":
		switch len(os.Args) {
		case 2:
			catalog("")
		case 3:
			if os.Args[2] != "--ids" && os.Args[2] != "--core" {
				usage("praxis catalog [--ids|--core]")
			}
			catalog(os.Args[2])
		default:
			usage("praxis catalog [--ids|--core]")
		}
	case "describe":
		if len(os.Args) != 3 {
			usage("praxis describe <id>")
		}
		describe(os.Args[2])
	case "attempt":
		if len(os.Args) != 3 {
			usage("praxis attempt <id>")
		}
		attempt(os.Args[2])
	case "record":
		if len(os.Args) != 5 {
			usage("praxis record <id> <moves> <time>")
		}
		record(os.Args[2], os.Args[3], os.Args[4])
	case "reset":
		switch len(os.Args) {
		case 2:
			reset(false)
		case 3:
			if os.Args[2] != "--yes" {
				usage("praxis reset [--yes]")
			}
			reset(true)
		default:
			usage("praxis reset [--yes]")
		}
	case "stats":
		switch len(os.Args) {
		case 2:
			statsSummary()
		case 3:
			statsForID(os.Args[2])
		default:
			usage("praxis stats [<id>]")
		}
	default:
		entrypoint()
	}
}

func tutorialProgress(m map[string]stats.Stats) (done, total int) {
	for _, id := range content.CoreTutorialIDs() {
		total++
		if s, ok := m[id]; ok && s.Completions > 0 {
			done++
		}
	}
	return done, total
}

func entrypoint() {
	banner()
	fmt.Println()

	m := loadOrFail()
	id := stats.NextChallenge(m, content.IDs())
	fmt.Println("Next:")
	if id != "" {
		fmt.Printf("  %s\n", id)
	} else {
		fmt.Println("  Curriculum complete.")
	}
	fmt.Println()
	fmt.Println("Run:")
	fmt.Println("  praxis next")
	fmt.Println("  praxis help")
}

func help() {
	banner()
	fmt.Println()
	fmt.Println("Start:")
	fmt.Println("  praxis next")
	fmt.Println()
	fmt.Println("Track:")
	fmt.Println("  praxis attempt <id>")
	fmt.Println("  praxis record <id> <moves> <time>")
	fmt.Println()
	fmt.Println("Progress:")
	fmt.Println("  praxis stats")
	fmt.Println("  praxis reset")
	fmt.Println()
	fmt.Println("Explore:")
	fmt.Println("  praxis catalog [--ids|--core]")
	fmt.Println()
	fmt.Println("Inspect:")
	fmt.Println("  praxis describe <id>")
}

func describe(id string) {
	d, ok := content.DescriptionFor(id)
	if !ok {
		unknown(id)
	}
	if err := json.NewEncoder(os.Stdout).Encode(d); err != nil {
		fmt.Fprintln(os.Stderr, "describe failed:", err)
		os.Exit(1)
	}
}

func catalog(flag string) {
	switch flag {
	case "--ids":
		for _, id := range content.IDs() {
			fmt.Println(id)
		}
	case "--core":
		for _, id := range content.CoreTutorialIDs() {
			fmt.Println(id)
		}
	default:
		for _, c := range content.All() {
			fmt.Println(c.Name)
		}
	}
}

func next() {
	id := stats.NextChallenge(loadOrFail(), content.IDs())
	if id != "" {
		fmt.Println(id)
	}
}

func reset(confirmed bool) {
	if !confirmed {
		fmt.Print("This will erase all Praxis progress.\n\nType RESET to continue: ")
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			fmt.Println("Reset cancelled.")
			os.Exit(1)
		}
		if input != "RESET" {
			fmt.Println("Reset cancelled.")
			os.Exit(1)
		}
	}
	if err := stats.Reset(); err != nil {
		fmt.Fprintln(os.Stderr, "reset failed:", err)
		os.Exit(1)
	}
	fmt.Println("Praxis has been reset.")
	fmt.Println("Start fresh with: praxis next")
}

func attempt(id string) {
	if !content.Exists(id) {
		unknown(id)
	}
	mutateOrFail(func(m map[string]stats.Stats) error {
		stats.Attempt(m, id)
		return nil
	})
}

func record(id, movesStr, timeStr string) {
	if !content.Exists(id) {
		unknown(id)
	}
	moves, ok := positiveInt(movesStr, "moves")
	if !ok {
		os.Exit(1)
	}
	timeMs, ok := positiveInt(timeStr, "time")
	if !ok {
		os.Exit(1)
	}
	mutateOrFail(func(m map[string]stats.Stats) error {
		stats.Update(m, id, moves, timeMs)
		return nil
	})
}

func statsForID(id string) {
	if !content.Exists(id) {
		unknown(id)
	}
	m := loadOrFail()
	s := m[id]
	fmt.Printf("Attempts: %d\n", s.Attempts)
	fmt.Printf("Completions: %d\n", s.Completions)
	fmt.Printf("Success Rate: %.0f%%\n", stats.SuccessRate(s)*100)
	fmt.Printf("Best Moves: %d\n", s.BestMoves)
	fmt.Printf("Best Time: %d ms\n", s.BestTimeMs)
	fmt.Printf("Mastery: %s\n", stats.MasteryTier(s))
	fmt.Printf("Confidence: %s\n", stats.Confidence(s))
}

func statsSummary() {
	m := loadOrFail()
	ids := content.IDs()
	valid := make(map[string]stats.Stats, len(ids))
	var completed int
	var totalAttempts int
	for _, id := range ids {
		s, ok := m[id]
		if !ok {
			continue
		}
		valid[id] = s
		if s.Completions > 0 {
			completed++
		}
		totalAttempts += s.Attempts
	}
	fmt.Printf("Challenges Completed: %d/%d\n", completed, len(ids))
	fmt.Printf("Total Attempts: %d\n", totalAttempts)
	fmt.Println()
	coreDone, coreTotal := tutorialProgress(m)
	if coreDone >= coreTotal {
		fmt.Println("Tutorial: Complete")
	} else {
		fmt.Printf("Tutorial: %d/%d core\n", coreDone, coreTotal)
	}
	fmt.Println()
	dist := stats.MasteryDistribution(valid, len(ids))
	fmt.Println("Mastery:")
	fmt.Printf("  Unseen: %d\n", dist["Unseen"])
	fmt.Printf("  Learning: %d\n", dist["Learning"])
	fmt.Printf("  Practiced: %d\n", dist["Practiced"])
	fmt.Printf("  Experienced: %d\n", dist["Experienced"])
	tiers := []string{"Learning", "Practiced", "Experienced"}
	highest := ""
	for _, t := range tiers {
		if dist[t] > 0 {
			highest = t
		}
	}
	if highest != "" {
		fmt.Printf("\nMost mastered: %s\n", highest)
	}

	next := stats.NextChallenge(m, ids)
	fmt.Println()
	if next == "" {
		fmt.Println("Next:")
		fmt.Println("  Curriculum complete.")
	} else {
		fmt.Println("Next:")
		fmt.Printf("  %s\n", next)
	}

	review := stats.RecommendedReview(m, ids)
	if review != "" {
		fmt.Println()
		fmt.Println("Recommended Review:")
		fmt.Printf("  %s\n", review)
	}
}
