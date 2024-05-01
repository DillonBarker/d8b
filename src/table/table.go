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
	tableQuery        = "SELECT * FROM %s.%s"
	schemasQuery      = "SELECT schema_name FROM information_schema.schemata"
	tablesQuery       = "SELECT table_name FROM information_schema.tables WHERE table_schema = '%s';"
)

func GetSchemas() *tview.List {
	schemaList := tview.NewList()
	schemaList.SetBorder(true).
		SetTitleColor(baseTableColour).
		SetBorderColor(baseTableColour).
		SetBorderPadding(0, 0, 1, 1)

	rows, headers, err := db.ExecuteQuery(schemasQuery)

	if headers == nil {
		panic(err)
	}

	schemaList.AddItem(strings.ToUpper(headers[0]), "", 0, nil).
		SetMainTextColor(tableHeaderColour)

	for _, tableName := range rows {
		schemaList.AddItem(tableName[0], "", 0, nil).
			SetMainTextColor(baseTableColour)
	}

	schemaList.
		ShowSecondaryText(false).
		SetTitle(fmt.Sprintf(" schemas [%d] ", len(rows)))

	return schemaList
}

func GetTable(schemaName string, tableName string) *tview.Table {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	rowData, headers, err := db.ExecuteQuery(fmt.Sprintf(tableQuery, schemaName, tableName))

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
		SetTitle(fmt.Sprintf(" %s.%s [%d] ", schemaName, tableName, len(rowData))).
		SetTitleColor(baseTableColour).
		SetBorder(true).
		SetBorderStyle(tcell.StyleDefault).
		SetBorderColor(baseTableColour).
		SetBorderPadding(0, 0, 1, 1)

	return table
}

func GetTables(schemaName string) *tview.List {
	tableList := tview.NewList()
	tableList.SetBorder(true).
		SetTitleColor(baseTableColour).
		SetBorderColor(baseTableColour).
		SetBorderPadding(0, 0, 1, 1)

	rows, headers, err := db.ExecuteQuery(fmt.Sprintf(tablesQuery, schemaName))

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
		SetTitle(fmt.Sprintf(" %s [%d] ", schemaName, len(rows)))

	return tableList
}
