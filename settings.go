package main

var IP string
var PORT string

type DBPathJSON struct {
	Name string `json:"name"`
}

type TableHeaderJSON struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TableJSON struct {
	Name    string            `json:"name"`
	Headers []TableHeaderJSON `json:"headers"`
	Values  [][]interface{}   `json:"values"`
}
