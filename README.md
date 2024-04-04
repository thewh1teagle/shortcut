# shortcut

<img src="https://github.com/thewh1teagle/shortcut/blob/main/design/logo.png?raw=true" width=180 >

# Introduction

Shortcut is a command-line tool that lets users create global keyboard shortcuts using a `JSON` file. It allows assigning key combinations to shell commands for tasks like opening files or launching URLs. The tool offers JSON schema-based auto-completion to simplify the customization process and boost user productivity.

# Downloads

Download one of the releases in [shortcut/releases](https://github.com/thewh1teagle/shortcut/releases)

# Usage

1. Put `shortcut` in known path such as `C:\bin`.
2. create JSON file alongside `shortcut` named `shortcut.config.json` with the following:
```json
{
  "version": "0.0.1",
  "shortcuts": [
    {"name" "test", "command": "echo hi", "keys": ["ctrl", "d"]}
  ]
}
```
3. Run it once. it will autostart on boot from then.

Now you can customize it, I recommend edit using [VSCode](https://code.visualstudio.com/download)

# Start on boot
For starting shortcut at boot, execute
```console
./shortcut --install
```

# Supported Platforms

`Windows`, `Linux`, `macOS`

# Todo
- [x] Start on boot with [go-autostart](https://github.com/emersion/go-autostart)
- [x] Hot reload with [fsnotify](https://github.com/fsnotify/fsnotify)
- [x] Releases with [goReleaser](https://goreleaser.com/quick-start/)
- [x] Keys autocomplete
- [ ] Installer like [oranda](https://github.com/axodotdev/oranda)
