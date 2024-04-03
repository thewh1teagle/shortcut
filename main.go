package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	hook "github.com/robotn/gohook"
	"github.com/xeipuuv/gojsonschema"
)

const configPath = "shortcuts.json"
const jsonSchemaPath = "schema.json"

type Shortcut struct {
	Name      string
	Shortcuts []string
	Action    string
}

func main() {

	shortcuts, err := parseShrotcuts()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(shortcuts)
	registerShortcuts(shortcuts)

	s := hook.Start()
	<-hook.Process(s)
}

func registerShortcuts(shortcuts []Shortcut) {
	for _, shortcut := range shortcuts {
		hook.Register(hook.KeyDown, shortcut.Shortcuts, func(e hook.Event) {
			fmt.Printf("Shortcut <%s> activated", shortcut.Name)
			command := strings.Split(shortcut.Action, " ")
			if len(command) == 1 {
				exec.Command(command[0])
			} else {
				exec.Command(command[0], command[1:]...)
			}
		})
	}
}

func parseShrotcuts() ([]Shortcut, error) {
	// Load the JSON schema
	schemaLoader := gojsonschema.NewStringLoader(jsonSchemaPath)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return nil, err
	}

	// Load the JSON data
	documentLoader := gojsonschema.NewStringLoader(configPath)

	// Validate the JSON data against the schema
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return nil, err
	}

	if !result.Valid() {
		fmt.Printf("The JSON data is invalid: %s\n", result.Errors())
		return nil, fmt.Errorf("invalid JSON data")
	}

	// Unmarshal the validated JSON data
	var data struct {
		Version   string     `json:"version"`
		Shortcuts []Shortcut `json:"shortcuts"`
	}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}

	return data.Shortcuts, nil
}
