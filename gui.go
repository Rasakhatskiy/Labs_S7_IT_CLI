package main

import (
	"CLI_DBMS_viewer/database"
	"CLI_DBMS_viewer/globvar"
	"CLI_DBMS_viewer/utils"
	"fmt"
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

func showInfoMessage(optionForm *tview.Flex, message string) {
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

	app.SetRoot(modal, false)
}

func showYesNoMessage(message string, fYes, fNo func()) {
	var modal *tview.Modal
	modal = tview.NewModal().
		SetText(message).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				fYes()
			}
			if buttonLabel == "No" {
				fNo()
			}
		})
	app.SetRoot(modal, false)
}

func showDBList(list DatabaseList) {
	var listForm *tview.List
	listForm = tview.NewList()

	for i, s := range list.Databases {
		listForm.AddItem(s, "", rune('a'+i), nil)
	}

	listForm.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		globvar.DBname = s
		showTablesList()
	})

	listForm.AddItem("Back", "Press to exit", 'q', func() {
		showConnectForm()
	})

	listForm.SetBorder(true).SetTitle(" Available databases ").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(listForm, true).SetFocus(listForm).Run(); err != nil {
		panic(err)
	}
}

func showTablesList() {
	list, err := getTablesList(globvar.DBname)
	if err != nil {
		showErrorMessage(nil, nil, "Exit", "Cooler Exit", err.Error())
	}

	var listForm *tview.List
	listForm = tview.NewList()
	for i, s := range list {
		listForm.AddItem(s, "", rune('a'+i), func() {
			table, err := getTable(globvar.DBname, s)
			if err != nil {
				showErrorMessage(nil, nil, "Exit", "Cooler Exit", err.Error())
			}
			showTableForm("", "", table)
		})
	}
	listForm.AddItem("Add new table", "", 'n', func() {
		showCreateTableForm()
	})
	listForm.AddItem("Back", "", 'q', func() {
		showConnectForm()
	})

	listForm.SetBorder(true).SetTitle(fmt.Sprintf(" Database %s ", globvar.DBname)).SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(listForm, true).SetFocus(listForm).Run(); err != nil {
		panic(err)
	}
}

func showTableForm(dbName, tableName string, tablet *TableJSON) {
	table := tview.NewTable().
		SetBorders(true)

	for i, header := range tablet.Headers {
		table.SetCell(0, i, tview.NewTableCell(header.Name).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter))

		table.SetCell(1, i, tview.NewTableCell(header.Type).
			SetTextColor(tcell.ColorBlue).
			SetAlign(tview.AlignCenter))
	}

	for i, row := range tablet.Values {
		for j, data := range row {
			table.SetCell(i+2, j, tview.NewTableCell(fmt.Sprintf("%v", data)).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter))
		}
	}

	table.Select(0, 0).
		SetSelectable(true, false).
		SetFixed(1, 0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				showTablesList()
			}
		}).
		SetSelectedFunc(func(row int, column int) {
			//table.GetCell(row, column).SetTextColor(tcell.ColorRed)
			//table.SetSelectable(false, false)
		})

	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}

func showConnectForm() {
	var form *tview.Form

	ipField := tview.NewInputField()
	ipField.SetLabel("IP")
	ipField.SetText("localhost")
	ipField.SetFieldWidth(16)

	portField := tview.NewInputField()
	portField.SetLabel("Port")
	portField.SetText("1323")
	portField.SetFieldWidth(8)

	form = tview.NewForm().
		AddFormItem(ipField).
		AddFormItem(portField).
		AddButton("Connect", func() {
			IP = ipField.GetText()
			PORT = portField.GetText()
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

func showCreateTableForm() {
	var rowFlex *tview.Flex
	var mainFlex *tview.Flex
	var mainForm *tview.Form
	var columnForm *tview.Form
	var columnList *tview.List

	columnList = tview.NewList()

	columnList.SetBorder(true)
	columnList.SetTitle("Columns")
	columnList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		globvar.Headers = utils.RemoveIndex(globvar.Headers, i)
		globvar.Types = utils.RemoveIndex(globvar.Types, i)
		setList(columnList)
		app.SetFocus(columnForm)
	})

	rowFlex = tview.NewFlex().SetDirection(tview.FlexRow)

	textBoxColName := tview.NewInputField()
	textBoxColName.SetTitle("Column name")
	textBoxColName.SetFieldWidth(32)

	dropDown := tview.NewDropDown()
	dropDown.SetTitle("Column type")
	dropDown.SetOptions(database.TypesListStr, nil)

	columnForm = tview.NewForm().
		AddFormItem(textBoxColName).
		AddFormItem(dropDown).
		AddButton("Add column", func() {
			colName := textBoxColName.GetText()
			_, colType := dropDown.GetCurrentOption()
			globvar.Headers = append(globvar.Headers, colName)
			globvar.Types = append(globvar.Types, colType)
			setList(columnList)
		}).
		AddButton("<- Back", func() {
			app.SetFocus(mainForm)
		}).
		AddButton("Remove column ->", func() {
			app.SetFocus(columnList)
		})
	columnForm.SetBorder(true)

	mainForm = tview.NewForm().
		AddInputField("Table name", "", 32, nil, nil).
		AddButton("Save", func() {
			err := sendCreatedDB()

			if err == nil {
				//showInfoMessage(mainFlex, "Table created successfully")
				showTablesList()
			} else {
				showInfoMessage(mainFlex, "хуйня: "+err.Error())
			}

		}).
		AddButton("Cancel", func() {
		}).
		AddButton("to create columns ->", func() {
			app.SetFocus(columnForm)
		})
	mainForm.SetBorder(true)

	rowFlex.
		AddItem(mainForm, 0, 1, true).
		AddItem(columnForm, 0, 1, true)

	mainFlex = tview.NewFlex().
		AddItem(rowFlex, 0, 2, false).
		AddItem(columnList, 0, 2, false)

	if err := app.SetRoot(mainFlex, true).SetFocus(mainForm).Run(); err != nil {
		panic(err)
	}
}

func setList(columnList *tview.List) {
	columnList.Clear()
	for i, name := range globvar.Headers {
		columnList.AddItem(name, globvar.Types[i], rune('0'+i), nil)
	}
}
