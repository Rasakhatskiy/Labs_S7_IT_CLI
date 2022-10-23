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

func showDBList(list DatabaseList) {
	var listForm *tview.List
	listForm = tview.NewList()
	for i, s := range list.Databases {
		listForm.AddItem(s, "", rune('a'+i), nil)
	}
	listForm.AddItem("Back", "Press to exit", 'q', func() {
		setConnectForm()
	})

	if err := app.SetRoot(listForm, true).SetFocus(listForm).Run(); err != nil {
		panic(err)
	}
}

func setConnectForm() {
	var form *tview.Form

	form = tview.NewForm().
		AddInputField("IP", "localhost", 16, nil, func(text string) { IP = text }).
		AddInputField("Port", "1323", 10, nil, func(text string) { PORT = text }).
		AddButton("Connect", func() {
			dblist, err := connectToServer()
			if err != nil {
				showErrorMessage(form, nil, "OK", "Exit", err.Error())
			} else {
				showDBList(dblist)
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Connect to server").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
}
