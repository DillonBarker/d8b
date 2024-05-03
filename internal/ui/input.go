package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Input() *tview.InputField {
	inputField := tview.NewInputField()

	inputField.SetBorder(true)

	inputField.SetFieldTextColor(tview.Styles.PrimaryTextColor)
	inputField.SetPlaceholderTextColor(tview.Styles.PrimaryTextColor)
	inputField.SetFieldBackgroundColor(tcell.ColorBlack)

	return inputField
}
