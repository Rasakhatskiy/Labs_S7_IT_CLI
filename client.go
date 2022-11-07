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

func connectToServer() ([]DBPathJSON, error) {
	sb, err := getHttpResponse(fmt.Sprintf("%s/databases", getUrl()))
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

func postCreateDB(name string) error {
	url := fmt.Sprintf("%s/databases/new_database", getUrl())
	data, err := json.Marshal(name)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(data))

	//err = json.NewDecoder(resp.Body).Decode(&resp)
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Println(resp)
	return nil
}

func makeHttpRequest(requestType, url string, data []byte) (string, error) {
	var resp *http.Response
	var err error

	switch requestType {
	case globvar.REQ_GET:
		resp, err = http.Get(url)
		break
	case globvar.REQ_POST:
		resp, err = http.Post(url, "application/json", bytes.NewBuffer(data))
		break
	case globvar.REQ_DELETE:
		//todo delete
		break
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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

func postNewRow(dbname, tableName string, values []string) error {
	url := fmt.Sprintf("%s/databases/%s/%s/new_row", getUrl(), dbname, tableName)
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}

	err = makeHttpRequest(url, data)

	// Display Results
	//fmt.Println("response Status : ", resp.Status)
	//fmt.Println("response Headers : ", resp.Header)
	//fmt.Println("response Body : ", string(respBody))

	return nil
}

func postEditRow(dbname, tableName string, values []string) error {

	return nil
}

func deleteRow(dbname, tableName string, data []interface{}) error {
	return nil
}
