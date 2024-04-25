package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/lib/pq"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type item struct {
	title       string
	description string
}

func (i item) Titled() string      { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type model struct {
	showTable bool
	table     *table.Model
	list      *list.Model
}

type Config struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
}

func loadConfig() (Config, error) {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		return config, err
	}
	return config, nil
}

func executeQuery(db *sql.DB, query string) ([][]string, []string) {
	var result [][]string

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	values := make([]interface{}, len(columns))
	valuePointers := make([]interface{}, len(columns))
	for i := range values {
		valuePointers[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(valuePointers...)
		if err != nil {
			panic(err)
		}

		var row []string
		for _, v := range values {
			switch val := v.(type) {
			case nil:
				row = append(row, "NULL")
			case []byte:
				row = append(row, string(val))
			default:
				row = append(row, fmt.Sprintf("%v", val))
			}
		}
		result = append(result, row)
	}

	return result, columns
}

func connectToDb() (*sql.DB, error) {
	config, err := loadConfig()

	if err != nil {
		panic(err)
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.list != nil {
				selectedItem := m.list.SelectedItem()
				if selectedItem != nil {
					m.showTable = true
				}
			}
		case "b":
			m.showTable = false
		}
	}

	if m.showTable {
		newTableModel, cmd := m.table.Update(msg)
		m.table = &newTableModel
		return m, cmd
	} else {
		newListModel, cmd := m.list.Update(msg)
		m.list = &newListModel
		return m, cmd
	}
}

func (m model) View() string {
	if m.showTable {
		return baseStyle.Render(m.table.View()) + "\n"
	} else {
		return docStyle.Render(m.list.View())
	}
}

func main() {
	items := []list.Item{
		item{title: "List organisations", description: "SELECT * FROM organisation"},
		item{title: "List organisation subscriptions", description: "SELECT * FROM organisation_subscriptions"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)

	m := model{list: &l}

	m.list.Title = "Pick a command to run"

	database, err := connectToDb()

	if err != nil {
		panic(err)
	}

	rowData, columnData := executeQuery(database, "SELECT * FROM organisation")

	var columns []table.Column
	var rows []table.Row

	for _, row := range rowData {
		rows = append(rows, row)
	}

	for _, header := range columnData {
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

	m = model{false, &t, &l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
