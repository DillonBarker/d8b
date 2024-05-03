package main

import (
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

	frame := tview.NewFrame(schemaList)

	schemaList.SetSelectedFunc(func(index int, schemaName, secondaryText string, shortcut rune) {
		tableList := table.GetTables(schemaName)

		frame.SetPrimitive(tableList)

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
			frame.SetPrimitive(schemaList)
		})

	})

	grid := tview.NewGrid().
		SetRows(4, 0, 1).
		SetColumns(0, 1, 0).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(frame, 1, 0, 1, 3, 0, 0, true).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	if err := app.SetRoot(grid, true).SetFocus(schemaList).Run(); err != nil {
		panic(err)
	}
}
