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
	table    table.Model
}

func InitialModel(choices Choices, queries map[string]string) model {
	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
		headers:  nil,
		output:   nil,
		queries:  queries,
		table:    table.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
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

		for _, row := range m.output {
			rows = append(rows, row)
		}

		for _, header := range m.headers {
			columns = append(columns, table.Column{
				Title: header,
				Width: len(header) + 5,
			})
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(10),
		)

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

		tableView := baseStyle.Render(t.View()) + "\n"

		return tableView
	}

	s += "\nPress q to quit.\n"
	return s
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
			if len(m.output) > 0 && m.cursor > 0 {
				m.cursor--
			} else if len(m.output) == 0 && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if len(m.output) > 0 && m.cursor < len(m.output)-1 {
				m.cursor++
			} else if len(m.output) == 0 && m.cursor < len(m.choices.Choice)-1 {
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

				m.table.Focus()
			}
		case "b":
			if len(m.output) > 0 {
				m.output = nil
				m.headers = nil
			}
		}
	}

	if m.table.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				m.table.MoveUp(1)
				print(m.table.SelectedRow())
			case "down", "j":
				m.table.MoveDown(1)
				print(m.table.SelectedRow())
			}

		}
		return m, nil
	}

	return m, nil
}
