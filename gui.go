package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func showErrorMessage(firstOptionForm, secondOptionForm *tview.Form, btnText1, btnText2, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{btnText1, btnText2}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == btnText1 {
				if firstOptionForm != nil {
					app.SetRoot(firstOptionForm, true)
				} else {
					app.Stop()
				}
			}
			if buttonLabel == btnText2 {
				if secondOptionForm != nil {
					app.SetRoot(secondOptionForm, true)
				} else {
					app.Stop()
				}
			}
		}).
		SetBackgroundColor(tcell.ColorRed)

	app.SetRoot(modal, true)
}

func showInfoMessage(optionForm *tview.Form, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				if optionForm != nil {
					app.SetRoot(optionForm, true)
				} else {
					app.Stop()
				}
			}

		})

	app.SetRoot(modal, true)
}
