package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
}

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

func loadConfig() (Config, error) {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		return config, err
	}
	return config, nil
}

func loadQueries() (map[string]string, Choices, error) {
	var queries map[string]string
	var choices Choices

	if _, err := toml.DecodeFile("queries.toml", &choices); err != nil {
		return queries, choices, err
	}

	queryMap := make(map[string]string)
	for _, choice := range choices.Choice {
		queryMap[choice.Name] = choice.Query
	}

	return queryMap, choices, nil
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

func doQuery(db *sql.DB, query string) ([][]string, []string) {
	var result [][]string

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Get the column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	// Create a slice to hold the values of one row
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

func initialModel(choices Choices, queries map[string]string) model {
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
	db, err := connectToDb()
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
			m.output, m.headers = doQuery(db, query)
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
		// Calculate column widths
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

		// Print headers
		for i, header := range m.headers {
			s += fmt.Sprintf("| %-*s ", colWidths[i], header)
		}
		s += "|\n"
		// Print separator
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
		// Print rows
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

func main() {
	queries, choices, err := loadQueries()
	if err != nil {
		fmt.Printf("Error loading queries: %v\n", err)
		os.Exit(1)
	}

	initial := initialModel(choices, queries)

	p := tea.NewProgram(initial)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
