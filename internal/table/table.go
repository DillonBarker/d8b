package table

import (
	"fmt"
	"strings"

	"github.com/DillonBarker/d8b/internal/db"
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

func GetSchemas() ([]string, []string, error) {
	rows, headers, err := db.ExecuteQuery(schemasQuery)

	if headers == nil {
		panic(err)
	}

	var flattened []string
	for _, row := range rows {
		flattened = append(flattened, row...)
	}

	return flattened, headers, err
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
				SetTextColor(tview.Styles.PrimaryTextColor).
				SetAlign(tview.AlignLeft))
	}

	for rowIndex, row := range rowData {
		for itemIndex, item := range row {
			table.SetCell(rowIndex+1, itemIndex,
				tview.NewTableCell(item).
					SetTextColor(tview.Styles.SecondaryTextColor).
					SetAlign(tview.AlignLeft))
		}
	}

	table.
		SetTitle(fmt.Sprintf(" table(%s.%s) [%d] ", schemaName, tableName, len(rowData))).
		SetTitleColor(tview.Styles.PrimaryTextColor).
		SetBorder(true).
		SetBorderColor(baseTableColour).
		SetBorderPadding(0, 0, 1, 1)

	return table
}

func GetTables(schemaName string) ([]string, error) {
	rows, headers, err := db.ExecuteQuery(fmt.Sprintf(tablesQuery, schemaName))

	if headers == nil {
		panic(err)
	}

	var flattened []string
	for _, row := range rows {
		flattened = append(flattened, row...)
	}

	return flattened, err
}
