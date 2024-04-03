package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	hook "github.com/robotn/gohook"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed schema.json
var jsonSchema string

type Shortcut struct {
	Name    string
	Keys    []string
	Command string
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
		fmt.Printf("Register shortcut %v\n", shortcut)
		hook.Register(hook.KeyDown, shortcut.Keys, func(e hook.Event) {
			fmt.Printf("Shortcut <%s> activated\n", shortcut.Name)
			command := strings.Split(shortcut.Command, " ")
			if len(command) == 1 {
				cmd := exec.Command(command[0])
				go cmd.Run()
			} else {
				cmd := exec.Command(command[0], command[1:]...)
				go cmd.Run()
			}
		})
	}
}

func parseShrotcuts() ([]Shortcut, error) {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return nil, err
	}

	// Construct the relative file path
	configPath := filepath.Join(filepath.Dir(exePath), "shortcut.conf.json")
	jsonConfig, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	schemaLoader := gojsonschema.NewBytesLoader([]byte(jsonSchema))
	documentLoader := gojsonschema.NewBytesLoader([]byte(jsonConfig))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
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

	err = json.Unmarshal([]byte(jsonConfig), &data)
	if err != nil {
		return nil, err
	}

	return data.Shortcuts, nil
}
