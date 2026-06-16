package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/divijg19/Praxis/internal/content"
)

var stageOrder = []string{
	"Movement",
	"Search",
	"Structural Navigation",
	"Editing",
	"Text Objects",
	"Registers",
}

func stageRank(stage string) int {
	for i, s := range stageOrder {
		if s == stage {
			return i
		}
	}
	return len(stageOrder)
}

func main() {
	layerOrder := []string{"Tutorial", "Training"}

	// Collect challenges grouped by layer then stage, preserving curriculum order
	type entry struct {
		d    content.Description
		rank int
	}
	grouped := make(map[string]map[string][]entry)

	for _, c := range content.All() {
		d, ok := content.DescriptionFor(c.ID)
		if !ok {
			fmt.Fprintf(os.Stderr, "error: DescriptionFor(%q) failed\n", c.ID)
			os.Exit(1)
		}
		if grouped[d.Layer] == nil {
			grouped[d.Layer] = make(map[string][]entry)
		}
		grouped[d.Layer][d.Stage] = append(grouped[d.Layer][d.Stage], entry{d, stageRank(d.Stage)})
	}

	// Count totals per layer
	layerCounts := make(map[string]int)
	for _, c := range content.All() {
		layerCounts[c.Layer]++
	}

	fmt.Println("# Challenge Catalog")
	fmt.Println()
	fmt.Printf("Total: **%d challenges** — ", len(content.All()))
	var parts []string
	for _, layer := range layerOrder {
		if n := layerCounts[layer]; n > 0 {
			parts = append(parts, fmt.Sprintf("**%d %s**", n, layer))
		}
	}
	fmt.Println(strings.Join(parts, " + ") + ".")
	fmt.Println()
	fmt.Println("Generated from `internal/content/describe.go` via `content.DescriptionFor`.")
	fmt.Println()

	for _, layer := range layerOrder {
		stages := grouped[layer]
		if len(stages) == 0 {
			continue
		}

		fmt.Printf("## %s\n\n", layer)

		// Render stages in pedagogical order
		for _, stage := range stageOrder {
			entries := stages[stage]
			if len(entries) == 0 {
				continue
			}

			fmt.Printf("### %s\n\n", stage)

			for _, e := range entries {
				d := e.d
				fmt.Printf("#### %s\n\n", d.Name)
				fmt.Printf("- **ID:** `%s`\n", d.ID)
				fmt.Printf("- **Verify:** `%s`\n", d.Verify)
				fmt.Printf("- **Layer:** `%s`\n", d.Layer)
				if d.Target != "" {
					fmt.Printf("- **Target:** `%s`\n", d.Target)
				}
				fmt.Printf("- **Primary Concept:** `%s`\n", d.Concept)
				fmt.Printf("- **Context:** `%s`\n", d.Context)
				fmt.Printf("- **Stage:** `%s`\n", d.Stage)
				if d.Evaluation != nil {
					fmt.Printf("- **Max Moves:** `%d`\n", d.Evaluation.MaxMoves)
				}
				fmt.Println()
				fmt.Println("#### Content")
				fmt.Println()
				fmt.Println("```text")
				for _, line := range d.Content {
					fmt.Println(line)
				}
				fmt.Println("```")
				if len(d.Result) > 0 {
					fmt.Println()
					fmt.Println("#### Result")
					fmt.Println()
					fmt.Println("```text")
					for _, line := range d.Result {
						fmt.Println(line)
					}
					fmt.Println("```")
				}
				fmt.Println()
				fmt.Println("---")
				fmt.Println()
			}
		}
	}
}
