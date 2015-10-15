package main

import (
	"github.com/go-martini/martini"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type inBody struct {
	In	string	`json:"in"`
}

func main() {
	var inputCollection []string

	m := martini.Classic()
	m.Get("/", func() string {
		collection := "Collected inputs:\n"
		for i, input := range inputCollection {
			collection += fmt.Sprintf("index: %v \t|\tinput was: {%v}\n", i, input)
		}
		return collection
	})
	m.Post("/in", func(res http.ResponseWriter, req *http.Request) string {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "something went wrong"
		}
		var in inBody
		err = json.Unmarshal(body, &in)
		if err != nil {
			return "something else went wrong"
		}
		inputCollection = append(inputCollection, in.In)
		return fmt.Sprintf("here's what you gave me: " + in.In + "\ncurrent collection size: %v", len(inputCollection))
	})
	m.Run()
}