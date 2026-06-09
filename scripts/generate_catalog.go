package main

import (
	"fmt"
	"os"

	"github.com/divijg19/Praxis/internal/challenge"
	"github.com/divijg19/Praxis/internal/content"
)

type pack struct {
	Name string
	IDs  []string
}

var packs = []pack{
	{"Movement", []string{"motion_rush", "grid_rush"}},
	{"Search", []string{"find_hunter", "word_hunter", "symbol_hunter", "line_hunter"}},
	{"Structural Navigation", []string{
		"paren_hunter", "sentence_hunter", "slash_hunter",
		"question_hunter", "repeat_hunter",
		"inner_paren_hunter", "around_paren_hunter",
		"inner_bracket_hunter", "around_bracket_hunter",
		"inner_quote_hunter", "around_quote_hunter",
		"paragraph_hunter", "match_hunter",
	}},
	{"Editing", []string{
		"delete_character_hunter", "replace_character_hunter",
		"toggle_case_hunter", "delete_word_hunter", "change_word_hunter",
	}},
	{"UTF-8 Proof", []string{"utf8_cursor_hunter"}},
	{"Structural Editing", []string{
		"delete_line_hunter", "delete_to_end_hunter",
		"delete_inner_word_hunter", "delete_around_word_hunter",
		"delete_inner_paren_hunter", "delete_around_paren_hunter",
		"delete_inner_quote_hunter", "delete_around_quote_hunter",
		"change_inner_word_hunter", "change_inner_paren_hunter",
		"change_inner_quote_hunter", "yank_line_hunter",
	}},
	{"Registers", []string{
		"named_register_hunter", "word_register_hunter",
		"register_replace_hunter", "register_duplicate_hunter",
	}},
}

func main() {
	all := make(map[string]challenge.Challenge)
	for _, c := range content.All() {
		all[c.ID] = c
	}

	var count int
	for _, p := range packs {
		for _, id := range p.IDs {
			if _, ok := all[id]; !ok {
				fmt.Fprintf(os.Stderr, "error: challenge %q not found in content\n", id)
				os.Exit(1)
			}
		}
		count += len(p.IDs)
	}

	if count != len(all) {
		fmt.Fprintf(os.Stderr, "error: catalog covers %d challenges but content has %d\n", count, len(all))
		os.Exit(1)
	}

	fmt.Println("# Challenge Catalog")
	fmt.Println()
	fmt.Printf("Total: **%d challenges** across **%d curriculum packs**.\n", count, len(packs))
	fmt.Println()
	fmt.Println("Generated from `internal/content/content.go`. Do not edit by hand.")
	fmt.Println()

	for _, p := range packs {
		fmt.Printf("## %s\n\n", p.Name)
		for _, id := range p.IDs {
			c := all[id]
			fmt.Printf("### %s\n\n", c.Name)
			fmt.Printf("- **ID:** `%s`  \n", c.ID)
			fmt.Printf("- **Verify:** `%s`  \n", c.Verify)
			if c.Target != "" {
				fmt.Printf("- **Target:** `%s`  \n", c.Target)
			}
			fmt.Println()
			fmt.Println("#### Content")
			fmt.Println()
			fmt.Println("```text")
			for _, line := range c.Content {
				fmt.Println(line)
			}
			fmt.Println("```")
			if len(c.Result) > 0 {
				fmt.Println()
				fmt.Println("#### Result")
				fmt.Println()
				fmt.Println("```text")
				for _, line := range c.Result {
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
