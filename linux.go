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

func getBuildCmd(sourceFile, outFile string) *exec.Cmd {
	cmd := exec.Command(gopherJS, "build", sourceFile, "-o", outFile)
	cmd.Env = append(os.Environ(), "GOOS=linux")
	return cmd
}
