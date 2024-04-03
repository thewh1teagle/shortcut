# shortcut

<img src="https://github.com/thewh1teagle/shortcut/blob/main/design/logo.png?raw=true" width=180 >

# Introduction

Shortcut is a command-line tool designed to simplify the creation of global keyboard shortcuts across various platforms. By defining keyboard shortcuts in a straightforward `JSON` file, users can effortlessly assign commands to specific key combinations, like `Ctrl+C`, linking them to versatile shell commands that can execute tasks ranging from opening files to launching URLs. With the added benefit of `JSON` schema-based auto-completion, Shortcut streamlines the process of customizing keyboard functionalities to enhance user efficiency.

# Downloads

Download one of the releases in [shortcut/releases](https://github.com/thewh1teagle/shortcuts)

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