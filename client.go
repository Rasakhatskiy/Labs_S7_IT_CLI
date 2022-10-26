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

func connectToServer() (DatabaseList, error) {
	url := getUrl()
	url += "/databases"
	resp, err := http.Get(url)
	if err != nil {
		return DatabaseList{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DatabaseList{}, err
	}
	//Convert the body to type string
	sb := string(body)

	var list DatabaseList

	err = json.Unmarshal([]byte(sb), &list)
	if err != nil {
		return DatabaseList{}, err
	}

	return list, nil
}

func getTablesList(dbname string) ([]string, error) {
	url := getUrl() + "/databases/" + dbname
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sb := string(body)

	var list []string

	err = json.Unmarshal([]byte(sb), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
