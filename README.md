# shortcut

<img src="https://github.com/thewh1teagle/shortcut/blob/main/design/logo.png?raw=true" width=180 >

# Introduction

Shortcut is a command-line tool designed to simplify the creation of global keyboard shortcuts across various platforms. By defining keyboard shortcuts in a straightforward `JSON` file, users can effortlessly assign commands to specific key combinations, like `Ctrl+C`, linking them to versatile shell commands that can execute tasks ranging from opening files to launching URLs. With the added benefit of `JSON` schema-based auto-completion, Shortcut streamlines the process of customizing keyboard functionalities to enhance user efficiency.

# Downloads

Download one of the releases in [shortcut/releases](https://github.com/thewh1teagle/shortcuts)

# Usage

Put `shortcut` in known path such as `C:\bin`.
create JSON file alongside `shortcut` named `shortcut.config.json` with the following:
```json
{
  "version": "0.0.1",
  "shortcuts": [
    {"name" "test", action: "echo hi", "keys": ["ctrl", "d"]}
  ]
}
```


# Supported Platforms

`Windows`, `Linux`, `macOS`
