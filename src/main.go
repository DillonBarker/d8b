package main

import (
	"fmt"
	"os"

	"github.com/DillonBarker/d8b/src/model"
	"github.com/DillonBarker/d8b/src/queries"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	choices, err := queries.LoadQueries()
	listKeys := model.NewListKeyMap()

	if err != nil {
		panic(err)
	}

	var items []list.Item
	for _, choice := range choices.Choices {
		items = append(items, model.Item{Name: choice.Name, Query: choice.Query})
	}

	var (
		defaultWidth    = 20
		listHeight      = 14
		titleStyle      = lipgloss.NewStyle().MarginLeft(2)
		paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
		helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	)

	l := list.New(items, model.ItemDelegate{}, defaultWidth, listHeight)
	l.Title = "What query do you want to run?"
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.AddNewQuery,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.ToggleTitleBar,
			listKeys.ToggleStatusBar,
			listKeys.TogglePagination,
			listKeys.ToggleHelpMenu,
		}
	}
	l.SetShowFilter(true)
	l.SetFilteringEnabled(true)

	m := model.Model(l, listKeys)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
