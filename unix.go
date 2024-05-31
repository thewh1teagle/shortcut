//go:build !windows
// +build !windows

package main

import "os/exec"

func attachConsoleIfPossible()                        {}
func prepareCommand(cmd *exec.Cmd, shortcut Shortcut) {}
