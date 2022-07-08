package main

import (
	"encoding/json"

	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string
type user struct {
	Id    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age,omitempty"`
}

var jsonUsers []user
var jsonUser user

func Perform(args Arguments, writer io.Writer) error {
	// Recieved parsed commands and write the related data
	// to console
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	// common error check
	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}

	// open file
	fileToReadWrite := args["fileName"]
	fileJson, error := os.OpenFile(fileToReadWrite, os.O_RDWR|os.O_CREATE, 0644)
	if error != nil {
		return fmt.Errorf("%w", error)
	}
	defer fileJson.Close()

	switch args["operation"] {
	case "add":

		// reading the content of file and put an empty slice to it
		tempBuffer, err := ioutil.ReadAll(fileJson)
		if err != nil {
			return fmt.Errorf("%w", error)
		}
		jsonUsers = make([]user, 0)
		if len(tempBuffer) != 0 {
			json.Unmarshal(tempBuffer, &jsonUsers)
		}
		if args["item"] == "" {
			return fmt.Errorf("-item flag has to be specified")
		}
		err = json.Unmarshal([]byte(args["item"]), &jsonUser)
		if err != nil {
			return fmt.Errorf("-json object is not valid %w", error)
		}

		if jsonUser.Id == "" || jsonUser.Email == "" || jsonUser.Age == 0 {
			return fmt.Errorf("-data is not full or incorrect")
		}
		for _, value := range jsonUsers {
			if value.Id == jsonUser.Id {
				message := "Item with id " + jsonUser.Id + " already exists"
				writer.Write([]byte(message))
				return nil
			}
		}

		jsonUsers = append(jsonUsers, jsonUser)

		data, error := json.Marshal(jsonUsers)
		if error != nil {
			return fmt.Errorf("%w", error)
		}

		if error := ioutil.WriteFile(fileToReadWrite, data, 0644); error != nil {
			fileJson.Close()
			return fmt.Errorf("%w", error)
		}
	case "list":
		contentOfFile, error := ioutil.ReadAll(fileJson)
		if error != nil {
			return fmt.Errorf("%w", error)
		}
		writer.Write(contentOfFile)
	case "findById":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
		contentOfFile, error := ioutil.ReadAll(fileJson)
		if error != nil {
			return fmt.Errorf("Cannot read the file - %w", error)
		}
		if error := json.Unmarshal(contentOfFile, &jsonUsers); error != nil {
			return fmt.Errorf("%w", error)
		}
		for _, value := range jsonUsers {
			if value.Id == args["id"] {
				buff, _ := json.Marshal(value)
				writer.Write(buff)
				return nil
			}
		}
		writer.Write([]byte(""))
	case "remove":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
		contentOfFile, error := ioutil.ReadAll(fileJson)
		if error != nil {
			return fmt.Errorf("Cannot read the file - %w", error)
		}
		if error := json.Unmarshal(contentOfFile, &jsonUsers); error != nil {
			return fmt.Errorf("%w", error)
		}
		for index, value := range jsonUsers {
			if value.Id == args["id"] {
				tempslice := removeElementOfSlice(jsonUsers, index)
				byteBuff, _ := json.Marshal(tempslice)
				ioutil.WriteFile(args["fileName"], byteBuff, 0644)
				return nil
			}
		}
		message := "Item with id " + args["id"] + " not found"
		writer.Write([]byte(message))
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
	return nil
}
func parseArgs() (output Arguments) {
	// Function which makes parsing arguments
	// recieved from command line to Arguments type
	var operation string
	flag.StringVar(&operation, "operation", "", "add, list, findById, remove")
	var item string
	flag.StringVar(&item, "item", "", "data add to JSON")
	var fileName string
	flag.StringVar(&fileName, "fileName", "", "path to the JSON file")
	var id string
	flag.StringVar(&id, "id", "", "id argument for using with FindById and Remove command")
	flag.Parse()

	output = make(map[string]string, 4)
	output["id"] = id
	output["operation"] = operation
	output["item"] = item
	output["fileName"] = fileName
	return
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}

}

func removeElementOfSlice(s []user, i int) []user {
	ret := make([]user, 0)
	ret = append(ret, s[:i]...)
	return append(ret, s[i+1:]...)
}
