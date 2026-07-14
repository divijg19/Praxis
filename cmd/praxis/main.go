package main

import (
	"encoding/json"
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

func nonNegInt(s, what string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		fmt.Fprintf(os.Stderr, "praxis: usage: praxis record <id> <moves> <time> (%s must be a non-negative integer)\n", what)
		return 0, false
	}
	return n, true
}

func banner() {
	fmt.Println("Praxis")
	fmt.Println("Mastery through practice.")
}

func loadOrFail() map[string]stats.Stats {
	m, err := stats.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "praxis: progress file is corrupt")
		fmt.Println("Your Praxis progress is corrupted.")
		fmt.Println("A backup was saved next to it. Run `praxis reset` to start fresh; your backup is preserved.")
		os.Exit(1)
	}
	return m
}

func main() {
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
		if len(os.Args) != 2 {
			usage("praxis catalog")
		}
		catalog()
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
			reset("")
		case 3:
			if os.Args[2] != "--yes" {
				usage("praxis reset [--yes]")
			}
			reset("--yes")
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
	fmt.Println("Progress:")
	fmt.Println("  praxis stats")
	fmt.Println("  praxis reset")
	fmt.Println()
	fmt.Println("Explore:")
	fmt.Println("  praxis catalog")
	fmt.Println()
	fmt.Println("Inspect:")
	fmt.Println("  praxis describe <id>")
}

func describe(id string) {
	d, ok := content.DescriptionFor(id)
	if !ok {
		unknown(id)
	}
	json.NewEncoder(os.Stdout).Encode(d)
}

func catalog() {
	for _, c := range content.All() {
		fmt.Println(c.Name)
	}
}

func next() {
	id := stats.NextChallenge(loadOrFail(), content.IDs())
	if id != "" {
		fmt.Println(id)
	}
}

func reset(yesArg string) {
	yes := yesArg == "--yes"
	if !yes {
		fmt.Print("This will erase all Praxis progress.\n\nType RESET to continue: ")
		var input string
		fmt.Scanln(&input)
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
	m := loadOrFail()
	stats.Attempt(m, id)
	stats.Save(m)
}

func record(id, movesStr, timeStr string) {
	if !content.Exists(id) {
		unknown(id)
	}
	moves, ok := nonNegInt(movesStr, "moves")
	if !ok {
		os.Exit(1)
	}
	timeMs, ok := nonNegInt(timeStr, "time")
	if !ok {
		os.Exit(1)
	}
	m := loadOrFail()
	stats.Update(m, id, moves, timeMs)
	stats.Save(m)
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
	var completed int
	var totalAttempts int
	for _, id := range ids {
		s := m[id]
		if s.Completions > 0 {
			completed++
		}
		totalAttempts += s.Attempts
	}
	fmt.Printf("Challenges Completed: %d/%d\n", completed, len(ids))
	fmt.Printf("Total Attempts: %d\n", totalAttempts)
	fmt.Println()
	dist := stats.MasteryDistribution(m, len(ids))
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
