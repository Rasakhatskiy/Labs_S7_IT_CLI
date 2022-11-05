package main

import (
	"CLI_DBMS_viewer/globvar"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//const url = "http://localhost:1323"

func getUrl() string {
	return fmt.Sprintf("http://%s:%s", IP, PORT)
}

func getHttpResponse(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//Convert the body to type string

	return string(body), nil
}

func connectToServer() (DatabaseList, error) {
	sb, err := getHttpResponse(fmt.Sprintf("%s/databases", getUrl()))
	if err != nil {
		return DatabaseList{}, err
	}

	var list DatabaseList

	err = json.Unmarshal([]byte(sb), &list)
	if err != nil {
		return DatabaseList{}, err
	}

	return list, nil
}

func getTablesList(dbname string) ([]string, error) {
	sb, err := getHttpResponse(getUrl() + "/databases/" + dbname)
	if err != nil {
		return nil, err
	}

	var list []string

	err = json.Unmarshal([]byte(sb), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func getTable(dbName, tableName string) (*TableJSON, error) {
	url := fmt.Sprintf("%s/databases/%s/%s", getUrl(), dbName, tableName)
	sb, err := getHttpResponse(url)
	if err != nil {
		return nil, err
	}

	var table TableJSON
	err = json.Unmarshal([]byte(sb), &table)
	if err != nil {
		return nil, err
	}

	return &table, nil
}

func sendCreatedDB() error {
	if len(globvar.Headers) != len(globvar.Types) {
		return errors.New("fatal: len of headers != len of types")
	}

	if len(globvar.Headers) == 0 {
		return errors.New("empty table")
	}

	var headersJSON []TableHeaderJSON
	for i, name := range globvar.Headers {
		headersJSON = append(headersJSON, TableHeaderJSON{
			Name: name,
			Type: globvar.Types[i],
		})
	}

	tableJSON := TableJSON{
		Headers: headersJSON,
		Values:  nil,
	}

	err := postTable("", &tableJSON)
	if err != nil {
		return err
	}

	return nil
}

func postTable(dbName string, tableJSON *TableJSON) error {
	url := fmt.Sprintf("%s/databases/%s/new_table", getUrl(), dbName)
	data, err := json.Marshal(tableJSON)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}

	fmt.Println(res["json"])
	return nil
}
