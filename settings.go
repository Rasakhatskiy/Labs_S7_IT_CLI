package main

var IP string
var PORT string

type DatabaseList struct {
	Databases []string `json:"databases"`
}

type TableJSON struct {
	Headers []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"headers"`
	Values [][]interface{} `json:"values"`
}
