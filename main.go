package main

import (
	"github.com/rivo/tview"
)

var app *tview.Application

func main() {
	app = tview.NewApplication()

	var form *tview.Form

	form = tview.NewForm().
		AddInputField("IP", "localhost", 16, nil, func(text string) { IP = text }).
		AddInputField("Port", "1323", 10, nil, func(text string) { PORT = text }).
		AddButton("Connect", func() {
			body, err := connectToServer()
			if err != nil {
				showErrorMessage(form, nil, "OK", "Exit", err.Error())
			} else {
				showInfoMessage(form, body)
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
