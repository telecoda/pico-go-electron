// +build linux

package main

import (
	"os"
	"os/exec"
	"syscall"
)

const (
	gopherJS = "gopherjs"
)

func getBuildCmd(cartFile, mainFile, outFile string) *exec.Cmd {
	cmd := exec.Command(getGopherJSPath(), "build", cartFile, mainFile, "-o", outFile)
	cmd.Env = append(os.Environ(), "GOOS=linux")
	return cmd
}