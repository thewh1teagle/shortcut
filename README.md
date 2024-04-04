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


# Supported Platforms

`Windows`, `Linux`, `macOS`

# Todo
- [ ] Start on boot with [go-autostart](https://github.com/emersion/go-autostart)
- [ ] Installer like [oranda](https://github.com/axodotdev/oranda)
- [ ] Hot reload with [fsnotify](https://github.com/fsnotify/fsnotify)
- [ ] Releases with [goReleaser](https://goreleaser.com/quick-start/)
- [ ] Keys autocomplete
