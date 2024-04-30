package main

import (
	"fmt"
	"math"

	"github.com/rivo/tview"
)

type TableData struct {
	tview.TableContentReadOnly
}

func (d *TableData) GetCell(row, column int) *tview.TableCell {
	letters := [...]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 'A' + byte(row%26)} // log(math.MaxInt64) / log(26) ~= 14
	start := len(letters) - 1
	row /= 26
	for row > 0 {
		start--
		row--
		letters[start] = 'A' + byte(row%26)
		row /= 26
	}
	return tview.NewTableCell(fmt.Sprintf("[red]%s[green]%d", letters[start:], column))
}

func (d *TableData) GetRowCount() int {
	return math.MaxInt64
}

func (d *TableData) GetColumnCount() int {
	return math.MaxInt64
}

func makeTable() *tview.Table {
	data := &TableData{}
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetContent(data)

	table.SetBorder(false)

	table.SetTitle(" table x y z ")
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
