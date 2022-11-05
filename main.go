package main

import (
	"CLI_DBMS_viewer/database"
	"github.com/rivo/tview"
)

var app *tview.Application

func initTypes() {
	database.TypesListStr = []string{
		database.TypeIntegerTS,
		database.TypeRealTS,
		database.TypeCharTS,
		database.TypeStringTS,
		database.TypeHTMLTS,
		database.TypeStringRangeTS,
	}
}

func main() {
	initTypes()
	app = tview.NewApplication()

	showConnectForm()
	app.Stop()
}
