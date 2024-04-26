package model

import (
	"fmt"
	"io"
	"strings"

	"github.com/DillonBarker/d8b/src/db"
	"github.com/DillonBarker/d8b/src/queries"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle         = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	queryStyle        = lipgloss.NewStyle().PaddingLeft(6).Faint(true)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("40"))
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type Item struct {
	Name  string
	Query string
}

func (i Item) FilterValue() string { return "" }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	nameStr := fmt.Sprintf("%d. %s", index+1, i.Name)
	nameFn := itemStyle.Render
	if index == m.Index() {
		nameFn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, nameFn(nameStr))

	queryStr := i.Query
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

func Model(l list.Model, listKeys *listKeyMap) model {
	return model{list: l, keys: listKeys}
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
		case key.Matches(msg, m.keys.AddNewQuery):
			choice := queries.Choice{Name: "test", Query: "test"}
			queries.AddQuery(choice)

			m.list.InsertItem(len(m.list.Items()), Item{Name: choice.Name, Query: choice.Query})

			return m, nil

		case key.Matches(msg, m.keys.ToggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			return m, nil

		case key.Matches(msg, m.keys.ToggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.TogglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.ToggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}

		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(Item)

			if ok {
				m.choice = string(i.Query)

				rowData, columnData := db.ExecuteQuery(m.choice)

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
					Bold(true)
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

type listKeyMap struct {
	AddNewQuery      key.Binding
	ToggleTitleBar   key.Binding
	ToggleStatusBar  key.Binding
	TogglePagination key.Binding
	ToggleHelpMenu   key.Binding
}

func NewListKeyMap() *listKeyMap {
	return &listKeyMap{
		AddNewQuery: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "add new query"),
		),
		ToggleTitleBar: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle title"),
		),
		ToggleStatusBar: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle status"),
		),
		TogglePagination: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "toggle pagination"),
		),
		ToggleHelpMenu: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "toggle help"),
		),
	}
}
