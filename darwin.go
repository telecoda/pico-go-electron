// +build darwin

package main

import (
	"os"
	"os/exec"
)

const (
	gopherJS = "gopherjs"
)

func getBuildCmd(sourceFile, outFile string) *exec.Cmd {
	cmd := exec.Command(getGopherJSPath(), "build", sourceFile, "-o", outFile)
	cmd.Env = append(os.Environ(), "GOOS=linux")
	return cmd
}
