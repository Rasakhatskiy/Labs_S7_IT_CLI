package main

import (
	"CLI_DBMS_viewer/database"
	"CLI_DBMS_viewer/globvar"
	"CLI_DBMS_viewer/utils"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"reflect"
)

func showErrorMessage(firstOptionForm, secondOptionForm *tview.Flex, btnText1, btnText2, message string) {
	flex := tview.NewFlex()
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{btnText1, btnText2}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == btnText1 {
				if firstOptionForm != nil {
					focusOnFlex(firstOptionForm)
				} else {
					app.Stop()
				}
			}
			if buttonLabel == btnText2 {
				if secondOptionForm != nil {
					focusOnFlex(secondOptionForm)
				} else {
					app.Stop()
				}
			}
		}).
		SetBackgroundColor(tcell.ColorRed)

	flex.AddItem(modal, 0, 1, true)

	focusOnFlex(flex)
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

func focusOnFlex(flex *tview.Flex) {
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

func prepareDBList() {
	dblist, err := connectToServer()
	if err != nil {
		app.Stop()
		showErrorMessage(getConnectForm(), nil, "OK", "Exit", err.Error())
	} else {
		focusOnFlex(getDBListForm(dblist))
	}
}

func getConnectForm() *tview.Flex {
	flex := tview.NewFlex()
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
			prepareDBList()
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Connect to server").SetTitleAlign(tview.AlignLeft)

	flex.AddItem(form, 0, 1, true)

	return flex
}

func getDBListForm(list []DBPathJSON) *tview.Flex {
	flex := tview.NewFlex()
	var listForm *tview.List
	listForm = tview.NewList()

	for i, s := range list {
		listForm.AddItem(s.Name, "", rune('0'+i), nil)
	}

	listForm.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		globvar.DBname = s
		focusOnFlex(getTablesListForm())
	})

	listForm.AddItem("Create database", "new file", 'n', func() {
		focusOnFlex(getCreateDBForm())
	})

	listForm.AddItem("Back", "Press to exit", 'q', func() {
		focusOnFlex(getConnectForm())
	})
	listForm.SetBorder(true).SetTitle(" Available databases ").SetTitleAlign(tview.AlignLeft)

	flex.AddItem(listForm, 0, 1, true)
	return flex
}

func getTablesListForm() *tview.Flex {
	flex := tview.NewFlex()

	list, err := getTablesList(globvar.DBname)
	if err != nil {
		showErrorMessage(nil, nil, "Exit", "Cooler Exit", err.Error())
		return nil
	}

	var listForm *tview.List
	listForm = tview.NewList()
	listForm.SetSelectedFunc(func(i int, tableName string, s2 string, r rune) {
		tableJSON, err := getTable(globvar.DBname, tableName)
		if err != nil {
			showErrorMessage(flex, nil, "Ok", ":(", err.Error())
		}
		globvar.TableName = tableName
		focusOnFlex(getTableForm(tableJSON))
	})

	for i, s := range list {
		listForm.AddItem(s, "", rune('a'+i), nil)
	}
	listForm.AddItem("Add new table", "", 'n', func() {
		focusOnFlex(getCreateTableForm())
	})
	listForm.AddItem("Back", "", 'q', func() {
		prepareDBList()
	})

	listForm.SetBorder(true).SetTitle(fmt.Sprintf(" Database %s ", globvar.DBname)).SetTitleAlign(tview.AlignLeft)

	flex.AddItem(listForm, 0, 1, true)
	return flex
}

func setTableValues(table *tview.Table, tableJSON *TableJSON) {
	table.Clear()
	for i, header := range tableJSON.Headers {
		table.SetCell(0, i, tview.NewTableCell(header.Name).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter))

		table.SetCell(1, i, tview.NewTableCell(header.Type).
			SetTextColor(tcell.ColorBlue).
			SetAlign(tview.AlignCenter))
	}

	for i, row := range tableJSON.Values {
		for j, data := range row {
			table.SetCell(i+2, j, tview.NewTableCell(fmt.Sprintf("%v", data)).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter))
		}
	}
}

func getTableForm(tableJSON *TableJSON) *tview.Flex {
	var (
		table       = tview.NewTable()
		listOptions = tview.NewList()
		addEditForm = tview.NewForm()
		flex        = tview.NewFlex()
		inputs      []*tview.InputField
	)

	addEditForm.SetBorder(true)
	addEditForm.AddButton("Save", func() {
		var values []string
		for _, input := range inputs {
			if len(input.GetText()) == 0 {
				app.SetFocus(input)
				return
			}
			values = append(values, input.GetText())
		}

		var err error = nil
		if globvar.TableOperationType == globvar.Create {
			err = postNewRow(globvar.DBname, globvar.TableName, values)
		}
		if globvar.TableOperationType == globvar.Update {
			err = postEditRow(globvar.DBname, globvar.TableName, values, globvar.SelectedRow)
		}
		if err != nil {
			showErrorMessage(flex, nil, "OK", "Exit", err.Error())
		} else {
			tableJSON, err = getTable(globvar.DBname, globvar.TableName)
			if err != nil {
				showErrorMessage(flex, nil, "OK", "Exit", err.Error())
			}
			focusOnFlex(getTableForm(tableJSON))
		}
	})
	addEditForm.AddButton("Cancel", func() {
		app.SetFocus(listOptions)
	})

	// set input fields for edit
	for _, header := range tableJSON.Headers {
		input := tview.NewInputField().
			SetLabel(header.Name)

		inputs = append(inputs, input)
		if header.Type == database.TypeStringRangeTS {
			input2 := tview.NewInputField().
				SetLabel(header.Name + " 2")
			inputs = append(inputs, input2)
		}
	}

	listOptions.
		AddItem("Navigate", "", 'n', func() {
			globvar.TableOperationType = globvar.Read
			app.SetFocus(table)
		}).
		AddItem("Add row", "", 'a', func() {
			globvar.TableOperationType = globvar.Create
			addEditForm.Clear(false)
			for i := range inputs {
				addEditForm.AddFormItem(inputs[i])
				app.SetFocus(addEditForm)
			}
		}).
		AddItem("Edit row", "", 'e', func() {
			globvar.TableOperationType = globvar.Update
			app.SetFocus(table)
		}).
		AddItem("Delete row", "", 'd', func() {
			globvar.TableOperationType = globvar.Delete
			app.SetFocus(table)
		}).
		AddItem("Back", "", 'b', func() {
			focusOnFlex(getTablesListForm())
		})

	listOptions.SetBorder(true)

	setTableValues(table, tableJSON)

	table.Select(0, 0).
		SetSelectable(true, false).
		SetFixed(1, 0).

		// ESCAPE
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				app.SetFocus(listOptions)
			}
		}).

		// ENTER
		SetSelectedFunc(func(row int, column int) {
			if row < 2 {
				return
			}
			switch globvar.TableOperationType {
			case globvar.Read:
				break
			case globvar.Update:
				offset := 0
				for i, data := range tableJSON.Values[row-2] {
					if tableJSON.Headers[i].Type == database.TypeStringRangeTS {
						list := reflect.ValueOf(data)
						inputs[i+offset+0].SetText(fmt.Sprintf("%v", list.Index(0)))
						inputs[i+offset+1].SetText(fmt.Sprintf("%v", list.Index(1)))
						offset++
					} else {
						inputs[i+offset].SetText(fmt.Sprintf("%v", data))
					}
				}
				globvar.SelectedRow = row - 2
				addEditForm.Clear(false)
				for i := range inputs {
					addEditForm.AddFormItem(inputs[i])
					app.SetFocus(addEditForm)
				}
				app.SetFocus(addEditForm)
				break
			case globvar.Delete:
				globvar.SelectedRow = row - 2
				err := deleteRow(globvar.DBname, globvar.TableName, row)
				if err != nil {
					showErrorMessage(flex, nil, "OK", "Exit", err.Error())
				} else {
					showInfoMessage(getTableForm(tableJSON), "Successfully deleted")
				}
				break
			}

			//table.GetCell(row, column).SetTextColor(tcell.ColorRed)
			//table.SetSelectable(false, false)
		})
	table.SetBorders(true)

	rowFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rowFlex.AddItem(listOptions, 0, 1, true)
	rowFlex.AddItem(addEditForm, 0, 3, true)

	flex.AddItem(rowFlex, 0, 1, true)
	flex.AddItem(table, 0, 3, true)
	return flex
}

func setList(columnList *tview.List) {
	columnList.Clear()
	for i, name := range globvar.Headers {
		columnList.AddItem(name, globvar.Types[i], rune('0'+i), nil)
	}
}

func getCreateTableForm() *tview.Flex {
	var rowFlex *tview.Flex
	var flex *tview.Flex
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
	textBoxColName.SetLabel("Column name")
	textBoxColName.SetFieldWidth(32)

	dropDown := tview.NewDropDown()
	dropDown.SetLabel("Column type")
	dropDown.SetOptions(database.TypesListStr, nil)

	columnForm = tview.NewForm().
		AddFormItem(textBoxColName).
		AddFormItem(dropDown).
		AddButton("Add column", func() {
			colName := textBoxColName.GetText()
			_, colType := dropDown.GetCurrentOption()

			globvar.Headers = append(globvar.Headers, colName)
			globvar.Types = append(globvar.Types, colType)

			textBoxColName.SetText("")
			dropDown.SetOptions(database.TypesListStr, nil)

			setList(columnList)
		}).
		AddButton("<- Back", func() {
			app.SetFocus(mainForm)
		}).
		AddButton("Remove column ->", func() {
			app.SetFocus(columnList)
		})
	columnForm.SetBorder(true)

	tableNameInput := tview.NewInputField().
		SetLabel("Table name").
		SetFieldWidth(32)

	mainForm = tview.NewForm().
		AddFormItem(tableNameInput).
		AddButton("to create columns ->", func() {
			app.SetFocus(columnForm)
		}).
		AddButton("Save", func() {
			globvar.TableName = tableNameInput.GetText()

			err := postTablePrep()

			if err == nil {
				showInfoMessage(getTablesListForm(), "Table created successfully")
			} else {
				showInfoMessage(flex, "error: "+err.Error())
			}

		}).
		AddButton("Cancel", func() {
			focusOnFlex(getTablesListForm())
		})
	mainForm.SetBorder(true)

	rowFlex.
		AddItem(mainForm, 0, 1, true).
		AddItem(columnForm, 0, 1, false)

	flex = tview.NewFlex().
		AddItem(rowFlex, 0, 2, true).
		AddItem(columnList, 0, 2, false)

	return flex
}

func getCreateDBForm() *tview.Flex {
	flex := tview.NewFlex()

	nameInput := tview.NewInputField()
	nameInput.SetLabel("DB name")
	nameInput.SetFieldWidth(16)

	form := tview.NewForm().
		AddFormItem(nameInput).
		AddButton("OK", func() {
			err := postCreateDB(nameInput.GetText())
			if err != nil {
				showErrorMessage(flex, flex, "OK", "OK", err.Error())
			}
			prepareDBList()
		}).
		AddButton("Cancel", func() {
			prepareDBList()
		})

	form.
		SetBorder(true).
		SetTitle("Create DB")

	flex.AddItem(form, 0, 1, true)

	return flex
}
