package model

import (
	"fmt"
	"strings"

	"github.com/DillonBarker/d8b/src/db"
	tea "github.com/charmbracelet/bubbletea"
)

type Choice struct {
	Name  string `toml:"name"`
	Query string `toml:"query"`
}

type Choices struct {
	Choice []Choice `toml:"choice"`
}

type model struct {
	choices  Choices
	cursor   int
	selected map[int]struct{}
	output   [][]string
	headers  []string
	queries  map[string]string
}

func InitialModel(choices Choices, queries map[string]string) model {
	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
		headers:  nil,
		output:   nil,
		queries:  queries,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	database, err := db.ConnectToDb()

	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices.Choice)-1 {
				m.cursor++
			}
		case "enter", " ":
			choice := m.choices.Choice[m.cursor].Name
			query, ok := m.queries[choice]
			if !ok {
				fmt.Printf("Error: Query for %s not found.\n", choice)
				return m, nil
			}

			if choice == "Create New Entry" {

			} else {
				m.output, m.headers = db.ExecuteQuery(database, query)
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "What database command to run?\n\n"

	for i, choice := range m.choices.Choice {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice.Name)
	}
	s += "\n"

	if len(m.output) > 0 {
		colWidths := make([]int, len(m.headers))

		for i, header := range m.headers {
			colWidths[i] = len(header)
		}
		for _, row := range m.output {
			for i, val := range row {
				if len(val) > colWidths[i] {
					colWidths[i] = len(val)
				}
			}
		}

		for i, header := range m.headers {
			s += fmt.Sprintf("| %-*s ", colWidths[i], header)
		}
		s += "|\n"

		for i, width := range colWidths {
			if i == 0 {
				s += "+"
			}
			s += strings.Repeat("-", width+2)
			if i < len(colWidths)-1 {
				s += "+"
			}
		}
		s += "+\n"

		for _, row := range m.output {
			for i, val := range row {
				s += fmt.Sprintf("| %-*s ", colWidths[i], val)
			}
			s += "|\n"
		}
	}

	s += "\nPress q to quit.\n"

	return s
}
