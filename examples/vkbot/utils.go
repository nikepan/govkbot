package main

import (
    "os"
    "encoding/json"
    "log"
    "io/ioutil"
)

func readJSON(fn string, v interface{}) {
    file, _ := os.Open(fn)
    defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(v)
	if err != nil {
		log.Println("error:", err)
	}
}

func readText(fn string) string {
    content, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println("error:", err)
	}
	return string(content)
}