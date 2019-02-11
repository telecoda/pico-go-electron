// +build windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

const (
	gopherJS = "gopherjs.exe"
)

func getBuildCmd(cartFile, mainFile, spritesFile, outFile string) *exec.Cmd {
	// we use GOOS=linux to compile to JS even on windows...
	cmd := exec.Command(getGopherJSPath(), "build", cartFile, mainFile, spritesFile, "-o", outFile)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Env = append(os.Environ(), "GOOS=linux")
	return cmd
}
