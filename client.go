package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func connectToServer() (DatabaseList, error) {
	//resp, err := http.Get(fmt.Sprintf("http://%s:%s/databases", IP, PORT))
	resp, err := http.Get("http://localhost:1323/databases")
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
