package main

import (
	"fmt"
	"strings"

	"github.com/DillonBarker/d8b/src/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func getTableData() *tview.Table {
	query := "SELECT * from organisation"
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	rowData, headers, err := db.ExecuteQuery(query)

	if err != nil || headers == nil {
		panic(err)
	}

	for headerIndex, header := range headers {
		table.SetCell(0, headerIndex,
			tview.NewTableCell(strings.ToUpper(header)).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft))
	}

	for rowIndex, row := range rowData {
		for itemIndex, item := range row {
			color := tcell.ColorWhite

			table.SetCell(rowIndex+1, itemIndex,
				tview.NewTableCell(item).
					SetTextColor(color).
					SetAlign(tview.AlignLeft))
		}
	}

	table.
		SetTitle(fmt.Sprintf(" %s  [%d] ", query, len(rowData))).
		SetBorder(true)
	return table
}

func makeTable() *tview.Table {
	table := getTableData()

	return table
}

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	table := makeTable

	frame := tview.NewFrame(table())

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0, 1, 0).
		SetBorders(false).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(frame, 1, 0, 1, 3, 0, 0, true).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	if err := tview.NewApplication().EnableMouse(false).SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
