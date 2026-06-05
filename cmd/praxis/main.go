package main

import (
	"fmt"
	"os"

	"github.com/divijg19/Praxis/internal/content"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "list" {
		list()
		return
	}
	fmt.Println("Praxis")
	fmt.Println("Mastery through practice.")
}

func list() {
	for _, c := range content.All() {
		fmt.Println(c.Name)
	}
}
