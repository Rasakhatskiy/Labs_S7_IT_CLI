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

func makeHttpRequest(requestType, url string, data []byte) (string, error) {
	var resp *http.Response
	var err error

	switch requestType {
	case globvar.REQ_GET:
		resp, err = http.Get(url)
		defer resp.Body.Close()
	case globvar.REQ_POST:
		resp, err = http.Post(url, "application/json", bytes.NewBuffer(data))
		defer resp.Body.Close()
	case globvar.REQ_DELETE:
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return "", err
		}

		return "", errors.New("not implemented yet")
		//todo delete
	}
	if err != nil {
		return "", err
	}

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", errors.New(string(respBody))
	}

	return string(respBody), nil
}

func connectToServer() ([]DBPathJSON, error) {
	url := fmt.Sprintf("%s/databases", getUrl())
	sb, err := makeHttpRequest(globvar.REQ_GET, url, nil)
	if err != nil {
		return []DBPathJSON{}, err
	}

	var list []DBPathJSON
	err = json.Unmarshal([]byte(sb), &list)
	if err != nil {
		return []DBPathJSON{}, err
	}

	return list, nil
}

func getTablesList(dbname string) ([]string, error) {
	url := fmt.Sprintf("%s/databases/%s", getUrl(), dbname)
	sb, err := makeHttpRequest(globvar.REQ_GET, url, nil)
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
	sb, err := makeHttpRequest(globvar.REQ_GET, url, nil)
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

func postTablePrep() error {
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
		Name:    globvar.TableName,
		Headers: headersJSON,
		Values:  nil,
	}

	err := postTable(globvar.DBname, &tableJSON)
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

	_, err = makeHttpRequest(globvar.REQ_POST, url, data)
	if err != nil {
		return err
	}

	return nil
}

func postCreateDB(name string) error {
	url := fmt.Sprintf("%s/databases/new_database", getUrl())
	data, err := json.Marshal(name)
	if err != nil {
		return err
	}

	_, err = makeHttpRequest(globvar.REQ_POST, url, data)
	if err != nil {
		return err
	}

	return nil
}

func postNewRow(dbname, tableName string, values []string) error {
	url := fmt.Sprintf("%s/databases/%s/%s/new_row", getUrl(), dbname, tableName)
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}

	_, err = makeHttpRequest(globvar.REQ_POST, url, data)

	if err != nil {
		return err
	}

	return nil
}

func postEditRow(dbname, tableName string, values []string, index int) error {
	url := fmt.Sprintf("%s/databases/%s/%s/%d", getUrl(), dbname, tableName, index)
	data, err := json.Marshal(values)

	if err != nil {
		return err
	}

	_, err = makeHttpRequest(globvar.REQ_POST, url, data)

	if err != nil {
		return err
	}

	return nil
}

func deleteRow(dbname, tableName string, index int) error {
	url := fmt.Sprintf("%s/databases/%s/%s/%d", getUrl(), dbname, tableName, index)
	_, err := makeHttpRequest(globvar.REQ_DELETE, url, nil)

	if err != nil {
		return err
	}

	return nil
}
