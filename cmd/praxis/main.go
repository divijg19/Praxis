package main

import (
	"fmt"
	"os"

	"github.com/divijg19/Praxis/internal/content"
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
