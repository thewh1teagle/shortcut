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
	"time"

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

func handleConfigError(err error) {
	fmt.Println("Failed to read config:", err)
	os.Exit(1)
}

func setupWatcher(configPath string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watcher.Add(configPath)
	if err != nil {
		return nil, err
	}
	return watcher, nil
}

func blockForever() {
	// Block main goroutine forever.
	<-make(chan struct{})
}

func run(cmd *cobra.Command, args []string) {
	if installFlag {
		install()
		os.Exit(0)
	}

	fmt.Printf("🚀 Shortcut version %v\n", version)
	if versionFlag {
		os.Exit(0)
	}

	configPath, err := getConfigPath()
	if err != nil {
		handleConfigError(err)
	}

	watcher, err := setupWatcher(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	jsonConfig, err := readConfig(*configPath)
	if err != nil {
		handleConfigError(err)
	}

	shortcuts, err := parseShortcuts(jsonConfig)
	if err != nil {
		fmt.Println("🚨 Unable to parse shortcuts, blocking forever...")
		blockForever()
	} else {
		registerAndWatch(shortcuts, watcher)
	}
}

func registerAndWatch(shortcuts []Shortcut, watcher *fsnotify.Watcher) {
	registerShortcuts(shortcuts)
	s := hook.Start()
	go func() {
		<-hook.Process(s)
	}()

	var debounceTimer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(50*time.Millisecond, func() {
				handleFileEvent(event, shortcuts)
			})
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func handleFileEvent(event fsnotify.Event, shortcuts []Shortcut) {
	if event.Has(fsnotify.Write) {
		fmt.Println("🔥 Hot Reload triggered")
		hook.End()
		registerShortcuts(shortcuts)
		s := hook.Start()
		go func() {
			<-hook.Process(s)
		}()
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
		hook.Register(hook.KeyDown, shortcut.Keys, func(e hook.Event) {
			fmt.Printf("Shortcut '%s' activated ✅\n", shortcut.Name)

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
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return nil, err
	}

	configPath := filepath.Join(filepath.Dir(exePath), "shortcut.conf.json")

	if _, err := os.Stat(configPath); err == nil {
		return &configPath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return nil, err
	}

	configPath = filepath.Join(cwd, "shortcut.conf.json")

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

func parseShortcuts(jsonConfig []byte) ([]Shortcut, error) {
	schemaLoader := gojsonschema.NewBytesLoader([]byte(jsonSchema))
	documentLoader := gojsonschema.NewBytesLoader([]byte(jsonConfig))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, err
	}

	if !result.Valid() {
		fmt.Printf("🚨 The JSON data is invalid: %s\n", result.Errors())
		return nil, fmt.Errorf("invalid JSON data")
	}

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
