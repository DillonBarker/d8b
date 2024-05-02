package main

import (
	"strings"

	"github.com/DillonBarker/d8b/internal/db"
	"github.com/DillonBarker/d8b/internal/table"
	"github.com/DillonBarker/d8b/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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

	schemaList := table.GetSchemas(nil)
	schemaListFiltered := tview.NewList()
	tableList := tview.NewList()

	frame := tview.NewFrame(schemaList)

	input := ui.Input()

	var inTables bool
	var inTable bool
	var inInput bool
	var isFiltering bool

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if inTables {
				frame.Clear()
				inTables = false
				inTable = false
				frame.SetPrimitive(schemaList)
			}
			if inTable {
				frame.Clear()
				inTables = true
				frame.SetPrimitive(tableList)
			}
			if isFiltering {
				frame.Clear()
				inInput = false
				input.SetText("")
				frame.SetPrimitive(schemaList)
			}
			if inInput {
				inInput = false
				input.SetText("")
				frame.SetPrimitive(schemaList) // or another basic list
				app.SetFocus(frame)
			}
		case tcell.KeyEnter:
			if inInput {
				filterText := strings.Replace(input.GetText(), "/", "", 1)

				schemaListFiltered = table.GetSchemas(&filterText)

				schemaList.FindItems("", "", false, false)

				frame.SetPrimitive(schemaListFiltered)
				app.SetFocus(frame)

				inInput = false

				isFiltering = true
			} else if isFiltering {
				index := schemaListFiltered.GetCurrentItem()
				text, _ := schemaListFiltered.GetItemText(index)

				tableList = table.GetTables(text)
				frame.SetPrimitive(tableList)
				app.SetFocus(frame)

				inTables = true
				input.SetText("")
				isFiltering = false
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case '/':
				inInput = true
				input.SetLabel("filter: ")
				input.SetText("")
				app.SetFocus(input)
			}
		}

		return event
	})

	schemaList.SetSelectedFunc(func(index int, schemaName, secondaryText string, shortcut rune) {
		if inInput {
			return
		}

		tableList = table.GetTables(schemaName)

		frame.Clear()
		inTables = true
		inTable = false
		frame.SetPrimitive(tableList)

		tableList.SetSelectedFunc(func(index int, tableName, secondaryText string, shortcut rune) {
			table := table.GetTable(schemaName, tableName)

			frame.Clear()
			inTables = false
			inTable = true
			frame.SetPrimitive(table)
		})
	})

	grid := tview.NewGrid().
		SetRows(4, 3, 0, 1).
		SetColumns(0, 1, 0).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(input, 1, 0, 1, 3, 0, 0, false).
		AddItem(frame, 2, 0, 1, 3, 0, 0, true).
		AddItem(newPrimitive("Footer"), 3, 0, 1, 3, 0, 0, false)

	if err := app.SetRoot(grid, true).SetFocus(schemaList).Run(); err != nil {
		panic(err)
	}
}
