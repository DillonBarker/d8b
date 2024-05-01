package main

import (
	"github.com/DillonBarker/d8b/src/table"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetText(text)
	}

	app := tview.NewApplication()

	tableList := table.GetTables()

	frame := tview.NewFrame(tableList)

	var inTable bool

	tableList.SetSelectedFunc(func(index int, tableName, secondaryText string, shortcut rune) {
		table := table.GetTable(tableName)

		frame.Clear()
		inTable = true
		frame.SetPrimitive(table)
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if inTable {
				frame.Clear()
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

	if err := app.SetRoot(grid, true).SetFocus(tableList).Run(); err != nil {
		panic(err)
	}
}
