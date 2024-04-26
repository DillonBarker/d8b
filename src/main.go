package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/lib/pq"
)

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

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	queryStyle        = lipgloss.NewStyle().PaddingLeft(6).Faint(true)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("40"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	name  string
	query string
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	nameStr := fmt.Sprintf("%d. %s", index+1, i.name)
	nameFn := itemStyle.Render
	if index == m.Index() {
		nameFn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, nameFn(nameStr))

	queryStr := i.query
	queryFn := queryStyle.Render
	fmt.Fprintln(w)
	fmt.Fprint(w, queryFn(queryStr))
}

type model struct {
	table    table.Model
	list     list.Model
	keys     *listKeyMap
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.query)

				database, err := connectToDb()

				if err != nil {
					panic(err)
				}

				rowData, columnData := executeQuery(database, m.choice)

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

				m.table = t
			}

		case "b":
			m.choice = ""
		}
	}

	var cmd tea.Cmd
	m.list, _ = m.list.Update(msg)
	m.table, _ = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return baseStyle.Render(m.table.View())
	}
	if m.quitting {
		return quitTextStyle.Render("Fare thee well.")
	}
	return "\n" + m.list.View()
}

type Choice struct {
	Name  string `toml:"name"`
	Query string `toml:"query"`
}

type Choices struct {
	Choices []Choice `toml:"choice"`
}

func loadQueries() (Choices, error) {
	var choices Choices
	if _, err := toml.DecodeFile("queries.toml", &choices); err != nil {
		return choices, err
	}
	return choices, nil
}

func main() {
	choices, err := loadQueries()
	listKeys := newListKeyMap()

	if err != nil {
		panic(err)
	}

	var items []list.Item
	for _, choice := range choices.Choices {
		items = append(items, item{name: choice.Name, query: choice.Query})
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What query do you want to run?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}
	l.SetShowFilter(true)
	l.SetFilteringEnabled(true)

	m := model{list: l, keys: listKeys}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type listKeyMap struct {
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}
