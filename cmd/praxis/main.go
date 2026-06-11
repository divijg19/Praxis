package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/divijg19/Praxis/internal/content"
	"github.com/divijg19/Praxis/internal/stats"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "list":
			list()
			return
		case "challenge":
			if len(os.Args) > 2 {
				challenge(os.Args[2])
				return
			}
		case "target":
			if len(os.Args) > 2 {
				target(os.Args[2])
				return
			}
		case "verify":
			if len(os.Args) > 2 {
				verify(os.Args[2])
				return
			}
		case "result":
			if len(os.Args) > 2 {
				result(os.Args[2])
				return
			}
		case "record":
			if len(os.Args) > 4 {
				record(os.Args[2], os.Args[3], os.Args[4])
				return
			}
		case "stats":
			if len(os.Args) > 2 {
				statsForID(os.Args[2])
				return
			}
			statsSummary()
			return
		}
	}
	fmt.Println("Praxis")
	fmt.Println("Mastery through practice.")
}

func list() {
	for _, c := range content.All() {
		fmt.Println(c.Name)
	}
}

func challenge(id string) {
	for _, c := range content.All() {
		if c.ID == id {
			for _, line := range c.Content {
				fmt.Println(line)
			}
			return
		}
	}
	fmt.Fprintln(os.Stderr, "unknown challenge:", id)
	os.Exit(1)
}

func target(id string) {
	for _, c := range content.All() {
		if c.ID == id {
			fmt.Println(c.Target)
			return
		}
	}
	fmt.Fprintln(os.Stderr, "unknown challenge:", id)
	os.Exit(1)
}

func verify(id string) {
	for _, c := range content.All() {
		if c.ID == id {
			fmt.Println(c.Verify)
			return
		}
	}
	fmt.Fprintln(os.Stderr, "unknown challenge:", id)
	os.Exit(1)
}

func result(id string) {
	for _, c := range content.All() {
		if c.ID == id {
			for _, line := range c.Result {
				fmt.Println(line)
			}
			return
		}
	}
	fmt.Fprintln(os.Stderr, "unknown challenge:", id)
	os.Exit(1)
}

func record(id, movesStr, timeStr string) {
	moves, _ := strconv.Atoi(movesStr)
	timeMs, _ := strconv.Atoi(timeStr)
	m, _ := stats.Load()
	stats.Update(m, id, moves, timeMs)
	stats.Save(m)
}

func statsForID(id string) {
	found := false
	for _, c := range content.All() {
		if c.ID == id {
			found = true
			break
		}
	}
	if !found {
		fmt.Fprintln(os.Stderr, "unknown challenge:", id)
		os.Exit(1)
	}
	m, _ := stats.Load()
	s := m[id]
	fmt.Printf("Attempts: %d\n", s.Attempts)
	fmt.Printf("Completions: %d\n", s.Completions)
	fmt.Printf("Best Moves: %d\n", s.BestMoves)
	fmt.Printf("Best Time: %dms\n", s.BestTimeMs)
	fmt.Printf("Mastery: %s\n", stats.MasteryTier(s))
}

func statsSummary() {
	m, _ := stats.Load()
	var completed int
	var totalAttempts int
	for _, c := range content.All() {
		s := m[c.ID]
		if s.Completions > 0 {
			completed++
		}
		totalAttempts += s.Attempts
	}
	fmt.Printf("Challenges Completed: %d/%d\n", completed, len(content.All()))
	fmt.Printf("Total Attempts: %d\n", totalAttempts)
	fmt.Println()
	dist := stats.MasteryDistribution(m, len(content.All()))
	fmt.Println("Mastery:")
	fmt.Printf("  Unseen: %d\n", dist["Unseen"])
	fmt.Printf("  Learning: %d\n", dist["Learning"])
	fmt.Printf("  Practiced: %d\n", dist["Practiced"])
	fmt.Printf("  Experienced: %d\n", dist["Experienced"])
	tiers := []string{"Unseen", "Learning", "Practiced", "Experienced"}
	highest := "Unseen"
	for _, t := range tiers {
		if dist[t] > 0 {
			highest = t
		}
	}
	fmt.Printf("\nHighest Tier: %s\n", highest)
}
