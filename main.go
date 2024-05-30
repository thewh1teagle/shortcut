package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/emersion/go-autostart"
	"github.com/fsnotify/fsnotify"
	hook "github.com/robotn/gohook"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed shortcut.schema.json
var jsonSchema string

type Shortcut struct {
	Name       string
	Keys       []string
	Command    string
	HideWindow *bool // Pointer to bool for optional field
}

var installFlag bool
var versionFlag bool
var version = "1.0.0"

func reloadApp() error {
	// Get the executable path
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}
	// Windows being special *again*
	if runtime.GOOS == "windows" {
		// Start a new process on Windows
		cmd := exec.Command(exe, os.Args...)
		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Start()
		if err == nil {
			os.Exit(0)
		} else {
			log.Fatal(err.Error())
		}
	} else {
		// Replace the current process with a new one
		err = syscall.Exec(exe, os.Args, os.Environ())
		if err != nil {
			return fmt.Errorf("executing new process: %w", err)
		}
	}

	// Exit the current process
	os.Exit(0)
	return nil
}

func hotReload(path string) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					// log.Println("modified file:", event.Name)
					reloadApp()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan struct{})
}

func main() {
	if runtime.GOOS == "windows" {
		attachConsoleIfPossible()
	}
	cobra.MousetrapHelpText = "" // allow running by clicking .exe file from GUI
	var rootCmd = &cobra.Command{
		Use:   "shortcut",
		Short: "Shortcuts manager",
		Run:   run,
	}

	rootCmd.Flags().BoolVarP(&installFlag, "install", "i", false, "Install the application")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Get version")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	if installFlag {
		install()
		os.Exit(0)
	}
	if versionFlag {
		fmt.Printf("shortcut version %v\n", version)
		os.Exit(0)
	}
	configPath, err := getConfigPath()
	if err != nil {
		fmt.Println("Not config found:", err)
		os.Exit(1)
	}
	jsonConfig, err := readConfig(*configPath)
	go hotReload(*configPath)
	if err != nil {
		// Handle error
		fmt.Println("Failed to read config:", err)
	}

	shortcuts, err := parseShrotcuts(jsonConfig)
	if err != nil {
		// Block main goroutine forever.
		<-make(chan struct{})
	} else {
		fmt.Println(shortcuts)
		registerShortcuts(shortcuts)

		s := hook.Start()
		<-hook.Process(s)
	}
}

func install() error {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return err
	}
	app := &autostart.App{
		Name:        "shortcuts",
		DisplayName: "Shortcuts Manager",
		Exec:        []string{exePath},
	}

	if app.IsEnabled() {
		log.Println("Shortcut already installed, removing it...")

		if err := app.Disable(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("Installing shortcut to be run at boot from %v...\n", exePath)

		if err := app.Enable(); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func registerShortcuts(shortcuts []Shortcut) {
	for _, shortcut := range shortcuts {
		fmt.Printf("Register shortcut %v\n", shortcut)
		hook.Register(hook.KeyDown, shortcut.Keys, func(e hook.Event) {
			fmt.Printf("Shortcut <%s> activated\n", shortcut.Name)
			command := strings.Split(shortcut.Command, " ")
			cmd := exec.Command(command[0], command[1:]...)
			if runtime.GOOS == "windows" && shortcut.HideWindow != nil {
				cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: *shortcut.HideWindow}
			}

			if err := cmd.Start(); err != nil {
				fmt.Printf("Error executing command for shortcut <%s>: %v\n", shortcut.Name, err)
			}
		})
	}
}

func getConfigPath() (*string, error) {
	// Try getting the executable path
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return nil, err
	}

	// Construct the relative file path based on executable directory
	configPath := filepath.Join(filepath.Dir(exePath), "shortcut.conf.json")

	// Check if the config file exists at the executable directory
	if _, err := os.Stat(configPath); err == nil {
		return &configPath, nil
	}

	// If the config file is not found in the executable directory, try the cwd
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return nil, err
	}

	configPath = filepath.Join(cwd, "shortcut.conf.json")

	// Check if the config file exists at the cwd
	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("config file not found")
	}

	return &configPath, nil
}

func readConfig(path string) ([]byte, error) {
	jsonConfig, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	return jsonConfig, nil
}

func parseShrotcuts(jsonConfig []byte) ([]Shortcut, error) {

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
