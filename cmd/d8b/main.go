package main

import (
	"github.com/DillonBarker/d8b/internal/table"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetText(text)
	}

	app := tview.NewApplication()

	schemaList := table.GetSchemas()
	tableList := tview.NewList()

	frame := tview.NewFrame(schemaList)

	var inTables bool
	var inTable bool

	schemaList.SetSelectedFunc(func(index int, schemaName, secondaryText string, shortcut rune) {
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
		}
		return event
	})

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0, 1, 0).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(frame, 1, 0, 1, 3, 0, 0, true).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	if err := app.SetRoot(grid, true).SetFocus(schemaList).Run(); err != nil {
		panic(err)
	}
}
