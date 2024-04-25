package model

import (
	"fmt"

	"github.com/DillonBarker/d8b/src/db"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
			if m.cursor > 0 && len(m.output) == 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices.Choice)-1 && len(m.output) == 0 {
				m.cursor++
			}
		case "enter", " ":
			if len(m.output) == 0 {
				choice := m.choices.Choice[m.cursor].Name
				query, ok := m.queries[choice]
				if !ok {
					fmt.Printf("Error: Query for %s not found.\n", choice)
					return m, nil
				}

				m.output, m.headers = db.ExecuteQuery(database, query)
			}
		case "b":
			if len(m.output) > 0 {
				// Clear output and headers to return to options
				m.output = nil
				m.headers = nil
			}
		}
	}

	return m, nil
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

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
		var columns []table.Column
		var rows []table.Row

		// Prepare rows
		for _, row := range m.output {
			rows = append(rows, row)
		}

		// Prepare columns
		for _, header := range m.headers {
			columns = append(columns, table.Column{
				Title: header,
				Width: len(header),
			})
		}

		// Create the table
		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(7),
		)

		// Define styles
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		// Render the table and apply baseStyle
		return baseStyle.Render(t.View()) + "\n"
	}

	// If no data to display, just return the regular view
	s += "\nPress q to quit.\n"
	return s
}
