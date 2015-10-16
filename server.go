package main

import (
	"github.com/go-martini/martini"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
)

type inBody struct {
	In	string	`json:"in"`
}

type Docker struct {
	Name string `json:"Name"`
	Image string `json:"Image"`
	Command string `json:"Command"`
}

type Tarball struct {
	Data []byte `json:"Data"`
	Container Docker `json:"Container"`
}

func main() {
	var inputCollection []string

	var uploads map[string]

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
	m.Post("/upload_container", func(res http.ResponseWriter, req *http.Request) string {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "something went wrong"
		}
		var tarball Tarball
		err = json.Unmarshal(body, &tarball)
		if err != nil {
			return "data did not read as Container/Data pair" + string(body)
		}
//		path := "./uploads/"+tarball.Container.Name+".tgz"
//		if _, err := os.Stat(path); !os.IsNotExist(err) {
//			return "this container has already been uploaded"
//		}
//		err = writeFileToDisk(path, tarball.Data)
//		if err != nil {
//			return "error writing file to disk" + err.Error()
//		}
		if _, haskey := uploads[tarball.Container.Name]; haskey {
			return "error: "+tarball.Container.Name+" has already been uploaded."
		}
		uploads[tarball.Container.Name] = tarball
		return tarball.Container.Name + " saved successfully"
	})

	m.Post("/download_container/:container_name", func(params martini.Params) []byte {
		containerName := params["container_name"]
		if _, haskey := uploads[containerName]; haskey {
			return "error: "+containerName+" does not exist."
		}
		tarball := uploads[containerName]
		response, err := json.Marshal(tarball)
		if err != nil {
			return "error converting tarball into a response"
		}
		delete(uploads, containerName)
		return response
	})



	m.Run()
}

//func writeFileToDisk(path string, data []byte) error {
//	err := ioutil.WriteFile(path, data, 0666)
//	return err
//}
//
//func readFileFromDisk(path string) ([]byte, error) {
//	data, err := ioutil.ReadFile(path)
//	return data, err
//}