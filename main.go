package main

import (
	"github.com/rivo/tview"
)

var app *tview.Application

func main() {
	app = tview.NewApplication()

	setConnectForm()
	app.Stop()
}
