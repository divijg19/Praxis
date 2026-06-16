package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/divijg19/Praxis/internal/content"
	"github.com/divijg19/Praxis/internal/stats"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			help()
			return
		case "describe":
			if len(os.Args) > 2 {
				describe(os.Args[2])
				return
			}
		case "catalog":
			catalog()
			return
		case "next":
			next()
			return
		case "attempt":
			if len(os.Args) > 2 {
				attempt(os.Args[2])
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
	entrypoint()
}

func entrypoint() {
	fmt.Println("Praxis")
	fmt.Println("Mastery through practice.")
	fmt.Println()

	var curriculumIDs []string
	for _, c := range content.All() {
		curriculumIDs = append(curriculumIDs, c.ID)
	}
	m, _ := stats.Load()
	id := stats.NextChallenge(m, curriculumIDs)
	fmt.Println("Next:")
	if id != "" {
		fmt.Printf("  %s\n", id)
	} else {
		fmt.Println("  All challenges complete!")
	}
	fmt.Println()
	fmt.Println("Run:")
	fmt.Println("  praxis next")
	fmt.Println("  praxis help")
}

func help() {
	fmt.Println("Praxis")
	fmt.Println("Mastery through practice.")
	fmt.Println()
	fmt.Println("Start:")
	fmt.Println("  praxis next")
	fmt.Println()
	fmt.Println("Progress:")
	fmt.Println("  praxis stats")
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
		fmt.Fprintln(os.Stderr, "unknown challenge:", id)
		os.Exit(1)
	}
	json.NewEncoder(os.Stdout).Encode(d)
}

func catalog() {
	for _, c := range content.All() {
		fmt.Println(c.Name)
	}
}

func next() {
	var curriculumIDs []string
	for _, c := range content.All() {
		curriculumIDs = append(curriculumIDs, c.ID)
	}
	m, _ := stats.Load()
	id := stats.NextChallenge(m, curriculumIDs)
	if id != "" {
		fmt.Println(id)
	}
}

func attempt(id string) {
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
	stats.Attempt(m, id)
	stats.Save(m)
}

func record(id, movesStr, timeStr string) {
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
	fmt.Printf("Success Rate: %.0f%%\n", stats.SuccessRate(s)*100)
	fmt.Printf("Best Moves: %d\n", s.BestMoves)
	fmt.Printf("Best Time: %dms\n", s.BestTimeMs)
	fmt.Printf("Mastery: %s\n", stats.MasteryTier(s))
	fmt.Printf("Confidence: %s\n", stats.Confidence(s))
}

func statsSummary() {
	m, _ := stats.Load()
	var completed int
	var totalAttempts int
	var curriculumIDs []string
	for _, c := range content.All() {
		curriculumIDs = append(curriculumIDs, c.ID)
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

	next := stats.NextChallenge(m, curriculumIDs)
	fmt.Println()
	if next == "" {
		fmt.Println("Next Challenge:")
		fmt.Println("  Complete")
	} else {
		fmt.Println("Next Challenge:")
		fmt.Printf("  %s\n", next)
	}

	review := stats.RecommendedReview(m, curriculumIDs)
	if review != "" {
		fmt.Println()
		fmt.Println("Recommended Review:")
		fmt.Printf("  %s\n", review)
	}
}
