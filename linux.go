// +build linux

package main

import (
	"os"
	"os/exec"
)

const (
	gopherJS = "gopherjs"
)

func getBuildCmd(cartFile, mainFile, spritesFile, outFile string) *exec.Cmd {
	cmd := exec.Command(getGopherJSPath(), "build", cartFile, mainFile, spritesFile, "-o", outFile)
	cmd.Env = append(os.Environ(), "GOOS=linux")
	return cmd
}
