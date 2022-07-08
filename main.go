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
	Id    string
	Email string
	Age   string
}

var jsonUsers []user
var jsonUser user

func Perform(args Arguments, writer io.Writer) error {
	// Recieved parsed commands and write the related data
	// to console

	// open file
	fileJson, error := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0644)
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

		if len(tempBuffer) != 0 {
			if err := json.Unmarshal(tempBuffer, &jsonUsers); err != nil {
				return fmt.Errorf("%w", error)
			}
		}
		if args["item"] == "" {
			return fmt.Errorf("-item flag should be specified")
		}
		err = json.Unmarshal([]byte(args["item"]), &jsonUser)
		if err != nil {
			return fmt.Errorf("-json object is not valid %w", error)
		}

		if jsonUser.Id == "" || jsonUser.Email == "" || jsonUser.Age == "" {
			return fmt.Errorf("-data is not full or incorrect")
		}
		for _, value := range jsonUsers {
			if value.Id == jsonUser.Id {
				return fmt.Errorf("Such Id already exists")
			}
		}

		jsonUsers = append(jsonUsers, jsonUser)

		fmt.Printf("%+v", jsonUsers)
		data, error := json.Marshal(jsonUsers)
		if error != nil {
			return fmt.Errorf("%w", error)
		}

		fmt.Println(data)
		if error := ioutil.WriteFile("users.json", data, 0644); error != nil {
			fileJson.Close()
			return fmt.Errorf("%w", error)
		}
	}

	//case "list":

	//	return 2
	//case "findById":
	//	return 3
	//case "remove":
	//	return 4
	//default:
	//	return 5
	//}
	//return nil */
	//}
	//}
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

	flag.Parse()

	output = make(map[string]string, 3)
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
