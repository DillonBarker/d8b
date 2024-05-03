package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/DillonBarker/d8b/internal/db"
	"github.com/DillonBarker/d8b/internal/table"
	"github.com/DillonBarker/d8b/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const baseTableColour = tcell.ColorMediumSpringGreen

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetText(text)
	}

	config, err := db.LoadConfig()

	if err != nil {
		panic(err)
	}

	header := ui.Header(config)

	app := tview.NewApplication()

	schemaListRows, _, err := table.GetSchemas(nil)

	if err != nil {
		panic(err)
	}

	schemaList := tview.NewList()
	tableList := tview.NewList()

	for _, tableName := range schemaListRows {
		schemaList.AddItem(tableName, "", 0, nil).
			SetMainTextColor(tview.Styles.SecondaryTextColor)
	}

	schemaList.
		SetBorder(true).
		SetTitleColor(baseTableColour).
		SetBorderColor(baseTableColour).
		SetBorderPadding(0, 0, 1, 1)
	schemaList.
		ShowSecondaryText(false).
		SetTitle(fmt.Sprintf(" schemas [%d] ", len(schemaListRows))).
		SetTitleColor(tview.Styles.PrimaryTextColor)

	frame := tview.NewFrame(schemaList)

	input := ui.Input()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '/' {
			app.SetFocus(input)
			input.SetLabel("> ")
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			input.SetText("")
			input.SetLabel("")
			updateList(schemaList, schemaListRows, nil)
		}
		return event
	})

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			input.SetText("")
			updateList(schemaList, schemaListRows, nil)
			app.SetFocus(frame)
		}
		if key == tcell.KeyEnter {
			app.SetFocus(frame)
			if app.GetFocus() == schemaList {
				filterText := input.GetText()
				updateList(schemaList, schemaListRows, &filterText)
			}
			if app.GetFocus() == tableList {
				filterText := input.GetText()
				title := tableList.GetTitle()
				re := regexp.MustCompile(`\(([^)]+)\)`)
				matches := re.FindStringSubmatch(title)

				var schemaName string
				if len(matches) > 1 {
					schemaName = matches[1]
				}

				tableListRows, err := table.GetTables(schemaName)

				if err != nil {
					panic(err)
				}

				updateList(tableList, tableListRows, &filterText)
			}
		}
	})

	schemaList.SetSelectedFunc(func(index int, schemaName, secondaryText string, shortcut rune) {
		tableListRows, err := table.GetTables(schemaName)

		if err != nil {
			panic(err)
		}

		for _, tableName := range tableListRows {
			tableList.AddItem(tableName, "", 0, nil).
				SetMainTextColor(tview.Styles.SecondaryTextColor)
		}

		tableList.
			ShowSecondaryText(false).
			SetTitle(fmt.Sprintf(" schema(%s) [%d] ", schemaName, len(tableListRows)))

		tableList.SetBorder(true).
			SetTitleColor(tview.Styles.PrimaryTextColor).
			SetBorderColor(baseTableColour).
			SetBorderPadding(0, 0, 1, 1)

		frame.SetPrimitive(tableList)

		input.SetText("")
		input.SetLabel("")

		app.SetFocus(frame)

		tableList.SetSelectedFunc(func(index int, tableName, secondaryText string, shortcut rune) {
			table := table.GetTable(schemaName, tableName)

			frame.SetPrimitive(table)

			table.SetDoneFunc(func(key tcell.Key) {
				switch key {
				case tcell.KeyEscape:
					frame.SetPrimitive(tableList)
				}
			})
		})

		tableList.SetDoneFunc(func() {
			tableList.Clear()
			frame.SetPrimitive(schemaList)
			app.SetFocus(frame)
		})
	})

	grid := tview.NewGrid().
		SetRows(4, 3, 0, 1).
		SetColumns(0, 1, 0).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(input, 1, 0, 1, 3, 0, 0, true).
		AddItem(frame, 2, 0, 1, 3, 0, 0, true).
		AddItem(newPrimitive("Footer"), 3, 0, 1, 3, 0, 0, false)

	if err := app.SetRoot(grid, true).SetFocus(schemaList).Run(); err != nil {
		panic(err)
	}
}

func filterData(data []string, filter string) []string {
	filteredData := make([]string, 0)
	filter = strings.ToLower(filter)
	for _, item := range data {
		if strings.Contains(strings.ToLower(item), filter) {
			filteredData = append(filteredData, item)
		}
	}
	return filteredData
}

func updateList(list *tview.List, data []string, filter *string) {
	list.Clear()
	var d []string

	if filter != nil {
		f := *filter
		d = filterData(data, f)
	} else {
		d = data
	}

	for _, item := range d {
		list.AddItem(item, "", 0, nil)
	}
}
