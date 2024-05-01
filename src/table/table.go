package table

import (
	"fmt"
	"strings"

	"github.com/DillonBarker/d8b/src/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	tableHeaderColour = tcell.ColorWhite
	baseTableColour   = tcell.ColorMediumSpringGreen
	tableQuery        = "SELECT * FROM %s"
	tablesQuery       = "SELECT DISTINCT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'"
)

func GetTable(query string) *tview.Table {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	rowData, headers, err := db.ExecuteQuery(fmt.Sprintf(tableQuery, query))

	if err != nil || headers == nil {
		panic(err)
	}

	for headerIndex, header := range headers {
		table.SetCell(0, headerIndex,
			tview.NewTableCell(strings.ToUpper(header)).
				SetTextColor(tableHeaderColour).
				SetAlign(tview.AlignLeft))
	}

	for rowIndex, row := range rowData {
		for itemIndex, item := range row {
			table.SetCell(rowIndex+1, itemIndex,
				tview.NewTableCell(item).
					SetTextColor(baseTableColour).
					SetAlign(tview.AlignLeft))
		}
	}

	table.
		SetTitle(fmt.Sprintf(" %s [%d] ", query, len(rowData))).
		SetTitleColor(baseTableColour).
		SetBorder(true).
		SetBorderStyle(tcell.StyleDefault).
		SetBorderColor(baseTableColour).
		SetBorderPadding(0, 0, 1, 1)

	return table
}

func GetTables() *tview.List {
	tableList := tview.NewList()
	tableList.SetBorder(true).
		SetTitleColor(baseTableColour).
		SetBorderColor(baseTableColour).
		SetBorderPadding(0, 0, 1, 1)

	rows, headers, err := db.ExecuteQuery(tablesQuery)

	if headers == nil {
		panic(err)
	}

	tableList.AddItem(strings.ToUpper(headers[0]), "", 0, nil).
		SetMainTextColor(tableHeaderColour)

	for _, tableName := range rows {
		tableList.AddItem(tableName[0], "", 0, nil).
			SetMainTextColor(baseTableColour)
	}

	tableList.
		ShowSecondaryText(false).
		SetTitle(fmt.Sprintf(" Tables [%d] ", len(rows)))

	return tableList
}
