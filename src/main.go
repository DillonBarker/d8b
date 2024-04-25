package main

import (
	"fmt"
	"os"

	"github.com/DillonBarker/d8b/src/model"
	"github.com/DillonBarker/d8b/src/query"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/lib/pq"
)

func main() {
	queries, choices, err := query.LoadQueries()
	if err != nil {
		fmt.Printf("Error loading queries: %v\n", err)
		os.Exit(1)
	}

	initial := model.InitialModel(choices, queries)

	p := tea.NewProgram(initial)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
