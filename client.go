package main

import (
	"io/ioutil"
	"net/http"
)

func connectToServer() (string, error) {
	//resp, err := http.Get(fmt.Sprintf("http://%s:%s/databases", IP, PORT))
	resp, err := http.Get("http://localhost:1323/databases")
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//Convert the body to type string
	sb := string(body)
	return sb, nil
}
