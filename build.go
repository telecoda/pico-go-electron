package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

func getGoPath() string {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		// default to user home + /go
		user, err := user.Current()
		if err != nil {
			panic(fmt.Errorf("Failed to fetch current user: %s", err))
		}
		goPath = filepath.Join(user.HomeDir, "go")
	}
	return goPath
}

func getGopherJSPath() string {
	return filepath.Join(getGoPath(), "bin", gopherJS)
}

func getVersionCmd() *exec.Cmd {
	cmd := exec.Command(getGopherJSPath(), "version")
	return cmd
}
