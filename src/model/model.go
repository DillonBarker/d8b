package model

import (
	"fmt"
	"io"
	"strings"

	"github.com/DillonBarker/d8b/src/db"
	"github.com/DillonBarker/d8b/src/queries"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle         = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	queryStyle        = lipgloss.NewStyle().PaddingLeft(6).Faint(true)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("40"))
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	errorStyle        = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("160"))
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("40"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = blurredStyle.Copy()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
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

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 64

		switch i {
		case 0:
			t.Placeholder = "Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Query"
			t.CharLimit = 64
		}

		m.inputs[i] = t
	}

	return m
}

func initialModelWithPlaceholder(name string, query string) model {
	m := model{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 64

		switch i {
		case 0:
			t.Placeholder = name
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = query
			t.CharLimit = 64
		}

		m.inputs[i] = t
	}

	return m
}

type model struct {
	table      table.Model
	list       list.Model
	keys       *listKeyMap
	choice     string
	quitting   bool
	inputting  bool
	editting   bool
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	error      error
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
		if m.error != nil {
			switch msg.String() {
			case "q", "b":
				m.error = nil
				return m, nil
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}
		}
		if m.inputting {
			switch msg.String() {
			case "ctrl+q":
				m.inputting = false
				return m, nil

			case "ctrl+c":
				return m, tea.Quit

			case "ctrl+r":
				m.cursorMode++
				if m.cursorMode > cursor.CursorHide {
					m.cursorMode = cursor.CursorBlink
				}
				cmds := make([]tea.Cmd, len(m.inputs))
				for i := range m.inputs {
					cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
				}
				return m, tea.Batch(cmds...)

			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				if s == "enter" && m.focusIndex == len(m.inputs) {
					if m.editting {
						i := m.list.Index()

						oldChoice := queries.Choice{Name: m.inputs[0].Placeholder, Query: m.inputs[1].Placeholder}
						choice := queries.Choice{Name: m.inputs[0].Value(), Query: m.inputs[1].Value()}

						queries.RemoveQuery(oldChoice)
						queries.AddQuery(choice)

						m.list.RemoveItem(i)
						m.list.InsertItem(len(m.list.Items()), Item{Name: choice.Name, Query: choice.Query})
					} else {
						choice := queries.Choice{Name: m.inputs[0].Value(), Query: m.inputs[1].Value()}

						queries.AddQuery(choice)

						m.list.InsertItem(len(m.list.Items()), Item{Name: choice.Name, Query: choice.Query})
					}

					m.inputting = false
					return m, nil
				}

				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}

				if m.focusIndex > len(m.inputs) {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs)
				}

				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusIndex {
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = focusedStyle
						m.inputs[i].TextStyle = focusedStyle
						continue
					}

					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}

				return m, tea.Batch(cmds...)
			}
		} else {
			switch {
			case key.Matches(msg, m.keys.AddNewQuery):
				m.inputs = initialModel().inputs
				m.inputting = true

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

			case "ctrl+d":
				item, ok := m.list.SelectedItem().(Item)
				i := m.list.Index()

				if ok {
					choice := queries.Choice{Name: item.Name, Query: item.Query}
					queries.RemoveQuery(choice)

					m.list.RemoveItem(i)

					return m, nil
				}

			case "e":
				item, ok := m.list.SelectedItem().(Item)

				if ok {
					m.inputs = initialModelWithPlaceholder(item.Name, item.Query).inputs
					m.editting = true
					m.inputting = true

					return m, nil
				}

			case "enter":
				if !m.table.Focused() {

					i, ok := m.list.SelectedItem().(Item)

					if ok {
						m.choice = string(i.Query)

						rowData, columnData, err := db.ExecuteQuery(m.choice)

						if err != nil {
							m.choice = ""
							m.error = err
						}

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
							Foreground(lipgloss.Color("240")).
							Background(lipgloss.Color("40")).
							Bold(true)
						t.SetStyles(s)

						m.table = t
					}
				}

			case "b":
				m.table.Blur()
				m.choice = ""
			}
		}
	}

	cmd := m.updateInputs(msg)
	m.list, _ = m.list.Update(msg)
	m.table, _ = m.table.Update(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	if m.error != nil {
		s := errorStyle.Render("Error: ", m.error.Error())
		s += helpStyle.Render("\n")
		s += helpStyle.Render("(q or b to go back)")
		return s
	}
	if m.inputting {
		var b strings.Builder

		for i := range m.inputs {
			b.WriteString(m.inputs[i].View())
			if i < len(m.inputs)-1 {
				b.WriteRune('\n')
			}
		}

		button := &blurredButton
		if m.focusIndex == len(m.inputs) {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)

		b.WriteString(helpStyle.Render("cursor mode is "))
		b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
		b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
		b.WriteString(helpStyle.Render(" (ctrl+q to go back)"))

		return b.String()
	}
	if m.choice != "" {
		s := helpStyle.Render("\n(b to go back)")

		return baseStyle.Render(m.table.View(), s)
	}
	if m.quitting {
		return quitTextStyle.Render("Fare thee well.")
	}
	return "\n" + m.list.View()
}

type listKeyMap struct {
	AddNewQuery      key.Binding
	EditQuery        key.Binding
	RemoveQuery      key.Binding
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
		EditQuery: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit query"),
		),
		RemoveQuery: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "remove query"),
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
