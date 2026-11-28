package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/divijg19/Praxis/internal/challenge"
	"github.com/divijg19/Praxis/internal/content"
	"github.com/divijg19/Praxis/internal/score"
	"github.com/divijg19/Praxis/internal/stats"
	"github.com/divijg19/Praxis/internal/verify"
)

func usage(msg string) {
	fmt.Fprintln(os.Stderr, "praxis: usage:", msg)
	os.Exit(1)
}

func unknown(id string) {
	fmt.Fprintln(os.Stderr, "unknown challenge:", id)
	os.Exit(1)
}

func banner() {
	fmt.Println("Praxis")
	fmt.Println("Mastery through practice.")
}

func corruptProgress() {
	fmt.Fprintln(os.Stderr, "praxis: progress file is corrupt")
	fmt.Println("Your Praxis progress is corrupted.")
	fmt.Println("A backup was saved next to it. Run `praxis reset` to start fresh.")
}

func loadOrFail() stats.Progress {
	p, err := stats.Load()
	if err != nil {
		if errors.Is(err, stats.ErrCorrupt) || errors.Is(err, stats.ErrInvalidStats) {
			corruptProgress()
		} else {
			fmt.Fprintln(os.Stderr, "praxis: could not read progress:", err)
		}
		os.Exit(1)
	}
	return p
}

func mutateOrFail(fn func(*stats.Progress) error) {
	if err := stats.Mutate(fn); err != nil {
		if errors.Is(err, stats.ErrCorrupt) || errors.Is(err, stats.ErrInvalidStats) {
			corruptProgress()
		} else {
			fmt.Fprintln(os.Stderr, "praxis: could not update progress:", err)
		}
		os.Exit(1)
	}
}

func devOnly() {
	if os.Getenv("PRAXIS_DEV") != "1" {
		fmt.Fprintln(os.Stderr, "praxis: development command; set PRAXIS_DEV=1")
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
			if os.Args[2] != "--ids" && os.Args[2] != "--core" && os.Args[2] != "--json" {
				usage("praxis catalog [--ids|--core|--json]")
			}
			catalog(os.Args[2])
		default:
			usage("praxis catalog [--ids|--core|--json]")
		}
	case "describe":
		if len(os.Args) != 3 {
			usage("praxis describe <id>")
		}
		describe(os.Args[2])
	case "trace":
		devOnly()
		if len(os.Args) != 3 {
			usage("praxis trace <id>")
		}
		trace(os.Args[2])
	case "start":
		if len(os.Args) != 3 {
			usage("praxis start <id>")
		}
		start(os.Args[2])
	case "hint":
		if len(os.Args) != 4 && len(os.Args) != 5 {
			usage("praxis hint <id> <level> [--preview]")
		}
		if len(os.Args) == 5 && os.Args[4] != "--preview" {
			usage("praxis hint <id> <level> [--preview]")
		}
		hint(os.Args[2], os.Args[3], len(os.Args) == 5)
	case "complete":
		if len(os.Args) != 3 {
			usage("praxis complete <id>")
		}
		complete(os.Args[2], true)
	case "verify":
		if len(os.Args) != 3 {
			usage("praxis verify <id>")
		}
		complete(os.Args[2], false)
	case "dev":
		devOnly()
		if len(os.Args) < 4 || os.Args[2] != "mark-complete" {
			usage("praxis dev mark-complete <id> [count]")
		}
		devMarkComplete(os.Args[3])
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
	case "report":
		if len(os.Args) != 2 && (len(os.Args) != 3 || os.Args[2] != "--json") {
			usage("praxis report [--json]")
		}
		report(len(os.Args) == 3)
	case "stats":
		switch len(os.Args) {
		case 2:
			statsSummary()
		case 3:
			if os.Args[2] == "--json" {
				statsJSON()
			} else {
				statsForID(os.Args[2])
			}
		default:
			usage("praxis stats [<id>|--json]")
		}
	default:
		usage("unknown command: " + os.Args[1])
	}
}

func tutorialProgress(p stats.Progress) (done, total int) {
	for _, id := range content.CoreTutorialIDs() {
		total++
		if p.Entries[id].Completions > 0 {
			done++
		}
	}
	return done, total
}

func entrypoint() {
	banner()
	fmt.Println()
	p := loadOrFail()
	id := stats.NextChallenge(p, content.IDs())
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
	fmt.Println("  praxis start <id>")
	fmt.Println()
	fmt.Println("Practice:")
	fmt.Println("  praxis complete <id>   (reads a verified key log on stdin)")
	fmt.Println("  praxis verify <id>     (replay only, no progress write)")
	fmt.Println("  praxis hint <id> <level> [--preview]")
	fmt.Println()
	fmt.Println("Track:")
	fmt.Println("  praxis stats")
	fmt.Println("  praxis report [--json]")
	fmt.Println("  praxis catalog [--ids|--core|--json]")
	fmt.Println("  praxis reset")
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

// hintCostFor derives the cost from persisted reveal state only; stdin is never
// trusted for hint cost, level, or marks.
func hintCostFor(id string, p stats.Progress) int {
	return content.HintCostThrough(id, stats.RevealedHintLevel(p, id))
}

func hint(id, levelArg string, preview bool) {
	if !content.Exists(id) {
		unknown(id)
	}
	level, err := strconv.Atoi(levelArg)
	if err != nil || level < 1 || level > content.HintCount(id) {
		usage("praxis hint <id> <level> [--preview]")
	}
	text := content.HintText(id, level)
	cost := content.HintCostThrough(id, level)
	if preview {
		fmt.Printf("Hint %d of %d costs +%d key-equivalents (total +%d).\n", level, content.HintCount(id), content.HintCost(id, level), cost)
		return
	}
	mutateOrFail(func(p *stats.Progress) error {
		stats.RevealHint(p, id, level)
		return nil
	})
	fmt.Println(text)
	fmt.Printf("Cost: +%d key-equivalents (total +%d).\n", content.HintCost(id, level), cost)
}

func trace(id string) {
	if !content.Exists(id) {
		unknown(id)
	}
	t, ok := content.TraceFor(id)
	if !ok {
		fmt.Fprintln(os.Stderr, "no canonical trace for challenge:", id)
		os.Exit(1)
	}
	if err := json.NewEncoder(os.Stdout).Encode(t); err != nil {
		fmt.Fprintln(os.Stderr, "trace failed:", err)
		os.Exit(1)
	}
}

func catalog(flag string) {
	if flag == "--json" {
		type row struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Layer   string `json:"layer"`
			Tier    string `json:"tier"`
			Stage   string `json:"stage"`
			Concept string `json:"concept"`
			Par     int    `json:"par"`
		}
		rows := make([]row, 0, len(content.All()))
		for _, c := range content.All() {
			d, _ := content.DescriptionFor(c.ID)
			rows = append(rows, row{
				ID: c.ID, Name: c.Name, Layer: c.Layer, Tier: content.Tier(c.ID),
				Stage: d.Stage, Concept: d.Concept, Par: d.Par,
			})
		}
		if err := json.NewEncoder(os.Stdout).Encode(rows); err != nil {
			fmt.Fprintln(os.Stderr, "catalog failed:", err)
			os.Exit(1)
		}
		return
	}
	if flag == "--core" {
		for _, id := range content.CoreTutorialIDs() {
			fmt.Println(id)
		}
		return
	}
	for _, c := range content.All() {
		if flag == "--ids" {
			fmt.Println(c.ID)
		} else {
			fmt.Println(c.Name)
		}
	}
}

func start(id string) {
	if !content.Exists(id) {
		unknown(id)
	}
	mutateOrFail(func(p *stats.Progress) error {
		stats.Start(p, id)
		return nil
	})
}

func complete(id string, write bool) {
	if !content.Exists(id) {
		unknown(id)
	}
	var in verify.Input
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		verb := "complete"
		if !write {
			verb = "verify"
		}
		usage("praxis " + verb + " <id> (invalid JSON on stdin)")
	}
	result, err := verify.Run(id, in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "replay failed:", err)
		os.Exit(1)
	}
	out := struct {
		Verified    bool   `json:"verified"`
		Reason      string `json:"reason"`
		Keys        int    `json:"keys"`
		Edits       int    `json:"edits"`
		HintCost    int    `json:"hint_cost"`
		Par         int    `json:"par"`
		Marks       int    `json:"marks"`
		Completions int    `json:"completions"`
	}{
		Verified: result.Verified, Reason: result.Reason, Keys: result.Keys,
		Edits: result.Edits, Par: result.Par,
	}
	hintCost := 0
	if result.Verified && !write {
		p := loadOrFail()
		hintCost = hintCostFor(id, p)
		out.Marks = score.Marks(result.Keys, result.Par, hintCost)
		out.HintCost = hintCost
	}
	if result.Verified && write {
		mutateOrFail(func(p *stats.Progress) error {
			hintCost = hintCostFor(id, *p)
			s, marks, err := stats.Complete(p, stats.Event{
				ID: id, Keys: result.Keys, Par: result.Par,
				Edits: result.Edits, HintCost: hintCost, Verified: true,
			})
			if err == nil {
				out.Marks = marks
				out.Completions = s.Completions
				out.HintCost = hintCost
			}
			return err
		})
	}
	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		fmt.Fprintln(os.Stderr, "complete failed:", err)
		os.Exit(1)
	}
}

func devMarkComplete(id string) {
	if !content.Exists(id) {
		unknown(id)
	}
	count := 3
	if len(os.Args) > 5 {
		usage("praxis dev mark-complete <id> [count]")
	}
	if len(os.Args) == 5 {
		parsed, err := strconv.Atoi(os.Args[4])
		count = parsed
		if err != nil || count < 1 {
			usage("praxis dev mark-complete <id> [count]")
		}
	}
	mutateOrFail(func(p *stats.Progress) error {
		stats.Start(p, id)
		for i := 0; i < count; i++ {
			if _, _, err := stats.Complete(p, stats.Event{
				ID: id, Keys: 1, Par: 1, Edits: 1, Verified: true,
			}); err != nil {
				return err
			}
		}
		return nil
	})
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
		if _, err := fmt.Scanln(&input); err != nil || input != "RESET" {
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

func statsForID(id string) {
	if !content.Exists(id) {
		unknown(id)
	}
	s := loadOrFail().Entries[id]
	fmt.Printf("Starts: %d\n", s.Starts)
	fmt.Printf("Completions: %d\n", s.Completions)
	fmt.Printf("Best Keys: %d\n", s.BestKeys)
	fmt.Printf("Best Edits: %d\n", s.BestEdits)
	fmt.Printf("Best Marks: %d/3\n", s.BestMarks)
	fmt.Printf("Mastery: %s\n", stats.MasteryTier(s))
}

func statsJSON() {
	p := loadOrFail()
	if err := json.NewEncoder(os.Stdout).Encode(p); err != nil {
		fmt.Fprintln(os.Stderr, "stats failed:", err)
		os.Exit(1)
	}
}

func statsSummary() {
	p := loadOrFail()
	ids := content.IDs()
	var completed, totalStarts int
	for _, id := range ids {
		s, ok := p.Entries[id]
		if !ok {
			continue
		}
		if s.Completions > 0 {
			completed++
		}
		totalStarts += s.Starts
	}
	fmt.Printf("Challenges Completed: %d/%d\n", completed, len(ids))
	fmt.Printf("Total Starts: %d\n", totalStarts)
	fmt.Println()
	coreDone, coreTotal := tutorialProgress(p)
	if coreDone >= coreTotal {
		fmt.Println("Tutorial: Complete")
	} else {
		fmt.Printf("Tutorial: %d/%d core\n", coreDone, coreTotal)
	}
	fmt.Println()
	dist := stats.MasteryDistribution(p, len(ids))
	fmt.Println("Mastery:")
	fmt.Printf("  Unseen: %d\n", dist["Unseen"])
	fmt.Printf("  Learning: %d\n", dist["Learning"])
	fmt.Printf("  Practiced: %d\n", dist["Practiced"])
	fmt.Printf("  Experienced: %d\n", dist["Experienced"])
	fmt.Println()
	next := stats.NextChallenge(p, ids)
	fmt.Println("Next:")
	if next != "" {
		d, _ := content.DescriptionFor(next)
		fmt.Printf("  %s\n", d.Name)
	}
	review := stats.RecommendedReview(p, ids)
	if review != "" {
		fmt.Println()
		fmt.Println("Recommended Review:")
		fmt.Printf("  %s\n", review)
	}
}

// A report is a private, descriptive view over verified runs. It adds no
// score, ranking, comparison, or progress state, and never gates practice.
type reportRun struct {
	ID           string `json:"id"`
	VerifiedKeys int    `json:"verified_keys"`
	Par          int    `json:"par"`
	Marks        int    `json:"marks"`
	HintCost     int    `json:"hint_cost"`
}

type reportData struct {
	CompletionsByLayer map[string]int `json:"completions_by_layer"`
	Runs               int            `json:"runs"`
	Marks              map[string]int `json:"marks"`
	HintRuns           int            `json:"hint_runs"`
	HintKeyEquivalents int            `json:"hint_key_equivalents"`
	Recent             []reportRun    `json:"recent"`
}

const recentRuns = 5

func buildReport(p stats.Progress) reportData {
	data := reportData{
		CompletionsByLayer: map[string]int{
			challenge.LayerTutorial: 0,
			challenge.LayerTraining: 0,
			challenge.LayerTrial:    0,
		},
		Marks:  map[string]int{"1": 0, "2": 0, "3": 0},
		Recent: []reportRun{},
	}
	for _, id := range content.IDs() {
		if layer := content.LayerFor(id); layer != "" {
			data.CompletionsByLayer[layer] += p.Entries[id].Completions
		}
	}
	for _, e := range p.Events {
		if !e.Verified {
			continue
		}
		data.Runs++
		if e.Marks >= 1 && e.Marks <= 3 {
			data.Marks[strconv.Itoa(e.Marks)]++
		}
		if e.HintCost > 0 {
			data.HintRuns++
			data.HintKeyEquivalents += e.HintCost
		}
	}
	start := len(p.Events) - recentRuns
	if start < 0 {
		start = 0
	}
	for _, e := range p.Events[start:] {
		data.Recent = append(data.Recent, reportRun{
			ID: e.ID, VerifiedKeys: e.Keys, Par: e.Par,
			Marks: e.Marks, HintCost: e.HintCost,
		})
	}
	return data
}

func report(jsonOutput bool) {
	data := buildReport(loadOrFail())
	if jsonOutput {
		if err := json.NewEncoder(os.Stdout).Encode(data); err != nil {
			fmt.Fprintln(os.Stderr, "report failed:", err)
			os.Exit(1)
		}
		return
	}
	fmt.Println("Private practice report")
	fmt.Println()
	fmt.Println("Completions by layer:")
	for _, layer := range []string{challenge.LayerTutorial, challenge.LayerTraining, challenge.LayerTrial} {
		fmt.Printf("  %s: %d\n", layer, data.CompletionsByLayer[layer])
	}
	fmt.Println()
	fmt.Println("Verified runs:")
	fmt.Printf("  Total: %d\n", data.Runs)
	fmt.Printf("  3 marks: %d   2 marks: %d   1 mark: %d\n", data.Marks["3"], data.Marks["2"], data.Marks["1"])
	fmt.Printf("  Runs with hints: %d\n", data.HintRuns)
	fmt.Printf("  Hint key-equivalents: %d\n", data.HintKeyEquivalents)
	fmt.Println()
	fmt.Println("Recent verified runs:")
	if len(data.Recent) == 0 {
		fmt.Println("  none")
		return
	}
	for _, r := range data.Recent {
		fmt.Printf("  %s: %d/%d keys, %d marks, +%d hint\n", r.ID, r.VerifiedKeys, r.Par, r.Marks, r.HintCost)
	}
}
