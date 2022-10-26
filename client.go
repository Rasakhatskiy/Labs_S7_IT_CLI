package main

import (
	"encoding/json"
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
