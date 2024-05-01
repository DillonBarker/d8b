package ui

import (
	"fmt"

	"github.com/DillonBarker/d8b/internal/db"
	"github.com/rivo/tview"
)

func Header(config db.Config) *tview.Flex {
	labels := tview.NewTextView().SetTextColor(tview.Styles.SecondaryTextColor)
	values := tview.NewTextView().SetTextColor(tview.Styles.PrimaryTextColor)

	labels.SetText("Host:\nPort:\nDatabase:\nUser:").
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	values.SetText(fmt.Sprintf("%s\n%d\n%s\n%s", config.Host, config.Port, config.DBName, config.User)).
		SetTextAlign(tview.AlignLeft)

	header := tview.NewFlex().
		AddItem(labels, 0, 1, false).
		AddItem(values, 0, 10, false)

	return header
}
